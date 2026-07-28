package set

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sliceCollection[T any] []T

type sliceSource[T any] []T

func (values sliceSource[T]) Length() int {
	return len(values)
}

func (values sliceSource[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func (values sliceCollection[T]) IsEmpty() bool {
	return len(values) == 0
}

func (values sliceCollection[T]) Length() int {
	return len(values)
}

func (values sliceCollection[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func (values sliceCollection[T]) Enumerate() iter.Seq2[int, T] {
	return slices.All(values)
}

func TestCollectHashKeys(t *testing.T) {
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectHashKeys(sliceSource[int]{1, 1, 2}))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectHashKeys(NewHashSetFromSlice([]int{1, 2})))
}

func TestCollectHashMembership(t *testing.T) {
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectHashMembership(sliceSource[int]{1, 1, 2}))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectHashMembership(NewHashSetFromSlice([]int{1, 2})))
}

func TestMatchHashMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, contained := MatchHashMembership(
		membership,
		sliceSource[int]{2, 2, 1},
	)

	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, contained = MatchHashMembership(
		map[int]bool{1: false},
		sliceSource[int]{1, 2},
	)
	assert.Zero(t, matched)
	assert.False(t, contained)

	matched, contained = MatchHashMembership(
		map[int]bool{1: false, 2: false},
		NewHashSetFromSlice([]int{2, 1}),
	)
	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	matched, contained = MatchHashMembership(
		map[int]bool{1: false},
		NewHashSetFromSlice([]int{1, 2}),
	)
	assert.Zero(t, matched)
	assert.False(t, contained)
}

func TestMatchHashSubset(t *testing.T) {
	assert.True(t, MatchHashSubset(
		map[int]bool{},
		sliceSource[int]{1},
	))
	assert.True(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		sliceSource[int]{2, 2, 3, 1},
	))
	assert.False(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		sliceSource[int]{1, 3},
	))
	assert.True(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		NewHashSetFromSlice([]int{2, 3, 1}),
	))
	assert.False(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		NewHashSetFromSlice([]int{1, 3}),
	))
}

func TestMatchHashContainedMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, hasExtra := MatchHashContainedMembership(
		membership,
		sliceSource[int]{1, 3, 2, 4},
	)

	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, hasExtra = MatchHashContainedMembership(
		map[int]bool{1: false, 2: false},
		sliceSource[int]{1},
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)

	matched, hasExtra = MatchHashContainedMembership(
		map[int]bool{1: false, 2: false},
		NewHashSetFromSlice([]int{1, 3, 2}),
	)
	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	matched, hasExtra = MatchHashContainedMembership(
		map[int]bool{1: false, 2: false},
		NewHashSetFromSlice([]int{1}),
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)
}

func TestCollectKeyedKeys(t *testing.T) {
	differentIdentity := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, recordNameLength)

	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(
		sliceSource[record]{{ID: 1}, {ID: 1}, {ID: 2}},
		recordID,
	))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(differentIdentity, recordID))
}

func TestCollectKeyedMembership(t *testing.T) {
	differentIdentity := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, recordNameLength)

	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(
		sliceSource[record]{{ID: 1}, {ID: 1}, {ID: 2}},
		recordID,
	))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(differentIdentity, recordID))
}

func TestMatchKeyedMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, contained := MatchKeyedMembership(
		membership,
		sliceSource[record]{{ID: 2}, {ID: 2}, {ID: 1}},
		recordID,
	)

	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false},
		sliceSource[record]{{ID: 1}, {ID: 2}},
		recordID,
	)
	assert.Zero(t, matched)
	assert.False(t, contained)

	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 2}, {ID: 1}}, recordID),
		recordID,
	)
	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	)
	assert.Zero(t, matched)
	assert.False(t, contained)
}

