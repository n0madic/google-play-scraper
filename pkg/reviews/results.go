package reviews

// Results of operation
type Results []*Review

// Append result
func (results *Results) Append(res ...Review) {
	for _, result := range res {
		if !results.searchDuplicate(result.ID) {
			results.append(result)
		}
	}
}

func (results *Results) append(result Review) {
	*results = append(*results, &result)
}

func (results *Results) searchDuplicate(id string) bool {
	for _, result := range *results {
		if id == result.ID {
			return true
		}
	}
	return false
}
