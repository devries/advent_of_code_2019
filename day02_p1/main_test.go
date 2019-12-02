package main

import (
	"testing"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{"1,0,0,0,99", "2,0,0,0,99"},
		{"2,3,0,3,99", "2,3,0,6,99"},
		{"2,4,4,5,99,0", "2,4,4,5,99,9801"},
		{"1,1,1,4,99,5,6,0,99", "30,1,1,4,2,5,6,0,99"},
	}

	for _, test := range tests {
		input_opcodes, err := ParseProgram(test.input)
		if err != nil {
			t.Errorf("Error parsing input: %s", err)
		}
		output_opcodes, err := ParseProgram(test.output)
		if err != nil {
			t.Errorf("Error parsing output: %s", err)
		}

		err = ExecuteProgram(input_opcodes)
		if err != nil {
			t.Errorf("Error executing program: %s", err)
		}

		if !OpcodesEqual(input_opcodes, output_opcodes) {
			t.Errorf("Got %v, expected %v", input_opcodes, output_opcodes)
		}
	}
}

func OpcodesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
