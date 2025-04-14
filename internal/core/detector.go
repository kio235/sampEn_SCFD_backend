package core

import (
	"fmt"
	"math"

	"github.com/kio235/sampEn_SCFD_backend/internal/core/sampen"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/mathmethod"
)

type Detector struct {
	windowSize   int
	step         int
	relThreshold float64
	absThreshold float64
	records      batteryRecord
	HasRecords   bool
	calc         *sampen.SampEnCalc
	SampEnMatrix [][]float64
	FaultInfos   []FaultInfo
}

// NewDetector 创建一个新 Detector
//
// Parameters:
//
//	m - 参数 m
//	rCoeff - r 系数
//	wd - 窗口大小 (window size)
//	st - 步长 (step)
//	rel - 相对差异阈值 (默认 20%)
//	abs - 绝对差异阈值 (默认 0.2)
func NewDetector(m int, rCoeff float64, wd int, st int, rel float64, abs float64) (*Detector, error) {
	var d Detector
	var err error
	d.calc, err = sampen.NewSampEnCalc(m, rCoeff)
	if err != nil {
		return nil, err
	}
	if wd <= 0 {
		return nil, fmt.Errorf("window size must > 0")
	}
	if st <= 0 {
		return nil, fmt.Errorf("step must > 0")
	}
	if wd < st {
		return nil, fmt.Errorf("windows size must > step")
	}
	if rel < 0 {
		return nil, fmt.Errorf("rel threshold must >= 0")
	}
	if abs < 0 {
		return nil, fmt.Errorf("abs threshold must >= 0")
	}
	d.windowSize = wd
	d.step = st
	d.relThreshold = rel
	d.absThreshold = abs
	return &d, nil
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
func (d *Detector) SetArgs(m int, rCoeff float64, wd int, st int, rel float64, abs float64) error {
	calc, err := sampen.NewSampEnCalc(m, rCoeff)
	if err != nil {
		return err
	}
	if wd <= 0 {
		return fmt.Errorf("window size must > 0")
	}
	if st <= 0 {
		return fmt.Errorf("step must > 0")
	}
	if wd < st {
		return fmt.Errorf("windows size must > step")
	}
	if rel < 0 {
		return fmt.Errorf("rel threshold must >= 0")
	}
	if abs < 0 {
		return fmt.Errorf("abs threshold must >= 0")
	}
	d.calc = calc
	d.windowSize = wd
	d.step = st
	d.relThreshold = rel
	d.absThreshold = abs
	return nil
}

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

func (d *Detector) Compute() error {
	vol := mathmethod.ConvertTo2DShared(d.records.voltage, d.records.cellCount, d.records.recordCount)
	sampEnResult := make([][]float64, 0)
	for cel := range d.records.cellCount {
		sampEnResult = append(sampEnResult, make([]float64, 0))
		for i := 0; i <= d.records.recordCount-d.windowSize; i += d.step {
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

func (d *Detector) Check() (faultCount int, err error) {
	if !d.HasRecords {
		return 0, fmt.Errorf("No records.")
	}
	err = d.Compute()
	if err != nil {
		return 0, err
	}

	for t := range d.SampEnMatrix[0] {
		for cel := range d.SampEnMatrix {
			var other []float64
			for c := range d.SampEnMatrix {
				if c == cel {
					continue
				}
				other = append(other, d.SampEnMatrix[c][t])
			}
			mean := mathmethod.Mean(other)
			if math.Abs(d.SampEnMatrix[cel][t]-mean) >= d.absThreshold {
				if math.Abs(d.SampEnMatrix[cel][t]-mean)/mean >= d.relThreshold {
					faultCount += 1
					d.FaultInfos = append(d.FaultInfos, FaultInfo{t, cel})
				}
			}
		}
	}
	return
}
