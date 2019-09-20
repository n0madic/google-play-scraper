# google-play-scraper

Golang scraper to get data from Google Play Store

This project is inspired by the [google-play-scraper](https://github.com/facundoolano/google-play-scraper) node.js project

## Instalation

```shell
go get -u github.com/n0madic/google-play-scraper/...
```

## Usage

### Get app details

Retrieves the full detail of an application.

```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/n0madic/google-play-scraper/pkg/app"
)

func main() {
	a := app.New("com.google.android.googlequicksearchbox")
	err := a.LoadDetails("ru", "us")
	if err != nil {
		panic(err)
	}
	spew.Dump(a)
}
```

### Search apps

Retrieves a list of apps that results of searching by the given term.

```go
package main

import (
	"fmt"

	"github.com/n0madic/google-play-scraper/pkg/search"
)

func main() {
	query := search.NewQuery(search.Options{
		Query:    "game",
		Country:  "ru",
		Language: "us",
		Number:   100,
		Discount: true,
		Price:    search.PricePaid,
		PriceMax: 100,
		ScoreMin: 4,
	})

	err := query.Do()
	if err != nil {
		panic(err)
	}

	errors := query.LoadMoreDetails(20)
	if len(errors) > 0 {
		panic(errors[0])
	}

	for _, app := range query.Results {
		if !app.IAPOffers {
			fmt.Println(app.Title, app.URL)
		}
	}
}
```
