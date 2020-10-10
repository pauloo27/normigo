package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
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

func drawCaption(canvas *image.RGBA, text string) *gg.Context {
	width := canvas.Bounds().Max.X
	dc := gg.NewContextForRGBA(canvas)
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/font.ttf", 32); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(strings.ToUpper(text), 5, 10, 0, 0, float64(width-5), 1.25, gg.AlignLeft)

	return dc
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Missing parameters.\nUsage: normigo <src> <caption>")
		os.Exit(-1)
	}

	path := os.Args[1]

	caption := strings.Join(os.Args[2:], " ")

	img := loadImage(path)

	height, color := captionBox(img)

	canvas := removeCaption(img, height, color)

	dc := drawCaption(canvas, caption)

	err := dc.SavePNG("meme.png")
	handleError(err)
}
