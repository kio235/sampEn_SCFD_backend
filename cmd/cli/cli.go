package main

import (
	"fmt"

	"github.com/kio235/sampEn_SCFD_backend/internal/core"
	"github.com/kio235/sampEn_SCFD_backend/internal/utils/datareader"
)

func processFile(file string) error {
	reader := datareader.SimpleCSVReader{}
	headers, records, err := reader.ReadName(file)
	if err != nil {
		return err
	}
	detector, err := core.NewDetector(2, 0.2, 100, 1, 0.2, 0.1, 5)
	if err != nil {
		return err
	}
	detector.LoadRecords(headers, records)
	err = detector.Check()
	// fmt.Println(detector.FaultInfos)
	return err
}

func main() {
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/big.csv"
	err := processFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
}
