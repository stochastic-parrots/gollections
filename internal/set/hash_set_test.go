package set

import (
	"encoding/json"
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashSet(t *testing.T) {
	set := NewHashSet[int](4)

	assert.NotNil(t, set.values)
	assert.Empty(t, set.values)
}

func TestNewHashSetFromSlice(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 1})

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestNewHashSetFromSeq(t *testing.T) {
	set := NewHashSetFromSeq(slices.Values([]int{1, 2, 1}))

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestHashSet_Length(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.Equal(t, 2, set.Length())
}

func TestHashSet_IsEmpty(t *testing.T) {
	set := NewHashSet[int](0)
	assert.True(t, set.IsEmpty())

	set.Add(1)
	assert.False(t, set.IsEmpty())
}

func TestHashSet_Contains(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.Contains(1))
	assert.False(t, set.Contains(3))
}

func TestHashSet_Add(t *testing.T) {
	var set HashSet[int]

	assert.True(t, set.Add(1))
	assert.False(t, set.Add(1))

	assert.Equal(t, 1, set.Length())
	assert.True(t, set.Contains(1))
}

func TestHashSet_Adds(t *testing.T) {
	set := NewHashSet[int](0)

	assert.Equal(t, 2, set.Adds(1, 2, 1))
	assert.Zero(t, set.Adds())
	assert.Zero(t, set.Adds(1, 2))

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestHashSet_Remove(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.Remove(1))
	assert.False(t, set.Remove(1))
	assert.False(t, set.Contains(1))
	assert.Equal(t, 1, set.Length())
}

func TestHashSet_Removes(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.Equal(t, 2, set.Removes(2, 2, 4, 3))
	assert.Zero(t, set.Removes())
	assert.Zero(t, set.Removes(2, 4))
	assert.Equal(t, []int{1}, slices.Collect(set.All()))
}

func TestHashSet_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2, 3})

		assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(set.All()))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2, 3})
		count := 0

		set.All()(func(int) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		assert.Empty(t, slices.Collect(set.All()))
	})
}

func TestHashSet_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{10, 20, 30})
		var indexes []int
		var values []int

		for idx, value := range set.Enumerate() {
			indexes = append(indexes, idx)
			values = append(values, value)
		}

		assert.Equal(t, []int{0, 1, 2}, indexes)
		assert.ElementsMatch(t, []int{10, 20, 30}, values)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{10, 20, 30})
		count := 0

		set.Enumerate()(func(int, int) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		assert.Empty(t, slices.Collect(func(yield func(int) bool) {
			for _, value := range set.Enumerate() {
				if !yield(value) {
					return
				}
			}
		}))
	})
}

func TestHashSet_Clear(t *testing.T) {
	first, second, third := 1, 2, 3
	set := NewHashSetFromSlice([]*int{&first, &second})
	storage := set.values

	set.Clear()

	assert.True(t, set.IsEmpty())
	assert.Empty(t, storage)

	set.Add(&third)
	_, reused := storage[&third]
	assert.True(t, reused)
}

func TestHashSet_MarshalJSON(t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2})

		data, err := set.MarshalJSON()
		assert.NoError(t, err)

		var values []int
		assert.NoError(t, json.Unmarshal(data, &values))
		assert.ElementsMatch(t, []int{1, 2}, values)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		data, err := set.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("UnsupportedValue", func(t *testing.T) {
		set := NewHashSet[chan int](0)
		set.Add(make(chan int))

		_, err := set.MarshalJSON()
		assert.Error(t, err)
	})
}

type hashSource[T any] []T

func (values hashSource[T]) Length() int {
	return len(values)
}

func (values hashSource[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func TestCollectHashKeys(t *testing.T) {
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectHashKeys(hashSource[int]{1, 1, 2}))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectHashKeys(NewHashSetFromSlice([]int{1, 2})))
}

func TestCollectHashMembership(t *testing.T) {
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectHashMembership(hashSource[int]{1, 1, 2}))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectHashMembership(NewHashSetFromSlice([]int{1, 2})))
}

func TestMatchHashMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, contained := MatchHashMembership(
		membership,
		hashSource[int]{2, 2, 1},
	)

	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, contained = MatchHashMembership(
		map[int]bool{1: false},
		hashSource[int]{1, 2},
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
		hashSource[int]{1},
	))
	assert.True(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		hashSource[int]{2, 2, 3, 1},
	))
	assert.False(t, MatchHashSubset(
		map[int]bool{1: false, 2: false},
		hashSource[int]{1, 3},
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
		hashSource[int]{1, 3, 2, 4},
	)

	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, hasExtra = MatchHashContainedMembership(
		map[int]bool{1: false, 2: false},
		hashSource[int]{1},
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

func TestNewHashSetIntersection(t *testing.T) {
	shortest := NewHashSetFromSlice([]int{2, 3})

	result := NewHashSetIntersection(
		hashSource[int]{1, 2, 3, 4},
		shortest,
		hashSource[int]{2, 2, 5},
	)

	assert.Equal(t, []int{2}, slices.Collect(result.All()))
	assert.True(t, NewHashSetIntersection[int]().IsEmpty())
	assert.True(t, NewHashSetIntersection(
		hashSource[int]{1},
		hashSource[int]{},
	).IsEmpty())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(NewHashSetIntersection(
		hashSource[int]{1, 2, 2},
	).All()))
	assert.Equal(t, []int{2}, slices.Collect(NewHashSetIntersection(
		NewHashSetFromSlice([]int{1, 2, 3}),
		NewHashSetFromSlice([]int{2, 3}),
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.True(t, NewHashSetIntersection(
		NewHashSetFromSlice([]int{1, 2}),
		hashSource[int]{3},
		hashSource[int]{1, 2, 3},
	).IsEmpty())
}

func TestNewHashSetDifference(t *testing.T) {
	assert.True(t, NewHashSetDifference[int]().IsEmpty())
	assert.Equal(t, []int{1}, slices.Collect(NewHashSetDifference(
		NewHashSetFromSlice([]int{1, 2, 3, 4}),
		NewHashSetFromSlice([]int{2}),
		NewHashSetFromSlice([]int{3, 4, 5, 6}),
		hashSource[int]{9},
	).All()))
	assert.True(t, NewHashSetDifference(
		hashSource[int]{1, 2, 3},
		NewHashSetFromSlice([]int{1, 2, 3, 4}),
		hashSource[int]{4},
	).IsEmpty())
	assert.Equal(t, []int{2}, slices.Collect(NewHashSetDifference(
		NewHashSetFromSlice([]int{1, 2}),
		hashSource[int]{1, 1},
	).All()))
}

func TestNewHashSetSymmetricDifference(t *testing.T) {
	first := NewHashSetFromSlice([]int{1, 2})

	result := NewHashSetSymmetricDifference(
		first,
		hashSource[int]{2, 3, 3},
		hashSource[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
	assert.True(t, NewHashSetSymmetricDifference[int]().IsEmpty())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(NewHashSetSymmetricDifference(
		hashSource[int]{1, 2, 2},
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
	result := original.Union(hashSource[int]{2, 3, 3})
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
		hashSource[int]{2, 3, 3, 4},
		hashSource[int]{3, 3},
	)

	assert.Equal(t, []int{3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Intersection().All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Intersection(original).All()))
	assert.Equal(t, []int{2}, slices.Collect(original.Intersection(
		NewHashSetFromSlice([]int{2}),
	).All()))
	assert.True(t, original.Intersection(hashSource[int]{}).IsEmpty())
	assert.Equal(t, []int{2}, slices.Collect(original.Intersection(
		NewHashSetFromSlice([]int{2, 3}),
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.True(t, original.Intersection(
		NewHashSetFromSlice([]int{1, 2}),
		hashSource[int]{3},
	).IsEmpty())

	var empty HashSet[int]
	assert.True(t, empty.Intersection(hashSource[int]{1}).IsEmpty())
}

func TestHashSet_Difference(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2, 3})
	result := original.Difference(
		hashSource[int]{2, 2},
		hashSource[int]{4},
	)

	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.All()))
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(original.Difference().All()))
	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(original.Difference(
		NewHashSetFromSlice([]int{2, 4}),
	).All()))
	assert.Equal(t, []int{3}, slices.Collect(original.Difference(
		hashSource[int]{2},
		NewHashSetFromSlice([]int{1}),
	).All()))
	assert.Equal(t, []int{1}, slices.Collect(original.Difference(
		hashSource[int]{2},
		NewHashSetFromSlice([]int{3, 4, 5, 6}),
	).All()))
	assert.True(t, original.Difference(
		hashSource[int]{4},
		original,
	).IsEmpty())
	assert.True(t, original.Difference(original).IsEmpty())
	assert.True(t, original.Difference(hashSource[int]{1, 2, 3}).IsEmpty())
}

func TestHashSet_SymmetricDifference(t *testing.T) {
	original := NewHashSetFromSlice([]int{1, 2})
	result := original.SymmetricDifference(
		hashSource[int]{2, 3, 3},
		hashSource[int]{3, 4},
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

	assert.True(t, set.Equal(hashSource[int]{1, 2, 2}))
	assert.True(t, set.Equal(NewHashSetFromSlice([]int{1, 2})))
	assert.True(t, set.Equal(set))
	assert.False(t, set.Equal(hashSource[int]{1}))
	assert.False(t, set.Equal(hashSource[int]{1, 3}))
	assert.False(t, set.Equal(NewHashSetFromSlice([]int{1, 3})))
}

func TestHashSet_IsSubset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.IsSubset(hashSource[int]{1, 2, 2, 3}))
	assert.True(t, set.IsSubset(NewHashSetFromSlice([]int{1, 2, 3})))
	assert.True(t, set.IsSubset(set))
	assert.False(t, set.IsSubset(hashSource[int]{1}))
	assert.False(t, set.IsSubset(hashSource[int]{1, 3}))
	assert.False(t, set.IsSubset(NewHashSetFromSlice([]int{1, 3, 4})))
}

func TestHashSet_IsProperSubset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.False(t, set.IsProperSubset(hashSource[int]{1}))
	assert.True(t, set.IsProperSubset(hashSource[int]{1, 2, 2, 3}))
	assert.True(t, set.IsProperSubset(NewHashSetFromSlice([]int{1, 2, 3})))
	assert.False(t, set.IsProperSubset(hashSource[int]{1, 2, 2}))
	assert.False(t, set.IsProperSubset(hashSource[int]{1, 3, 4}))
	assert.False(t, set.IsProperSubset(NewHashSetFromSlice([]int{1, 3, 4})))
}

func TestHashSet_IsSuperset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.True(t, set.IsSuperset(hashSource[int]{1, 2, 2}))
	assert.True(t, set.IsSuperset(NewHashSetFromSlice([]int{1, 2})))
	assert.True(t, set.IsSuperset(set))
	assert.False(t, set.IsSuperset(hashSource[int]{1, 4}))
	assert.False(t, set.IsSuperset(NewHashSetFromSlice([]int{1, 4})))
}

func TestHashSet_IsProperSuperset(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.True(t, set.IsProperSuperset(hashSource[int]{1, 2, 2}))
	assert.True(t, set.IsProperSuperset(NewHashSetFromSlice([]int{1, 2})))
	assert.False(t, set.IsProperSuperset(hashSource[int]{1, 2, 3, 3}))
	assert.False(t, set.IsProperSuperset(hashSource[int]{1, 4}))
	assert.False(t, set.IsProperSuperset(NewHashSetFromSlice([]int{1, 4})))
}

func TestHashSet_IsDisjoint(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.IsDisjoint(hashSource[int]{3, 4, 4}))
	assert.False(t, set.IsDisjoint(hashSource[int]{2, 3}))
	assert.True(t, set.IsDisjoint(NewHashSetFromSlice([]int{3})))
	assert.False(t, set.IsDisjoint(NewHashSetFromSlice([]int{2, 3, 4})))
	assert.False(t, NewHashSetFromSlice([]int{1, 2, 3}).IsDisjoint(set))
	assert.False(t, set.IsDisjoint(set))
	assert.True(t, NewHashSet[int](0).IsDisjoint(NewHashSet[int](0)))
}

