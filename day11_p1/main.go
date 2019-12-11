package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %s", err))
	}

	program := string(content)
	program = strings.TrimSpace(program)
	startingOpcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	r := NewRobot()
	panels := make(map[Point]int64)

	r.Run(startingOpcodes, panels)

	fmt.Printf("Panels painted at least once: %d\n", len(panels))
}

type Point struct {
	X int
	Y int
}

type Robot struct {
	Position  Point
	Direction Point
}

func NewRobot() *Robot {
	r := Robot{Point{0, 0}, Point{0, -1}}

	return &r
}

func (r *Robot) Left() {
	r.Direction = Point{r.Direction.Y, -r.Direction.X}
}

func (r *Robot) Right() {
	r.Direction = Point{-r.Direction.Y, r.Direction.X}
}

func (r *Robot) Forward() {
	r.Position.X += r.Direction.X
	r.Position.Y += r.Direction.Y
}

func (r *Robot) Run(program map[int64]int64, panels map[Point]int64) {
	input := make(chan int64)
	output := make(chan int64)
	done := make(chan bool)

	go func() {
		if err := ExecuteProgram(program, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
		done <- true
	}()

	for {
		select {
		case input <- panels[r.Position]:
			color := <-output
			turn := <-output

			panels[r.Position] = color
			switch turn {
			case 0:
				r.Left()
			case 1:
				r.Right()
			default:
				panic(fmt.Errorf("Unexpected turn output: %d", turn))
			}
			r.Forward()
		case <-done:
			return
		}
	}
}
