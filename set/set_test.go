package set_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

type member struct {
	ID     int
	Name   string
	Labels []string
}

func memberID(value member) int {
	return value.ID
}

func TestFactoriesImplementSet(t *testing.T) {
	var _ *set.HashSet[int] = set.NewHashSet[int](0)
	var _ *set.HashSet[int] = set.HashSetFrom[int]([]int{1})
	var _ *set.HashSet[int] = set.HashSetFromSeq[int](slices.Values([]int{1}))
	var _ *set.KeyedHashSet[member, int] = set.NewKeyedHashSet(memberID, 0)
	var _ *set.KeyedHashSet[member, int] = set.KeyedHashSetFrom(memberID, []member{{ID: 1}})
	var _ *set.KeyedHashSet[member, int] = set.KeyedHashSetFromSeq(memberID, slices.Values([]member{{ID: 1}}))
	var _ set.Set[int] = set.NewHashSet[int](0)
	var _ set.Set[member] = set.NewKeyedHashSet(memberID, 0)
	var _ set.Source[int] = set.NewHashSet[int](0)
	var _ set.Source[member] = set.NewKeyedHashSet(memberID, 0)
	var _ set.Algebra[int, *set.HashSet[int]] = set.NewHashSet[int](0)
	var _ set.Algebra[member, *set.KeyedHashSet[member, int]] = set.NewKeyedHashSet(memberID, 0)
	var _ set.InPlaceAlgebra[int] = set.NewHashSet[int](0)
	var _ set.InPlaceAlgebra[member] = set.NewKeyedHashSet(memberID, 0)
}

func TestConcreteZeroValues(t *testing.T) {
	var values set.HashSet[int]

	assert.Equal(t, 1, values.Add(1))
	assert.Equal(t, 1, values.Add(2, 1))

	assert.Equal(t, 2, values.Len())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() set.Set[int]{
		"HashSetFromNil":        func() set.Set[int] { return set.HashSetFrom([]int(nil)) },
		"HashSetFromEmpty":      func() set.Set[int] { return set.HashSetFrom(make([]int, 0, 4)) },
		"KeyedHashSetFromNil":   func() set.Set[int] { return set.KeyedHashSetFrom(func(x int) int { return x }, []int(nil)) },
		"KeyedHashSetFromEmpty": func() set.Set[int] { return set.KeyedHashSetFrom(func(x int) int { return x }, make([]int, 0, 4)) },
	} {
		t.Run(name, func(t *testing.T) {
			values := construct()
			assert.True(t, values.IsEmpty())
			assert.Equal(t, 1, values.Add(1))
			assert.True(t, values.Contains(1))
		})
	}
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, set.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := set.HashSetFrom[int]([]int{1, 2})
		view := set.AsReadonly[int](mutable)

		_, mutableView := any(view).(set.Set[int])
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, mutableView)
		assert.False(t, unmarshals)
		assert.True(t, view.Contains(1))
		assert.False(t, view.IsEmpty())
		assert.Equal(t, 2, view.Len())
		assert.ElementsMatch(t, []int{1, 2}, slices.Collect(view.All()))

		var indexes []int
		for idx := range view.Enumerate() {
			indexes = append(indexes, idx)
		}
		assert.Equal(t, []int{0, 1}, indexes)

		data, err := view.MarshalJSON()
		assert.NoError(t, err)
		var decoded []int
		assert.NoError(t, json.Unmarshal(data, &decoded))
		assert.ElementsMatch(t, []int{1, 2}, decoded)

		mutable.Add(3)
		assert.True(t, view.Contains(3))
	})
}

func TestSourceFromSlice(t *testing.T) {
	values := []int{1, 2}
	var source *set.SliceSource[int] = set.SourceFromSlice(values)
	values[1] = 3

	assert.Equal(t, 2, source.Len())
	assert.Equal(t, []int{1, 3}, slices.Collect(source.All()))
}

func TestSliceSource_Len(t *testing.T) {
	var empty set.SliceSource[int]
	assert.Zero(t, empty.Len())
	assert.Equal(t, 3, set.SourceFromSlice([]int{1, 2, 2}).Len())
}

func TestSliceSource_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		source := set.SourceFromSlice([]int{1, 2, 2})
		assert.Equal(t, []int{1, 2, 2}, slices.Collect(source.All()))
	})
	t.Run("PartialIteration", func(t *testing.T) {
		source := set.SourceFromSlice([]int{1, 2, 2})
		count := 0
		for value := range source.All() {
			assert.Equal(t, 1, value)
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
	t.Run("EmptySource", func(t *testing.T) {
		var source set.SliceSource[int]
		assert.Empty(t, slices.Collect(source.All()))
	})
}

func TestConcreteZeroValue_SymmetricDifferenceWith(t *testing.T) {
	var values set.HashSet[int]
	other := set.HashSetFrom([]int{1})

	assert.Equal(t, 1, values.SymmetricDifferenceWith(other))
	assert.True(t, values.Contains(1))
}

func TestInPlaceAlgebra_ReadonlyAlias(t *testing.T) {
	t.Run("HashSet", func(t *testing.T) {
		values := set.HashSetFrom([]int{1, 2})
		view := set.AsReadonly[int](values)

		assert.Zero(t, values.UnionWith(view))
		assert.Zero(t, values.IntersectWith(view))
		assert.Zero(t, values.SymmetricDifferenceWith(view, view))
		assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))

		assert.Equal(t, 2, values.SymmetricDifferenceWith(view))
		assert.True(t, values.IsEmpty())

		values.Add(1, 2)
		assert.Equal(t, 2, values.DifferenceWith(view))
		assert.True(t, values.IsEmpty())
	})

	t.Run("KeyedHashSet", func(t *testing.T) {
		values := set.KeyedHashSetFrom(memberID, []member{{ID: 1}, {ID: 2}})
		view := set.AsReadonly[member](values)

		assert.Zero(t, values.UnionWith(view))
		assert.Zero(t, values.IntersectWith(view))
		assert.Zero(t, values.SymmetricDifferenceWith(view, view))
		assert.ElementsMatch(t, []member{{ID: 1}, {ID: 2}}, slices.Collect(values.All()))

		assert.Equal(t, 2, values.SymmetricDifferenceWith(view))
		assert.True(t, values.IsEmpty())

		values.Add(member{ID: 1}, member{ID: 2})
		assert.Equal(t, 2, values.DifferenceWith(view))
		assert.True(t, values.IsEmpty())
	})
}
