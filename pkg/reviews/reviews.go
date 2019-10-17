package reviews

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/n0madic/google-play-scraper/internal/parse"
	"github.com/n0madic/google-play-scraper/internal/util"
	"github.com/n0madic/google-play-scraper/pkg/store"
)

const (
	initialRequest   = `f.req=%5B%5B%5B%22UsvDTd%22%2C%22%5Bnull%2Cnull%2C%5B2%2C{{sort}}%2C%5B{{numberOfReviewsPerRequest}}%2Cnull%2Cnull%5D%2Cnull%2C%5B%5D%5D%2C%5B%5C%22{{appId}}%5C%22%2C7%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D`
	paginatedRequest = `f.req=%5B%5B%5B%22UsvDTd%22%2C%22%5Bnull%2Cnull%2C%5B2%2C{{sort}}%2C%5B{{numberOfReviewsPerRequest}}%2Cnull%2C%5C%22{{withToken}}%5C%22%5D%2Cnull%2C%5B%5D%5D%2C%5B%5C%22{{appId}}%5C%22%2C7%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D`
)

var numberOfReviewsPerRequest = 40

// Options of reviews
type Options struct {
	Country  string
	Language string
	Number   int
	Sorting  store.Sort
}

// Review of app
type Review struct {
	Avatar         string
	Criterias      map[string]int64
	ID             string
	Score          int
	Reviewer       string
	Reply          string
	ReplyTimestamp time.Time
	Respondent     string
	Text           string
	Timestamp      time.Time
	Useful         int
	Version        string
}

// URL of review
func (r *Review) URL(appID string) string {
	if r.ID != "" {
		return fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&reviewId=%s", appID, r.ID)
	}
	return ""
}

// Reviews instance
type Reviews struct {
	appID   string
	options *Options
	Results Results
}

// New return similar list instance
func New(appID string, options Options) *Reviews {
	if options.Number == 0 {
		options.Number = numberOfReviewsPerRequest
	}
	if options.Sorting == 0 {
		options.Sorting = store.SortHelpfulness
	}
	return &Reviews{
		appID:   appID,
		options: &options,
	}
}

func (reviews *Reviews) batchexecute(payload string) ([]Review, string, error) {
	js, err := util.BatchExecute(reviews.options.Country, reviews.options.Language, payload)
	if err != nil {
		return nil, "", err
	}

	nextToken := util.GetJSONValue(js, "1.1")

	var results []Review
	rev := util.GetJSONArray(js, "0")
	for _, review := range rev {
		result := Parse(review.String())
		if result != nil {
			results = append(results, *result)
		}
	}

	return results, nextToken, nil
}

// Run reviews scraping
func (reviews *Reviews) Run() error {
	if numberOfReviewsPerRequest > reviews.options.Number {
		numberOfReviewsPerRequest = reviews.options.Number
	}

	r := strings.NewReplacer("{{sort}}", strconv.Itoa(int(reviews.options.Sorting)),
		"{{numberOfReviewsPerRequest}}", strconv.Itoa(numberOfReviewsPerRequest),
		"{{appId}}", string(reviews.appID),
	)
	payload := r.Replace(initialRequest)

	results, token, err := reviews.batchexecute(payload)
	if err != nil {
		return err
	}

	if len(results) > reviews.options.Number {
		reviews.Results.Append(results[:reviews.options.Number]...)
	} else {
		reviews.Results.Append(results...)
	}

	for len(reviews.Results) != reviews.options.Number {
		r := strings.NewReplacer("{{sort}}", strconv.Itoa(int(reviews.options.Sorting)),
			"{{numberOfReviewsPerRequest}}", strconv.Itoa(numberOfReviewsPerRequest),
			"{{withToken}}", token,
			"{{appId}}", string(reviews.appID),
		)
		payload := r.Replace(paginatedRequest)

		results, token, err = reviews.batchexecute(payload)

		if len(results) == 0 {
			break
		}

		if len(reviews.Results)+len(results) > reviews.options.Number {
			reviews.Results.Append(results[:reviews.options.Number-len(reviews.Results)]...)
		} else {
			reviews.Results.Append(results...)
		}
	}
	return nil
}

// Parse app review
func Parse(review string) *Review {
	text := util.GetJSONValue(review, "4")
	if text != "" {
		criteriasList := util.GetJSONArray(review, "12.0")
		criterias := make(map[string]int64, len(criteriasList))
		for _, criteria := range criteriasList {
			var rating int64
			if len(criteria.Array()) > 2 {
				rating = criteria.Array()[2].Array()[0].Int()
			}
			criterias[criteria.Array()[0].String()] = rating
		}
		return &Review{
			Avatar:         util.GetJSONValue(review, "1.1.3.2"),
			Criterias:      criterias,
			ID:             util.GetJSONValue(review, "0"),
			Reply:          util.GetJSONValue(review, "7.1"),
			ReplyTimestamp: time.Unix(parse.Int64(util.GetJSONValue(review, "7.2.0")), 0),
			Respondent:     util.GetJSONValue(review, "7.0"),
			Reviewer:       util.GetJSONValue(review, "1.0"),
			Score:          parse.Int(util.GetJSONValue(review, "2")),
			Text:           text,
			Timestamp:      time.Unix(parse.Int64(util.GetJSONValue(review, "5.0")), 0),
			Useful:         parse.Int(util.GetJSONValue(review, "6")),
			Version:        util.GetJSONValue(review, "10"),
		}
	}
	return nil
}
