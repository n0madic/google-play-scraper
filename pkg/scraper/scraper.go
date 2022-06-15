package scraper

import (
	"net/http"
	"strings"
	"sync"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

// BaseURL of Google Play Store
const BaseURL = "https://play.google.com/store/apps"

// Options of scraper
type Options struct {
	Country  string
	Discount bool
	Language string
	Number   int
	PriceMin float64
	PriceMax float64
	ScoreMin float64
	ScoreMax float64
}

// Scraper instance
type Scraper struct {
	options *Options
	Results Results
	url     string
}

func (scraper *Scraper) initialRequest() ([]app.App, string, error) {
	req, err := http.NewRequest("GET", scraper.url, nil)
	if err != nil {
		return nil, "", err
	}
	q := req.URL.Query()
	q.Add("gl", scraper.options.Country)
	q.Add("hl", scraper.options.Language)
	req.URL.RawQuery = q.Encode()

	data, err := util.GetInitData(req)
	if err != nil {
		return nil, "", err
	}

	// path for develeoper page by DevName is "0.1.0.22.0"
	// path for developer page by DevId is "0.1.0.21.0"
	// return results with next token
	return scraper.parseResult(data["ds:3"], "0.1.0.22.0", "0.1.0.21.0"),
		util.GetJSONValue(data["ds:3"], "0.1.0.22.1.3.1", "0.1.0.21.1.3.1"),
		nil
}

func (scraper *Scraper) batchexecute(token string) ([]app.App, string, error) {
	payload := strings.Replace("f.req=%5B%5B%5B%22qnKhOb%22%2C%22%5B%5Bnull%2C%5B%5B10%2C%5B10%2C50%5D%5D%2Ctrue%2Cnull%2C%5B96%2C27%2C4%2C8%2C57%2C30%2C110%2C79%2C11%2C16%2C49%2C1%2C3%2C9%2C12%2C104%2C55%2C56%2C51%2C10%2C34%2C77%5D%5D%2Cnull%2C%5C%22{{token}}%5C%22%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D", "{{token}}", token, 1)

	js, err := util.BatchExecute(scraper.options.Country, scraper.options.Language, payload)
	if err != nil {
		return nil, "", err
	}

	nextToken := util.GetJSONValue(js, "0.0.7.1")
	return scraper.parseResult(js, "0.0.0"), nextToken, nil
}

// Run scraping
func (scraper *Scraper) Run() error {
	scraper.Results = Results{}

	results, token, err := scraper.initialRequest()
	if err != nil {
		return err
	}

	if len(results) > scraper.options.Number {
		scraper.Results.Append(results[:scraper.options.Number]...)
	} else {
		scraper.Results.Append(results...)
	}

	for len(scraper.Results) != scraper.options.Number {
		results, token, err = scraper.batchexecute(token)

		if len(results) == 0 || err != nil {
			break
		}

		if len(scraper.Results)+len(results) > scraper.options.Number {
			scraper.Results.Append(results[:scraper.options.Number-len(scraper.Results)]...)
		} else {
			scraper.Results.Append(results...)
		}
	}

	return nil
}

// LoadMoreDetails for all results (in concurrency)
func (scraper *Scraper) LoadMoreDetails(maxWorkers int) (errors []error) {
	if maxWorkers < 1 {
		maxWorkers = 10
	}
	semaphore := make(chan struct{}, maxWorkers)

	mutex := &sync.Mutex{}
	var wg sync.WaitGroup

	for _, result := range scraper.Results {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(result *app.App) {
			defer wg.Done()
			err := result.LoadDetails()
			if err != nil {
				mutex.Lock()
				errors = append(errors, err)
				mutex.Unlock()
			}
			<-semaphore
		}(result)
	}
	wg.Wait()

	return
}

func (scraper *Scraper) parseResult(data string, paths ...string) (results []app.App) {
	for _, path := range paths {
		appData := util.GetJSONArray(data, path)
		for _, ap := range appData {
			price := app.Price{
				Currency: util.GetJSONValue(ap.String(), "0.8.1.0.1", "8.1.0.1", "7.0.3.2.1.0.1"),
				Value:    parse.Float(util.GetJSONValue(ap.String(), "0.8.1.0.2", "8.1.0.2", "7.0.3.2.1.0.2")),
			}
			if price.Value < scraper.options.PriceMin ||
				(scraper.options.PriceMax > scraper.options.PriceMin && price.Value > scraper.options.PriceMax) {
				continue
			}

			priceFull := app.Price{
				Currency: util.GetJSONValue(ap.String(), "0.8.1.0.1", "8.1.0.1", "7.0.3.2.1.1.1"),
				Value:    parse.Float(util.GetJSONValue(ap.String(), "0.8.1.0.0", "8.1.0.0", "7.0.3.2.1.1.2")),
			}
			if scraper.options.Discount && priceFull.Value < price.Value {
				continue
			}

			score := parse.Float(util.GetJSONValue(ap.String(), "0.4.0", "4.0", "6.0.2.1.1"))
			if score < scraper.options.ScoreMin ||
				(scraper.options.ScoreMax > scraper.options.ScoreMin && score > scraper.options.ScoreMax) {
				continue
			}

			application := app.New(util.GetJSONValue(ap.String(), "0.0.0", "0.0", "12.0"), app.Options{
				Country:  scraper.options.Country,
				Language: scraper.options.Language,
			})

			application.DeveloperURL, _ = util.AbsoluteURL(scraper.url, util.GetJSONValue(ap.String(), "4.0.0.1.4.2"))
			application.Developer = util.GetJSONValue(ap.String(), "0.14", "14", "4.0.0.0")
			application.DeveloperID = parse.ID(application.DeveloperURL)
			application.Free = price.Value == 0
			application.Icon = util.GetJSONValue(ap.String(), "0.1.3.2", "1.3.2", "1.1.0.3.2")
			application.Price = price
			application.PriceFull = priceFull
			application.Score = score
			application.Summary = util.GetJSONValue(ap.String(), "0.13.1", "13.1", "4.1.1.1.1")
			application.Title = util.GetJSONValue(ap.String(), "0.3", "3", "2")
			application.URL, _ = util.AbsoluteURL(scraper.url, util.GetJSONValue(ap.String(), "0.10.4.2", "10.4.2", "9.4.2"))
			results = append(results, *application)
		}
	}
	return results
}

// New return new Scraper instance
func New(url string, options *Options) *Scraper {
	scraper := &Scraper{
		Results: Results{},
		options: options,
		url:     url,
	}
	if scraper.options.Number == 0 {
		scraper.options.Number = 50
	}
	return scraper
}
