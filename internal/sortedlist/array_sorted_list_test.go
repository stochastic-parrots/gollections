package sortedlist

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	ilist "github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stretchr/testify/assert"
)

func TestNewArraySortedList(t *testing.T) {
	list := NewArraySortedList[int](10, cmp.Compare[int])

	assert.True(t, list.IsEmpty())
	assert.Equal(t, 0, len(list.data))
	assert.Equal(t, 10, cap(list.data))
}

func TestNewArraySortedListFromSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
}

func TestNewArraySortedListCloneSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewArraySortedListCloneSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{3, 1, 2}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
}

func TestNewArraySortedListFromSeq(t *testing.T) {
	seq := slices.Values([]int{3, 1, 2})

	list := NewArraySortedListFromSeq(seq, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
}

func TestArraySortedList_IsEmpty(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])

	assert.True(t, list.IsEmpty())

	list.Add(1)

	assert.False(t, list.IsEmpty())
}

func TestArraySortedList_Length(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(3, 1, 2)

	assert.Equal(t, 3, list.Length())
}

func TestArraySortedList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2)

		for idx, value := range []int{1, 2, 3} {
			x, err := list.Get(idx)
			assert.NoError(t, err)
			assert.Equal(t, value, x)
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(1)

		for _, idx := range []int{-1, 1, 2} {
			_, err := list.Get(idx)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, ilist.ErrIndexOutOfBound))
		}
	})
}

func TestArraySortedList_Find(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2, 2)

		idx, ok := list.Find(2)

		assert.True(t, ok)
		assert.Equal(t, 1, idx)
	})

	t.Run("NonExistent", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(1, 2, 3)

		idx, ok := list.Find(4)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		idx, ok := list.Find(1)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})
}

func TestArraySortedList_Contains(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(3, 1, 2)

	assert.True(t, list.Contains(2))
	assert.False(t, list.Contains(4))
}

func TestArraySortedList_First(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		x, ok := list.First()

		assert.False(t, ok)
		assert.Equal(t, 0, x)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2)

		x, ok := list.First()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
	})
}

func TestArraySortedList_Last(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		x, ok := list.Last()

		assert.False(t, ok)
		assert.Equal(t, 0, x)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2)

		x, ok := list.Last()

		assert.True(t, ok)
		assert.Equal(t, 3, x)
	})
}

func TestArraySortedList_Add(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		list.Add()

		assert.Nil(t, list.ToSlice())
	})

	t.Run("SingleValue", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(2)
		list.Add(1)
		list.Add(3)

		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("MultipleValues", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		list.Add(3, 1, 2, 2)

		assert.Equal(t, []int{1, 2, 2, 3}, list.ToSlice())
	})
}

func TestArraySortedList_Remove(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2, 2)

		ok := list.Remove(2)

		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("NonExistent", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(1, 2, 3)

		ok := list.Remove(4)

		assert.False(t, ok)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		assert.False(t, list.Remove(1))
	})
}

func TestArraySortedList_All(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(3, 1, 2)

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(list.All()))

	var values []int
	list.All()(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{1}, values)
}

func TestArraySortedList_Backward(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(3, 1, 2)

	assert.Equal(t, []int{3, 2, 1}, slices.Collect(list.Backward()))

	var values []int
	list.Backward()(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{3}, values)
}

func TestArraySortedList_Enumerate(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(3, 1, 2)

	indexes := make([]int, 0, list.Length())
	values := make([]int, 0, list.Length())
	for idx, value := range list.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}

	assert.Equal(t, []int{0, 1, 2}, indexes)
	assert.Equal(t, []int{1, 2, 3}, values)

	var stopped []int
	list.Enumerate()(func(_ int, x int) bool {
		stopped = append(stopped, x)
		return false
	})
	assert.Equal(t, []int{1}, stopped)
}

func TestArraySortedList_ToSlice(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	assert.Nil(t, list.ToSlice())

	list.Add(2, 1)
	slice := list.ToSlice()
	slice[0] = 99

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestArraySortedList_Clear(t *testing.T) {
	list := NewArraySortedList[*int](0, func(a, b *int) int {
		return cmp.Compare(*a, *b)
	})
	x := 1
	y := 2
	list.Add(&x, &y)

	list.Clear()

	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.ToSlice())
	assert.Nil(t, list.data[:cap(list.data)][0])
}

func TestArraySortedList_JSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(3, 1, 2)

		data, err := json.Marshal(list)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})

	t.Run("MarshalEmpty", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])

		data, err := json.Marshal(list)

		assert.NoError(t, err)
		assert.JSONEq(t, `[]`, string(data))
	})

	t.Run("Unmarshal", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(10)

		err := json.Unmarshal([]byte(`[3,1,2]`), list)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		list := NewArraySortedList[int](0, cmp.Compare[int])
		list.Add(1)

		err := json.Unmarshal([]byte(`{}`), list)

		assert.Error(t, err)
		assert.Equal(t, []int{1}, list.ToSlice())
	})
}

func TestArraySortedList_Format(t *testing.T) {
	list := NewArraySortedList[int](0, cmp.Compare[int])
	list.Add(5, 4, 3, 2, 1, 0)

	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", list.String())
	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", fmt.Sprintf("%v", list))
	assert.Contains(t, fmt.Sprintf("%+v", list), "len:6")
	assert.Contains(t, fmt.Sprintf("%#v", list), "ArraySortedList")
}

func TestArraySortedList_DescendingComparator(t *testing.T) {
	list := NewArraySortedList[int](0, func(a, b int) int {
		return cmp.Compare(b, a)
	})

	list.Add(1, 3, 2)

	assert.Equal(t, []int{3, 2, 1}, list.ToSlice())
	first, ok := list.First()
	assert.True(t, ok)
	assert.Equal(t, 3, first)
	last, ok := list.Last()
	assert.True(t, ok)
	assert.Equal(t, 1, last)
}
