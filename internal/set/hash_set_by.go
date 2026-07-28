package set

import (
	"iter"
	"maps"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// KeyedHashSet stores arbitrary values indexed by a comparable key derived
// from each value. The first value inserted for a key remains its
// representative.
type KeyedHashSet[T any, K comparable] struct {
	values map[K]T
	keyOf  func(T) K
}

// NewKeyedHashSet creates an empty keyed hash set with the requested initial
// capacity and identity function.
func NewKeyedHashSet[T any, K comparable](capacity int, keyOf func(T) K) *KeyedHashSet[T, K] {
	return &KeyedHashSet[T, K]{
		values: make(map[K]T, capacity),
		keyOf:  keyOf,
	}
}

// NewKeyedHashSetFromSlice creates a keyed hash set containing values from
// data. The first value for each derived key is preserved.
func NewKeyedHashSetFromSlice[T any, K comparable](data []T, keyOf func(T) K) *KeyedHashSet[T, K] {
	set := NewKeyedHashSet(len(data), keyOf)
	set.Adds(data...)
	return set
}

// NewKeyedHashSetFromSeq collects values yielded by seq into a keyed hash set.
// The first value for each derived key is preserved.
func NewKeyedHashSetFromSeq[T any, K comparable](seq iter.Seq[T], keyOf func(T) K) *KeyedHashSet[T, K] {
	set := NewKeyedHashSet(0, keyOf)
	for value := range seq {
		set.Add(value)
	}
	return set
}

// NewKeyedHashSetIntersection builds the intersection of sources while
// preserving representatives from the first source.
func NewKeyedHashSetIntersection[T any, K comparable](
	keyOf func(T) K,
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	if len(sources) == 0 {
		return NewKeyedHashSet(0, keyOf)
	}
	for _, source := range sources {
		if source.Length() == 0 {
			return NewKeyedHashSet(0, keyOf)
		}
	}
	if len(sources) == 1 {
		result := NewKeyedHashSet(sources[0].Length(), keyOf)
		result.unionWithSource(sources[0])
		return result
	}

	filterIdx := 1
	for idx := 2; idx < len(sources); idx++ {
		if sources[idx].Length() < sources[filterIdx].Length() {
			filterIdx = idx
		}
	}

	var result *KeyedHashSet[T, K]
	if sources[0].Length() <= sources[filterIdx].Length() {
		result = NewKeyedHashSet(sources[0].Length(), keyOf)
		result.unionWithSource(sources[0])
		result.intersectWith(sources[filterIdx])
	} else {
		keys := collectKeyedKeys(sources[filterIdx], keyOf)
		result = NewKeyedHashSet(min(sources[0].Length(), len(keys)), keyOf)
		for value := range sources[0].All() {
			key := keyOf(value)
			_, exists := keys[key]
			if !exists {
				continue
			}
			if _, duplicate := result.values[key]; !duplicate {
				result.values[key] = value
			}
		}
	}

	for idx, source := range sources[1:] {
		sourceIdx := idx + 1
		if sourceIdx != filterIdx {
			result.intersectWith(source)
		}
		if len(result.values) == 0 {
			break
		}
	}
	return result
}

// NewKeyedHashSetDifference builds the difference of sources directly in the
// returned set while preserving representatives from the first source.
func NewKeyedHashSetDifference[T any, K comparable](
	keyOf func(T) K,
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	if len(sources) == 0 {
		return NewKeyedHashSet(0, keyOf)
	}

	result := NewKeyedHashSet(sources[0].Length(), keyOf)
	result.unionWithSource(sources[0])
	for _, source := range sources[1:] {
		if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
			for _, value := range sourceSet.values {
				delete(result.values, keyOf(value))
			}
		} else {
			for value := range source.All() {
				delete(result.values, keyOf(value))
			}
		}
		if len(result.values) == 0 {
			break
		}
	}
	return result
}

// NewKeyedHashSetSymmetricDifference builds the symmetric difference of
// sources directly in the returned set. A surviving identity uses the
// representative from the last source that toggles it into the result.
func NewKeyedHashSetSymmetricDifference[T any, K comparable](
	keyOf func(T) K,
	sources ...Source[T],
) *KeyedHashSet[T, K] {
	if len(sources) == 0 {
		return NewKeyedHashSet(0, keyOf)
	}

	result := NewKeyedHashSet(sources[0].Length(), keyOf)
	result.unionWithSource(sources[0])
	toggleKeyedSources(result, sources[1:])
	return result
}

