package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// KeyedHashSet is a map-backed Set using a derived comparable key. Its zero value
// is invalid; use NewKeyedHashSet. Derived keys must have reflexive equality.
type KeyedHashSet[T any, K comparable] = set.KeyedHashSet[T, K]

var _ Set[any] = &set.KeyedHashSet[any, int]{}

func checkKeyOf[T any, K comparable](keyOf func(T) K) {
	if keyOf == nil {
		panic("set: nil identity function")
	}
}

// NewKeyedHashSet creates an empty keyed set. keyOf must be non-nil and stable
// for every stored value's lifetime in the set.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Len()              O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average + keyOf
//	Add(x)             O(1) on average + keyOf
//	Adds(xs... T)      O(len(xs)) on average + len(xs) keyOf calls
//	Remove(x)          O(1) on average + keyOf
//	Removes(xs... T)   O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(capacity).
func NewKeyedHashSet[T any, K comparable](keyOf func(T) K, capacity int) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSet(capacity, keyOf)
}

// KeyedHashSetFrom copies values from data into map storage without retaining the slice.
// The first value for each derived key is preserved. keyOf must be non-nil and stable.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Len()              O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average + keyOf
//	Add(x)             O(1) on average + keyOf
//	Adds(xs... T)      O(len(xs)) on average + len(xs) keyOf calls
//	Remove(x)          O(1) on average + keyOf
//	Removes(xs... T)   O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(len(data)) on average plus len(data) keyOf calls.
func KeyedHashSetFrom[T any, K comparable](keyOf func(T) K, data []T) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSetFromSlice(data, keyOf)
}

// KeyedHashSetFromSeq consumes seq once into map storage without guaranteeing iteration order. The first value for each
// derived key is preserved. keyOf must be non-nil and stable.
//
// Performance Summary (Time Complexity):
//
//	Operation          Time Complexity
//	----------------   ---------------
//	IsEmpty()          O(1)
//	Len()              O(1)
//	All()              O(N)
//	Enumerate()        O(N)
//	Contains(x)        O(1) on average + keyOf
//	Add(x)             O(1) on average + keyOf
//	Adds(xs... T)      O(len(xs)) on average + len(xs) keyOf calls
//	Remove(x)          O(1) on average + keyOf
//	Removes(xs... T)   O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(N) on average plus N keyOf calls, where N is the number of values yielded by seq.
func KeyedHashSetFromSeq[T any, K comparable](keyOf func(T) K, seq iter.Seq[T]) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSetFromSeq(seq, keyOf)
}
