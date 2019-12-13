package main

import (
	"fmt"
	"io/ioutil"
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

	output := make(chan int64)
	screen := make(map[Point]int64)

	go func() {
		if err := ExecuteProgram(startingOpcodes, nil, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	for {
		x, more := <-output
		if more == false {
			break
		}
		y := <-output
		tileID := <-output

		screen[Point{x, y}] = tileID
	}

	blockCount := 0
	for _, v := range screen {
		if v == int64(2) {
			blockCount++
		}
	}

	fmt.Println(blockCount)
}

type Point struct {
	X int64
	Y int64
}
