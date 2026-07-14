package sortedlist

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stretchr/testify/assert"
)

func TestNewArraySortedList(t *testing.T) {
	l := NewArraySortedList(10, cmp.Compare[int])

	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, len(l.data))
	assert.Equal(t, 10, cap(l.data))
}

func TestNewArraySortedListFromSlice(t *testing.T) {
	data := []int{3, 1, 2}

	l := NewArraySortedListFromSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, data)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestNewArraySortedListCloneSlice(t *testing.T) {
	data := []int{3, 1, 2}

	l := NewArraySortedListCloneSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{3, 1, 2}, data)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestNewArraySortedListFromSeq(t *testing.T) {
	seq := slices.Values([]int{3, 1, 2})

	l := NewArraySortedListFromSeq(seq, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestArraySortedList_IsEmpty(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])

	assert.True(t, l.IsEmpty())

	l.Add(1)

	assert.False(t, l.IsEmpty())
}

func TestArraySortedList_Length(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(3, 1, 2)

	assert.Equal(t, 3, l.Length())
}

func TestArraySortedList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2)

		for idx, value := range []int{1, 2, 3} {
			x, err := l.Get(idx)
			assert.NoError(t, err)
			assert.Equal(t, value, x)
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		for _, idx := range []int{-1, 1, 2} {
			_, err := l.Get(idx)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, list.ErrIndexOutOfBound))
		}
	})
}

func TestArraySortedList_LowerBound(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(1, 3, 3, 5)

	tests := []struct {
		name  string
		value int
		want  int
	}{
		{name: "BeforeFirst", value: 0, want: 0},
		{name: "ExistingDuplicate", value: 3, want: 1},
		{name: "BetweenValues", value: 4, want: 3},
		{name: "AfterLast", value: 6, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, l.LowerBound(tt.value))
		})
	}

	empty := NewArraySortedList(0, cmp.Compare[int])
	assert.Zero(t, empty.LowerBound(1))
}

func TestArraySortedList_UpperBound(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(1, 3, 3, 5)

	tests := []struct {
		name  string
		value int
		want  int
	}{
		{name: "BeforeFirst", value: 0, want: 0},
		{name: "ExistingDuplicate", value: 3, want: 3},
		{name: "BetweenValues", value: 4, want: 3},
		{name: "ExistingLast", value: 5, want: 4},
		{name: "AfterLast", value: 6, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, l.UpperBound(tt.value))
		})
	}

	empty := NewArraySortedList(0, cmp.Compare[int])
	assert.Zero(t, empty.UpperBound(1))
}

func TestArraySortedList_EqualRange(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(1, 3, 3, 5)

	tests := []struct {
		name      string
		value     int
		wantStart int
		wantEnd   int
	}{
		{name: "BeforeFirst", value: 0, wantStart: 0, wantEnd: 0},
		{name: "ExistingDuplicate", value: 3, wantStart: 1, wantEnd: 3},
		{name: "BetweenValues", value: 2, wantStart: 1, wantEnd: 1},
		{name: "AfterLast", value: 6, wantStart: 4, wantEnd: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := l.EqualRange(tt.value)
			assert.Equal(t, tt.wantStart, start)
			assert.Equal(t, tt.wantEnd, end)
		})
	}
}

func TestArraySortedList_Count(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(1, 3, 3, 5)

	assert.Equal(t, 2, l.Count(3))
	assert.Zero(t, l.Count(0))
	assert.Zero(t, l.Count(2))
	assert.Zero(t, l.Count(6))
}

func TestArraySortedList_Find(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2, 2)

		idx, ok := l.Find(2)

		assert.True(t, ok)
		assert.Equal(t, 1, idx)
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1, 2, 3)

		idx, ok := l.Find(4)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		idx, ok := l.Find(1)

		assert.False(t, ok)
		assert.Equal(t, -1, idx)
	})
}

func TestArraySortedList_Navigation(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
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

func TestArraySortedList_Contains(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(3, 1, 2)

	assert.True(t, l.Contains(2))
	assert.False(t, l.Contains(4))
}

func TestArraySortedList_First(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		x, ok := l.First()

		assert.False(t, ok)
		assert.Equal(t, 0, x)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2)

		x, ok := l.First()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
	})
}

func TestArraySortedList_Last(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		x, ok := l.Last()

		assert.False(t, ok)
		assert.Equal(t, 0, x)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2)

		x, ok := l.Last()

		assert.True(t, ok)
		assert.Equal(t, 3, x)
	})
}

func TestArraySortedList_Add(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		l.Add()

		assert.Nil(t, l.ToSlice())
	})

	t.Run("SingleValue", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(2)
		l.Add(1)
		l.Add(3)

		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("MultipleValues", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		l.Add(3, 1, 2, 2)

		assert.Equal(t, []int{1, 2, 2, 3}, l.ToSlice())
	})
}

