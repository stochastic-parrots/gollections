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

func TestNewHashSet(t *testing.T) {
	assertSetBehavior(t, set.NewHashSet[int](2))
}

func TestHashSet_From(t *testing.T) {
	data := []int{1, 2, 1}
	values := set.HashSetFrom(data)
	data[0] = 3

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestHashSet_FromSeq(t *testing.T) {
	values := set.HashSetFromSeq(slices.Values([]int{1, 2, 1}))

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

func TestHashSet_Union(t *testing.T) {
	left := set.HashSetFrom([]int{1, 1, 2})
	result := left.Union(hashValueSource[int]{2, 3})

	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(result.All()))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(left.All()))
}

func TestHashSet_Intersection(t *testing.T) {
	left := set.HashSetFrom([]int{1, 2, 2, 3})
	result := left.Intersection(
		hashValueSource[int]{2, 3, 4},
		hashValueSource[int]{3, 3},
	)

	assert.Equal(t, []int{3}, slices.Collect(result.All()))
}

func TestHashSet_Difference(t *testing.T) {
	left := set.HashSetFrom([]int{1, 2, 2, 3})
	result := left.Difference(
		hashValueSource[int]{2},
		hashValueSource[int]{4},
	)

	assert.ElementsMatch(t, []int{1, 3}, slices.Collect(result.All()))
}

func TestHashSet_SymmetricDifference(t *testing.T) {
	left := set.HashSetFrom([]int{1, 1, 2})
	result := left.SymmetricDifference(
		hashValueSource[int]{2, 3},
		hashValueSource[int]{3, 4},
	)

	assert.ElementsMatch(t, []int{1, 4}, slices.Collect(result.All()))
}

func readonlyHashSet(values ...int) set.Readonly[int] {
	return set.AsReadonly(set.HashSetFrom(values))
}

func TestRelations_Equal(t *testing.T) {
	left := set.HashSetFrom([]int{1, 2})
	right := readonlyHashSet(2, 1)

	assert.True(t, set.Equal(left, right))
	assert.False(t, set.Equal(left, readonlyHashSet(1, 3)))
}

func TestRelations_IsSubset(t *testing.T) {
	assert.True(t, set.IsSubset(readonlyHashSet(1), readonlyHashSet(1, 2)))
	assert.False(t, set.IsSubset(readonlyHashSet(3), readonlyHashSet(1, 2)))
	assert.True(t, set.IsSubset(readonlyHashSet(), readonlyHashSet(1)))
}

func TestRelations_IsProperSubset(t *testing.T) {
	assert.True(t, set.IsProperSubset(readonlyHashSet(1), readonlyHashSet(1, 2)))
	assert.False(t, set.IsProperSubset(readonlyHashSet(1), readonlyHashSet(1)))
	assert.False(t, set.IsProperSubset(readonlyHashSet(), readonlyHashSet()))
}

func TestRelations_IsSuperset(t *testing.T) {
	assert.True(t, set.IsSuperset(readonlyHashSet(1, 2), readonlyHashSet(1)))
	assert.False(t, set.IsSuperset(readonlyHashSet(1, 2), readonlyHashSet(3)))
}

func TestRelations_IsProperSuperset(t *testing.T) {
	assert.True(t, set.IsProperSuperset(readonlyHashSet(1, 2), readonlyHashSet(1)))
	assert.False(t, set.IsProperSuperset(readonlyHashSet(1), readonlyHashSet(1)))
}

func TestRelations_IsDisjoint(t *testing.T) {
	assert.True(t, set.IsDisjoint(readonlyHashSet(1), readonlyHashSet(2, 3)))
	assert.False(t, set.IsDisjoint(readonlyHashSet(1, 2), readonlyHashSet(2)))
}
