package list

import (
	"github.com/stochastic-parrots/gollections/internal/list"
)

type ArrayList[T any] = *list.ArrayList[T]
type LinkedList[T any] = *list.LinkedList[T]
type DoubleLinkedList[T any] = *list.DoubleLinkedList[T]

var _ List[any] = &list.ArrayList[any]{}
var _ List[any] = &list.LinkedList[any]{}
var _ List[any] = &list.DoubleLinkedList[any]{}

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
func NewArray[T any](size int) ArrayList[T] {
	return list.NewArrayList[T](size)
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
func NewLinked[T any]() LinkedList[T] {
	return list.NewLinkedList[T]()
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
func NewDoubleLinked[T any]() DoubleLinkedList[T] {
	return list.NewDoubleLinkedList[T]()
}
