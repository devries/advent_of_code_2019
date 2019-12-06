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

	transfers := om.Transfers("YOU", "SAN")
	transfers-- // Don't go into orbit around Santa
	fmt.Printf("Transfers: %d\n", transfers)
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

func (om OrbitMap) FirstCommonObject(obj1 string, obj2 string) string {
	orbits1 := om.Orbits(obj1)
	orbits2 := om.Orbits(obj2)

	for _, o := range orbits1 {
		if obj2 == o {
			return obj2
		}
	}

	for _, o := range orbits2 {
		if obj1 == o {
			return obj1
		}
	}

	for _, o1 := range orbits1 {
		for _, o2 := range orbits2 {
			if o1 == o2 {
				return o1
			}
		}
	}

	return ""
}

// Number of transfers to go from outer to inner
func (om OrbitMap) Transfers(object string, destination string) int {
	common := om.FirstCommonObject(object, destination)

	descent := om.Orbits(object)
	var down int

	for i, v := range descent {
		if v == common {
			down = i
			break
		}
	}

	ascent := om.Orbits(destination)
	var up int

	for i, v := range ascent {
		if v == common {
			up = i + 1
			break
		}
	}

	return down + up
}
