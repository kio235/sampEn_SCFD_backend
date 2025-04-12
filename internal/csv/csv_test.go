package csv

import (
	"fmt"
	"testing"
)

func TestReadCSVData(t *testing.T) {
	const FILE string = "/home/kio/code/sampEn_SCFD/sampEn_SCFD_backend/data/2ohm.csv"
	headers, records, length, err := ReadCSVData(FILE)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(headers)
	fmt.Println(length)
	fmt.Println(records[:30])
}
