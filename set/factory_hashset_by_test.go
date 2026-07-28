package set_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

type keyedValueSource[T any] []T

func (values keyedValueSource[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func (values keyedValueSource[T]) Length() int {
	return len(values)
}

func TestHashSetBy_NilIdentityFunction(t *testing.T) {
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.HashSetBy[int, int](nil)
	})
}

func TestHashSetByFactory_New(t *testing.T) {
	values := set.HashSetBy(memberID).New(2)
	first := member{ID: 1, Name: "first", Labels: []string{"a"}}

	assert.True(t, values.Add(first))
	assert.False(t, values.Add(member{ID: 1, Name: "duplicate"}))

	assert.Equal(t, 1, values.Length())
	assert.Equal(t, []member{first}, slices.Collect(values.All()))
}

func TestHashSetByFactory_From(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	data := []member{first, {ID: 1, Name: "duplicate"}, {ID: 2, Name: "second"}}
	values := set.HashSetBy(memberID).From(data)
	data[0] = member{ID: 3, Name: "changed"}

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestHashSetByFactory_FromSeq(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	values := set.HashSetBy(memberID).FromSeq(slices.Values([]member{
		first,
		{ID: 1, Name: "duplicate"},
		{ID: 2, Name: "second"},
	}))

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestHashSetByFactory_Union(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}

	result := factory.Union(
		keyedValueSource[member]{first, {ID: 1, Name: "duplicate"}},
		keyedValueSource[member]{{ID: 1, Name: "later"}, {ID: 2, Name: "second"}},
	)

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(result.All()))
	assert.True(t, factory.Union().IsEmpty())
}

func TestHashSetByFactory_Intersection(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 2, Name: "from first source"}

	result := factory.Intersection(
		keyedValueSource[member]{{ID: 1}, first, {ID: 2, Name: "duplicate"}},
		keyedValueSource[member]{{ID: 2, Name: "other"}, {ID: 3}},
	)

	assert.Equal(t, []member{first}, slices.Collect(result.All()))
	assert.True(t, factory.Intersection().IsEmpty())
}

func TestHashSetByFactory_Difference(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}

	result := factory.Difference(
		keyedValueSource[member]{first, {ID: 2}, {ID: 2, Name: "duplicate"}},
		keyedValueSource[member]{{ID: 2, Name: "other"}},
	)

	assert.Equal(t, []member{first}, slices.Collect(result.All()))
	assert.True(t, factory.Difference().IsEmpty())
}

func TestHashSetByFactory_SymmetricDifference(t *testing.T) {
	factory := set.HashSetBy(memberID)
	first := member{ID: 1, Name: "first"}
	last := member{ID: 4, Name: "last odd representative"}

	result := factory.SymmetricDifference(
		keyedValueSource[member]{first, {ID: 1, Name: "duplicate"}, {ID: 2}},
		keyedValueSource[member]{{ID: 2}, {ID: 3}},
		keyedValueSource[member]{{ID: 3}, last},
	)

	assert.ElementsMatch(t, []member{first, last}, slices.Collect(result.All()))
	assert.True(t, factory.SymmetricDifference().IsEmpty())
}

func TestHashSetByFactory_Equal(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.Equal(
		keyedValueSource[member]{{ID: 1, Name: "left"}},
		keyedValueSource[member]{{ID: 1, Name: "right"}, {ID: 1, Name: "duplicate"}},
	))
	assert.False(t, factory.Equal(
		keyedValueSource[member]{{ID: 1}},
		keyedValueSource[member]{{ID: 2}},
	))
}

func TestHashSetByFactory_IsSubset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsSubset(
		keyedValueSource[member]{{ID: 1}, {ID: 1}},
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
	))
	assert.False(t, factory.IsSubset(
		keyedValueSource[member]{{ID: 3}},
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
	))
}

func TestHashSetByFactory_IsProperSubset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsProperSubset(
		keyedValueSource[member]{{ID: 1}},
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
	))
	assert.False(t, factory.IsProperSubset(
		keyedValueSource[member]{{ID: 1}, {ID: 1}},
		keyedValueSource[member]{{ID: 1}},
	))
}

func TestHashSetByFactory_IsSuperset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsSuperset(
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
		keyedValueSource[member]{{ID: 1}, {ID: 1}},
	))
	assert.False(t, factory.IsSuperset(
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
		keyedValueSource[member]{{ID: 3}},
	))
}

func TestHashSetByFactory_IsProperSuperset(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsProperSuperset(
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
		keyedValueSource[member]{{ID: 1}},
	))
	assert.False(t, factory.IsProperSuperset(
		keyedValueSource[member]{{ID: 1}},
		keyedValueSource[member]{{ID: 1}, {ID: 1}},
	))
}

func TestHashSetByFactory_IsDisjoint(t *testing.T) {
	factory := set.HashSetBy(memberID)

	assert.True(t, factory.IsDisjoint(
		keyedValueSource[member]{{ID: 1}},
		keyedValueSource[member]{{ID: 2}, {ID: 3}},
	))
	assert.False(t, factory.IsDisjoint(
		keyedValueSource[member]{{ID: 1}, {ID: 2}},
		keyedValueSource[member]{{ID: 2}},
	))
}

func TestHashSetByFactory_RelationOperandStrategies(t *testing.T) {
	factory := set.HashSetBy(memberID)
	one := member{ID: 1}
	two := member{ID: 2}
	three := member{ID: 3}

	assert.True(t, factory.Equal(
		keyedValueSource[member]{one, one, two},
		keyedValueSource[member]{two, one},
	))

	assert.True(t, factory.IsSubset(
		keyedValueSource[member]{one, one, one},
		keyedValueSource[member]{one, two},
	))
	assert.False(t, factory.IsSubset(
		keyedValueSource[member]{one, three, three},
		keyedValueSource[member]{one, two},
	))
	assert.True(t, factory.IsSubset(
		keyedValueSource[member]{},
		keyedValueSource[member]{one},
	))

	assert.False(t, factory.IsProperSubset(
		keyedValueSource[member]{},
		keyedValueSource[member]{},
	))

	assert.True(t, factory.IsSuperset(
		keyedValueSource[member]{one, two},
		keyedValueSource[member]{one, one, one},
	))
	assert.False(t, factory.IsSuperset(
		keyedValueSource[member]{one, two},
		keyedValueSource[member]{one, three, three},
	))
}
