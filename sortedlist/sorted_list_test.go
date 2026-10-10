package sortedlist_test

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stochastic-parrots/gollections/sortedlist"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementSortedList(t *testing.T) {
	assert.IsType(t, (*sortedlist.ArraySortedList[int])(nil), sortedlist.NewOrderedArray[int](sortedlist.Asc, 0))
	assert.IsType(t, (*sortedlist.ArraySortedList[int])(nil), sortedlist.NewArray(cmp.Compare[int], 0))
	assert.IsType(t, (*sortedlist.ArraySortedList[int])(nil), sortedlist.NewOrderedArray[int](sortedlist.Asc, 0))
	var _ sortedlist.SortedList[int] = sortedlist.NewArray(cmp.Compare[int], 0)
	var _ sortedlist.SortedList[int] = sortedlist.ArrayFrom(cmp.Compare[int], []int{1})
	var _ sortedlist.SortedList[int] = sortedlist.ArrayClone(cmp.Compare[int], []int{1})
	var _ sortedlist.SortedList[int] = sortedlist.ArrayFromSeq(cmp.Compare[int], slices.Values([]int{1}))
	var _ sortedlist.SortedList[int] = sortedlist.NewOrderedArray[int](sortedlist.Asc, 0)
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArrayFrom(sortedlist.Asc, []int{1})
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArrayClone(sortedlist.Asc, []int{1})
	var _ sortedlist.SortedList[int] = sortedlist.OrderedArrayFromSeq(sortedlist.Asc, slices.Values([]int{1}))
}

func TestArrayConstructors_NilComparator(t *testing.T) {
	constructors := map[string]func(){
		"New":     func() { sortedlist.NewArray[int](nil, 0) },
		"From":    func() { sortedlist.ArrayFrom[int](nil, nil) },
		"Clone":   func() { sortedlist.ArrayClone[int](nil, nil) },
		"FromSeq": func() { sortedlist.ArrayFromSeq(nil, slices.Values([]int(nil))) },
	}
	for name, construct := range constructors {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithValue(t, "sortedlist: nil comparator", construct)
		})
	}
}

func TestNewArray(t *testing.T) {
	assertSortedListBehavior(t, sortedlist.NewArray(cmp.Compare[int], 0))
}

func TestArrayFrom(t *testing.T) {
	data := []int{3, 1, 2}

	l := sortedlist.ArrayFrom(cmp.Compare[int], data)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestArrayClone(t *testing.T) {
	data := []int{3, 1, 2}

	l := sortedlist.ArrayClone(cmp.Compare[int], data)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, data)
	data[0] = 9
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	assert.NoError(t, l.Replace(0, 0))
	assert.Equal(t, []int{9, 1, 2}, data)
}

func TestArrayFromSeq(t *testing.T) {
	source := list.NewArray[int](0)
	source.Append(3, 1, 2)

	l := sortedlist.ArrayFromSeq(cmp.Compare[int], source.All())

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	assert.Equal(t, []int{3, 1, 2}, source.ToSlice())
}

func TestOrderedArray(t *testing.T) {
	t.Run("Ascending", func(t *testing.T) {
		l := sortedlist.OrderedArrayClone(sortedlist.Asc, []int{3, 1, 2})

		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("Descending", func(t *testing.T) {
		l := sortedlist.OrderedArrayClone(sortedlist.Desc, []int{3, 1, 2})

		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	})

	t.Run("FloatNaN", func(t *testing.T) {
		nan := math.NaN()

		ascending := sortedlist.OrderedArrayClone(sortedlist.Asc, []float64{1, nan, 2})
		descending := sortedlist.OrderedArrayClone(sortedlist.Desc, []float64{1, nan, 2})

		ascendingValues := ascending.ToSlice()
		descendingValues := descending.ToSlice()

		assert.True(t, math.IsNaN(ascendingValues[0]))
		assert.Equal(t, []float64{1, 2}, ascendingValues[1:])
		assert.Equal(t, []float64{2, 1}, descendingValues[:2])
		assert.True(t, math.IsNaN(descendingValues[2]))
	})
}

func TestArraySortedList_GetError(t *testing.T) {
	l := sortedlist.NewArray(cmp.Compare[int], 0)

	_, err := l.Get(0)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sortedlist.ErrIndexOutOfBounds))
	assert.EqualError(t, err, "cannot get index 0 because list is empty")

	var bounds *sortedlist.IndexOutOfBoundsError
	assert.True(t, errors.As(err, &bounds))
	assert.Equal(t, 0, bounds.Index())
	assert.Equal(t, -1, bounds.Limit())
}