// Length returns the number of distinct keys in the set.
//
// Complexity: O(1).
func (set *KeyedHashSet[T, K]) Length() int {
	return len(set.values)
}

// IsEmpty returns true if the set contains no values.
//
// Complexity: O(1).
func (set *KeyedHashSet[T, K]) IsEmpty() bool {
	return len(set.values) == 0
}

// Contains returns true if the key derived from x belongs to the set.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Contains(x T) bool {
	_, ok := set.values[set.keyOf(x)]
	return ok
}

// Add inserts x unless its derived key is already present, reports whether the
// set changed, and preserves existing representatives.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Add(x T) bool {
	key := set.keyOf(x)
	if _, ok := set.values[key]; ok {
		return false
	}
	if set.values == nil {
		set.values = make(map[K]T)
	}
	set.values[key] = x
	return true
}

// Adds inserts every value in xs whose derived key is not already present and
// returns the number of values added.
//
// Complexity: Expected O(len(xs)) plus the identity-function calls.
func (set *KeyedHashSet[T, K]) Adds(xs ...T) (added int) {
	for _, x := range xs {
		if set.Add(x) {
			added++
		}
	}
	return added
}

// Remove deletes the value with the same derived key as x and reports whether
// it was present.
//
// Complexity: Expected O(1) plus the cost of the identity function.
func (set *KeyedHashSet[T, K]) Remove(x T) bool {
	key := set.keyOf(x)
	if _, ok := set.values[key]; !ok {
		return false
	}
	delete(set.values, key)
	return true
}

// Removes deletes every value in xs whose derived key is present and returns
// the number of values removed.
//
// Complexity: Expected O(len(xs)) plus the identity-function calls.
func (set *KeyedHashSet[T, K]) Removes(xs ...T) (removed int) {
	for _, x := range xs {
		if set.Remove(x) {
			removed++
		}
	}
	return removed
}

// Clone returns an independent keyed hash set containing the receiver's
// representative values and preserving its identity function.
//
// Complexity: O(N).
func (set *KeyedHashSet[T, K]) Clone() *KeyedHashSet[T, K] {
	clone := NewKeyedHashSet(len(set.values), set.keyOf)
	maps.Copy(clone.values, set.values)
	return clone
}

// Union returns a new set containing identities present in the receiver or any
// other operand. The receiver's identity function applies to every operand.
//
// Complexity: Expected O(N + sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) Union(
	others ...Source[T],
) *KeyedHashSet[T, K] {
	if len(others) == 1 {
		if otherSet, ok := others[0].(*KeyedHashSet[T, K]); ok && otherSet == set {
			return set.Clone()
		}
	}

	capacity := len(set.values)
	for _, other := range others {
		capacity += other.Length()
	}
	result := NewKeyedHashSet(capacity, set.keyOf)
	maps.Copy(result.values, set.values)
	result.UnionWith(others...)
	return result
}

// Intersection returns a new set containing identities present in the receiver
// and every other operand. Receiver representatives are preserved.
//
// Complexity: Expected O(N * len(others) + sum of operand lengths), plus keyOf
// calls.
func (set *KeyedHashSet[T, K]) Intersection(
	others ...Source[T],
) *KeyedHashSet[T, K] {
	if len(others) == 0 {
		return set.Clone()
	}
	if len(others) == 1 {
		if otherSet, ok := others[0].(*KeyedHashSet[T, K]); ok && otherSet == set {
			return set.Clone()
		}
	}
	if len(set.values) == 0 {
		return NewKeyedHashSet(0, set.keyOf)
	}
	for _, other := range others {
		if other.Length() == 0 {
			return NewKeyedHashSet(0, set.keyOf)
		}
	}

	filterIdx := smallestSource(others)
	keys := collectKeyedMatches(others[filterIdx], set.keyOf, set.values)
	matches := 0
	for key := range set.values {
		if _, exists := keys[key]; exists {
			matches++
		}
	}
	result := NewKeyedHashSet(matches, set.keyOf)
	for key, value := range set.values {
		_, exists := keys[key]
		if exists {
			result.values[key] = value
		}
	}

	for idx, other := range others {
		if idx != filterIdx {
			result.intersectWith(other)
		}
		if len(result.values) == 0 {
			break
		}
	}
	return result
}

