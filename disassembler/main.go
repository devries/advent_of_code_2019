package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Operation struct {
	Name string
	Args int64
}

// Operation assembly names
var instructions = map[int64]Operation{
	1:  Operation{"add", 3}, // add
	2:  Operation{"mul", 3}, // multiply
	3:  Operation{"in", 1},  // read input
	4:  Operation{"out", 1}, // write output
	5:  Operation{"jt", 2},  // jump if true
	6:  Operation{"jf", 2},  // jump if false
	7:  Operation{"lt", 3},  // less than
	8:  Operation{"eq", 3},  // equal
	9:  Operation{"rpo", 1}, // move relative base by parameter
	99: Operation{"hlt", 0}, // halt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <intcode_filename>\n", os.Args[0])
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(fmt.Errorf("Error opening file: %s", err))
	}

	program := string(content)
	program = strings.TrimSpace(program)
	opcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	maxAddress := findMaxAddress(opcodes)

	for i := int64(0); i <= maxAddress; {
		step, inst := parseInstruction(opcodes, i)
		if inst != "" {
			fmt.Printf("%6d: %s\n", i, inst)
		}
		i += step
	}
}

// Return length of instruction, and instruction text
func parseInstruction(opcodes map[int64]int64, ptr int64) (int64, string) {
	inst := opcodes[ptr]

	if inst == 0 {
		return 1, ""
	}

	op := inst % 100 // First two digits are op code
	oper := instructions[op]
	if oper.Name == "" {
		return 1, fmt.Sprintf("dat %d", inst)
	}

	args := []string{oper.Name}

	for i, j := int64(1), int64(100); i <= oper.Args; i, j = i+1, j*10 {
		ptype := (inst / j) % 10
		var prefix string

		switch ptype {
		case 0:
			prefix = ""
		case 1:
			prefix = "!"
		case 2:
			prefix = "%"
		default:
			// This doesn't look like it was actually an instruction
			return 1, fmt.Sprintf("dat %d", inst)
		}
		args = append(args, fmt.Sprintf("%s%d", prefix, opcodes[ptr+i]))
	}

	return oper.Args + 1, strings.Join(args, " ")
}

func findMaxAddress(opcodes map[int64]int64) int64 {
	var max int64

	for k, _ := range opcodes {
		if k > max {
			max = k
		}
	}

	return max
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