func TestMatchKeyedSubset(t *testing.T) {
	assert.True(t, MatchKeyedSubset(
		map[int]bool{},
		sliceSource[record]{{ID: 1}},
		recordID,
	))
	assert.True(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		sliceSource[record]{{ID: 2}, {ID: 2}, {ID: 3}, {ID: 1}},
		recordID,
	))
	assert.False(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		sliceSource[record]{{ID: 1}, {ID: 3}},
		recordID,
	))
	assert.True(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice(
			[]record{{ID: 2}, {ID: 3}, {ID: 1}},
			recordID,
		),
		recordID,
	))
	assert.False(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 3}}, recordID),
		recordID,
	))
}

func TestMatchKeyedContainedMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, hasExtra := MatchKeyedContainedMembership(
		membership,
		sliceSource[record]{{ID: 1}, {ID: 3}, {ID: 2}, {ID: 4}},
		recordID,
	)

	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		sliceSource[record]{{ID: 1}},
		recordID,
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)

	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice(
			[]record{{ID: 1}, {ID: 3}, {ID: 2}},
			recordID,
		),
		recordID,
	)
	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID),
		recordID,
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)
}

func TestNewHashSetIntersection(t *testing.T) {
	shortest := NewHashSetFromSlice([]int{2, 3})

	result := NewHashSetIntersection(
		sliceSource[int]{1, 2, 3, 4},
		shortest,
		sliceSource[int]{2, 2, 5},
	)

	assert.Equal(t, []int{2}, slices.Collect(result.All()))
	assert.True(t, NewHashSetIntersection[int]().IsEmpty())
	assert.True(t, NewHashSetIntersection(
		sliceCollection[int]{1},
		sliceCollection[int]{},
	).IsEmpty())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(NewHashSetIntersection(
		sliceCollection[int]{1, 2, 2},
	).All()))
	assert.Equal(t, []int{2}, slices.Collect(NewHashSetIntersection(
		NewHashSetFromSlice([]int{1, 2, 3}),
		NewHashSetFromSlice([]int{2, 3}),
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.True(t, NewHashSetIntersection(
		NewHashSetFromSlice([]int{1, 2}),
		sliceCollection[int]{3},
		sliceCollection[int]{1, 2, 3},
	).IsEmpty())
}

func TestNewHashSetDifference(t *testing.T) {
	assert.True(t, NewHashSetDifference[int]().IsEmpty())
	assert.Equal(t, []int{1}, slices.Collect(NewHashSetDifference(
		NewHashSetFromSlice([]int{1, 2, 3, 4}),
		NewHashSetFromSlice([]int{2}),
		NewHashSetFromSlice([]int{3, 4, 5, 6}),
		sliceSource[int]{9},
	).All()))
	assert.True(t, NewHashSetDifference(
		sliceCollection[int]{1, 2, 3},
		NewHashSetFromSlice([]int{1, 2, 3, 4}),
		sliceCollection[int]{4},
	).IsEmpty())
	assert.Equal(t, []int{2}, slices.Collect(NewHashSetDifference(
		NewHashSetFromSlice([]int{1, 2}),
		sliceCollection[int]{1, 1},
	).All()))
}

func TestNewHashSetSymmetricDifference(t *testing.T) {
	first := NewHashSetFromSlice([]int{1, 2})

	result := NewHashSetSymmetricDifference(
		first,
		sliceSource[int]{2, 3, 3},
		sliceSource[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
	assert.True(t, NewHashSetSymmetricDifference[int]().IsEmpty())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(NewHashSetSymmetricDifference(
		sliceCollection[int]{1, 2, 2},
	).All()))
}

func TestHashSet_Clone(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2})
	clone := original.Clone()

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(clone.All()))
	assert.True(t, clone.Add(3))
	assert.False(t, original.Contains(3))

	var zero HashSet[int]
	assert.True(t, zero.Clone().Add(1))
}

func TestHashSet_Union(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2})
	result := original.Union(sliceCollection[int]{2, 3, 3})
	self := original.Union(original)

	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(original.Union().All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(self.All()))
	assert.True(t, self.Add(3))
	assert.False(t, original.Contains(3))
}

