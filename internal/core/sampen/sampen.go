// Package sampen provides functionality for calculating Sample Entropy (SampEn) of time series data.
// SampEn is a measure of complexity that quantifies the unpredictability of fluctuations in a time series.
package sampen

import (
	"fmt"
	"math"

	"github.com/kio235/sampEn_SCFD_backend/internal/utils/mathmethod"
)

// SampEnCalc represents a Sample Entropy calculator.
// It holds the parameters and data needed for SampEn calculation.
type SampEnCalc struct {
	m           int       // embedding dimension
	rCoeff      float64   // tolerance coefficient
	voltageData []float64 // time series data
	n           int       // length of data
	sampEn      float64   // calculated sample entropy value
	computed    bool      // flag indicating if SampEn has been computed
}

// Clone creates a deep copy of the SampEnCalc instance.
// This is useful when you need to preserve the original calculator state.
func (c *SampEnCalc) Clone() *SampEnCalc {
	vol := make([]float64, len(c.voltageData))
	copy(vol, c.voltageData)
	return &SampEnCalc{c.m, c.rCoeff, vol, c.n, c.sampEn, c.computed}
}

// NewSampEnCalc creates a new Sample Entropy calculator with specified parameters.
// Parameters:
//   - m: embedding dimension, must be >= 1
//   - rCoeff: tolerance coefficient, must be > 0
//
// Returns a new SampEnCalc instance or an error if parameters are invalid.
func NewSampEnCalc(m int, rCoeff float64) (*SampEnCalc, error) {
	if m < 1 {
		return nil, fmt.Errorf("m must >=1")
	}
	if rCoeff <= 0 {
		return nil, fmt.Errorf("r coefficient must >0")
	}
	return &SampEnCalc{m: m, rCoeff: rCoeff}, nil
}

// LoadData loads time series data into the calculator.
// The data length must be greater than m+1 for valid SampEn calculation.
// This method resets the computed flag to false.
//
// Returns an error if the data length is insufficient.
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

// Compute calculates the Sample Entropy of the loaded time series data.
// The calculation uses the parameters m (embedding dimension) and r (tolerance)
// specified during initialization.
//
// Sample Entropy is defined as the negative natural logarithm of the conditional
// probability that two sequences similar for m points remain similar at the next point.
//
// Returns:
//   - sampEn: the calculated Sample Entropy value
//   - error: if data is not loaded, has invalid standard deviation, or if B=0
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
			diff := mathmethod.ChebyshevDistance(c.voltageData[i:i+c.m], c.voltageData[j:j+c.m])
			if diff <= r {
				b += 1
			}
			diff = mathmethod.ChebyshevDistance(c.voltageData[i:i+c.m+1], c.voltageData[j:j+c.m+1])
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

// func (c *SampEnCalc) Compute() (sampEn float64, err error) {
// 	if c.n < c.m+1 {
// 		return 0, fmt.Errorf("no data, use LoadData() first")
// 	}

// 	r := c.rCoeff * mathmethod.StandardDeviation(c.voltageData)
// 	if r <= 0 {
// 		return 0, fmt.Errorf("r<=0, data not valid")
// 	}

// 	windowCount := c.n - c.m
// 	totalA := 0
// 	totalB := 0

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	numWorkers := runtime.NumCPU()
// 	chunkSize := (windowCount + numWorkers - 1) / numWorkers

// 	for w := 0; w < numWorkers; w++ {
// 		start := w * chunkSize
// 		end := start + chunkSize
// 		if end > windowCount {
// 			end = windowCount
// 		}
// 		if start >= end {
// 			continue
// 		}

// 		wg.Add(1)
// 		go func(s, e int) {
// 			defer wg.Done()
// 			localA, localB := 0, 0

// 			for i := s; i < e; i++ {
// 				for j := i + 1; j < windowCount; j++ {
// 					// Calculate m-dimensional distance
// 					diffM := mathmethod.ChebyshevDistance(c.voltageData[i:i+c.m], c.voltageData[j:j+c.m])

// 					if diffM <= r {
// 						localB += 2 // Account for both (i,j) and (j,i)

// 						// Calculate (m+1)-dimensional distance using the previous result
// 						diffNew := math.Abs(c.voltageData[i+c.m] - c.voltageData[j+c.m])
// 						if diffMNew := math.Max(diffM, diffNew); diffMNew <= r {
// 							localA += 2
// 						}
// 					}
// 				}
// 			}

// 			mu.Lock()
// 			totalA += localA
// 			totalB += localB
// 			mu.Unlock()
// 		}(start, end)
// 	}

// 	wg.Wait()

// 	if totalB == 0 {
// 		return 0, fmt.Errorf("B=0, can't compute sample entropy")
// 	}

// 	return -math.Log(float64(totalA) / float64(totalB)), nil
// }
