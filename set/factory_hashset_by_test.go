package set_test

import (
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyedHashSet(t *testing.T) {
	values := set.NewKeyedHashSet(memberID, 2)
	first := member{ID: 1, Name: "first", Labels: []string{"a"}}

	assert.Equal(t, 1, values.Add(first))
	assert.Zero(t, values.Add(member{ID: 1, Name: "duplicate"}))

	assert.Equal(t, 1, values.Len())
	assert.Equal(t, []member{first}, slices.Collect(values.All()))
}

func TestKeyedHashSetFrom(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	data := []member{first, {ID: 1, Name: "duplicate"}, {ID: 2, Name: "second"}}
	values := set.KeyedHashSetFrom(memberID, data)
	data[0] = member{ID: 3, Name: "changed"}

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestKeyedHashSetFromSeq(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	values := set.KeyedHashSetFromSeq(memberID, slices.Values([]member{
		first,
		{ID: 1, Name: "duplicate"},
		{ID: 2, Name: "second"},
	}))

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestKeyedHashSet_FromMap(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	input := map[int]member{1: first}
	values := set.KeyedHashSetFromMap(memberID, input)

	assert.Equal(t, 1, values.Add(member{ID: 2, Name: "set mutation"}))
	assert.Equal(t, member{ID: 2, Name: "set mutation"}, input[2])

	input[3] = member{ID: 3, Name: "map mutation"}
	assert.True(t, values.Contains(member{ID: 3}))
}

func TestKeyedHashSet_CloneMap(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	input := map[int]member{1: first}
	values := set.KeyedHashSetCloneMap(memberID, input)

	values.Add(member{ID: 2, Name: "set mutation"})
	input[3] = member{ID: 3, Name: "map mutation"}

	assert.NotContains(t, input, 2)
	assert.True(t, values.Contains(member{ID: 2}))
	assert.False(t, values.Contains(member{ID: 3}))
	assert.Equal(t, 2, values.Len())
}

func TestKeyedHashSetConstructors_NilIdentityFunction(t *testing.T) {
	constructors := map[string]func(){
		"New":      func() { set.NewKeyedHashSet[int, int](nil, 0) },
		"From":     func() { set.KeyedHashSetFrom[int, int](nil, nil) },
		"FromMap":  func() { set.KeyedHashSetFromMap[int, int](nil, nil) },
		"CloneMap": func() { set.KeyedHashSetCloneMap[int, int](nil, nil) },
		"FromSeq":  func() { set.KeyedHashSetFromSeq[int, int](nil, slices.Values([]int(nil))) },
	}
	for name, construct := range constructors {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithValue(t, "set: nil identity function", construct)
		})
	}
}

func TestKeyedHashSet_FromMap_SymmetricDifferenceWithPreservesMap(t *testing.T) {
	input := map[int]member{}
	values := set.KeyedHashSetFromMap(memberID, input)

	assert.Equal(t, 1, values.SymmetricDifferenceWith(
		set.SourceFromSlice([]member{{ID: 1, Name: "inserted"}}),
	))
	assert.Equal(t, member{ID: 1, Name: "inserted"}, input[1])

	input[2] = member{ID: 2, Name: "external"}
	assert.True(t, values.Contains(member{ID: 2}))
}
