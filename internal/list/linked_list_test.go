package list

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinkedList(t *testing.T) {
	l := NewLinkedList[any]()

	assert.Equal(t, 0, l.Length())
	assert.True(t, l.IsEmpty())
	assert.Nil(t, l.first)
	assert.Nil(t, l.last)
}

func TestLinkedList_Length(t *testing.T) {
	l := NewLinkedList[any]()

	assert.Equal(t, 0, l.Length())
	l.Append(1, 2, 3)
	assert.Equal(t, 3, l.Length())
}

func TestLinkedList_IsEmpty(t *testing.T) {
	l := NewLinkedList[any]()

	assert.True(t, l.IsEmpty())
	l.Append(1, 2, 3)
	assert.False(t, l.IsEmpty())
}

func TestLinkedList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)

		for i := range 4 {
			x, err := l.Get(i)
			assert.Equal(t, values[i], x)
			assert.Nil(t, err)
		}
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)
		l.Reverse()

		inversedIndex := 3
		for i := range 4 {
			x, err := l.Get(i)
			assert.Equal(t, values[inversedIndex], x)
			assert.Nil(t, err)
			inversedIndex--
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(10, 1, 9, 100)

		for _, i := range []int{-1, 4, 5} {
			_, err := l.Get(i)
			target := NewIndexOutOfBoundError(i, l.length)
			template := "index %d is out of bounds; maximum valid index is %d"
			assert.ErrorAsf(t, err, &target, template, i, l.length)
			assert.EqualErrorf(t, err, err.Error(), template, i, l.length)
		}
	})
}

func TestLinkedList_Find(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)

		for idx, value := range values {
			fidx, exists := l.Find(value, cmp.Compare[int])
			assert.Equal(t, idx, fidx)
			assert.True(t, exists)
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		idx, exists := l.Find(4, cmp.Compare[int])

		assert.Equal(t, -1, idx)
		assert.False(t, exists)
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		idx, exists := l.Find(2, cmp.Compare[int])

		assert.Equal(t, -1, idx)
		assert.False(t, exists)
	})
}

func TestLinkedList_Contains(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)

		for _, value := range values {
			assert.True(t, l.Contains(value, cmp.Compare[int]))
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)

		assert.False(t, l.Contains(4, cmp.Compare[int]))
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()

		assert.False(t, l.Contains(4, cmp.Compare[int]))
	})
}

func TestLinkedList_Set(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)

		for i := range 4 {
			err := l.Set(i, values[i]+1)
			x, _ := l.Get(i)
			assert.Equal(t, values[i]+1, x)
			assert.Nil(t, err)
		}
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)
		l.Reverse()

		inversedIndex := 3
		for i := range 4 {
			err := l.Set(i, values[inversedIndex]+1)
			x, _ := l.Get(i)
			assert.Equal(t, values[inversedIndex]+1, x)
			assert.Nil(t, err)
		}
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		l := NewLinkedList[int]()
		values := []int{10, 1, 9, 100}
		l.Append(values...)

		for _, i := range []int{-1, 4, 5} {
			err := l.Set(i, 0)
			target := NewIndexOutOfBoundError(i, l.length)
			template := "index %d is out of bounds; maximum valid index is %d"
			assert.ErrorAsf(t, err, &target, template, i, l.length)
			assert.EqualErrorf(t, err, err.Error(), template, i, l.length)
		}

		for i := range l.length {
			x, _ := l.Get(i)
			assert.Equal(t, values[i], x)
		}
	})
}

func TestLinkedList_Insert(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		err := l.Insert(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, l.Length())
		assert.Equal(t, []int{0}, l.ToSlice())
	})

	t.Run("AtBeginning", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2)
		err := l.Insert(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, 3, l.Length())
		assert.Equal(t, []int{0, 1, 2}, l.ToSlice())
	})

	t.Run("InMiddle", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 3, 4, 5)
		err := l.Insert(2, 2)
		assert.NoError(t, err)
		assert.Equal(t, 5, l.Length())
		assert.Equal(t, []int{1, 3, 2, 4, 5}, l.ToSlice())
	})

	t.Run("AtBeginning_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2)
		l.Reverse()

		err := l.Insert(0, 3)
		assert.NoError(t, err)
		assert.Equal(t, 3, l.Length())
		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
		assert.Equal(t, 3, l.first.Value)
	})

	t.Run("InMiddle_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 3)
		l.Reverse()

		err := l.Insert(1, 2)
		assert.NoError(t, err)
		assert.Equal(t, 3, l.Length())
		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	})

	t.Run("AtEnd", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2)
		err := l.Insert(2, 3)
		assert.NoError(t, err)
		assert.Equal(t, 3, l.Length())
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1)
		err := l.Insert(5, 10)
		assert.Error(t, err)
		assert.Equal(t, 1, l.Length())
		assert.Equal(t, []int{1}, l.ToSlice())
	})
}

func TestLinkedList_Remove(t *testing.T) {
	t.Run("OnlyElement", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1)
		val, err := l.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.True(t, l.IsEmpty())
		assert.Nil(t, l.first)
		assert.Nil(t, l.last)
	})

	t.Run("Head", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		val, err := l.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.Equal(t, []int{2, 3}, l.ToSlice())
	})

	t.Run("Tail", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		val, err := l.Remove(2)
		assert.NoError(t, err)
		assert.Equal(t, 3, val)
		assert.Equal(t, []int{1, 2}, l.ToSlice())
	})

	t.Run("Middle", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		val, err := l.Remove(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
		assert.Equal(t, []int{1, 3}, l.ToSlice())
	})

	t.Run("Head_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		val, err := l.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 3, val)
		assert.Equal(t, []int{2, 1}, l.ToSlice())
	})

	t.Run("Tail_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		val, err := l.Remove(2)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.Equal(t, []int{3, 2}, l.ToSlice())
	})

	t.Run("Middle_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		val, err := l.Remove(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
		assert.Equal(t, []int{3, 1}, l.ToSlice())
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1)
		_, err := l.Remove(1)
		assert.Error(t, err)
	})
}

