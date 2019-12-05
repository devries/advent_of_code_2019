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
	starting_opcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	input := []int{5}
	var output []int

	if err := ExecuteProgram(starting_opcodes, input, &output); err != nil {
		panic(fmt.Errorf("Error executing program: %s", err))
	}

	fmt.Println(output)
}

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

func ExecuteProgram(opcodes []int, input []int, output *[]int) error {
	for ptr := 0; ptr < len(opcodes); {
		// Perform instruction parsing
		instruction := opcodes[ptr]
		opcode := instruction % 100
		param_modes := make([]int, 3)
		for i, j := 0, 100; i < 3; i, j = i+1, j*10 {
			param_modes[i] = (instruction / j) % 10
		}

		switch opcode {
		case 1:
			// ADD
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, param_modes[0], ptr+1) + ParameterMode(opcodes, param_modes[1], ptr+2)
			ptr += 4
		case 2:
			// MULTIPLY
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, param_modes[0], ptr+1) * ParameterMode(opcodes, param_modes[1], ptr+2)
			ptr += 4
		case 3:
			// INPUT
			opcodes[opcodes[ptr+1]] = input[0]
			input = input[1:]
			ptr += 2
		case 4:
			// OUTPUT
			*output = append(*output, ParameterMode(opcodes, param_modes[0], ptr+1))
			ptr += 2
		case 5:
			// JUMP if TRUE
			if ParameterMode(opcodes, param_modes[0], ptr+1) == 0 {
				ptr += 3
			} else {
				ptr = ParameterMode(opcodes, param_modes[1], ptr+2)
			}
		case 6:
			// JUMP if FALSE
			if ParameterMode(opcodes, param_modes[0], ptr+1) == 0 {
				ptr = ParameterMode(opcodes, param_modes[1], ptr+2)
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			if ParameterMode(opcodes, param_modes[0], ptr+1) < ParameterMode(opcodes, param_modes[1], ptr+2) {
				opcodes[opcodes[ptr+3]] = 1
			} else {
				opcodes[opcodes[ptr+3]] = 0
			}
			ptr += 4
		case 8:
			// EQUALS
			if ParameterMode(opcodes, param_modes[0], ptr+1) == ParameterMode(opcodes, param_modes[1], ptr+2) {
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

func ParameterMode(opcodes []int, parameter_mode int, ptr int) int {
	switch parameter_mode {
	case 0:
		return opcodes[opcodes[ptr]]
	case 1:
		return opcodes[ptr]
	default:
		panic(fmt.Errorf("Unexpected parameter mode: %d", parameter_mode))
	}
	return 0
}
