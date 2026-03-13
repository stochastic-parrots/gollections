package heap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	constructor "github.com/stochastic-parrots/gollections/internal/heap"
)

var _ Heap[any] = &constructor.BinaryHeap[any]{}

type BinaryHeap[T any] = *constructor.BinaryHeap[T]

// NewBinary creates a new empty Binary Heap.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewBinary[T any](capacity int, cmp func(T, T) bool) BinaryHeap[T] {
	return constructor.NewBinaryHeap(capacity, cmp)
}

// BinaryFrom creates a Binary Heap with custom comparator.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func BinaryFrom[T any](data []T, cmp func(T, T) bool) BinaryHeap[T] {
	return constructor.NewBinaryHeapFromSlice(data, cmp)
}

// BinaryClone creates a Binary Heap from a clone of the provided slice with custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func BinaryClone[T any](data []T, cmp func(T, T) bool) BinaryHeap[T] {
	return constructor.NewBinaryHeapCloneSlice(data, cmp)
}

// NewMinBinary creates a new empty Min Binary Heap for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewMinBinary[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return constructor.NewBinaryHeap(capacity, comparator.Min[T]())
}

// NewMaxBinary creates a new empty Max Binary Heap for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func NewMaxBinary[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return constructor.NewBinaryHeap(capacity, comparator.Max[T]())
}

// MinBinaryFrom creates a Min Binary Heap for ordered types.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MinBinaryFrom[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructor.NewBinaryHeapFromSlice(data, comparator.Min[T]())
}

// MinBinaryClone creates a Min Binary Heap from a clone of the provided slice for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Space Complexity    O(N) (Clone)
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MinBinaryClone[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructor.NewBinaryHeapCloneSlice(data, comparator.Min[T]())
}

// MaxBinaryFrom creates a Max Binary Heap.
//
// WARNING: This operation is In-Place and WILL modify the original slice order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MaxBinaryFrom[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructor.NewBinaryHeapFromSlice(data, comparator.Max[T]())
}

// MaxBinaryClone creates a Max Binary Heap from a clone of the provided slice for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N)
//	Space Complexity    O(N) (Clone)
//	Push(xs... T)       O(K log N) or O(N+K)*
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
//	IsEmpty()           O(1)
func MaxBinaryClone[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructor.NewBinaryHeapCloneSlice(data, comparator.Max[T]())
}
