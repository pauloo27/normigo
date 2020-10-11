package translate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
)

type TranslateResult struct {
	OriginalText, TranslatedText, From, To string
}

func Translate(text, from, to string) (TranslateResult, error) {
	res, err := http.Get(fmt.Sprintf(
		"https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s&ie=UTF-8&oe=UTF-8",
		url.QueryEscape(from),
		url.QueryEscape(to),
		url.QueryEscape(text),
	))

	if err != nil {
		return TranslateResult{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TranslateResult{}, err
	}

	defer res.Body.Close()

	finalText := ""

	data, _, _, err := jsonparser.Get(body, "[0]")
	if err != nil {
		return TranslateResult{}, err
	}

	var funcErr error
	_, err = jsonparser.ArrayEach(data, func(sentence []byte, t jsonparser.ValueType, i int, err error) {
		if err != nil {
			funcErr = err
			return
		}
		text, _, _, err := jsonparser.Get(sentence, "[0]")
		if err != nil {
			funcErr = err
			return
		}
		finalText += string(text)
	})

	if funcErr != nil {
		return TranslateResult{}, funcErr
	}

	if err != nil {
		return TranslateResult{}, err
	}

	return TranslateResult{OriginalText: text, From: from, To: to, TranslatedText: finalText}, nil
}
