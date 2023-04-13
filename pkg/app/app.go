package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/reviews"
)

const (
	detailURL = "https://play.google.com/store/apps/details?id="
	playURL   = "https://play.google.com"
)

// Price of app
type Price struct {
	Currency string
	Value    float64
}

// App of search
type App struct {
	AdSupported              bool
	AndroidVersion           string
	AndroidVersionMin        float64
	Available                bool
	ContentRating            string
	ContentRatingDescription string
	Description              string
	DescriptionHTML          string
	Developer                string
	DeveloperAddress         string
	DeveloperEmail           string
	DeveloperID              string
	DeveloperInternalID      string
	DeveloperURL             string
	DeveloperWebsite         string
	FamilyGenre              string
	FamilyGenreID            string
	Free                     bool
	Genre                    string
	GenreID                  string
	HeaderImage              string
	IAPOffers                bool
	IAPRange                 string
	Icon                     string
	ID                       string
	Installs                 string
	InstallsMin              int
	InstallsMax              int
	Permissions              map[string][]string
	Price                    Price
	PriceFull                Price
	PrivacyPolicy            string
	Ratings                  int
	RatingsHistogram         map[int]int
	RecentChanges            string
	RecentChangesHTML        string
	Released                 string
	Reviews                  []*reviews.Review
	ReviewsTotalCount        int
	Score                    float64
	ScoreText                string
	Screenshots              []string
	SimilarURL               string
	Summary                  string
	Title                    string
	Updated                  time.Time
	URL                      string
	Version                  string
	Video                    string
	VideoImage               string
	options                  *Options
}

// Options of app
type Options struct {
	Country  string
	Language string
}

