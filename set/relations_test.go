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

	assert.True(t, set.EqualBy(left, equal, byID))
	assert.True(t, set.IsSubsetBy(left, right, byID))
	assert.True(t, set.IsProperSubsetBy(left, right, byID))
	assert.True(t, set.IsSupersetBy(right, left, byID))
	assert.True(t, set.IsProperSupersetBy(right, left, byID))
	assert.False(t, set.IsDisjointBy(left, right, byID))
	assert.True(t, set.IsDisjointBy(left, set.KeyedHashSetFrom([]member{
		{ID: 3, Name: "unrelated"},
	}, byName), byID))
	assert.False(t, set.IsSubsetBy(set.KeyedHashSetFrom([]member{
		{ID: 3, Name: "missing"},
	}, byName), right, byID))
}

func TestEqualBy_NilIdentityFunction(t *testing.T) {
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.EqualBy[member, int](nil, nil, nil)
	})
}
