package list

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoubleLinkedList(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.first)
	assert.Nil(t, list.last)
}

func TestDoubleLinkedList_Length(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.Equal(t, 0, list.Length())
	list.Append(1, 2, 3)
	assert.Equal(t, 3, list.Length())
}

func TestDoubleLinkedList_IsEmpty(t *testing.T) {
	list := NewDoubleLinkedList[any]()

	assert.True(t, list.IsEmpty())
	list.Append(1, 2, 3)
	assert.False(t, list.IsEmpty())
}

func TestDoubleLinkedList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		values := []int{10, 1, 9, 100}
		list.Append(values...)

		for i := range 4 {
			x, err := list.Get(i)
			assert.Equal(t, values[i], x)
			assert.Nil(t, err)
		}
	})

	t.Run("Reversed", func(t *testing.T) {
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
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(10, 1, 9, 100)

		for _, i := range []int{-1, 4, 5} {
			_, err := list.Get(i)
			target := NewIndexOutOfBoundError(i, list.length)
			template := "index %d is out of bounds; maximum valid index is %d"
			assert.ErrorAsf(t, err, &target, template, i, list.length)
			assert.EqualErrorf(t, err, err.Error(), template, i, list.length)
		}
	})
}

func TestDoubleLinkedList_Find(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		values := []int{10, 1, 9, 100}
		list.Append(values...)

		for idx, value := range values {
			fidx, exists := list.Find(value, cmp.Compare[int])
			assert.Equal(t, idx, fidx)
			assert.True(t, exists)
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		idx, exists := list.Find(4, cmp.Compare[int])

		assert.Equal(t, -1, idx)
		assert.False(t, exists)
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		idx, exists := list.Find(2, cmp.Compare[int])

		assert.Equal(t, -1, idx)
		assert.False(t, exists)
	})
}

func TestDoubleLinkedList_Contains(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		values := []int{10, 1, 9, 100}
		list.Append(values...)

		for _, value := range values {
			assert.True(t, list.Contains(value, cmp.Compare[int]))
		}
	})

	t.Run("NonExistent", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)

		assert.False(t, list.Contains(4, cmp.Compare[int]))
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()

		assert.False(t, list.Contains(4, cmp.Compare[int]))
	})
}

func TestDoubleLinkedList_Set(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		values := []int{10, 1, 9, 100}
		list.Append(values...)

		for i := range 4 {
			err := list.Set(i, values[i]+1)
			x, _ := list.Get(i)
			assert.Equal(t, values[i]+1, x)
			assert.Nil(t, err)
		}
	})

	t.Run("Reversed", func(t *testing.T) {
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
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
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
	})
}

func TestDoubleLinkedList_Insert(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		err := list.Insert(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, []int{0}, list.ToSlice())
	})

	t.Run("AtBeginning", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2)
		err := list.Insert(0, 0)
		assert.NoError(t, err)
		assert.Equal(t, []int{0, 1, 2}, list.ToSlice())
	})

	t.Run("InMiddle", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 3)
		err := list.Insert(1, 2)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("AtBeginning_Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2)
		l.Reverse()

		err := l.Insert(0, 3)
		assert.NoError(t, err)
		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
		assert.Equal(t, 3, l.first.Value)
	})

	t.Run("InMiddle_Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 3)
		l.Reverse()

		err := l.Insert(1, 2)
		assert.NoError(t, err)
		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	})

	t.Run("AtEnd", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2)
		err := list.Insert(2, 3)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1)
		err := list.Insert(5, 10)
		assert.Error(t, err)
		assert.Equal(t, []int{1}, list.ToSlice())
	})
}

func TestDoubleLinkedList_Remove(t *testing.T) {
	t.Run("OnlyElement", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1)
		val, err := list.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.True(t, list.IsEmpty())
		assert.Nil(t, list.first)
		assert.Nil(t, list.last)
	})

	t.Run("Head", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		val, err := list.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.Equal(t, []int{2, 3}, list.ToSlice())
	})

	t.Run("Tail", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		val, err := list.Remove(2)
		assert.NoError(t, err)
		assert.Equal(t, 3, val)
		assert.Equal(t, []int{1, 2}, list.ToSlice())
	})

	t.Run("Middle", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		val, err := list.Remove(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
		assert.Equal(t, []int{1, 3}, list.ToSlice())
	})

	t.Run("Head_Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		val, err := l.Remove(0)
		assert.NoError(t, err)
		assert.Equal(t, 3, val)
		assert.Equal(t, []int{2, 1}, l.ToSlice())
	})

	t.Run("Tail_Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		val, err := l.Remove(2)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
		assert.Equal(t, []int{3, 2}, l.ToSlice())
	})

	t.Run("Middle_Reversed", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		list.Reverse()

		val, err := list.Remove(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
		assert.Equal(t, []int{3, 1}, list.ToSlice())
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1)
		_, err := list.Remove(1)
		assert.Error(t, err)
	})
}

func TestDoubleLinkedList_Append(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	list.Append(1, 2, 3)

	assert.False(t, list.IsEmpty())
	assert.Equal(t, 3, list.Length())

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
}

