package list

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/list"
)

// LinkedList is a doubly linked [List].
type LinkedList[T any] = list.DoubleLinkedList[T]

var _ List[any] = &list.DoubleLinkedList[any]{}

// LinkedFactory constructs doubly linked lists.
//
// Linked lists provide O(1) reversal and boundary insertion or removal, while
// indexed operations require O(N) traversal.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New()               O(1)
//	From/FromSeq        O(N)
//	Append(xs...T)      O(len(xs))
//	Insert/Remove       O(N)
//	Get/Set             O(N)
//	Find/Contains       O(N)
//	Reverse             O(1)
//	Clear               O(N)
type LinkedFactory[T any] struct{}

// Linked returns a factory for doubly linked lists.
func Linked[T any]() LinkedFactory[T] {
	return LinkedFactory[T]{}
}

// New creates an empty doubly linked list.
func (LinkedFactory[T]) New() *LinkedList[T] {
	return list.NewDoubleLinkedList[T]()
}

// From creates a linked list containing data.
//
// Values are copied into newly allocated nodes; the source slice is not
// modified or retained.
func (LinkedFactory[T]) From(data []T) *LinkedList[T] {
	return list.NewDoubleLinkedListFromSlice(data)
}

// FromSeq creates a linked list containing values from seq.
func (LinkedFactory[T]) FromSeq(seq iter.Seq[T]) *LinkedList[T] {
	return list.NewDoubleLinkedListFromSeq(seq)
}
