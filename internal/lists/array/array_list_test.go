package array

import (
	"testing"

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

func TestArrayListAppend(t *testing.T) {
	list := NewArrayList[int](0)
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, list.Length())
	assert.Equal(t, 3, len(list.data))
	// Golang Slice Growth Strategy
	assert.Equal(t, 4, cap(list.data))
}

func TestDoubleLinkedListIterator(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)

	index := 0
	for it := list.Iterator(); it.HasNext(); index++ {
		assert.Equal(t, items[index], it.Next())
	}
}

func TestDoubleLinkedListReversedIterator(t *testing.T) {
	list := NewArrayList[int](3)
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()

	index := len(items) - 1
	for it := list.Iterator(); it.HasNext(); index-- {
		assert.Equal(t, items[index], it.Next())
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
