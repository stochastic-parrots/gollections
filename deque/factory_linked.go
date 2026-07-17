package deque

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/deque"
)

// LinkedDeque is a doubly linked [Deque]. Its zero value is ready for use.
type LinkedDeque[T any] = deque.DoubleLinkedDeque[T]

var _ Deque[any] = &deque.DoubleLinkedDeque[any]{}

// LinkedFactory constructs doubly linked deques.
//
// Linked deques provide O(1) insertion and removal at both ends without moving
// existing elements or reallocating a backing array.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New()               O(1)
//	From/FromSeq        O(N)
//	Append/Prepend      O(1)
//	Appends/Prepends    O(len(xs))
//	Shift/Pop           O(1)
//	Front/Back          O(1)
//	Clear               O(N)
type LinkedFactory[T any] struct{}

// Linked returns a factory for doubly linked deques.
func Linked[T any]() LinkedFactory[T] {
	return LinkedFactory[T]{}
}

// New creates an empty doubly linked deque.
func (LinkedFactory[T]) New() *LinkedDeque[T] {
	return deque.NewDoubleLinkedDeque[T]()
}

// From creates a linked deque containing data in front-to-back order.
//
// Values are copied into newly allocated nodes; the source slice is not
// modified or retained.
func (LinkedFactory[T]) From(data []T) *LinkedDeque[T] {
	return deque.NewDoubleLinkedDequeFromSlice(data)
}

// FromSeq creates a linked deque containing values from seq in front-to-back order.
func (LinkedFactory[T]) FromSeq(seq iter.Seq[T]) *LinkedDeque[T] {
	return deque.NewDoubleLinkedDequeFromSeq(seq)
}
