package search

import (
	"testing"
)

var resultsCount = 20

func TestSearch(t *testing.T) {
	q := NewQuery("match", PriceAll, Options{
		Country:  "us",
		Language: "en",
		Number:   resultsCount,
	})
	err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(q.Results) != resultsCount {
		t.Errorf("Expected Results length is %d, got %d", resultsCount, len(q.Results))
	}

	errors := q.LoadMoreDetails(0)
	for err := range errors {
		t.Error(err)
	}

	dupCheck := map[string]bool{}

	for _, result := range q.Results {
		if result.Description == "" {
			t.Error("Expected Description", result.ID)
		}
		if _, exist := dupCheck[result.ID]; exist {
			t.Errorf("Duplicate ID %s found in results", result.ID)
		} else {
			dupCheck[result.ID] = true
		}
	}

}
