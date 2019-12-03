package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
It is worth a little explanation to describe the fairly complex structure holding
the wire distances here. I have made a

	map[Point]map[string]int

which is a map with Point keys and map values, the value maps have string keys and
integer values. I am storing at each point a travel distance for each wire, indexed
by point in the top map, and by wire identifier in the second level maps, I store the
wire travel distance.

A complication is that I check if a value is already there for that point and wire, as
only the lowest travel distance should be saved. I then look for any intersection by
looking for points with more than two wires.
*/

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

	grid := make(map[Point]map[string]int)

	for i, path := range wirePaths {
		identifier := strconv.Itoa(i)
		if err := ParsePath(path, identifier, grid); err != nil {
			panic(fmt.Errorf("Error parsing path %s: %s", identifier, err))
		}
	}

	distances := FindTravelDistances(grid)
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

func ParsePath(wirePath string, wireIdentifier string, grid map[Point]map[string]int) error {
	pathComponents := strings.Split(wirePath, ",")
	position := Point{0, 0}
	travelDistance := 0

	for _, c := range pathComponents {
		componentSlice := []rune(c)
		direction := componentSlice[0]

		distance, err := strconv.Atoi(string(componentSlice[1:]))
		if err != nil {
			return err
		}

		for i := 0; i < distance; i++ {
			travelDistance += 1

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
			if ptVal == nil {
				grid[position] = make(map[string]int)
				grid[position][wireIdentifier] = travelDistance
			} else {
				if ptVal[wireIdentifier] == 0 {
					ptVal[wireIdentifier] = travelDistance
				}
			}
		}
	}

	return nil
}

func FindTravelDistances(grid map[Point]map[string]int) []int {
	var result []int

	for _, v := range grid {
		if len(v) > 1 {
			totalTravel := 0
			for _, wireDistance := range v {
				totalTravel += wireDistance
			}
			result = append(result, totalTravel)
		}
	}

	return result
}
