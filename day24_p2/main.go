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

	state := make(map[int]uint32)

	state[0] = parseInput(f, 5)
	f.Close()

	adjacents := CreateAdjacentsMap(5, 5)

	for i, v := range adjacents {
		fmt.Printf("%d, %d\n", i, v)
	}

	for i := 0; i < 200; i++ {
		state = step(state, 5, 5, adjacents)
	}

	drawLevels(state, 5, 5)
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

func (pt Point) FindAdjacents(width, height int) (uint32, uint32, uint32) {
	var levelOutside uint32
	var levelSame uint32
	var levelInside uint32

	center := Point{width / 2, height / 2}

	for _, d := range directions {
		test := Point{pt.X + d.X, pt.Y + d.Y}
		if test == center {
			if d == (Point{1, 0}) {
				for i := 0; i < height; i++ {
					levelInside |= Point{0, i}.Encode(width)
				}
			} else if d == (Point{-1, 0}) {
				for i := 0; i < height; i++ {
					levelInside |= Point{width - 1, i}.Encode(width)
				}
			} else if d == (Point{0, 1}) {
				for i := 0; i < width; i++ {
					levelInside |= Point{i, 0}.Encode(width)
				}
			} else if d == (Point{0, -1}) {
				for i := 0; i < width; i++ {
					levelInside |= Point{i, height - 1}.Encode(width)
				}
			}
		} else if test.X < 0 {
			levelOutside |= Point{center.X - 1, center.Y}.Encode(width)
		} else if test.X >= width {
			levelOutside |= Point{center.X + 1, center.Y}.Encode(width)
		} else if test.Y < 0 {
			levelOutside |= Point{center.X, center.Y - 1}.Encode(width)
		} else if test.Y >= height {
			levelOutside |= Point{center.X, center.Y + 1}.Encode(width)
		} else {
			levelSame |= test.Encode(width)
		}
	}

	return levelOutside, levelSame, levelInside
}

func CreateAdjacentsMap(width, height int) map[uint32][]uint32 {
	result := make(map[uint32][]uint32)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			p := Point{i, j}
			pe := p.Encode(width)
			pu, ps, pd := p.FindAdjacents(width, height)
			result[pe] = []uint32{pu, ps, pd}
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

func step(state map[int]uint32, width, height int, adjacentMap map[uint32][]uint32) map[int]uint32 {
	newstate := make(map[int]uint32)
	lmax := 0
	lmin := 0

	// Find populated levels
	for i, v := range state {
		if v != 0 {
			if i < lmin {
				lmin = i
			}
			if i > lmax {
				lmax = i
			}
		}
	}

	centerPoint := height/2*width + width/2

	for level := lmin - 1; level <= lmax+1; level++ {
		newstate[level] = state[level]
		for i := 0; i < width*height; i++ {
			if i == centerPoint {
				continue
			}
			var pos uint32 = 1 << i
			bug := state[level] & pos
			au := adjacentMap[pos][0]
			as := adjacentMap[pos][1]
			ad := adjacentMap[pos][2]
			adjCount := CountBits(au&state[level-1]) + CountBits(as&state[level]) + CountBits(ad&state[level+1])

			if bug == 0 && (adjCount == 1 || adjCount == 2) {
				newstate[level] |= pos
			} else if bug > 0 && adjCount != 1 {
				newstate[level] &^= pos
			}
		}
	}

	return newstate
}

func drawLevels(state map[int]uint32, width, height int) {
	lmax := 0
	lmin := 0

	// Find populated levels
	for i, v := range state {
		if v != 0 {
			if i < lmin {
				lmin = i
			}
			if i > lmax {
				lmax = i
			}
		}
	}

	bugcount := 0

	for level := lmin; level <= lmax; level++ {
		fmt.Printf("Depth: %d\n", level)
		for j := 0; j < height; j++ {
			for i := 0; i < width; i++ {
				if i == width/2 && j == height/2 {
					fmt.Printf("?")
					continue
				}
				var pos uint32 = 1 << uint32(j*width+i)
				if state[level]&pos > 0 {
					fmt.Printf("#")
					bugcount++
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	fmt.Printf("%d bugs found\n", bugcount)
}
