package store

// Sort of apps
type Sort int

const (
	// SortHelpfulness method
	SortHelpfulness Sort = iota + 1
	// SortNewest method
	SortNewest
	// SortRating method
	SortRating
)
