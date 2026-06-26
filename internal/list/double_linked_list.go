package list

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
	"github.com/stochastic-parrots/gollections/internal/shared/node"
)

// DoubleLinkedList represents a doubly linked list data structure.
// It stores elements in a series of nodes where each node points to both
// its predecessor and its successor, allowing for efficient bi-directional traversal.
type DoubleLinkedList[T any] struct {
	first, last *node.DoubleLinkedNode[T]
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

// forward returns the logically "next" node based on the list's orientation.
// If the list is reversed, it moves toward the 'previous' pointer.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) forward(n *node.DoubleLinkedNode[T]) *node.DoubleLinkedNode[T] {
	if l.reversed {
		return n.Previous
	}
	return n.Next
}

// backward returns the logically "previous" node based on the list's orientation.
// If the list is reversed, it moves toward the 'next' pointer.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) backward(n *node.DoubleLinkedNode[T]) *node.DoubleLinkedNode[T] {
	if l.reversed {
		return n.Next
	}
	return n.Previous
}

// get traverses the list to find the node at the specific index.
// It optimizes search time by starting from either the head or the tail,
// depending on which is closer to the requested index.
//
// Complexity: O(n/2) which simplifies to O(n).
func (l *DoubleLinkedList[T]) get(idx int) *node.DoubleLinkedNode[T] {
	size := l.Length()
	var current *node.DoubleLinkedNode[T]

	if idx < size/2 {
		current = l.first
		for range idx {
			current = l.forward(current)
		}
	} else {
		current = l.last
		for i := 0; i < (size - 1 - idx); i++ {
			current = l.backward(current)
		}
	}
	return current
}

// Get retrieves the value at the specified index.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *DoubleLinkedList[T]) Get(idx int) (T, error) {
	if idx < 0 || idx >= l.Length() {
		var zero T
		return zero, NewIndexOutOfBoundError(idx, l.Length()-1)
	}

	return l.get(idx).Value, nil
}

// Find locates the index of an element using a linear search.
//
// Complexity: O(n).
func (l *DoubleLinkedList[T]) Find(x T, cmp func(a, b T) int) (idx int, ok bool) {
	if l.IsEmpty() {
		return -1, false
	}

	for idx, value := range l.Enumerate() {
		if cmp(x, value) == 0 {
			return idx, true
		}
	}

	return -1, false
}

// Contains returns true if the element exists in the list according to cmp.
//
// Complexity: O(n).
func (l *DoubleLinkedList[T]) Contains(x T, cmp func(a, b T) int) bool {
	if l.IsEmpty() {
		return false
	}

	for value := range l.All() {
		if cmp(x, value) == 0 {
			return true
		}
	}

	return false
}

// Set updates the value at the specified index.
//
// Complexity: O(n).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *DoubleLinkedList[T]) Set(idx int, x T) error {
	if idx < 0 || idx >= l.Length() {
		return NewIndexOutOfBoundError(idx, l.Length()-1)
	}

	l.get(idx).Value = x
	return nil
}

// append is the internal implementation for adding a value to the logical end of the list.
// It handles pointer updates for both standard and reversed list states.
func (l *DoubleLinkedList[T]) append(x T) {
	if l.IsEmpty() {
		l.first = node.NewDoubleLinkedNode(x)
		l.last = l.first
		l.length++
		return
	}

	if !l.reversed {
		new := node.NewDoubleLinkedNode(x)
		l.last.Next = new
		new.Previous = l.last
		l.last = new
		l.length++
		return
	}

	new := node.NewDoubleLinkedNode(x)
	l.last.Previous = new
	new.Next = l.last
	l.last = new
	l.length++
}

// prepend is the internal implementation for adding a value to the logical start of the list.
// It handles pointer updates for both standard and reversed list states.
func (l *DoubleLinkedList[T]) prepend(x T) {
	if l.IsEmpty() {
		l.append(x)
		return
	}

	new := node.NewDoubleLinkedNode(x)
	if !l.reversed {
		new.Next = l.first
		l.first.Previous = new
	} else {
		new.Previous = l.first
		l.first.Next = new
	}
	l.first = new
	l.length++
}

