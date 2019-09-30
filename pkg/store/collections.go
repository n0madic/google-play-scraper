package store

// Collection of apps
type Collection string

const (
	// TopFree apps
	TopFree Collection = "topselling_free"
	// TopPaid apps
	TopPaid Collection = "topselling_paid"
	// TopNewFree apps
	TopNewFree Collection = "topselling_new_free"
	// TopNewPaid apps
	TopNewPaid Collection = "topselling_new_paid"
	// TopGrossing apps
	TopGrossing Collection = "topgrossing"
	// TopTrending apps
	TopTrending Collection = "movers_shakers"
)
