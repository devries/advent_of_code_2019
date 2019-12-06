package main

import (
	"strings"
	"testing"
)

func TestExecution(t *testing.T) {
	var tests = []struct {
		input     string
		transfers int
	}{
		{"COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L\nK)YOU\nI)SAN", 4},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		om, err := NewOrbitMap(r)
		if err != nil {
			t.Errorf("Error reading input: %s", err)
			continue
		}

		transfers := om.Transfers("YOU", "SAN")
		transfers-- // don't orbit santa
		if transfers != test.transfers {
			t.Errorf("Got %d, expected %d", transfers, test.transfers)
		}
	}
}
