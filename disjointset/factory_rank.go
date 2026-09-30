package disjointset

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/disjointset"
)

// IntsRangeByRank partitions all integers in a fixed, inclusive range using
// union by rank and path compression. Its zero value is invalid; use
// [NewIntsRangeByRank]. Rank reports the upper-bound tree height heuristic for
// the set containing a value, not the number of values in that set.
type IntsRangeByRank[T constraint.Integer] = disjointset.FlatDisjointSetUnionByRank[T]

var _ RankDisjointSet[int] = (*disjointset.FlatDisjointSetUnionByRank[int])(nil)

// NewIntsRangeByRank creates singleton sets for every integer from min through
// max, inclusive. Choose it when union by rank is preferred and set sizes are
// not needed. It panics when min > max or the interval length cannot fit in an
// int. Allocation may also fail if the requested interval exceeds available
// memory.
// N is the number of values in the range, and α is the inverse Ackermann
// function. The amortized bounds apply to a sequence of operations between
// resets; an individual path traversal can take O(log N) time.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   -----------------
//	IsEmpty()          O(1)
//	Len()              O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Find(x)            O(α(N)) Amortized
//	Connected(a, b)    O(α(N)) Amortized
//	Disjoints()        O(1)
//	Union(a, b)        O(α(N)) Amortized
//	Reset()            O(N)
//	Rank(x)            O(α(N)) Amortized
//
// Complexity: O(N).
func NewIntsRangeByRank[T constraint.Integer](min, max T) *IntsRangeByRank[T] {
	return disjointset.NewFlatDisjointSetUnionByRank(min, max)
}
