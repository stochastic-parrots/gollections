package set

import (
	"iter"
	"maps"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// HashSet stores unique comparable values as keys in a Go map.
type HashSet[T comparable] struct {
	values map[T]struct{}
}

// NewHashSet creates an empty hash set with the requested initial capacity.
func NewHashSet[T comparable](capacity int) *HashSet[T] {
	return &HashSet[T]{values: make(map[T]struct{}, capacity)}
}

// NewHashSetFromSlice creates a hash set containing the unique values in data.
func NewHashSetFromSlice[T comparable](data []T) *HashSet[T] {
	set := NewHashSet[T](len(data))
	set.Adds(data...)
	return set
}

// NewHashSetFromSeq collects the unique values yielded by seq into a hash set.
func NewHashSetFromSeq[T comparable](seq iter.Seq[T]) *HashSet[T] {
	set := NewHashSet[T](0)
	for value := range seq {
		set.Add(value)
	}
	return set
}

// NewHashSetIntersection builds the intersection of sources while starting
// with the shortest source.
func NewHashSetIntersection[T comparable](
	sources ...Source[T],
) *HashSet[T] {
	if len(sources) == 0 {
		return NewHashSet[T](0)
	}
	for _, source := range sources {
		if source.Length() == 0 {
			return NewHashSet[T](0)
		}
	}

	first := smallestSource(sources)
	if len(sources) > 1 {
		allHashSets := true
		for _, source := range sources {
			if _, ok := source.(*HashSet[T]); !ok {
				allHashSets = false
				break
			}
		}
		if allHashSets {
			candidate := sources[first].(*HashSet[T])
			matches := 0
			for value := range candidate.values {
				if hashSourcesContain(sources, value) {
					matches++
				}
			}
			result := NewHashSet[T](matches)
			for value := range candidate.values {
				if hashSourcesContain(sources, value) {
					result.values[value] = struct{}{}
				}
			}
			return result
		}
	}

	var result *HashSet[T]
	if sourceSet, ok := sources[first].(*HashSet[T]); ok {
		result = sourceSet.Clone()
	} else {
		result = NewHashSet[T](sources[first].Length())
		result.unionWithSource(sources[first])
	}
	for idx, source := range sources {
		if idx != first {
			result.intersectWith(source)
		}
		if len(result.values) == 0 {
			break
		}
	}
	return result
}

// NewHashSetDifference builds the difference of sources directly in the
// returned set while preserving values from the first source.
func NewHashSetDifference[T comparable](
	sources ...Source[T],
) *HashSet[T] {
	if len(sources) == 0 {
		return NewHashSet[T](0)
	}

	var result *HashSet[T]
	if sourceSet, ok := sources[0].(*HashSet[T]); ok {
		result = sourceSet.Clone()
	} else {
		result = NewHashSet[T](sources[0].Length())
		result.unionWithSource(sources[0])
	}
	for _, source := range sources[1:] {
		if sourceSet, ok := source.(*HashSet[T]); ok {
			if len(sourceSet.values) < len(result.values) {
				for value := range sourceSet.values {
					delete(result.values, value)
				}
			} else {
				for value := range result.values {
					if _, exists := sourceSet.values[value]; exists {
						delete(result.values, value)
					}
				}
			}
		} else {
			for value := range source.All() {
				delete(result.values, value)
			}
		}
		if len(result.values) == 0 {
			break
		}
	}
	return result
}

// NewHashSetSymmetricDifference builds the symmetric difference of sources
// directly in the returned set.
func NewHashSetSymmetricDifference[T comparable](
	sources ...Source[T],
) *HashSet[T] {
	if len(sources) == 0 {
		return NewHashSet[T](0)
	}

	var result *HashSet[T]
	if sourceSet, ok := sources[0].(*HashSet[T]); ok {
		result = sourceSet.Clone()
	} else {
		result = NewHashSet[T](sources[0].Length())
		result.unionWithSource(sources[0])
	}
	toggleHashSources(result, sources[1:])
	return result
}

// Length returns the number of values in the set.
//
// Complexity: O(1).
func (set *HashSet[T]) Length() int {
	return len(set.values)
}

// IsEmpty returns true if the set contains no values.
//
// Complexity: O(1).
func (set *HashSet[T]) IsEmpty() bool {
	return len(set.values) == 0
}

// Contains returns true if x belongs to the set.
//
// Complexity: Expected O(1).
func (set *HashSet[T]) Contains(x T) bool {
	_, ok := set.values[x]
	return ok
}

// Add inserts x if it is not already present and reports whether the set
// changed.
//
// Complexity: Expected O(1).
func (set *HashSet[T]) Add(x T) bool {
	if _, ok := set.values[x]; ok {
		return false
	}
	if set.values == nil {
		set.values = make(map[T]struct{})
	}
	set.values[x] = struct{}{}
	return true
}

// Adds inserts every value in xs that is not already present and returns the
// number of values added.
//
// Complexity: Expected O(len(xs)).
func (set *HashSet[T]) Adds(xs ...T) (added int) {
	for _, x := range xs {
		if set.Add(x) {
			added++
		}
	}
	return added
}

// Remove deletes x and reports whether it was present.
//
// Complexity: Expected O(1).
func (set *HashSet[T]) Remove(x T) bool {
	if _, ok := set.values[x]; !ok {
		return false
	}
	delete(set.values, x)
	return true
}

// Removes deletes every value in xs that is present and returns the number of
// values removed.
//
// Complexity: Expected O(len(xs)).
func (set *HashSet[T]) Removes(xs ...T) (removed int) {
	for _, x := range xs {
		if set.Remove(x) {
			removed++
		}
	}
	return removed
}

// Clone returns an independent hash set containing the receiver's values.
//
// Complexity: O(N).
func (set *HashSet[T]) Clone() *HashSet[T] {
	clone := NewHashSet[T](len(set.values))
	maps.Copy(clone.values, set.values)
	return clone
}

// Union returns a new set containing values present in the receiver or any
// other operand. The receiver is not modified.
//
// Complexity: Expected O(N + sum of operand lengths).
func (set *HashSet[T]) Union(others ...Source[T]) *HashSet[T] {
	if len(others) == 1 {
		if otherSet, ok := others[0].(*HashSet[T]); ok && otherSet == set {
			return set.Clone()
		}
	}

	capacity := len(set.values)
	for _, other := range others {
		capacity += other.Length()
	}
	result := NewHashSet[T](capacity)
	maps.Copy(result.values, set.values)
	result.UnionWith(others...)
	return result
}

// Intersection returns a new set containing values present in the receiver and
// every other operand. The receiver is not modified.
//
// Complexity: Expected O(N * len(others) + sum of operand lengths).
func (set *HashSet[T]) Intersection(others ...Source[T]) *HashSet[T] {
	if len(others) == 0 {
		return set.Clone()
	}
	if len(others) == 1 {
		if otherSet, ok := others[0].(*HashSet[T]); ok && otherSet == set {
			return set.Clone()
		}
	}
	if len(set.values) == 0 {
		return NewHashSet[T](0)
	}
	for _, other := range others {
		if other.Length() == 0 {
			return NewHashSet[T](0)
		}
	}

	if len(others) > 1 {
		allHashSets := true
		candidate := set
		for _, other := range others {
			otherSet, ok := other.(*HashSet[T])
			if !ok {
				allHashSets = false
				break
			}
			if len(otherSet.values) < len(candidate.values) {
				candidate = otherSet
			}
		}
		if allHashSets {
			matches := 0
			for value := range candidate.values {
				if _, exists := set.values[value]; exists &&
					hashSourcesContain(others, value) {
					matches++
				}
			}
			result := NewHashSet[T](matches)
			for value := range candidate.values {
				if _, exists := set.values[value]; exists &&
					hashSourcesContain(others, value) {
					result.values[value] = struct{}{}
				}
			}
			return result
		}
	}

	filterIdx := smallestSource(others)
	var filter map[T]struct{}
	if filterSet, ok := others[filterIdx].(*HashSet[T]); ok {
		filter = filterSet.values
	} else {
		filter = collectHashMatches(others[filterIdx], set.values)
	}
	candidates, membership := set.values, filter
	if len(candidates) > len(membership) {
		candidates, membership = membership, candidates
	}
	matches := 0
	for value := range candidates {
		if _, exists := membership[value]; exists {
			matches++
		}
	}
	result := NewHashSet[T](matches)
	for value := range candidates {
		if _, exists := membership[value]; exists {
			result.values[value] = struct{}{}
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

// Difference returns a new set containing receiver values absent from every
// other operand. The receiver is not modified.
//
// Complexity: Expected O(N * len(others) + sum of operand lengths).
func (set *HashSet[T]) Difference(others ...Source[T]) *HashSet[T] {
	if len(others) == 0 {
		return set.Clone()
	}
	for _, other := range others {
		if otherSet, ok := other.(*HashSet[T]); ok && otherSet == set {
			return NewHashSet[T](0)
		}
	}

	allHashSets := true
	for _, other := range others {
		if _, ok := other.(*HashSet[T]); !ok {
			allHashSets = false
			break
		}
	}
	if allHashSets {
		survivors := 0
		for value := range set.values {
			excluded := false
			for _, other := range others {
				otherSet := other.(*HashSet[T])
				if _, exists := otherSet.values[value]; exists {
					excluded = true
					break
				}
			}
			if !excluded {
				survivors++
			}
		}

		result := NewHashSet[T](survivors)
		for value := range set.values {
			excluded := false
			for _, other := range others {
				otherSet := other.(*HashSet[T])
				if _, exists := otherSet.values[value]; exists {
					excluded = true
					break
				}
			}
			if !excluded {
				result.values[value] = struct{}{}
			}
		}
		return result
	}

	excluded := make(map[T]struct{})
	for _, other := range others {
		if otherSet, ok := other.(*HashSet[T]); ok {
			if len(otherSet.values) < len(set.values) {
				for value := range otherSet.values {
					if _, exists := set.values[value]; exists {
						excluded[value] = struct{}{}
					}
				}
			} else {
				for value := range set.values {
					if _, exists := otherSet.values[value]; exists {
						excluded[value] = struct{}{}
					}
				}
			}
		} else {
			for value := range other.All() {
				if _, exists := set.values[value]; exists {
					excluded[value] = struct{}{}
				}
			}
		}
		if len(excluded) == len(set.values) {
			return NewHashSet[T](0)
		}
	}

	result := NewHashSet[T](len(set.values) - len(excluded))
	for value := range set.values {
		if _, removed := excluded[value]; !removed {
			result.values[value] = struct{}{}
		}
	}
	return result
}

// SymmetricDifference returns a new set containing values present in an odd
// number of the receiver and other operands. The receiver is not modified.
//
// Complexity: Expected O(N + sum of operand lengths).
func (set *HashSet[T]) SymmetricDifference(
	others ...Source[T],
) *HashSet[T] {
	if len(others) == 1 {
		if otherSet, ok := others[0].(*HashSet[T]); ok && otherSet == set {
			return NewHashSet[T](0)
		}
	}
	result := set.Clone()
	toggleHashSources(result, others)
	return result
}

// Equal returns true if the receiver and other contain the same values.
//
// Complexity: Expected O(N + other.Length()).
func (set *HashSet[T]) Equal(other Source[T]) bool {
	if len(set.values) > other.Length() {
		return false
	}
	if otherSet, ok := other.(*HashSet[T]); ok {
		if otherSet == set {
			return true
		}
		return len(set.values) == len(otherSet.values) &&
			containsAll(set.values, otherSet.values)
	}
	matched := make(map[T]struct{}, min(len(set.values), other.Length()))
	for value := range other.All() {
		if _, exists := set.values[value]; !exists {
			return false
		}
		matched[value] = struct{}{}
	}
	return len(matched) == len(set.values)
}

// IsSubset returns true if every receiver value is present in other.
//
// Complexity: Expected O(N + other.Length()).
func (set *HashSet[T]) IsSubset(other Source[T]) bool {
	if len(set.values) > other.Length() {
		return false
	}
	if otherSet, ok := other.(*HashSet[T]); ok {
		if otherSet == set {
			return true
		}
		return len(set.values) <= len(otherSet.values) &&
			containsAll(set.values, otherSet.values)
	}
	matched := make(map[T]struct{}, min(len(set.values), other.Length()))
	for value := range other.All() {
		if _, exists := set.values[value]; exists {
			matched[value] = struct{}{}
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
// Complexity: Expected O(N + other.Length()).
func (set *HashSet[T]) IsProperSubset(other Source[T]) bool {
	if len(set.values) >= other.Length() {
		return false
	}
	if otherSet, ok := other.(*HashSet[T]); ok {
		return len(set.values) < len(otherSet.values) &&
			containsAll(set.values, otherSet.values)
	}
	matched := make(map[T]struct{}, min(len(set.values), other.Length()))
	hasExtra := false
	for value := range other.All() {
		if _, exists := set.values[value]; exists {
			matched[value] = struct{}{}
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
// Complexity: Expected O(other.Length()).
func (set *HashSet[T]) IsSuperset(other Source[T]) bool {
	if otherSet, ok := other.(*HashSet[T]); ok {
		if otherSet == set {
			return true
		}
		return len(set.values) >= len(otherSet.values) &&
			containsAll(otherSet.values, set.values)
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
// Complexity: Expected O(N + other.Length()).
func (set *HashSet[T]) IsProperSuperset(other Source[T]) bool {
	if otherSet, ok := other.(*HashSet[T]); ok {
		return len(set.values) > len(otherSet.values) &&
			containsAll(otherSet.values, set.values)
	}
	matched := make(map[T]struct{}, min(len(set.values), other.Length()))
	for value := range other.All() {
		if _, exists := set.values[value]; !exists {
			return false
		}
		matched[value] = struct{}{}
		if len(matched) == len(set.values) {
			return false
		}
	}
	return len(matched) < len(set.values)
}

// IsDisjoint returns true if the receiver and other share no value.
//
// Complexity: Expected O(other.Length()).
func (set *HashSet[T]) IsDisjoint(other Source[T]) bool {
	if otherSet, ok := other.(*HashSet[T]); ok {
		if otherSet == set {
			return len(set.values) == 0
		}
		left, right := set.values, otherSet.values
		if len(left) > len(right) {
			left, right = right, left
		}
		for value := range left {
			if _, exists := right[value]; exists {
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
// Complexity: Expected O(sum of operand lengths).
func (set *HashSet[T]) UnionWith(others ...Source[T]) (added int) {
	for _, other := range others {
		if otherSet, ok := other.(*HashSet[T]); ok {
			if otherSet == set {
				continue
			}
			if len(set.values) == 0 {
				if set.values == nil {
					set.values = make(map[T]struct{}, len(otherSet.values))
				}
				maps.Copy(set.values, otherSet.values)
				added += len(otherSet.values)
				continue
			}
			for value := range otherSet.values {
				if set.Add(value) {
					added++
				}
			}
			continue
		}
		for value := range other.All() {
			if set.Add(value) {
				added++
			}
		}
	}
	return added
}

// IntersectWith removes receiver identities absent from any operand and
// returns how many were removed.
//
// Complexity: Expected O(N * len(others) + sum of operand lengths).
func (set *HashSet[T]) IntersectWith(
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
// Complexity: Expected O(sum of operand lengths).
func (set *HashSet[T]) DifferenceWith(
	others ...Source[T],
) (removed int) {
	if len(set.values) == 0 {
		return 0
	}
	for _, other := range others {
		if otherSet, ok := other.(*HashSet[T]); ok {
			if otherSet == set {
				removed += len(set.values)
				clear(set.values)
				return removed
			}
			for value := range otherSet.values {
				if _, exists := set.values[value]; exists {
					delete(set.values, value)
					removed++
				}
			}
			if len(set.values) == 0 {
				return removed
			}
			continue
		}
		for value := range collectHashMatches(other, set.values) {
			if _, ok := set.values[value]; ok {
				delete(set.values, value)
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
// operands and returns how many receiver memberships changed.
//
// Complexity: Expected O(sum of operand lengths).
func (set *HashSet[T]) SymmetricDifferenceWith(
	others ...Source[T],
) int {
	if len(others) == 0 {
		return 0
	}
	if len(others) == 1 {
		if otherSet, ok := others[0].(*HashSet[T]); ok {
			if otherSet == set {
				changed := len(set.values)
				clear(set.values)
				return changed
			}
			for value := range otherSet.values {
				if _, exists := set.values[value]; exists {
					delete(set.values, value)
				} else {
					set.values[value] = struct{}{}
				}
			}
			return len(otherSet.values)
		}
	}

	toggles := collectHashIdentities(others[0])
	var seen map[T]struct{}
	for _, other := range others[1:] {
		if otherSet, ok := other.(*HashSet[T]); ok {
			for value := range otherSet.values {
				if _, exists := toggles[value]; exists {
					delete(toggles, value)
				} else {
					toggles[value] = struct{}{}
				}
			}
			continue
		}
		if seen == nil {
			seen = make(map[T]struct{}, other.Length())
		} else {
			clear(seen)
		}
		for value := range other.All() {
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			if _, ok := toggles[value]; ok {
				delete(toggles, value)
			} else {
				toggles[value] = struct{}{}
			}
		}
	}

	if len(set.values) == 0 && len(toggles) > 0 {
		set.values = make(map[T]struct{}, len(toggles))
	}
	for value := range toggles {
		if _, ok := set.values[value]; ok {
			delete(set.values, value)
		} else {
			set.values[value] = struct{}{}
		}
	}
	return len(toggles)
}

func (set *HashSet[T]) intersectWith(other Source[T]) (removed int) {
	if otherSet, ok := other.(*HashSet[T]); ok {
		if otherSet == set {
			return 0
		}
		for value := range set.values {
			if _, exists := otherSet.values[value]; !exists {
				delete(set.values, value)
				removed++
			}
		}
		return removed
	}

	identities := collectHashMatches(other, set.values)
	for value := range set.values {
		if _, exists := identities[value]; !exists {
			delete(set.values, value)
			removed++
		}
	}
	return removed
}

func (set *HashSet[T]) unionWithSource(source Source[T]) (added int) {
	for value := range source.All() {
		if set.Add(value) {
			added++
		}
	}
	return added
}

func toggleHashSources[T comparable, S Source[T]](
	set *HashSet[T],
	others []S,
) {
	var seen map[T]struct{}
	for _, other := range others {
		if otherSet, ok := any(other).(*HashSet[T]); ok {
			for value := range otherSet.values {
				if _, exists := set.values[value]; exists {
					delete(set.values, value)
				} else {
					set.values[value] = struct{}{}
				}
			}
			continue
		}

		if seen == nil {
			seen = make(map[T]struct{}, other.Length())
		} else {
			clear(seen)
		}
		for value := range other.All() {
			if _, duplicate := seen[value]; duplicate {
				continue
			}
			seen[value] = struct{}{}
			if _, exists := set.values[value]; exists {
				delete(set.values, value)
			} else {
				set.values[value] = struct{}{}
			}
		}
	}
}

// All returns a lazy iterator over the set in unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *HashSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range set.values {
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate returns a lazy iterator over iteration positions and values in
// unspecified order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (set *HashSet[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for value := range set.values {
			if !yield(idx, value) {
				return
			}
			idx++
		}
	}
}

// Clear removes all values while keeping the set initialized for reuse.
//
// Complexity: O(allocated map storage).
func (set *HashSet[T]) Clear() {
	clear(set.values)
}

// MarshalJSON encodes the set as a JSON array in unspecified iteration order.
//
// Complexity: O(N).
func (set *HashSet[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(set)
}

func collectHashIdentities[T comparable](
	source Source[T],
) map[T]struct{} {
	if set, ok := source.(*HashSet[T]); ok {
		identities := make(map[T]struct{}, len(set.values))
		maps.Copy(identities, set.values)
		return identities
	}
	identities := make(map[T]struct{}, source.Length())
	for value := range source.All() {
		if _, ok := identities[value]; !ok {
			identities[value] = struct{}{}
		}
	}
	return identities
}

func collectHashMatches[T comparable](
	source Source[T],
	membership map[T]struct{},
) map[T]struct{} {
	matches := make(map[T]struct{}, min(source.Length(), len(membership)))
	for value := range source.All() {
		if _, exists := membership[value]; exists {
			matches[value] = struct{}{}
		}
	}
	return matches
}

func smallestSource[T any, S Source[T]](sources []S) int {
	smallest := 0
	for idx := 1; idx < len(sources); idx++ {
		if sources[idx].Length() < sources[smallest].Length() {
			smallest = idx
		}
	}
	return smallest
}

func hashSourcesContain[T comparable, S Source[T]](
	sources []S,
	value T,
) bool {
	for _, source := range sources {
		sourceSet := any(source).(*HashSet[T])
		if _, exists := sourceSet.values[value]; !exists {
			return false
		}
	}
	return true
}

func containsAll[T comparable](subset, superset map[T]struct{}) bool {
	for value := range subset {
		if _, ok := superset[value]; !ok {
			return false
		}
	}
	return true
}

// CollectHashKeys collects the unique values yielded by source as map keys.
func CollectHashKeys[T comparable](source Source[T]) map[T]struct{} {
	if sourceSet, ok := source.(*HashSet[T]); ok {
		keys := make(map[T]struct{}, len(sourceSet.values))
		maps.Copy(keys, sourceSet.values)
		return keys
	}
	keys := make(map[T]struct{}, source.Length())
	for value := range source.All() {
		keys[value] = struct{}{}
	}
	return keys
}

// CollectHashMembership collects unique values yielded by source as unseen
// membership entries.
func CollectHashMembership[T comparable](source Source[T]) map[T]bool {
	if sourceSet, ok := source.(*HashSet[T]); ok {
		membership := make(map[T]bool, len(sourceSet.values))
		for value := range sourceSet.values {
			membership[value] = false
		}
		return membership
	}
	membership := make(map[T]bool, source.Length())
	for value := range source.All() {
		membership[value] = false
	}
	return membership
}

// MatchHashMembership marks membership entries yielded by source and reports
// whether every yielded value belongs to membership.
func MatchHashMembership[T comparable](
	membership map[T]bool,
	source Source[T],
) (matched int, contained bool) {
	if sourceSet, ok := source.(*HashSet[T]); ok {
		for value := range sourceSet.values {
			seen, exists := membership[value]
			if !exists {
				return 0, false
			}
			if !seen {
				membership[value] = true
				matched++
			}
		}
		return matched, true
	}
	for value := range source.All() {
		seen, exists := membership[value]
		if !exists {
			return 0, false
		}
		if !seen {
			membership[value] = true
			matched++
		}
	}
	return matched, true
}

// MatchHashSubset reports whether source yields every membership entry at
// least once and marks entries as they are matched.
func MatchHashSubset[T comparable](
	membership map[T]bool,
	source Source[T],
) bool {
	if len(membership) == 0 {
		return true
	}
	matched := 0
	if sourceSet, ok := source.(*HashSet[T]); ok {
		for value := range sourceSet.values {
			seen, exists := membership[value]
			if exists && !seen {
				membership[value] = true
				matched++
				if matched == len(membership) {
					return true
				}
			}
		}
		return false
	}
	for value := range source.All() {
		seen, exists := membership[value]
		if exists && !seen {
			membership[value] = true
			matched++
			if matched == len(membership) {
				return true
			}
		}
	}
	return false
}

// MatchHashContainedMembership marks membership entries yielded by source and
// reports whether source also yields a value outside membership.
func MatchHashContainedMembership[T comparable](
	membership map[T]bool,
	source Source[T],
) (matched int, hasExtra bool) {
	if sourceSet, ok := source.(*HashSet[T]); ok {
		for value := range sourceSet.values {
			seen, exists := membership[value]
			if !exists {
				hasExtra = true
			} else if !seen {
				membership[value] = true
				matched++
			}
			if matched == len(membership) && hasExtra {
				return matched, true
			}
		}
		return matched, hasExtra
	}
	for value := range source.All() {
		seen, exists := membership[value]
		if !exists {
			hasExtra = true
		} else if !seen {
			membership[value] = true
			matched++
		}
		if matched == len(membership) && hasExtra {
			return matched, true
		}
	}
	return matched, hasExtra
}
