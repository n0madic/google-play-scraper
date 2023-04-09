package developer

import (
	"net/url"

	"github.com/realchandan/google-play-scraper/pkg/scraper"
)

// Options type alias
type Options = scraper.Options

// New return the list of applications by the given developer name
func New(name string, options Options) *scraper.Scraper {
	return new("/developer", name, &options)
}

// NewByID return the list of applications by the given developer ID
func NewByID(devID string, options Options) *scraper.Scraper {
	return new("/dev", devID, &options)
}

func new(path, name string, options *Options) *scraper.Scraper {
	u, err := url.Parse(scraper.BaseURL + path)
	if err != nil {
		return nil
	}
	q := u.Query()
	q.Set("id", name)
	u.RawQuery = q.Encode()
	return scraper.New(u.String(), options)
}
