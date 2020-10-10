package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	"github.com/Pauloo27/normigo/ocr"
	"github.com/Pauloo27/normigo/reddit"
	"github.com/Pauloo27/normigo/utils"
	"github.com/fogleman/gg"
)

func loadImage(path string) image.Image {
	imgFile, err := os.Open(path)
	utils.HandleError(err, "Cannot open file")

	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	utils.HandleError(err, "Cannot decode image")

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

func drawCaption(canvas *image.RGBA, text string, height int, fontH float64) *gg.Context {
	width := canvas.Bounds().Max.X
	dc := gg.NewContextForRGBA(canvas)

	if fontH == 0 {
		textLen := len(text)

		fontH = (float64(width) * float64(height)) / (float64(textLen) * 45.0)
		fmt.Println("Current font size:", fontH)
		fmt.Println("If the font size sucks, pass a better one after the caption parameter.")
	}

	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/font.ttf", fontH); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(strings.ToUpper(text), 5, 15, 0, 0, float64(width-5), 1.35, gg.AlignLeft)

	return dc
}

func printUsage() {
	fmt.Println("Usage: normigo <reddit url>")
	fmt.Println(`Usage: normigo <image path> "<caption>" [font size (optional)]`)
}

func fromImage() {
	fontSize := 0.0
	if len(os.Args) != 3 {
		if len(os.Args) == 4 {
			value, err := strconv.ParseFloat(os.Args[3], 64)
			utils.HandleError(err, "Cannot parse font size")
			fontSize = value
		} else {
			printUsage()
			os.Exit(-1)
		}
	}

	path := os.Args[1]

	caption := os.Args[2]

	img := loadImage(path)

	height, color := captionBox(img)

	canvas := removeCaption(img, height, color)

	dc := drawCaption(canvas, caption, height, fontSize)

	err := dc.SavePNG("meme.png")
	utils.HandleError(err, "Cannot save output file")
}

func fromReddit() {
	postURL := os.Args[1]
	imageURL := reddit.GetImageURL(postURL)
	text := ocr.GetTextFromImageURL(imageURL)
	fmt.Println(text)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(-1)
	}
	if len(os.Args) == 2 {
		fromReddit()
	} else {
		fromImage()
	}
}
