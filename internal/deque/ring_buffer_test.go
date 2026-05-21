package deque

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRingBufferDeque(t *testing.T) {
	deque := NewRingBufferDeque[int](10)

	assert.Equal(t, 0, deque.Length())
	assert.True(t, deque.IsEmpty())
	assert.Len(t, deque.data, 10)
	assert.Equal(t, 0, deque.read)
	assert.Equal(t, 0, deque.write)
}

func TestRingBufferDeque_Length(t *testing.T) {
	deque := NewRingBufferDeque[int](1)
	assert.Equal(t, 0, deque.Length())

	deque.Append(1, 2)
	assert.Equal(t, 2, deque.Length())
}

func TestRingBufferDeque_IsEmpty(t *testing.T) {
	deque := NewRingBufferDeque[int](1)
	assert.True(t, deque.IsEmpty())

	deque.Append(1)
	assert.False(t, deque.IsEmpty())
}

func TestRingBufferDeque_Front(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		x, ok := deque.Front()

		assert.False(t, ok)
		assert.Zero(t, x)
	})

	t.Run("Populated", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1, 2)

		x, ok := deque.Front()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
	})
}

func TestRingBufferDeque_Back(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		x, ok := deque.Back()

		assert.False(t, ok)
		assert.Zero(t, x)
	})

	t.Run("Populated", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1, 2)

		x, ok := deque.Back()

		assert.True(t, ok)
		assert.Equal(t, 2, x)
	})
}

func TestRingBufferDeque_Append(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		deque.Append(1)

		assert.Equal(t, []int{1}, deque.ToSlice())
	})

	t.Run("ManyWithGrowth", func(t *testing.T) {
		deque := NewRingBufferDeque[int](0)

		deque.Append(1, 2, 3)

		assert.Equal(t, []int{1, 2, 3}, deque.ToSlice())
		assert.GreaterOrEqual(t, len(deque.data), 3)
	})
}

func TestRingBufferDeque_Prepend(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		deque.Prepend(1)

		assert.Equal(t, []int{1}, deque.ToSlice())
	})

	t.Run("ManyPreservesInputOrder", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(4)

		deque.Prepend(1, 2, 3)

		assert.Equal(t, []int{1, 2, 3, 4}, deque.ToSlice())
	})
}

func TestRingBufferDeque_Shift(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		x, ok := deque.Shift()

		assert.False(t, ok)
		assert.Zero(t, x)
	})

	t.Run("SingleElement", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1)

		x, ok := deque.Shift()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
		assert.True(t, deque.IsEmpty())
		assert.Equal(t, 0, deque.read)
		assert.Equal(t, 0, deque.write)
	})

	t.Run("ManyElements", func(t *testing.T) {
		deque := NewRingBufferDeque[int](2)
		deque.Append(1, 2, 3)

		x, ok := deque.Shift()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
		assert.Equal(t, []int{2, 3}, deque.ToSlice())
	})
}

func TestRingBufferDeque_Pop(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		x, ok := deque.Pop()

		assert.False(t, ok)
		assert.Zero(t, x)
	})

	t.Run("SingleElement", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1)

		x, ok := deque.Pop()

		assert.True(t, ok)
		assert.Equal(t, 1, x)
		assert.True(t, deque.IsEmpty())
		assert.Equal(t, 0, deque.read)
		assert.Equal(t, 0, deque.write)
	})

	t.Run("ManyElements", func(t *testing.T) {
		deque := NewRingBufferDeque[int](2)
		deque.Append(1, 2, 3)

		x, ok := deque.Pop()

		assert.True(t, ok)
		assert.Equal(t, 3, x)
		assert.Equal(t, []int{1, 2}, deque.ToSlice())
	})
}

