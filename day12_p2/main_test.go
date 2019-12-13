package main

import (
	"fmt"
	"strings"
	"testing"
)

const inputA string = `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`

const inputB string = `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`

func TestExecution(t *testing.T) {
	tests := []struct {
		input string
		steps int64
	}{
		{inputA, 2772},
		{inputB, 4686774924},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		moons, err := parseInput(r)
		if err != nil {
			t.Errorf("Error reading input: %s", err)
		}

		cx := multicycle(moons)
		fmt.Println(cx)
		steps := cx
		if steps != test.steps {
			t.Errorf("Got %d steps, expected %d", steps, test.steps)
		}
	}
}
