package array

import (
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewArrayList(t *testing.T) {
	list := NewArrayList[int](100)

	assert.True(t, list.IsEmpty())
	assert.Equal(t, 0, list.Length())
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
	assert.Equal(t, 4, list.Length())
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
		target := lists.NewIndexOutOfBoundError(i, list.length)
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, list.length)
		assert.EqualErrorf(t, err, err.Error(), template, i, list.length)
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
		target := lists.NewIndexOutOfBoundError(i, list.length)
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, list.length)
		assert.EqualErrorf(t, err, err.Error(), template, i, list.length)
	}

	for i := range list.length {

		assert.Equal(t, values[i], list.data[i])
	}
}

func TestArrayListAppend(t *testing.T) {
	list := NewArrayList[int](0)
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, list.Length())
	assert.Equal(t, 3, len(list.data))
	// Golang Slice Growth Strategy
	assert.Equal(t, 4, cap(list.data))
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

	assert.Equal(t, 3, list.Length())
	assert.False(t, list.IsEmpty())
	ArrayListContains(t, []int{3, 2, 1}, list)
}

func TestEmptyArrayListString(t *testing.T) {
	list := NewArrayList[int](3)
	assert.Equal(t, "[]", list.String())
}

func TestArrayListString(t *testing.T) {
	list := NewArrayList[string](3)
	list.Append("a", "b", "c")
	assert.Equal(t, "[a, b, c]", list.String())
}

func ArrayListContains[T any](t *testing.T, items []T, list *ArrayList[T]) {
	for index, x := range items {
		assert.Equal(t, x, list.data[index])
	}
}
