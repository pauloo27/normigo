package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Pauloo27/normigo/utils"
)

type OCRResult struct {
	ParsedResults []struct {
		ParsedText string
	}
}

func GetTextFromImageURL(imageURL, apiKey string) string {
	path := fmt.Sprintf("https://api.ocr.space/parse/imageurl?apikey=%sd&url=%s", apiKey, imageURL)

	res, err := http.Get(path)
	utils.HandleError(err, "Cannot get "+path)

	bodyB, err := ioutil.ReadAll(res.Body)

	utils.HandleError(err, "Cannot read body")

	defer res.Body.Close()

	var result OCRResult

	err = json.Unmarshal(bodyB, &result)
	utils.HandleError(err, "Cannot parser json")

	return result.ParsedResults[0].ParsedText
}
