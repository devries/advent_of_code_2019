package main

import (
	"testing"
)

func TestCycle(t *testing.T) {
	tests := []struct {
		path  []int64
		cycle bool
	}{
		{[]int64{1, 2}, true},
		{[]int64{1, 3, 2, 4, 1}, true},
		{[]int64{1, 1, 1}, false},
	}

	for _, test := range tests {
		cycle := detectCycle(test.path)
		if cycle != test.cycle {
			t.Errorf("Got %t for path %v, expected %t", cycle, test.path, test.cycle)
		}
	}
}
