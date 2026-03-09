package list

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
)

type LinkedNode[T any] struct {
	value T
	next  *LinkedNode[T]
}

func NewLinkedNode[T any](value T) *LinkedNode[T] {
	return &LinkedNode[T]{value: value, next: nil}
}

func (node LinkedNode[T]) Value() T {
	return node.value
}

func (node LinkedNode[T]) Next() *LinkedNode[T] {
	return node.next
}

func (node LinkedNode[T]) HasNext() bool {
	return node.next != nil
}

func (node *LinkedNode[T]) Append(x T) *LinkedNode[T] {
	new := NewLinkedNode(x)
	node.next = new
	return new
}

// LinkedList represents a singly linked list data structure.
// It consists of nodes where each element points to the next, making it efficient
// for sequential insertion and deletion at the ends, but requiring O(n) for random access.
type LinkedList[T any] struct {
	first, last *LinkedNode[T]
	length      int
}

// NewLinkedList creates an empty LinkedList.
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		first:  nil,
		last:   nil,
		length: 0,
	}
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (l *LinkedList[T]) Length() int {
	return l.length
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (l *LinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

// Get retrieves the value at the specified index by traversing the list from the start.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.Length() {
		var zero T
		return zero, NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		current = current.Next()
	}

	return current.Value(), nil
}

// Set updates the value at the specified index.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *LinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= l.Length() {
		return NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		current = current.Next()
	}

	current.value = x
	return nil
}

func (l *LinkedList[T]) append(x T) {
	new := NewLinkedNode(x)

	if l.IsEmpty() {
		l.first = new
		l.last = new
		l.length++
		return
	} else {
		l.last.next = new
		l.last = new
	}
	l.length++
}

// Append adds one or more elements to the end of the list.
//
// Complexity: O(k) where k is the number of elements provided.
func (l *LinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		l.append(x)
	}
}

// Reverse inverts the order of the elements in the list in-place.
//
// Complexity: O(n).
// Note: Unlike DoubleLinkedList, this requires a full traversal to update pointers.
func (l *LinkedList[T]) Reverse() {
	if l.Length() <= 1 {
		return
	}

	var previous *LinkedNode[T]
	current := l.first

	for current != nil {
		next := current.next
		current.next = previous
		previous = current
		current = next
	}

	tmp := l.first
	l.first = previous
	l.last = tmp
}

// All returns a sequence that yields elements from the first to the last node.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *LinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for current := l.first; current != nil; current = current.next {
			if !yield(current.value) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the index and value of each element.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *LinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for current, index := l.first, 0; current != nil; current = current.next {
			if !yield(index, current.value) {
				return
			}
			index++
		}
	}
}

// Format implements the fmt.Formatter interface for custom string formatting.
func (l *LinkedList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
func (l *LinkedList[T]) String() string {
	return fmt.Sprint(l)
}
