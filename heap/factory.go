package heap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stochastic-parrots/gollections/internal/heap"
)

var _ Heap[any] = &heap.BinaryHeap[any]{}

// BinaryHeap is a comparator-backed binary [Heap].
type BinaryHeap[T any] = *heap.BinaryHeap[T]

// NewBinary creates an empty binary heap with a custom priority comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewBinary[T any](capacity int, cmp func(T, T) bool) BinaryHeap[T] {
	return heap.NewBinaryHeap(capacity, cmp)
}

// BinaryFrom creates a binary heap with a custom priority comparator.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func BinaryFrom[T any](data []T, cmp func(T, T) bool) BinaryHeap[T] {
	return heap.NewBinaryHeapFromSlice(data, cmp)
}

// BinaryClone creates a binary heap from a clone of the provided slice with a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Extra Space         O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func BinaryClone[T any](data []T, cmp func(T, T) bool) BinaryHeap[T] {
	return heap.NewBinaryHeapCloneSlice(data, cmp)
}

// NewMinBinary creates an empty min binary heap for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewMinBinary[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return heap.NewBinaryHeap(capacity, comparator.Min[T]())
}

// NewMaxBinary creates an empty max binary heap for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewMaxBinary[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return heap.NewBinaryHeap(capacity, comparator.Max[T]())
}

// MinBinaryFrom creates a min binary heap for ordered types.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MinBinaryFrom[T cmp.Ordered](data []T) BinaryHeap[T] {
	return heap.NewBinaryHeapFromSlice(data, comparator.Min[T]())
}

// MinBinaryClone creates a min binary heap from a clone of the provided slice for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Extra Space         O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MinBinaryClone[T cmp.Ordered](data []T) BinaryHeap[T] {
	return heap.NewBinaryHeapCloneSlice(data, comparator.Min[T]())
}

// MaxBinaryFrom creates a max binary heap.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MaxBinaryFrom[T cmp.Ordered](data []T) BinaryHeap[T] {
	return heap.NewBinaryHeapFromSlice(data, comparator.Max[T]())
}

// MaxBinaryClone creates a max binary heap from a clone of the provided slice for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Extra Space         O(N)
//	Push(xs... T)       O(K log N) or O(N+K)
//	Pop()               O(log N)
//	Replace(x)          O(log N)
//	Peek()              O(1)
//	Drain()             O(N log N)
//	Clear()             O(N)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MaxBinaryClone[T cmp.Ordered](data []T) BinaryHeap[T] {
	return heap.NewBinaryHeapCloneSlice(data, comparator.Max[T]())
}