// Insert adds a value at the specified index, shifting existing elements to the right.
//
// If the index is equal to the current length, the value is appended to the end.
// If the index is 0, the value becomes the new first element.
//
// Complexity: O(n) in the worst case; O(1) if inserting at the boundaries (0 or Length).
//
// Returns an IndexOutOfBounds error if the index is out of range [0, Length].
func (l *DoubleLinkedList[T]) Insert(idx int, x T) error {
	size := l.Length()
	if idx < 0 || idx > size {
		return NewIndexOutOfBoundError(idx, size)
	}

	if idx == 0 {
		l.prepend(x)
		return nil
	}
	if idx == size {
		l.Append(x)
		return nil
	}

	target := l.get(idx)
	node := node.NewDoubleLinkedNode(x)

	if !l.reversed {
		prev := target.Previous
		node.Next = target
		node.Previous = prev
		if prev != nil {
			prev.Next = node
		}
		target.Previous = node
	} else {
		nxt := target.Next
		node.Previous = target
		node.Next = nxt
		if nxt != nil {
			nxt.Previous = node
		}
		target.Next = node
	}

	l.length++
	return nil
}

// Remove deletes the element at the specified index and returns its value.
// It optimizes removal by checking if the index is at the boundaries (0 or length-1).
//
// Complexity: O(n) in the worst case; O(1) if removing from the start or end.
//
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *DoubleLinkedList[T]) Remove(idx int) (T, error) {
	size := l.Length()
	if idx < 0 || idx >= size {
		var zero T
		return zero, NewIndexOutOfBoundError(idx, size-1)
	}

	current := l.get(idx)
	val := current.Value
	p, n := current.Previous, current.Next

	if p != nil {
		p.Next = n
	}
	if n != nil {
		n.Previous = p
	}

	if current == l.first {
		l.first = l.forward(current)
	}

	if current == l.last {
		l.last = l.backward(current)
	}

	current.Next = nil
	current.Previous = nil
	var zero T
	current.Value = zero
	l.length--
	return val, nil
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
//
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
			if !yield(current.Value) {
				return
			}

			if !l.reversed {
				current = current.Next
			} else {
				current = current.Previous
			}
		}
	}
}

// Backward returns a sequence that yields elements in order from index length-1 to 0.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *DoubleLinkedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := l.last
		for current != nil {
			if !yield(current.Value) {
				return
			}
			if !l.reversed {
				current = current.Previous
			} else {
				current = current.Next
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
			if !yield(index, current.Value) {
				return
			}

			if !l.reversed {
				current = current.Next
			} else {
				current = current.Previous
			}
		}
	}
}

// ToSlice exports the list elements into a native Go slice.
// It pre-allocates the slice based on the current list length for efficiency.
//
// Complexity: O(n).
func (l *DoubleLinkedList[T]) ToSlice() []T {
	if l.length == 0 {
		return nil
	}

	slice := make([]T, l.length)
	for idx, value := range l.Enumerate() {
		slice[idx] = value
	}
	return slice
}

// Clear removes all elements from the list.
//
// After calling Clear, the list will be empty and its length will be zero.
// This operation is typically more efficient than creating a new list
// as it may reuse the underlying storage.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (l *DoubleLinkedList[T]) Clear() {
	var zero T
	current := l.first
	for current != nil {
		next := current.Next
		if l.reversed {
			next = current.Previous
		}
		current.Previous = nil
		current.Next = nil
		current.Value = zero
		current = next
	}
	l.first = nil
	l.last = nil
	l.length = 0
	l.reversed = false
}

// MarshalJSON converts the list into a JSON array.
// It uses the internal serialization utility to ensure elements are
// encoded in their current logical order.
//
// Complexity: O(n).
func (l *DoubleLinkedList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(l)
}

// UnmarshalJSON populates the list from a JSON array.
// It clears any existing elements before appending the new ones from the JSON data.
//
// Note: This operation is destructive; it calls Clear() to remove all existing
// elements before appending the ones from the JSON data.
//
// Complexity: O(n + k) where k is the number of elements in the JSON.
func (l *DoubleLinkedList[T]) UnmarshalJSON(data []byte) error {
	return collection.Unmarshal(data, l.Clear, l.Append)
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *DoubleLinkedList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *DoubleLinkedList[T]) String() string {
	return fmt.Sprint(l)
}
