package list

import (
	"github.com/stochastic-parrots/gollections/internal/list"
)

type ArrayList[T any] = *list.ArrayList[T]
type LinkedList[T any] = *list.DoubleLinkedList[T]

var _ List[any] = &list.ArrayList[any]{}
var _ List[any] = &list.DoubleLinkedList[any]{}
var _ List[any] = &list.LinkedList[any]{}

// NewArray creates a new ArrayList backed by a contiguous slice.
//
// It provides constant-time O(1) random access and is highly cache-efficient
// due to contiguous memory allocation. It is the preferred general-purpose
// list for most use cases where random access is frequent.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Append(xs...T)      O(len(xs)) Amortized
//	Insert(idx, x)      O(N)
//	Remove(idx)         O(N)
//	Get(index)          O(1)
//	Set(index, x)       O(1)
//	Find(x)             O(N)
//	Contains(x)         O(N)
//	Reverse()           O(N)
//	Clear()             O(N)
func NewArray[T any](size int) ArrayList[T] {
	return list.NewArrayList[T](size)
}

// NewLinked creates a new doubly linked list.
//
// Unlike the array-based list, the linked list is a pointer-based structure
// where each element points to its successor and predecessor. This makes it
// ideal for scenarios with frequent insertions and removals at the ends
// of the list, as these operations do not require memory shifting.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Append(xs...T)      O(len(xs))
//	Insert(idx, x)      O(N)
//	Remove(idx)         O(N)
//	Get(index)          O(N)
//	Set(index, x)       O(N)
//	Find(x)             O(N)
//	Contains(x)         O(N)
//	Reverse()           O(1)
//	Clear()             O(N)
func NewLinked[T any]() LinkedList[T] {
	return list.NewDoubleLinkedList[T]()
}
