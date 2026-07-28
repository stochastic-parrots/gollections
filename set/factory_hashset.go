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
// HashSet's zero value is ready for use. A HashSet value must not be copied
// after first use; share its pointer instead.
type HashSet[T comparable] = set.HashSet[T]

var (
	_ Set[int]                    = &set.HashSet[int]{}
	_ Algebra[int, *HashSet[int]] = &set.HashSet[int]{}
	_ InPlaceAlgebra[int]         = &set.HashSet[int]{}
)

// HashFactory constructs hash sets using Go equality for value identity.
//
// The zero value is ready for use.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation                    Time Complexity
//	-------------------------    ---------------
//	New(capacity)                O(capacity)
//	From/FromSeq                 O(N)
//	Add/Contains/Remove          O(1)
//	Adds/Removes(xs...T)         O(len(xs))
//	Union/Difference/SymDiff     O(sum of operand lengths)
//	Intersection                 O(N * sources + sum of source lengths)
//	Set relations                O(N + M)
//	All/Enumerate                O(N)
//	Clear                        O(allocated map storage)
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

// Union returns a new set containing identities present in any source. With no
// sources, it returns an empty set.
func (factory HashFactory[T]) Union(
	sources ...Source[T],
) *HashSet[T] {
	if len(sources) == 0 {
		return factory.New(0)
	}
	capacity := 0
	for _, source := range sources {
		capacity += source.Length()
	}
	result := factory.New(capacity)
	result.UnionWith(sources...)
	return result
}

// Intersection returns a new set containing identities present in every
// source. Representatives come from the first source. With no sources, it
// returns an empty set.
func (factory HashFactory[T]) Intersection(
	sources ...Source[T],
) *HashSet[T] {
	return set.NewHashSetIntersection(sources...)
}

// Difference returns a new set containing identities from the first source
// that are absent from every subsequent source. With no sources, it returns an
// empty set.
func (factory HashFactory[T]) Difference(
	sources ...Source[T],
) *HashSet[T] {
	return set.NewHashSetDifference(sources...)
}

// SymmetricDifference returns a new set containing identities present in an
// odd number of sources. With no sources, it returns an empty set.
func (factory HashFactory[T]) SymmetricDifference(
	sources ...Source[T],
) *HashSet[T] {
	return set.NewHashSetSymmetricDifference(sources...)
}

// Equal returns true if left and right contain the same identities.
func (factory HashFactory[T]) Equal(
	left, right Source[T],
) bool {
	if leftSet, ok := left.(*HashSet[T]); ok {
		return leftSet.Equal(right)
	}
	if rightSet, ok := right.(*HashSet[T]); ok {
		return rightSet.Equal(left)
	}
	if left.Length() > right.Length() {
		left, right = right, left
	}
	membership := set.CollectHashMembership(left)
	matched, contained := set.MatchHashMembership(membership, right)
	return contained && matched == len(membership)
}

// IsSubset returns true if every identity in subset is present in superset.
func (factory HashFactory[T]) IsSubset(
	subset, superset Source[T],
) bool {
	if subsetSet, ok := subset.(*HashSet[T]); ok {
		return subsetSet.IsSubset(superset)
	}
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		for value := range subset.All() {
			if !supersetSet.Contains(value) {
				return false
			}
		}
		return true
	}
	if subset.Length() <= superset.Length() {
		return set.MatchHashSubset(set.CollectHashMembership(subset), superset)
	}
	keys := set.CollectHashKeys(superset)
	for value := range subset.All() {
		if _, exists := keys[value]; !exists {
			return false
		}
	}
	return true
}

// IsProperSubset returns true if subset is contained in superset and the two
// sources are not equal.
func (factory HashFactory[T]) IsProperSubset(
	subset, superset Source[T],
) bool {
	if subsetSet, ok := subset.(*HashSet[T]); ok {
		return subsetSet.IsProperSubset(superset)
	}
	if subset.Length() <= superset.Length() {
		membership := set.CollectHashMembership(subset)
		matched, hasExtra := set.MatchHashContainedMembership(membership, superset)
		return matched == len(membership) && hasExtra
	}
	membership := set.CollectHashMembership(superset)
	matched, contained := set.MatchHashMembership(membership, subset)
	return contained && matched < len(membership)
}

// IsSuperset returns true if every identity in subset is present in superset.
func (factory HashFactory[T]) IsSuperset(
	superset, subset Source[T],
) bool {
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		return supersetSet.IsSuperset(subset)
	}
	if subset.Length() <= superset.Length() {
		return set.MatchHashSubset(set.CollectHashMembership(subset), superset)
	}
	keys := set.CollectHashKeys(superset)
	for value := range subset.All() {
		if _, exists := keys[value]; !exists {
			return false
		}
	}
	return true
}

// IsProperSuperset returns true if superset contains subset and the two sources
// are not equal.
func (factory HashFactory[T]) IsProperSuperset(
	superset, subset Source[T],
) bool {
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		return supersetSet.IsProperSuperset(subset)
	}
	if subset.Length() <= superset.Length() {
		membership := set.CollectHashMembership(subset)
		matched, hasExtra := set.MatchHashContainedMembership(membership, superset)
		return matched == len(membership) && hasExtra
	}
	membership := set.CollectHashMembership(superset)
	matched, contained := set.MatchHashMembership(membership, subset)
	return contained && matched < len(membership)
}

// IsDisjoint returns true if left and right share no identity.
func (factory HashFactory[T]) IsDisjoint(
	left, right Source[T],
) bool {
	if leftSet, ok := left.(*HashSet[T]); ok {
		return leftSet.IsDisjoint(right)
	}
	if rightSet, ok := right.(*HashSet[T]); ok {
		return rightSet.IsDisjoint(left)
	}
	if left.Length() > right.Length() {
		left, right = right, left
	}
	keys := set.CollectHashKeys(left)
	for value := range right.All() {
		if _, exists := keys[value]; exists {
			return false
		}
	}
	return true
}
