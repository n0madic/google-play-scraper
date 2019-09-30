package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
)

const (
	detailURL = "https://play.google.com/store/apps/details?id="
	playURL   = "https://play.google.com"
)

// Comment of app
type Comment struct {
	Answer          string
	AnswerTimestamp time.Time
	Answerer        string
	Avatar          string
	Commentator     string
	Rating          int
	Text            string
	Timestamp       time.Time
	Useful          int
}

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
	Comments                 []*Comment
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
	Price                    Price
	PriceFull                Price
	PrivacyPolicy            string
	Ratings                  int
	RatingsHistogram         map[int]int
	RecentChanges            string
	RecentChangesHTML        string
	Released                 string
	Reviews                  string
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
}

// LoadDetails of app
func (app *App) LoadDetails(country, language string) error {
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
	q.Add("gl", country)
	q.Add("hl", language)
	req.URL.RawQuery = q.Encode()

	appData, err := util.GetInitData(req)
	if err != nil {
		return err
	}

	if app.ID == "" {
		app.ID = parse.ID(app.URL)
	}

	app.AdSupported = util.GetJSONValue(appData["ds:5"], "0.12.14.0") != ""

	app.AndroidVersion = util.GetJSONValue(appData["ds:8"], "2")
	app.AndroidVersionMin = parse.Float(app.AndroidVersion)

	comments := util.GetJSONArray(appData["ds:15"], "0")
	for _, comment := range comments {
		text := util.GetJSONValue(comment.String(), "4")
		if text != "" {
			app.Comments = append(app.Comments, &Comment{
				Answer:          util.GetJSONValue(comment.String(), "7.1"),
				AnswerTimestamp: time.Unix(parse.Int64(util.GetJSONValue(comment.String(), "7.2.0")), 0),
				Answerer:        util.GetJSONValue(comment.String(), "7.0"),
				Avatar:          util.GetJSONValue(comment.String(), "1.1.3.2"),
				Commentator:     util.GetJSONValue(comment.String(), "1.0"),
				Rating:          parse.Int(util.GetJSONValue(comment.String(), "2")),
				Text:            text,
				Timestamp:       time.Unix(parse.Int64(util.GetJSONValue(comment.String(), "5.0")), 0),
				Useful:          parse.Int(util.GetJSONValue(comment.String(), "6")),
			})
		}
	}

	app.ContentRating = util.GetJSONValue(appData["ds:5"], "0.12.4.0")
	app.ContentRatingDescription = util.GetJSONValue(appData["ds:5"], "0.12.4.2.1")

	app.DescriptionHTML = util.GetJSONValue(appData["ds:5"], "0.10.0.1")
	app.Description = util.HTMLToText(app.DescriptionHTML)

	relativeDevURL := util.GetJSONValue(appData["ds:5"], "0.12.5.5.4.2")
	devURL, _ := util.AbsoluteURL(playURL, relativeDevURL)
	app.Developer = util.GetJSONValue(appData["ds:5"], "0.12.5.1")
	app.DeveloperAddress = util.GetJSONValue(appData["ds:5"], "0.12.5.4.0")
	app.DeveloperEmail = util.GetJSONValue(appData["ds:5"], "0.12.5.2.0")
	app.DeveloperID = util.GetJSONValue(appData["ds:5"], "0.12.5.0.0")
	app.DeveloperURL = devURL
	app.DeveloperWebsite = util.GetJSONValue(appData["ds:5"], "0.12.5.3.5.2")

	app.Genre = util.GetJSONValue(appData["ds:5"], "0.12.13.0.0")
	app.GenreID = util.GetJSONValue(appData["ds:5"], "0.12.13.0.2")
	app.FamilyGenre = util.GetJSONValue(appData["ds:5"], "0.12.13.1.0")
	app.FamilyGenreID = util.GetJSONValue(appData["ds:5"], "0.12.13.1.2")

	app.HeaderImage = util.GetJSONValue(appData["ds:5"], "0.12.2.3.2")

	app.IAPRange = util.GetJSONValue(appData["ds:5"], "0.12.12.0")
	app.IAPOffers = app.IAPRange != ""

	app.Icon = util.GetJSONValue(appData["ds:5"], "0.12.1.3.2")

	app.Installs = util.GetJSONValue(appData["ds:5"], "0.12.9.0")
	app.InstallsMin = parse.Int(app.Installs)

	price := Price{
		Currency: util.GetJSONValue(appData["ds:3"], "0.2.0.0.0.1.0.1"),
		Value:    parse.Float(util.GetJSONValue(appData["ds:3"], "0.2.0.0.0.1.0.2")),
	}
	app.Free = price.Value == 0
	app.Price = price
	app.PriceFull = Price{
		Currency: util.GetJSONValue(appData["ds:3"], "0.2.0.0.0.1.1.1"),
		Value:    parse.Float(util.GetJSONValue(appData["ds:3"], "0.2.0.0.0.1.1.2")),
	}

	app.PrivacyPolicy = util.GetJSONValue(appData["ds:5"], "0.12.7.2")

	app.Ratings = parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.2.1"))
	app.RatingsHistogram = map[int]int{
		1: parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.1.1.1")),
		2: parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.1.2.1")),
		3: parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.1.3.1")),
		4: parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.1.4.1")),
		5: parse.Int(util.GetJSONValue(appData["ds:6"], "0.6.1.5.1")),
	}

	screenshots := util.GetJSONArray(appData["ds:5"], "0.12.0")
	for _, screen := range screenshots {
		app.Screenshots = append(app.Screenshots, util.GetJSONValue(screen.String(), "3.2"))
	}

	similarURL := util.GetJSONValue(appData["ds:7"], "1.1.0.0.3.4.2")
	if similarURL != "" {
		app.SimilarURL, _ = util.AbsoluteURL(playURL, similarURL)
	}

	app.RecentChangesHTML = util.GetJSONValue(appData["ds:5"], "0.12.6.1")
	app.RecentChanges = util.HTMLToText(app.RecentChangesHTML)
	app.Released = util.GetJSONValue(appData["ds:5"], "0.12.36")
	app.Reviews = util.GetJSONValue(appData["ds:6"], "0.6.3.1")
	app.Score = parse.Float(util.GetJSONValue(appData["ds:6"], "0.6.0.1"))
	app.Size = util.GetJSONValue(appData["ds:8"], "0")
	app.Summary = util.GetJSONValue(appData["ds:5"], "0.10.1.1")
	app.Title = util.GetJSONValue(appData["ds:5"], "0.0.0")
	app.Updated = time.Unix(parse.Int64(util.GetJSONValue(appData["ds:5"], "0.12.8.0")), 0)
	app.Version = util.GetJSONValue(appData["ds:8"], "1")
	app.Video = util.GetJSONValue(appData["ds:5"], "0.12.3.0.3.2")
	app.VideoImage = util.GetJSONValue(appData["ds:5"], "0.12.3.1.3.2")

	return nil
}

// New return App instance
func New(id string) *App {
	return &App{
		ID:  id,
		URL: detailURL + id,
	}
}
