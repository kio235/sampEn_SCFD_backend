package core

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/kio235/sampEn_SCFD_backend/internal/core/sampen"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/mathmethod"
	"golang.org/x/sync/errgroup"
)

// Detector 使用样品熵(package sampen)检测一个电压序列是否存在故障。
type Detector struct {
	m                 int                // 样品熵参数 m
	rCoeff            float64            // 样品熵参数 r coefficient
	windowSize        int                // 滑动窗口大小
	stepSize          int                // 滑动窗口步长
	relativeThreshold float64            // 检测目标样品熵与其他样品熵的相对差异阈值
	absoluteThreshold float64            // 检测目标样品熵与其他样品熵的相对差异阈值
	faultThreshold    int                // 出现连续faultThreshold个异常点则视为错误
	records           batteryRecord      // 电池数据
	HasRecords        bool               // 是否已加载电池数据
	calc              *sampen.SampEnCalc // 样品熵检测结构体
	SampEnMatrix      [][]float64        // 样品熵结果
	FaultInfos        []FaultInfo        // 错误检测结果
}

// NewDetector 创建一个新 Detector
//
//	m - 参数 m (默认2)
//	rCoeff - r 系数 (默认0.2)
//	windowSize - 窗口大小 (window size) (默认100)
//	stepSize - 步长 (step) (默认1)
//	relativeThreshold - 相对差异阈值 (默认 30%)
//	absoluteThreshold - 绝对差异阈值 (默认 0.4)
//	faultThreshold - 出现连续faultThreshold个异常点则视为错误 (默认5)
func NewDetector(opts ...DetectorOption) (*Detector, error) {
	d := &Detector{
		m:                 2,
		rCoeff:            0.2,
		windowSize:        100,
		stepSize:          1,
		relativeThreshold: 0.3,
		absoluteThreshold: 0.4,
		faultThreshold:    5,
		SampEnMatrix:      make([][]float64, 0),
		FaultInfos:        make([]FaultInfo, 0),
	}
	for _, opt := range opts {
		opt(d)
	}

	if d.m < 1 {
		return nil, fmt.Errorf("m must >=1")
	}
	if d.rCoeff <= 0 {
		return nil, fmt.Errorf("r must >0")
	}
	if d.windowSize <= 0 {
		return nil, fmt.Errorf("window size must > 0")
	}
	if d.stepSize <= 0 {
		return nil, fmt.Errorf("step must > 0")
	}
	if d.relativeThreshold < 0 {
		return nil, fmt.Errorf("relative threshold must >= 0")
	}
	if d.absoluteThreshold < 0 {
		return nil, fmt.Errorf("absolute threshold must >= 0")
	}
	if d.faultThreshold < 1 {
		return nil, fmt.Errorf("fault threshold must >=1")
	}
	var err error
	d.calc, err = sampen.NewSampEnCalc(d.m, d.rCoeff)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// LoadRecords 加载电池数据
//
//	cellsName - 每个电池编号
//	voltages - 电池电压的一维数组
func (d *Detector) LoadRecords(cellsName []string, voltages []float64) error {
	if len(cellsName) < 3 {
		return fmt.Errorf("Must have more than 2 cells.")
	}
	if len(voltages) == 0 {
		return fmt.Errorf("Record is empty.")
	}
	if !mathmethod.IsMultiple(len(voltages), len(cellsName)) {
		return fmt.Errorf("Length of record is not an integer multiple of cell length.")
	}

	d.records.voltage = make([]float64, len(voltages))
	copy(d.records.voltage, voltages)
	d.records.cell = make([]string, len(cellsName))
	copy(d.records.cell, cellsName)

	d.records.recordCount = len(voltages) / len(cellsName)
	d.records.cellCount = len(cellsName)
	d.HasRecords = true
	return nil
}

// ComputeSingle 使用滑动窗口计算电压序列的样品熵, 没有进行多核优化
func (d *Detector) ComputeSingle() error {
	vol := mathmethod.ConvertTo2DShared(d.records.voltage, d.records.cellCount, d.records.recordCount)
	sampEnResult := make([][]float64, 0)
	for cel := range d.records.cellCount {
		sampEnResult = append(sampEnResult, make([]float64, 0))
		for i := 0; i <= d.records.recordCount-d.windowSize; i += d.stepSize {
			err := d.calc.LoadData(vol[cel][i : i+d.windowSize])
			if err != nil {
				return err
			}
			en, err := d.calc.Compute()
			if err != nil {
				return err
			}
			sampEnResult[cel] = append(sampEnResult[cel], en)
		}
	}
	d.SampEnMatrix = sampEnResult
	return nil
}

// ComputeSingle 使用滑动窗口计算电压序列的样品熵, 使用 goroutine 进行多核优化
func (d *Detector) Compute() error {
	vol := mathmethod.ConvertTo2DShared(d.records.voltage, d.records.cellCount, d.records.recordCount)
	sampEnResult := make([][]float64, d.records.cellCount)

	g, ctx := errgroup.WithContext(context.Background())

	for cel := range d.records.cellCount {
		cel := cel
		g.Go(func() error {
			cellVol := vol[cel]
			recordCount := len(cellVol)
			step := d.stepSize
			windowSize := d.windowSize

			// 预计算结果切片大小
			numWindows := (recordCount - windowSize + step) / step
			results := make([]float64, numWindows)

			// 错误处理通道（缓冲大小为1）
			errCh := make(chan error, 1)

			// 任务通道
			taskCh := make(chan struct {
				index int
				start int
			}, numWindows)

			// 生成任务
			go func() {
				defer close(taskCh)
				for idx, i := 0, 0; i <= recordCount-windowSize; i, idx = i+step, idx+1 {
					select {
					case <-ctx.Done():
						return
					case taskCh <- struct {
						index int
						start int
					}{index: idx, start: i}:
					}
				}
			}()

			// 启动Worker池
			var wg sync.WaitGroup
			for w := 0; w < runtime.NumCPU(); w++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for task := range taskCh {
						select {
						case <-ctx.Done():
							return
						default:
							calc := d.calc.Clone()
							window := cellVol[task.start : task.start+windowSize]
							if err := calc.LoadData(window); err != nil {
								select {
								case errCh <- fmt.Errorf("cell %d window %d: %w", cel, task.start, err):
								default: // 保证不阻塞
								}
								return
							}
							en, err := calc.Compute()
							if err != nil {
								select {
								case errCh <- fmt.Errorf("cell %d window %d: %w", cel, task.start, err):
								default:
								}
								return
							}
							results[task.index] = en
						}
					}
				}()
			}

			// 等待所有worker完成
			wg.Wait()
			close(errCh)

			// 检查错误
			if err := <-errCh; err != nil {
				return err
			}

			// 保存结果
			sampEnResult[cel] = results
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	d.SampEnMatrix = sampEnResult
	return nil
}

func (d *Detector) Check() error {
	if !d.HasRecords {
		return fmt.Errorf("No records.")
	}
	err := d.Compute()
	if err != nil {
		return err
	}

	for cel := range d.SampEnMatrix {
		count := 0 // 连续出现 count 个异常点
		for t := range d.SampEnMatrix[0] {
			var other []float64
			for c := range d.SampEnMatrix {
				if c == cel {
					continue
				}
				other = append(other, d.SampEnMatrix[c][t])
			}
			mean := mathmethod.Mean(other)
			if math.Abs(d.SampEnMatrix[cel][t]-mean) >= d.absoluteThreshold {
				if math.Abs(d.SampEnMatrix[cel][t]-mean)/mean >= d.relativeThreshold {
					count++
					if count >= d.faultThreshold {
						d.FaultInfos = append(d.FaultInfos, FaultInfo{t*d.stepSize + d.windowSize/2, cel})
					}
					continue
				}
			}
			continue
		}
	}
	return nil
}
