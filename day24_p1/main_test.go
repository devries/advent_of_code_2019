package main

import (
	"strings"
	"testing"
)

const seenTwice string = `.....
.....
.....
#....
.#...`

const initialState string = `....#
#..#.
#..##
..#..
#....`

const stepOne string = `#..#.
####.
###.#
##.##
.##..`

func TestEncoding(t *testing.T) {
	tests := []struct {
		input  string
		result uint32
	}{
		{seenTwice, 2129920},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		encoded := parseInput(r, 5)

		if encoded != test.result {
			t.Errorf("Got %d, expected %d", encoded, test.result)
		}
	}
}

func TestCounting(t *testing.T) {
	tests := []struct {
		input  string
		result int
	}{
		{seenTwice, 2},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		encoded := parseInput(r, 5)

		bitcount := CountBits(encoded)

		if bitcount != test.result {
			t.Errorf("Got %d, expected %d", bitcount, test.result)
		}
	}
}

func TestStep(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{initialState, stepOne},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		encodedIn := parseInput(r, 5)

		r2 := strings.NewReader(test.output)
		encodedOut := parseInput(r2, 5)

		m := CreateAdjacentsMap(5, 5)

		testOut := step(encodedIn, 5, 5, m)

		if testOut != encodedOut {
			t.Errorf("Got %d, expected %d", testOut, encodedOut)
		}
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{initialState, seenTwice},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		encodedIn := parseInput(r, 5)

		r2 := strings.NewReader(test.output)
		encodedOut := parseInput(r2, 5)

		_, finalState := findRepeat(encodedIn, 5, 5)

		if finalState != encodedOut {
			t.Errorf("Got %d, expected %d", finalState, encodedOut)
		}
	}
}
