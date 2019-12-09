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

	output := make(chan int64)
	input := make(chan int64)
	go func() {
		if err := ExecuteProgram(startingOpcodes, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	input <- int64(2)

	for oline := range output {
		fmt.Println(oline)
	}
}

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
			opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] =
				ParameterValue(opcodes, ptr, 1, relativeBase) + ParameterValue(opcodes, ptr, 2, relativeBase)
			ptr += 4
		case 2:
			// MULTIPLY
			opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] =
				ParameterValue(opcodes, ptr, 1, relativeBase) * ParameterValue(opcodes, ptr, 2, relativeBase)
			ptr += 4
		case 3:
			// INPUT
			opcodes[ParameterSetAddress(opcodes, ptr, 1, relativeBase)] = <-input
			ptr += 2
		case 4:
			// OUTPUT
			output <- ParameterValue(opcodes, ptr, 1, relativeBase)
			ptr += 2
		case 5:
			// JUMP if TRUE
			if ParameterValue(opcodes, ptr, 1, relativeBase) == 0 {
				ptr += 3
			} else {
				ptr = ParameterValue(opcodes, ptr, 2, relativeBase)
			}
		case 6:
			// JUMP if FALSE
			if ParameterValue(opcodes, ptr, 1, relativeBase) == 0 {
				ptr = ParameterValue(opcodes, ptr, 2, relativeBase)
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			if ParameterValue(opcodes, ptr, 1, relativeBase) < ParameterValue(opcodes, ptr, 2, relativeBase) {
				opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] = 1
			} else {
				opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] = 0
			}
			ptr += 4
		case 8:
			// EQUALS
			if ParameterValue(opcodes, ptr, 1, relativeBase) == ParameterValue(opcodes, ptr, 2, relativeBase) {
				opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] = 1
			} else {
				opcodes[ParameterSetAddress(opcodes, ptr, 3, relativeBase)] = 0
			}
			ptr += 4
		case 9:
			// Adjust relative base by Parameter
			relativeBase += ParameterValue(opcodes, ptr, 1, relativeBase)
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
func ParameterValue(opcodes map[int64]int64, ptr int64, parameter int64, relativeBase int64) int64 {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := (opcodes[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// Position mode (return value at the position of parameter)
		// return GetMemory(opcodes, GetMemory(opcodes, ptr+parameter))
		return opcodes[opcodes[ptr+parameter]]
	case 1:
		// Immediate mode (return value of parameter)
		// return GetMemory(opcodes, ptr+parameter)
		return opcodes[ptr+parameter]
	case 2:
		// Relative mode
		// return GetMemory(opcodes, GetMemory(opcodes, ptr+parameter)+relativeBase)
		return opcodes[opcodes[ptr+parameter]+relativeBase]
	default:
		panic(fmt.Errorf("Unexpected parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr))
	}
}

func ParameterSetAddress(opcodes map[int64]int64, ptr int64, parameter int64, relativeBase int64) int64 {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := (opcodes[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// return GetMemory(opcodes, ptr+parameter)
		return opcodes[ptr+parameter]
	case 2:
		// return GetMemory(opcodes, ptr+parameter) + relativeBase
		return opcodes[ptr+parameter] + relativeBase
	default:
		panic(fmt.Errorf("Unexpected set parameter mode %d for opcode %d at position %d", parameterMode, opcodes[ptr], ptr))
	}
}
