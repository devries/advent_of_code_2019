package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}

	positions, width, height := ParseMap(string(content))

	fmt.Printf("Asteroid Count: %d\n", len(positions))

	hiddens := HiddenAsteroidDetector(positions, width, height)

	min := len(positions)
	var minPt Point
	for pt, unseen := range hiddens {
		if len(unseen) < min {
			min = len(unseen)
			minPt = pt
		}
	}
	seen := len(positions) - len(hiddens[minPt]) - 1 // Subtract one for yourself
	fmt.Printf("Best position: %v\n", minPt)
	fmt.Printf("Seen asteroids: %d\n", seen)

	counter := 1
	for {
		var vaporizers []PointAngle
		var nextPositions []Point

		for _, v := range positions {
			if hiddens[minPt][v] == false && !minPt.Equals(v) {
				offset := minPt.Offset(v)
				angleMath := math.Atan2(float64(offset.Y), float64(offset.X))
				angleProblem := math.Pi/2.0 + angleMath
				if angleProblem < 0.0 {
					angleProblem = 5.0*math.Pi/2.0 + angleMath
				}

				vaporizers = append(vaporizers, PointAngle{angleProblem, v})
			} else {
				nextPositions = append(nextPositions, v)
			}
		}

		sort.Sort(ByAngle(vaporizers))

		for _, v := range vaporizers {
			fmt.Printf("%d: Asteroid %v, Angle: %f\n", counter, v.Target, v.Angle)
			counter++
		}

		positions = nextPositions
		if len(positions) <= 1 {
			break
		}
		hiddens = HiddenAsteroidDetector(positions, width, height)
	}
}

type PointAngle struct {
	Angle  float64
	Target Point
}

type ByAngle []PointAngle

func (a ByAngle) Len() int           { return len(a) }
func (a ByAngle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAngle) Less(i, j int) bool { return a[i].Angle < a[j].Angle }

func ParseMap(input string) ([]Point, int, int) {
	var ret []Point
	x := 0
	y := 0
	var width int
	var height int
	inputRunes := []rune(input)
	for _, v := range inputRunes {
		switch v {
		case '#':
			ret = append(ret, Point{x, y})
			x++
		case '.':
			x++
		case '\n':
			width = x
			y++
			x = 0
		}
	}

	if inputRunes[len(inputRunes)-1] != '\n' {
		y++
	}

	height = y
	return ret, width, height
}

func HiddenAsteroidDetector(positions []Point, width int, height int) map[Point]map[Point]bool {
	hiddens := make(map[Point]map[Point]bool)

	for _, v := range positions {
		hiddens[v] = make(map[Point]bool)
	}

	for i, origin := range positions {
		for j := i + 1; j < len(positions); j++ {
			destination := positions[j]
			if hiddens[origin][destination] == true {
				continue
			}

			offset := origin.Offset(destination)
			min := offset.AbsMin()

			var slope Point
			var factor int
			for gcf := min; gcf > 0; gcf-- {
				slope = offset.Factor(gcf)
				if !slope.Equals(Point{0, 0}) {
					factor = gcf
					break
				}
			}

			for scale := factor + 1; ; scale++ {
				test := origin.Add(slope.Scale(scale))
				if test.OutOfBounds(width, height) {
					break
				}
				if hiddens[test] != nil {
					hiddens[origin][test] = true
					hiddens[test][origin] = true
				}
			}
		}
	}

	return hiddens
}

type Point struct {
	X int
	Y int
}

func (pt Point) Offset(pt2 Point) Point {
	off := Point{pt2.X - pt.X, pt2.Y - pt.Y}

	return off
}

func (pt Point) Scale(m int) Point {
	r := Point{m * pt.X, m * pt.Y}

	return r
}

func (pt Point) Add(pt2 Point) Point {
	n := Point{pt.X + pt2.X, pt.Y + pt2.Y}

	return n
}

func (pt Point) OutOfBounds(width int, height int) bool {
	if pt.X < 0 || pt.X >= width {
		return true
	} else if pt.Y < 0 || pt.Y >= height {
		return true
	} else {
		return false
	}
}

func (pt Point) AbsMin() int {
	if Abs(pt.X) < Abs(pt.Y) {
		if Abs(pt.X) == 0 {
			return Abs(pt.Y)
		} else {
			return Abs(pt.X)
		}
	} else {
		if Abs(pt.Y) == 0 {
			return Abs(pt.X)
		} else {
			return Abs(pt.Y)
		}
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
func (pt Point) Factor(f int) Point {
	if pt.X%f == 0 && pt.Y%f == 0 {
		return Point{pt.X / f, pt.Y / f}
	} else {
		return Point{0, 0}
	}
}

func (pt Point) Equals(pt2 Point) bool {
	if pt.X == pt2.X && pt.Y == pt2.Y {
		return true
	} else {
		return false
	}
}
