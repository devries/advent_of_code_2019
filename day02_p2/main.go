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

	desired_output := 19690720

	// Iterate through looking for desired result
	for noun := 0; noun < 100; noun += 1 {
		for verb := 0; verb < 100; verb += 1 {
			opcodes := make([]int, len(starting_opcodes))
			copy(opcodes, starting_opcodes)

			// Enter trial noun and verb
			opcodes[1] = noun
			opcodes[2] = verb

			err = ExecuteProgram(opcodes)
			if err != nil {
				panic(fmt.Errorf("Error executing program: %s", err))
			}

			if opcodes[0] == desired_output {
				fmt.Printf("Noun = %d, Verb = %d, Solution = %d\n", noun, verb, noun*100+verb)
				return
			}
		}
	}
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

func ExecuteProgram(opcodes []int) error {
	for ptr := 0; ptr < len(opcodes); ptr += 4 {
		switch opcodes[ptr] {
		case 1:
			opcodes[opcodes[ptr+3]] = opcodes[opcodes[ptr+1]] + opcodes[opcodes[ptr+2]]
		case 2:
			opcodes[opcodes[ptr+3]] = opcodes[opcodes[ptr+1]] * opcodes[opcodes[ptr+2]]
		case 99:
			return nil
		default:
			return fmt.Errorf("Unexpected opcode: %d", opcodes[ptr])
		}
	}
	return fmt.Errorf("Ran out of program without halt")
}