func TestHashSet_Intersection(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2, 3})
	result := original.Intersection(
		sliceCollection[int]{2, 3, 3, 4},
		sliceCollection[int]{3, 3},
	)

	assert.Equal(t, []int{3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Intersection().All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Intersection(original).All()))
	assert.Equal(t, []int{2}, slices.Collect(original.Intersection(
		NewHashSetFromSlice([]int{2}),
	).All()))
	assert.True(t, original.Intersection(sliceCollection[int]{}).IsEmpty())
	assert.Equal(t, []int{2}, slices.Collect(original.Intersection(
		NewHashSetFromSlice([]int{2, 3}),
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.True(t, original.Intersection(
		NewHashSetFromSlice([]int{1, 2}),
		sliceCollection[int]{3},
	).IsEmpty())

	var empty HashSet[int]
	assert.True(t, empty.Intersection(sliceCollection[int]{1}).IsEmpty())
}

func TestHashSet_Difference(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2, 3})
	result := original.Difference(
		sliceCollection[int]{2, 2},
		sliceCollection[int]{4},
	)

	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Difference().All()))
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(original.Difference(
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.Equal(t, []int{3}, slices.Collect(original.Difference(
		sliceCollection[int]{2},
		NewHashSetFromSlice([]int{1}),
	).All()))
	assert.Equal(t, []int{1}, slices.Collect(original.Difference(
		sliceCollection[int]{2},
		NewHashSetFromSlice([]int{3, 4, 5, 6}),
	).All()))
	assert.True(t, original.Difference(
		sliceCollection[int]{4},
		original,
	).IsEmpty())
	assert.True(t, original.Difference(original).IsEmpty())
	assert.True(t, original.Difference(sliceCollection[int]{1, 2, 3}).IsEmpty())
}

func TestHashSet_SymmetricDifference(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2})
	result := original.SymmetricDifference(
		sliceCollection[int]{2, 3, 3},
		sliceCollection[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(original.SymmetricDifference().All()))
	assert.True(t, original.SymmetricDifference(original).IsEmpty())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(
		original.SymmetricDifference(original, original).All(),
	))
	assert.Equal(t, []int{1}, slices.Collect(original.SymmetricDifference(
		NewHashSetFromSlice([]int{2}),
	).All()))
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(original.SymmetricDifference(
		NewHashSetFromSlice([]int{2, 3}),
	).All()))
}

func TestHashSet_Equal(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.Equal(sliceCollection[int]{1, 2, 2}))
	assert.True(t, set.Equal(NewHashSetFromSlice([]int{1, 2})))
	assert.True(t, set.Equal(set))
	assert.False(t, set.Equal(sliceCollection[int]{1}))
	assert.False(t, set.Equal(sliceCollection[int]{1, 3}))
	assert.False(t, set.Equal(NewHashSetFromSlice([]int{1, 3})))
}

func TestHashSet_IsSubset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.IsSubset(sliceCollection[int]{1, 2, 2, 3}))
	assert.True(t, set.IsSubset(NewHashSetFromSlice([]int{1, 2, 3})))
	assert.True(t, set.IsSubset(set))
	assert.False(t, set.IsSubset(sliceCollection[int]{1}))
	assert.False(t, set.IsSubset(sliceCollection[int]{1, 3}))
	assert.False(t, set.IsSubset(NewHashSetFromSlice([]int{1, 3, 4})))
}

func TestHashSet_IsProperSubset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.False(t, set.IsProperSubset(sliceCollection[int]{1}))
	assert.True(t, set.IsProperSubset(sliceCollection[int]{1, 2, 2, 3}))
	assert.True(t, set.IsProperSubset(NewHashSetFromSlice([]int{1, 2, 3})))
	assert.False(t, set.IsProperSubset(sliceCollection[int]{1, 2, 2}))
	assert.False(t, set.IsProperSubset(sliceCollection[int]{1, 3, 4}))
	assert.False(t, set.IsProperSubset(NewHashSetFromSlice([]int{1, 3, 4})))
}

