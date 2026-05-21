package deque

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// RingBufferDeque is a dynamic double-ended queue backed by a circular array.
type RingBufferDeque[T any] struct {
	data  []T
	read  int
	write int
	count int
}

// NewRingBufferDeque creates an empty RingBufferDeque with an initial pre-allocated capacity.
// Pre-allocating capacity can significantly improve performance by reducing
// the number of memory reallocations as the deque grows.
func NewRingBufferDeque[T any](size int) *RingBufferDeque[T] {
	data := make([]T, size)
	return &RingBufferDeque[T]{data: data}
}

// next returns the index that follows idx in the ring, wrapping around
// to zero when the end of the underlying array is reached.
func (rb *RingBufferDeque[T]) next(idx int) int {
	return (idx + 1) % len(rb.data)
}

// prev returns the index that precedes idx in the ring, wrapping around
// to the last position when idx is zero.
func (rb *RingBufferDeque[T]) prev(idx int) int {
	return (idx - 1 + len(rb.data)) % len(rb.data)
}

// grow doubles the capacity of the ring buffer, reorganizing all elements
// into contiguous logical order starting at index 0.
//
// This is called automatically by append and prepend when the buffer is full.
// Complexity: O(n).
func (rb *RingBufferDeque[T]) grow() {
	oldCap := len(rb.data)
	newCap := 1
	if oldCap > 0 {
		newCap = oldCap * 2
	}
	newData := make([]T, newCap)
	for i := range rb.count {
		newData[i] = rb.data[(rb.read+i)%oldCap]
	}
	rb.data = newData
	rb.read = 0
	rb.write = rb.count
}

// Front returns the element at the beginning of the deque without removing it.
//
// It returns the zero value of T and false if the deque is empty.
// Complexity: O(1).
func (rb *RingBufferDeque[T]) Front() (T, bool) {
	if rb.count == 0 {
		var zero T
		return zero, false
	}

	return rb.data[rb.read], true
}

// Back returns the element at the end of the deque without removing it.
//
// It returns the zero value of T and false if the deque is empty.
// Complexity: O(1).
func (rb *RingBufferDeque[T]) Back() (T, bool) {
	if rb.count == 0 {
		var zero T
		return zero, false
	}

	idx := rb.prev(rb.write)
	return rb.data[idx], true
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (rb *RingBufferDeque[T]) IsEmpty() bool {
	return rb.count == 0
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (rb *RingBufferDeque[T]) Length() int {
	return rb.count
}

// All returns a sequence that yields elements in order from index 0 to length-1.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (rb *RingBufferDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range rb.count {
			idx := (rb.read + i) % len(rb.data)
			if !yield(rb.data[idx]) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the index and value of each element.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (rb *RingBufferDeque[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range rb.count {
			idx := (rb.read + i) % len(rb.data)
			if !yield(i, rb.data[idx]) {
				return
			}
		}
	}
}

// ToSlice exports the list elements into a native Go slice.
// It pre-allocates the slice based on the current list length for efficiency.
//
// Complexity: O(n).
func (rb *RingBufferDeque[T]) ToSlice() []T {
	if rb.count == 0 {
		return nil
	}

	slice := make([]T, rb.count)
	for i := range rb.count {
		slice[i] = rb.data[(rb.read+i)%len(rb.data)]
	}
	return slice
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
//
// Complexity: O(1) as it respects a fixed display limit.
func (rb *RingBufferDeque[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, rb, rb.Length())
}

// String returns a string representation of the deque.
//
// Complexity: O(1) as it respects a fixed display limit.
func (rb *RingBufferDeque[T]) String() string {
	return fmt.Sprint(rb)
}

// MarshalJSON converts the deque into a JSON array.
// It uses the internal serialization utility to ensure elements are
// encoded in their current logical order.
//
// Complexity: O(n).
func (rb *RingBufferDeque[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(rb)
}

// Clear removes all elements from the list.
//
// After calling Clear, the list will be empty and its length will be zero.
// This operation is typically more efficient than creating a new array list
// as it may reuse the underlying storage.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (rb *RingBufferDeque[T]) Clear() {
	clear(rb.data)
	rb.count = 0
	rb.read = 0
	rb.write = 0
}

func (rb *RingBufferDeque[T]) append(x T) {
	if rb.count == len(rb.data) {
		rb.grow()
	}

	rb.data[rb.write] = x
	rb.write = rb.next(rb.write)
	rb.count++
}

// Append adds the given elements to the end of the deque.
//
// Complexity: O(k) where k is the number of elements.
func (rb *RingBufferDeque[T]) Append(xs ...T) {
	for _, x := range xs {
		rb.append(x)
	}
}

func (rb *RingBufferDeque[T]) prepend(x T) {
	if rb.count == len(rb.data) {
		rb.grow()
	}

	rb.read = rb.prev(rb.read)
	rb.data[rb.read] = x
	rb.count++
}

// Prepend adds the given elements to the beginning of the deque.
// The relative order of the provided elements is preserved at the front.
//
// Complexity: O(k) where k is the number of elements.
func (rb *RingBufferDeque[T]) Prepend(xs ...T) {
	for i := len(xs) - 1; i >= 0; i-- {
		rb.prepend(xs[i])
	}
}

// Shift removes and returns the element from the beginning of the deque.
//
// Returns the zero value and false if the deque is empty.
// Complexity: O(1).
func (rb *RingBufferDeque[T]) Shift() (T, bool) {
	if rb.count == 0 {
		var zero T
		return zero, false
	}

	x := rb.data[rb.read]
	var zero T
	rb.data[rb.read] = zero
	rb.read = rb.next(rb.read)
	rb.count--
	if rb.count == 0 {
		rb.read = 0
		rb.write = 0
	}
	return x, true
}

// Pop removes and returns the element from the end of the deque.
//
// Returns the zero value and false if the deque is empty.
// Complexity: O(1).
func (rb *RingBufferDeque[T]) Pop() (T, bool) {
	if rb.count == 0 {
		var zero T
		return zero, false
	}

	rb.write = rb.prev(rb.write)
	x := rb.data[rb.write]
	var zero T
	rb.data[rb.write] = zero
	rb.count--
	if rb.count == 0 {
		rb.read = 0
		rb.write = 0
	}
	return x, true
}

// UnmarshalJSON populates the deque from a JSON array.
// It clears any existing elements before appending the new ones from the JSON data.
//
// Note: This operation is destructive; it calls Clear() to remove all existing
// elements before appending the ones from the JSON data.
//
// Complexity: O(n + k) where k is the number of elements in the JSON.
func (rb *RingBufferDeque[T]) UnmarshalJSON(data []byte) error {
	return collection.Unmarshal(data, rb.Clear, rb.Append)
}
