package main

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"net/http"
	"os"
	"strings"

	"github.com/Pauloo27/normigo/ocr"
	"github.com/Pauloo27/normigo/reddit"
	"github.com/Pauloo27/normigo/translate"
	"github.com/Pauloo27/normigo/utils"
	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
)

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

func applyTranslation(img image.Image, caption string, fontSize float64) {
	height, color := captionBox(img)

	canvas := removeCaption(img, height, color)

	dc := drawCaption(canvas, caption, height, fontSize)

	err := dc.SavePNG("www/assets/meme.png")
	utils.HandleError(err, "Cannot save output file")
}

type FillPage struct {
	ImageURL string
	Err      error
}

type GeneratePage struct {
	OriginalImageURL, ImageURL, Caption string
	Err                                 error
}

func main() {
	err := godotenv.Load()
	utils.HandleError(err, "Cannot load .env")

	apiKey := os.Getenv("OCR_APIKEY")

	templates := template.Must(template.ParseFiles(
		"www/main.html",
		"www/home.html",
		"www/fill.html",
		"www/generate.html",
	))

	homeTemplate := templates.Lookup("home.html")
	fillTemplate := templates.Lookup("fill.html")
	generateTemplate := templates.Lookup("generate.html")

	fs := http.FileServer(http.Dir("www/assets/"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/generate/", func(w http.ResponseWriter, r *http.Request) {
		var finalErr error

		originalImageURL := r.FormValue("url")
		caption := strings.ReplaceAll(r.FormValue("caption"), "\\n", "\n")

		res, err := http.Get(originalImageURL)
		if err != nil {
			finalErr = err
		} else {
			img, _, err := image.Decode(res.Body)
			if err != nil {
				finalErr = err
			}

			defer res.Body.Close()

			applyTranslation(img, caption, 0)
		}
		imageURL := "/assets/meme.png"

		err = generateTemplate.Execute(w, GeneratePage{originalImageURL, imageURL, caption, finalErr})
		utils.HandleError(err, "Cannot run template")
	})

	http.HandleFunc("/tr/", func(w http.ResponseWriter, r *http.Request) {
		text := r.FormValue("text")
		result, err := translate.Translate(text, "en", "pt")

		if err != nil {
			fmt.Fprintln(w, "error:", err)
		}
		fmt.Fprintln(w, result.TranslatedText)
	})

	http.HandleFunc("/ocr/", func(w http.ResponseWriter, r *http.Request) {
		imageURL := r.FormValue("url")
		text, err := ocr.GetTextFromImageURL(imageURL, apiKey)
		if err != nil {
			fmt.Fprintln(w, "error:", err)
		}
		fmt.Fprintln(w, text)
	})

	http.HandleFunc("/fill/", func(w http.ResponseWriter, r *http.Request) {
		imageURL, err := reddit.GetImageURL(r.FormValue("url"))
		err = fillTemplate.Execute(w, FillPage{imageURL, err})
		utils.HandleError(err, "Cannot run template")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := homeTemplate.Execute(w, nil)
		utils.HandleError(err, "Cannot run template")
	})

	fmt.Println("Running server at localhost:25555")
	err = http.ListenAndServe(":25555", nil)
	utils.HandleError(err, "Cannot start server")
}