// Difference returns a new set containing receiver identities absent from
// every other operand. Receiver representatives are preserved.
//
// Complexity: Expected O(N + sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) Difference(
	others ...Source[T],
) *KeyedHashSet[T, K] {
	if len(others) == 0 {
		return set.Clone()
	}

	excluded := make(map[K]struct{})
	for _, other := range others {
		if otherSet, ok := other.(*KeyedHashSet[T, K]); ok && otherSet == set {
			return NewKeyedHashSet(0, set.keyOf)
		}
		if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
			for _, value := range otherSet.values {
				key := set.keyOf(value)
				if _, exists := set.values[key]; exists {
					excluded[key] = struct{}{}
				}
			}
		} else {
			for value := range other.All() {
				key := set.keyOf(value)
				if _, exists := set.values[key]; exists {
					excluded[key] = struct{}{}
				}
			}
		}
		if len(excluded) == len(set.values) {
			return NewKeyedHashSet(0, set.keyOf)
		}
	}

	result := NewKeyedHashSet(len(set.values)-len(excluded), set.keyOf)
	for key, value := range set.values {
		if _, removed := excluded[key]; !removed {
			result.values[key] = value
		}
	}
	return result
}

// SymmetricDifference returns a new set containing identities present in an odd
// number of the receiver and other operands. A surviving identity uses the
// representative from the last operand that toggles it into the result.
//
// Complexity: Expected O(N + sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) SymmetricDifference(
	others ...Source[T],
) *KeyedHashSet[T, K] {
	if len(others) == 1 {
		if otherSet, ok := others[0].(*KeyedHashSet[T, K]); ok && otherSet == set {
			return NewKeyedHashSet(0, set.keyOf)
		}
	}
	result := set.Clone()
	toggleKeyedSources(result, others)
	return result
}

// Equal returns true if the receiver and other contain the same derived keys.
//
// Complexity: Expected O(N + other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) Equal(other Source[T]) bool {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok && otherSet == set {
		return true
	}
	if len(set.values) > other.Length() {
		return false
	}
	matched := make(map[K]struct{}, min(len(set.values), other.Length()))
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		for _, value := range otherSet.values {
			key := set.keyOf(value)
			if _, exists := set.values[key]; !exists {
				return false
			}
			matched[key] = struct{}{}
		}
		return len(matched) == len(set.values)
	}
	for value := range other.All() {
		key := set.keyOf(value)
		if _, exists := set.values[key]; !exists {
			return false
		}
		matched[key] = struct{}{}
	}
	return len(matched) == len(set.values)
}

// IsSubset returns true if every receiver key is present in other.
//
// Complexity: Expected O(N + other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) IsSubset(other Source[T]) bool {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok && otherSet == set {
		return true
	}
	if len(set.values) > other.Length() {
		return false
	}
	matched := make(map[K]struct{}, min(len(set.values), other.Length()))
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		for _, value := range otherSet.values {
			key := set.keyOf(value)
			if _, exists := set.values[key]; exists {
				matched[key] = struct{}{}
				if len(matched) == len(set.values) {
					return true
				}
			}
		}
		return len(matched) == len(set.values)
	}
	for value := range other.All() {
		key := set.keyOf(value)
		if _, exists := set.values[key]; exists {
			matched[key] = struct{}{}
			if len(matched) == len(set.values) {
				return true
			}
		}
	}
	return len(matched) == len(set.values)
}

// IsProperSubset returns true if the receiver is a subset of other and the two
// operands are not equal.
//
// Complexity: Expected O(N + other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) IsProperSubset(
	other Source[T],
) bool {
	if len(set.values) >= other.Length() {
		return false
	}
	matched := make(map[K]struct{}, min(len(set.values), other.Length()))
	hasExtra := false
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		for _, value := range otherSet.values {
			key := set.keyOf(value)
			if _, exists := set.values[key]; exists {
				matched[key] = struct{}{}
			} else {
				hasExtra = true
			}
			if hasExtra && len(matched) == len(set.values) {
				return true
			}
		}
		return hasExtra && len(matched) == len(set.values)
	}
	for value := range other.All() {
		key := set.keyOf(value)
		if _, exists := set.values[key]; exists {
			matched[key] = struct{}{}
		} else {
			hasExtra = true
		}
		if hasExtra && len(matched) == len(set.values) {
			return true
		}
	}
	return hasExtra && len(matched) == len(set.values)
}

