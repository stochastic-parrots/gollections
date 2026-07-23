package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// HashSet is a map-backed [Set] using Go equality for value identity. Values
// must have reflexive equality; in particular, floating-point NaN values cannot
// be found after insertion. When T is an interface, every dynamic value used as
// a key must also be comparable. Use [HashSetBy] to canonicalize such values.
//
// HashSet's zero value is ready for use.
type HashSet[T comparable] = set.HashSet[T]

var _ Set[int] = &set.HashSet[int]{}

// HashFactory constructs hash sets using Go equality for value identity.
//
// The zero value is ready for use.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation               Time Complexity
//	--------------------    ---------------
//	New(capacity)           O(capacity)
//	From/FromSeq            O(N)
//	Add/Contains/Remove     O(1)
//	Adds/Removes(xs...T)    O(len(xs))
//	All/Enumerate/Clear     O(N)
type HashFactory[T comparable] struct{}

// HashSetOf returns a factory for hash sets using Go equality.
func HashSetOf[T comparable]() HashFactory[T] {
	return HashFactory[T]{}
}

// New creates an empty hash set with the requested initial capacity.
func (HashFactory[T]) New(capacity int) *HashSet[T] {
	return set.NewHashSet[T](capacity)
}

// From creates a hash set containing the unique values in data.
//
// From copies values into map storage and does not retain data.
func (HashFactory[T]) From(data []T) *HashSet[T] {
	return set.NewHashSetFromSlice(data)
}

// FromSeq collects the unique values yielded by seq into a new hash set.
func (HashFactory[T]) FromSeq(seq iter.Seq[T]) *HashSet[T] {
	return set.NewHashSetFromSeq(seq)
}
