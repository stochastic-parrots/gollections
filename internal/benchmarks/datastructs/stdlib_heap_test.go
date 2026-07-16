package datastructs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdLibHeap_Push(t *testing.T) {
	heap := NewStdLibHeap(0)

	heap.Push(3)
	heap.Push(1)
	heap.Push(2)

	value, ok := heap.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, value)
	assert.Equal(t, 3, heap.Length())
}

func TestStdLibHeap_Pushes(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		heap := NewStdLibHeap(0)

		heap.Pushes()

		assert.Zero(t, heap.Length())
	})

	t.Run("Single", func(t *testing.T) {
		heap := NewStdLibHeap(0)

		heap.Pushes(2)

		value, ok := heap.Peek()
		assert.True(t, ok)
		assert.Equal(t, 2, value)
	})

	t.Run("Many", func(t *testing.T) {
		heap := NewStdLibHeap(0)
		heap.Push(3)

		heap.Pushes(4, 1, 2)

		values := make([]int, 0, heap.Length())
		for heap.Length() > 0 {
			value, ok := heap.Pop()
			assert.True(t, ok)
			values = append(values, value)
		}
		assert.Equal(t, []int{1, 2, 3, 4}, values)
	})
}
