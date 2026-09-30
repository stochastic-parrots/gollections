package list

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/list"
)

// LinkedList is a doubly linked list. Its zero value is ready for use.
type LinkedList[T any] = list.DoubleLinkedList[T]

var _ List[int] = &list.DoubleLinkedList[int]{}

// NewLinked creates an empty doubly linked list. Choose it for O(1)
// reversal and boundary insertion or removal; indexed operations traverse O(N) nodes.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()                 O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Get(idx)              O(N)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(N)
//	Remove(idx)           O(N)
//	Reverse()             O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(1).
func NewLinked[T any]() *LinkedList[T] { return list.NewDoubleLinkedList[T]() }

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
//	Get(idx)              O(N)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(N)
//	Remove(idx)           O(N)
//	Reverse()             O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(len(data)).
func LinkedFrom[T any](data []T) *LinkedList[T] { return list.NewDoubleLinkedListFromSlice(data) }

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
//	Get(idx)              O(N)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1)
//	Appends(xs... T)      O(len(xs))
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(N)
//	Remove(idx)           O(N)
//	Reverse()             O(1)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func LinkedFromSeq[T any](seq iter.Seq[T]) *LinkedList[T] {
	return list.NewDoubleLinkedListFromSeq(seq)
}
