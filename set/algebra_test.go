package set_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

type valuesCollection[T any] []T

func (values valuesCollection[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func (values valuesCollection[T]) Enumerate() iter.Seq2[int, T] {
	return slices.All(values)
}

func (values valuesCollection[T]) IsEmpty() bool {
	return len(values) == 0
}

func (values valuesCollection[T]) Length() int {
	return len(values)
}

func TestHashFactory_Union(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Union(
		valuesCollection[int]{1, 1, 2},
		valuesCollection[int]{2, 3},
	)

	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(result.All()))
	assert.True(t, factory.Union().IsEmpty())
}

func TestHashFactory_Intersection(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Intersection(
		valuesCollection[int]{1, 2, 2, 3},
		valuesCollection[int]{2, 3, 4},
		valuesCollection[int]{3, 3},
	)

	assert.Equal(t, []int{3}, slices.Collect(result.All()))
	assert.True(t, factory.Intersection().IsEmpty())
}

func TestHashFactory_Difference(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Difference(
		valuesCollection[int]{1, 2, 2, 3},
		valuesCollection[int]{2},
		valuesCollection[int]{4},
	)

	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(result.All()))
	assert.True(t, factory.Difference().IsEmpty())
}

func TestHashFactory_SymmetricDifference(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.SymmetricDifference(
		valuesCollection[int]{1, 1, 2},
		valuesCollection[int]{2, 3},
		valuesCollection[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
	assert.True(t, factory.SymmetricDifference().IsEmpty())
}

func TestHashFactory_Equal(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.Equal(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{2, 1, 1},
	))
	assert.False(t, factory.Equal(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{1, 3},
	))
}

func TestHashFactory_IsSubset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsSubset(
		valuesCollection[int]{1, 1},
		valuesCollection[int]{1, 2},
	))
	assert.False(t, factory.IsSubset(
		valuesCollection[int]{1, 3},
		valuesCollection[int]{1, 2},
	))
}

func TestHashFactory_IsProperSubset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsProperSubset(
		valuesCollection[int]{1},
		valuesCollection[int]{1, 2},
	))
	assert.False(t, factory.IsProperSubset(
		valuesCollection[int]{1, 1},
		valuesCollection[int]{1},
	))
}

func TestHashFactory_IsSuperset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsSuperset(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{1, 1},
	))
	assert.False(t, factory.IsSuperset(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{3},
	))
}

func TestHashFactory_IsProperSuperset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsProperSuperset(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{1},
	))
	assert.False(t, factory.IsProperSuperset(
		valuesCollection[int]{1},
		valuesCollection[int]{1, 1},
	))
}

func TestHashFactory_IsDisjoint(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsDisjoint(
		valuesCollection[int]{1},
		valuesCollection[int]{2, 3},
	))
	assert.False(t, factory.IsDisjoint(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{2},
	))
}

func TestHashFactory_RelationOperandStrategies(t *testing.T) {
	factory := set.HashSetOf[int]()
	oneTwo := factory.From([]int{1, 2})
	twoThree := factory.From([]int{2, 3})

	assert.True(t, factory.Equal(oneTwo, valuesCollection[int]{2, 1}))
	assert.True(t, factory.Equal(valuesCollection[int]{2, 1}, oneTwo))
	assert.True(t, factory.Equal(
		valuesCollection[int]{1, 1, 2},
		valuesCollection[int]{2, 1},
	))

	assert.True(t, factory.IsSubset(oneTwo, valuesCollection[int]{1, 2, 3}))
	assert.True(t, factory.IsSubset(valuesCollection[int]{1}, oneTwo))
	assert.False(t, factory.IsSubset(valuesCollection[int]{3}, oneTwo))
	assert.True(t, factory.IsSubset(
		valuesCollection[int]{1, 1, 1},
		valuesCollection[int]{1, 2},
	))
	assert.False(t, factory.IsSubset(
		valuesCollection[int]{1, 3, 3},
		valuesCollection[int]{1, 2},
	))
	assert.True(t, factory.IsSubset(
		valuesCollection[int]{},
		valuesCollection[int]{1},
	))

	assert.True(t, factory.IsProperSubset(
		factory.From([]int{1}),
		valuesCollection[int]{1, 2},
	))
	assert.False(t, factory.IsProperSubset(
		valuesCollection[int]{},
		valuesCollection[int]{},
	))

	assert.True(t, factory.IsSuperset(oneTwo, valuesCollection[int]{1}))
	assert.True(t, factory.IsSuperset(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{1, 1, 1},
	))
	assert.False(t, factory.IsSuperset(
		valuesCollection[int]{1, 2},
		valuesCollection[int]{1, 3, 3},
	))
	assert.True(t, factory.IsProperSuperset(
		oneTwo,
		valuesCollection[int]{1},
	))

	assert.False(t, factory.IsDisjoint(oneTwo, twoThree))
	assert.False(t, factory.IsDisjoint(valuesCollection[int]{2}, oneTwo))
}

