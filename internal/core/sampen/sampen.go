package sampen

import (
	"fmt"
	"math"

	"github.com/kio235/sampEn_SCFD_backend/internal/utils/mathmethod"
)

type SampEnCalc struct {
	m           int
	rCoeff      float64
	voltageData []float64
	n           int
	sampEn      float64
	computed    bool
}

func NewSampEnCalc(m int, rCoeff float64) *SampEnCalc {
	return &SampEnCalc{m: m, rCoeff: rCoeff}
}

func (c *SampEnCalc) LoadData(voltageData []float64) error {
	n := len(voltageData)
	if n < c.m+1 {
		return fmt.Errorf("data length must be higher than m+1=%d", c.m+1)
	}
	c.n = n
	c.voltageData = make([]float64, n)
	copy(c.voltageData, voltageData)
	c.computed = false
	return nil
}

func (c *SampEnCalc) Compute() (sampEn float64, err error) {
	if c.n < c.m+1 {
		return 0, fmt.Errorf("no data, use LoadData() first")
	}

	r := c.rCoeff * mathmethod.StandardDeviation(c.voltageData)
	if r <= 0 {
		return 0, fmt.Errorf("r<=0, data not valid")
	}

	a := 0
	b := 0

	windowCount := c.n - c.m
	for i := range windowCount {
		for j := range windowCount {
			if i == j {
				continue
			}
			diff, err := mathmethod.ChebyshevDistance(c.voltageData[i:i+c.m], c.voltageData[j:j+c.m])
			if err != nil {
				return 0, err
			}
			if diff <= r {
				b += 1
			}
			diff, err = mathmethod.ChebyshevDistance(c.voltageData[i:i+c.m+1], c.voltageData[j:j+c.m+1])
			if err != nil {
				return 0, err
			}
			if diff <= r {
				a += 1
			}
		}
	}
	if b == 0 {
		return 0, fmt.Errorf("B=0, can't process valid sampEn")
	}
	c.sampEn = -math.Log(float64(a) / float64(b))
	return c.sampEn, nil
}
