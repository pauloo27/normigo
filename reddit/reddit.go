package reddit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Pauloo27/normigo/utils"
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

func GetImageURL(postURL string) string {
	client := &http.Client{}

	path := strings.Split(postURL, "?")[0]

	if !strings.HasSuffix(path, ".json") {
		path += ".json"
	}

	req, err := http.NewRequest("GET", path, nil)
	utils.HandleError(err, "Cannot create GET request")

	req.Header.Add("User-Agent", `NormIGo`)

	res, err := client.Do(req)
	utils.HandleError(err, "Cannot request to "+path)

	bodyB, err := ioutil.ReadAll(res.Body)
	utils.HandleError(err, "Cannot read body")

	defer res.Body.Close()

	var results []RedditPost

	err = json.Unmarshal(bodyB, &results)
	utils.HandleError(err, "Cannot parse json")

	return results[0].Data.Children[0].Data.URL
}
