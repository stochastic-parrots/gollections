package lists

import (
	"github.com/stochastic-parrots/gollections/internal/lists/array"
	"github.com/stochastic-parrots/gollections/internal/lists/doublelinked"
	"github.com/stochastic-parrots/gollections/internal/lists/linked"
)

type ArrayList[T any] = *array.ArrayList[T]
type LinkedList[T any] = *linked.LinkedList[T]
type DoubleLinkedList[T any] = *doublelinked.DoubleLinkedList[T]

var _ List[any] = &array.ArrayList[any]{}
var _ List[any] = &linked.LinkedList[any]{}
var _ List[any] = &doublelinked.DoubleLinkedList[any]{}

// NewArrayList creates and returns a new empty array list.
//
// It is a wrapper for golang slices.
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
func NewArrayList[T any](size int) ArrayList[T] {
	return array.NewArrayList[T](size)
}

// NewLinkedList creates and returns a new empty linked list.
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
func NewLinkedList[T any]() LinkedList[T] {
	return linked.NewLinkedList[T]()
}

// NewDoubleLinkedList creates and returns a new empty double linked list.
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
func NewDoubleLinkedList[T any]() DoubleLinkedList[T] {
	return doublelinked.NewDoubleLinkedList[T]()
}
