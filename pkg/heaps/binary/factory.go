package binary

import (
	"cmp"

	internal "github.com/stochastic-parrots/gollections/internal/heaps"
	constructors "github.com/stochastic-parrots/gollections/internal/heaps/binary"
	"github.com/stochastic-parrots/gollections/pkg/heaps"
)

var _ heaps.Heap[any] = &constructors.BinaryHeap[any]{}

type BinaryHeap[T any] = *constructors.BinaryHeap[T]

// New creates and returns a new empty Binary Heap with a custom comparator.
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

// Min creates and returns a new empty Binary Min-Heap for ordered types.
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

// MinFromSlice creates a Min-Heap from an existing slice for ordered types.
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

// MinCloneSlice creates a Min-Heap from a Clone of the provided slice for ordered types.
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

// Max creates and returns a new empty Binary Max-Heap for ordered types.
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

// MaxFromSlice creates a Max-Heap from an existing slice for ordered types.
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

// MaxCloneSlice creates a Max-Heap from a Clone of the provided slice for ordered types.
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