func TestRingBufferDeque_WrapAround(t *testing.T) {
	deque := NewRingBufferDeque[int](3)
	deque.Append(1, 2, 3)

	x, ok := deque.Shift()
	assert.True(t, ok)
	assert.Equal(t, 1, x)

	x, ok = deque.Shift()
	assert.True(t, ok)
	assert.Equal(t, 2, x)

	deque.Append(4, 5, 6)

	assert.Equal(t, []int{3, 4, 5, 6}, deque.ToSlice())
}

func TestRingBufferDeque_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		deque := NewRingBufferDeque[string](2)
		deque.Append("a", "b", "c")

		assert.Equal(t, []string{"a", "b", "c"}, slices.Collect(deque.All()))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		deque := NewRingBufferDeque[string](2)
		deque.Append("a", "b", "c")

		var collected []string
		for value := range deque.All() {
			collected = append(collected, value)
			break
		}

		assert.Equal(t, []string{"a"}, collected)
	})
}

func TestRingBufferDeque_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		deque := NewRingBufferDeque[string](2)
		deque.Append("a", "b", "c")

		var indexes []int
		var values []string
		for idx, value := range deque.Enumerate() {
			indexes = append(indexes, idx)
			values = append(values, value)
		}

		assert.Equal(t, []int{0, 1, 2}, indexes)
		assert.Equal(t, []string{"a", "b", "c"}, values)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		deque := NewRingBufferDeque[string](2)
		deque.Append("a", "b", "c")

		var indexes []int
		var values []string
		for idx, value := range deque.Enumerate() {
			indexes = append(indexes, idx)
			values = append(values, value)
			break
		}

		assert.Equal(t, []int{0}, indexes)
		assert.Equal(t, []string{"a"}, values)
	})
}

func TestRingBufferDeque_ToSlice(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		assert.Nil(t, deque.ToSlice())
	})

	t.Run("Populated", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1, 2)

		got := deque.ToSlice()
		got[0] = 100

		assert.Equal(t, []int{1, 2}, deque.ToSlice())
	})
}

func TestRingBufferDeque_String(t *testing.T) {
	deque := NewRingBufferDeque[int](2)
	deque.Append(1, 2)

	assert.Equal(t, "[1 2]", deque.String())
}

func TestRingBufferDeque_Format(t *testing.T) {
	deque := NewRingBufferDeque[int](2)
	deque.Append(1, 2)

	assert.Equal(t, "*deque.RingBufferDeque[int]{size:2, cap:2}", fmt.Sprintf("%#v", deque))
	assert.Equal(t, "*deque.RingBufferDeque[int]{len:2, cap:2} [1 2]", fmt.Sprintf("%+v", deque))
}

func TestRingBufferDeque_MarshalJSON(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)

		data, err := json.Marshal(deque)

		assert.NoError(t, err)
		assert.JSONEq(t, `[]`, string(data))
	})

	t.Run("Populated", func(t *testing.T) {
		deque := NewRingBufferDeque[int](2)
		deque.Append(2, 3)
		deque.Prepend(1)

		data, err := json.Marshal(deque)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})
}

func TestRingBufferDeque_UnmarshalJSON(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1, 2)

		err := json.Unmarshal([]byte(`[8,9]`), deque)

		assert.NoError(t, err)
		assert.Equal(t, []int{8, 9}, deque.ToSlice())
	})

	t.Run("Invalid", func(t *testing.T) {
		deque := NewRingBufferDeque[int](1)
		deque.Append(1)

		err := json.Unmarshal([]byte(`invalid`), deque)

		assert.Error(t, err)
		assert.Equal(t, []int{1}, deque.ToSlice())
	})
}

func TestRingBufferDeque_Clear(t *testing.T) {
	deque := NewRingBufferDeque[int](2)
	deque.Append(1, 2, 3)

	deque.Clear()

	assert.True(t, deque.IsEmpty())
	assert.Nil(t, deque.ToSlice())

	deque.Prepend(5, 6)
	assert.Equal(t, []int{5, 6}, deque.ToSlice())
}
