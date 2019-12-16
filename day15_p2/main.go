package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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

	compartment := make(map[Point]rune)
	compartment[Point{0, 0}] = ' '

	fmt.Printf("\033[H\033[J\033[H")

	explore(nil, 1, input, output, compartment)
	explore(nil, 2, input, output, compartment)
	explore(nil, 3, input, output, compartment)
	explore(nil, 4, input, output, compartment)

	// fmt.Printf("\033[H\033[J\033[H")
	fmt.Printf("\033[3J\033[H")
	printMap(compartment)

	minutes := 0
	for {
		if isFull(compartment) {
			break
		}
		oxygenStep(compartment)
		minutes++
		fmt.Printf("\033[3J\033[H")
		printMap(compartment)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Printf("Compartment filled in %d minutes\n", minutes)

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

func getLocation(path []int64) Point {
	position := Point{0, 0}

	for _, v := range path {
		motion := directions[v]
		position.X += motion.X
		position.Y += motion.Y
	}

	return position
}

func explore(path []int64, direction int64, input chan int64, output chan int64, compartment map[Point]rune) {
	input <- direction
	status := <-output

	var newpath []int64
	if path == nil {
		newpath = []int64{direction}
	} else {
		newpath = make([]int64, len(path)+1)
		copy(newpath, path)
		newpath[len(path)] = direction
	}

	if status == 0 { // hit wall
		compartment[getLocation(newpath)] = '#'
		return
	} else { // successful step
		if status == 1 {
			compartment[getLocation(newpath)] = ' '
		} else {
			compartment[getLocation(newpath)] = 'O'
		}
		fmt.Printf("\033[3J\033[H")
		printMap(compartment)
		time.Sleep(50 * time.Millisecond)

		if !detectCycle(newpath) {
			if direction != 2 {
				explore(newpath, 1, input, output, compartment)
			}
			if direction != 1 {
				explore(newpath, 2, input, output, compartment)
			}
			if direction != 4 {
				explore(newpath, 3, input, output, compartment)
			}
			if direction != 3 {
				explore(newpath, 4, input, output, compartment)
			}
		}
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

func printMap(c map[Point]rune) {
	minX, maxX, minY, maxY := int64(0), int64(0), int64(0), int64(0)

	for v, _ := range c {
		if v.X < minX {
			minX = v.X
		} else if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		} else if v.Y > maxY {
			maxY = v.Y
		}
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			p := c[Point{i, j}]
			if p == 0 {
				fmt.Printf("?")
			} else {
				fmt.Printf("%c", p)
			}
		}
		fmt.Printf("\n")
	}
}

// Check if oxygen has filled the compartment
func isFull(c map[Point]rune) bool {
	for _, v := range c {
		if v == ' ' {
			return false
		}
	}

	return true
}

// Step oxygen
func oxygenStep(c map[Point]rune) {
	var fillPoints []Point

	for k, v := range c {
		if v == 'O' {
			for _, d := range directions {
				test := Point{k.X + d.X, k.Y + d.Y}
				if c[test] == ' ' {
					fillPoints = append(fillPoints, test)
				}
			}
		}
	}

	for _, v := range fillPoints {
		c[v] = 'O'
	}
}
