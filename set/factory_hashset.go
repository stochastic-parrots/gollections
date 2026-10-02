package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// HashSet is a map-backed [Set] using Go equality for value identity. Values
// must have reflexive equality; in particular, floating-point NaN values cannot
// be found after insertion. When T is an interface, every dynamic value used as
// a key must also be comparable. Use [NewKeyedHashSet] to canonicalize such
// values.
//
// HashSet's zero value is ready for use. A HashSet value must not be copied
// after first use; share its pointer instead.
type HashSet[T comparable] = set.HashSet[T]

var (
	_ Set[int]                    = &set.HashSet[int]{}
	_ Algebra[int, *HashSet[int]] = &set.HashSet[int]{}
	_ InPlaceAlgebra[int]         = &set.HashSet[int]{}
)

// NewHashSet creates an empty hash set with the requested initial capacity.
// Choose it when Go equality defines value identity and an initial capacity is
// useful.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Length()           O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average
//	Add(x)             O(1) on average
//	Adds(xs... T)      O(len(xs)) on average
//	Remove(x)          O(1) on average
//	Removes(xs... T)   O(len(xs)) on average
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(capacity).
func NewHashSet[T comparable](capacity int) *HashSet[T] {
	return set.NewHashSet[T](capacity)
}

// HashSetFrom creates a hash set containing the unique values in data.
//
// HashSetFrom copies values into map storage and does not retain data.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Length()           O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average
//	Add(x)             O(1) on average
//	Adds(xs... T)      O(len(xs)) on average
//	Remove(x)          O(1) on average
//	Removes(xs... T)   O(len(xs)) on average
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(len(data)) on average.
func HashSetFrom[T comparable](data []T) *HashSet[T] {
	return set.NewHashSetFromSlice(data)
}

// HashSetFromSeq collects the unique values yielded by seq into a new hash set.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Length()           O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average
//	Add(x)             O(1) on average
//	Adds(xs... T)      O(len(xs)) on average
//	Remove(x)          O(1) on average
//	Removes(xs... T)   O(len(xs)) on average
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(N) on average, where N is the number of values yielded by seq.
func HashSetFromSeq[T comparable](seq iter.Seq[T]) *HashSet[T] {
	return set.NewHashSetFromSeq(seq)
}
