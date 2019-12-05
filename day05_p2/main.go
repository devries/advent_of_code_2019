package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %s", err))
	}

	program := string(content)
	program = strings.TrimSpace(program)
	startingOpcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	input := []int{5}
	var output []int

	if err := ExecuteProgram(startingOpcodes, input, &output); err != nil {
		panic(fmt.Errorf("Error executing program: %s", err))
	}

	fmt.Println(output)
}

// ParseProgram parses an intcode program for the emulator.
func ParseProgram(program string) ([]int, error) {
	positions := strings.Split(program, ",")

	opcodes := make([]int, len(positions))

	for i, s := range positions {
		op, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		opcodes[i] = op
	}

	return opcodes, nil
}

// ExecuteProgram is an Intcode computer emulator. An array of integers representing
// the program is provided by the opcodes argument. Inputs are provided via the input
// argument, and a pointer to outputs is provided by the output argument.
func ExecuteProgram(opcodes []int, input []int, output *[]int) error {
	for ptr := 0; ptr < len(opcodes); {
		// Perform instruction parsing
		opcode := opcodes[ptr] % 100

		switch opcode {
		case 1:
			// ADD
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, ptr, 1) + ParameterMode(opcodes, ptr, 2)
			ptr += 4
		case 2:
			// MULTIPLY
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, ptr, 1) * ParameterMode(opcodes, ptr, 2)
			ptr += 4
		case 3:
			// INPUT
			opcodes[opcodes[ptr+1]] = input[0]
			input = input[1:]
			ptr += 2
		case 4:
			// OUTPUT
			*output = append(*output, ParameterMode(opcodes, ptr, 1))
			ptr += 2
		case 5:
			// JUMP if TRUE
			if ParameterMode(opcodes, ptr, 1) == 0 {
				ptr += 3
			} else {
				ptr = ParameterMode(opcodes, ptr, 2)
			}
		case 6:
			// JUMP if FALSE
			if ParameterMode(opcodes, ptr, 1) == 0 {
				ptr = ParameterMode(opcodes, ptr, 2)
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			if ParameterMode(opcodes, ptr, 1) < ParameterMode(opcodes, ptr, 2) {
				opcodes[opcodes[ptr+3]] = 1
			} else {
				opcodes[opcodes[ptr+3]] = 0
			}
			ptr += 4
		case 8:
			// EQUALS
			if ParameterMode(opcodes, ptr, 1) == ParameterMode(opcodes, ptr, 2) {
				opcodes[opcodes[ptr+3]] = 1
			} else {
				opcodes[opcodes[ptr+3]] = 0
			}
			ptr += 4
		case 99:
			// HALT
			return nil
		default:
			return fmt.Errorf("Unexpected opcode: %d", opcodes[ptr])
		}
	}
	return fmt.Errorf("Ran out of program without halt")
}

// ParameterMode returns the appropriate value for parameters. The opcodes are the current
// program opcodes, the ptr is the current program pointer, and the parameter is the number
// of the parameter starting with 1 (i.e. 1 is the first parameter, 2 is the second...)
func ParameterMode(opcodes []int, ptr int, parameter int) int {
	j := 10
	for i := 0; i < parameter; i++ {
		j *= 10
	}
	parameterMode := (opcodes[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// Position mode (return value at the position of parameter)
		return opcodes[opcodes[ptr+parameter]]
	case 1:
		// Immediate mode (return value of parameter)
		return opcodes[ptr+parameter]
	default:
		panic(fmt.Errorf("Unexpected parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr))
	}
}
