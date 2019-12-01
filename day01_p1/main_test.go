package main

import (
	"testing"
)

func TestConversion(t *testing.T) {
	var tests = []struct {
		mass int
		fuel int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}

	for _, test := range tests {
		result := moduleFuel(test.mass)
		if result != test.fuel {
			t.Errorf("For mass %d calculated %d, expected %d", test.mass, result, test.fuel)
		}
	}
}
