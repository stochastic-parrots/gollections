package set

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// KeyedHashSet is a map-backed [Set] for arbitrary values whose identity is a
// derived comparable key. Derived keys must have reflexive equality;
// floating-point NaN is therefore not a valid key without canonicalization.
//
// KeyedHashSet's zero value is invalid; construct one with [HashSetBy]. A
// KeyedHashSet value must not be copied after first use; share its pointer
// instead.
type KeyedHashSet[T any, K comparable] = set.KeyedHashSet[T, K]

var (
	_ Set[any]                              = &set.KeyedHashSet[any, int]{}
	_ Algebra[any, *KeyedHashSet[any, int]] = &set.KeyedHashSet[any, int]{}
	_ InPlaceAlgebra[any]                   = &set.KeyedHashSet[any, int]{}
)

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
func (factory HashSetByFactory[T, K]) New(capacity int) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSet(capacity, factory.keyOf)
}

// From creates a keyed hash set containing values from data. The first value
// encountered for each derived key is preserved.
//
// From copies values into map storage and does not retain data.
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
func (factory HashSetByFactory[T, K]) From(data []T) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetFromSlice(data, factory.keyOf)
}

// FromSeq collects values yielded by seq into a new keyed hash set. The first
// value yielded for each derived key is preserved.
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
func (factory HashSetByFactory[T, K]) FromSeq(seq iter.Seq[T]) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetFromSeq(seq, factory.keyOf)
}

// Union returns a new set containing keys present in any source. With no
// sources, it returns an empty set.
func (factory HashSetByFactory[T, K]) Union(
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	// A largest-source hint avoids reserving the sum of overlapping sources.
	capacity := 0
	for _, source := range sources {
		capacity = max(capacity, source.Length())
	}
	result := factory.New(capacity)
	result.UnionWith(sources...)
	return result
}

// Intersection returns a new set containing keys present in every source.
// Representatives come from the first source. With no sources, it returns an
// empty set.
func (factory HashSetByFactory[T, K]) Intersection(
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetIntersection(factory.keyOf, sources...)
}

// Difference returns a new set containing keys from the first source that are
// absent from every subsequent source. With no sources, it returns an empty
// set.
func (factory HashSetByFactory[T, K]) Difference(
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetDifference(factory.keyOf, sources...)
}

// SymmetricDifference returns a new set containing keys present in an odd
// number of sources. A surviving key uses the representative from the last
// source that toggles it into the result. With no sources, it returns an empty
// set.
func (factory HashSetByFactory[T, K]) SymmetricDifference(
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	return set.NewKeyedHashSetSymmetricDifference(factory.keyOf, sources...)
}

// Equal returns true if left and right contain the same derived keys.
func (factory HashSetByFactory[T, K]) Equal(
	left, right Source[T],
) bool {
	if left.Length() > right.Length() {
		left, right = right, left
	}
	membership := set.CollectKeyedMembership(left, factory.keyOf)
	matched, contained := set.MatchKeyedMembership(
		membership,
		right,
		factory.keyOf,
	)
	return contained && matched == len(membership)
}

// IsSubset returns true if every key in subset is present in superset.
func (factory HashSetByFactory[T, K]) IsSubset(
	subset, superset Source[T],
) bool {
	if subset.Length() <= superset.Length() {
		return set.MatchKeyedSubset(
			set.CollectKeyedMembership(subset, factory.keyOf),
			superset,
			factory.keyOf,
		)
	}
	keys := set.CollectKeyedKeys(superset, factory.keyOf)
	for value := range subset.All() {
		if _, exists := keys[factory.keyOf(value)]; !exists {
			return false
		}
	}
	return true
}

// IsProperSubset returns true if subset is contained in superset and the two
// sources are not equal under the factory's identity function.
func (factory HashSetByFactory[T, K]) IsProperSubset(
	subset, superset Source[T],
) bool {
	if subset.Length() <= superset.Length() {
		membership := set.CollectKeyedMembership(subset, factory.keyOf)
		matched, hasExtra := set.MatchKeyedContainedMembership(
			membership,
			superset,
			factory.keyOf,
		)
		return matched == len(membership) && hasExtra
	}
	membership := set.CollectKeyedMembership(superset, factory.keyOf)
	matched, contained := set.MatchKeyedMembership(
		membership,
		subset,
		factory.keyOf,
	)
	return contained && matched < len(membership)
}

// IsSuperset returns true if every key in subset is present in superset.
func (factory HashSetByFactory[T, K]) IsSuperset(
	superset, subset Source[T],
) bool {
	if subset.Length() <= superset.Length() {
		return set.MatchKeyedSubset(
			set.CollectKeyedMembership(subset, factory.keyOf),
			superset,
			factory.keyOf,
		)
	}
	keys := set.CollectKeyedKeys(superset, factory.keyOf)
	for value := range subset.All() {
		if _, exists := keys[factory.keyOf(value)]; !exists {
			return false
		}
	}
	return true
}

// IsProperSuperset returns true if superset contains subset and the two sources
// are not equal under the factory's identity function.
func (factory HashSetByFactory[T, K]) IsProperSuperset(
	superset, subset Source[T],
) bool {
	if subset.Length() <= superset.Length() {
		membership := set.CollectKeyedMembership(subset, factory.keyOf)
		matched, hasExtra := set.MatchKeyedContainedMembership(
			membership,
			superset,
			factory.keyOf,
		)
		return matched == len(membership) && hasExtra
	}
	membership := set.CollectKeyedMembership(superset, factory.keyOf)
	matched, contained := set.MatchKeyedMembership(
		membership,
		subset,
		factory.keyOf,
	)
	return contained && matched < len(membership)
}

// IsDisjoint returns true if left and right have no derived key in common.
func (factory HashSetByFactory[T, K]) IsDisjoint(
	left, right Source[T],
) bool {
	if left.Length() > right.Length() {
		left, right = right, left
	}
	keys := set.CollectKeyedKeys(left, factory.keyOf)
	for value := range right.All() {
		if _, exists := keys[factory.keyOf(value)]; exists {
			return false
		}
	}
	return true
}
