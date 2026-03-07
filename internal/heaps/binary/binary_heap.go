package binary

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
)

type BinaryHeap[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewBinaryHeap[T any](capacity int, less func(a, b T) bool) *BinaryHeap[T] {
	data := make([]T, 0, capacity)
	return &BinaryHeap[T]{data, less}
}
func NewBinaryHeapFromSlice[T any](data []T, less func(a, b T) bool) *BinaryHeap[T] {
	heap := &BinaryHeap[T]{data, less}
	heap.heapify()
	return heap
}

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

func (heap *BinaryHeap[T]) Peek() (T, bool) {
	if len(heap.data) == 0 {
		var zero T
		return zero, false
	}

	return heap.data[0], true
}

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

func (heap *BinaryHeap[T]) IsEmpty() bool {
	return len(heap.data) == 0
}

func (heap *BinaryHeap[T]) Length() int {
	return len(heap.data)
}

func (heap *BinaryHeap[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range heap.data {
			if !yield(value) {
				return
			}
		}
	}
}

func (heap *BinaryHeap[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range heap.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

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

func (heap *BinaryHeap[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, heap, cap(heap.data))
}

func (heap *BinaryHeap[T]) String() string {
	return fmt.Sprint(heap)
}
