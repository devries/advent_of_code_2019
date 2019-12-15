package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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

	output := make(chan int64)
	input := make(chan int64)
	screen := make(map[Point]int64)

	go func() {
		if err := ExecuteProgram(startingOpcodes, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	score := int64(0)
	paddle := Point{0, 0}
	ball := Point{0, 0}
	prev_ball := Point{0, 0}

mainloop:
	for {
		pdl := predictTrajectory(prev_ball, ball, paddle, screen)
		motion := pdl - paddle.X
		if motion < 0 {
			motion = -1
		} else if motion > 0 {
			motion = 1
		}

		var x int64
		var more bool
		loopit := true

		// This is where I got stuck, how to check if input is wanted.
		// I used select to check if input is desired, but I also need to loop
		// so I get the output.
		for loopit {
			select {
			case x, more = <-output:
				if more == false {
					break mainloop
				}
				loopit = false
			case input <- motion:
				// fmt.Printf("Joystick: %d\n", motion)
				fmt.Printf("\033[3J\033[H")
				printScreen(screen)
				fmt.Printf("\n%d\n", score)
				time.Sleep(50 * time.Millisecond)
			}
		}

		var y int64
		y = <-output

		var tileID int64
		tileID = <-output

		if x == -1 && y == 0 {
			score = tileID
			if CountBlocks(screen) == 0 {
				break mainloop
			}
			//fmt.Printf("\033[3J\033[H")
			//printScreen(screen)
		} else {
			screen[Point{x, y}] = tileID
		}
		if tileID == int64(3) {
			paddle = Point{x, y}
		} else if tileID == int64(4) {
			prev_ball = ball
			ball = Point{x, y}
		}
	}
	fmt.Printf("\033[3J\033[H")
	printScreen(screen)
	fmt.Printf("\n%d\n", score)

	// fmt.Println(score)
}

type Point struct {
	X int64
	Y int64
}

func CountBlocks(screen map[Point]int64) int {
	blockCount := 0
	for _, v := range screen {
		if v == int64(2) {
			blockCount++
		}
	}

	return blockCount
}

func predictTrajectory(prev_ball Point, ball Point, paddle Point, screen map[Point]int64) int64 {
	// fmt.Printf("Sending paddle to %d\n", ball.X)
	return ball.X
}

func printScreen(screen map[Point]int64) {
	minY := int64(0)
	maxY := int64(0)
	minX := int64(0)
	maxX := int64(0)

	for k, _ := range screen {
		if k.X < minX {
			minX = k.X
		} else if k.X > maxX {
			maxX = k.X
		}

		if k.Y < minY {
			minY = k.Y
		} else if k.Y > maxY {
			maxY = k.Y
		}
	}
	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			p := screen[Point{i, j}]
			switch p {
			case 0:
				fmt.Printf(" ")
			case 1:
				fmt.Printf("█")
			case 2:
				fmt.Printf("█")
			case 3:
				fmt.Printf("—")
			case 4:
				fmt.Printf("●")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
