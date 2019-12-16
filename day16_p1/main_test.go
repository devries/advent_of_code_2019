package main

import (
	"testing"
)

func TestPhasing(t *testing.T) {
	tests := []struct {
		input  string
		output []int
	}{
		{"80871224585914546619083218645595", []int{2, 4, 1, 7, 6, 1, 7, 6}},
		{"19617804207202209144916044189917", []int{7, 3, 7, 4, 5, 4, 1, 8}},
		{"69317163492948606335995924319873", []int{5, 2, 4, 3, 2, 1, 3, 3}},
	}

	for _, test := range tests {
		n, err := parseNumber(test.input)
		if err != nil {
			t.Errorf("Unable to parse number: %s", err)
		}

		o := RepeatPhase(n, 100)

		for i := 0; i < len(test.output); i++ {
			if o[i] != test.output[i] {
				t.Errorf("Got %v, expected %v", o[:len(test.output)], test.output)
			}
		}
	}
}
