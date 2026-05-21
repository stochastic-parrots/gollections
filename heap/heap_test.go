package heap_test

import (
	"cmp"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/heap"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementHeap(t *testing.T) {
	var _ heap.Heap[int] = heap.NewBinary[int](0, cmp.Less[int])
	var _ heap.Heap[int] = heap.NewMinBinary[int](0)
	var _ heap.Heap[int] = heap.NewMaxBinary[int](0)
	var _ heap.Heap[int] = heap.BinaryClone([]int{1}, cmp.Less[int])
}

func TestNewBinary(t *testing.T) {
	assertHeapBehavior(t, heap.NewBinary[int](0, cmp.Less[int]), []int{1, 2, 3})
}

func TestNewMinBinary(t *testing.T) {
	assertHeapBehavior(t, heap.NewMinBinary[int](0), []int{1, 2, 3})
}

func TestNewMaxBinary(t *testing.T) {
	assertHeapBehavior(t, heap.NewMaxBinary[int](0), []int{3, 2, 1})
}

func TestBinaryFrom(t *testing.T) {
	data := []int{3, 1, 2}
	h := heap.BinaryFrom(data, cmp.Less[int])

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(func(yield func(int) bool) {
		for _, value := range h.Drain() {
			if !yield(value) {
				return
			}
		}
	}))
	assert.NotEqual(t, []int{3, 1, 2}, data)
}

func TestBinaryClone(t *testing.T) {
	data := []int{3, 1, 2}
	h := heap.BinaryClone(data, cmp.Less[int])

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(func(yield func(int) bool) {
		for _, value := range h.Drain() {
			if !yield(value) {
				return
			}
		}
	}))
	assert.Equal(t, []int{3, 1, 2}, data)
}

func assertHeapBehavior(t *testing.T, h heap.Heap[int], expectedDrain []int) {
	t.Helper()

	assert.True(t, h.IsEmpty())

	h.Push(3, 1, 2)

	top, ok := h.Peek()
	assert.True(t, ok)
	assert.Equal(t, expectedDrain[0], top)
	assert.Equal(t, 3, h.Length())
	assert.Len(t, slices.Collect(h.All()), 3)

	old, ok := h.Replace(4)
	assert.True(t, ok)
	assert.Equal(t, expectedDrain[0], old)
	assert.Equal(t, 3, h.Length())

	data, err := json.Marshal(h)
	assert.NoError(t, err)

	var values []int
	err = json.Unmarshal(data, &values)
	assert.NoError(t, err)
	assert.Len(t, values, 3)

	err = json.Unmarshal([]byte(`[8,9]`), h)
	assert.NoError(t, err)
	assert.Equal(t, 2, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())

	h.Push(3, 1, 2)
	assert.Equal(t, expectedDrain, drain(h))
	assert.True(t, h.IsEmpty())
}

func drain(h heap.Heap[int]) []int {
	var values []int
	for _, value := range h.Drain() {
		values = append(values, value)
	}
	return values
}