func TestHashSet_IsSuperset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.True(t, set.IsSuperset(sliceCollection[int]{1, 2, 2}))
	assert.True(t, set.IsSuperset(NewHashSetFromSlice([]int{1, 2})))
	assert.True(t, set.IsSuperset(set))
	assert.False(t, set.IsSuperset(sliceCollection[int]{1, 4}))
	assert.False(t, set.IsSuperset(NewHashSetFromSlice([]int{1, 4})))
}

func TestHashSet_IsProperSuperset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.True(t, set.IsProperSuperset(sliceCollection[int]{1, 2, 2}))
	assert.True(t, set.IsProperSuperset(NewHashSetFromSlice([]int{1, 2})))
	assert.False(t, set.IsProperSuperset(sliceCollection[int]{1, 2, 3, 3}))
	assert.False(t, set.IsProperSuperset(sliceCollection[int]{1, 4}))
	assert.False(t, set.IsProperSuperset(NewHashSetFromSlice([]int{1, 4})))
}

func TestHashSet_IsDisjoint(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.IsDisjoint(sliceCollection[int]{3, 4, 4}))
	assert.False(t, set.IsDisjoint(sliceCollection[int]{2, 3}))
	assert.True(t, set.IsDisjoint(NewHashSetFromSlice([]int{3})))
	assert.False(t, set.IsDisjoint(NewHashSetFromSlice([]int{2, 3, 4})))
	assert.False(t, NewHashSetFromSlice([]int{1, 2, 3}).IsDisjoint(set))
	assert.False(t, set.IsDisjoint(set))
	assert.True(t, NewHashSet[int](0).IsDisjoint(NewHashSet[int](0)))
}

func TestHashSet_UnionWith(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.Equal(t, 2, set.UnionWith(
		sliceCollection[int]{2, 3, 3},
		sliceCollection[int]{4},
	))
	assert.Zero(t, set.UnionWith())
	assert.Zero(t, set.UnionWith(set))
	assert.Zero(t, set.UnionWith(set, set))
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, slices.Collect(set.All()))

	var empty HashSet[int]
	assert.Equal(t, 2, empty.UnionWith(NewHashSetFromSlice([]int{5, 6})))
	assert.ElementsMatch(t, []int{5, 6}, slices.Collect(empty.All()))
	assert.Equal(t, 1, empty.UnionWith(NewHashSetFromSlice([]int{6, 7})))
	assert.ElementsMatch(t, []int{5, 6, 7}, slices.Collect(empty.All()))

	preallocated := NewHashSet[int](4)
	assert.Equal(t, 2, preallocated.UnionWith(NewHashSetFromSlice([]int{8, 9})))
	assert.ElementsMatch(t, []int{8, 9}, slices.Collect(preallocated.All()))
}

func TestHashSet_IntersectWith(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3, 4})

	assert.Equal(t, 3, set.IntersectWith(
		sliceCollection[int]{2, 3, 4, 4},
		sliceCollection[int]{3},
	))
	assert.Zero(t, set.IntersectWith())
	assert.Zero(t, set.IntersectWith(set))
	assert.Equal(t, []int{3}, slices.Collect(set.All()))

	assert.Equal(t, 1, set.IntersectWith(sliceCollection[int]{}))
	assert.True(t, set.IsEmpty())

	set.Adds(1, 2)
	assert.Equal(t, 2, set.IntersectWith(
		NewHashSetFromSlice([]int{1}),
		NewHashSetFromSlice([]int{2}),
	))
	assert.True(t, set.IsEmpty())

	set.Adds(1, 2, 3)
	assert.Equal(t, 2, set.IntersectWith(NewHashSetFromSlice([]int{2})))
	assert.Equal(t, []int{2}, slices.Collect(set.All()))
}

func TestHashSet_DifferenceWith(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3, 4})

	assert.Equal(t, 2, set.DifferenceWith(
		sliceCollection[int]{2, 2},
		sliceCollection[int]{4, 5},
	))
	assert.Zero(t, set.DifferenceWith())
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(set.All()))

	assert.Equal(t, 2, set.DifferenceWith(set))
	assert.True(t, set.IsEmpty())
	assert.Zero(t, set.DifferenceWith(NewHashSetFromSlice([]int{1})))

	set.Adds(1, 2, 3)
	assert.Equal(t, 1, set.DifferenceWith(NewHashSetFromSlice([]int{2, 4})))
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(set.All()))
	assert.Equal(t, 2, set.DifferenceWith(NewHashSetFromSlice([]int{1, 3})))
	assert.True(t, set.IsEmpty())

	set.Adds(1, 2)
	assert.Equal(t, 2, set.DifferenceWith(sliceCollection[int]{1, 2}))
	assert.True(t, set.IsEmpty())
}

