package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// HashSet stores unique comparable values as keys in a Go map.
type HashSet[T comparable] struct {
	values map[T]struct{}
}

// NewHashSet creates an empty hash set with the requested initial capacity.
func NewHashSet[T comparable](capacity int) *HashSet[T] {
	return &HashSet[T]{values: make(map[T]struct{}, capacity)}
}

// NewHashSetFromSlice creates a hash set containing the unique values in data.
func NewHashSetFromSlice[T comparable](data []T) *HashSet[T] {
	set := NewHashSet[T](len(data))
	set.Add(data...)
	return set
}

// NewHashSetFromSeq collects the unique values yielded by seq into a hash set.
func NewHashSetFromSeq[T comparable](seq iter.Seq[T]) *HashSet[T] {
	set := NewHashSet[T](0)
	for value := range seq {
		set.Add(value)
	}
	return set
}

// Len returns the number of values in the set.
//
// Complexity: O(1).
func (set *HashSet[T]) Len() int {
	return len(set.values)
}

// IsEmpty returns true if the set contains no values.
//
// Complexity: O(1).
func (set *HashSet[T]) IsEmpty() bool {
	return len(set.values) == 0
}

// Contains returns true if x belongs to the set.
//
// Complexity: O(1) on average.
func (set *HashSet[T]) Contains(x T) bool {
	_, ok := set.values[x]
	return ok
}

// Add inserts every value in xs that is not already present and returns the
// number of values added. Empty input leaves the set unchanged.
//
// Complexity: O(len(xs)) on average.
func (set *HashSet[T]) Add(xs ...T) (added int) {
	for _, x := range xs {
		if _, ok := set.values[x]; ok {
			continue
		}
		if set.values == nil {
			set.values = make(map[T]struct{})
		}
		set.values[x] = struct{}{}
		added++
	}
	return added
}

// Remove deletes every value in xs that is present and returns the number of
// values removed. Empty input leaves the set unchanged.
//
// Complexity: O(len(xs)) on average.
func (set *HashSet[T]) Remove(xs ...T) (removed int) {
	for _, x := range xs {
		if _, ok := set.values[x]; !ok {
			continue
		}
		delete(set.values, x)
		removed++
	}
	return removed
}

// All returns a lazy iterator over the set in unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *HashSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range set.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate returns a lazy iterator over iteration positions and values in
// unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *HashSet[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for value := range set.values {
			if !yield(idx, value) {
				return
			}
			idx++
		}
	}
}

// Clear removes all values while keeping the set initialized for reuse.
//
// Complexity: O(N).
func (set *HashSet[T]) Clear() {
	clear(set.values)
}

// MarshalJSON encodes the set as a JSON array in unspecified iteration order.
//
// Complexity: O(N).
func (set *HashSet[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(set)
}
