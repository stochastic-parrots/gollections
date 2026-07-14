package heap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stochastic-parrots/gollections/internal/heap"
)

// BinaryHeap is a comparator-backed binary [Heap].
type BinaryHeap[T any] = heap.BinaryHeap[T]

var _ Heap[any] = &heap.BinaryHeap[any]{}

// BinaryFactory constructs binary heaps with a fixed priority comparator.
//
// The comparator returns true when its first argument has higher priority than
// its second argument. Binary heaps provide O(1) Peek, O(log N) Pop and Replace,
// and O(N) in-place or cloned construction from a slice.
//
// The zero value is invalid. Create a factory with [Binary] or [OrderedBinary].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	From/Clone          O(N)
//	Push(xs...T)        O(K log N) or O(N+K)
//	Pop/Replace         O(log N)
//	Peek                O(1)
//	Drain               O(N log N)
//	Clear               O(N)
type BinaryFactory[T any] struct {
	hasPriority func(T, T) bool
}

// Binary returns a binary-heap factory using hasPriority to order values.
func Binary[T any](hasPriority func(T, T) bool) BinaryFactory[T] {
	return BinaryFactory[T]{hasPriority: hasPriority}
}

// OrderedBinary returns a binary-heap factory using the natural order of T.
func OrderedBinary[T cmp.Ordered](order Order) BinaryFactory[T] {
	if order == Max {
		return Binary(comparator.Max[T]())
	}
	return Binary(comparator.Min[T]())
}

// New creates an empty binary heap with space preallocated for capacity values.
func (factory BinaryFactory[T]) New(capacity int) *BinaryHeap[T] {
	return heap.NewBinaryHeap(capacity, factory.hasPriority)
}

// From creates a binary heap using data as its backing storage.
//
// WARNING: From reorders data in place. Use [BinaryFactory.Clone] when the
// original slice order must be preserved.
func (factory BinaryFactory[T]) From(data []T) *BinaryHeap[T] {
	return heap.NewBinaryHeapFromSlice(data, factory.hasPriority)
}

// Clone creates a binary heap from a copy of data.
//
// Clone does not modify or retain the provided slice.
func (factory BinaryFactory[T]) Clone(data []T) *BinaryHeap[T] {
	return heap.NewBinaryHeapCloneSlice(data, factory.hasPriority)
}
