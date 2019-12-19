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

	j := int64(100)
	for i := int64(0); true; i++ {
		fmt.Printf("Testing: %d,%d\n", i, j-99)
		v := testPoint(Point{i, j}, startingOpcodes)
		if v == 1 {
			v2 := testPoint(Point{i + 99, j - 99}, startingOpcodes)
			if v2 == 1 {
				// success
				fmt.Printf("Point found: %d,%d\n", i, j-99)
				return
			} else {
				j++
				i--
			}
		}
	}

	fmt.Printf("Beam size: %d\n", count)
}

func testPoint(p Point, opcodes map[int64]int64) int64 {
	output := make(chan int64)
	input := make(chan int64)
	program := CopyProgram(opcodes)
	go func() {
		if err := ExecuteProgram(program, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	input <- p.X
	input <- p.Y
	v := <-output

	return v
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