// IsSuperset returns true if every identity in other is present in the
// receiver.
//
// Complexity: Expected O(other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) IsSuperset(other Source[T]) bool {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		if otherSet == set {
			return true
		}
		for _, value := range otherSet.values {
			if !set.Contains(value) {
				return false
			}
		}
		return true
	}
	for value := range other.All() {
		if !set.Contains(value) {
			return false
		}
	}
	return true
}

// IsProperSuperset returns true if the receiver is a superset of other and the
// two operands are not equal.
//
// Complexity: Expected O(N + other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) IsProperSuperset(
	other Source[T],
) bool {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok && otherSet == set {
		return false
	}
	matched := make(map[K]struct{}, min(len(set.values), other.Length()))
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		for _, value := range otherSet.values {
			key := set.keyOf(value)
			if _, exists := set.values[key]; !exists {
				return false
			}
			matched[key] = struct{}{}
			if len(matched) == len(set.values) {
				return false
			}
		}
		return len(matched) < len(set.values)
	}
	for value := range other.All() {
		key := set.keyOf(value)
		if _, exists := set.values[key]; !exists {
			return false
		}
		matched[key] = struct{}{}
		if len(matched) == len(set.values) {
			return false
		}
	}
	return len(matched) < len(set.values)
}

// IsDisjoint returns true if the receiver and other share no derived key.
//
// Complexity: Expected O(other.Length()), plus keyOf calls.
func (set *KeyedHashSet[T, K]) IsDisjoint(
	other Source[T],
) bool {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
		if otherSet == set {
			return len(set.values) == 0
		}
		for _, value := range otherSet.values {
			if set.Contains(value) {
				return false
			}
		}
		return true
	}
	for value := range other.All() {
		if set.Contains(value) {
			return false
		}
	}
	return true
}

// UnionWith adds identities from every operand and returns how many were added.
//
// Complexity: Expected O(sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) UnionWith(
	others ...Source[T],
) (added int) {
	for _, other := range others {
		if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
			if otherSet == set {
				continue
			}
			for _, value := range otherSet.values {
				if set.Add(value) {
					added++
				}
			}
		} else {
			for value := range other.All() {
				if set.Add(value) {
					added++
				}
			}
		}
	}
	return added
}

// IntersectWith removes receiver identities absent from any operand and
// returns how many were removed.
//
// Complexity: Expected O(N * len(others) + sum of operand lengths), plus keyOf
// calls.
func (set *KeyedHashSet[T, K]) IntersectWith(
	others ...Source[T],
) (removed int) {
	if len(set.values) == 0 || len(others) == 0 {
		return 0
	}

	first := smallestSource(others)
	if others[first].Length() == 0 {
		removed := len(set.values)
		clear(set.values)
		return removed
	}
	removed += set.intersectWith(others[first])
	for idx, other := range others {
		if idx != first {
			removed += set.intersectWith(other)
		}
		if len(set.values) == 0 {
			return removed
		}
	}
	return removed
}

// DifferenceWith removes identities present in any operand and returns how
// many were removed.
//
// Complexity: Expected O(sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) DifferenceWith(
	others ...Source[T],
) (removed int) {
	if len(set.values) == 0 {
		return 0
	}
	for _, other := range others {
		if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
			if otherSet == set {
				removed += len(set.values)
				clear(set.values)
				return removed
			}
			for _, value := range otherSet.values {
				key := set.keyOf(value)
				if _, exists := set.values[key]; exists {
					delete(set.values, key)
					removed++
				}
			}
			if len(set.values) == 0 {
				return removed
			}
			continue
		}
		for key := range collectKeyedMatches(other, set.keyOf, set.values) {
			if _, ok := set.values[key]; ok {
				delete(set.values, key)
				removed++
			}
		}
		if len(set.values) == 0 {
			return removed
		}
	}
	return removed
}

