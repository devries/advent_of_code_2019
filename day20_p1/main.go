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
	if err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}

	printMap(grid)

	portals := findPortals(grid)

	// fmt.Println(portals)

	maze := createMaze(grid, portals)

	// fmt.Println(maze)

	queue := NewStateQueue()
	seen := make(map[Point]bool)

	startingPoint := portals["AA"][0]
	endingPoint := portals["ZZ"][0]

	startState := State{startingPoint, 0}

	queue.Add(startState)

	for queue.Available() {
		state := queue.Pop()
		seen[state.Position] = true
		for _, nextPosition := range maze[state.Position] {
			newstate := State{nextPosition, state.Steps + 1}

			if nextPosition == endingPoint {
				fmt.Printf("Completed in %d steps\n", newstate.Steps)
				return
			} else if seen[nextPosition] {
				continue
			} else {
				queue.Add(newstate)
			}
		}
	}

	fmt.Printf("Didn't find the end\n")
}

type Point struct {
	X int
	Y int
}

var directions = []Point{Point{0, 1}, Point{0, -1}, Point{1, 0}, Point{-1, 0}}

type State struct {
	Position Point
	Steps    int
}

type StateQueue []State

func NewStateQueue() *StateQueue {
	r := []State{}

	return (*StateQueue)(&r)
}

func (sq *StateQueue) Add(s State) {
	*sq = append(*sq, s)
}

func (sq *StateQueue) Pop() State {
	var r State

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

func findPortals(grid map[Point]rune) map[string][]Point {
	result := make(map[string][]Point)

	for p, r := range grid {
		if r >= 'A' && r <= 'Z' {
			nextR := grid[Point{p.X + 1, p.Y}]
			nextD := grid[Point{p.X, p.Y + 1}]

			var p1, p2 Point
			var l2 rune

			if nextR >= 'A' && nextR <= 'Z' {
				// Horizontal
				p1 = Point{p.X - 1, p.Y}
				p2 = Point{p.X + 2, p.Y}
				l2 = nextR
			} else if nextD >= 'A' && nextD <= 'Z' {
				p1 = Point{p.X, p.Y - 1}
				p2 = Point{p.X, p.Y + 2}
				l2 = nextD
			} else {
				continue
			}

			portalName := string([]rune{r, l2})

			var pNext Point
			if grid[p1] == '.' {
				pNext = p1
			} else if grid[p2] == '.' {
				pNext = p2
			} else {
				panic(fmt.Errorf("A portal with no point?"))
			}

			result[portalName] = append(result[portalName], pNext)
		}
	}

	return result
}

func createMaze(grid map[Point]rune, portals map[string][]Point) map[Point]map[Point]Point {
	result := make(map[Point]map[Point]Point)

	for p, r := range grid {
		if r == '.' {
			result[p] = make(map[Point]Point)

			for _, d := range directions {
				pn := Point{p.X + d.X, p.Y + d.Y}
				rn := grid[pn]
				if rn == '.' {
					result[p][d] = pn
				} else if rn >= 'A' && rn <= 'Z' {
					pnn := Point{p.X + 2*d.X, p.Y + 2*d.Y}
					rnn := grid[pnn]

					var portalName string
					if d.X < 0 || d.Y < 0 {
						portalName = string([]rune{rnn, rn})
					} else {
						portalName = string([]rune{rn, rnn})
					}

					portalPoints := portals[portalName]

					for _, pp := range portalPoints {
						if p != pp {
							result[p][d] = pp
						}
					}
				}
			}
		}
	}

	return result
}
