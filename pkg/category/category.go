package category

import (
	"net/http"

	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/scraper"
	"github.com/n0madic/google-play-scraper/pkg/store"
)

// Options type alias
type Options = scraper.Options

// List of clusters
type List map[string]*scraper.Scraper

// New return category list instance
func New(category store.Category, sort store.Sort, age store.Age, options Options) (List, error) {
	path := ""
	switch sort {
	case store.SortRating:
		path = "/top"
	case store.SortNewest:
		path = "/new"
	}
	if category != "" {
		path += "/category/" + string(category)
	}

	req, err := http.NewRequest("GET", scraper.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	if age != "" {
		q := req.URL.Query()
		q.Add("age", string(age))
		req.URL.RawQuery = q.Encode()
	}

	data, err := util.GetInitData(req)
	if err != nil {
		return nil, err
	}

	list := make(List)

	clusterList := util.GetJSONArray(data["ds:3"], "0.1")
	for _, cluster := range clusterList {
		key := util.GetJSONValue(cluster.String(), "0.1")
		url, err := util.AbsoluteURL(scraper.BaseURL, util.GetJSONValue(cluster.String(), "0.3.4.2"))
		if key == "" {
			key = util.GetJSONValue(cluster.String(), "20.0")
			url, err = util.AbsoluteURL(scraper.BaseURL, util.GetJSONValue(cluster.String(), "20.2.4.2"))
		}
		if key != "" && err == nil {
			list[key] = scraper.New(url, &options)
		}
	}

	return list, nil
}