func TestHashSet_SymmetricDifferenceWith(t *testing.T) {
	var set HashSet[int]

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		sliceCollection[int]{1, 2, 2},
		sliceCollection[int]{2, 3, 3},
	))
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(set.All()))
	assert.Zero(t, set.SymmetricDifferenceWith())

	assert.Equal(t, 2, set.SymmetricDifferenceWith(&set))
	assert.True(t, set.IsEmpty())

	set.Add(1)
	assert.Equal(t, 2, set.SymmetricDifferenceWith(NewHashSetFromSlice([]int{1, 2})))
	assert.Equal(t, []int{2}, slices.Collect(set.All()))

	set.Clear()
	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		sliceCollection[int]{1},
		NewHashSetFromSlice([]int{1, 2}),
	))
	assert.Equal(t, []int{2}, slices.Collect(set.All()))
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		sliceCollection[int]{1},
		sliceCollection[int]{2},
		sliceCollection[int]{1, 3, 3},
	))
	assert.Equal(t, []int{3}, slices.Collect(set.All()))

	set.Clear()
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		NewHashSetFromSlice([]int{1}),
		sliceCollection[int]{2},
	))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestNewKeyedHashSetIntersection(t *testing.T) {
	result := NewKeyedHashSetIntersection(
		recordID,
		sliceSource[record]{{ID: 1}, {ID: 2, Name: "first"}, {ID: 2}},
		sliceSource[record]{{ID: 2}, {ID: 3}},
		sliceSource[record]{{ID: 2}},
	)

	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(result.All()))
	assert.True(t, NewKeyedHashSetIntersection(recordID).IsEmpty())
	assert.True(t, NewKeyedHashSetIntersection(
		recordID,
		sliceCollection[record]{},
		sliceCollection[record]{{ID: 1}},
	).IsEmpty())
	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 2}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}},
		).All(),
	))
	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			sliceCollection[record]{{ID: 2, Name: "first"}},
			sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 3}},
		).All(),
	))
	filter := NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID)
	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			sliceCollection[record]{{ID: 1}, {ID: 2, Name: "first"}},
			filter,
		).All(),
	))
	assert.True(t, NewKeyedHashSetIntersection(
		recordID,
		sliceCollection[record]{{ID: 1}, {ID: 2}},
		sliceCollection[record]{{ID: 1}},
		sliceCollection[record]{{ID: 2}},
	).IsEmpty())
}

func TestNewKeyedHashSetDifference(t *testing.T) {
	assert.True(t, NewKeyedHashSetDifference(recordID).IsEmpty())
	assert.Equal(t, []record{{ID: 1, Name: "first"}}, slices.Collect(
		NewKeyedHashSetDifference(
			recordID,
			sliceSource[record]{{ID: 1, Name: "first"}, {ID: 2}},
			NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID),
		).All(),
	))
	assert.True(t, NewKeyedHashSetDifference(
		recordID,
		sliceCollection[record]{{ID: 1}, {ID: 2}},
		sliceCollection[record]{{ID: 1}, {ID: 2}},
		sliceCollection[record]{{ID: 3}},
	).IsEmpty())
}

func TestNewKeyedHashSetSymmetricDifference(t *testing.T) {
	result := NewKeyedHashSetSymmetricDifference(
		recordID,
		sliceSource[record]{{ID: 1}, {ID: 2}, {ID: 2}},
		sliceSource[record]{{ID: 2}, {ID: 3, Name: "three"}},
	)

	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 3, Name: "three"}}, slices.Collect(result.All()))
	assert.True(t, NewKeyedHashSetSymmetricDifference(recordID).IsEmpty())
}

