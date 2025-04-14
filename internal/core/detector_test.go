package core

import (
	"fmt"
	"testing"

	"github.com/kio235/sampEn_SCFD_backend/internal/utils/datareader"
)

func TestLoadRecords(t *testing.T) {
	reader := datareader.SimpleCSVReader{}
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv"
	headers, records, err := reader.ReadName(file)
	if err != nil {
		t.Error(err)
	}
	detector, err := NewDetector(2, 0.2, 100, 30, 0.2, 0.2)
	if err != nil {
		t.Error(err)
	}
	err = detector.LoadRecords(headers, records)
	if err != nil {
		t.Error(err)
	}
}

func TestCompute(t *testing.T) {
	reader := datareader.SimpleCSVReader{}
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv"
	headers, records, err := reader.ReadName(file)
	if err != nil {
		t.Error(err)
	}
	detector, err := NewDetector(2, 0.2, 100, 30, 0.2, 0.2)
	if err != nil {
		t.Error(err)
	}
	err = detector.LoadRecords(headers, records)
	if err != nil {
		t.Error(err)
	}
	err = detector.Compute()
	if err != nil {
		t.Error(err)
	}
}

func TestCheck(t *testing.T) {
	reader := datareader.SimpleCSVReader{}
	file := "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv"
	headers, records, err := reader.ReadName(file)
	if err != nil {
		t.Error(err)
	}
	detector, err := NewDetector(2, 0.2, 100, 30, 0.2, 0.3)
	if err != nil {
		t.Error(err)
	}
	err = detector.LoadRecords(headers, records)
	if err != nil {
		t.Error(err)
	}
	faultCount, err := detector.Check()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(faultCount)
}
