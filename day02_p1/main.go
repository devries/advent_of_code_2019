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
	opcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	// Enter 1202 program alarm
	opcodes[1] = 12
	opcodes[2] = 2

	err = ExecuteProgram(opcodes)
	if err != nil {
		panic(fmt.Errorf("Error executing program: %s", err))
	}

	fmt.Println(opcodes[0])
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
