package binary

import (
	"fmt"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/heaps"
	"github.com/stretchr/testify/assert"
)

func TestNewBinaryHeap(t *testing.T) {
	heap := NewBinaryHeap(10, heaps.MinFunc[int]())

	assert.Equal(t, 0, heap.Length())
	assert.True(t, heap.IsEmpty())
	assert.Empty(t, heap.data)
}

func TestNewBinaryHeapFromSlice(t *testing.T) {
	data := []int{10, 5, 8, 2, 7}
	heap := NewBinaryHeapFromSlice(data, heaps.MinFunc[int]())

	assert.Equal(t, 5, heap.Length())

	val, ok := heap.Peek()
	assert.True(t, ok)
	assert.Equal(t, 2, val)
}

func TestNewBinaryHeapCloneSlice(t *testing.T) {
	src := []int{10, 5, 8}
	heap := NewBinaryHeapCloneSlice(src, heaps.MinFunc[int]())

	heap.Pop()
	assert.Equal(t, 3, len(src))
	assert.Equal(t, 2, heap.Length())
}

func TestBinaryHeapPush(t *testing.T) {
	heap := NewBinaryHeap(0, heaps.MinFunc[int]())

	heap.Push(10)
	heap.Push(5)
	val, _ := heap.Peek()
	assert.Equal(t, 5, val)

	heap.Push(2, 20, 1)
	val, _ = heap.Peek()
	assert.Equal(t, 1, val)
	assert.Equal(t, 5, heap.Length())
}

func TestBinaryHeapPushNothing(t *testing.T) {
	heap := NewBinaryHeap(0, heaps.MinFunc[int]())
	heap.Push()
	_, ok := heap.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, heap.Length())
}

func TestBinaryHeapPop(t *testing.T) {
	heap := NewBinaryHeapFromSlice([]int{10, 2, 8, 1}, heaps.MinFunc[int]())

	expected := []int{1, 2, 8, 10}
	for _, exp := range expected {
		val, ok := heap.Pop()
		assert.True(t, ok)
		assert.Equal(t, exp, val)
	}

	val, ok := heap.Pop()
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestBinaryHeapPeek(t *testing.T) {
	heap := NewBinaryHeap(0, heaps.MinFunc[int]())

	_, ok := heap.Peek()
	assert.False(t, ok)

	heap.Push(42)
	val, ok := heap.Peek()
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestBinaryHeapDrain(t *testing.T) {
	items := []int{5, 1, 9, 3}
	heap := NewBinaryHeapFromSlice(items, heaps.MinFunc[int]())

	expected := []int{1, 3, 5, 9}
	count := 0
	for idx, val := range heap.Drain() {
		assert.Equal(t, expected[idx], val)
		count++
	}

	assert.Equal(t, 4, count)
	assert.True(t, heap.IsEmpty())
}

func TestBinaryHeapAll(t *testing.T) {
	items := []int{1, 2, 3}
	heap := NewBinaryHeapFromSlice(items, heaps.MinFunc[int]())

	var collected []int
	for val := range heap.All() {
		collected = append(collected, val)
	}

	assert.Equal(t, 3, len(collected))
	for _, item := range items {
		assert.Contains(t, collected, item)
	}
}

func TestBinaryHeapEnumerate(t *testing.T) {
	heap := NewBinaryHeapFromSlice([]int{10, 20}, heaps.MinFunc[int]())

	for idx, val := range heap.Enumerate() {
		assert.Equal(t, heap.data[idx], val)
	}
}

func TestBinaryHeapString(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		heap := NewBinaryHeap(3, heaps.MinFunc[int]())
		heap.Push(30, 10, 20)
		got := heap.String()
		want := "[10 30 20]"
		assert.Equal(t, want, got)
	})

	t.Run("String Many Elements", func(t *testing.T) {
		heap := NewBinaryHeap(10, heaps.MinFunc[int]())
		heap.Push(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := heap.String()
		want := "[-99 -10 -1 0 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestBinaryHeapFormat(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		heap := NewBinaryHeap(3, heaps.MinFunc[int]())
		heap.Push(30, 10, 20)
		got := fmt.Sprintf("%v", heap)
		want := "[10 30 20]"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose", func(t *testing.T) {
		heap := NewBinaryHeap(10, heaps.MinFunc[int]())
		heap.Push(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%#v", heap)
		want := "*binary.BinaryHeap[int]{size:10, cap:10}"
		assert.Equal(t, want, got)
	})

	t.Run("Verbose + String", func(t *testing.T) {
		heap := NewBinaryHeap(10, heaps.MinFunc[int]())
		heap.Push(30, 10, 20, 40, 1, 0, -1, -10, 0, -99)
		got := fmt.Sprintf("%+v", heap)
		want := "*binary.BinaryHeap[int]{len:10, cap:10} [-99 -10 -1 0 1 ...(+5 more)]"
		assert.Equal(t, want, got)
	})
}

func TestBinaryHeapLargePush(t *testing.T) {
	heap := NewBinaryHeap(0, heaps.MinFunc[int]())
	largeSlice := make([]int, 100)
	for i := range largeSlice {
		largeSlice[i] = 100 - i
	}

	heap.Push(largeSlice...)
	val, _ := heap.Peek()
	assert.Equal(t, 1, val)
	assert.Equal(t, 100, heap.Length())
}
