package translate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Pauloo27/normigo/utils"
	"github.com/buger/jsonparser"
)

type TranslateResult struct {
	OriginalText, TranslatedText, From, To string
}

func Translate(text, from, to string) TranslateResult {
	res, err := http.Get(fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s&ie=UTF-8&oe=UTF-8", url.QueryEscape(from), url.QueryEscape(to), url.QueryEscape(text)))

	utils.HandleError(err, "Cannot do get request")

	body, err := ioutil.ReadAll(res.Body)
	utils.HandleError(err, "Cannot read body")

	defer res.Body.Close()

	finalText := ""

	data, _, _, err := jsonparser.Get(body, "[0]")
	utils.HandleError(err, "Cannot parse json")

	_, err = jsonparser.ArrayEach(data, func(sentence []byte, t jsonparser.ValueType, i int, err error) {
		utils.HandleError(err, "Cannot parse json")
		text, _, _, err := jsonparser.Get(sentence, "[0]")
		utils.HandleError(err, "Cannot parse json")
		finalText += string(text)
	})

	utils.HandleError(err, "Cannot parse json")

	return TranslateResult{OriginalText: text, From: from, To: to, TranslatedText: finalText}
}
