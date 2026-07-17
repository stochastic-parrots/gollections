package deque

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/deque"
)

// ArrayDeque is a circular-array backed [Deque]. Its zero value is ready for use.
type ArrayDeque[T any] = deque.RingBufferDeque[T]

var _ Deque[any] = &deque.RingBufferDeque[any]{}

// ArrayFactory constructs circular-array backed deques.
//
// Array deques provide amortized O(1) insertion and removal at both ends with
// cache-friendly traversal.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	From(data)          O(1)
//	Clone/FromSeq       O(N)
//	Append/Prepend      O(1) amortized
//	Appends/Prepends    O(len(xs)) amortized
//	Shift/Pop           O(1)
//	Front/Back          O(1)
//	Clear               O(N)
type ArrayFactory[T any] struct{}

// Array returns a factory for circular-array backed deques.
func Array[T any]() ArrayFactory[T] {
	return ArrayFactory[T]{}
}

// New creates an empty array deque with the requested initial capacity.
func (ArrayFactory[T]) New(capacity int) *ArrayDeque[T] {
	return deque.NewRingBufferDeque[T](capacity)
}

// From creates an array deque using data as its backing storage.
//
// Values initially appear in front-to-back order. The caller transfers
// ownership of data to the returned deque and must not use the slice afterward.
// Use [ArrayFactory.Clone] to preserve the source slice.
func (ArrayFactory[T]) From(data []T) *ArrayDeque[T] {
	return deque.NewRingBufferDequeFromSlice(data)
}

// Clone creates an array deque from a shallow copy of data.
//
// Clone does not modify or retain the provided slice.
func (ArrayFactory[T]) Clone(data []T) *ArrayDeque[T] {
	return deque.NewRingBufferDequeCloneSlice(data)
}

// FromSeq collects seq into a new array deque in front-to-back order.
func (ArrayFactory[T]) FromSeq(seq iter.Seq[T]) *ArrayDeque[T] {
	return deque.NewRingBufferDequeFromSeq(seq)
}
