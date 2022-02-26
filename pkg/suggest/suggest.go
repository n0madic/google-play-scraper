package suggest

import (
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

// Options type alias
type Options = app.Options

// Get returns up to five suggestion to complete a search query term
func Get(term string, options Options) (list []string, err error) {
	payload := "f.req=%5B%5B%5B%22IJ4APc%22%2C%22%5B%5Bnull%2C%5B%5C%22" + term + "%20e%5C%22%5D%2C%5B10%5D%2C%5B2%5D%2C4%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D"

	js, err := util.BatchExecute(options.Country, options.Language, payload)
	if err != nil {
		return nil, err
	}

	for _, sugg := range util.GetJSONArray(js, "0.0.#.0") {
		list = append(list, sugg.String())
	}

	return list, nil
}
