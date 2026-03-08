package array

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewArrayList(t *testing.T) {
	list := NewArrayList[int](100)

	assert.True(t, list.IsEmpty())
	assert.Equal(t, 0, len(list.data))
	assert.Equal(t, 0, len(list.data))
	assert.Equal(t, 100, cap(list.data))
}

func TestArrayListIsEmpty(t *testing.T) {
	list := NewArrayList[int](100)
	assert.True(t, list.IsEmpty())
}

func TestArrayListLength(t *testing.T) {
	list := NewArrayList[int](100)
	list.Append(10, 1, 9, 100)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 4, len(list.data))
}

func TestArrayListGet(t *testing.T) {
	list := NewArrayList[int](100)
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		x, err := list.Get(i)
		assert.Equal(t, values[i], x)
		assert.Nil(t, err)
	}
}

func TestArrayListGetInvalidIndex(t *testing.T) {
	list := NewArrayList[int](100)
	list.Append(10, 1, 9, 100)

	for _, i := range []int{-1, 4, 5} {
		_, err := list.Get(i)
		target := lists.NewIndexOutOfBoundError(i, len(list.data))
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, len(list.data))
		assert.EqualErrorf(t, err, err.Error(), template, i, len(list.data))
	}
}

func TestArrayListSet(t *testing.T) {
	list := NewArrayList[int](100)
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		err := list.Set(i, values[i]+1)
		assert.Equal(t, values[i]+1, list.data[i])
		assert.Nil(t, err)
	}
}

func TestArrayListSetInvalidIndex(t *testing.T) {
	list := NewArrayList[int](100)
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for _, i := range []int{-1, 4, 5} {
		err := list.Set(i, 0)
		target := lists.NewIndexOutOfBoundError(i, len(list.data))
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, len(list.data))
		assert.EqualErrorf(t, err, err.Error(), template, i, len(list.data))
	}

	for i := range len(list.data) {

		assert.Equal(t, values[i], list.data[i])
	}
}

func TestArrayListAppend(t *testing.T) {
	list := NewArrayList[int](0)
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, len(list.data))
}

func TestArrayListListIterator(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)

	idx := 0
	for value := range list.Iterator() {
		assert.Equal(t, items[idx], value)
		idx++
	}
}

func TestArrayListListReversedIterator(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()
	slices.Reverse(items)

	idx := 0
	for value := range list.Iterator() {
		assert.Equal(t, items[idx], value)
		idx++
	}
}

func TestArrayListEnumerate(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestArrayListReversedEnumerate(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()
	slices.Reverse(items)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestArrayListReverse(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()

	assert.Equal(t, 3, len(list.data))
	assert.False(t, list.IsEmpty())
	ArrayListContains(t, []int{3, 2, 1}, list)
}
func TestArrayListString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		array := NewArrayList[int](3)
		array.Append(30, 10, 20)
		got := array.String()
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("String Many Elements", func(t *testing.T) {
		array := NewArrayList[int](10)
		array.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := array.String()
		want := "[30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestArrayListFormat(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		array := NewArrayList[int](3)
		array.Append(30, 10, 20)
		got := fmt.Sprintf("%v", array)
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose", func(t *testing.T) {
		array := NewArrayList[int](10)
		array.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%#v", array)
		want := "*array.ArrayList[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose + String", func(t *testing.T) {
		array := NewArrayList[int](10)
		array.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", array)
		want := "*array.ArrayList[int]{len:10, cap:10} [30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func ArrayListContains[T any](t *testing.T, items []T, list *ArrayList[T]) {
	for index, x := range items {
		assert.Equal(t, x, list.data[index])
	}
}
