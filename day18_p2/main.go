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

	// Each quadrant's path is independant
	q1 := Search(grid, Point{0, 0}, Point{40, 40})
	fmt.Println(q1)
	q2 := Search(grid, Point{40, 0}, Point{80, 40})
	fmt.Println(q2)
	q3 := Search(grid, Point{0, 40}, Point{40, 80})
	fmt.Println(q3)
	q4 := Search(grid, Point{40, 40}, Point{80, 80})

	fmt.Printf("Total: %d\n", q1+q2+q3+q4)
}

func Search(grid map[Point]rune, min Point, max Point) int {
	queue := NewStateQueue()
	seen := make(map[State]bool)
	all := allKeys(grid, min, max)

	sp := getStartingPosition(grid, min, max)

	startState := State{sp, KeyStore(0)}
	startStep := StateStep{startState, 0}

	seen[startState] = true
	queue.Add(startStep)

	for queue.Available() {
		state := queue.Pop()
		for _, d := range directions {
			p := Point{state.Position.X + d.X, state.Position.Y + d.Y}
			keys := state.Keys
			steps := state.Steps + 1

			newstate := State{p, keys}
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
					return steps
				}
				newstate.Keys = keys
				newstep.Keys = keys

				queue.Add(newstep)
				seen[newstate] = true
			} else if grid[p] >= 'A' && grid[p] <= 'Z' {
				if !all.Contains(grid[p]-'A'+'a') || keys.Contains(grid[p]-'A'+'a') {
					seen[newstate] = true
					queue.Add(newstep)
				} else {
					continue
				}
			} else if grid[p] == '.' || grid[p] == '@' {
				seen[newstate] = true
				queue.Add(newstep)
			} else {
				fmt.Printf("Unexpected grid point: %c\n", grid[p])
			}
		}
	}
	return 0
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
	Position Point
	Keys     KeyStore
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

func allKeys(grid map[Point]rune, min Point, max Point) KeyStore {
	r := KeyStore(0)

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			v := grid[Point{i, j}]
			if v >= 'a' && v <= 'z' {
				r = r.Add(v)
			}
		}
	}

	return r
}

func getStartingPosition(grid map[Point]rune, min Point, max Point) Point {
	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			v := grid[Point{i, j}]
			if v == '@' {
				return Point{i, j}
			}
		}
	}

	return Point{-1, -1}
}
