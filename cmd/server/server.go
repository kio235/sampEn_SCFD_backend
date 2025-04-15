package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kio235/sampEn_SCFD_backend/internal/core"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/datareader"
)

// 已实现的 CSV 处理函数（示例）
// func ProcessCSV(filePath string) ([]DiagnosisResult, error) {
// 	// 你的 CSV 处理逻辑
// 	// 返回诊断结果和可能的错误
// }

// // 已实现的诊断函数（示例）
// func RunDiagnosis(csvData [][]string) ([]DiagnosisResult, error) {
// 	// 你的诊断逻辑
// 	// 返回诊断结果和可能的错误
// }

func main() {
	router := gin.Default()

	// 限制上传文件大小为 10MB
	router.MaxMultipartMemory = 10 << 20 // 10MB

	router.POST("/diagnose", func(c *gin.Context) {
		// 1. 接收上传文件
		fileHeader, err := c.FormFile("csv_file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing or invalid file",
			})
			return
		}

		// 2. 验证文件类型
		if !strings.HasSuffix(fileHeader.Filename, ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Only CSV files are allowed",
			})
			return
		}

		// 3. 打开上传文件
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to open file",
			})
			return
		}
		defer file.Close()

		// 4. 读取CSV内容
		reader := datareader.SimpleCSVReader{}
		cels, records, err := reader.Read(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid CSV format",
			})
			return
		}

		detector, err := core.NewDetector(2, 0.2, 100, 1, 0.2, 0.4, 5)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create detector",
			})
			return
		}
		err = detector.LoadRecords(cels, records)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Load data error. Invalid data",
			})
			return
		}
		err = detector.Check()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Check data error. Invalid data",
			})
			return
		}

		// 5. 执行诊断（调用已有函数）

		// 6. 返回诊断结果
		c.JSON(http.StatusOK, gin.H{
			"filename": fileHeader.Filename,
			"results":  detector.FaultInfos,
		})
	})

	router.Run(":8080")
}
