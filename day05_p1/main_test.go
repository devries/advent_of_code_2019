package main

import (
	"testing"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		program string
		input   []int
		output  []int
	}{
		{"3,0,4,0,99", []int{1}, []int{1}},
		{"3,0,4,0,99", []int{2}, []int{2}},
	}

	for _, test := range tests {
		input_opcodes, err := ParseProgram(test.program)
		if err != nil {
			t.Errorf("Error parsing input: %s", err)
			continue
		}

		var test_output []int
		err = ExecuteProgram(input_opcodes, test.input, &test_output)
		if err != nil {
			t.Errorf("Error executing program: %s", err)
			continue
		}

		if test_output[0] != test.output[0] {
			t.Errorf("Got %d, expected %d", test_output[0], test.output[0])
		}
	}
}
