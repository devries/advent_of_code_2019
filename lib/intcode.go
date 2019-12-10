package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseProgram parses an intcode program for the emulator.
func ParseProgram(program string) (map[int64]int64, error) {
	positions := strings.Split(program, ",")

	opcodes := make(map[int64]int64)

	for i, s := range positions {
		op, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		opcodes[int64(i)] = op
	}

	return opcodes, nil
}

// ExecuteProgram is an Intcode computer emulator. An array of integers representing
// the program is provided by the opcodes argument. Inputs and outputs are provided
// by integer channels
func ExecuteProgram(opcodes map[int64]int64, input <-chan int64, output chan<- int64) error {
	relativeBase := int64(0)
	// counter := int64(0)
	for ptr := int64(0); ; /* ptr < int64(len(opcodes)) */ {
		// Perform instruction parsing
		opcode := opcodes[ptr] % 100
		// counter++

		// fmt.Printf("%d --> Opcode: %d, ptr: %d, fullop: %d, relativeBase: %d, Memory Length: %d\n", counter, opcode, ptr, opcodes[ptr], relativeBase, len(opcodes))
		// fmt.Printf("Program: %v (len: %d)\n\n", opcodes, len(opcodes))
		switch opcode {
		case 1:
			// ADD
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p', 's'}, relativeBase)
			if err != nil {
				return err
			}

			opcodes[p[2]] = p[0] + p[1]
			ptr += 4
		case 2:
			// MULTIPLY
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p', 's'}, relativeBase)
			if err != nil {
				return err
			}

			opcodes[p[2]] = p[0] * p[1]
			ptr += 4
		case 3:
			// INPUT
			p, err := GetAllParameters(opcodes, ptr, []rune{'s'}, relativeBase)
			if err != nil {
				return err
			}

			opcodes[p[0]] = <-input
			ptr += 2
		case 4:
			// OUTPUT
			p, err := GetAllParameters(opcodes, ptr, []rune{'p'}, relativeBase)
			if err != nil {
				return err
			}

			output <- p[0]
			ptr += 2
		case 5:
			// JUMP if TRUE
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p'}, relativeBase)
			if err != nil {
				return err
			}

			if p[0] == 0 {
				ptr += 3
			} else {
				ptr = p[1]
			}
		case 6:
			// JUMP if FALSE
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p'}, relativeBase)
			if err != nil {
				return err
			}

			if p[0] == 0 {
				ptr = p[1]
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p', 's'}, relativeBase)
			if err != nil {
				return err
			}
			if p[0] < p[1] {
				opcodes[p[2]] = 1
			} else {
				opcodes[p[2]] = 0
			}
			ptr += 4
		case 8:
			// EQUALS
			p, err := GetAllParameters(opcodes, ptr, []rune{'p', 'p', 's'}, relativeBase)
			if err != nil {
				return err
			}
			if p[0] == p[1] {
				opcodes[p[2]] = 1
			} else {
				opcodes[p[2]] = 0
			}
			ptr += 4
		case 9:
			// Adjust relative base by Parameter
			p, err := GetAllParameters(opcodes, ptr, []rune{'p'}, relativeBase)
			if err != nil {
				return err
			}
			relativeBase += p[0]
			ptr += 2
		case 99:
			// HALT
			if output != nil {
				close(output)
			}
			return nil
		default:
			if output != nil {
				close(output)
			}
			return fmt.Errorf("Unexpected opcode: %d", opcodes[ptr])
		}
	}
	return fmt.Errorf("Ran out of program without halt")
}

// ParameterValue returns the appropriate value for parameters. The opcodes are the current
// program opcodes, the ptr is the current program pointer, and the parameter is the number
// of the parameter starting with 1 (i.e. 1 is the first parameter, 2 is the second...)
func ParameterValue(opcodes map[int64]int64, ptr int64, parameter int64, relativeBase int64) (int64, error) {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := (opcodes[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// Position mode (return value at the position of parameter)
		return opcodes[opcodes[ptr+parameter]], nil
	case 1:
		// Immediate mode (return value of parameter)
		return opcodes[ptr+parameter], nil
	case 2:
		// Relative mode
		return opcodes[opcodes[ptr+parameter]+relativeBase], nil
	default:
		return 0, fmt.Errorf("Unexpected parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr)
	}
}

// This returns the address to write to using the parameter mode
func ParameterSetAddress(opcodes map[int64]int64, ptr int64, parameter int64, relativeBase int64) (int64, error) {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := (opcodes[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// Address in position mode
		return opcodes[ptr+parameter], nil
	case 2:
		// Address in relative mode
		return opcodes[ptr+parameter] + relativeBase, nil
	default:
		return 0, fmt.Errorf("Unexpected set parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr)
	}
}

func GetAllParameters(opcodes map[int64]int64, ptr int64, format []rune, relativeBase int64) ([]int64, error) {
	ret := make([]int64, len(format))

	for i, f := range format {
		var err error
		if f == 'p' {
			ret[i], err = ParameterValue(opcodes, ptr, int64(i+1), relativeBase)
		} else {
			ret[i], err = ParameterSetAddress(opcodes, ptr, int64(i+1), relativeBase)
		}
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
