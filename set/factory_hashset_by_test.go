package set_test

import (
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyedHashSet_NilIdentityFunction(t *testing.T) {
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.NewKeyedHashSet[int, member](0, nil)
	})
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.KeyedHashSetFrom[int, member](nil, nil)
	})
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.KeyedHashSetFromSeq[int, member](nil, nil)
	})
}

func TestNewKeyedHashSet(t *testing.T) {
	values := set.NewKeyedHashSet(2, memberID)
	first := member{ID: 1, Name: "first", Labels: []string{"a"}}

	assert.True(t, values.Add(first))
	assert.False(t, values.Add(member{ID: 1, Name: "duplicate"}))

	assert.Equal(t, 1, values.Length())
	assert.Equal(t, []member{first}, slices.Collect(values.All()))
}

func TestKeyedHashSet_From(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	data := []member{first, {ID: 1, Name: "duplicate"}, {ID: 2, Name: "second"}}
	values := set.KeyedHashSetFrom(data, memberID)
	data[0] = member{ID: 3, Name: "changed"}

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestKeyedHashSet_FromSeq(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	values := set.KeyedHashSetFromSeq(slices.Values([]member{
		first,
		{ID: 1, Name: "duplicate"},
		{ID: 2, Name: "second"},
	}), memberID)

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestKeyedHashSet_FromMap(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	input := map[int]member{1: first}
	values := set.KeyedHashSetFromMap(input, memberID)

	assert.True(t, values.Add(member{ID: 2, Name: "set mutation"}))
	assert.Equal(t, member{ID: 2, Name: "set mutation"}, input[2])

	input[3] = member{ID: 3, Name: "map mutation"}
	assert.True(t, values.Contains(member{ID: 3}))
}

func TestKeyedHashSet_CloneMap(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	input := map[int]member{1: first}
	values := set.KeyedHashSetCloneMap(input, memberID)

	values.Add(member{ID: 2, Name: "set mutation"})
	input[3] = member{ID: 3, Name: "map mutation"}

	assert.NotContains(t, input, 2)
	assert.True(t, values.Contains(member{ID: 2}))
	assert.False(t, values.Contains(member{ID: 3}))
	assert.Equal(t, 2, values.Length())
}
