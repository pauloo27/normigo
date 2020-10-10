package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadImage(path string) image.Image {
	imgFile, err := os.Open(path)
	handleError(err)

	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	handleError(err)

	return img
}

func getColorDiff(a color.Color, b color.Color) (diff uint32) {
	diff = 0

	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()

	check := func(x, y uint32) {
		if x >= y {
			diff += x - y
		} else {
			diff += y - x
		}
	}

	check(aR, bR)
	check(aG, bG)
	check(aB, bB)

	return
}

func captionBox(img image.Image) (y int, boxColor color.Color) {
	max := img.Bounds().Max.Y

	boxColor = img.At(0, 0)

	for y = 0; y < max; y++ {
		color := img.At(0, y)
		if color != boxColor {
			if getColorDiff(color, boxColor) >= 8000 {
				break
			}
		}
	}

	return
}

func removeCaption(img image.Image, height int, color color.Color) *image.RGBA {
	width := img.Bounds().Max.X
	canvas := image.NewRGBA(img.Bounds())

	draw.Draw(canvas, canvas.Bounds(), img, image.Point{0, 0}, draw.Src)
	draw.Draw(canvas, image.Rect(0, 0, width, height), &image.Uniform{color}, image.Point{0, 0}, draw.Src)

	return canvas
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing meme path parameter")
		os.Exit(-1)
	}
	path := os.Args[1]
	img := loadImage(path)
	height, color := captionBox(img)
	canvas := removeCaption(img, height, color)

	outputFile, err := os.Create("meme.jpg")
	handleError(err)

	err = jpeg.Encode(outputFile, canvas, nil)
	handleError(err)
}
