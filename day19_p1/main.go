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

	count := 0
	grid := make(map[Point]rune)

	for j := int64(0); j < 50; j++ {
		for i := int64(0); i < 50; i++ {
			output := make(chan int64)
			input := make(chan int64)
			opcodes := CopyProgram(startingOpcodes)
			go func() {
				if err := ExecuteProgram(opcodes, input, output); err != nil {
					panic(fmt.Errorf("Error executing program: %s", err))
				}
			}()

			fmt.Printf("Sending %d,%d\n", i, j)
			input <- i
			input <- j
			v := <-output

			fmt.Printf("Received: %d\n", v)

			if v == 1 {
				count++
				grid[Point{i, j}] = '#'
			} else {
				grid[Point{i, j}] = '.'
			}
		}
	}

	printGrid(grid)
	fmt.Printf("Beam size: %d\n", count)
}

type Point struct {
	X int64
	Y int64
}

var directions = map[int64]Point{
	'^': Point{0, -1}, // North
	'v': Point{0, 1},  // South
	'<': Point{-1, 0}, // West
	'>': Point{1, 0},  // East
}

func printGrid(grid map[Point]rune) {
	var imax, jmax int64

	for k, _ := range grid {
		if k.X > imax {
			imax = k.X
		}
		if k.Y > jmax {
			jmax = k.Y
		}
	}

	for j := int64(0); j <= jmax; j++ {
		for i := int64(0); i <= imax; i++ {
			fmt.Printf("%c", grid[Point{i, j}])
		}
		fmt.Printf("\n")
	}
}

func CopyProgram(in map[int64]int64) map[int64]int64 {
	out := make(map[int64]int64)

	for k, v := range in {
		out[k] = v
	}

	return out
}
