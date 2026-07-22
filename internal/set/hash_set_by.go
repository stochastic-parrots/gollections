package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// KeyedHashSet stores arbitrary values indexed by a comparable key derived
// from each value. The first value inserted for a key remains its
// representative.
type KeyedHashSet[T any, K comparable] struct {
	values map[K]T
	keyOf  func(T) K
}

// NewKeyedHashSet creates an empty keyed hash set with the requested initial
// capacity and identity function.
func NewKeyedHashSet[T any, K comparable](capacity int, keyOf func(T) K) *KeyedHashSet[T, K] {
	return &KeyedHashSet[T, K]{
		values: make(map[K]T, capacity),
		keyOf:  keyOf,
	}
}

// NewKeyedHashSetFromSlice creates a keyed hash set containing values from
// data. The first value for each derived key is preserved.
func NewKeyedHashSetFromSlice[T any, K comparable](data []T, keyOf func(T) K) *KeyedHashSet[T, K] {
	set := NewKeyedHashSet(len(data), keyOf)
	set.Adds(data...)
	return set
}

// NewKeyedHashSetFromSeq collects values yielded by seq into a keyed hash set.
// The first value for each derived key is preserved.
func NewKeyedHashSetFromSeq[T any, K comparable](seq iter.Seq[T], keyOf func(T) K) *KeyedHashSet[T, K] {
	set := NewKeyedHashSet(0, keyOf)
	for value := range seq {
		set.Add(value)
	}
	return set
}

// Length returns the number of distinct keys in the set.
//
// Complexity: O(1).
func (set *KeyedHashSet[T, K]) Length() int {
	return len(set.values)
}

// IsEmpty returns true if the set contains no values.
//
// Complexity: O(1).
func (set *KeyedHashSet[T, K]) IsEmpty() bool {
	return len(set.values) == 0
}

// Contains returns true if the key derived from x belongs to the set.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Contains(x T) bool {
	_, ok := set.values[set.keyOf(x)]
	return ok
}

// Add inserts x unless its derived key is already present, reports whether the
// set changed, and preserves existing representatives.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Add(x T) bool {
	key := set.keyOf(x)
	if _, ok := set.values[key]; ok {
		return false
	}
	if set.values == nil {
		set.values = make(map[K]T)
	}
	set.values[key] = x
	return true
}

// Adds inserts every value in xs whose derived key is not already present and
// returns the number of values added.
//
// Complexity: Expected O(len(xs)) plus the identity-function calls.
func (set *KeyedHashSet[T, K]) Adds(xs ...T) (added int) {
	for _, x := range xs {
		if set.Add(x) {
			added++
		}
	}
	return added
}

// Remove deletes the value with the same derived key as x and reports whether
// it was present.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Remove(x T) bool {
	key := set.keyOf(x)
	if _, ok := set.values[key]; !ok {
		return false
	}
	delete(set.values, key)
	return true
}

// Removes deletes every value in xs whose derived key is present and returns
// the number of values removed.
//
// Complexity: Expected O(len(xs)) plus the identity-function calls.
func (set *KeyedHashSet[T, K]) Removes(xs ...T) (removed int) {
	for _, x := range xs {
		if set.Remove(x) {
			removed++
		}
	}
	return removed
}

// All returns a lazy iterator over the representative values in unspecified
// order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *KeyedHashSet[T, K]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range set.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate returns a lazy iterator over iteration positions and representative
// values in unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *KeyedHashSet[T, K]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for _, value := range set.values {
			if !yield(idx, value) {
				return
			}
			idx++
		}
	}
}

// Clear removes all values while preserving the identity function and keeping
// the set initialized for reuse.
//
// Complexity: O(N).
func (set *KeyedHashSet[T, K]) Clear() {
	clear(set.values)
}

// MarshalJSON encodes the representative values as a JSON array in unspecified
// iteration order.
//
// Complexity: O(N).
func (set *KeyedHashSet[T, K]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(set)
}
