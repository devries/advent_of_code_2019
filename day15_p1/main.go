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
	input := make(chan int64)

	go func() {
		if err := ExecuteProgram(startingOpcodes, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	var results []int

	explore(nil, 1, input, output, &results)
	explore(nil, 2, input, output, &results)
	explore(nil, 3, input, output, &results)
	explore(nil, 4, input, output, &results)

	best := Minimum(results)
	fmt.Printf("Shortest Path Length: %d\n", best)
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

func detectCycle(path []int64) bool {
	visited := make(map[Point]bool)

	position := Point{0, 0}
	visited[position] = true

	for _, v := range path {
		motion := directions[v]
		position.X += motion.X
		position.Y += motion.Y
		if visited[position] {
			return true
		} else {
			visited[position] = true
		}
	}

	return false
}

func explore(path []int64, direction int64, input chan int64, output chan int64, results *[]int) {
	input <- direction
	status := <-output

	switch status {
	case 0: // hit wall
		return
	case 1: // successful step
		var newpath []int64
		if path == nil {
			newpath = []int64{direction}
		} else {
			newpath = make([]int64, len(path)+1)
			copy(newpath, path)
			newpath[len(path)] = direction
		}
		if !detectCycle(newpath) {
			if direction != 2 {
				explore(newpath, 1, input, output, results)
			}
			if direction != 1 {
				explore(newpath, 2, input, output, results)
			}
			if direction != 4 {
				explore(newpath, 3, input, output, results)
			}
			if direction != 3 {
				explore(newpath, 4, input, output, results)
			}
		}
	case 2:
		var newpath []int64
		if path == nil {
			newpath = []int64{direction}
		} else {
			newpath = make([]int64, len(path)+1)
			copy(newpath, path)
			newpath[len(path)] = direction
		}

		pathLength := len(newpath)
		*results = append(*results, pathLength)
		fmt.Printf("Path length: %d\n", len(newpath))
	}
	switch direction {
	case 1:
		input <- 2
		<-output
	case 2:
		input <- 1
		<-output
	case 3:
		input <- 4
		<-output
	case 4:
		input <- 3
		<-output
	}
}

func Minimum(a []int) int {
	var r int
	for i, v := range a {
		if i == 0 {
			r = v
		} else {
			if v < r {
				r = v
			}
		}
	}

	return r
}
