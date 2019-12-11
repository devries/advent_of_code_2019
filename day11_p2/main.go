package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
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
	panels[Point{0, 0}] = 1

	r.Run(startingOpcodes, panels)

	outputImg := PanelImage(panels)
	f, _ := os.Create("image.png")
	defer f.Close()
	png.Encode(f, outputImg)
}

func PanelImage(panels map[Point]int64) image.Image {
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for k, _ := range panels {
		if k.X < minX {
			minX = k.X
		}
		if k.X > maxX {
			maxX = k.X
		}
		if k.Y < minY {
			minY = k.Y
		}
		if k.Y > maxY {
			maxY = k.Y
		}
	}

	ul := image.Point{minX, minY}
	br := image.Point{maxX + 1, maxY + 1}
	outputImg := image.NewRGBA(image.Rectangle{ul, br})

	colorMap := map[int64]color.Color{
		0: color.Black,
		1: color.White,
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			c := panels[Point{i, j}]
			outputImg.Set(i, j, colorMap[c])
		}
	}

	return outputImg
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
