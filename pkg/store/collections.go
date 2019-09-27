package store

// Collection of apps
type Collection string

const (
	// TopFree apps
	TopFree Collection = "topselling_free"
	// TopPaid apps
	TopPaid Collection = "topselling_paid"
	// NewFree apps
	NewFree Collection = "topselling_new_free"
	// NewPaid apps
	NewPaid Collection = "topselling_new_paid"
	// Grossing apps
	Grossing Collection = "topgrossing"
	// Trending apps
	Trending Collection = "movers_shakers"
)