func TestHashSet_UnionWith(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.Equal(t, 2, set.UnionWith(
		hashSource[int]{2, 3, 3},
		hashSource[int]{4},
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
		hashSource[int]{2, 3, 4, 4},
		hashSource[int]{3},
	))
	assert.Zero(t, set.IntersectWith())
	assert.Zero(t, set.IntersectWith(set))
	assert.Equal(t, []int{3}, slices.Collect(set.All()))

	assert.Equal(t, 1, set.IntersectWith(hashSource[int]{}))
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
		hashSource[int]{2, 2},
		hashSource[int]{4, 5},
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
	assert.Equal(t, 2, set.DifferenceWith(hashSource[int]{1, 2}))
	assert.True(t, set.IsEmpty())
}

func TestHashSet_SymmetricDifferenceWith(t *testing.T) {
	var set HashSet[int]

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		hashSource[int]{1, 2, 2},
		hashSource[int]{2, 3, 3},
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
		hashSource[int]{1},
		NewHashSetFromSlice([]int{1, 2}),
	))
	assert.Equal(t, []int{2}, slices.Collect(set.All()))
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		hashSource[int]{1},
		hashSource[int]{2},
		hashSource[int]{1, 3, 3},
	))
	assert.Equal(t, []int{3}, slices.Collect(set.All()))

	set.Clear()
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		NewHashSetFromSlice([]int{1}),
		hashSource[int]{2},
	))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}
