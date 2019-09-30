package app

import (
	"net/url"
	"testing"
)

func TestLoadDetails(t *testing.T) {
	app := New("com.disney.WMW")
	err := app.LoadDetails("", "")
	if err != nil {
		t.Error(err)
	}

	if !app.AdSupported {
		t.Error("Expected AdSupported is true, got", app.AdSupported)
	}
	if app.AndroidVersion == "" {
		t.Error("Expected Android version")
	}
	if app.AndroidVersionMin == 0 {
		t.Error("Expected AndroidVersionMin is greater than zero")
	}
	if len(app.Comments) == 0 {
		t.Error("Expected Comments length is greater than zero")
	} else {
		for i, comment := range app.Comments {
			if comment.Avatar == "" {
				t.Errorf("Expected commentator Avatar in Comments[%d]: %+v", i, comment)
			}
			if comment.Commentator == "" {
				t.Errorf("Expected Commentator in Comments[%d]: %+v", i, comment)
			}
			if comment.Rating < 1 {
				t.Errorf("Expected Rating is greater than zero in Comments[%d]: %+v", i, comment)
			}
			if comment.Text == "" {
				t.Errorf("Expected comment Text in Comments[%d]: %+v", i, comment)
			}
			if comment.Timestamp.IsZero() {
				t.Errorf("Expected comment Timestamp in Comments[%d]: %+v", i, comment)
			}
		}
	}
	if app.ContentRating == "" {
		t.Error("Expected ContentRating")
	}
	if app.Description == "" {
		t.Error("Expected Description")
	}
	if app.DescriptionHTML == "" {
		t.Error("Expected DescriptionHTML")
	}
	if app.Developer == "" {
		t.Error("Expected Developer")
	}
	if app.DeveloperAddress == "" {
		t.Error("Expected DeveloperAddress")
	}
	if app.DeveloperEmail == "" {
		t.Error("Expected DeveloperEmail")
	}
	if app.DeveloperID == "" {
		t.Error("Expected DeveloperID")
	}
	if _, err = url.ParseRequestURI(app.DeveloperURL); err != nil {
		t.Error("Expected valid DeveloperURL, got", app.DeveloperURL)
	}
	if _, err = url.ParseRequestURI(app.DeveloperWebsite); err != nil {
		t.Error("Expected valid DeveloperWebsite, got", app.DeveloperWebsite)
	}
	if app.FamilyGenre == "" {
		t.Error("Expected FamilyGenre")
	}
	if app.FamilyGenreID == "" {
		t.Error("Expected FamilyGenreID")
	}
	if app.Genre == "" {
		t.Error("Expected Genre")
	}
	if app.GenreID == "" {
		t.Error("Expected GenreID")
	}
	if app.HeaderImage == "" {
		t.Error("Expected HeaderImage")
	}
	if !app.IAPOffers {
		t.Error("Expected IAPOffers is true, got", app.IAPOffers)
	}
	if app.IAPRange == "" {
		t.Error("Expected IAPRange")
	}
	if _, err = url.ParseRequestURI(app.Icon); err != nil {
		t.Error("Expected valid Icon url, got", app.Icon)
	}
	if app.ID == "" {
		t.Error("Expected ID")
	}
	if app.Installs == "" {
		t.Error("Expected Installs")
	}
	if app.InstallsMin == 0 {
		t.Error("Expected InstallsMin is greater than zero")
	}
	if app.Price.Currency == "" {
		t.Error("Expected Price.Currency")
	}
	if app.Price.Value == 0 {
		t.Error("Expected Price.Value is greater than zero")
	}
	if _, err = url.ParseRequestURI(app.Icon); err != nil {
		t.Error("Expected valid Icon url, got", app.Icon)
	}
	if app.PrivacyPolicy == "" {
		t.Error("Expected PrivacyPolicy")
	}
	if app.Ratings == 0 {
		t.Error("Expected Ratings is greater than zero")
	}
	if len(app.RatingsHistogram) != 5 {
		t.Error("Expected RatingsHistogram lenght if 5, got", len(app.RatingsHistogram))
	}
	for i := 1; i <= 5; i++ {
		if val, ok := app.RatingsHistogram[i]; ok {
			if val == 0 {
				t.Errorf("Expected RatingsHistogram[%d] is greater than zero", i)
			}
		} else {
			t.Error("Expected RatingsHistogram with key", i)
		}
	}
	if app.RecentChanges == "" {
		t.Error("Expected RecentChanges")
	}
	if app.RecentChangesHTML == "" {
		t.Error("Expected RecentChangesHTML")
	}
	if app.Released == "" {
		t.Error("Expected Released date")
	}
	if app.Reviews == "" {
		t.Error("Expected Reviews")
	}
	if app.Score == 0 {
		t.Error("Expected Score is greater than zero")
	}
	if len(app.Screenshots) == 0 {
		t.Error("Expected Screenshots length is greater than zero")
	} else {
		for i, screen := range app.Screenshots {
			if _, err = url.ParseRequestURI(screen); err != nil {
				t.Errorf("Expected valid Screenshots[%d] url, got %s", i, screen)
			}
		}
	}
	if _, err = url.ParseRequestURI(app.SimilarURL); err != nil {
		t.Error("Expected valid SimilarURL, got", app.SimilarURL)
	}
	if app.Size == "" {
		t.Error("Expected Size")
	}
	if app.Summary == "" {
		t.Error("Expected Summary")
	}
	if app.Title == "" {
		t.Error("Expected Title")
	}
	if app.Updated.IsZero() {
		t.Error("Expected Updated date")
	}
	if _, err = url.ParseRequestURI(app.URL); err != nil {
		t.Error("Expected valid URL, got", app.URL)
	}
	if app.Version == "" {
		t.Error("Expected Version")
	}
	if _, err = url.ParseRequestURI(app.Video); err != nil {
		t.Error("Expected valid Video url, got", app.Video)
	}
	if _, err = url.ParseRequestURI(app.VideoImage); err != nil {
		t.Error("Expected valid VideoImage url, got", app.VideoImage)
	}
}
