package suggest

import (
	"testing"
)

var resultsCount = 5

func TestSuggest(t *testing.T) {
	list, err := Get("test", Options{})
	if err != nil {
		t.Error(err)
	}

	if len(list) != resultsCount {
		t.Errorf("Expected suggest list length is %d, got %d", resultsCount, len(list))
	}
}
