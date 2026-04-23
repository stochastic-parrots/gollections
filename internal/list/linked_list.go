package list

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
	"github.com/stochastic-parrots/gollections/internal/shared/node"
)

// LinkedList represents a singly linked list data structure.
// It consists of nodes where each element points to the next, making it efficient
// for sequential insertion and deletion at the ends, but requiring O(n) for random access.
type LinkedList[T any] struct {
	first, last *node.LinkedNode[T]
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
		current = current.Next
	}

	return current.Value, nil
}

// Find locates the index of an element using a linear search.
//
// Complexity: O(n).
func (l *LinkedList[T]) Find(x T, cmp func(a, b T) int) (idx int, ok bool) {
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

// Contains returns true if the element exists in the list according.
//
// Complexity: O(n).
func (l *LinkedList[T]) Contains(x T, cmp func(a, b T) int) bool {
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
func (l *LinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= l.Length() {
		return NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		current = current.Next
	}

	current.Value = x
	return nil
}

func (l *LinkedList[T]) append(x T) {
	new := node.NewLinkedNode(x)
	if l.first == nil {
		l.first = new
	} else {
		l.last.Next = new
	}
	l.last = new
	l.length++
}

// Append adds one or more elements to the end of the list.
//
// Complexity: O(len(xs)) the number of elements provided.
func (l *LinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		l.append(x)
	}
}

// Insert adds a value at the specified index, shifting existing elements to the right.
//
// If the index is equal to the current length, the value is appended to the end.
// If the index is 0, the value becomes the new first element.
//
// Complexity: O(n) in the worst case; O(1) if inserting at the boundaries (0 or Length).
// Returns an IndexOutOfBounds error if the index is out of range [0, Length].
func (l *LinkedList[T]) Insert(idx int, x T) error {
	size := l.Length()
	if idx < 0 || idx > size {
		return NewIndexOutOfBoundError(idx, size)
	}

	if idx == size {
		l.Append(x)
		return nil
	}

	node := node.NewLinkedNode(x)

	if idx == 0 {
		node.Next = l.first
		l.first = node
		l.length++
		return nil
	}

	current := l.first
	for range idx - 1 {
		current = current.Next
	}

	node.Next = current.Next
	current.Next = node

	l.length++
	return nil
}

// Remove deletes the element at the specified index and returns its value.
// It optimizes removal by checking if the index is at the boundaries (0 or length-1).
//
// Complexity: O(n) in the worst case; O(1) if removing from the start or end.
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *LinkedList[T]) Remove(idx int) (T, error) {
	size := l.Length()
	if idx < 0 || idx >= size {
		var zero T
		return zero, NewIndexOutOfBoundError(idx, size-1)
	}

	var val T
	if idx == 0 {
		val = l.first.Value
		old := l.first
		l.first = l.first.Next

		if l.first == nil {
			l.last = nil
		}

		old.Next = nil
		var zero T
		old.Value = zero
	} else {
		prev := l.first
		for range idx - 1 {
			prev = prev.Next
		}
		removed := prev.Next
		val = removed.Value
		prev.Next = removed.Next

		if idx == size-1 {
			l.last = prev
		}
		removed.Next = nil
		var zero T
		removed.Value = zero
	}

	l.length--
	return val, nil
}

// Reverse inverts the order of the elements in the list in-place.
//
// Complexity: O(n).
// Note: Unlike DoubleLinkedList, this requires a full traversal to update pointers.
func (l *LinkedList[T]) Reverse() {
	if l.Length() <= 1 {
		return
	}

	var previous *node.LinkedNode[T]
	current := l.first
	l.last = l.first

	for current != nil {
		next := current.Next
		current.Next = previous
		previous = current
		current = next
	}

	l.first = previous
	l.last.Next = nil
}

// All returns a sequence that yields elements from the first to the last node.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *LinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for current := l.first; current != nil; current = current.Next {
			if !yield(current.Value) {
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
		for current, index := l.first, 0; current != nil; current = current.Next {
			if !yield(index, current.Value) {
				return
			}
			index++
		}
	}
}

// Backward returns a sequence that yields elements in order from index length-1 to 0.
//
// Complexity: O(n) for a full traversal, O(1) per step. O(n) space.
func (l *LinkedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		slice := l.ToSlice()
		for idx := len(slice) - 1; idx >= 0; idx-- {
			if !yield(slice[idx]) {
				return
			}
		}
	}
}

// ToSlice exports the list elements into a native Go slice.
// It pre-allocates the slice based on the current list length for efficiency.
//
// Complexity: O(n).
func (l *LinkedList[T]) ToSlice() []T {
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
func (l *LinkedList[T]) Clear() {
	var zero T
	current := l.first
	for current != nil {
		next := current.Next
		current.Next = nil
		current.Value = zero
		current = next
	}
	l.first = nil
	l.last = nil
	l.length = 0
}

// MarshalJSON converts the list into a JSON array.
// It uses the internal serialization utility to ensure elements are
// encoded in their current logical order.
//
// Complexity: O(n).
func (l *LinkedList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(l)
}

// UnmarshalJSON populates the list from a JSON array.
// It clears any existing elements before appending the new ones from the JSON data.
//
// Note: This operation is destructive; it calls Clear() to remove all existing
// elements before appending the ones from the JSON data.
//
// Complexity: O(n + k) where k is the number of elements in the JSON.
func (l *LinkedList[T]) UnmarshalJSON(data []byte) error {
	return collection.Unmarshal(data, l.Clear, l.Append)
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *LinkedList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *LinkedList[T]) String() string {
	return fmt.Sprint(l)
}
