package app

import (
	"net/url"
	"testing"
)

func TestLoadDetails(t *testing.T) {
	app := New("com.nekki.vector.paid", Options{"us", "en"})
	err := app.LoadDetails()
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
	if !app.Available {
		t.Error("Expected Available is true, got", app.Available)
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
	// if app.FamilyGenre == "" {
	// 	t.Error("Expected FamilyGenre")
	// }
	// if app.FamilyGenreID == "" {
	// 	t.Error("Expected FamilyGenreID")
	// }
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
	if app.InstallsMax == 0 {
		t.Error("Expected InstallsMax is greater than zero")
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
		t.Error("Expected RatingsHistogram length if 5, got", len(app.RatingsHistogram))
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
	if len(app.Reviews) == 0 {
		t.Error("Expected Reviews length is greater than zero")
	} else {
		for i, comment := range app.Reviews {
			if comment.Text == "" {
				t.Errorf("Expected comment Text in Reviews[%d]: %+v", i, comment)
			}
		}
	}
	if app.ReviewsTotalCount == 0 {
		t.Error("Expected ReviewsTotalCount is greater than zero")
	}
	if app.Score == 0 {
		t.Error("Expected Score is greater than zero")
	}
	if app.ScoreText == "" {
		t.Error("Expected ScoreText")
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

func TestLoadPermissions(t *testing.T) {
	app := New("com.android.chrome", Options{"us", "us"})
	err := app.LoadPermissions()
	if err != nil {
		t.Error(err)
	}

	if len(app.Permissions) == 0 {
		t.Fatal("Expected Permissions map length is greater than zero")
	}

	for key, permission := range app.Permissions {
		if key == "" {
			t.Error("Expected permission key is not empty")
		}

		if len(permission) == 0 {
			t.Fatal("Expected permission list length is greater than zero")
		}

		for _, perm := range permission {
			if perm == "" {
				t.Error("Expected permission in the list is not empty")
			}
		}

	}
}
