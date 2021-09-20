package category

import (
	"fmt"
	"testing"

	"github.com/n0madic/google-play-scraper/pkg/store"
)

var resultsCount = 10

func TestCategory(t *testing.T) {
	sortList := []store.Sort{store.SortHelpfulness, store.SortNewest, store.SortRating}
	for _, sort := range sortList {
		l, err := New(store.Game, sort, store.AgeFiveUnder, Options{
			Country:  "us",
			Language: "us",
			Number:   resultsCount,
		})
		if err != nil {
			t.Error(err)
		}

		fmt.Println(l)
		if len(l) < 1 {
			t.Errorf("No empty clusters expected")
		} else {
			for key, cluster := range l {
				err := cluster.Run()
				if err != nil {
					t.Error(err)
				}

				if len(cluster.Results) == 0 {
					t.Errorf("[%s] Expected non-zero Results length", key)
				}
			}
		}
	}
}