func TestArraySortedList_Replace(t *testing.T) {
	t.Run("PreservesOrder", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1, 3, 5)

		err := l.Replace(1, 4)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 4, 5}, l.ToSlice())
	})

	t.Run("EquivalentValue", func(t *testing.T) {
		type item struct {
			priority int
			name     string
		}
		compare := func(a, b item) int {
			return cmp.Compare(a.priority, b.priority)
		}
		items := NewArraySortedList(0, compare)
		items.Add(item{priority: 1, name: "old"})

		err := items.Replace(0, item{priority: 1, name: "new"})

		assert.NoError(t, err)
		value, getErr := items.Get(0)
		assert.NoError(t, getErr)
		assert.Equal(t, "new", value.name)
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1)

		for _, idx := range []int{-1, 1} {
			err := l.Replace(idx, 1)
			assert.ErrorIs(t, err, list.ErrIndexOutOfBound)
		}
	})

	t.Run("BeforePrevious", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1, 3, 5)

		err := l.Replace(1, 0)

		assert.ErrorIs(t, err, ErrOrderViolation)
		assert.Equal(t, []int{1, 3, 5}, l.ToSlice())
	})

	t.Run("AfterNext", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1, 3, 5)

		err := l.Replace(1, 6)

		assert.ErrorIs(t, err, ErrOrderViolation)
		assert.Equal(t, []int{1, 3, 5}, l.ToSlice())
	})
}

func TestArraySortedList_Remove(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2, 2)

		ok := l.Remove(2)

		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(1, 2, 3)

		ok := l.Remove(4)

		assert.False(t, ok)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		assert.False(t, l.Remove(1))
	})

	t.Run("ClearsDiscardedReference", func(t *testing.T) {
		l := NewArraySortedList(0, func(a, b *int) int {
			return cmp.Compare(*a, *b)
		})
		x := 1
		y := 2
		l.Add(&x, &y)

		assert.True(t, l.Remove(&x))
		assert.Nil(t, l.data[:cap(l.data)][l.Length()])
	})
}

func TestArraySortedList_All(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(3, 1, 2)

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(l.All()))

	var values []int
	l.All()(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{1}, values)

	empty := NewArraySortedList(0, cmp.Compare[int])
	assert.Empty(t, slices.Collect(empty.All()))
}

func TestArraySortedList_Backward(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(3, 1, 2)

	assert.Equal(t, []int{3, 2, 1}, slices.Collect(l.Backward()))

	var values []int
	l.Backward()(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{3}, values)

	empty := NewArraySortedList(0, cmp.Compare[int])
	assert.Empty(t, slices.Collect(empty.Backward()))
}

func TestArraySortedList_Range(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
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

func TestArraySortedList_Enumerate(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	l.Add(3, 1, 2)

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

	empty := NewArraySortedList(0, cmp.Compare[int])
	count := 0
	for range empty.Enumerate() {
		count++
	}
	assert.Zero(t, count)
}

func TestArraySortedList_ToSlice(t *testing.T) {
	l := NewArraySortedList(0, cmp.Compare[int])
	assert.Nil(t, l.ToSlice())

	l.Add(2, 1)
	slice := l.ToSlice()
	slice[0] = 99

	assert.Equal(t, []int{1, 2}, l.ToSlice())
}

func TestArraySortedList_Clear(t *testing.T) {
	l := NewArraySortedList(0, func(a, b *int) int {
		return cmp.Compare(*a, *b)
	})
	x := 1
	y := 2
	l.Add(&x, &y)

	l.Clear()

	assert.True(t, l.IsEmpty())
	assert.Nil(t, l.ToSlice())
	assert.Nil(t, l.data[:cap(l.data)][0])
}

func TestArraySortedList_JSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])
		l.Add(3, 1, 2)

		data, err := json.Marshal(l)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})

	t.Run("MarshalEmpty", func(t *testing.T) {
		l := NewArraySortedList(0, cmp.Compare[int])

		data, err := json.Marshal(l)

		assert.NoError(t, err)
		assert.JSONEq(t, `[]`, string(data))
	})

}

func TestArraySortedList_Format(t *testing.T) {
	l := NewArraySortedList(10, cmp.Compare[int])
	l.Add(5, 4, 3, 2, 1, 0)

	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", l.String())
	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", fmt.Sprintf("%v", l))
	assert.Contains(t, fmt.Sprintf("%+v", l), "len:6, cap:10")
	assert.Contains(t, fmt.Sprintf("%#v", l), "size:6, cap:10")
}

func TestArraySortedList_DescendingComparator(t *testing.T) {
	l := NewArraySortedList(0, func(a, b int) int {
		return cmp.Compare(b, a)
	})

	l.Add(1, 3, 2)

	assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	first, ok := l.First()
	assert.True(t, ok)
	assert.Equal(t, 3, first)
	last, ok := l.Last()
	assert.True(t, ok)
	assert.Equal(t, 1, last)
	assert.Equal(t, []int{3, 2}, slices.Collect(l.Range(3, 1)))
}
