package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	// "sort"
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

	opcodes := CopyProgram(startingOpcodes)

	output := make(chan int64)

	go func() {
		if err := ExecuteProgram(opcodes, nil, output); err != nil {
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
	l, d := findRobotAndDirection(grid)
	p := explore(grid, l, d)
	fmt.Println(string(p))

	cs := compactSequence(p)
	a, b, c := findTriplet(cs)
	fmt.Println(a, b, c)

	mainPath := stringifySequence(cs)
	mainPath = strings.ReplaceAll(mainPath, a, "A")
	mainPath = strings.ReplaceAll(mainPath, b, "B")
	mainPath = strings.ReplaceAll(mainPath, c, "C")

	fmt.Println(mainPath)

	toSend := strings.Join([]string{mainPath, a, b, c, "n\n"}, "\n")

	opcodes = CopyProgram(startingOpcodes)
	opcodes[0] = 2

	output = make(chan int64)
	input := make(chan int64)
	done := make(chan bool)
	go func() {
		if err := ExecuteProgram(opcodes, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
		done <- true
	}()

	letterCounter := 0
	nextLetter := int64(toSend[letterCounter])
ioloop:
	for {
		// IO LOOP
		select {
		case c := <-output:
			if c < 256 {
				fmt.Printf("%c", c)
			} else {
				fmt.Printf("%d\n", c)
			}
		case input <- nextLetter:
			if letterCounter < len(toSend)-1 {
				letterCounter++
			}
			fmt.Printf("%c", nextLetter)
			nextLetter = int64(toSend[letterCounter])
		case <-done:
			break ioloop
		}
	}

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

func Left(d Point) Point {
	return Point{d.Y, -d.X}
}

func Right(d Point) Point {
	return Point{-d.Y, d.X}
}

func findRobotAndDirection(grid map[Point]int64) (Point, Point) {
	for k, v := range grid {
		if v == '^' {
			return k, directions['^']
			break
		}
	}
	return Point{0, 0}, Point{0, 0}
}

func explore(grid map[Point]int64, location Point, direction Point) []rune {
	var r []rune

	for {
		check := Point{location.X + direction.X, location.Y + direction.Y}
		if grid[check] != '#' {
			// End of line, check if the path continues left or right
			leftward := Left(direction)
			rightward := Right(direction)
			if grid[Point{location.X + leftward.X, location.Y + leftward.Y}] == '#' {
				direction = leftward
				r = append(r, 'L')
			} else if grid[Point{location.X + rightward.X, location.Y + rightward.Y}] == '#' {
				direction = rightward
				r = append(r, 'R')
			} else {
				// End of line
				break
			}
		} else {
			location = Point{location.X + direction.X, location.Y + direction.Y}
			r = append(r, '1')
		}
	}

	return r
}

func compactSequence(sequence []rune) []string {
	ctr := 0
	var out []string

	for i := 0; i < len(sequence); i++ {
		if sequence[i] <= 'Z' && sequence[i] >= 'A' {
			if ctr > 0 {
				out = append(out, strconv.Itoa(ctr))
				ctr = 0
			}
			out = append(out, string(sequence[i]))
		} else {
			ctr++
		}
	}
	if ctr > 0 {
		out = append(out, strconv.Itoa(ctr))
	}

	return out
}

func stringifySequence(sequence []string) string {
	return strings.Join(sequence, ",")
}

func makeSubsequences(sequences [][]string) map[string]int {
	r := make(map[string]int)

	for _, sequence := range sequences {
		for l := 2; l <= len(sequence); l += 2 {
			s := stringifySequence(sequence[:l])
			if len(s) < 20 {
				if _, ok := r[s]; !ok {
					r[s] = l
				}
			}
		}
	}

	return r
}

func cutSequence(sequence []string, segments [][]string) [][]string {
	var r [][]string

	for _, path := range segments {
		pcut := 0
		for i := 0; i <= len(path)-len(sequence); {
			if EqualSequence(sequence, path[i:i+len(sequence)]) {
				if i > pcut {
					r = append(r, path[pcut:i])
				}
				i += len(sequence)
				pcut = i
			} else {
				i += 2
			}
		}
		if pcut <= len(path)-len(sequence) {
			r = append(r, path[pcut:])
		}
	}

	return r
}

func EqualSequence(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func findTriplet(path []string) (string, string, string) {
	subsequences := makeSubsequences([][]string{path})
	for sa, _ := range subsequences {
		s := cutSequence(strings.Split(sa, ","), [][]string{path})
		subsequencesTwo := makeSubsequences(s)
		for sb, _ := range subsequencesTwo {
			s2 := cutSequence(strings.Split(sb, ","), s)
			sc := s2[0]
			success := true
			for i := 1; i < len(s2); i++ {
				if !EqualSequence(sc, s2[i]) {
					success = false
					break
				}
			}
			if success {
				return sa, sb, strings.Join(sc, ",")
			}
		}
	}

	return "", "", ""
}

func CopyProgram(in map[int64]int64) map[int64]int64 {
	out := make(map[int64]int64)

	for k, v := range in {
		out[k] = v
	}

	return out
}
