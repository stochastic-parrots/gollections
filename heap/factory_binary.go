package heap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/heap"
	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

// BinaryHeap is a comparator-backed binary [Heap]. Its zero value is invalid;
// construct one with [Binary] or [OrderedBinary].
type BinaryHeap[T any] = heap.BinaryHeap[T]

var _ Heap[any] = &heap.BinaryHeap[any]{}

// BinaryFactory constructs binary heaps with a fixed priority ordering.
//
// The comparator returns true when its first argument has higher priority than
// its second argument. Binary heaps provide O(1) Peek, O(log N) Pop and Replace,
// and O(N) in-place or cloned construction from a slice.
// The comparator must be non-nil, define a strict weak ordering, and remain
// stable for the lifetime of every heap created by the factory.
//
// The zero value is invalid. Create a factory with [Binary] or [OrderedBinary].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	From/Clone          O(N)
//	Push(x)             O(log N)
//	Pushes(xs...T)      O(len(xs) log N) or O(N + len(xs))
//	Pop/Replace         O(log N)
//	Peek                O(1)
//	Drain               O(N log N)
//	Clear               O(N)
//
// Pushes uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses O(len(xs) log N) insertion otherwise.
type BinaryFactory[T any] struct {
	hasPriority func(T, T) bool
}

// Binary returns a binary-heap factory using hasPriority to order values.
// It panics if hasPriority is nil.
func Binary[T any](hasPriority func(T, T) bool) BinaryFactory[T] {
	if hasPriority == nil {
		panic("heap: nil priority comparator")
	}
	return BinaryFactory[T]{hasPriority: hasPriority}
}

// OrderedBinary returns a binary-heap factory using the natural order of T.
func OrderedBinary[T cmp.Ordered](order Order) BinaryFactory[T] {
	if order == Max {
		return Binary(ordering.Max[T]())
	}
	return Binary(ordering.Min[T]())
}

// New creates an empty binary heap with space preallocated for capacity values.
func (factory BinaryFactory[T]) New(capacity int) *BinaryHeap[T] {
	return heap.NewBinaryHeap(capacity, factory.hasPriority)
}

// From creates a binary heap using data as its backing storage.
//
// WARNING: From reorders data in place and transfers ownership of its backing
// storage. The caller must not use data or aliases of its backing array after
// this call. Use [BinaryFactory.Clone] to preserve the source slice.
func (factory BinaryFactory[T]) From(data []T) *BinaryHeap[T] {
	return heap.NewBinaryHeapFromSlice(data, factory.hasPriority)
}

// Clone creates a binary heap from a copy of data.
//
// Clone does not modify or retain the provided slice.
func (factory BinaryFactory[T]) Clone(data []T) *BinaryHeap[T] {
	return heap.NewBinaryHeapCloneSlice(data, factory.hasPriority)
}
