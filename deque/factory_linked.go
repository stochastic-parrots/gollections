package deque

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/deque"
)

// LinkedDeque is a doubly linked deque. Its zero value is ready for use.
type LinkedDeque[T any] = deque.DoubleLinkedDeque[T]

var _ Deque[int] = &deque.DoubleLinkedDeque[int]{}

// NewLinked creates an empty doubly linked deque. Choose it for O(1) updates
// at either end without moving existing values or reallocating a backing array.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()                 O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1)
//	Prepends(xs... T)     O(len(xs))
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(1).
func NewLinked[T any]() *LinkedDeque[T] { return deque.NewDoubleLinkedDeque[T]() }

// LinkedFrom copies values from data into new nodes without retaining the source slice.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()                 O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1)
//	Prepends(xs... T)     O(len(xs))
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(len(data)).
func LinkedFrom[T any](data []T) *LinkedDeque[T] { return deque.NewDoubleLinkedDequeFromSlice(data) }

// LinkedFromSeq consumes seq once into new storage in iteration order.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()                 O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1)
//	Prepends(xs... T)     O(len(xs))
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func LinkedFromSeq[T any](seq iter.Seq[T]) *LinkedDeque[T] {
	return deque.NewDoubleLinkedDequeFromSeq(seq)
}
