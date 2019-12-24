package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open file\n"))
	}

	state := parseInput(f, 5)
	f.Close()

	steps, finalState := findRepeat(state, 5, 5)
	fmt.Printf("State %d found in %d steps\n", finalState, steps)
}

func parseInput(r io.Reader, width int) uint32 {
	var result uint32

	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		characters := []rune(line)

		for x, c := range characters {
			if c == '#' {
				pos := y*width + x
				result |= 1 << pos
			}
		}
		y++
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Unable to parse input: %s", err))
	}

	return result
}

type Point struct {
	X int
	Y int
}

var directions []Point = []Point{Point{1, 0}, Point{-1, 0}, Point{0, 1}, Point{0, -1}}

func (pt Point) Encode(width int) uint32 {
	var result uint32

	pos := pt.Y*width + pt.X
	result = 1 << pos

	return result
}

func (pt Point) FindAdjacents(width, height int) uint32 {
	var result uint32

	for _, d := range directions {
		test := Point{pt.X + d.X, pt.Y + d.Y}
		if test.X >= 0 && test.X < width && test.Y >= 0 && test.Y < height {
			result |= test.Encode(width)
		}
	}

	return result
}

func CreateAdjacentsMap(width, height int) map[uint32]uint32 {
	result := make(map[uint32]uint32)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			p := Point{i, j}
			pe := p.Encode(width)
			pa := p.FindAdjacents(width, height)
			result[pe] = pa
		}
	}

	return result
}

func CountBits(n uint32) int {
	// This is Kernighan's algorithm
	count := 0
	for n > 0 {
		n = n & (n - 1)
		count++
	}

	return count
}

func step(state uint32, width, height int, adjacentMap map[uint32]uint32) uint32 {
	newstate := state

	for i := 0; i < width*height; i++ {
		var pos uint32 = 1 << i
		bug := state & pos
		adjacentPos := adjacentMap[pos]
		adjCount := CountBits(adjacentPos & state)

		if bug == 0 && (adjCount == 1 || adjCount == 2) {
			newstate |= pos
		} else if bug > 0 && adjCount != 1 {
			newstate &^= pos
		}
	}

	return newstate
}

func findRepeat(state uint32, width, height int) (int, uint32) {
	m := CreateAdjacentsMap(width, height)
	seen := make(map[uint32]bool)
	count := 0

	seen[state] = true

	for {
		state = step(state, width, height, m)
		count++
		if seen[state] {
			return count, state
		} else {
			seen[state] = true
		}
	}
}
