package algorithms

import (
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

// RangeQueryCount counts values matched by half-open range queries.
func RangeQueryCount(list datastructs.SortedList, queries []models.RangeQuery) int {
	count := 0
	for _, query := range queries {
		count += list.LowerBound(query.To) - list.LowerBound(query.From)
	}
	return count
}
