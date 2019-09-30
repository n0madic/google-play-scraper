package collection

import (
	"testing"

	"github.com/n0madic/google-play-scraper/pkg/store"
)

var resultsCount = 100

func TestCollection(t *testing.T) {
	q := New(store.TopNewPaid, Options{
		Number: resultsCount,
	})
	err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(q.Results) != resultsCount {
		t.Errorf("Expected Results length is %d, got %d", resultsCount, len(q.Results))
	}

	for _, app := range q.Results {
		if app.Price.Value == 0 {
			t.Error("Expected Price is greater than zero")
		}
	}
}
