package mathmethod

import (
	"fmt"
	"math"
)

// mean 计算均值
func Mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// standardDeviation 计算标准差
func StandardDeviation(data []float64) float64 {
	m := Mean(data)
	var sumSquares float64
	for _, v := range data {
		sumSquares += (v - m) * (v - m)
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

func ChebyshevDistance(vec1 []float64, vec2 []float64) (res float64, err error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("vectors have different length ")
	}
	if len(vec1) == 0 {
		return 0, fmt.Errorf("vectors are empty")
	}
	for i := range vec1 {
		diff := math.Abs(vec1[i] - vec2[i])
		if diff > res {
			res = diff
		}
	}
	return res, nil
}
