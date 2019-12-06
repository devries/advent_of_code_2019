package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	om, err := NewOrbitMap(f)
	if err != nil {
		panic(fmt.Errorf("Unable to parse input: %s", err))
	}

	checksum := om.Checksum()
	fmt.Printf("Checksum: %d\n", checksum)
}

type OrbitMap map[string]string

func NewOrbitMap(r io.Reader) (OrbitMap, error) {
	om := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		segments := strings.Split(line, ")")

		if len(segments) != 2 {
			return nil, fmt.Errorf("Unable to parse line %s", line)
		}

		om[segments[1]] = segments[0]
	}

	return om, nil
}

// Orbits returns the objects object is orbiting
func (om OrbitMap) Orbits(object string) []string {
	var r []string

	for {
		object = om[object]
		if object == "" {
			break
		}
		r = append(r, object)
	}

	return r
}

func (om OrbitMap) Checksum() int {
	total := 0

	for k, _ := range om {
		total += len(om.Orbits(k))
	}

	return total
}
