package similar

import (
	"github.com/realchandan/google-play-scraper/pkg/app"
	"github.com/realchandan/google-play-scraper/pkg/scraper"
)

// Options type alias
type Options = scraper.Options

// New return similar list instance
func New(appID string, options Options) *scraper.Scraper {
	a := app.New(appID, app.Options{
		Country:  options.Country,
		Language: options.Language,
	})
	err := a.LoadDetails()
	if err != nil || a.SimilarURL == "" {
		return nil
	}
	return scraper.New(a.SimilarURL, &options)
}
