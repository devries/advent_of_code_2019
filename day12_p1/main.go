package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Errorf("Unable to open input file"))
	}
	defer f.Close()

	moons, err := parseInput(f)
	if err != nil {
		panic(fmt.Errorf("Unable to parse input: %s\n", err))
	}

	for i := 0; i < 1000; i++ {
		Step(moons)
	}

	e := TotalEnergy(moons)
	fmt.Printf("%d\n", e)
}

type Planet struct {
	P Vector
	V Vector
}

type Vector struct {
	X int
	Y int
	Z int
}

func parseInput(r io.Reader) ([]Planet, error) {
	var result []Planet

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		p := Planet{}
		n, err := fmt.Sscanf(line, "<x=%d, y=%d, z=%d>", &p.P.X, &p.P.Y, &p.P.Z)

		if n != 3 {
			return nil, fmt.Errorf("Did not parse enough coordinates in line %s", line)
		} else if err != nil {
			return nil, fmt.Errorf("Error parsing input: %s", err)
		}

		result = append(result, p)
	}

	return result, nil
}

func Step(moons []Planet) {
	for i := range moons {
		for j := i + 1; j < len(moons); j++ {
			if moons[i].P.X > moons[j].P.X {
				moons[i].V.X--
				moons[j].V.X++
			} else if moons[i].P.X < moons[j].P.X {
				moons[i].V.X++
				moons[j].V.X--
			}

			if moons[i].P.Y > moons[j].P.Y {
				moons[i].V.Y--
				moons[j].V.Y++
			} else if moons[i].P.Y < moons[j].P.Y {
				moons[i].V.Y++
				moons[j].V.Y--
			}

			if moons[i].P.Z > moons[j].P.Z {
				moons[i].V.Z--
				moons[j].V.Z++
			} else if moons[i].P.Z < moons[j].P.Z {
				moons[i].V.Z++
				moons[j].V.Z--
			}
		}
	}

	for i := range moons {
		moons[i].P.X += moons[i].V.X
		moons[i].P.Y += moons[i].V.Y
		moons[i].P.Z += moons[i].V.Z
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func KineticEnergy(p Planet) int {
	r := Abs(p.V.X) + Abs(p.V.Y) + Abs(p.V.Z)

	return r
}

func PotentialEnergy(p Planet) int {
	r := Abs(p.P.X) + Abs(p.P.Y) + Abs(p.P.Z)

	return r
}

func TotalEnergy(moons []Planet) int {
	r := 0

	for i := range moons {
		r += KineticEnergy(moons[i]) * PotentialEnergy(moons[i])
	}

	return r
}
