package test

import "testing"

func BenchmarkAppend(b *testing.B) {
	vol := make([]float64, 1200)
	for i := 0; i < b.N; i++ {
		var dst []float64
		dst = append(dst, vol...)
	}
}

func BenchmarkCopy(b *testing.B) {
	vol := make([]float64, 1200)
	for i := 0; i < b.N; i++ {
		dst := make([]float64, len(vol))
		copy(dst, vol)
	}
}
