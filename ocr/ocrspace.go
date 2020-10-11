package ocr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OCRResult struct {
	ParsedResults []struct {
		ParsedText string
	}
}

func GetTextFromImageURL(imageURL, apiKey string) (string, error) {
	path := fmt.Sprintf("https://api.ocr.space/parse/imageurl?apikey=%sd&url=%s", apiKey, imageURL)

	res, err := http.Get(path)
	if err != nil {
		return "", err
	}

	bodyB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var result OCRResult

	err = json.Unmarshal(bodyB, &result)
	if err != nil {
		return "", err
	}

	return result.ParsedResults[0].ParsedText, nil
}
