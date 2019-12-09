package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	// "sync"
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
	_ = startingOpcodes
}

// ParseProgram parses an intcode program for the emulator.
func ParseProgram(program string) ([]int64, error) {
	positions := strings.Split(program, ",")

	opcodes := make([]int64, len(positions))

	for i, s := range positions {
		op, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		opcodes[i] = op
	}

	return opcodes, nil
}

// ExecuteProgram is an Intcode computer emulator. An array of integers representing
// the program is provided by the opcodes argument. Inputs and outputs are provided
// by integer channels
func ExecuteProgram(opcodes []int64, input <-chan int64, output chan<- int64) error {
	relativeBase := int64(0)
	for ptr := int64(0); ptr < int64(len(opcodes)); {
		// Perform instruction parsing
		opcode := opcodes[ptr] % 100

		switch opcode {
		case 1:
			// ADD
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, ptr, 1, relativeBase) + ParameterMode(opcodes, ptr, 2, relativeBase)
			ptr += 4
		case 2:
			// MULTIPLY
			opcodes[opcodes[ptr+3]] = ParameterMode(opcodes, ptr, 1, relativeBase) * ParameterMode(opcodes, ptr, 2, relativeBase)
			ptr += 4
		case 3:
			// INPUT
			opcodes[opcodes[ptr+1]] = <-input
			ptr += 2
		case 4:
			// OUTPUT
			output <- ParameterMode(opcodes, ptr, 1, relativeBase)
			ptr += 2
		case 5:
			// JUMP if TRUE
			if ParameterMode(opcodes, ptr, 1, relativeBase) == 0 {
				ptr += 3
			} else {
				ptr = ParameterMode(opcodes, ptr, 2, relativeBase)
			}
		case 6:
			// JUMP if FALSE
			if ParameterMode(opcodes, ptr, 1, relativeBase) == 0 {
				ptr = ParameterMode(opcodes, ptr, 2, relativeBase)
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			if ParameterMode(opcodes, ptr, 1, relativeBase) < ParameterMode(opcodes, ptr, 2, relativeBase) {
				opcodes[opcodes[ptr+3]] = 1
			} else {
				opcodes[opcodes[ptr+3]] = 0
			}
			ptr += 4
		case 8:
			// EQUALS
			if ParameterMode(opcodes, ptr, 1, relativeBase) == ParameterMode(opcodes, ptr, 2, relativeBase) {
				opcodes[opcodes[ptr+3]] = 1
			} else {
				opcodes[opcodes[ptr+3]] = 0
			}
			ptr += 4
		case 99:
			// HALT
			close(output)
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
func ParameterMode(opcodes []int64, ptr int64, parameter int64, relativeBase int64) int64 {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
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
	case 2:
		// Relative mode
		return opcodes[opcodes[ptr+parameter]+relativeBase]
	default:
		panic(fmt.Errorf("Unexpected parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr))
	}
}

func GetMemory(memory *[]int64, loc int) int64 {
	if loc > len(*memory) {
		expansion := loc - len(*memory) + 1
		*memory = append(*memory, make([]int64, expansion)...)
		return 0
	} else {
		return (*memory)[loc]
	}
}

func SetMemory(memory *[]int64, loc int, value int64) {
	if loc > len(*memory) {
		expansion := loc - len(*memory) + 1
		*memory = append(*memory, make([]int64, expansion)...)
	}

	(*memory)[loc] = value
}
