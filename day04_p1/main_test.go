package main

import (
	"testing"
)

func TestSolution(t *testing.T) {
	var tests = []struct {
		trial   int
		success bool
	}{
		{111111, true},
		{223450, false},
		{123789, false},
	}

	for _, test := range tests {
		result := ValidateCode(test.trial)
		if result != test.success {
			t.Errorf("For %d Expected %t got %t", test.trial, test.success, result)
		}
	}
}