func TestLinkedList_Append(t *testing.T) {
	l := NewLinkedList[int]()
	l.Append(1, 2, 3)

	assert.False(t, l.IsEmpty())
	assert.Equal(t, 3, l.Length())
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestLinkedList_Reverse(t *testing.T) {
	t.Run("NonEmptyList", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()
		l.Append(4)

		assert.Equal(t, []int{3, 2, 1, 4}, l.ToSlice())
	})

	t.Run("EmptyList", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Reverse()

		assert.Equal(t, 0, l.Length())
		assert.True(t, l.IsEmpty())
		assert.Nil(t, l.first)
		assert.Nil(t, l.last)
	})
}

func TestLinkedList_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		idx := 0
		for value := range l.All() {
			assert.Equal(t, items[idx], value)
			idx++
		}
	})

	t.Run("PartialIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		idx := 0
		for value := range l.All() {
			assert.Equal(t, items[idx], value)
			idx++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 1, idx)
		assert.Equal(t, 3, l.Length())
	})

	t.Run("FullIteration_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)
		l.Reverse()
		slices.Reverse(items)

		idx := 0
		for value := range l.All() {
			assert.Equal(t, items[idx], value)
			idx++
		}
	})

	t.Run("PartialIteration_Reverserd", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)
		l.Reverse()
		slices.Reverse(items)

		idx := 0
		for value := range l.All() {
			assert.Equal(t, items[idx], value)
			idx++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 1, idx)
		assert.Equal(t, 3, l.Length())
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		count := 0

		for range l.All() {
			count++
		}

		assert.Equal(t, 0, l.Length())
		assert.Equal(t, 0, count)
	})
}

func TestLinkedList_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		for idx, value := range l.Enumerate() {
			assert.Equal(t, items[idx], value)
		}

		assert.Equal(t, 3, l.Length())
	})

	t.Run("PartialIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		count := 0
		for idx, value := range l.Enumerate() {
			assert.Equal(t, items[idx], value)
			count++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 2, count)
		assert.Equal(t, 3, l.Length())
	})

	t.Run("FullIteration_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)
		l.Reverse()
		slices.Reverse(items)

		for idx, value := range l.Enumerate() {
			assert.Equal(t, items[idx], value)
		}

		assert.Equal(t, 3, l.Length())
	})

	t.Run("PartialIteration_Reverserd", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)
		l.Reverse()
		slices.Reverse(items)

		count := 0
		for idx, value := range l.Enumerate() {
			assert.Equal(t, items[idx], value)
			count++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 3, l.Length())
		assert.Equal(t, 2, count)
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		count := 0

		for range l.Enumerate() {
			count++
		}

		assert.Equal(t, 0, l.Length())
		assert.Equal(t, 0, count)
	})
}

func TestLinkedList_Backward(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		var got []int
		for v := range l.Backward() {
			got = append(got, v)
		}

		assert.Equal(t, []int{3, 2, 1}, got)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3, 4)

		count := 0
		for v := range l.Backward() {
			count++
			if v == 3 {
				break
			}
		}

		assert.Equal(t, 2, count)
	})

	t.Run("FullIteration_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		var got []int
		for v := range l.Backward() {
			got = append(got, v)
		}

		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("PartialIteration_Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3, 4)
		l.Reverse()

		count := 0
		for v := range l.Backward() {
			count++
			if v == 3 {
				break
			}
		}

		assert.Equal(t, 3, count)
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		count := 0
		for range l.Backward() {
			count++
		}

		assert.Equal(t, 0, count)
	})
}

func TestLinkedList_String(t *testing.T) {
	t.Run("InLimit", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(30, 10, 20)
		got := l.String()
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("ExceedsLimit", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := l.String()
		want := "[30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestLinkedList_Format(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(30, 10, 20)
		got := fmt.Sprintf("%v", l)
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%#v", l)
		want := "*list.LinkedList[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("GoSyntax", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", l)
		want := "*list.LinkedList[int]{len:10, cap:10} [30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestLinkedList_MarshalJSON(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(10, 20, 30)

		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[10,20,30]", string(got))
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[3,2,1]", string(got))
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(got))
	})
}

func TestLinkedList_UnmarshalJSON(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		data := []byte("[4,5,6]")
		l := NewLinkedList[int]()
		l.Append(1)

		err := l.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.Equal(t, 3, l.Length())
		assert.Equal(t, []int{4, 5, 6}, l.ToSlice())

		assert.Equal(t, 4, l.first.Value)
		assert.Equal(t, 6, l.last.Value)
	})

	t.Run("Empty", func(t *testing.T) {
		data := []byte("[]")
		l := NewLinkedList[int]()
		l.Append(1, 2)

		err := l.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.True(t, l.IsEmpty())
		assert.Nil(t, l.first)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		data := []byte("[1, 2, 'error']")
		l := NewLinkedList[int]()
		err := l.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}

func TestLinkedList_ToSlice(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewLinkedList[int]()
		assert.Nil(t, l.ToSlice())
	})
}

func TestLinkedList_Clear(t *testing.T) {
	l := NewLinkedList[int]()
	l.Append(1, 2, 3)
	l.Clear()

	assert.Equal(t, 0, l.Length())
	assert.Nil(t, l.first)
	assert.Nil(t, l.last)
}
