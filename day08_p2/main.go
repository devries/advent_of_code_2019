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

	imageData := string(content)
	imageData = strings.TrimSpace(imageData)

	im := &SpaceImage{
		Width:  25,
		Height: 6,
		Data:   []rune(imageData),
	}

	// Make drawable image
	zoom := 10
	ul := image.Point{0, 0}
	br := image.Point{im.Width * zoom, im.Height * zoom}
	outputImg := image.NewRGBA(image.Rectangle{ul, br})

	for y := 0; y < im.Height; y++ {
		for x := 0; x < im.Width; x++ {
			var pixel color.Color
			for l := 0; l < im.NumLayers(); l++ {
				r := im.Get(x, y, l)
				if r == '0' {
					pixel = color.Black
					break
				} else if r == '1' {
					pixel = color.White
					break
				} else {
					continue
				}
			}
			for j := y * zoom; j < (y+1)*zoom; j++ {
				for i := x * zoom; i < (x+1)*zoom; i++ {
					outputImg.Set(i, j, pixel)
				}
			}

		}
	}

	f, _ := os.Create("image.png")
	defer f.Close()
	png.Encode(f, outputImg)
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
