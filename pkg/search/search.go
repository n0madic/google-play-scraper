package search

import (
	"net/url"
	"strconv"

	"github.com/realchandan/google-play-scraper/pkg/scraper"
)

const (
	searchURL = "https://play.google.com/store/search"
)

// PriceQuery value
type PriceQuery int

// Options type alias
type Options = scraper.Options

const (
	// PriceAll - all prices
	PriceAll PriceQuery = iota
	// PriceFree - only free
	PriceFree
	// PricePaid - only paid
	PricePaid
)

// NewQuery return Query instance
func NewQuery(query string, price PriceQuery, options Options) *scraper.Scraper {
	baseURL, err := url.Parse(searchURL)
	if err != nil {
		return nil
	}

	// Query params
	params := url.Values{}
	params.Add("q", url.QueryEscape(query))
	params.Add("c", "apps")
	params.Add("fpr", "false")
	params.Add("price", strconv.Itoa(int(price)))
	baseURL.RawQuery = params.Encode()

	return scraper.New(baseURL.String(), &options)
}
