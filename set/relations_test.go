package set_test

import (
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

func TestRelationsBy(t *testing.T) {
	byName := func(value member) string { return value.Name }
	byID := func(value member) int { return value.ID }

	left := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "left one"},
		{ID: 1, Name: "left duplicate by ID"},
	}, byName)
	right := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "right one"},
		{ID: 2, Name: "right two"},
	}, byName)
	equal := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "other one"},
		{ID: 1, Name: "other duplicate by ID"},
	}, byName)

	assert.True(t, set.EqualBy(byID, left, equal))
	assert.True(t, set.EqualBy(byID, left, equal, set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "third one"},
		{ID: 1, Name: "third duplicate"},
	}, byName)))
	assert.False(t, set.EqualBy(byID, left, equal, right))
	assert.True(t, set.IsSubsetBy(left, right, byID))
	assert.True(t, set.IsProperSubsetBy(left, right, byID))
	assert.True(t, set.IsSupersetBy(right, left, byID))
	assert.True(t, set.IsProperSupersetBy(right, left, byID))
	assert.False(t, set.IsDisjointBy(byID, left, right))
	assert.True(t, set.IsDisjointBy(byID, left, set.KeyedHashSetFrom([]member{
		{ID: 3, Name: "unrelated"},
	}, byName)))
	assert.False(t, set.IsDisjointBy(
		byID,
		set.KeyedHashSetFrom([]member{{ID: 1, Name: "first"}}, byName),
		set.KeyedHashSetFrom([]member{{ID: 2, Name: "second"}}, byName),
		set.KeyedHashSetFrom([]member{{ID: 2, Name: "overlap with second"}}, byName),
	))
	assert.False(t, set.IsSubsetBy(set.KeyedHashSetFrom([]member{
		{ID: 3, Name: "missing"},
	}, byName), right, byID))
}

func TestRelationsBy_ReadonlyOperands(t *testing.T) {
	byName := func(value member) string { return value.Name }
	byID := func(value member) int { return value.ID }
	asReadonly := func(values set.Set[member]) set.Readonly[member] {
		return set.AsReadonly[member](values)
	}

	subset := set.KeyedHashSetFrom([]member{{ID: 1, Name: "a"}}, byName)
	superset := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "b"},
		{ID: 2, Name: "c"},
	}, byName)
	assert.True(t, set.EqualBy(byID, asReadonly(subset), asReadonly(set.KeyedHashSetFrom([]member{{ID: 1, Name: "different"}}, byName))))
	assert.True(t, set.IsSubsetBy(asReadonly(subset), asReadonly(superset), byID))
	assert.True(t, set.IsProperSubsetBy(asReadonly(subset), asReadonly(superset), byID))
	assert.True(t, set.IsSupersetBy(asReadonly(superset), asReadonly(subset), byID))
	assert.True(t, set.IsProperSupersetBy(asReadonly(superset), asReadonly(subset), byID))
	assert.True(t, set.IsDisjointBy(
		byID,
		asReadonly(subset),
		asReadonly(set.KeyedHashSetFrom([]member{{ID: 3, Name: "unrelated"}}, byName)),
	))

	longSubset := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "one"},
		{ID: 1, Name: "one duplicate identity"},
		{ID: 2, Name: "two"},
		{ID: 2, Name: "two duplicate identity"},
	}, byName)
	shortSuperset := set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "other one"},
		{ID: 2, Name: "other two"},
		{ID: 3, Name: "three"},
	}, byName)
	assert.True(t, set.IsSubsetBy(asReadonly(longSubset), asReadonly(shortSuperset), byID))
	assert.True(t, set.IsProperSubsetBy(asReadonly(longSubset), asReadonly(shortSuperset), byID))
	assert.True(t, set.EqualBy(
		byID,
		asReadonly(longSubset),
		asReadonly(set.KeyedHashSetFrom([]member{
			{ID: 1, Name: "same first ID"},
			{ID: 2, Name: "same second ID"},
		}, byName)),
	))
	assert.False(t, set.IsProperSubsetBy(asReadonly(longSubset), asReadonly(superset), byID))
	assert.False(t, set.IsSubsetBy(asReadonly(shortSuperset), asReadonly(longSubset), byID))
	assert.False(t, set.IsSubsetBy(asReadonly(longSubset), asReadonly(set.KeyedHashSetFrom([]member{
		{ID: 1, Name: "one"},
	}, byName)), byID))
}

func TestRelations_ReadonlyOperands(t *testing.T) {
	left := set.HashSetFrom([]int{1, 2})
	right := set.HashSetFrom([]int{2, 1})
	leftView := readonlyHashSet(1, 2)
	rightView := readonlyHashSet(2, 1)

	assert.True(t, set.Equal(left, rightView))
	assert.True(t, set.Equal(leftView, right))
	assert.True(t, set.Equal(leftView, rightView))
	assert.True(t, set.Equal(left, rightView, readonlyHashSet(1, 2)))
	assert.False(t, set.Equal(left, rightView, readonlyHashSet(1, 3)))
	assert.False(t, set.Equal(readonlyHashSet(1), rightView))

	assert.True(t, set.IsSubset(left, rightView))
	assert.True(t, set.IsSubset(readonlyHashSet(1), right))
	assert.False(t, set.IsSubset(readonlyHashSet(3), right))
	assert.True(t, set.IsSubset(readonlyHashSet(1), rightView))
	assert.False(t, set.IsSubset(readonlyHashSet(1, 3, 4), rightView))

	assert.True(t, set.IsProperSubset(left, readonlyHashSet(1, 2, 3)))
	assert.True(t, set.IsProperSubset(readonlyHashSet(1), rightView))
	assert.False(t, set.IsProperSubset(readonlyHashSet(1, 2), rightView))
	assert.False(t, set.IsProperSubset(readonlyHashSet(1, 2, 3), rightView))

	assert.True(t, set.IsSuperset(right, leftView))
	assert.True(t, set.IsSuperset(rightView, readonlyHashSet(1)))
	assert.False(t, set.IsSuperset(readonlyHashSet(1), rightView))

	assert.True(t, set.IsProperSuperset(right, readonlyHashSet(1)))
	assert.True(t, set.IsProperSuperset(rightView, readonlyHashSet(1)))
	assert.False(t, set.IsProperSuperset(readonlyHashSet(1), rightView))

	assert.False(t, set.IsDisjoint(left, rightView))
	assert.False(t, set.IsDisjoint(leftView, right))
	assert.True(t, set.IsDisjoint(readonlyHashSet(1), readonlyHashSet(2, 3)))
	assert.True(t, set.IsDisjoint(readonlyHashSet(1), readonlyHashSet(2), readonlyHashSet(3)))
	assert.False(t, set.IsDisjoint(readonlyHashSet(1), readonlyHashSet(2), readonlyHashSet(2, 3)))
}

func TestEqualBy_NilIdentityFunction(t *testing.T) {
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.EqualBy[member, int](nil, nil, nil)
	})
}
