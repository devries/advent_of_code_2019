package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	// f, err := os.Open("simple.txt")
	// f, err := os.Open("medium.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	grid, err := parseInput(f)

	printMap(grid)

	queue := NewStateQueue()
	seen := make(map[State]bool)
	all := allKeys(grid)

	maxsteps := 0

	startPositions := getStartingPosition(grid)
	spa, spb, spc, spd := startPositions[0], startPositions[1], startPositions[2], startPositions[3]
	startState := State{spa, spb, spc, spd, KeyStore(0)}
	startStep := StateStep{startState, 0}

	fmt.Println(startState)
	seen[startState] = true
	queue.Add(startStep)

	for queue.Available() {
		state := queue.Pop()
		for _, d := range directions {

			for robot := 0; robot < 4; robot++ {
				positions := []Point{state.PosA, state.PosB, state.PosC, state.PosD}
				positions[robot] = Point{positions[robot].X + d.X, positions[robot].Y + d.Y}
				p := positions[robot]
				keys := state.Keys
				steps := state.Steps + 1
				if steps > maxsteps {
					maxsteps = steps
					fmt.Printf("%d steps %026b\n", maxsteps, keys)
				}

				spa, spb, spc, spd := positions[0], positions[1], positions[2], positions[3]
				newstate := State{spa, spb, spc, spd, keys}
				newstep := StateStep{newstate, steps}

				if grid[p] == '#' {
					// Don't continue into a wall
					continue
				} else if seen[newstate] {
					continue
				} else if grid[p] >= 'a' && grid[p] <= 'z' {
					// Pick up key
					keys = keys.Add(grid[p])
					if keys == all {
						fmt.Printf("Steps: %d\n", steps)
						return
					}
					newstate.Keys = keys
					newstep.Keys = keys

					queue.Add(newstep)
					seen[newstate] = true
				} else if grid[p] >= 'A' && grid[p] <= 'Z' {
					if keys.Contains(grid[p] - 'A' + 'a') {
						seen[newstate] = true
						queue.Add(newstep)
					} else {
						continue
					}
				} else if grid[p] == '.' || grid[p] == '@' {
					for isHallway(grid, p, d) {
						p = nextStep(grid, p, d)
						positions[robot] = p
						steps++
						spa, spb, spc, spd := positions[0], positions[1], positions[2], positions[3]
						newstate = State{spa, spb, spc, spd, keys}
						newstep = StateStep{newstate, steps}
					}
					if !seen[newstate] {
						seen[newstate] = true
						queue.Add(newstep)
					}
				} else {
					fmt.Printf("Unexpected grid point at %v: %c\n", p, grid[p])
					fmt.Printf("OLD  : %+v\n", state)
					fmt.Printf("State: %+v\n", newstep)
					return
				}
			}
		}
	}

}

type Point struct {
	X int
	Y int
}

var directions = []Point{Point{0, 1}, Point{0, -1}, Point{1, 0}, Point{-1, 0}}

type KeyStore uint32

func (k KeyStore) Add(r rune) KeyStore {
	var ret KeyStore

	bit := r - 'a'
	newKey := 1 << bit

	ret = k | KeyStore(newKey)

	return ret
}

func (k KeyStore) Contains(r rune) bool {
	bit := r - 'a'
	newKey := 1 << bit

	if k&KeyStore(newKey) > 0 {
		return true
	} else {
		return false
	}
}

type State struct {
	PosA Point
	PosB Point
	PosC Point
	PosD Point
	Keys KeyStore
}

type StateStep struct {
	State
	Steps int
}

type StateQueue []StateStep

func NewStateQueue() *StateQueue {
	r := []StateStep{}

	return (*StateQueue)(&r)
}

func (sq *StateQueue) Add(s StateStep) {
	*sq = append(*sq, s)
}

func (sq *StateQueue) Pop() StateStep {
	var r StateStep

	if len(*sq) > 0 {
		r = (*sq)[0]
		*sq = (*sq)[1:]
	}

	return r
}

func (sq *StateQueue) Available() bool {
	if len(*sq) > 0 {
		return true
	} else {
		return false
	}
}

func parseInput(r io.Reader) (map[Point]rune, error) {
	result := make(map[Point]rune)

	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		squares := []rune(line)

		for i, v := range squares {
			result[Point{i, y}] = v
		}
		y++
	}

	return result, nil
}

func printMap(c map[Point]rune) {
	minX, maxX, minY, maxY := 0, 0, 0, 0

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

func allKeys(grid map[Point]rune) KeyStore {
	r := KeyStore(0)

	for _, v := range grid {
		if v >= 'a' && v <= 'z' {
			r = r.Add(v)
		}
	}

	return r
}

func getStartingPosition(grid map[Point]rune) []Point {
	r := []Point{}

	for k, v := range grid {
		if v == '@' {
			r = append(r, k)
		}
	}

	return r
}

func isHallway(grid map[Point]rune, p Point, direction Point) bool {
	backtrack := Point{p.X - direction.X, p.Y - direction.Y}

	openCount := 0
	for j := p.Y - 1; j <= p.Y+1; j++ {
		for i := p.X - 1; i <= p.X+1; i++ {
			test := Point{i, j}
			if test == backtrack {
				continue
			}
			if grid[test] == '.' || grid[test] == '@' {
				openCount++
			} else if grid[test] != '#' {
				return false
			}
		}
	}

	if openCount == 2 {
		return true
	} else {
		return false
	}
}

func nextStep(grid map[Point]rune, p Point, direction Point) Point {
	backtrack := Point{p.X - direction.X, p.Y - direction.Y}

	for j := p.Y - 1; j <= p.Y+1; j++ {
		for i := p.X - 1; i <= p.X+1; i++ {
			test := Point{i, j}
			if test == backtrack || test == p {
				continue
			}
			if grid[test] == '.' || grid[test] == '@' {
				return test
			}
		}
	}
	return Point{-1, -1}
}
