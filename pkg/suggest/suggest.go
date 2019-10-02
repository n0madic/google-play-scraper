package suggest

import (
	"encoding/json"
	"net/http"

	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

const suggURL = "https://market.android.com/suggest/SuggRequest"

// Options type alias
type Options = app.Options

// Get returns up to five suggestion to complete a search query term
func Get(term string, options Options) (list []string, err error) {
	req, err := http.NewRequest("GET", suggURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("json", "1")
	q.Set("query", term)
	q.Add("gl", options.Country)
	q.Add("hl", options.Language)
	req.URL.RawQuery = q.Encode()

	body, err := util.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var jslist []struct {
		S string `json:"s"`
	}
	err = json.Unmarshal(body, &jslist)
	if err != nil {
		return nil, err
	}

	for _, sugg := range jslist {
		if sugg.S != "" {
			list = append(list, sugg.S)
		}
	}

	return list, nil
}
