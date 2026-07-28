package set_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

type hashValueSource[T any] []T

func (values hashValueSource[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func (values hashValueSource[T]) Length() int {
	return len(values)
}

func TestHashFactory_New(t *testing.T) {
	assertSetBehavior(t, set.HashSetOf[int]().New(2))
}

func TestHashFactory_From(t *testing.T) {
	data := []int{1, 2, 1}
	values := set.HashSetOf[int]().From(data)
	data[0] = 3

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestHashFactory_FromSeq(t *testing.T) {
	values := set.HashSetOf[int]().FromSeq(slices.Values([]int{1, 2, 1}))

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func assertSetBehavior(t *testing.T, values set.Set[int]) {
	t.Helper()

	assert.True(t, values.IsEmpty())
	assert.True(t, values.Add(1))
	assert.Equal(t, 1, values.Adds(2, 1))

	assert.False(t, values.IsEmpty())
	assert.Equal(t, 2, values.Length())
	assert.True(t, values.Contains(1))
	assert.False(t, values.Contains(3))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))

	assert.True(t, values.Remove(1))
	assert.False(t, values.Remove(1))
	assert.Equal(t, 2, values.Adds(3, 4))
	assert.Equal(t, 2, values.Removes(2, 2, 4, 5))

	data, err := values.MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `[3]`, string(data))

	values.Clear()
	assert.True(t, values.IsEmpty())

	assert.True(t, values.Add(3))
	assert.True(t, values.Contains(3))
}

func TestHashFactory_Union(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Union(
		hashValueSource[int]{1, 1, 2},
		hashValueSource[int]{2, 3},
	)

	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(result.All()))
	assert.True(t, factory.Union().IsEmpty())
}

func TestHashFactory_Intersection(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Intersection(
		hashValueSource[int]{1, 2, 2, 3},
		hashValueSource[int]{2, 3, 4},
		hashValueSource[int]{3, 3},
	)

	assert.Equal(t, []int{3}, slices.Collect(result.All()))
	assert.True(t, factory.Intersection().IsEmpty())
}

func TestHashFactory_Difference(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.Difference(
		hashValueSource[int]{1, 2, 2, 3},
		hashValueSource[int]{2},
		hashValueSource[int]{4},
	)

	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(result.All()))
	assert.True(t, factory.Difference().IsEmpty())
}

func TestHashFactory_SymmetricDifference(t *testing.T) {
	factory := set.HashSetOf[int]()

	result := factory.SymmetricDifference(
		hashValueSource[int]{1, 1, 2},
		hashValueSource[int]{2, 3},
		hashValueSource[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
	assert.True(t, factory.SymmetricDifference().IsEmpty())
}

func TestHashFactory_Equal(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.Equal(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{2, 1, 1},
	))
	assert.False(t, factory.Equal(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{1, 3},
	))
}

func TestHashFactory_IsSubset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsSubset(
		hashValueSource[int]{1, 1},
		hashValueSource[int]{1, 2},
	))
	assert.False(t, factory.IsSubset(
		hashValueSource[int]{1, 3},
		hashValueSource[int]{1, 2},
	))
}

func TestHashFactory_IsProperSubset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsProperSubset(
		hashValueSource[int]{1},
		hashValueSource[int]{1, 2},
	))
	assert.False(t, factory.IsProperSubset(
		hashValueSource[int]{1, 1},
		hashValueSource[int]{1},
	))
}

func TestHashFactory_IsSuperset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsSuperset(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{1, 1},
	))
	assert.False(t, factory.IsSuperset(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{3},
	))
}

func TestHashFactory_IsProperSuperset(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsProperSuperset(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{1},
	))
	assert.False(t, factory.IsProperSuperset(
		hashValueSource[int]{1},
		hashValueSource[int]{1, 1},
	))
}

func TestHashFactory_IsDisjoint(t *testing.T) {
	factory := set.HashSetOf[int]()

	assert.True(t, factory.IsDisjoint(
		hashValueSource[int]{1},
		hashValueSource[int]{2, 3},
	))
	assert.False(t, factory.IsDisjoint(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{2},
	))
}

func TestHashFactory_RelationOperandStrategies(t *testing.T) {
	factory := set.HashSetOf[int]()
	oneTwo := factory.From([]int{1, 2})
	twoThree := factory.From([]int{2, 3})

	assert.True(t, factory.Equal(oneTwo, hashValueSource[int]{2, 1}))
	assert.True(t, factory.Equal(hashValueSource[int]{2, 1}, oneTwo))
	assert.True(t, factory.Equal(
		hashValueSource[int]{1, 1, 2},
		hashValueSource[int]{2, 1},
	))

	assert.True(t, factory.IsSubset(oneTwo, hashValueSource[int]{1, 2, 3}))
	assert.True(t, factory.IsSubset(hashValueSource[int]{1}, oneTwo))
	assert.False(t, factory.IsSubset(hashValueSource[int]{3}, oneTwo))
	assert.True(t, factory.IsSubset(
		hashValueSource[int]{1, 1, 1},
		hashValueSource[int]{1, 2},
	))
	assert.False(t, factory.IsSubset(
		hashValueSource[int]{1, 3, 3},
		hashValueSource[int]{1, 2},
	))
	assert.True(t, factory.IsSubset(
		hashValueSource[int]{},
		hashValueSource[int]{1},
	))

	assert.True(t, factory.IsProperSubset(
		factory.From([]int{1}),
		hashValueSource[int]{1, 2},
	))
	assert.False(t, factory.IsProperSubset(
		hashValueSource[int]{},
		hashValueSource[int]{},
	))

	assert.True(t, factory.IsSuperset(oneTwo, hashValueSource[int]{1}))
	assert.True(t, factory.IsSuperset(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{1, 1, 1},
	))
	assert.False(t, factory.IsSuperset(
		hashValueSource[int]{1, 2},
		hashValueSource[int]{1, 3, 3},
	))
	assert.True(t, factory.IsProperSuperset(
		oneTwo,
		hashValueSource[int]{1},
	))

	assert.False(t, factory.IsDisjoint(oneTwo, twoThree))
	assert.False(t, factory.IsDisjoint(hashValueSource[int]{2}, oneTwo))
}