func TestHashSetByFactory_Union(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}

	result := factory.Union(
		valuesCollection[member]{first, {ID: 1, Name: "duplicate"}},
		valuesCollection[member]{{ID: 1, Name: "later"}, {ID: 2, Name: "second"}},
	)

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(result.All()))
	assert.True(t, factory.Union().IsEmpty())
}

func TestHashSetByFactory_Intersection(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 2, Name: "from first source"}

	result := factory.Intersection(
		valuesCollection[member]{{ID: 1}, first, {ID: 2, Name: "duplicate"}},
		valuesCollection[member]{{ID: 2, Name: "other"}, {ID: 3}},
	)

	assert.Equal(t, []member{first}, slices.Collect(result.All()))
	assert.True(t, factory.Intersection().IsEmpty())
}

func TestHashSetByFactory_Difference(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}

	result := factory.Difference(
		valuesCollection[member]{first, {ID: 2}, {ID: 2, Name: "duplicate"}},
		valuesCollection[member]{{ID: 2, Name: "other"}},
	)

	assert.Equal(t, []member{first}, slices.Collect(result.All()))
	assert.True(t, factory.Difference().IsEmpty())
}

func TestHashSetByFactory_SymmetricDifference(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}
	last := member{ID: 4, Name: "last odd representative"}

	result := factory.SymmetricDifference(
		valuesCollection[member]{first, {ID: 1, Name: "duplicate"}, {ID: 2}},
		valuesCollection[member]{{ID: 2}, {ID: 3}},
		valuesCollection[member]{{ID: 3}, last},
	)

	assert.ElementsMatch(t, []member{first, last}, slices.Collect(result.All()))
	assert.True(t, factory.SymmetricDifference().IsEmpty())
}

func TestHashSetByFactory_Equal(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.Equal(
		valuesCollection[member]{{ID: 1, Name: "left"}},
		valuesCollection[member]{{ID: 1, Name: "right"}, {ID: 1, Name: "duplicate"}},
	))
	assert.False(t, factory.Equal(
		valuesCollection[member]{{ID: 1}},
		valuesCollection[member]{{ID: 2}},
	))
}

func TestHashSetByFactory_IsSubset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsSubset(
		valuesCollection[member]{{ID: 1}, {ID: 1}},
		valuesCollection[member]{{ID: 1}, {ID: 2}},
	))
	assert.False(t, factory.IsSubset(
		valuesCollection[member]{{ID: 3}},
		valuesCollection[member]{{ID: 1}, {ID: 2}},
	))
}

func TestHashSetByFactory_IsProperSubset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsProperSubset(
		valuesCollection[member]{{ID: 1}},
		valuesCollection[member]{{ID: 1}, {ID: 2}},
	))
	assert.False(t, factory.IsProperSubset(
		valuesCollection[member]{{ID: 1}, {ID: 1}},
		valuesCollection[member]{{ID: 1}},
	))
}

func TestHashSetByFactory_IsSuperset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsSuperset(
		valuesCollection[member]{{ID: 1}, {ID: 2}},
		valuesCollection[member]{{ID: 1}, {ID: 1}},
	))
	assert.False(t, factory.IsSuperset(
		valuesCollection[member]{{ID: 1}, {ID: 2}},
		valuesCollection[member]{{ID: 3}},
	))
}

func TestHashSetByFactory_IsProperSuperset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsProperSuperset(
		valuesCollection[member]{{ID: 1}, {ID: 2}},
		valuesCollection[member]{{ID: 1}},
	))
	assert.False(t, factory.IsProperSuperset(
		valuesCollection[member]{{ID: 1}},
		valuesCollection[member]{{ID: 1}, {ID: 1}},
	))
}

func TestHashSetByFactory_IsDisjoint(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsDisjoint(
		valuesCollection[member]{{ID: 1}},
		valuesCollection[member]{{ID: 2}, {ID: 3}},
	))
	assert.False(t, factory.IsDisjoint(
		valuesCollection[member]{{ID: 1}, {ID: 2}},
		valuesCollection[member]{{ID: 2}},
	))
}

func TestHashSetByFactory_RelationOperandStrategies(t *testing.T) {
	factory := set.HashSetBy(memberID)
	one := member{ID: 1}
	two := member{ID: 2}
	three := member{ID: 3}

	assert.True(t, factory.Equal(
		valuesCollection[member]{one, one, two},
		valuesCollection[member]{two, one},
	))

	assert.True(t, factory.IsSubset(
		valuesCollection[member]{one, one, one},
		valuesCollection[member]{one, two},
	))
	assert.False(t, factory.IsSubset(
		valuesCollection[member]{one, three, three},
		valuesCollection[member]{one, two},
	))
	assert.True(t, factory.IsSubset(
		valuesCollection[member]{},
		valuesCollection[member]{one},
	))

	assert.False(t, factory.IsProperSubset(
		valuesCollection[member]{},
		valuesCollection[member]{},
	))

	assert.True(t, factory.IsSuperset(
		valuesCollection[member]{one, two},
		valuesCollection[member]{one, one, one},
	))
	assert.False(t, factory.IsSuperset(
		valuesCollection[member]{one, two},
		valuesCollection[member]{one, three, three},
	))
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
