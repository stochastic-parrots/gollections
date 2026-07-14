package sortedlist_test

import (
	"cmp"
	"encoding/json"
	"errors"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stochastic-parrots/gollections/sortedlist"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementSortedList(t *testing.T) {
	var _ *sortedlist.ArraySortedList[int] = sortedlist.Array(cmp.Compare[int]).New(0)
	var _ *sortedlist.OrderedArraySortedList[int] = sortedlist.OrderedArray[int]().New(0)
	var _ sortedlist.SortedList[int] = sortedlist.Array(cmp.Compare[int]).New(0)
	var _ sortedlist.SortedList[int] = sortedlist.Array(cmp.Compare[int]).From([]int{1})
	var _ sortedlist.SortedList[int] = sortedlist.Array(cmp.Compare[int]).Clone([]int{1})
	var _ sortedlist.SortedList[int] = sortedlist.Array(cmp.Compare[int]).FromSeq(slices.Values([]int{1}))
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArray[int]().New(0)
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArray[int]().From([]int{1})
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArray[int]().Clone([]int{1})
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArray[int]().FromSeq(slices.Values([]int{1}))
}

func TestArrayFactory_New(t *testing.T) {
	assertSortedListBehavior(t, sortedlist.Array(cmp.Compare[int]).New(0))
}

func TestOrderedArrayFactory_New(t *testing.T) {
	assertSortedListBehavior(t, sortedlist.OrderedArray[int]().New(0))
}

func TestArrayFactory_From(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.Array(cmp.Compare[int]).From(data)

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{1, 2, 3}, data)
}

func TestArrayFactory_Clone(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.Array(cmp.Compare[int]).Clone(data)

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, data)
}

func TestArrayFactory_FromSeq(t *testing.T) {
	source := list.Array[int]().New(0)
	source.Append(3, 1, 2)

	list := sortedlist.Array(cmp.Compare[int]).FromSeq(source.All())

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, source.ToSlice())
}

func TestOrderedArrayFactory_From(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.OrderedArray[int]().From(data)

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{1, 2, 3}, data)
}

func TestOrderedArrayFactory_Clone(t *testing.T) {
	data := []int{3, 1, 2}

	list := sortedlist.OrderedArray[int]().Clone(data)

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, data)
}

func TestOrderedArrayFactory_FromSeq(t *testing.T) {
	source := list.Array[int]().New(0)
	source.Append(3, 1, 2)

	list := sortedlist.OrderedArray[int]().FromSeq(source.All())

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, source.ToSlice())
}

func TestArraySortedList_GetError(t *testing.T) {
	list := sortedlist.Array(cmp.Compare[int]).New(0)

	_, err := list.Get(0)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sortedlist.ErrIndexOutOfBound))

	var bounds *sortedlist.IndexOutOfBoundError
	assert.True(t, errors.As(err, &bounds))
	assert.Equal(t, 0, bounds.Index())
	assert.Equal(t, -1, bounds.Limit())
}

func TestArraySortedList_ReplaceError(t *testing.T) {
	list := sortedlist.Array(cmp.Compare[int]).New(0)
	list.Add(1, 3, 5)

	err := list.Replace(1, 6)

	assert.ErrorIs(t, err, sortedlist.ErrOrderViolation)
	assert.Equal(t, []int{1, 3, 5}, list.ToSlice())
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, sortedlist.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := sortedlist.Array(cmp.Compare[int]).New(0)
		mutable.Add(2, 1)

		view := sortedlist.AsReadonly(mutable)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{2, 1}, slices.Collect(view.Backward()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Length())
		assert.False(t, view.IsEmpty())
		assert.Equal(t, 1, view.LowerBound(2))
		assert.Equal(t, 2, view.UpperBound(2))
		start, end := view.EqualRange(2)
		assert.Equal(t, 1, start)
		assert.Equal(t, 2, end)
		assert.Equal(t, 1, view.Count(2))
		assert.Equal(t, []int{1}, slices.Collect(view.Range(1, 2)))

		x, err := view.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, x)

		idx, ok := view.Find(2)
		assert.True(t, ok)
		assert.Equal(t, 1, idx)

		value, idx, ok := view.Ceiling(2)
		assert.True(t, ok)
		assert.Equal(t, 2, value)
		assert.Equal(t, 1, idx)

		value, idx, ok = view.Floor(2)
		assert.True(t, ok)
		assert.Equal(t, 2, value)
		assert.Equal(t, 1, idx)

		value, idx, ok = view.Higher(1)
		assert.True(t, ok)
		assert.Equal(t, 2, value)
		assert.Equal(t, 1, idx)

		value, idx, ok = view.Lower(2)
		assert.True(t, ok)
		assert.Equal(t, 1, value)
		assert.Equal(t, 0, idx)

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
	assert.Equal(t, 1, list.LowerBound(2))
	assert.Equal(t, 3, list.UpperBound(2))
	start, end := list.EqualRange(2)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	assert.Equal(t, 2, list.Count(2))
	assert.Equal(t, []int{2, 2}, slices.Collect(list.Range(2, 3)))

	value, idx, ok := list.Ceiling(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)
	assert.Equal(t, 1, idx)

	value, idx, ok = list.Floor(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = list.Higher(2)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 3, idx)

	value, idx, ok = list.Lower(2)
	assert.True(t, ok)
	assert.Equal(t, 1, value)
	assert.Equal(t, 0, idx)

	err = list.Replace(2, 2)
	assert.NoError(t, err)

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