func TestArraySortedList_ReplaceError(t *testing.T) {
	l := sortedlist.NewArray(cmp.Compare[int], 0)
	l.Add(1, 3, 5)

	err := l.Replace(1, 6)

	assert.ErrorIs(t, err, sortedlist.ErrOrderViolation)
	assert.Equal(t, []int{1, 3, 5}, l.ToSlice())
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() sortedlist.SortedList[int]{
		"ArrayFromNil":    func() sortedlist.SortedList[int] { return sortedlist.ArrayFrom(cmp.Compare[int], []int(nil)) },
		"ArrayFromEmpty":  func() sortedlist.SortedList[int] { return sortedlist.ArrayFrom(cmp.Compare[int], make([]int, 0, 4)) },
		"ArrayCloneNil":   func() sortedlist.SortedList[int] { return sortedlist.ArrayClone(cmp.Compare[int], []int(nil)) },
		"ArrayCloneEmpty": func() sortedlist.SortedList[int] { return sortedlist.ArrayClone(cmp.Compare[int], make([]int, 0, 4)) },
	} {
		t.Run(name, func(t *testing.T) {
			values := construct()
			assert.True(t, values.IsEmpty())
			values.Add(1)
			assert.Equal(t, []int{1}, values.ToSlice())
		})
	}
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, sortedlist.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := sortedlist.NewArray(cmp.Compare[int], 0)
		mutable.Add(2, 1)

		view := sortedlist.AsReadonly(mutable)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{2, 1}, slices.Collect(view.Backward()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Len())
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

func assertSortedListBehavior(t *testing.T, l sortedlist.SortedList[int]) {
	t.Helper()

	_, unmarshals := any(l).(json.Unmarshaler)
	assert.False(t, unmarshals)
	assert.True(t, l.IsEmpty())

	l.Add(3, 1, 2, 2)

	assert.Equal(t, []int{1, 2, 2, 3}, l.ToSlice())
	assert.Equal(t, []int{1, 2, 2, 3}, slices.Collect(l.All()))
	assert.Equal(t, []int{3, 2, 2, 1}, slices.Collect(l.Backward()))

	x, err := l.Get(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, x)

	idx, ok := l.Find(2)
	assert.True(t, ok)
	assert.Equal(t, 1, idx)
	assert.Equal(t, 1, l.LowerBound(2))
	assert.Equal(t, 3, l.UpperBound(2))
	start, end := l.EqualRange(2)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	assert.Equal(t, 2, l.Count(2))
	assert.Equal(t, []int{2, 2}, slices.Collect(l.Range(2, 3)))

	value, idx, ok := l.Ceiling(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)
	assert.Equal(t, 1, idx)

	value, idx, ok = l.Floor(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = l.Higher(2)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 3, idx)

	value, idx, ok = l.Lower(2)
	assert.True(t, ok)
	assert.Equal(t, 1, value)
	assert.Equal(t, 0, idx)

	err = l.Replace(2, 2)
	assert.NoError(t, err)

	first, ok := l.First()
	assert.True(t, ok)
	assert.Equal(t, 1, first)

	last, ok := l.Last()
	assert.True(t, ok)
	assert.Equal(t, 3, last)

	assert.True(t, l.Contains(3))
	assert.False(t, l.Contains(4))

	ok = l.Remove(2)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())

	data, err := json.Marshal(l)
	assert.NoError(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(data))
	assert.Equal(t, "[1 2 3]", fmt.Sprint(l))

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
