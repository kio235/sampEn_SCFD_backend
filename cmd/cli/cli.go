package main

import (
	"fmt"

	"github.com/kio235/sampEn_SCFD_backend/internal/core"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/datareader"
)

func processFile(file string) (int, error) {
	reader := datareader.SimpleCSVReader{}
	headers, records, err := reader.ReadName(file)
	if err != nil {
		return 0, err
	}
	detector, err := core.NewDetector(2, 0.2, 100, 30, 0.2, 0.3)
	if err != nil {
		return 0, err
	}
	detector.LoadRecords(headers, records)
	faultCount, err := detector.Check()
	return faultCount, err
}

func main() {
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv"
	count, err := processFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(count)
}