func TestDoubleLinkedList_Reverse(t *testing.T) {
	t.Run("NonEmptyList", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(1, 2, 3)
		list.Reverse()
		list.Append(4)

		assert.Equal(t, []int{3, 2, 1, 4}, list.ToSlice())
	})

	t.Run("EmptyList", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Reverse()

		assert.Equal(t, 0, list.Length())
		assert.True(t, list.IsEmpty())
		assert.Nil(t, list.first)
		assert.Nil(t, list.last)
	})
}

func TestDoubleLinkedList_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)

		idx := 0
		for value := range list.All() {
			assert.Equal(t, items[idx], value)
			idx++
		}
	})

	t.Run("PartialIteration", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)

		idx := 0
		for value := range list.All() {
			assert.Equal(t, items[idx], value)
			idx++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 1, idx)
		assert.Equal(t, 3, list.Length())
	})

	t.Run("FullIteration_Reversed", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)
		list.Reverse()
		slices.Reverse(items)

		idx := 0
		for value := range list.All() {
			assert.Equal(t, items[idx], value)
			idx++
		}
	})

	t.Run("PartialIteration_Reverserd", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)
		list.Reverse()
		slices.Reverse(items)

		idx := 0
		for value := range list.All() {
			assert.Equal(t, items[idx], value)
			idx++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 1, idx)
		assert.Equal(t, 3, list.Length())
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		count := 0

		for range list.All() {
			count++
		}

		assert.Equal(t, 0, list.Length())
		assert.Equal(t, 0, count)
	})
}

func TestDoubleLinkedList_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)

		for idx, value := range list.Enumerate() {
			assert.Equal(t, items[idx], value)
		}

		assert.Equal(t, 3, list.Length())
	})

	t.Run("PartialIteration", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)

		count := 0
		for idx, value := range list.Enumerate() {
			assert.Equal(t, items[idx], value)
			count++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 2, count)
		assert.Equal(t, 3, list.Length())
	})

	t.Run("FullIteration_Reversed", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)
		list.Reverse()
		slices.Reverse(items)

		for idx, value := range list.Enumerate() {
			assert.Equal(t, items[idx], value)
		}

		assert.Equal(t, 3, list.Length())
	})

	t.Run("PartialIteration_Reverserd", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		list.Append(items...)
		list.Reverse()
		slices.Reverse(items)

		count := 0
		for idx, value := range list.Enumerate() {
			assert.Equal(t, items[idx], value)
			count++
			if idx == 1 {
				break
			}
		}

		assert.Equal(t, 3, list.Length())
		assert.Equal(t, 2, count)
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		count := 0

		for range list.Enumerate() {
			count++
		}

		assert.Equal(t, 0, list.Length())
		assert.Equal(t, 0, count)
	})
}

func TestDoubleLinkedList_Backward(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		items := []int{1, 2, 3}
		l.Append(items...)

		var got []int
		for v := range l.Backward() {
			got = append(got, v)
		}

		assert.Equal(t, []int{3, 2, 1}, got)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
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
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		var got []int
		for v := range l.Backward() {
			got = append(got, v)
		}

		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("PartialIteration_Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
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
		l := NewDoubleLinkedList[int]()
		count := 0
		for range l.Backward() {
			count++
		}

		assert.Equal(t, 0, count)
	})
}

func TestDoubleLinkedList_String(t *testing.T) {
	t.Run("InLimit", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20)
		got := list.String()
		want := "[30 10 20]"
		assert.Equal(t, want, got)
	})

	t.Run("ExceedsLimit", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := list.String()
		want := "[30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestDoubleLinkedList_Format(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
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
		want := "*list.DoubleLinkedList[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("GoSyntax", func(t *testing.T) {
		list := NewDoubleLinkedList[int]()
		list.Append(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", list)
		want := "*list.DoubleLinkedList[int]{len:10, cap:10} [30 10 20 40 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestDoubleLinkedList_MarshalJSON(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(10, 20, 30)

		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[10,20,30]", string(got))
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[3,2,1]", string(got))
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		got, err := l.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(got))
	})
}

func TestDoubleLinkedList_UnmarshalJSON(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		data := []byte("[4,5,6]")
		l := NewDoubleLinkedList[int]()
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
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2)

		err := l.UnmarshalJSON(data)
		assert.NoError(t, err)
		assert.True(t, l.IsEmpty())
		assert.Nil(t, l.first)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		data := []byte("[1, 2, 'error']")
		l := NewDoubleLinkedList[int]()
		err := l.UnmarshalJSON(data)
		assert.Error(t, err)
	})
}

func TestDoubleLinkedList_ToSlice(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()

		assert.Equal(t, []int{3, 2, 1}, l.ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		assert.Nil(t, l.ToSlice())
	})
}

func TestDoubleLinkedList_Clear(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Clear()

		assert.Equal(t, 0, l.Length())
		assert.Nil(t, l.first)
		assert.Nil(t, l.last)
		assert.False(t, l.reversed)
	})

	t.Run("Reversed", func(t *testing.T) {
		l := NewDoubleLinkedList[int]()
		l.Append(1, 2, 3)
		l.Reverse()
		l.Clear()

		assert.Equal(t, 0, l.Length())
		assert.Nil(t, l.first)
		assert.Nil(t, l.last)
		assert.False(t, l.reversed)

		l.Append(10)
		assert.Equal(t, []int{10}, l.ToSlice())
	})
}
