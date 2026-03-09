package list

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinkedNode(t *testing.T) {
	node := NewLinkedNode(5)

	assert.Equal(t, 5, node.Value())
	assert.Nil(t, node.Next())
	assert.False(t, node.HasNext())
}

func TestLinkedNodeHasNext(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.True(t, node.HasNext())
	assert.False(t, next.HasNext())
}

func TestLinkedNodeAppend(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Nil(t, next.Next())
	assert.Equal(t, 6, next.Value())
}

func TestLinkedNodeNext(t *testing.T) {
	node := NewLinkedNode(5)
	next := node.Append(6)

	assert.Same(t, next, node.Next())
	assert.Nil(t, next.Next())
}

func TestLinkedNodeValue(t *testing.T) {
	node := NewLinkedNode(10)
	other := NewLinkedNode("some")

	assert.Equal(t, 10, node.Value())
	assert.Equal(t, "some", other.Value())
}

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
		target := NewIndexOutOfBoundError(i, list.length)
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
		target := NewIndexOutOfBoundError(i, list.length)
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

func TestLinkedListAll(t *testing.T) {
	list := NewLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	idx := 0
	for value := range list.All() {
		assert.Equal(t, items[idx], value)
		idx++
	}
}

func TestLinkedListReversedAll(t *testing.T) {
	list := NewLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()
	slices.Reverse(items)

	idx := 0
	for value := range list.All() {
		assert.Equal(t, items[idx], value)
		idx++
	}
}

func TestLinkedListEnumerate(t *testing.T) {
	list := NewLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestLinkedListReversedEnumerate(t *testing.T) {
	list := NewLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()
	slices.Reverse(items)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestLinkedListString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(30, 10, 20)
		got := list.String()
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("String Many Elements", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := list.String()
		want := "[30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestLinkedListFormat(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(30, 10, 20)
		got := fmt.Sprintf("%v", list)
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%#v", list)
		want := "*list.LinkedList[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose + String", func(t *testing.T) {
		list := NewLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", list)
		want := "*list.LinkedList[int]{len:10, cap:10} [30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func LinkedListContains[T any](t *testing.T, items []T, list *LinkedList[T]) {
	current := list.first
	for _, x := range items {
		assert.Equal(t, x, current.value)
		current = current.next
	}
}
