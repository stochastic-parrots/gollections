package list

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
)

type DoubleLinkedNode[T any] struct {
	previous *DoubleLinkedNode[T]
	value    T
	next     *DoubleLinkedNode[T]
}

func NewDoubleLinkedNode[T any](value T) *DoubleLinkedNode[T] {
	return &DoubleLinkedNode[T]{
		previous: nil,
		value:    value,
		next:     nil,
	}
}

func (node *DoubleLinkedNode[T]) Unlink() {
	previous := node.previous
	next := node.next

	if previous != nil {
		previous.next = next
	}
	if next != nil {
		next.previous = previous
	}

	node.previous = nil
	node.next = nil
}

func (node DoubleLinkedNode[T]) Value() T {
	return node.value
}

func (node DoubleLinkedNode[T]) Next() *DoubleLinkedNode[T] {
	return node.next
}

func (node *DoubleLinkedNode[T]) Append(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.next = new
	new.previous = node
	return new
}

func (node DoubleLinkedNode[T]) HasNext() bool {
	return node.next != nil
}

func (node DoubleLinkedNode[T]) Previous() *DoubleLinkedNode[T] {
	return node.previous
}

func (node *DoubleLinkedNode[T]) PreAppend(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.previous = new
	new.next = node
	return new
}

func (node DoubleLinkedNode[T]) HasPrevious() bool {
	return node.previous != nil
}

// DoubleLinkedList represents a doubly linked list data structure.
// It stores elements in a series of nodes where each node points to both
// its predecessor and its successor, allowing for efficient bi-directional traversal.
type DoubleLinkedList[T any] struct {
	first, last *DoubleLinkedNode[T]
	length      int
	reversed    bool
}

// NewDoubleLinkedList creates an empty DoubleLinkedList ready for use.
func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		first:    nil,
		last:     nil,
		length:   0,
		reversed: false,
	}
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) Length() int {
	return l.length
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

// Get retrieves the value at the specified index.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *DoubleLinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.Length() {
		var zero T
		return zero, NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		if l.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	return current.Value(), nil
}

// Set updates the value at the specified index.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *DoubleLinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= l.Length() {
		return NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		if l.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	current.value = x
	return nil
}

func (l *DoubleLinkedList[T]) append(x T) {
	if l.IsEmpty() {
		l.first = NewDoubleLinkedNode(x)
		l.last = l.first
		l.length++
		return
	}

	if !l.reversed {
		new := NewDoubleLinkedNode(x)
		l.last.next = new
		new.previous = l.last
		l.last = new
		l.length++
		return
	}

	new := NewDoubleLinkedNode(x)
	l.last.previous = new
	new.next = l.last
	l.last = new
	l.length++
}

// Append inserts one or more elements at the end of the list.
//
// Complexity: O(k) where k is the number of elements provided.
func (l *DoubleLinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		l.append(x)
	}
}

// Reverse inverts the logical order of the list in O(1) time.
// Instead of physically rearranging nodes, it toggles an internal flag
// that affects traversal and insertion logic.
func (l *DoubleLinkedList[T]) Reverse() {
	temp := l.first
	l.first = l.last
	l.last = temp
	l.reversed = !l.reversed
}

// All returns a sequence that yields elements in their logical order.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *DoubleLinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := l.first
		for current != nil {
			if !yield(current.value) {
				return
			}

			if !l.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

// Enumerate returns a sequence that yields the index and value of elements
// in their logical order.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *DoubleLinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := l.first
		for index := 0; current != nil; index++ {
			if !yield(index, current.value) {
				return
			}

			if !l.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
func (l *DoubleLinkedList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
func (l *DoubleLinkedList[T]) String() string {
	return fmt.Sprint(l)
}
