package collection

import (
	"github.com/realchandan/google-play-scraper/pkg/scraper"
	"github.com/realchandan/google-play-scraper/pkg/store"
)

// Options type alias
type Options = scraper.Options

// New return collection list instance
func New(collection store.Collection, options Options) *scraper.Scraper {
	return scraper.New(scraper.BaseURL+"/collection/"+string(collection), &options)
}
