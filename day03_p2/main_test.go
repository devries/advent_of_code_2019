package main

import (
	"sort"
	"testing"
)

func TestSolution(t *testing.T) {
	var tests = []struct {
		pathA    string
		pathB    string
		distance int
	}{
		{"R75,D30,R83,U83,L12,D49,R71,U7,L72", "U62,R66,U55,R34,D71,R55,D58,R83", 610},
		{"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51", "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7", 410},
	}

	for _, test := range tests {
		grid := make(map[Point]map[string]int)

		if err := ParsePath(test.pathA, "A", grid); err != nil {
			t.Errorf("Error parsing input A: %s", err)
			continue
		}

		if err := ParsePath(test.pathB, "B", grid); err != nil {
			t.Errorf("Error parsing input B: %s", err)
			continue
		}

		distances := FindTravelDistances(grid)
		if len(distances) == 0 {
			t.Errorf("No distances found")
			continue
		}
		sort.Ints(distances)
		if distances[0] != test.distance {
			t.Errorf("Expected %d got %d", test.distance, distances[0])
		}
	}
}