// LoadDetails of app
func (app *App) LoadDetails() error {
	if app.URL == "" {
		if app.ID != "" {
			app.URL = detailURL + app.ID
		} else {
			return fmt.Errorf("App ID or URL required")
		}
	}

	req, err := http.NewRequest("GET", app.URL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("gl", app.options.Country)
	q.Add("hl", app.options.Language)
	req.URL.RawQuery = q.Encode()

	appData, err := util.GetInitData(req)
	if err != nil {
		return err
	}

	if app.ID == "" {
		app.ID = parse.ID(app.URL)
	}

	for dsAppInfo := range appData {
		relativeDevURL := util.GetJSONValue(appData[dsAppInfo], "1.2.68.1.4.2")
		if relativeDevURL == "" {
			continue
		}

		app.AdSupported = util.GetJSONValue(appData[dsAppInfo], "1.2.48") != ""

		app.AndroidVersion = util.GetJSONValue(appData[dsAppInfo], "1.2.140.1.1.0.0.1")
		app.AndroidVersionMin = parse.Float(app.AndroidVersion)

		app.Available = util.GetJSONValue(appData[dsAppInfo], "1.2.18.0") != ""

		app.ContentRating = util.GetJSONValue(appData[dsAppInfo], "1.2.9.0")
		app.ContentRatingDescription = util.GetJSONValue(appData[dsAppInfo], "1.2.9.2.1")

		app.DescriptionHTML = util.GetJSONValue(appData[dsAppInfo], "1.2.72.0.1")
		app.Description = util.HTMLToText(app.DescriptionHTML)

		devURL, _ := util.AbsoluteURL(playURL, relativeDevURL)
		app.Developer = util.GetJSONValue(appData[dsAppInfo], "1.2.68.0")
		app.DeveloperAddress = util.GetJSONValue(appData[dsAppInfo], "1.2.69.2.0")
		app.DeveloperEmail = util.GetJSONValue(appData[dsAppInfo], "1.2.69.1.0")
		app.DeveloperID = parse.ID(util.GetJSONValue(appData[dsAppInfo], "1.2.68.1.4.2"))
		app.DeveloperInternalID = util.GetJSONValue(appData[dsAppInfo], "1.2.68.2")
		app.DeveloperURL = devURL
		app.DeveloperWebsite = util.GetJSONValue(appData[dsAppInfo], "1.2.69.0.5.2")

		app.Genre = util.GetJSONValue(appData[dsAppInfo], "1.2.79.0.0.0")
		app.GenreID = util.GetJSONValue(appData[dsAppInfo], "1.2.79.0.0.2")
		app.FamilyGenre = util.GetJSONValue(appData[dsAppInfo], "1.12.13.1.0")
		app.FamilyGenreID = util.GetJSONValue(appData[dsAppInfo], "1.12.13.1.2")

		app.HeaderImage = util.GetJSONValue(appData[dsAppInfo], "1.2.96.0.3.2")

		app.IAPRange = util.GetJSONValue(appData[dsAppInfo], "1.2.19.0")
		app.IAPOffers = app.IAPRange != ""

		app.Icon = util.GetJSONValue(appData[dsAppInfo], "1.2.95.0.3.2")

		app.Installs = util.GetJSONValue(appData[dsAppInfo], "1.2.13.0")
		app.InstallsMin = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.13.1"))
		app.InstallsMax = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.13.2"))

		price := Price{
			Currency: util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.0.1"),
			Value:    parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.0.2")),
		}
		app.Free = price.Value == 0
		app.Price = price
		app.PriceFull = Price{
			Currency: util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.1.1"),
			Value:    parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.57.0.0.0.0.1.1.2")),
		}

		app.PrivacyPolicy = util.GetJSONValue(appData[dsAppInfo], "1.2.99.0.5.2")

		app.Ratings = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.2.1"))
		app.RatingsHistogram = map[int]int{
			1: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.1.1")),
			2: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.2.1")),
			3: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.3.1")),
			4: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.4.1")),
			5: parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.1.5.1")),
		}

		for dsAppReview := range appData {
			reviewList := util.GetJSONArray(appData[dsAppReview], "0")
			check := util.GetJSONValue(appData[dsAppReview], "2.0.1.2.0")
			if len(reviewList) > 2 && check != "" {
				for _, review := range reviewList {
					r := reviews.Parse(review.String())
					if r != nil {
						app.Reviews = append(app.Reviews, r)
					}
				}
				break
			}
		}
		app.ReviewsTotalCount = parse.Int(util.GetJSONValue(appData[dsAppInfo], "1.2.51.2.1"))

		screenshots := util.GetJSONArray(appData[dsAppInfo], "1.2.78.0")
		for _, screen := range screenshots {
			app.Screenshots = append(app.Screenshots, util.GetJSONValue(screen.String(), "3.2"))
		}

		for dsAppSimilar := range appData {
			similarURL := util.GetJSONValue(appData[dsAppSimilar], "1.1.1.21.1.2.4.2")
			if similarURL != "" {
				app.SimilarURL, _ = util.AbsoluteURL(playURL, similarURL)
				break
			}
		}

		app.RecentChangesHTML = util.GetJSONValue(appData[dsAppInfo], "1.2.144.1.1", "1.2.145.0.0")
		app.RecentChanges = util.HTMLToText(app.RecentChangesHTML)
		app.Released = util.GetJSONValue(appData[dsAppInfo], "1.2.10.0")
		app.Score = parse.Float(util.GetJSONValue(appData[dsAppInfo], "1.2.51.0.1"))
		app.ScoreText = util.GetJSONValue(appData[dsAppInfo], "1.2.51.0.0")
		app.Summary = util.GetJSONValue(appData[dsAppInfo], "1.2.73.0.1")
		app.Title = util.GetJSONValue(appData[dsAppInfo], "1.2.0.0")
		app.Updated = time.Unix(parse.Int64(util.GetJSONValue(appData[dsAppInfo], "1.2.145.0.1.0")), 0)
		app.Version = util.GetJSONValue(appData[dsAppInfo], "1.2.140.0.0.0")
		app.Video = util.GetJSONValue(appData[dsAppInfo], "1.2.100.0.0.3.2")
		app.VideoImage = util.GetJSONValue(appData[dsAppInfo], "1.2.100.1.0.3.2")
	}
	return nil
}

// LoadPermissions get the list of perms an app has access to
func (app *App) LoadPermissions() error {
	payload := strings.Replace("f.req=%5B%5B%5B%22xdSrCf%22%2C%22%5B%5Bnull%2C%5B%5C%22{{appID}}%5C%22%2C7%5D%2C%5B%5D%5D%5D%22%2Cnull%2C%221%22%5D%5D%5D", "{{appID}}", app.ID, 1)

	js, err := util.BatchExecute(app.options.Country, app.options.Language, payload)
	if err != nil {
		return err
	}

	app.Permissions = make(map[string][]string)
	for _, perm := range util.GetJSONArray(js, "0") {
		key := util.GetJSONValue(perm.String(), "0")
		for _, permission := range util.GetJSONArray(perm.String(), "2") {
			app.Permissions[key] = append(app.Permissions[key], util.GetJSONValue(permission.String(), "1"))
		}
	}

	return nil
}

// New return App instance
func New(id string, options Options) *App {
	return &App{
		ID:      id,
		URL:     detailURL + id,
		options: &options,
	}
}
