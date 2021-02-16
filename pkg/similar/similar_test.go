package similar

import (
	"testing"
)

var resultsCount = 100

func TestSimilar(t *testing.T) {
	q := New("com.google.android.googlequicksearchbox", Options{
		Country:  "us",
		Language: "us",
		Number:   resultsCount,
	})
	err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(q.Results) == 0 {
		t.Errorf("Expected Results length is greater than zero")
	}
}