func TestKeyedHashSet_Clone(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
	}, recordID)
	clone := original.Clone()

	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(clone.All()))
	assert.True(t, clone.Add(record{ID: 3, Name: "three"}))
	assert.False(t, original.Contains(record{ID: 3}))
	assert.Equal(t, 3, clone.keyOf(record{ID: 3}))
}

func TestKeyedHashSet_Union(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
	}, recordID)
	result := original.Union(sliceCollection[record]{
		{ID: 2, Name: "right"},
		{ID: 3, Name: "three"},
		{ID: 3, Name: "duplicate"},
	})
	self := original.Union(original)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Union().All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(self.All()))
	assert.True(t, self.Add(record{ID: 3}))
	assert.False(t, original.Contains(record{ID: 3}))
}

func TestKeyedHashSet_Intersection(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
	}, recordID)
	result := original.Intersection(
		sliceCollection[record]{{ID: 2, Name: "right"}, {ID: 3}},
		sliceCollection[record]{{ID: 2, Name: "again"}},
	)

	assert.Equal(t, []record{{ID: 2, Name: "left"}}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Intersection().All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(
		original.Intersection(original).All(),
	))
	assert.True(t, original.Intersection(sliceCollection[record]{}).IsEmpty())
	assert.True(t, original.Intersection(
		sliceCollection[record]{{ID: 1}},
		sliceCollection[record]{{ID: 2}},
	).IsEmpty())
	assert.Equal(t, []record{{ID: 2, Name: "left"}}, slices.Collect(
		original.Intersection(NewKeyedHashSetFromSlice(
			[]record{{ID: 2}},
			recordID,
		)).All(),
	))

	empty := NewKeyedHashSet(0, recordID)
	assert.True(t, empty.Intersection(sliceCollection[record]{{ID: 1}}).IsEmpty())
}

func TestKeyedHashSet_Difference(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
		{ID: 3, Name: "three"},
	}, recordID)
	result := original.Difference(
		sliceCollection[record]{{ID: 2}, {ID: 2}},
		sliceCollection[record]{{ID: 4}},
	)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "three"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Difference().All()))
	assert.ElementsMatch(t, []record{{ID: 1, Name: "one"}, {ID: 3, Name: "three"}}, slices.Collect(
		original.Difference(NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID)).All(),
	))
	assert.True(t, original.Difference(original).IsEmpty())
	assert.True(t, original.Difference(
		sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 3}},
	).IsEmpty())
}

func TestKeyedHashSet_SymmetricDifference(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
	}, recordID)
	result := original.SymmetricDifference(
		sliceCollection[record]{{ID: 2}, {ID: 3, Name: "first"}},
		sliceCollection[record]{{ID: 3, Name: "second"}},
		sliceCollection[record]{{ID: 3, Name: "third"}, {ID: 4, Name: "four"}},
	)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "third"},
		{ID: 4, Name: "four"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.SymmetricDifference().All()))
	assert.True(t, original.SymmetricDifference(original).IsEmpty())
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(
		original.SymmetricDifference(original, original).All(),
	))
	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3},
	}, slices.Collect(original.SymmetricDifference(
		NewKeyedHashSetFromSlice([]record{{ID: 2}, {ID: 3}}, recordID),
	).All()))
	assert.Equal(t, []record{
		{ID: 1, Name: "one"},
	}, slices.Collect(original.SymmetricDifference(
		sliceCollection[record]{{ID: 2}, {ID: 2}},
	).All()))
	byNameLength := func(value record) int {
		return len(value.Name)
	}
	assert.Equal(t, 1, original.SymmetricDifference(
		NewKeyedHashSetFromSlice([]record{
			{ID: 1, Name: "a"},
			{ID: 1, Name: "bb"},
		}, byNameLength),
	).Length())
}

