package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// KeyedHashSet is a map-backed [Set] for arbitrary values whose identity is a
// derived comparable key. Derived keys must have reflexive equality;
// floating-point NaN is therefore not a valid key without canonicalization.
//
// KeyedHashSet's zero value is invalid; construct one with [HashSetBy].
type KeyedHashSet[T any, K comparable] = set.KeyedHashSet[T, K]

var _ Set[any] = &set.KeyedHashSet[any, int]{}

// HashSetByFactory constructs hash sets using a derived comparable key for value
// identity.
//
// The identity function must be non-nil and remain stable for the lifetime of
// every set created by the factory. If mutable values are stored, the key
// derived from a value must not change while that value belongs to a set. Keys
// must have reflexive equality, so floating-point NaN values must be
// canonicalized before they are returned as keys.
//
// The zero value is invalid. Create a factory with [HashSetBy].
//
// Performance Summary (Expected Time Complexity):
//
//	Operation               Time Complexity
//	--------------------    ---------------
//	New(capacity)           O(capacity)
//	From/FromSeq            O(N)
//	Add/Contains/Remove     O(1) + keyOf
//	Adds/Removes(xs...T)    O(len(xs)) + keyOf calls
//	All/Enumerate/Clear     O(N)
type HashSetByFactory[T any, K comparable] struct {
	keyOf func(T) K
}

// HashSetBy returns a factory that identifies values by the key returned by keyOf.
// It panics if keyOf is nil.
func HashSetBy[T any, K comparable](keyOf func(T) K) HashSetByFactory[T, K] {
	if keyOf == nil {
		panic("set: nil identity function")
	}
	return HashSetByFactory[T, K]{keyOf: keyOf}
}

// New creates an empty keyed hash set with the requested initial capacity.
func (factory HashSetByFactory[T, K]) New(capacity int) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSet(capacity, factory.keyOf)
}

// From creates a keyed hash set containing values from data. The first value
// encountered for each derived key is preserved.
//
// From copies values into map storage and does not retain data.
func (factory HashSetByFactory[T, K]) From(data []T) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetFromSlice(data, factory.keyOf)
}

// FromSeq collects values yielded by seq into a new keyed hash set. The first
// value yielded for each derived key is preserved.
func (factory HashSetByFactory[T, K]) FromSeq(seq iter.Seq[T]) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetFromSeq(seq, factory.keyOf)
}
