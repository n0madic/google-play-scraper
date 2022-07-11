package developer

import (
	"testing"
)

var resultsCount = 47

func TestDeveloper(t *testing.T) {
	q := New("Google LLC", Options{
		Number: resultsCount,
	})
	err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(q.Results) != resultsCount {
		t.Errorf("Expected Results length is %d, got %d", resultsCount, len(q.Results))
	}
}

func TestDeveloperByID(t *testing.T) {
	// Test on Google LLC
	q := NewByID("5700313618786177705", Options{
		Number: resultsCount,
	})
	err := q.Run()
	if err != nil {
		t.Error(err)
	}

	if len(q.Results) != resultsCount {
		t.Errorf("Expected Results length is %d, got %d", resultsCount, len(q.Results))
	}
}