func TestKeyedHashSet_RelationsReapplyReceiverIdentity(t *testing.T) {
	byID := NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID)
	byNameLength := func(value record) int {
		return len(value.Name)
	}
	sameIDs := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
	}, byNameLength)

	assert.True(t, byID.Equal(sameIDs))
	assert.True(t, byID.IsSubset(sameIDs))
	assert.False(t, byID.IsProperSubset(sameIDs))
	assert.True(t, byID.IsSuperset(sameIDs))
	assert.False(t, byID.IsProperSuperset(sameIDs))
	assert.False(t, byID.IsDisjoint(sameIDs))

	withExtra := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "bb"},
	}, byNameLength)
	assert.False(t, byID.Equal(withExtra))
	assert.True(t, byID.IsProperSubset(withExtra))
	assert.False(t, byID.IsSuperset(withExtra))
	byIDWithExtra := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)
	assert.True(t, byIDWithExtra.IsProperSuperset(sameIDs))
	assert.False(t, byIDWithExtra.IsSubset(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}, {ID: 3}},
		recordID,
	)))
	assert.False(t, byIDWithExtra.IsProperSuperset(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}, {ID: 3}},
		recordID,
	)))
	assert.True(t, byID.IsDisjoint(NewKeyedHashSetFromSlice(
		[]record{{ID: 2, Name: "a"}},
		byNameLength,
	)))
}

func TestKeyedHashSet_Equal(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.Equal(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.True(t, set.Equal(set))
	assert.False(t, set.Equal(sliceCollection[record]{{ID: 1}}))
	assert.False(t, set.Equal(sliceCollection[record]{{ID: 1}, {ID: 3}}))
}

func TestKeyedHashSet_IsSubset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.IsSubset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}}))
	assert.True(t, set.IsSubset(set))
	assert.False(t, set.IsSubset(sliceCollection[record]{{ID: 1}}))
	assert.False(t, set.IsSubset(sliceCollection[record]{{ID: 1}, {ID: 3}}))
}

func TestKeyedHashSet_IsProperSubset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.False(t, set.IsProperSubset(sliceCollection[record]{{ID: 1}}))
	assert.True(t, set.IsProperSubset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 3}}))
	assert.False(t, set.IsProperSubset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.False(t, set.IsProperSubset(sliceCollection[record]{{ID: 1}, {ID: 3}, {ID: 4}}))
}

func TestKeyedHashSet_IsSuperset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}}, recordID)

	assert.True(t, set.IsSuperset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.True(t, set.IsSuperset(set))
	assert.False(t, set.IsSuperset(sliceCollection[record]{{ID: 1}, {ID: 4}}))
}

func TestKeyedHashSet_IsProperSuperset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}}, recordID)

	assert.True(t, set.IsProperSuperset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.False(t, set.IsProperSuperset(sliceCollection[record]{{ID: 1}, {ID: 2}, {ID: 3}}))
	assert.False(t, set.IsProperSuperset(set))
	assert.False(t, set.IsProperSuperset(sliceCollection[record]{{ID: 1}, {ID: 4}}))
}

func TestKeyedHashSet_IsDisjoint(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.IsDisjoint(sliceCollection[record]{{ID: 3}, {ID: 4}, {ID: 4}}))
	assert.False(t, set.IsDisjoint(sliceCollection[record]{{ID: 2}, {ID: 3}}))
	assert.False(t, set.IsDisjoint(set))
	empty := NewKeyedHashSet(0, recordID)
	assert.True(t, empty.IsDisjoint(empty))
}

func TestKeyedHashSet_UnionWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2, Name: "left"}}, recordID)

	assert.Equal(t, 2, set.UnionWith(
		sliceCollection[record]{{ID: 2, Name: "right"}, {ID: 3, Name: "three"}},
		sliceCollection[record]{{ID: 4, Name: "four"}},
	))
	assert.Zero(t, set.UnionWith())
	assert.Zero(t, set.UnionWith(set))
	assert.Zero(t, set.UnionWith(set, set))
	assert.Equal(t, 1, set.UnionWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 5, Name: "five"}},
		recordID,
	)))
	assert.ElementsMatch(t, []record{
		{ID: 1},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
		{ID: 4, Name: "four"},
		{ID: 5, Name: "five"},
	}, slices.Collect(set.All()))
}

