package category

import (
	"testing"

	"github.com/n0madic/google-play-scraper/pkg/store"
)

var resultsCount = 10

func TestCategory(t *testing.T) {
	sortList := []store.Sort{store.SortHelpfulness, store.SortNewest, store.SortRating}
	for _, sort := range sortList {
		l, err := New(store.Business, sort, store.AgeFiveUnder, Options{
			Country:  "us",
			Language: "us",
			Number:   resultsCount,
		})
		if err != nil {
			t.Error(err)
		}

		for key, cluster := range l {
			err := cluster.Run()
			if err != nil {
				t.Error(err)
			}

			if len(cluster.Results) != resultsCount {
				t.Errorf("[%s] Expected Results length is %d, got %d", key, resultsCount, len(cluster.Results))
			}
		}
	}
}
