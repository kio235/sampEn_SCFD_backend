package core

import (
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
	detector, err := NewDetector()
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
	detector, err := NewDetector()
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
	detector, err := NewDetector()
	if err != nil {
		t.Error(err)
	}
	err = detector.LoadRecords(headers, records)
	if err != nil {
		t.Error(err)
	}
	err = detector.Check()
	if err != nil {
		t.Error(err)
	}
}
