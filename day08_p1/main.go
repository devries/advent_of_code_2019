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

	imageData := string(content)
	imageData = strings.TrimSpace(imageData)

	image := &SpaceImage{
		Width:  25,
		Height: 6,
		Data:   []rune(imageData),
	}

	minZero := -1
	minCheck := 0
	for l := 0; l < image.NumLayers(); l++ {
		runeCount := make(map[rune]int)
		for y := 0; y < image.Height; y++ {
			for x := 0; x < image.Width; x++ {
				r := image.Get(x, y, l)
				runeCount[r] = runeCount[r] + 1
			}
		}
		if minZero == -1 || runeCount['0'] < minZero {
			minZero = runeCount['0']
			minCheck = runeCount['1'] * runeCount['2']
		}
	}

	fmt.Printf("Check: %d\n", minCheck)
}

type SpaceImage struct {
	Width  int
	Height int
	Data   []rune
}

func (im *SpaceImage) Get(x int, y int, layer int) rune {
	loc := im.Width*im.Height*layer + im.Width*y + x

	return im.Data[loc]
}

func (im *SpaceImage) NumLayers() int {
	res := len(im.Data) / (im.Width * im.Height)

	return res
}
