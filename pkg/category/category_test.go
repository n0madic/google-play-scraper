package category

import (
	"testing"

	"github.com/realchandan/google-play-scraper/pkg/store"
)

var resultsCount = 10

func TestCategory(t *testing.T) {
	l, err := New(store.Game, store.AgeFiveUnder, Options{
		Country:  "us",
		Language: "us",
		Number:   resultsCount,
	})
	if err != nil {
		t.Error(err)
	}

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
