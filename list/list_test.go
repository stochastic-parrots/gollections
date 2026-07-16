package list_test

import (
	"cmp"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementList(t *testing.T) {
	var _ *list.ArrayList[int] = list.Array[int]().New(0)
	var _ *list.ArrayList[int] = list.Array[int]().From([]int{1})
	var _ *list.ArrayList[int] = list.Array[int]().Clone([]int{1})
	var _ *list.ArrayList[int] = list.Array[int]().FromSeq(slices.Values([]int{1}))
	var _ *list.LinkedList[int] = list.Linked[int]().New()
	var _ *list.LinkedList[int] = list.Linked[int]().From([]int{1})
	var _ *list.LinkedList[int] = list.Linked[int]().FromSeq(slices.Values([]int{1}))
	var _ list.List[int] = list.Array[int]().New(0)
	var _ list.List[int] = list.Linked[int]().New()
}

func TestConcreteZeroValues(t *testing.T) {
	var array list.ArrayList[int]
	array.Appends(1, 2)
	assert.Equal(t, []int{1, 2}, array.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &array))
	assert.Equal(t, []int{3, 4}, array.ToSlice())

	var linked list.LinkedList[int]
	linked.Appends(1, 2)
	linked.Reverse()
	assert.Equal(t, []int{2, 1}, linked.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &linked))
	assert.Equal(t, []int{3, 4}, linked.ToSlice())
}

func TestArrayFactory_From(t *testing.T) {
	data := []int{1, 2}
	list := list.Array[int]().From(data)

	assert.NoError(t, list.Set(0, 10))
	assert.Equal(t, []int{10, 2}, data)
}

func TestArrayFactory_Clone(t *testing.T) {
	data := []int{1, 2}
	list := list.Array[int]().Clone(data)

	assert.NoError(t, list.Set(0, 10))
	assert.Equal(t, []int{1, 2}, data)
}

func TestArrayFactory_FromSeq(t *testing.T) {
	list := list.Array[int]().FromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestLinkedFactory_From(t *testing.T) {
	data := []int{1, 2}
	list := list.Linked[int]().From(data)
	data[0] = 10

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestLinkedFactory_FromSeq(t *testing.T) {
	list := list.Linked[int]().FromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestArrayFactory_New(t *testing.T) {
	assertListBehavior(t, list.Array[int]().New(1))
}

func TestLinkedFactory_New(t *testing.T) {
	assertListBehavior(t, list.Linked[int]().New())
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, list.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := list.Array[int]().New(0)
		mutable.Appends(1, 2)

		view := list.AsReadonly[int](mutable)
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, unmarshals)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{2, 1}, slices.Collect(view.Backward()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Length())
		assert.False(t, view.IsEmpty())

		x, err := view.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, x)

		idx, ok := view.Find(2, cmp.Compare[int])
		assert.True(t, ok)
		assert.Equal(t, 1, idx)

		assert.True(t, view.Contains(2, cmp.Compare[int]))
		assert.Equal(t, "[1 2]", view.String())

		mutable.Append(3)
		assert.Equal(t, []int{1, 2, 3}, view.ToSlice())

		data, err := json.Marshal(view)
		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})
}

func assertListBehavior(t *testing.T, l list.List[int]) {
	t.Helper()

	assert.True(t, l.IsEmpty())

	l.Appends(1, 3)
	err := l.Insert(1, 2)
	assert.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	assert.Equal(t, []int{1, 2, 3}, slices.Collect(l.All()))
	assert.Equal(t, []int{3, 2, 1}, slices.Collect(l.Backward()))

	x, err := l.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, x)

	idx, ok := l.Find(3, cmp.Compare[int])
	assert.True(t, ok)
	assert.Equal(t, 2, idx)

	assert.True(t, l.Contains(2, cmp.Compare[int]))

	err = l.Set(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 20, 3}, l.ToSlice())

	l.Reverse()
	assert.Equal(t, []int{3, 20, 1}, l.ToSlice())

	removed, err := l.Remove(1)
	assert.NoError(t, err)
	assert.Equal(t, 20, removed)
	assert.Equal(t, []int{3, 1}, l.ToSlice())

	data, err := json.Marshal(l)
	assert.NoError(t, err)
	assert.JSONEq(t, `[3,1]`, string(data))

	err = json.Unmarshal([]byte(`[8,9]`), l)
	assert.NoError(t, err)
	assert.Equal(t, []int{8, 9}, l.ToSlice())

	l.Clear()
	assert.True(t, l.IsEmpty())
	assert.Nil(t, l.ToSlice())
}

func collectIndexes(seq func(func(int, int) bool)) []int {
	var indexes []int
	for idx := range seq {
		indexes = append(indexes, idx)
	}
	return indexes
}
