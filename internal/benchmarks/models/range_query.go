package models

import "math/rand/v2"

// RangeQuery represents a half-open integer range [From, To).
type RangeQuery struct {
	From int
	To   int
}

// NewRandomRangeQueries creates random half-open integer range queries.
func NewRandomRangeQueries(size, maxValue, maxWidth int) []RangeQuery {
	queries := make([]RangeQuery, size)
	for idx := range queries {
		from := rand.IntN(maxValue)
		to := min(from+1+rand.IntN(maxWidth), maxValue)
		queries[idx] = RangeQuery{From: from, To: to}
	}
	return queries
}
