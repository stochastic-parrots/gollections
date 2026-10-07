package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// KeyedHashSet is a map-backed Set using a derived comparable key. Its zero value
// is invalid; use NewKeyedHashSet. Derived keys must have reflexive equality.
// A KeyedHashSet value must not be copied after first use; share its pointer instead.
type KeyedHashSet[T any, K comparable] = set.KeyedHashSet[T, K]

var (
	_ Set[any]                              = &set.KeyedHashSet[any, int]{}
	_ Algebra[any, *KeyedHashSet[any, int]] = &set.KeyedHashSet[any, int]{}
	_ InPlaceAlgebra[any]                   = &set.KeyedHashSet[any, int]{}
)

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
//	Add(xs... T)       O(len(xs)) on average + len(xs) keyOf calls
//	Remove(xs... T)    O(len(xs)) on average + len(xs) keyOf calls
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
//	Add(xs... T)       O(len(xs)) on average + len(xs) keyOf calls
//	Remove(xs... T)    O(len(xs)) on average + len(xs) keyOf calls
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
//	Add(xs... T)       O(len(xs)) on average + len(xs) keyOf calls
//	Remove(xs... T)    O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(N) on average plus N keyOf calls, where N is the number of values yielded by seq.
func KeyedHashSetFromSeq[T any, K comparable](keyOf func(T) K, seq iter.Seq[T]) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSetFromSeq(seq, keyOf)
}

// KeyedHashSetFromMap uses values as the set's backing map without copying.
// The caller must ensure that every map key equals keyOf(value). The set and
// caller share the map, so mutations through either are visible to both.
// keyOf must be non-nil and remain stable while values belong to the set. It
// panics if keyOf is nil.
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
//	Add(xs... T)       O(len(xs)) on average + len(xs) keyOf calls
//	Remove(xs... T)    O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(1).
func KeyedHashSetFromMap[T any, K comparable](keyOf func(T) K, values map[K]T) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSetFromMap(values, keyOf)
}

// KeyedHashSetCloneMap creates a keyed hash set with a copy of values.
// The caller must ensure that every map key equals keyOf(value). The set does
// not retain the input map. keyOf must be non-nil and remain stable while
// values belong to the set. It panics if keyOf is nil.
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
//	Add(xs... T)       O(len(xs)) on average + len(xs) keyOf calls
//	Remove(xs... T)    O(len(xs)) on average + len(xs) keyOf calls
//	Clear()            O(N)
//	MarshalJSON()      O(N)
//
// Complexity: O(len(values)).
func KeyedHashSetCloneMap[T any, K comparable](keyOf func(T) K, values map[K]T) *KeyedHashSet[T, K] {
	checkKeyOf(keyOf)
	return set.NewKeyedHashSetCloneMap(values, keyOf)
}
