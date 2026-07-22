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
	var _ *set.HashSet[int] = set.Hash[int]().New(0)
	var _ *set.HashSet[int] = set.Hash[int]().From([]int{1})
	var _ *set.HashSet[int] = set.Hash[int]().FromSeq(slices.Values([]int{1}))
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).New(0)
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).From([]member{{ID: 1}})
	var _ *set.KeyedHashSet[member, int] = set.HashSetBy(memberID).FromSeq(slices.Values([]member{{ID: 1}}))
	var _ set.Set[int] = set.Hash[int]().New(0)
	var _ set.Set[member] = set.HashSetBy(memberID).New(0)
}

func TestConcreteZeroValues(t *testing.T) {
	var values set.HashSet[int]

	assert.True(t, values.Add(1))
	assert.Equal(t, 1, values.Adds(2, 1))

	assert.Equal(t, 2, values.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestHashSetBy_NilIdentityFunction(t *testing.T) {
	assert.PanicsWithValue(t, "set: nil identity function", func() {
		set.HashSetBy[int, int](nil)
	})
}

func TestHashFactory_New(t *testing.T) {
	assertSetBehavior(t, set.Hash[int]().New(2))
}

func TestHashFactory_From(t *testing.T) {
	data := []int{1, 2, 1}
	values := set.Hash[int]().From(data)
	data[0] = 3

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestHashFactory_FromSeq(t *testing.T) {
	values := set.Hash[int]().FromSeq(slices.Values([]int{1, 2, 1}))

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
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

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, set.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := set.Hash[int]().From([]int{1, 2})
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
