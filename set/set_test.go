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
	var _ *set.HashSet[int] = set.HashSetOf[int]().New(0)
	var _ *set.HashSet[int] = set.HashSetOf[int]().From([]int{1})
	var _ *set.HashSet[int] = set.HashSetOf[int]().FromSeq(slices.Values([]int{1}))
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).New(0)
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).From([]member{{ID: 1}})
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).FromSeq(slices.Values([]member{{ID: 1}}))
	var _ set.Set[int] = set.HashSetOf[int]().New(0)
	var _ set.Set[member] = set.HashSetBy(memberID).New(0)
	var _ set.Algebra[int, *set.HashSet[int]] = set.HashSetOf[int]().New(0)
	var _ set.Algebra[member, *set.KeyedHashSet[member, int]] = set.HashSetBy(memberID).New(0)
	var _ set.InPlaceAlgebra[int] = set.HashSetOf[int]().New(0)
	var _ set.InPlaceAlgebra[member] = set.HashSetBy(memberID).New(0)
}

func TestConcreteZeroValues(t *testing.T) {
	var values set.HashSet[int]

	assert.True(t, values.Add(1))
	assert.Equal(t, 1, values.Adds(2, 1))

	assert.Equal(t, 2, values.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, set.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := set.HashSetOf[int]().From([]int{1, 2})
		view := set.AsReadonly[int](mutable)

		_, mutableView := any(view).(set.Set[int])
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, mutableView)
		assert.False(t, unmarshals)
		assert.True(t, view.Contains(1))
		assert.False(t, view.IsEmpty())
		assert.Equal(t, 2, view.Length())
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

func TestInPlaceAlgebra_ReadonlyAlias(t *testing.T) {
	t.Run("HashSet", func(t *testing.T) {
		values := set.HashSetOf[int]().From([]int{1, 2})
		view := set.AsReadonly[int](values)

		assert.Zero(t, values.UnionWith(view))
		assert.Zero(t, values.IntersectWith(view))
		assert.Zero(t, values.SymmetricDifferenceWith(view, view))
		assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))

		assert.Equal(t, 2, values.SymmetricDifferenceWith(view))
		assert.True(t, values.IsEmpty())

		values.Adds(1, 2)
		assert.Equal(t, 2, values.DifferenceWith(view))
		assert.True(t, values.IsEmpty())
	})

	t.Run("KeyedHashSet", func(t *testing.T) {
		values := set.HashSetBy(memberID).From([]member{{ID: 1}, {ID: 2}})
		view := set.AsReadonly[member](values)

		assert.Zero(t, values.UnionWith(view))
		assert.Zero(t, values.IntersectWith(view))
		assert.Zero(t, values.SymmetricDifferenceWith(view, view))
		assert.ElementsMatch(t, []member{{ID: 1}, {ID: 2}}, slices.Collect(values.All()))

		assert.Equal(t, 2, values.SymmetricDifferenceWith(view))
		assert.True(t, values.IsEmpty())

		values.Adds(member{ID: 1}, member{ID: 2})
		assert.Equal(t, 2, values.DifferenceWith(view))
		assert.True(t, values.IsEmpty())
	})
}
