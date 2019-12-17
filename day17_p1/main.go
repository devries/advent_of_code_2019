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

	go func() {
		if err := ExecuteProgram(startingOpcodes, nil, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	grid := make(map[Point]int64)
	var i, j int64
	for v := range output {
		if v == 10 {
			i = 0
			j++
		} else {
			grid[Point{i, j}] = v
			i++
		}
	}

	printGrid(grid)
	ip := findIntersections(grid)
	fmt.Println(ip)

	var apSum = int64(0)
	for _, p := range ip {
		ap := p.X * p.Y
		apSum += ap
	}

	fmt.Printf("Alignment Parameter Sum: %d\n", apSum)
}

type Point struct {
	X int64
	Y int64
}

var directions = map[int64]Point{
	1: Point{0, -1}, // North
	2: Point{0, 1},  // South
	3: Point{-1, 0}, // West
	4: Point{1, 0},  // East
}

func printGrid(grid map[Point]int64) {
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

func findIntersections(grid map[Point]int64) []Point {
	var r []Point

	for k, v := range grid {
		if v == 35 {
			intersection := true
			for _, p := range directions {
				if grid[Point{k.X + p.X, k.Y + p.Y}] != 35 {
					intersection = false
					break
				}
			}
			if intersection {
				r = append(r, k)
			}
		}
	}

	return r
}
