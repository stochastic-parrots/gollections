package disjointset

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/disjointset"
)

// IntsRangeBySize partitions all integers in a fixed, inclusive range using
// union by size and path compression. Its zero value is invalid; use
// [NewIntsRangeBySize]. Size reports the number of values in a value's set.
type IntsRangeBySize[T constraint.Integer] = disjointset.FlatDisjointSetUnionBySize[T]

var _ SizeDisjointSet[int] = (*disjointset.FlatDisjointSetUnionBySize[int])(nil)

// NewIntsRangeBySize creates singleton sets for every integer from min through
// max, inclusive. Choose it when callers need the size of a value's set. It
// panics when min > max or the interval length cannot fit in an int. Allocation
// may also fail if the requested interval exceeds available memory.
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
//	Size(x)            O(α(N)) Amortized
//
// Complexity: O(N).
func NewIntsRangeBySize[T constraint.Integer](min, max T) *IntsRangeBySize[T] {
	return disjointset.NewFlatDisjointSetUnionBySize(min, max)
}
