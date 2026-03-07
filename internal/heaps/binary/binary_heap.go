package binary

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
)

// BinaryHeap represents a priority queue implemented as a binary tree in a contiguous slice.
// It is generic over type T and requires a 'less' function to define priority.
type BinaryHeap[T any] struct {
	data []T
	less func(a, b T) bool
}

// NewBinaryHeap creates an empty BinaryHeap with the specified initial capacity.
// The 'less' function determines if 'a' has higher priority than 'b' (e.g., a < b for a Min-Heap).
func NewBinaryHeap[T any](capacity int, less func(a, b T) bool) *BinaryHeap[T] {
	data := make([]T, 0, capacity)
	return &BinaryHeap[T]{data, less}
}

// NewBinaryHeapFromSlice creates a BinaryHeap using the provided slice directly as its storage.
// Warning: This method modifies the original slice in-place to satisfy the heap property (O(n)).
func NewBinaryHeapFromSlice[T any](data []T, less func(a, b T) bool) *BinaryHeap[T] {
	heap := &BinaryHeap[T]{data, less}
	heap.heapify()
	return heap
}

// NewBinaryHeapCloneSlice creates a BinaryHeap by cloning the provided slice.
// The original slice remains unchanged. The new heap is built in O(n) time.
func NewBinaryHeapCloneSlice[T any](src []T, less func(a, b T) bool) *BinaryHeap[T] {
	data := make([]T, len(src))
	copy(data, src)
	heap := &BinaryHeap[T]{data, less}
	heap.heapify()
	return heap
}

func (heap *BinaryHeap[T]) heapify() {
	n := len(heap.data)
	if n < 2 {
		return
	}

	for i := (n / 2) - 1; i >= 0; i-- {
		heap.fixdown(i)
	}
}

func (heap *BinaryHeap[T]) fixdown(idx int) {
	length, i := len(heap.data), idx
	target := heap.data[i]

	for {
		left := 2*i + 1 // left children
		if left >= length || left < 0 {
			break
		}

		right := 2*i + 2 // right children
		smallest := left

		if right < length && heap.less(heap.data[right], heap.data[smallest]) {
			smallest = right
		}

		if !heap.less(heap.data[smallest], target) {
			break
		}

		heap.data[i] = heap.data[smallest]
		i = smallest
	}

	heap.data[i] = target
}

func (heap *BinaryHeap[T]) fixup(idx int) {
	i := idx
	target := heap.data[i]

	for i > 0 {
		parent := (i - 1) / 2

		if !heap.less(target, heap.data[parent]) {
			break
		}

		heap.data[i] = heap.data[parent]
		i = parent
	}

	heap.data[i] = target
}

// Pop removes and returns the element with the highest priority.
//
// Complexity: O(log n).
// Returns the zero value of T and false if the heap is empty.
func (heap *BinaryHeap[T]) Pop() (T, bool) {
	var zero T
	length := len(heap.data)

	if length == 0 {
		return zero, false
	}

	last := length - 1
	root := heap.data[0]
	heap.data[0] = heap.data[last]
	heap.data[last] = zero
	heap.data = heap.data[:last]
	length--

	if length >= 2 {
		heap.fixdown(0)
	}
	return root, true
}

// Peek returns the element with the highest priority without removing it.
//
// Complexity: O(1).
// Returns the zero value of T and false if the heap is empty.
func (heap *BinaryHeap[T]) Peek() (T, bool) {
	if len(heap.data) == 0 {
		var zero T
		return zero, false
	}

	return heap.data[0], true
}

// Push inserts one or more elements into the heap.
//
// Complexity: O(log n) per element inserted.
// If multiple elements are provided and the heap is small or empty,
// it uses an optimized O(n) heapify approach.
func (heap *BinaryHeap[T]) Push(xs ...T) {
	k := len(xs)
	if k == 0 {
		return
	}
	n := len(heap.data)

	if n == 0 || (k > 64 && k > n) {
		heap.data = append(heap.data, xs...)
		heap.heapify()
		return
	}

	for _, x := range xs {
		heap.data = append(heap.data, x)
		heap.fixup(len(heap.data) - 1)
	}
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1).
func (heap *BinaryHeap[T]) IsEmpty() bool {
	return len(heap.data) == 0
}

// Length returns the current number of elements in the heap.
//
// Complexity: O(1).
func (heap *BinaryHeap[T]) Length() int {
	return len(heap.data)
}

// Iterator returns a sequence that yields elements in the underlying slice order.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (heap *BinaryHeap[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range heap.data {
			if !yield(value) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the index and value of elements
// as stored in the underlying slice.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (heap *BinaryHeap[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range heap.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

// Drain returns a destructive iterator that removes and yields elements
// in priority order (highest to lowest).
//
// Complexity: O(n log n) for a full traversal.
func (heap *BinaryHeap[T]) Drain() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		index := 0
		for {
			val, ok := heap.Pop()
			if !ok {
				return
			}
			if !yield(index, val) {
				return
			}
			index++
		}
	}
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
func (heap *BinaryHeap[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, heap, cap(heap.data))
}

// String returns a string representation of the heap.
func (heap *BinaryHeap[T]) String() string {
	return fmt.Sprint(heap)
}
