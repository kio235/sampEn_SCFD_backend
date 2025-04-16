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

type Detector struct {
	m                 int
	rCoeff            float64
	windowSize        int
	stepSize          int
	relativeThreshold float64
	absoluteThreshold float64
	faultThreshold    int
	records           batteryRecord
	HasRecords        bool
	calc              *sampen.SampEnCalc
	SampEnMatrix      [][]float64
	FaultInfos        []FaultInfo
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

// SetArgs 调整 Detector 参数
//
// Parameters:
//
//	m - 参数 m
//	rCoeff - r 系数
//	wd - 窗口大小 (window size)
//	st - 步长 (step)
//	rel - 相对差异阈值 (默认 20%)
//	abs - 绝对差异阈值 (默认 0.2)
// func (d *Detector) SetArgs(m int, rCoeff float64, wd int, st int, rel float64, abs float64, cnt int) error {
// 	calc, err := sampen.NewSampEnCalc(m, rCoeff)
// 	if err != nil {
// 		return err
// 	}
// 	if wd <= 0 {
// 		return fmt.Errorf("window size must > 0")
// 	}
// 	if st <= 0 {
// 		return fmt.Errorf("step must > 0")
// 	}
// 	if wd < st {
// 		return fmt.Errorf("windows size must > step")
// 	}
// 	if rel < 0 {
// 		return fmt.Errorf("rel threshold must >= 0")
// 	}
// 	if abs < 0 {
// 		return fmt.Errorf("abs threshold must >= 0")
// 	}
// 	if cnt < 1 {
// 		return fmt.Errorf("cnt must >=1")
// 	}
// 	d.calc = calc
// 	d.windowSize = wd
// 	d.step = st
// 	d.relThreshold = rel
// 	d.absThreshold = abs
// 	d.FaultCount = cnt
// 	return nil
// }

func (d *Detector) LoadRecords(cel []string, vol []float64) error {
	if len(cel) < 3 {
		return fmt.Errorf("Must have more than 2 cells.")
	}
	if len(vol) == 0 {
		return fmt.Errorf("Record is empty.")
	}
	if !mathmethod.IsMultiple(len(vol), len(cel)) {
		return fmt.Errorf("Length of record is not an integer multiple of cell length.")
	}

	d.records.voltage = make([]float64, len(vol))
	copy(d.records.voltage, vol)
	d.records.cell = make([]string, len(cel))
	copy(d.records.cell, cel)

	d.records.recordCount = len(vol) / len(cel)
	d.records.cellCount = len(cel)
	d.HasRecords = true
	return nil
}

// func (d *Detector) Compute() error {
// 	vol := mathmethod.ConvertTo2DShared(d.records.voltage, d.records.cellCount, d.records.recordCount)
// 	sampEnResult := make([][]float64, 0)
// 	for cel := range d.records.cellCount {
// 		sampEnResult = append(sampEnResult, make([]float64, 0))
// 		for i := 0; i <= d.records.recordCount-d.windowSize; i += d.step {
// 			err := d.calc.LoadData(vol[cel][i : i+d.windowSize])
// 			if err != nil {
// 				return err
// 			}
// 			en, err := d.calc.Compute()
// 			if err != nil {
// 				return err
// 			}
// 			sampEnResult[cel] = append(sampEnResult[cel], en)
// 		}
// 	}
// 	d.SampEnMatrix = sampEnResult
// 	return nil
// }

func (d *Detector) Compute() error {
	vol := mathmethod.ConvertTo2DShared(d.records.voltage, d.records.cellCount, d.records.recordCount)
	sampEnResult := make([][]float64, d.records.cellCount)

	g, ctx := errgroup.WithContext(context.Background())

	for cel := 0; cel < d.records.cellCount; cel++ {
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
