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

	steps := multicycle(moons)
	fmt.Printf("%d\n", steps)
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

func multicycle(moons []Planet) int64 {
	p := make([]int, len(moons))
	for i := range moons {
		p[i] = moons[i].P.X
	}
	v := []int{0, 0, 0, 0}

	cx := cycle(p, v)

	for i := range moons {
		p[i] = moons[i].P.Y
	}
	v = []int{0, 0, 0, 0}

	cy := cycle(p, v)

	for i := range moons {
		p[i] = moons[i].P.Z
	}
	v = []int{0, 0, 0, 0}

	cz := cycle(p, v)

	lyz := lcm(cy, cz)
	lxyz := lcm(cx, lyz)

	return lxyz
}

func cycle(positions []int, velocities []int) int64 {
	init_positions := make([]int, len(positions))
	copy(init_positions, positions)

	init_velocities := make([]int, len(velocities))
	copy(init_velocities, velocities)

	steps := int64(0)
	for {
		for i := range positions {
			for j := i + 1; j < len(positions); j++ {
				if positions[i] > positions[j] {
					velocities[i]--
					velocities[j]++
				} else if positions[i] < positions[j] {
					velocities[i]++
					velocities[j]--
				}
			}
		}
		for i := range positions {
			positions[i] += velocities[i]
		}
		steps++
		if CompareArrays(init_positions, positions) && CompareArrays(init_velocities, velocities) {
			break
		}
	}

	return steps
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

func MinElement(a []int64) int {
	m := int64(0)
	e := 0
	for i := range a {
		if i == 0 || a[i] < m {
			m = a[i]
			e = i
		}
	}

	return e
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

func (v Vector) Equals(a Vector) bool {
	if v.X != a.X {
		return false
	}
	if v.Y != a.Y {
		return false
	}
	if v.Z != a.Z {
		return false
	}

	return true
}

func (p Planet) Equals(a Planet) bool {
	if !p.P.Equals(a.P) {
		return false
	}
	if !p.V.Equals(a.V) {
		return false
	}

	return true
}

func CopyPlanets(moons []Planet) []Planet {
	var r []Planet

	for i := range moons {
		r = append(r, moons[i])
	}

	return r
}

func CompareArrays(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

/*
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	} else {
		return gcd(b, a%b)
	}
}
*/

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}
