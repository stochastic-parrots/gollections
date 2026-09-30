package deque

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/deque"
)

// ArrayDeque is a circular-array deque. Its zero value is ready for use.
type ArrayDeque[T any] = deque.RingBufferDeque[T]

var _ Deque[int] = &deque.RingBufferDeque[int]{}

// NewArray creates an empty circular-array deque. Choose it for memory-local
// traversal and amortized O(1) insertion or removal at either end.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Length()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1) Amortized
//	Prepends(xs... T)     O(len(xs)) Amortized
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(capacity).
func NewArray[T any](capacity int) *ArrayDeque[T] { return deque.NewRingBufferDeque[T](capacity) }

// ArrayFrom transfers ownership of data and its backing array; the caller must stop using the slice and every alias afterward. Use [ArrayClone] to preserve the source.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Length()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1) Amortized
//	Prepends(xs... T)     O(len(xs)) Amortized
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(1).
func ArrayFrom[T any](data []T) *ArrayDeque[T] { return deque.NewRingBufferDequeFromSlice(data) }

// ArrayClone makes an independent shallow copy of data without modifying or retaining its backing array.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Length()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1) Amortized
//	Prepends(xs... T)     O(len(xs)) Amortized
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(len(data)).
func ArrayClone[T any](data []T) *ArrayDeque[T] { return deque.NewRingBufferDequeCloneSlice(data) }

// ArrayFromSeq consumes seq once into new storage in iteration order.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Length()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Front()               O(1)
//	Back()                O(1)
//	ToSlice()             O(N)
//	Shift()               O(1)
//	Pop()                 O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	Prepend(x)            O(1) Amortized
//	Prepends(xs... T)     O(len(xs)) Amortized
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func ArrayFromSeq[T any](seq iter.Seq[T]) *ArrayDeque[T] { return deque.NewRingBufferDequeFromSeq(seq) }
