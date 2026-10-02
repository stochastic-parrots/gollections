package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// KeyedHashSet is a map-backed [Set] for arbitrary values whose identity is a
// derived comparable key. Derived keys must have reflexive equality;
// floating-point NaN is therefore not a valid key without canonicalization.
//
// Its zero value is invalid. The identity function must remain stable while
// values belong to the set. A KeyedHashSet value must not be copied after first
// use; share its pointer instead.
type KeyedHashSet[K comparable, T any] = set.KeyedHashSet[K, T]

var (
	_ Set[any]                              = &set.KeyedHashSet[int, any]{}
	_ Algebra[any, *KeyedHashSet[int, any]] = &set.KeyedHashSet[int, any]{}
	_ InPlaceAlgebra[any]                   = &set.KeyedHashSet[int, any]{}
)

// NewKeyedHashSet creates an empty keyed hash set with the requested initial
// capacity. keyOf must be non-nil and remain stable while values belong to the
// set. It panics if keyOf is nil.
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
func NewKeyedHashSet[K comparable, T any](capacity int, keyOf func(T) K) *KeyedHashSet[K, T] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return set.NewKeyedHashSet(capacity, keyOf)
}

// KeyedHashSetFrom copies data into a new keyed hash set without retaining the
// source slice. The first value encountered for each derived key is preserved.
// keyOf must be non-nil and remain stable while values belong to the set. It
// panics if keyOf is nil.
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
func KeyedHashSetFrom[K comparable, T any](data []T, keyOf func(T) K) *KeyedHashSet[K, T] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return set.NewKeyedHashSetFromSlice(data, keyOf)
}

// KeyedHashSetFromSeq collects values yielded by seq into a new keyed hash set.
// The first value yielded for each derived key is preserved. keyOf must be
// non-nil and remain stable while values belong to the set. It panics if keyOf
// is nil.
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
func KeyedHashSetFromSeq[K comparable, T any](seq iter.Seq[T], keyOf func(T) K) *KeyedHashSet[K, T] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return set.NewKeyedHashSetFromSeq(seq, keyOf)
}

// KeyedHashSetFromMap uses values as the set's backing map without copying.
// The caller must ensure that every map key equals keyOf(value). The set and
// caller share the map, so mutations through either are visible to both.
// keyOf must be non-nil and remain stable while values belong to the set. It
// panics if keyOf is nil.
//
// Complexity: O(1).
func KeyedHashSetFromMap[K comparable, T any](values map[K]T, keyOf func(T) K) *KeyedHashSet[K, T] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return set.NewKeyedHashSetFromMap(values, keyOf)
}

// KeyedHashSetCloneMap creates a keyed hash set with a copy of values.
// The caller must ensure that every map key equals keyOf(value). The set does
// not retain the input map. keyOf must be non-nil and remain stable while
// values belong to the set. It panics if keyOf is nil.
//
// Complexity: O(len(values)).
func KeyedHashSetCloneMap[K comparable, T any](values map[K]T, keyOf func(T) K) *KeyedHashSet[K, T] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return set.NewKeyedHashSetCloneMap(values, keyOf)
}
