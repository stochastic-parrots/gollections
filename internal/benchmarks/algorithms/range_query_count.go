package algorithms

import (
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
)

// RangeQueryCount counts values matched by one half-open range query.
func RangeQueryCount(list datastructs.SortedList, from, to int) int {
	return list.LowerBound(to) - list.LowerBound(from)
}
