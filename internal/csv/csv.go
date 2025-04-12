package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadCSVData(f string) (headers []string, records []float64, length int, err error) {
	// 打开CSV文件
	file, err := os.Open(f)
	if err != nil {
		fmt.Println("打开文件出错:", err)
		return nil, nil, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 读取标题行
	header, err := reader.Read()
	if err != nil {
		fmt.Println("读取标题行出错:", err)
		return nil, nil, 0, err
	}
	fmt.Println("标题:", header)
	headers = header

	// 逐行读取数据
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取记录出错:", err)
			return nil, nil, 0, err
		}

		// 检查列数一致性
		if len(record) != len(headers) {
			return nil, nil, 0, fmt.Errorf("列数不匹配: 标题%d列, 记录%d列", len(headers), len(record))
		}

		// 转换每个值为浮点数
		for _, v := range record {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				fmt.Println("转换数值出错:", err)
				return nil, nil, 0, fmt.Errorf("'%s' 转换失败: %v", v, err)
			}
			records = append(records, num)
		}
		length++
	}
	return headers, records, length, nil
}