func TestKeyedHashSet_IntersectWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{
		{ID: 1},
		{ID: 2},
		{ID: 3, Name: "stored"},
		{ID: 4},
	}, recordID)

	assert.Equal(t, 3, set.IntersectWith(
		sliceCollection[record]{{ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}},
		sliceCollection[record]{{ID: 3, Name: "query"}},
	))
	assert.Zero(t, set.IntersectWith())
	assert.Zero(t, set.IntersectWith(set))
	assert.Equal(t, []record{{ID: 3, Name: "stored"}}, slices.Collect(set.All()))

	assert.Equal(t, 1, set.IntersectWith(sliceCollection[record]{}))
	assert.True(t, set.IsEmpty())

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 2, set.IntersectWith(
		sliceCollection[record]{{ID: 1}},
		sliceCollection[record]{{ID: 2}},
	))
	assert.True(t, set.IsEmpty())
}

func TestKeyedHashSet_DifferenceWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}, recordID)

	assert.Equal(t, 2, set.DifferenceWith(
		sliceCollection[record]{{ID: 2}, {ID: 2}},
		sliceCollection[record]{{ID: 4}, {ID: 5}},
	))
	assert.Zero(t, set.DifferenceWith())
	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 3}}, slices.Collect(set.All()))

	assert.Equal(t, 2, set.DifferenceWith(set))
	assert.True(t, set.IsEmpty())
	assert.Zero(t, set.DifferenceWith(sliceCollection[record]{{ID: 1}}))

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 2, set.DifferenceWith(sliceCollection[record]{{ID: 1}, {ID: 2}}))
	assert.True(t, set.IsEmpty())

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 1, set.DifferenceWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 2}},
		recordID,
	)))
	assert.Equal(t, []record{{ID: 1}}, slices.Collect(set.All()))
	assert.Equal(t, 1, set.DifferenceWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}},
		recordID,
	)))
	assert.True(t, set.IsEmpty())
}

func TestKeyedHashSet_SymmetricDifferenceWith(t *testing.T) {
	set := KeyedHashSet[record, int]{keyOf: recordID}

	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		NewKeyedHashSetFromSlice([]record{{ID: 4}}, recordID),
	))
	assert.Equal(t, []record{{ID: 4}}, slices.Collect(set.All()))
	set.Clear()

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		sliceCollection[record]{{ID: 1, Name: "one"}, {ID: 2}, {ID: 2}},
		sliceCollection[record]{{ID: 2}, {ID: 3, Name: "three"}, {ID: 3}},
	))
	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "three"},
	}, slices.Collect(set.All()))
	assert.Zero(t, set.SymmetricDifferenceWith())

	assert.Equal(t, 2, set.SymmetricDifferenceWith(&set))
	assert.True(t, set.IsEmpty())

	otherKey := func(value record) int {
		return len(value.Name)
	}
	other := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, otherKey)

	single := NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID)
	assert.Equal(t, 2, single.SymmetricDifferenceWith(other))
	assert.Equal(t, []record{{ID: 2, Name: "ccc"}}, slices.Collect(single.All()))

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		other,
		sliceCollection[record]{{ID: 2}},
		sliceCollection[record]{{ID: 3}},
	))
	assert.Equal(t, 2, set.Length())
	assert.True(t, set.Contains(record{ID: 1}))
	assert.True(t, set.Contains(record{ID: 3}))

	set.Clear()
	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		sliceCollection[record]{{ID: 1}},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
	))
	assert.Equal(t, []record{{ID: 2}}, slices.Collect(set.All()))

	set.Clear()
	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		sliceCollection[record]{{ID: 2}},
		other,
	))
	assert.Equal(t, 1, set.Length())
	assert.True(t, set.Contains(record{ID: 1}))

	set.Clear()
	set.Add(record{ID: 1})
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		sliceCollection[record]{{ID: 1}},
		sliceCollection[record]{{ID: 2}},
	))
	assert.Equal(t, []record{{ID: 2}}, slices.Collect(set.All()))
}
