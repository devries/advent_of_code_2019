package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const springscript = `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
WALK
`

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

	opcodes := CopyProgram(startingOpcodes)

	output := make(chan int64)
	input := make(chan int64)
	done := make(chan bool)
	go func() {
		if err := ExecuteProgram(opcodes, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
		done <- true
	}()

	toSend := []rune(springscript)
	letterCounter := 0
	nextLetter := int64(toSend[letterCounter])
ioloop:
	for {
		// IO LOOP
		select {
		case c := <-output:
			if c < 256 {
				fmt.Printf("%c", c)
			} else {
				fmt.Printf("%d\n", c)
			}
		case input <- nextLetter:
			if letterCounter < len(toSend)-1 {
				letterCounter++
			}
			fmt.Printf("%c", nextLetter)
			nextLetter = int64(toSend[letterCounter])
		case <-done:
			break ioloop
		}
	}

}

func CopyProgram(in map[int64]int64) map[int64]int64 {
	out := make(map[int64]int64)

	for k, v := range in {
		out[k] = v
	}

	return out
}
