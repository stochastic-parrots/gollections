package sortedlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderedArraySortedList(t *testing.T) {
	l := NewOrderedArraySortedList[int](10)

	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, len(l.data))
	assert.Equal(t, 10, cap(l.data))
}

func TestNewOrderedArraySortedListFromSlice(t *testing.T) {
	data := []int{3, 1, 2}

	l := NewOrderedArraySortedListFromSlice(data)

	assert.Equal(t, []int{1, 2, 3}, data)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestNewOrderedArraySortedListCloneSlice(t *testing.T) {
	data := []int{3, 1, 2}

	l := NewOrderedArraySortedListCloneSlice(data)

	assert.Equal(t, []int{3, 1, 2}, data)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestNewOrderedArraySortedListFromSeq(t *testing.T) {
	seq := slices.Values([]int{3, 1, 2})

	l := NewOrderedArraySortedListFromSeq(seq)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestOrderedArraySortedList_IsEmpty(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)

	assert.True(t, l.IsEmpty())

	l.Add(1)

	assert.False(t, l.IsEmpty())
}

func TestOrderedArraySortedList_Length(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(3, 1, 2)

	assert.Equal(t, 3, l.Length())
}

func TestOrderedArraySortedList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(3, 1, 2)

		for idx, value := range []int{1, 2, 3} {
			x, err := l.Get(idx)
			assert.NoError(t, err)
			assert.Equal(t, value, x)
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		for _, idx := range []int{-1, 1, 2} {
			_, err := l.Get(idx)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, list.ErrIndexOutOfBound))
		}
	})
}

func TestOrderedArraySortedList_Bounds(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(1, 3, 3, 5)

	assert.Equal(t, 0, l.LowerBound(0))
	assert.Equal(t, 1, l.LowerBound(3))
	assert.Equal(t, 3, l.LowerBound(4))
	assert.Equal(t, 4, l.LowerBound(6))

	assert.Equal(t, 0, l.UpperBound(0))
	assert.Equal(t, 3, l.UpperBound(3))
	assert.Equal(t, 3, l.UpperBound(4))
	assert.Equal(t, 4, l.UpperBound(5))
	assert.Equal(t, 4, l.UpperBound(6))

	start, end := l.EqualRange(3)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	assert.Equal(t, 2, l.Count(3))

	start, end = l.EqualRange(0)
	assert.Equal(t, 0, start)
	assert.Equal(t, 0, end)
	assert.Zero(t, l.Count(0))

	start, end = l.EqualRange(2)
	assert.Equal(t, 1, start)
	assert.Equal(t, 1, end)
	assert.Zero(t, l.Count(2))

	start, end = l.EqualRange(6)
	assert.Equal(t, 4, start)
	assert.Equal(t, 4, end)
	assert.Zero(t, l.Count(6))
}

func TestOrderedArraySortedList_Find(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(3, 1, 2, 2)

		idx, ok := l.Find(2)

		assert.True(t, ok)
		assert.Equal(t, 1, idx)
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1, 2, 3)

		idx, ok := l.Find(4)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		idx, ok := l.Find(1)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})
}

func TestOrderedArraySortedList_Navigation(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(1, 3, 3, 5)

	value, idx, ok := l.Ceiling(2)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 1, idx)

	value, idx, ok = l.Ceiling(6)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = l.Floor(4)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = l.Floor(0)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = l.Higher(3)
	assert.True(t, ok)
	assert.Equal(t, 5, value)
	assert.Equal(t, 3, idx)

	value, idx, ok = l.Higher(5)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = l.Lower(4)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = l.Lower(1)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)
}

func TestOrderedArraySortedList_Contains(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(3, 1, 2)

	assert.True(t, l.Contains(2))
	assert.False(t, l.Contains(4))
}

func TestOrderedArraySortedList_FirstLast(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		first, firstOK := l.First()
		last, lastOK := l.Last()

		assert.False(t, firstOK)
		assert.Zero(t, first)
		assert.False(t, lastOK)
		assert.Zero(t, last)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(3, 1, 2)

		first, firstOK := l.First()
		last, lastOK := l.Last()

		assert.True(t, firstOK)
		assert.Equal(t, 1, first)
		assert.True(t, lastOK)
		assert.Equal(t, 3, last)
	})
}

func TestOrderedArraySortedList_Add(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		l.Add()

		assert.Nil(t, l.ToSlice())
	})

	t.Run("SingleValue", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(2)
		l.Add(1)
		l.Add(3)

		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("MultipleValues", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		l.Add(3, 1, 2, 2)

		assert.Equal(t, []int{1, 2, 2, 3}, l.ToSlice())
	})
}

