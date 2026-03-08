package doublelinked

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewDoubleLinkedList(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.first)
	assert.Nil(t, list.last)
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

func TestDoubleLinkedListGet(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		x, err := list.Get(i)
		assert.Equal(t, values[i], x)
		assert.Nil(t, err)
	}
}

func TestReversedDoubleLinkedListGet(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)
	list.Reverse()

	inversedIndex := 3
	for i := range 4 {
		x, err := list.Get(i)
		assert.Equal(t, values[inversedIndex], x)
		assert.Nil(t, err)
		inversedIndex--
	}
}

func TestDoubleLinkedListGetInvalidIndex(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Append(10, 1, 9, 100)

	for _, i := range []int{-1, 4, 5} {
		_, err := list.Get(i)
		target := lists.NewIndexOutOfBoundError(i, list.length)
		template := "index %d is out of bounds; maximum valid index is %d"
		assert.ErrorAsf(t, err, &target, template, i, list.length)
		assert.EqualErrorf(t, err, err.Error(), template, i, list.length)
	}
}

func TestDoubleLinkedListSet(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)

	for i := range 4 {
		err := list.Set(i, values[i]+1)
		x, _ := list.Get(i)
		assert.Equal(t, values[i]+1, x)
		assert.Nil(t, err)
	}
}

func TestReversedDoubleLinkedListSet(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	values := []int{10, 1, 9, 100}
	list.Append(values...)
	list.Reverse()

	inversedIndex := 3
	for i := range 4 {
		err := list.Set(i, values[inversedIndex]+1)
		x, _ := list.Get(i)
		assert.Equal(t, values[inversedIndex]+1, x)
		assert.Nil(t, err)
	}
}

func TestDoubleLinkedListSetInvalidIndex(t *testing.T) {
	list := NewDoubleLinkedList[int]()
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
	assert.Nil(t, list.first)
	assert.Nil(t, list.last)
}

func TestDoubleLinkedListIterator(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	idx := 0
	for value := range list.Iterator() {
		assert.Equal(t, items[idx], value)
		idx++
	}
}

func TestDoubleLinkedListReversedIterator(t *testing.T) {
	list := NewDoubleLinkedList[int]()
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

func TestDoubleLinkedListEnumerate(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestDoubleLinkedListReversedEnumerate(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	items := []int{1, 2, 3}
	list.Append(items...)
	list.Reverse()
	slices.Reverse(items)

	for idx, value := range list.Enumerate() {
		assert.Equal(t, items[idx], value)
	}
}

func TestDoubleLinkedListString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20)
		got := list.String()
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("String Many Elements", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := list.String()
		want := "[30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestDoubleLinkedListFormat(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20)
		got := fmt.Sprintf("%v", list)
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%#v", list)
		want := "*doublelinked.DoubleLinkedList[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose + String", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", list)
		want := "*doublelinked.DoubleLinkedList[int]{len:10, cap:10} [30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func DoubleLinkedListContains[T any](t *testing.T, items []T, list *DoubleLinkedList[T]) {
	current := list.first
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
