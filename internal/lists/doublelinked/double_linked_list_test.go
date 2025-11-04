package doublelinked

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoubleLinkedList(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.First)
	assert.Nil(t, list.Last)
}

func TestDoubleLinkedListLength(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	list.Append(1, 2, 3)
	assert.Equal(t, 3, list.Length())
}

func TestDoubleLinkedListIsEmpty(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.True(t, list.IsEmpty())
	list.Append(1, 2, 3)
	assert.False(t, list.IsEmpty())
}

func TestDoubleLinkedListAppend(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, list.Length())

	DoubleLinkedListContains(t, []int{1, 2, 3}, list)
}

func TestDoubleLinkedListReverse(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Append(1, 2, 3)
	list.Reverse()
	list.Append(4)

	DoubleLinkedListContains(t, []int{3, 2, 1, 4}, list)
}

func TestEmptyDoubleLinkedListReverse(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Reverse()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.First)
	assert.Nil(t, list.Last)
}

func TestDoubleLinkedListIterator(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	index := 0
	for it := list.Iterator(); it.HasNext(); index++ {
		assert.Equal(t, items[index], it.Next())
	}
}

func TestDoubleLinkedListReversedIterator(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()

	index := len(items) - 1
	for it := list.Iterator(); it.HasNext(); index-- {
		assert.Equal(t, items[index], it.Next())
	}
}

func TestEmptyLinkedListString(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	assert.Equal(t, "[]", list.String())
}

func TestLinkedListString(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Append(1, 10, 11)

	assert.Equal(t, "[1, 10, 11]", list.String())
}

func DoubleLinkedListContains[T any](t *testing.T, items []T, list *DoubleLinkedList[T]) {
	current := list.First
	for index, x := range items {
		assert.Equal(t, x, current.value)

		if index == 0 && list.reversed {
			current = current.previous
			continue
		}

		if index == 0 && !list.reversed {
			current = current.next
			continue
		}

		if !list.reversed {
			assert.Equal(t, items[index-1], current.previous.value)
			current = current.next
			continue
		}

		if list.reversed {
			assert.Equal(t, items[index-1], current.next.value)
			current = current.previous
			continue
		}
	}
}
