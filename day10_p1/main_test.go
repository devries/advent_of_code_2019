package main

import (
	"fmt"
	"testing"
)

const m1 string = `......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`

const m2 string = `#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`

const m3 string = `.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`

const m4 string = `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`

func TestExecution(t *testing.T) {
	var tests = []struct {
		input string
		seen  int
	}{
		{m1, 33},
		{m2, 35},
		{m3, 41},
		{m4, 210},
	}

	for _, test := range tests {
		positions, width, height := ParseMap(test.input)
		hiddens := HiddenAsteroidDetector(positions, width, height)

		min := len(positions)
		var minPt Point
		for pt, unseen := range hiddens {
			if len(unseen) < min {
				min = len(unseen)
				minPt = pt
			}
		}
		seen := len(positions) - len(hiddens[minPt]) - 1
		fmt.Printf("Best point: %d, width: %d, height: %d\n", minPt, width, height)
		for k, _ := range hiddens[minPt] {
			fmt.Printf("%v,", k)
		}
		fmt.Printf("\n")

		if seen != test.seen {
			t.Errorf("Expected %d, got %d", test.seen, seen)
		}
	}
}
