package main

import (
	"testing"
)

func TestSolution(t *testing.T) {
	var tests = []struct {
		trial   int
		success bool
	}{
		{112233, true},
		{123444, false},
		{111122, true},
	}

	for _, test := range tests {
		result := ValidateCode(test.trial)
		if result != test.success {
			t.Errorf("For %d Expected %t got %t", test.trial, test.success, result)
		}
	}
}
