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

	phases := []int{0, 1, 2, 3, 4}
	ch := make(chan []int)

	// Start the generator
	go IntPermutations(phases, ch)

	maxOutput := 0
	maxPhase := make([]int, 5)

	for phasePerm := range ch {
		output := Amplifiers(startingOpcodes, phasePerm)
		if output > maxOutput {
			maxOutput = output
			copy(maxPhase, phasePerm)
		}
	}

	fmt.Printf("Max output: %d\n", maxOutput)
	fmt.Printf("Phases: %v\n", maxPhase)
}

// Run program through 5 amplifiers with phases
func Amplifiers(program []int, phases []int) int {
	signal := 0

	for i := 0; i < 5; i++ {
		ampOpcodes := make([]int, len(program))
		copy(ampOpcodes, program)

		input := []int{phases[i], signal}
		var output []int
		if err := ExecuteProgram(ampOpcodes, input, &output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
		signal = output[0]
	}

	return signal
}

// Permutations of an integer slice (generator)
// Uses Heap's Algorithm (thanks wikipedia)
func IntPermutations(a []int, ch chan<- []int) {
	k := len(a)
	intPermutationsRecursor(k, a, ch)
	close(ch)
}

func intPermutationsRecursor(k int, a []int, ch chan<- []int) {
	if k == 1 {
		output := make([]int, len(a))
		copy(output, a)
		ch <- output
	} else {
		intPermutationsRecursor(k-1, a, ch)

		for i := 0; i < k-1; i++ {
			if k%2 == 0 {
				a[i], a[k-1] = a[k-1], a[i]
			} else {
				a[0], a[k-1] = a[k-1], a[0]
			}
			intPermutationsRecursor(k-1, a, ch)
		}
	}
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
