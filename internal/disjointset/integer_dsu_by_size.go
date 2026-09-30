package disjointset

import (
	"iter"

	"github.com/stochastic-parrots/gollections/constraint"
)

// FlatDisjointSetUnionBySize stores one parent and size per value in an
// inclusive integer range. Sizes are authoritative only at roots.
// Find, Connected, Union, and Size take O(α(N)) amortized time between resets,
// where N is the range length and α is the inverse Ackermann function. A single
// call can traverse O(log N) nodes.
type FlatDisjointSetUnionBySize[T constraint.Integer] struct {
	min, max  T
	disjoints int
	parents   []int
	sizes     []uint
}

// NewFlatDisjointSetUnionBySize allocates singleton sets for the range.
func NewFlatDisjointSetUnionBySize[T constraint.Integer](min, max T) *FlatDisjointSetUnionBySize[T] {
	interval := flatInterval(min, max)
	parents, sizes := make([]int, interval), make([]uint, interval)
	disjoints := interval

	for idx := range interval {
		parents[idx], sizes[idx] = idx, 1
	}
	return &FlatDisjointSetUnionBySize[T]{min, max, disjoints, parents, sizes}
}

func (dsu *FlatDisjointSetUnionBySize[T]) normalize(x T) (int, bool) {
	if x < dsu.min || x > dsu.max {
		return 0, false
	}
	return flatIndex(x, dsu.min), true
}

func (dsu *FlatDisjointSetUnionBySize[T]) find(normalized int) int {
	if dsu.parents[normalized] != normalized {
		dsu.parents[normalized] = dsu.find(dsu.parents[normalized])
	}
	return dsu.parents[normalized]
}

func (dsu *FlatDisjointSetUnionBySize[T]) link(a, b int) {
	if dsu.sizes[b] > dsu.sizes[a] {
		dsu.parents[a] = b
		dsu.sizes[b] += dsu.sizes[a]
	} else {
		dsu.parents[b] = a
		dsu.sizes[a] += dsu.sizes[b]
	}
}

// Find returns x's current representative and compresses its path. Invalid
// values return the zero value and false.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionBySize[T]) Find(x T) (T, bool) {
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
func (dsu *FlatDisjointSetUnionBySize[T]) Connected(a, b T) (connected, ok bool) {
	normA, okA := dsu.normalize(a)
	normB, okB := dsu.normalize(b)
	if !okA || !okB {
		return false, false
	}
	return dsu.find(normA) == dsu.find(normB), true
}

// Union links distinct roots by size. It returns false for invalid values or
// values already in the same set.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionBySize[T]) Union(a, b T) bool {
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

// Size returns the size of x's set, or zero and false for an invalid value.
//
// Complexity: O(α(N)) Amortized.
func (dsu *FlatDisjointSetUnionBySize[T]) Size(x T) (size int, ok bool) {
	normalized, ok := dsu.normalize(x)
	if !ok {
		return 0, false
	}
	root := dsu.find(normalized)
	return int(dsu.sizes[root]), true
}

// Disjoints returns the current number of sets.
func (dsu *FlatDisjointSetUnionBySize[T]) Disjoints() int {
	return dsu.disjoints
}

// Len returns the constant number of values in the range.
func (dsu *FlatDisjointSetUnionBySize[T]) Len() int { return len(dsu.parents) }

// IsEmpty is false for every constructed instance because the range is nonempty.
func (dsu *FlatDisjointSetUnionBySize[T]) IsEmpty() bool { return len(dsu.parents) == 0 }

// All lazily yields the values in increasing order, regardless of unions.
//
// Complexity: O(N) for a full traversal.
func (dsu *FlatDisjointSetUnionBySize[T]) All() iter.Seq[T] {
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
func (dsu *FlatDisjointSetUnionBySize[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx := range dsu.parents {
			if !yield(idx, flatValue(idx, dsu.min)) {
				return
			}
		}
	}
}

// Reset restores singleton sets and reuses the allocated parent and size slices.
//
// Complexity: O(N).
func (dsu *FlatDisjointSetUnionBySize[T]) Reset() {
	for idx := range dsu.parents {
		dsu.parents[idx] = idx
		dsu.sizes[idx] = 1
	}
	dsu.disjoints = len(dsu.parents)
}
