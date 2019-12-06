package main

import (
	"strings"
	"testing"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		input    string
		checksum int
	}{
		{"COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L", 42},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		om, err := NewOrbitMap(r)
		if err != nil {
			t.Errorf("Error reading input: %s", err)
			continue
		}
		checksum := om.Checksum()
		if checksum != test.checksum {
			t.Errorf("Got %d, expected %d", checksum, test.checksum)
		}
	}
}
