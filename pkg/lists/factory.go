package lists

import (
	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stochastic-parrots/gollections/internal/lists/array"
	"github.com/stochastic-parrots/gollections/internal/lists/doublelinked"
	"github.com/stochastic-parrots/gollections/internal/lists/linked"
)

type List[T any] = lists.List[T]

// NewArrayList creates and returns a new empty array list.
//
// It is a wrapper for golang slices, returns the list as a List[T] interface.
//
// Performance Summary (Time Complexity):
//
//	Operation                    Time Complexity
//	--------------------------   ---------------
//	Append(xs...T)               O(xs) Amortized
//	Get(index int)               O(1)
//	Set(index int, x T)          O(1)
//	Reverse()                    O(N)
//	Length()                     O(1)
//	IsEmpty()                    O(1)
//	String()                     O(N)
//	Iterator()                   O(1)
func NewArrayList[T any](size int) List[T] {
	return array.NewArrayList[T](size)
}

// NewLinkedList creates and returns a new empty linked list.
//
// It returns the list as a List[T] interface.
//
// Performance Summary (Time Complexity):
//
//	Operation                    Time Complexity
//	--------------------------   ---------------
//	Append(xs...T)               O(xs)
//	Get(index int)               O(N)
//	Set(index int, x T)          O(N)
//	Reverse()                    O(N)
//	Length()                     O(1)
//	IsEmpty()                    O(1)
//	String()                     O(N)
//	Iterator()                   O(1)
func NewLinkedList[T any]() List[T] {
	return linked.NewLinkedList[T]()
}

// NewDoubleLinkedList creates and returns a new empty double linked list.
//
// It returns the list as a List[T] interface.
//
// Performance Summary (Time Complexity):
//
//	Operation                    Time Complexity
//	--------------------------   ---------------
//	Append(xs...T)               O(xs)
//	Get(index int)               O(N)
//	Set(index int, x T)          O(N)
//	Reverse()                    O(1)
//	Length()                     O(1)
//	IsEmpty()                    O(1)
//	String()                     O(N)
//	Iterator()                   O(1)
func NewDoubleLinkedList[T any]() List[T] {
	return doublelinked.NewDoubleLinkedList[T]()
}
