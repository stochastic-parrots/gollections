package set

import (
	"maps"

	"github.com/stochastic-parrots/gollections/internal/set"
)

// Equal reports whether first and every additional set contain the same
// values.
//
// Complexity: Expected O(Q*N + M), where Q is the number of additional sets,
// N is first.Length(), and M is the sum of their lengths.
func Equal[T comparable](first Readonly[T], others ...Readonly[T]) bool {
	for _, other := range others {
		if !equal(first, other) {
			return false
		}
	}
	return true
}

func equal[T comparable](left, right Readonly[T]) bool {
	if left.Length() != right.Length() {
		return false
	}
	if leftSet, ok := left.(*HashSet[T]); ok {
		return leftSet.Equal(right)
	}
	if rightSet, ok := right.(*HashSet[T]); ok {
		return rightSet.Equal(left)
	}
	membership := set.CollectHashMembership(left)
	matched, contained := set.MatchHashMembership(membership, right)
	return contained && matched == len(membership)
}

// IsSubset reports whether every value in subset is present in superset.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsSubset[T comparable](subset, superset Readonly[T]) bool {
	if subset.Length() > superset.Length() {
		return false
	}
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
	return set.MatchHashSubset(
		set.CollectHashMembership(subset),
		superset,
	)
}

// IsProperSubset reports whether subset is contained in superset and the two
// sources are not equal.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsProperSubset[T comparable](subset, superset Readonly[T]) bool {
	if subset.Length() >= superset.Length() {
		return false
	}
	if subsetSet, ok := subset.(*HashSet[T]); ok {
		return subsetSet.IsProperSubset(superset)
	}
	membership := set.CollectHashMembership(subset)
	matched, hasExtra := set.MatchHashContainedMembership(membership, superset)
	return matched == len(membership) && hasExtra
}

// IsSuperset reports whether every value in subset is present in superset.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsSuperset[T comparable](superset, subset Readonly[T]) bool {
	if subset.Length() > superset.Length() {
		return false
	}
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		return supersetSet.IsSuperset(subset)
	}
	return set.MatchHashSubset(
		set.CollectHashMembership(subset),
		superset,
	)
}

// IsProperSuperset reports whether superset contains subset and the two
// sources are not equal.
//
// Complexity: Expected O(N + M), where N = superset.Length() and M =
// subset.Length().
func IsProperSuperset[T comparable](superset, subset Readonly[T]) bool {
	if subset.Length() >= superset.Length() {
		return false
	}
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		return supersetSet.IsProperSuperset(subset)
	}
	membership := set.CollectHashMembership(subset)
	matched, hasExtra := set.MatchHashContainedMembership(membership, superset)
	return matched == len(membership) && hasExtra
}

// IsDisjoint reports whether first and every additional set are pairwise
// disjoint.
//
// Complexity: Expected O(M), where M is the sum of all operand lengths.
func IsDisjoint[T comparable](first Readonly[T], others ...Readonly[T]) bool {
	values := set.CollectHashKeys(first)
	for _, other := range others {
		if !addHashValues(values, other) {
			return false
		}
	}
	return true
}

func addHashValues[T comparable](known map[T]struct{}, source Readonly[T]) bool {
	current := set.CollectHashKeys(source)
	for value := range current {
		if _, exists := known[value]; exists {
			return false
		}
	}
	maps.Copy(known, current)
	return true
}

// EqualBy reports whether first and every additional set contain the same
// identities derived by keyOf. It panics if keyOf is nil.
//
// Complexity: Expected O(Q*N + M), where Q is the number of additional sets,
// N is first.Length(), and M is the sum of their lengths.
func EqualBy[T any, K comparable](
	keyOf func(T) K,
	first Readonly[T],
	others ...Readonly[T],
) bool {
	requireKeyOf(keyOf)
	for _, other := range others {
		if !equalBy(first, other, keyOf) {
			return false
		}
	}
	return true
}

func equalBy[T any, K comparable](left, right Readonly[T], keyOf func(T) K) bool {
	if left.Length() > right.Length() {
		left, right = right, left
	}
	membership := set.CollectKeyedMembership(left, keyOf)
	matched, contained := set.MatchKeyedMembership(membership, right, keyOf)
	return contained && matched == len(membership)
}

// IsSubsetBy reports whether every identity derived from subset is present in
// superset. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsSubsetBy[T any, K comparable](subset, superset Readonly[T], keyOf func(T) K) bool {
	requireKeyOf(keyOf)
	if subset.Length() <= superset.Length() {
		return set.MatchKeyedSubset(set.CollectKeyedMembership(subset, keyOf), superset, keyOf)
	}
	keys := set.CollectKeyedKeys(superset, keyOf)
	for value := range subset.All() {
		if _, exists := keys[keyOf(value)]; !exists {
			return false
		}
	}
	return true
}

// IsProperSubsetBy reports whether subset is contained in superset and their
// derived identities are not equal. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsProperSubsetBy[T any, K comparable](subset, superset Readonly[T], keyOf func(T) K) bool {
	requireKeyOf(keyOf)
	if subset.Length() <= superset.Length() {
		membership := set.CollectKeyedMembership(subset, keyOf)
		matched, hasExtra := set.MatchKeyedContainedMembership(membership, superset, keyOf)
		return matched == len(membership) && hasExtra
	}
	membership := set.CollectKeyedMembership(superset, keyOf)
	matched, contained := set.MatchKeyedMembership(membership, subset, keyOf)
	return contained && matched < len(membership)
}

// IsSupersetBy reports whether every identity derived from subset is present
// in superset. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = superset.Length() and M =
// subset.Length().
func IsSupersetBy[T any, K comparable](superset, subset Readonly[T], keyOf func(T) K) bool {
	return IsSubsetBy(subset, superset, keyOf)
}

// IsProperSupersetBy reports whether superset contains subset and their
// derived identities are not equal. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = superset.Length() and M =
// subset.Length().
func IsProperSupersetBy[T any, K comparable](superset, subset Readonly[T], keyOf func(T) K) bool {
	return IsProperSubsetBy(subset, superset, keyOf)
}

// IsDisjointBy reports whether first and every additional set are pairwise
// disjoint under identities derived by keyOf. It panics if keyOf is nil.
//
// Complexity: Expected O(M), where M is the sum of all operand lengths.
func IsDisjointBy[T any, K comparable](
	keyOf func(T) K,
	first Readonly[T],
	others ...Readonly[T],
) bool {
	requireKeyOf(keyOf)
	keys := set.CollectKeyedKeys(first, keyOf)
	for _, other := range others {
		if !addKeyedIdentities(keys, other, keyOf) {
			return false
		}
	}
	return true
}

func addKeyedIdentities[T any, K comparable](
	known map[K]struct{},
	source Readonly[T],
	keyOf func(T) K,
) bool {
	current := set.CollectKeyedKeys(source, keyOf)
	for key := range current {
		if _, exists := known[key]; exists {
			return false
		}
	}
	maps.Copy(known, current)
	return true
}

func requireKeyOf[T any, K comparable](keyOf func(T) K) {
	if keyOf == nil {
		panic("set: nil identity function")
	}
}
