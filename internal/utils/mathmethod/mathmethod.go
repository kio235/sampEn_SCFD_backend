package mathmethod

import (
	"math"
)

// Mean 计算均值
func Mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// StandardDeviation 计算标准差
func StandardDeviation(data []float64) float64 {
	m := Mean(data)
	var sumSquares float64
	for _, v := range data {
		sumSquares += (v - m) * (v - m)
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

// func ChebyshevDistance(vec1 []float64, vec2 []float64) (res float64, err error) {
// 	if len(vec1) != len(vec2) {
// 		return 0, fmt.Errorf("vectors have different length ")
// 	}
// 	if len(vec1) == 0 {
// 		return 0, fmt.Errorf("vectors are empty")
// 	}
// 	for i := range vec1 {
// 		diff := math.Abs(vec1[i] - vec2[i])
// 		if diff > res {
// 			res = diff
// 		}
// 	}
// 	return res, nil
// }

// ChebyshevDistance calculates the maximum absolute difference between elements of two vectors.
func ChebyshevDistance(vec1, vec2 []float64) float64 {
	maxDiff := 0.0
	for i := range vec1 {
		if diff := math.Abs(vec1[i] - vec2[i]); diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}

func IsMultiple(a, b int) bool {
	// 避免除以零的错误
	if b == 0 {
		return false
	}
	return a%b == 0
}

func ConvertTo2DShared[T any](s []T, cols, rows int) [][]T {
	if rows*cols != len(s) {
		return nil // 或处理错误
	}
	result := make([][]T, cols)
	for i := range cols {
		result[i] = make([]T, rows)
	}
	for i, val := range s {
		col := i % cols
		row := i / cols
		result[col][row] = val
	}
	return result
}
