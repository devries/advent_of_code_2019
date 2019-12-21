package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("adventure.ic")
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

	keyboard := make(chan rune)

	go func() {
		// Read keyboard input
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			v := scanner.Text() + "\n"
			for _, r := range []rune(v) {
				keyboard <- r
			}
		}
	}()

	toSend := []rune{}

ioloop:
	for {
		// IO LOOP

		if len(toSend) > 0 {
			select {
			case c := <-output:
				if c < 256 {
					fmt.Printf("%c", c)
				} else {
					fmt.Printf("%d\n", c)
				}
			case c := <-keyboard:
				toSend = append(toSend, c)
			case input <- int64(toSend[0]):
				toSend = toSend[1:]
				// fmt.Printf("%c", nextLetter)
				// nextLetter = int64(toSend[letterCounter])
			case <-done:
				break ioloop
			}
		} else {
			select {
			case c := <-output:
				if c < 256 {
					fmt.Printf("%c", c)
				} else {
					fmt.Printf("%d\n", c)
				}
			case c := <-keyboard:
				toSend = append(toSend, c)
			case <-done:
				break ioloop
			}
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
