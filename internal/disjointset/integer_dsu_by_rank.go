package disjointset

import (
	"iter"

	"github.com/stochastic-parrots/gollections/constraint"
)

// FlatDisjointSetUnionByRank stores one parent and rank per value in an
// inclusive integer range. The parent slice index is the normalized value.
// Find, Connected, Union, and Rank take O(α(N)) amortized time between resets,
// where N is the range length and α is the inverse Ackermann function. A single
// call can traverse O(log N) nodes.
type FlatDisjointSetUnionByRank[T constraint.Integer] struct {
	min, max  T
	disjoints int
	parents   []int
	ranks     []uint8
}

// NewFlatDisjointSetUnionByRank allocates singleton sets for the range.
func NewFlatDisjointSetUnionByRank[T constraint.Integer](min, max T) *FlatDisjointSetUnionByRank[T] {
	interval := flatInterval(min, max)
	parents, ranks := make([]int, interval), make([]uint8, interval)
	disjoints := interval

	for idx := range interval {
		parents[idx] = idx
	}
	return &FlatDisjointSetUnionByRank[T]{min, max, disjoints, parents, ranks}
}

func (dsu *FlatDisjointSetUnionByRank[T]) normalize(x T) (int, bool) {
	if x < dsu.min || x > dsu.max {
		return 0, false
	}
	return flatIndex(x, dsu.min), true
}

func (dsu *FlatDisjointSetUnionByRank[T]) find(normalized int) int {
	if dsu.parents[normalized] != normalized {
		dsu.parents[normalized] = dsu.find(dsu.parents[normalized])
	}
	return dsu.parents[normalized]
}

func (dsu *FlatDisjointSetUnionByRank[T]) link(a, b int) {
	if dsu.ranks[a] > dsu.ranks[b] {
		dsu.parents[b] = a
	} else if dsu.ranks[b] > dsu.ranks[a] {
		dsu.parents[a] = b
	} else {
		dsu.parents[b] = a
		dsu.ranks[a] += 1
	}
}

// Find returns x's current representative and compresses its path. Invalid
// values return the zero value and false.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionByRank[T]) Find(x T) (T, bool) {
	normalized, ok := dsu.normalize(x)
	if !ok {
		var zero T
		return zero, false
	}

	root := dsu.find(normalized)
	return flatValue(root, dsu.min), true
}

// Connected reports membership in the same set, with ok false for invalid
// values. It compresses paths for valid values.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionByRank[T]) Connected(a, b T) (connected, ok bool) {
	normA, okA := dsu.normalize(a)
	normB, okB := dsu.normalize(b)
	if !okA || !okB {
		return false, false
	}
	return dsu.find(normA) == dsu.find(normB), true
}

// Union links distinct roots by rank. It returns false for invalid values or
// values already in the same set.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionByRank[T]) Union(a, b T) (ok bool) {
	normA, okA := dsu.normalize(a)
	normB, okB := dsu.normalize(b)

	if !okA || !okB {
		return false
	}

	rootA, rootB := dsu.find(normA), dsu.find(normB)

	if rootA == rootB {
		return false
	}

	dsu.link(rootA, rootB)
	dsu.disjoints--
	return true
}

// Rank returns the rank of x's root, or zero and false for an invalid value.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionByRank[T]) Rank(x T) (rank int, ok bool) {
	normalized, ok := dsu.normalize(x)
	if !ok {
		return 0, false
	}
	root := dsu.find(normalized)
	return int(dsu.ranks[root]), true
}

// Disjoints returns the current number of sets.
func (dsu *FlatDisjointSetUnionByRank[T]) Disjoints() int {
	return dsu.disjoints
}

// Len returns the constant number of values in the range.
func (dsu *FlatDisjointSetUnionByRank[T]) Len() int { return len(dsu.parents) }

// IsEmpty is false for every constructed instance because the range is nonempty.
func (dsu *FlatDisjointSetUnionByRank[T]) IsEmpty() bool { return len(dsu.parents) == 0 }

// All lazily yields the values in increasing order, regardless of unions.
//
// Complexity: O(N) for a full traversal.
func (dsu *FlatDisjointSetUnionByRank[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for idx := range dsu.parents {
			if !yield(flatValue(idx, dsu.min)) {
				return
			}
		}
	}
}

// Enumerate lazily yields each range offset and value in increasing order.
//
// Complexity: O(N) for a full traversal.
func (dsu *FlatDisjointSetUnionByRank[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx := range dsu.parents {
			if !yield(idx, flatValue(idx, dsu.min)) {
				return
			}
		}
	}
}

// Reset restores singleton sets and reuses the allocated parent and rank slices.
//
// Complexity: O(N).
func (dsu *FlatDisjointSetUnionByRank[T]) Reset() {
	for idx := range dsu.parents {
		dsu.parents[idx] = idx
		dsu.ranks[idx] = 0
	}
	dsu.disjoints = len(dsu.parents)
}
