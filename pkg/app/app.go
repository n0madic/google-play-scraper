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
	ContentRating            string
	ContentRatingDescription string
	Description              string
	DescriptionHTML          string
	Developer                string
	DeveloperAddress         string
	DeveloperEmail           string
	DeveloperID              string
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
	Screenshots              []string
	SimilarURL               string
	Size                     string
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

const (
	dsAppInfo    = "ds:5"
	dsAppPrice   = "ds:3"
	dsAppRating  = "ds:6"
	dsAppSimilar = "ds:7"
	dsAppVersion = "ds:8"
)

var dsAppReview = [...]string{"ds:17", "ds:18", "ds:19"}

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

	app.AdSupported = util.GetJSONValue(appData[dsAppInfo], "0.12.14.0") != ""

	app.AndroidVersion = util.GetJSONValue(appData[dsAppVersion], "2")
	app.AndroidVersionMin = parse.Float(app.AndroidVersion)

	app.ContentRating = util.GetJSONValue(appData[dsAppInfo], "0.12.4.0")
	app.ContentRatingDescription = util.GetJSONValue(appData[dsAppInfo], "0.12.4.2.1")

	app.DescriptionHTML = util.GetJSONValue(appData[dsAppInfo], "0.10.0.1")
	app.Description = util.HTMLToText(app.DescriptionHTML)

	relativeDevURL := util.GetJSONValue(appData[dsAppRating], "0.12.5.5.4.2")
	devURL, _ := util.AbsoluteURL(playURL, relativeDevURL)
	app.Developer = util.GetJSONValue(appData[dsAppInfo], "0.12.5.1")
	app.DeveloperAddress = util.GetJSONValue(appData[dsAppInfo], "0.12.5.4.0")
	app.DeveloperEmail = util.GetJSONValue(appData[dsAppInfo], "0.12.5.2.0")
	app.DeveloperID = util.GetJSONValue(appData[dsAppInfo], "0.12.5.0.0")
	app.DeveloperURL = devURL
	app.DeveloperWebsite = util.GetJSONValue(appData[dsAppInfo], "0.12.5.3.5.2")

	app.Genre = util.GetJSONValue(appData[dsAppInfo], "0.12.13.0.0")
	app.GenreID = util.GetJSONValue(appData[dsAppInfo], "0.12.13.0.2")
	app.FamilyGenre = util.GetJSONValue(appData[dsAppInfo], "0.12.13.1.0")
	app.FamilyGenreID = util.GetJSONValue(appData[dsAppInfo], "0.12.13.1.2")

	app.HeaderImage = util.GetJSONValue(appData[dsAppInfo], "0.12.2.3.2")

	app.IAPRange = util.GetJSONValue(appData[dsAppInfo], "0.12.12.0")
	app.IAPOffers = app.IAPRange != ""

	app.Icon = util.GetJSONValue(appData[dsAppInfo], "0.12.1.3.2")

	app.Installs = util.GetJSONValue(appData[dsAppInfo], "0.12.9.0")
	app.InstallsMin = parse.Int(app.Installs)

	price := Price{
		Currency: util.GetJSONValue(appData[dsAppPrice], "0.2.0.0.0.1.0.1"),
		Value:    parse.Float(util.GetJSONValue(appData[dsAppPrice], "0.2.0.0.0.1.0.2")),
	}
	app.Free = price.Value == 0
	app.Price = price
	app.PriceFull = Price{
		Currency: util.GetJSONValue(appData[dsAppPrice], "0.2.0.0.0.1.1.1"),
		Value:    parse.Float(util.GetJSONValue(appData[dsAppPrice], "0.2.0.0.0.1.1.2")),
	}

	app.PrivacyPolicy = util.GetJSONValue(appData[dsAppInfo], "0.12.7.2")

	app.Ratings = parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.2.1"))
	app.RatingsHistogram = map[int]int{
		1: parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.1.1.1")),
		2: parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.1.2.1")),
		3: parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.1.3.1")),
		4: parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.1.4.1")),
		5: parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.1.5.1")),
	}

	for _, section := range dsAppReview {
		reviewList := util.GetJSONArray(appData[section], "0")
		if len(reviewList) > 2 {
			for _, review := range reviewList {
				r := reviews.Parse(review.String())
				if r != nil {
					app.Reviews = append(app.Reviews, r)
				}
			}
		}
	}
	app.ReviewsTotalCount = parse.Int(util.GetJSONValue(appData[dsAppRating], "0.6.3.1"))

	screenshots := util.GetJSONArray(appData[dsAppInfo], "0.12.0")
	for _, screen := range screenshots {
		app.Screenshots = append(app.Screenshots, util.GetJSONValue(screen.String(), "3.2"))
	}

	similarURL := util.GetJSONValue(appData[dsAppSimilar], "1.1.0.0.3.4.2")
	if similarURL != "" {
		app.SimilarURL, _ = util.AbsoluteURL(playURL, similarURL)
	}

	app.RecentChangesHTML = util.GetJSONValue(appData[dsAppInfo], "0.12.6.1")
	app.RecentChanges = util.HTMLToText(app.RecentChangesHTML)
	app.Released = util.GetJSONValue(appData[dsAppInfo], "0.12.36")
	app.Score = parse.Float(util.GetJSONValue(appData[dsAppRating], "0.6.0.1"))
	app.Size = util.GetJSONValue(appData[dsAppVersion], "0")
	app.Summary = util.GetJSONValue(appData[dsAppInfo], "0.10.1.1")
	app.Title = util.GetJSONValue(appData[dsAppInfo], "0.0.0")
	app.Updated = time.Unix(parse.Int64(util.GetJSONValue(appData[dsAppInfo], "0.12.8.0")), 0)
	app.Version = util.GetJSONValue(appData[dsAppVersion], "1")
	app.Video = util.GetJSONValue(appData[dsAppInfo], "0.12.3.0.3.2")
	app.VideoImage = util.GetJSONValue(appData[dsAppInfo], "0.12.3.1.3.2")

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
