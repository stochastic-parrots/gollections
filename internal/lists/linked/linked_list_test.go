package linked

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewLinkedList(t *testing.T) {
	list := NewLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.first)
	assert.Nil(t, list.last)
}

func TestLinkedListLength(t *testing.T) {
	list := NewLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	list.Append(1, 2, 3)
	assert.Equal(t, 3, list.Length())
}

func TestLinkedListIsEmpty(t *testing.T) {
	list := NewLinkedList[any]()

	assert.True(t, list.IsEmpty())
	list.Append(1, 2, 3)
	assert.False(t, list.IsEmpty())
}

func TestLinkedListGet(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		x, err := list.Get(i)
		assert.Equal(t, values[i], x)
		assert.Nil(t, err)
	}
}

func TestLinkedListGetInvalidIndex(t *testing.T) {
	list := NewLinkedList[int]()
	list.Append(10, 1, 9, 100)

	for _, i := range []int{-1, 4, 5} {
		_, err := list.Get(i)
		target := lists.NewIndexOutOfBoundError(i, list.length)
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, list.length)
		assert.EqualErrorf(t, err, err.Error(), template, i, list.length)
	}
}

func TestLinkedListSet(t *testing.T) {
	list := NewLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		err := list.Set(i, values[i]+1)
		x, _ := list.Get(i)
		assert.Equal(t, values[i]+1, x)
		assert.Nil(t, err)
	}
}

func TestLinkedListSetInvalidIndex(t *testing.T) {
	list := NewLinkedList[int]()
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
		x, _ := list.Get(i)
		assert.Equal(t, values[i], x)
	}
}

func TestLinkedListAdd(t *testing.T) {
	list := NewLinkedList[int]()
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, list.Length())

	LinkedListContains(t, []int{1, 2, 3}, list)
}

func TestLinkedListReverse(t *testing.T) {
	list := NewLinkedList[int]()
	list.Append(1, 2, 3)
	list.Reverse()

	LinkedListContains(t, []int{3, 2, 1}, list)
}

func TestEmptyLinkedListReverse(t *testing.T) {
	list := NewLinkedList[int]()
	list.Reverse()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.first)
	assert.Nil(t, list.last)
}

func TestLinkedListIterator(t *testing.T) {
	list := NewLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	index := 0
	for it := list.Iterator(); it.HasNext(); index++ {
		assert.Equal(t, items[index], it.Next())
	}
}

func TestEmptyLinkedListString(t *testing.T) {
	list := NewLinkedList[int]()
	assert.Equal(t, "[]", list.String())
}

func TestLinkedListString(t *testing.T) {
	list := NewLinkedList[int]()
	list.Append(1, 10, 11)

	assert.Equal(t, "[1, 10, 11]", list.String())
}

func LinkedListContains[T any](t *testing.T, items []T, list *LinkedList[T]) {
	current := list.first
	for _, x := range items {
		assert.Equal(t, x, current.value)
		current = current.next
	}
}
