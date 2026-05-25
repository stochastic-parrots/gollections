package sortedlist_test

import (
	"cmp"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stochastic-parrots/gollections/sortedlist"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementSortedList(t *testing.T) {
	var _ sortedlist.SortedList[int] = sortedlist.NewArray(0, cmp.Compare[int])
	var _ sortedlist.SortedList[int] = sortedlist.ArrayClone([]int{1}, cmp.Compare[int])
	var _ sortedlist.SortedList[int] = sortedlist.ArrayFromSeq(slices.Values([]int{1}), cmp.Compare[int])
}

func TestNewArray(t *testing.T) {
	assertSortedListBehavior(t, sortedlist.NewArray(0, cmp.Compare[int]))
}

func TestArrayFrom(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.ArrayFrom(data, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{1, 2, 3}, data)
}

func TestArrayClone(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.ArrayClone(data, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, data)
}

func TestArrayFromSeq(t *testing.T) {
	source := list.NewArray[int](0)
	source.Append(3, 1, 2)

	list := sortedlist.ArrayFromSeq(source.All(), cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, source.ToSlice())
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, sortedlist.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := sortedlist.NewArray(0, cmp.Compare[int])
		mutable.Add(2, 1)

		view := sortedlist.AsReadonly(mutable)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{2, 1}, slices.Collect(view.Backward()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Length())
		assert.False(t, view.IsEmpty())

		x, err := view.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, x)

		idx, ok := view.Find(2)
		assert.True(t, ok)
		assert.Equal(t, 1, idx)

		first, ok := view.First()
		assert.True(t, ok)
		assert.Equal(t, 1, first)

		last, ok := view.Last()
		assert.True(t, ok)
		assert.Equal(t, 2, last)

		assert.True(t, view.Contains(2))
		assert.Equal(t, "[1 2]", view.String())

		mutable.Add(0)
		assert.Equal(t, []int{0, 1, 2}, view.ToSlice())

		data, err := json.Marshal(view)
		assert.NoError(t, err)
		assert.JSONEq(t, `[0,1,2]`, string(data))
	})
}

func assertSortedListBehavior(t *testing.T, list sortedlist.SortedList[int]) {
	t.Helper()

	assert.True(t, list.IsEmpty())

	list.Add(3, 1, 2, 2)

	assert.Equal(t, []int{1, 2, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{1, 2, 2, 3}, slices.Collect(list.All()))
	assert.Equal(t, []int{3, 2, 2, 1}, slices.Collect(list.Backward()))

	x, err := list.Get(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, x)

	idx, ok := list.Find(2)
	assert.True(t, ok)
	assert.Equal(t, 1, idx)

	first, ok := list.First()
	assert.True(t, ok)
	assert.Equal(t, 1, first)

	last, ok := list.Last()
	assert.True(t, ok)
	assert.Equal(t, 3, last)

	assert.True(t, list.Contains(3))
	assert.False(t, list.Contains(4))

	ok = list.Remove(2)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())

	data, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(data))

	err = json.Unmarshal([]byte(`[9,7,8]`), list)
	assert.NoError(t, err)
	assert.Equal(t, []int{7, 8, 9}, list.ToSlice())

	assert.Equal(t, "[7 8 9]", list.String())

	list.Clear()
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.ToSlice())
}

func collectIndexes(seq func(func(int, int) bool)) []int {
	var indexes []int
	for idx := range seq {
		indexes = append(indexes, idx)
	}
	return indexes
}
