package heaps

import (
	"cmp"

	internal "github.com/stochastic-parrots/gollections/internal/heaps"
	constructors "github.com/stochastic-parrots/gollections/internal/heaps/binary"
)

var _ Heap[any] = &constructors.BinaryHeap[any]{}

type BinaryHeap[T any] = *constructors.BinaryHeap[T]

// New creates a new empty Binary Heap (Default Implementation).
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
func New[T any](capacity int, cmp func(T, T) bool) BinaryHeap[T] {
	return constructors.NewBinaryHeap(capacity, cmp)
}

// Min creates a new empty Binary (Default Implementation) Min-Heap for ordered types.
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
func Min[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return constructors.NewBinaryHeap(capacity, internal.MinFunc[T]())
}

// Max creates a new empty Binary (Default Implementation) Max-Heap for ordered types.
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
func Max[T cmp.Ordered](capacity int) BinaryHeap[T] {
	return constructors.NewBinaryHeap(capacity, internal.MaxFunc[T]())
}

// MinFromSlice creates a Binary (Default Implementation) Min-Heap for ordered types.
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
func MinFromSlice[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructors.NewBinaryHeapFromSlice(data, internal.MinFunc[T]())
}

// MinCloneSlice creates a Binary (Default Implementation) Min-Heap for ordered types
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
func MinCloneSlice[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructors.NewBinaryHeapCloneSlice(data, internal.MinFunc[T]())
}

// MaxFromSlice creates a Binary (Default Implementation) Min-Heap.
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
func MaxFromSlice[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructors.NewBinaryHeapFromSlice(data, internal.MaxFunc[T]())
}

// MaxCloneSlice creates a Binary Min-Heap from a clone of the provided slice for ordered types
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
func MaxCloneSlice[T cmp.Ordered](data []T) BinaryHeap[T] {
	return constructors.NewBinaryHeapCloneSlice(data, internal.MaxFunc[T]())
}
