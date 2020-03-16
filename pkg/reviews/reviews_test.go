package reviews

import (
	"testing"
)

var resultsCount = 88

func TestReviews(t *testing.T) {
	r := New("com.viber.voip", Options{
		Number: resultsCount,
	})
	err := r.Run()
	if err != nil {
		t.Error(err)
	}

	if len(r.Results) != resultsCount {
		t.Errorf("Expected Results length is %d, got %d", resultsCount, len(r.Results))
	} else {
		for i, review := range r.Results {
			if review.Avatar == "" {
				t.Errorf("Expected reviewer Avatar in Results[%d]: %+v", i, review)
			}
			if review.ID == "" {
				t.Errorf("Expected review ID in Results[%d]: %+v", i, review)
			}
			if review.Reviewer == "" {
				t.Errorf("Expected Reviewer in Results[%d]: %+v", i, review)
			}
			if review.Score < 1 {
				t.Errorf("Expected Score is greater than zero in Results[%d]: %+v", i, review)
			}
			if review.Text == "" {
				t.Errorf("Expected review Text in Results[%d]: %+v", i, review)
			}
			if review.Timestamp.IsZero() {
				t.Errorf("Expected review Timestamp in Results[%d]: %+v", i, review)
			}
			if review.URL("") == "" {
				t.Errorf("Expected review URL in Results[%d]: %+v", i, review)
			}
		}
	}
}
