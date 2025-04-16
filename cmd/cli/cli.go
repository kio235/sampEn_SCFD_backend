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
	detector, err := core.NewDetector()
	if err != nil {
		return err
	}
	detector.LoadRecords(headers, records)
	err = detector.Check()
	fmt.Println(detector.FaultInfos)
	return err
}

func main() {
	file := "../../data/big.csv"
	err := processFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
}
