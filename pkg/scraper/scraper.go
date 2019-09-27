package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

// BaseURL of Google Play Store
const BaseURL = "https://play.google.com/store/apps"

const batchexecuteURL = "https://play.google.com/_/PlayStoreUi/data/batchexecute"

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
	req, err := http.NewRequest("POST", scraper.url, nil)
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

	// return results with next token
	return scraper.parseResult(data["ds:3"], "0.1.0.0.0"), util.GetJSONValue(data["ds:3"], "0.1.0.0.7.1"), nil
}

func (scraper *Scraper) batchexecute(token string) ([]app.App, string, error) {
	data := strings.Replace("f.req=%5B%5B%5B%22qnKhOb%22%2C%22%5B%5Bnull%2C%5B%5B10%2C%5B10%2C50%5D%5D%2Ctrue%2Cnull%2C%5B96%2C27%2C4%2C8%2C57%2C30%2C110%2C79%2C11%2C16%2C49%2C1%2C3%2C9%2C12%2C104%2C55%2C56%2C51%2C10%2C34%2C77%5D%5D%2Cnull%2C%5C%22{{token}}%5C%22%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D", "{{token}}", token, 1)

	req, err := http.NewRequest("POST", batchexecuteURL, strings.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	q := req.URL.Query()
	q.Add("authuser", "0")
	q.Add("bl", "boq_playuiserver_20190424.04_p0")
	q.Add("gl", scraper.options.Country)
	q.Add("hl", scraper.options.Language)
	q.Add("soc-app", "121")
	q.Add("soc-platform", "1")
	q.Add("soc-device", "1")
	q.Add("rpcids", "qnKhOb")
	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("response error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var js [][]interface{}
	err = json.Unmarshal(bytes.TrimLeft(body, ")]}'"), &js)
	if err != nil {
		return nil, "", err
	}
	if len(js) < 1 || len(js[0]) < 2 {
		return nil, "", fmt.Errorf("Invalid size of the resulting array")
	}
	if js[0][2] == nil {
		return nil, "", nil
	}

	reqData := js[0][2].(string)
	nextToken := util.GetJSONValue(reqData, "0.0.7.1")
	return scraper.parseResult(reqData, "0.0.0"), nextToken, nil
}

// Run scraping
func (scraper *Scraper) Run() error {
	scraper.Results = nil

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

		if len(results) == 0 {
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

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for _, result := range scraper.Results {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(result *app.App) {
			defer wg.Done()
			err := result.LoadDetails(scraper.options.Country, scraper.options.Language)
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

func (scraper *Scraper) parseResult(data, path string) (results []app.App) {
	appData := util.GetJSONArray(data, path)
	for _, ap := range appData {
		price := app.Price{
			Currency: util.GetJSONValue(ap.String(), "7.0.3.2.1.0.1"),
			Value:    parse.Float(util.GetJSONValue(ap.String(), "7.0.3.2.1.0.2")),
		}
		if price.Value < scraper.options.PriceMin ||
			(scraper.options.PriceMax > scraper.options.PriceMin && price.Value > scraper.options.PriceMax) {
			continue
		}

		priceFull := app.Price{
			Currency: util.GetJSONValue(ap.String(), "7.0.3.2.1.1.1"),
			Value:    parse.Float(util.GetJSONValue(ap.String(), "7.0.3.2.1.1.2")),
		}
		if scraper.options.Discount && priceFull.Value < price.Value {
			continue
		}

		score := parse.Float(util.GetJSONValue(ap.String(), "6.0.2.1.1"))
		if score < scraper.options.ScoreMin ||
			(scraper.options.ScoreMax > scraper.options.ScoreMin && score > scraper.options.ScoreMax) {
			continue
		}

		relativeAppURL := util.GetJSONValue(ap.String(), "9.4.2")
		appURL, _ := util.AbsoluteURL(scraper.url, relativeAppURL)
		relativeDevURL := util.GetJSONValue(ap.String(), "4.0.0.1.4.2")
		devURL, _ := util.AbsoluteURL(scraper.url, relativeDevURL)
		results = append(results, app.App{
			Developer:    util.GetJSONValue(ap.String(), "4.0.0.0"),
			DeveloperID:  parse.ID(relativeDevURL),
			DeveloperURL: devURL,
			Free:         price.Value == 0,
			Icon:         util.GetJSONValue(ap.String(), "1.1.0.3.2"),
			ID:           util.GetJSONValue(ap.String(), "12.0"),
			Price:        price,
			PriceFull:    priceFull,
			Score:        score,
			Summary:      util.GetJSONValue(ap.String(), "4.1.1.1.1"),
			Title:        util.GetJSONValue(ap.String(), "2"),
			URL:          appURL,
		})
	}
	return
}

// New return new Scraper instance
func New(url string, options *Options) *Scraper {
	scraper := &Scraper{
		options: options,
		url:     url,
	}
	if scraper.options.Number == 0 {
		scraper.options.Number = 50
	}
	return scraper
}