// SymmetricDifferenceWith toggles identities present in an odd number of
// operands and returns how many receiver memberships changed. An identity
// toggled into the receiver uses the representative from the last operand that
// makes it present.
//
// Complexity: Expected O(sum of operand lengths), plus keyOf calls.
func (set *KeyedHashSet[T, K]) SymmetricDifferenceWith(
	others ...Source[T],
) int {
	if len(others) == 0 {
		return 0
	}
	if len(others) == 1 {
		if otherSet, ok := others[0].(*KeyedHashSet[T, K]); ok {
			if otherSet == set {
				changed := len(set.values)
				clear(set.values)
				return changed
			}
			if set.values == nil {
				set.values = make(map[K]T, len(otherSet.values))
			}
			seen := make(map[K]struct{}, len(otherSet.values))
			for _, value := range otherSet.values {
				key := set.keyOf(value)
				if _, duplicate := seen[key]; duplicate {
					continue
				}
				seen[key] = struct{}{}
				if _, exists := set.values[key]; exists {
					delete(set.values, key)
				} else {
					set.values[key] = value
				}
			}
			return len(seen)
		}
	}

	toggles := collectKeyedIdentities(others[0], set.keyOf)
	var seen map[K]struct{}
	for _, other := range others[1:] {
		if seen == nil {
			seen = make(map[K]struct{}, other.Length())
		} else {
			clear(seen)
		}
		if otherSet, ok := other.(*KeyedHashSet[T, K]); ok {
			for _, value := range otherSet.values {
				key := set.keyOf(value)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
				if _, exists := toggles[key]; exists {
					delete(toggles, key)
				} else {
					toggles[key] = value
				}
			}
		} else {
			for value := range other.All() {
				key := set.keyOf(value)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
				if _, exists := toggles[key]; exists {
					delete(toggles, key)
				} else {
					toggles[key] = value
				}
			}
		}
	}

	if len(set.values) == 0 && len(toggles) > 0 {
		set.values = make(map[K]T, len(toggles))
	}
	for key, value := range toggles {
		if _, ok := set.values[key]; ok {
			delete(set.values, key)
		} else {
			set.values[key] = value
		}
	}
	return len(toggles)
}

func (set *KeyedHashSet[T, K]) intersectWith(
	other Source[T],
) (removed int) {
	if otherSet, ok := other.(*KeyedHashSet[T, K]); ok && otherSet == set {
		return 0
	}

	keys := collectKeyedMatches(other, set.keyOf, set.values)
	for key := range set.values {
		if _, exists := keys[key]; !exists {
			delete(set.values, key)
			removed++
		}
	}
	return removed
}

func (set *KeyedHashSet[T, K]) unionWithSource(source Source[T]) (added int) {
	for value := range source.All() {
		if set.Add(value) {
			added++
		}
	}
	return added
}

func toggleKeyedSources[T any, K comparable, S Source[T]](
	set *KeyedHashSet[T, K],
	others []S,
) {
	var seen map[K]struct{}
	for _, other := range others {
		if seen == nil {
			seen = make(map[K]struct{}, other.Length())
		} else {
			clear(seen)
		}
		if otherSet, ok := any(other).(*KeyedHashSet[T, K]); ok {
			for _, value := range otherSet.values {
				key := set.keyOf(value)
				if _, duplicate := seen[key]; duplicate {
					continue
				}
				seen[key] = struct{}{}
				if _, exists := set.values[key]; exists {
					delete(set.values, key)
				} else {
					set.values[key] = value
				}
			}
		} else {
			for value := range other.All() {
				key := set.keyOf(value)
				if _, duplicate := seen[key]; duplicate {
					continue
				}
				seen[key] = struct{}{}
				if _, exists := set.values[key]; exists {
					delete(set.values, key)
				} else {
					set.values[key] = value
				}
			}
		}
	}
}

// All returns a lazy iterator over the representative values in unspecified
// order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *KeyedHashSet[T, K]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range set.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate returns a lazy iterator over iteration positions and representative
// values in unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *KeyedHashSet[T, K]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for _, value := range set.values {
			if !yield(idx, value) {
				return
			}
			idx++
		}
	}
}

// Clear removes all values while preserving the identity function and keeping
// the set initialized for reuse.
//
// Complexity: O(allocated map storage).
func (set *KeyedHashSet[T, K]) Clear() {
	clear(set.values)
}

// MarshalJSON encodes the representative values as a JSON array in unspecified
// iteration order.
//
// Complexity: O(N).
func (set *KeyedHashSet[T, K]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(set)
}

func collectKeyedIdentities[T any, K comparable](
	source Source[T],
	keyOf func(T) K,
) map[K]T {
	identities := make(map[K]T, source.Length())
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			key := keyOf(value)
			if _, exists := identities[key]; !exists {
				identities[key] = value
			}
		}
	} else {
		for value := range source.All() {
			key := keyOf(value)
			if _, exists := identities[key]; !exists {
				identities[key] = value
			}
		}
	}
	return identities
}