func TestOrderedArraySortedList_Replace(t *testing.T) {
	t.Run("PreservesOrder", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1, 3, 5)

		err := l.Replace(1, 4)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 4, 5}, l.ToSlice())
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1)

		for _, idx := range []int{-1, 1} {
			err := l.Replace(idx, 1)
			assert.ErrorIs(t, err, list.ErrIndexOutOfBound)
		}
	})

	t.Run("OrderViolation", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1, 3, 5)

		assert.ErrorIs(t, l.Replace(1, 0), ErrOrderViolation)
		assert.ErrorIs(t, l.Replace(1, 6), ErrOrderViolation)
		assert.Equal(t, []int{1, 3, 5}, l.ToSlice())
	})
}

func TestOrderedArraySortedList_Remove(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(3, 1, 2, 2)

		ok := l.Remove(2)

		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1, 2, 3)

		ok := l.Remove(4)

		assert.False(t, ok)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		assert.False(t, l.Remove(1))
	})
}

func TestOrderedArraySortedList_Range(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(1, 3, 3, 5)

	assert.Equal(t, []int{3, 3}, slices.Collect(l.Range(2, 5)))
	assert.Equal(t, []int{3, 3}, slices.Collect(l.Range(3, 4)))
	assert.Empty(t, slices.Collect(l.Range(5, 2)))

	var values []int
	l.Range(2, 5)(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{3}, values)
}

func TestOrderedArraySortedList_Iterators(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(3, 1, 2)

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(l.All()))
	assert.Equal(t, []int{3, 2, 1}, slices.Collect(l.Backward()))

	var allValues []int
	l.All()(func(x int) bool {
		allValues = append(allValues, x)
		return false
	})
	assert.Equal(t, []int{1}, allValues)

	var backwardValues []int
	l.Backward()(func(x int) bool {
		backwardValues = append(backwardValues, x)
		return false
	})
	assert.Equal(t, []int{3}, backwardValues)

	indexes := make([]int, 0, l.Length())
	values := make([]int, 0, l.Length())
	for idx, value := range l.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}

	assert.Equal(t, []int{0, 1, 2}, indexes)
	assert.Equal(t, []int{1, 2, 3}, values)

	var stopped []int
	l.Enumerate()(func(_ int, x int) bool {
		stopped = append(stopped, x)
		return false
	})
	assert.Equal(t, []int{1}, stopped)

	empty := NewOrderedArraySortedList[int](0)
	assert.Empty(t, slices.Collect(empty.All()))
	assert.Empty(t, slices.Collect(empty.Backward()))
	count := 0
	for range empty.Enumerate() {
		count++
	}
	assert.Zero(t, count)
}

func TestOrderedArraySortedList_ToSlice(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	assert.Nil(t, l.ToSlice())

	l.Add(2, 1)
	slice := l.ToSlice()
	slice[0] = 99

	assert.Equal(t, []int{1, 2}, l.ToSlice())
}

func TestOrderedArraySortedList_Clear(t *testing.T) {
	l := NewOrderedArraySortedList[string](0)
	l.Add("a", "b")

	l.Clear()

	assert.True(t, l.IsEmpty())
	assert.Nil(t, l.ToSlice())
	assert.Empty(t, l.data[:cap(l.data)][0])
}

func TestOrderedArraySortedList_JSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(3, 1, 2)

		data, err := json.Marshal(l)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})

	t.Run("MarshalEmpty", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		data, err := json.Marshal(l)

		assert.NoError(t, err)
		assert.JSONEq(t, `[]`, string(data))
	})

	t.Run("Unmarshal", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(10)

		err := json.Unmarshal([]byte(`[3,1,2]`), l)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("UnmarshalReusesCapacityAndClearsValues", func(t *testing.T) {
		l := NewOrderedArraySortedList[string](4)
		l.Add("b", "a")
		backing := &l.data[:cap(l.data)][0]

		err := json.Unmarshal([]byte(`["c"]`), l)

		assert.NoError(t, err)
		assert.Same(t, backing, &l.data[:cap(l.data)][0])
		assert.Equal(t, 4, cap(l.data))
		assert.Equal(t, "c", l.data[0])
		assert.Empty(t, l.data[:cap(l.data)][1])
	})

	t.Run("UnmarshalReplacesStorageWhenCapacityIsSmall", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)

		err := json.Unmarshal([]byte(`[3,1,2]`), l)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		l := NewOrderedArraySortedList[int](0)
		l.Add(1)

		err := json.Unmarshal([]byte(`{}`), l)

		assert.Error(t, err)
		assert.Equal(t, []int{1}, l.ToSlice())
	})
}

func TestOrderedArraySortedList_Format(t *testing.T) {
	l := NewOrderedArraySortedList[int](0)
	l.Add(5, 4, 3, 2, 1, 0)

	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", l.String())
	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", fmt.Sprintf("%v", l))
	assert.Contains(t, fmt.Sprintf("%+v", l), "len:6")
	assert.Contains(t, fmt.Sprintf("%#v", l), "OrderedArraySortedList")
}
