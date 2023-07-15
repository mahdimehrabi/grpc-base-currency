package data

import (
	"fmt"
	"testing"
)

func TestNewRates(t *testing.T) {
	tr, err := NewRates()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", tr.rates)
}
