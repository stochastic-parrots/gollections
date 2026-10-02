package set

import "github.com/stochastic-parrots/gollections/internal/set"

// Equal reports whether left and right contain the same values.
//
// Complexity: Expected O(N + M), where N = left.Length() and M =
// right.Length().
func Equal[T comparable](left, right Readonly[T]) bool {
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

// IsSubset reports whether every value in subset is present in superset.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsSubset[T comparable](subset, superset Readonly[T]) bool {
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
		return set.MatchHashSubset(
			set.CollectHashMembership(subset),
			superset,
		)
	}
	keys := set.CollectHashKeys(superset)
	for value := range subset.All() {
		if _, exists := keys[value]; !exists {
			return false
		}
	}
	return true
}

// IsProperSubset reports whether subset is contained in superset and the two
// sources are not equal.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsProperSubset[T comparable](subset, superset Readonly[T]) bool {
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

// IsSuperset reports whether every value in subset is present in superset.
//
// Complexity: Expected O(N + M), where N = subset.Length() and M =
// superset.Length().
func IsSuperset[T comparable](superset, subset Readonly[T]) bool {
	if supersetSet, ok := superset.(*HashSet[T]); ok {
		return supersetSet.IsSuperset(subset)
	}
	if subset.Length() <= superset.Length() {
		return set.MatchHashSubset(
			set.CollectHashMembership(subset),
			superset,
		)
	}
	keys := set.CollectHashKeys(superset)
	for value := range subset.All() {
		if _, exists := keys[value]; !exists {
			return false
		}
	}
	return true
}

// IsProperSuperset reports whether superset contains subset and the two
// sources are not equal.
//
// Complexity: Expected O(N + M), where N = superset.Length() and M =
// subset.Length().
func IsProperSuperset[T comparable](superset, subset Readonly[T]) bool {
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

// IsDisjoint reports whether left and right share no values.
//
// Complexity: Expected O(N + M), where N = left.Length() and M =
// right.Length().
func IsDisjoint[T comparable](left, right Readonly[T]) bool {
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

// EqualBy reports whether left and right contain the same identities derived
// by keyOf. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = left.Length() and M =
// right.Length().
func EqualBy[T any, K comparable](left, right Readonly[T], keyOf func(T) K) bool {
	requireKeyOf(keyOf)
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

// IsDisjointBy reports whether left and right share no identity derived by
// keyOf. It panics if keyOf is nil.
//
// Complexity: Expected O(N + M), where N = left.Length() and M =
// right.Length().
func IsDisjointBy[T any, K comparable](left, right Readonly[T], keyOf func(T) K) bool {
	requireKeyOf(keyOf)
	if left.Length() > right.Length() {
		left, right = right, left
	}
	keys := set.CollectKeyedKeys(left, keyOf)
	for value := range right.All() {
		if _, exists := keys[keyOf(value)]; exists {
			return false
		}
	}
	return true
}

func requireKeyOf[T any, K comparable](keyOf func(T) K) {
	if keyOf == nil {
		panic("set: nil identity function")
	}
}
