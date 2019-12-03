package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %s", err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	wirePaths := make([]string, 0, 2)
	for scanner.Scan() {
		wirePaths = append(wirePaths, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}

	grid := make(map[Point]string)

	for i, path := range wirePaths {
		identifier := strconv.Itoa(i)
		if err := ParsePath(path, identifier, grid); err != nil {
			panic(fmt.Errorf("Error parsing path %s: %s", identifier, err))
		}
	}

	distances := FindIntersectionDistances(grid)
	if len(distances) == 0 {
		panic(fmt.Errorf("No intersections found!"))
	}
	sort.Ints(distances)
	fmt.Printf("Smallest distance: %d\n", distances[0])
}

type Point struct {
	X int
	Y int
}

func ParsePath(wirePath string, wireIdentifier string, grid map[Point]string) error {
	pathComponents := strings.Split(wirePath, ",")
	position := Point{0, 0}

	for _, c := range pathComponents {
		componentSlice := []rune(c)
		direction := componentSlice[0]

		distance, err := strconv.Atoi(string(componentSlice[1:]))
		if err != nil {
			return err
		}

		for i := 0; i < distance; i++ {
			switch direction {
			case 'U':
				position.Y += 1
			case 'D':
				position.Y -= 1
			case 'R':
				position.X += 1
			case 'L':
				position.X -= 1
			}

			ptVal := grid[position]
			if ptVal == wireIdentifier {
				continue
			}
			grid[position] = ptVal + wireIdentifier
		}
	}

	return nil
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func FindIntersectionDistances(grid map[Point]string) []int {
	var result []int

	for k, v := range grid {
		if len(v) > 1 {
			d := Abs(k.X) + Abs(k.Y)
			result = append(result, d)
		}
	}

	return result
}
