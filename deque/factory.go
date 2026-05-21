package deque

import (
	"github.com/stochastic-parrots/gollections/internal/deque"
)

// ArrayDeque is a circular-array backed [Deque].
type ArrayDeque[T any] = *deque.RingBufferDeque[T]

// LinkedDeque is a doubly linked [Deque].
type LinkedDeque[T any] = *deque.DoubleLinkedDeque[T]

var _ Deque[any] = &deque.RingBufferDeque[any]{}
var _ Deque[any] = &deque.DoubleLinkedDeque[any]{}

// NewArray creates an empty deque backed by a circular array.
//
// It provides amortized O(1) insertion and removal at both ends while keeping
// elements in contiguous memory for cache-friendly traversal. The size parameter
// pre-allocates the initial buffer capacity.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Append(xs...T)      O(len(xs)) Amortized
//	Prepend(xs...T)     O(len(xs)) Amortized
//	Shift()             O(1)
//	Pop()               O(1)
//	Front()             O(1)
//	Back()              O(1)
//	Clear()             O(N)
func NewArray[T any](size int) ArrayDeque[T] {
	return deque.NewRingBufferDeque[T](size)
}

// NewLinked creates an empty deque backed by a doubly linked list.
//
// It provides O(1) insertion and removal at both ends without moving existing
// elements. It is useful when growth is unpredictable or when avoiding backing
// array reallocations matters more than memory locality.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Append(xs...T)      O(len(xs))
//	Prepend(xs...T)     O(len(xs))
//	Shift()             O(1)
//	Pop()               O(1)
//	Front()             O(1)
//	Back()              O(1)
//	Clear()             O(N)
func NewLinked[T any]() LinkedDeque[T] {
	return deque.NewDoubleLinkedDeque[T]()
}
