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

	input <- int64(1)

	for oline := range output {
		fmt.Println(oline)
	}
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

		// fmt.Printf("Opcode: %d, ptr: %d, fullop: %d, relativeBase: %d\n", opcode, ptr, opcodes[ptr], relativeBase)
		// fmt.Printf("Program: %v (len: %d)\n\n", opcodes, len(opcodes))
		switch opcode {
		case 1:
			// ADD
			SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase),
				ParameterValue(&opcodes, ptr, 1, relativeBase)+ParameterValue(&opcodes, ptr, 2, relativeBase))
			ptr += 4
		case 2:
			// MULTIPLY
			SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase), // GetMemory(&opcodes, ptr+3),
				ParameterValue(&opcodes, ptr, 1, relativeBase)*ParameterValue(&opcodes, ptr, 2, relativeBase))
			ptr += 4
		case 3:
			// INPUT
			SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 1, relativeBase), //GetMemory(&opcodes, ptr+1),
				<-input)
			ptr += 2
		case 4:
			// OUTPUT
			output <- ParameterValue(&opcodes, ptr, 1, relativeBase)
			ptr += 2
		case 5:
			// JUMP if TRUE
			if ParameterValue(&opcodes, ptr, 1, relativeBase) == 0 {
				ptr += 3
			} else {
				ptr = ParameterValue(&opcodes, ptr, 2, relativeBase)
			}
		case 6:
			// JUMP if FALSE
			if ParameterValue(&opcodes, ptr, 1, relativeBase) == 0 {
				ptr = ParameterValue(&opcodes, ptr, 2, relativeBase)
			} else {
				ptr += 3
			}
		case 7:
			// LESS THAN
			if ParameterValue(&opcodes, ptr, 1, relativeBase) < ParameterValue(&opcodes, ptr, 2, relativeBase) {
				SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase), //GetMemory(&opcodes, ptr+3),
					1)
			} else {
				SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase), //GetMemory(&opcodes, ptr+3),
					0)
			}
			ptr += 4
		case 8:
			// EQUALS
			if ParameterValue(&opcodes, ptr, 1, relativeBase) == ParameterValue(&opcodes, ptr, 2, relativeBase) {
				SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase), //GetMemory(&opcodes, ptr+3),
					1)
			} else {
				SetMemory(&opcodes, ParameterSetAddress(&opcodes, ptr, 3, relativeBase), //GetMemory(&opcodes, ptr+3),
					0)
			}
			ptr += 4
		case 9:
			// Adjust relative base by Parameter
			relativeBase += ParameterValue(&opcodes, ptr, 1, relativeBase)
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
func ParameterValue(opcodes *[]int64, ptr int64, parameter int64, relativeBase int64) int64 {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := ((*opcodes)[ptr] / j) % 10

	switch parameterMode {
	case 0:
		// Position mode (return value at the position of parameter)
		return GetMemory(opcodes, GetMemory(opcodes, ptr+parameter))
	case 1:
		// Immediate mode (return value of parameter)
		return GetMemory(opcodes, ptr+parameter)
	case 2:
		// Relative mode
		return GetMemory(opcodes, GetMemory(opcodes, ptr+parameter)+relativeBase)
	default:
		panic(fmt.Errorf("Unexpected parameter mode %d for opcode %d at position %d", parameterMode, (*opcodes)[ptr], ptr))
	}
}

func ParameterSetAddress(opcodes *[]int64, ptr int64, parameter int64, relativeBase int64) int64 {
	j := int64(10)
	for i := int64(0); i < parameter; i++ {
		j *= 10
	}
	parameterMode := ((*opcodes)[ptr] / j) % 10

	switch parameterMode {
	case 0:
		return GetMemory(opcodes, ptr+parameter)
	case 2:
		return GetMemory(opcodes, ptr+parameter) + relativeBase
	default:
		panic(fmt.Errorf("Unexpected set parameter mode %d for opcode %d at position %d", parameterMode, (*opcodes)[ptr], ptr))
	}
}

func GetMemory(memory *[]int64, loc int64) int64 {
	if loc >= int64(len(*memory)) {
		expansion := loc - int64(len(*memory)) + 1
		*memory = append(*memory, make([]int64, expansion)...)
		return 0
	} else {
		return (*memory)[loc]
	}
}

func SetMemory(memory *[]int64, loc int64, value int64) {
	if loc >= int64(len(*memory)) {
		expansion := loc - int64(len(*memory)) + 1
		*memory = append(*memory, make([]int64, expansion)...)
	}

	(*memory)[loc] = value
}
