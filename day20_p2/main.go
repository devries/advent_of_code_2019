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
	seen := make(map[SeenState]bool)

	startingPoint := portals["AA"][0]
	endingPoint := portals["ZZ"][0]

	startState := State{startingPoint, 0, 0}

	queue.Add(startState)

	for queue.Available() {
		state := queue.Pop()
		seen[SeenState{state.Position, state.Level}] = true
		for _, connection := range maze[state.Position] {
			nextPosition := connection.Position
			nextLevel := state.Level + connection.LevelChange
			newstate := State{nextPosition, state.Steps + 1, nextLevel}

			if nextPosition == endingPoint && nextLevel == 0 {
				fmt.Printf("Completed in %d steps\n", newstate.Steps)
				return
			} else if seen[SeenState{nextPosition, nextLevel}] {
				continue
			} else if nextLevel < 0 {
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

type Connection struct {
	Position    Point
	LevelChange int
}

var directions = []Point{Point{0, 1}, Point{0, -1}, Point{1, 0}, Point{-1, 0}}

type State struct {
	Position Point
	Steps    int
	Level    int
}

type SeenState struct {
	Position Point
	Level    int
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

func createMaze(grid map[Point]rune, portals map[string][]Point) map[Point]map[Point]Connection {
	result := make(map[Point]map[Point]Connection)
	min, max := findMazeSize(grid)

	for p, r := range grid {
		if r == '.' {
			result[p] = make(map[Point]Connection)

			for _, d := range directions {
				pn := Point{p.X + d.X, p.Y + d.Y}
				rn := grid[pn]
				if rn == '.' {
					result[p][d] = Connection{pn, 0}
				} else if rn >= 'A' && rn <= 'Z' {
					var levelChange int
					if pn.X < min.X || pn.X > max.X || pn.Y < min.Y || pn.Y > max.Y {
						levelChange = -1
					} else {
						levelChange = 1
					}

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
							result[p][d] = Connection{pp, levelChange}
						}
					}
				}
			}
		}
	}

	return result
}

func findMazeSize(grid map[Point]rune) (Point, Point) {
	min := Point{-1, -1}
	max := Point{-1, -1}

	for p, r := range grid {
		if r == '#' {
			if min.X == -1 || p.X < min.X {
				min.X = p.X
			}
			if min.Y == -1 || p.Y < min.Y {
				min.Y = p.Y
			}
			if max.X == -1 || p.X > max.X {
				max.X = p.X
			}
			if max.Y == -1 || p.Y > max.Y {
				max.Y = p.Y
			}
		}
	}

	return min, max
}
