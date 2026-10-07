package disjointset

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/constraint"
)

// FlatDisjointSetUnionBySize stores an inclusive integer range in one slice.
// Nonnegative entries are parent indexes; negative entries are root sizes.
// Find, Connected, Union, and Size take O(α(N)) amortized time between resets,
// where N is the range length and α is the inverse Ackermann function. A single
// call can traverse O(log N) nodes.
type FlatDisjointSetUnionBySize[T constraint.Integer] struct {
	min, max  T
	disjoints int
	parents   []int
}

// NewFlatDisjointSetUnionBySize allocates singleton sets for the range.
func NewFlatDisjointSetUnionBySize[T constraint.Integer](min, max T) *FlatDisjointSetUnionBySize[T] {
	interval := flatInterval(min, max)
	parents := make([]int, interval)
	disjoints := interval

	for idx := range interval {
		parents[idx] = -1
	}
	return &FlatDisjointSetUnionBySize[T]{min, max, disjoints, parents}
}

func (dsu *FlatDisjointSetUnionBySize[T]) normalize(x T) (int, bool) {
	if x < dsu.min || x > dsu.max {
		return 0, false
	}
	return flatIndex(x, dsu.min), true
}

func (dsu *FlatDisjointSetUnionBySize[T]) find(normalized int) int {
	parent := dsu.parents[normalized]
	if parent < 0 {
		return normalized
	}
	root := dsu.find(parent)
	dsu.parents[normalized] = root
	return root
}

func (dsu *FlatDisjointSetUnionBySize[T]) link(a, b int) {
	if dsu.parents[b] < dsu.parents[a] {
		dsu.parents[b] += dsu.parents[a]
		dsu.parents[a] = b
	} else {
		dsu.parents[a] += dsu.parents[b]
		dsu.parents[b] = a
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
	return -dsu.parents[root], true
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

// Reset restores singleton sets and reuses the allocated parent slice.
//
// Complexity: O(N).
func (dsu *FlatDisjointSetUnionBySize[T]) Reset() {
	for idx := range dsu.parents {
		dsu.parents[idx] = -1
	}
	dsu.disjoints = len(dsu.parents)
}
