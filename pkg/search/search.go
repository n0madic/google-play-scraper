package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

const (
	batchexecuteURL = "https://play.google.com/_/PlayStoreUi/data/batchexecute"
	searchURL       = "https://play.google.com/store/search"
)

// PriceQuery value
type PriceQuery int

const (
	// PriceAll - all prices
	PriceAll PriceQuery = iota
	// PriceFree - only free
	PriceFree
	// PricePaid - only paid
	PricePaid
)

// Options of query
type Options struct {
	Country  string
	Discount bool
	Language string
	Number   int
	Price    PriceQuery
	PriceMin float64
	PriceMax float64
	Query    string
	ScoreMin float64
	ScoreMax float64
}

// Query instance
type Query struct {
	options Options
	Results Results
}

func (query *Query) initialRequest() ([]app.App, string, error) {
	req, err := http.NewRequest("POST", searchURL, nil)
	if err != nil {
		return nil, "", err
	}
	q := req.URL.Query()
	q.Add("q", url.QueryEscape(query.options.Query))
	q.Add("c", "apps")
	q.Add("price", strconv.Itoa(int(query.options.Price)))
	q.Add("gl", query.options.Country)
	q.Add("hl", query.options.Language)
	req.URL.RawQuery = q.Encode()

	data, err := util.GetInitData(req)
	if err != nil {
		return nil, "", err
	}

	// return results with next token
	return query.parseSearch(data["ds:3"], "0.1.0.0.0"), util.GetJSONValue(data["ds:3"], "0.1.0.0.7.1"), nil
}

func (query *Query) batchexecute(token string) ([]app.App, string, error) {
	data := strings.Replace("f.req=%5B%5B%5B%22qnKhOb%22%2C%22%5B%5Bnull%2C%5B%5B10%2C%5B10%2C50%5D%5D%2Ctrue%2Cnull%2C%5B96%2C27%2C4%2C8%2C57%2C30%2C110%2C79%2C11%2C16%2C49%2C1%2C3%2C9%2C12%2C104%2C55%2C56%2C51%2C10%2C34%2C77%5D%5D%2Cnull%2C%5C%22{{token}}%5C%22%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D", "{{token}}", token, 1)

	req, err := http.NewRequest("POST", batchexecuteURL, strings.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	q := req.URL.Query()
	q.Add("authuser", "0")
	q.Add("bl", "boq_playuiserver_20190424.04_p0")
	q.Add("gl", query.options.Country)
	q.Add("hl", query.options.Language)
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
		return nil, "", fmt.Errorf("search error: %s", resp.Status)
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
	return query.parseSearch(reqData, "0.0.0"), nextToken, nil
}

// Do query
func (query *Query) Do() error {
	query.Results = nil

	results, token, err := query.initialRequest()
	if err != nil {
		return err
	}

	if len(results) > query.options.Number {
		query.Results.Append(results[:query.options.Number]...)
	} else {
		query.Results.Append(results...)
	}

	for len(query.Results) != query.options.Number {
		results, token, err = query.batchexecute(token)

		if len(results) == 0 {
			break
		}

		if len(query.Results)+len(results) > query.options.Number {
			query.Results.Append(results[:query.options.Number-len(query.Results)]...)
		} else {
			query.Results.Append(results...)
		}
	}

	return nil
}

// LoadMoreDetails for all search results (in concurrency)
func (query *Query) LoadMoreDetails(maxWorkers int) (errors []error) {
	if maxWorkers < 1 {
		maxWorkers = 10
	}
	semaphore := make(chan struct{}, maxWorkers)

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for _, result := range query.Results {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(result *app.App) {
			defer wg.Done()
			err := result.LoadDetails(query.options.Country, query.options.Language)
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

func (query *Query) parseSearch(data, path string) (results []app.App) {
	appData := util.GetJSONArray(data, path)
	for _, ap := range appData {
		price := app.Price{
			Currency: util.GetJSONValue(ap.String(), "7.0.3.2.1.0.1"),
			Value:    parse.Float(util.GetJSONValue(ap.String(), "7.0.3.2.1.0.2")),
		}
		if price.Value < query.options.PriceMin ||
			(query.options.PriceMax > query.options.PriceMin && price.Value > query.options.PriceMax) {
			continue
		}

		priceFull := app.Price{
			Currency: util.GetJSONValue(ap.String(), "7.0.3.2.1.1.1"),
			Value:    parse.Float(util.GetJSONValue(ap.String(), "7.0.3.2.1.1.2")),
		}
		if query.options.Discount && priceFull.Value < price.Value {
			continue
		}

		score := parse.Float(util.GetJSONValue(ap.String(), "6.0.2.1.1"))
		if score < query.options.ScoreMin ||
			(query.options.ScoreMax > query.options.ScoreMin && score > query.options.ScoreMax) {
			continue
		}

		relativeAppURL := util.GetJSONValue(ap.String(), "9.4.2")
		appURL, _ := util.AbsoluteURL(searchURL, relativeAppURL)
		relativeDevURL := util.GetJSONValue(ap.String(), "4.0.0.1.4.2")
		devURL, _ := util.AbsoluteURL(searchURL, relativeDevURL)
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

// NewQuery return Query instance
func NewQuery(options Options) *Query {
	query := &Query{options: options}
	if query.options.Number == 0 {
		query.options.Number = 50
	}
	return query
}
