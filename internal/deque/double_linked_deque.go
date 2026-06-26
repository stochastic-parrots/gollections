package deque

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
	"github.com/stochastic-parrots/gollections/internal/shared/node"
)

// DoubleLinkedDeque represents a doubly linked deque data structure.
// It stores elements in a series of nodes where each node points to both
// its predecessor and its successor, allowing for efficient bidirectional traversal.
type DoubleLinkedDeque[T any] struct {
	first, last *node.DoubleLinkedNode[T]
	length      int
}

// NewDoubleLinkedDeque creates an empty DoubleLinkedDeque ready for use.
func NewDoubleLinkedDeque[T any]() *DoubleLinkedDeque[T] {
	return &DoubleLinkedDeque[T]{
		first:  nil,
		last:   nil,
		length: 0,
	}
}

// Length returns the current number of elements in the deque.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) Length() int {
	return d.length
}

// IsEmpty returns true if the deque contains no elements.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) IsEmpty() bool {
	return d.length == 0
}

// append is the internal implementation for adding a value to the end of the deque.
func (d *DoubleLinkedDeque[T]) append(x T) {
	new := node.NewDoubleLinkedNode(x)
	if d.IsEmpty() {
		d.first = new
		d.last = d.first
		d.length++
		return
	}

	new.Previous = d.last
	d.last.Next = new
	d.last = new
	d.length++
}

// prepend is the internal implementation for adding a value to the start of the deque.
func (d *DoubleLinkedDeque[T]) prepend(x T) {
	if d.IsEmpty() {
		d.append(x)
		return
	}

	new := node.NewDoubleLinkedNode(x)
	new.Next = d.first
	d.first.Previous = new
	d.first = new
	d.length++
}

// Prepend inserts one or more elements at the start of the deque.
//
// Complexity: O(k) where k is the number of elements provided.
func (d *DoubleLinkedDeque[T]) Prepend(xs ...T) {
	for i := len(xs) - 1; i >= 0; i-- {
		d.prepend(xs[i])
	}
}

// Append inserts one or more elements at the end of the deque.
//
// Complexity: O(k) where k is the number of elements provided.
func (d *DoubleLinkedDeque[T]) Append(xs ...T) {
	for _, x := range xs {
		d.append(x)
	}
}

// Front returns the element at the beginning of the deque without removing it.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) Front() (x T, ok bool) {
	if d.IsEmpty() {
		return x, false
	}
	return d.first.Value, true
}

// Back returns the element at the end of the deque without removing it.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) Back() (x T, ok bool) {
	if d.IsEmpty() {
		return x, false
	}
	return d.last.Value, true
}

// Shift removes and returns the element at the beginning of the deque.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) Shift() (T, bool) {
	var zero T
	if d.IsEmpty() {
		return zero, false
	}

	first := d.first
	data := first.Value
	d.first = first.Next
	if d.first == nil {
		d.last = nil
	} else {
		d.first.Previous = nil
	}
	first.Next, first.Previous, first.Value = nil, nil, zero
	d.length--
	return data, true
}

// Pop removes and returns the element at the end of the deque.
//
// Complexity: O(1).
func (d *DoubleLinkedDeque[T]) Pop() (T, bool) {
	var zero T
	if d.IsEmpty() {
		return zero, false
	}

	last := d.last
	data := last.Value
	d.last = last.Previous
	if d.last == nil {
		d.first = nil
	} else {
		d.last.Next = nil
	}
	last.Next, last.Previous, last.Value = nil, nil, zero
	d.length--
	return data, true
}

// All returns a sequence that yields elements in their logical order.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (d *DoubleLinkedDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := d.first
		for current != nil {
			if !yield(current.Value) {
				return
			}

			current = current.Next
		}
	}
}

// Enumerate returns a sequence that yields the index and value of elements
// in their logical order.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (d *DoubleLinkedDeque[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := d.first
		for index := 0; current != nil; index++ {
			if !yield(index, current.Value) {
				return
			}

			current = current.Next
		}
	}
}

// ToSlice exports the deque elements into a native Go slice.
// It pre-allocates the slice based on the current deque length for efficiency.
//
// Complexity: O(n).
func (d *DoubleLinkedDeque[T]) ToSlice() []T {
	if d.length == 0 {
		return nil
	}

	slice := make([]T, d.length)
	for idx, value := range d.Enumerate() {
		slice[idx] = value
	}
	return slice
}

// Clear removes all elements from the deque.
//
// After calling Clear, the deque will be empty and its length will be zero.
// This operation is typically more efficient than creating a new deque
// as it may reuse the underlying storage.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (d *DoubleLinkedDeque[T]) Clear() {
	var zero T
	current := d.first
	for current != nil {
		next := current.Next
		current.Previous = nil
		current.Next = nil
		current.Value = zero
		current = next
	}
	d.first = nil
	d.last = nil
	d.length = 0
}

// MarshalJSON converts the deque into a JSON array.
// It uses the internal serialization utility to ensure elements are
// encoded in their current logical order.
//
// Complexity: O(n).
func (d *DoubleLinkedDeque[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(d)
}

// UnmarshalJSON populates the deque from a JSON array.
// It clears any existing elements before appending the new ones from the JSON data.
//
// Note: This operation is destructive; it calls Clear() to remove all existing
// elements before appending the ones from the JSON data.
//
// Complexity: O(n + k) where k is the number of elements in the JSON.
func (d *DoubleLinkedDeque[T]) UnmarshalJSON(data []byte) error {
	return collection.Unmarshal(data, d.Clear, d.Append)
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
//
// Complexity: O(1) as it respects a fixed display limit.
func (d *DoubleLinkedDeque[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, d, d.Length())
}

// String returns a string representation of the deque.
//
// Complexity: O(1) as it respects a fixed display limit.
func (d *DoubleLinkedDeque[T]) String() string {
	return fmt.Sprint(d)
}
