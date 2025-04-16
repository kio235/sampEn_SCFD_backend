package main

import (
	"testing"

	"github.com/kio235/sampEn_SCFD_backend/internal/core"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/datareader"
)

// 基准测试整个流程，包括文件读取
// func BenchmarkProcessFileWithIO(b *testing.B) {
// 	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv" // 确保路径正确
// 	for i := 0; i < b.N; i++ {
// 		_, _ = processFile(file)
// 	}
// }

// 基准测试处理逻辑，排除文件IO
func BenchmarkProcessWithoutIO(b *testing.B) {
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/big.csv"
	// 预先读取数据
	reader := datareader.SimpleCSVReader{}
	headers, records, err := reader.ReadName(file)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer() // 重置计时器，排除准备阶段
	for i := 0; i < b.N; i++ {
		detector, err := core.NewDetector()
		if err != nil {
			b.Fatal(err)
		}
		detector.LoadRecords(headers, records)
		err = detector.Check()
		if err != nil {
			b.Fatal(err)
		}
	}
}