func collectKeyedKeys[T any, K comparable](
	source Source[T],
	keyOf func(T) K,
) map[K]struct{} {
	keys := make(map[K]struct{}, source.Length())
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			keys[keyOf(value)] = struct{}{}
		}
	} else {
		for value := range source.All() {
			keys[keyOf(value)] = struct{}{}
		}
	}
	return keys
}

func collectKeyedMatches[T any, K comparable, V any](
	source Source[T],
	keyOf func(T) K,
	membership map[K]V,
) map[K]struct{} {
	matches := make(map[K]struct{}, min(source.Length(), len(membership)))
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			key := keyOf(value)
			if _, exists := membership[key]; exists {
				matches[key] = struct{}{}
			}
		}
	} else {
		for value := range source.All() {
			key := keyOf(value)
			if _, exists := membership[key]; exists {
				matches[key] = struct{}{}
			}
		}
	}
	return matches
}

// CollectKeyedKeys collects the distinct keys derived from values yielded by
// source.
func CollectKeyedKeys[T any, K comparable](
	source Source[T],
	keyOf func(T) K,
) map[K]struct{} {
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		keys := make(map[K]struct{}, len(sourceSet.values))
		for _, value := range sourceSet.values {
			keys[keyOf(value)] = struct{}{}
		}
		return keys
	}
	keys := make(map[K]struct{}, source.Length())
	for value := range source.All() {
		keys[keyOf(value)] = struct{}{}
	}
	return keys
}

// CollectKeyedMembership collects distinct derived keys as unseen membership
// entries.
func CollectKeyedMembership[T any, K comparable](
	source Source[T],
	keyOf func(T) K,
) map[K]bool {
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		membership := make(map[K]bool, len(sourceSet.values))
		for _, value := range sourceSet.values {
			membership[keyOf(value)] = false
		}
		return membership
	}
	membership := make(map[K]bool, source.Length())
	for value := range source.All() {
		membership[keyOf(value)] = false
	}
	return membership
}

// MatchKeyedMembership marks derived membership keys yielded by source and
// reports whether every yielded key belongs to membership.
func MatchKeyedMembership[T any, K comparable](
	membership map[K]bool,
	source Source[T],
	keyOf func(T) K,
) (matched int, contained bool) {
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			key := keyOf(value)
			seen, exists := membership[key]
			if !exists {
				return 0, false
			}
			if !seen {
				membership[key] = true
				matched++
			}
		}
		return matched, true
	}
	for value := range source.All() {
		key := keyOf(value)
		seen, exists := membership[key]
		if !exists {
			return 0, false
		}
		if !seen {
			membership[key] = true
			matched++
		}
	}
	return matched, true
}

// MatchKeyedSubset reports whether source yields every derived membership key
// at least once and marks keys as they are matched.
func MatchKeyedSubset[T any, K comparable](
	membership map[K]bool,
	source Source[T],
	keyOf func(T) K,
) bool {
	if len(membership) == 0 {
		return true
	}
	matched := 0
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			key := keyOf(value)
			seen, exists := membership[key]
			if exists && !seen {
				membership[key] = true
				matched++
				if matched == len(membership) {
					return true
				}
			}
		}
		return false
	}
	for value := range source.All() {
		key := keyOf(value)
		seen, exists := membership[key]
		if exists && !seen {
			membership[key] = true
			matched++
			if matched == len(membership) {
				return true
			}
		}
	}
	return false
}

// MatchKeyedContainedMembership marks derived membership keys yielded by source
// and reports whether source also yields a key outside membership.
func MatchKeyedContainedMembership[T any, K comparable](
	membership map[K]bool,
	source Source[T],
	keyOf func(T) K,
) (matched int, hasExtra bool) {
	if sourceSet, ok := source.(*KeyedHashSet[T, K]); ok {
		for _, value := range sourceSet.values {
			key := keyOf(value)
			seen, exists := membership[key]
			if !exists {
				hasExtra = true
			} else if !seen {
				membership[key] = true
				matched++
			}
			if matched == len(membership) && hasExtra {
				return matched, true
			}
		}
		return matched, hasExtra
	}
	for value := range source.All() {
		key := keyOf(value)
		seen, exists := membership[key]
		if !exists {
			hasExtra = true
		} else if !seen {
			membership[key] = true
			matched++
		}
		if matched == len(membership) && hasExtra {
			return matched, true
		}
	}
	return matched, hasExtra
}
