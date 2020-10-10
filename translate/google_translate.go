package translate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Pauloo27/normigo/utils"
)

type TranslateResult struct {
	OriginalText, TranslatedText, From, To string
}

func Translate(text, from, to string) (result TranslateResult, err error) {
	res, err := http.Get(fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s&ie=UTF-8&oe=UTF-8", url.QueryEscape(from), url.QueryEscape(to), url.QueryEscape(text)))
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	jsonData := string(body)
	var parsedJSON [][][]interface{}
	err = json.Unmarshal([]byte(jsonData), &parsedJSON)
	utils.HandleError(err, "Cannot parse json")

	if from == "auto" {
		from = parsedJSON[8][0][0].(string)
	}

	phrases := parsedJSON[0]
	sentence := ""

	for _, item := range phrases {
		sentence = fmt.Sprintf("%s%s", sentence, item[0].(string))
	}

	result = TranslateResult{OriginalText: text, TranslatedText: sentence, From: from, To: to}
	return
}
