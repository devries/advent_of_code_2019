package main

import (
	"strconv"
	"testing"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		program string
		output  []int64
	}{
		{"109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99",
			[]int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}},
		{"104,1125899906842624,99", []int64{1125899906842624}},
	}

	for _, test := range tests {
		inputOpcodes, err := ParseProgram(test.program)
		if err != nil {
			t.Errorf("Error parsing input: %s", err)
			continue
		}

		output := make(chan int64)
		go func() {
			if err := ExecuteProgram(inputOpcodes, nil, output); err != nil {
				t.Errorf("Execution failed: %s", err)
			}
		}()

		var results []int64
		for r := range output {
			results = append(results, r)
		}

		if !OpcodesEqual(results, test.output) {
			t.Errorf("Got %d, expected %d", output, test.output)
		}
	}
}

func TestOutputLength(t *testing.T) {
	var tests = []struct {
		program      string
		outputLength int
	}{
		{"1102,34915192,34915192,7,4,7,99,0", 16},
	}

	for _, test := range tests {
		inputOpcodes, err := ParseProgram(test.program)
		if err != nil {
			t.Errorf("Error parsing input: %s", err)
			continue
		}

		output := make(chan int64)
		go func() {
			if err := ExecuteProgram(inputOpcodes, nil, output); err != nil {
				t.Errorf("Execution failed: %s", err)
			}
		}()

		var results []int64
		for r := range output {
			results = append(results, r)
		}

		if len(strconv.FormatInt(results[0], 10)) != test.outputLength {
			t.Errorf("Got %d, expected %d", output, test.outputLength)
		}
	}
}

func OpcodesEqual(a, b []int64) bool {
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
