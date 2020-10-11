package reddit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type RedditPost struct {
	Data struct {
		Children []struct {
			Data struct {
				URL string `json:"url_overridden_by_dest"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func GetImageURL(postURL string) (string, error) {
	client := &http.Client{}

	path := strings.Split(postURL, "?")[0]

	if !strings.HasSuffix(path, ".json") {
		path += ".json"
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", `NormIGo`)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var results []RedditPost

	err = json.Unmarshal(bodyB, &results)
	if err != nil {
		return "", err
	}

	return results[0].Data.Children[0].Data.URL, nil
}
