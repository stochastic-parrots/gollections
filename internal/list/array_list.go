package list

import (
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// ArrayList represents a dynamic array (slice-backed) data structure.
// It provides constant-time O(1) random access and is highly cache-efficient
// due to contiguous memory allocation. It is the preferred general-purpose
// list for most use cases.
type ArrayList[T any] struct {
	data []T
}

// NewArrayList creates an empty ArrayList with an initial pre-allocated capacity.
// Pre-allocating capacity can significantly improve performance by reducing
// the number of memory reallocations as the list grows.
func NewArrayList[T any](size int) *ArrayList[T] {
	data := make([]T, 0, size)
	return &ArrayList[T]{data: data}
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (l *ArrayList[T]) Length() int {
	return len(l.data)
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (l *ArrayList[T]) IsEmpty() bool {
	return len(l.data) == 0
}

// Get retrieves the value at the specified index.
//
// Complexity: O(1).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *ArrayList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	return l.data[index], nil
}

// Find locates the index of an element using a linear search.
//
// Complexity: O(n).
func (l *ArrayList[T]) Find(x T, cmp func(a, b T) int) (idx int, ok bool) {
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
func (l *ArrayList[T]) Contains(x T, cmp func(a, b T) int) bool {
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
// Complexity: O(1).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *ArrayList[T]) Set(index int, x T) error {
	if index < 0 || index >= len(l.data) {
		return NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	l.data[index] = x
	return nil
}

// Append adds one or more elements to the end of the list.
//
// Complexity: Amortized O(1) per element.
// If the underlying capacity is exceeded, a new, larger array is allocated
// and all elements are copied (O(n)).
func (l *ArrayList[T]) Append(xs ...T) {
	l.data = append(l.data, xs...)
}

// Insert places an element at the specified index.
//
// Complexity: O(n) as it requires shifting elements to the right.
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *ArrayList[T]) Insert(index int, x T) error {
	size := len(l.data)
	if index < 0 || index > size {
		return NewIndexOutOfBoundError(index, size)
	}

	if index == size {
		l.Append(x)
		return nil
	}

	var zero T
	l.data = append(l.data, zero)
	copy(l.data[index+1:], l.data[index:size])
	l.data[index] = x

	return nil
}

// Remove deletes the element at the specified index and returns its value.
//
// Complexity: O(n) as it requires shifting elements to the left.
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *ArrayList[T]) Remove(index int) (T, error) {
	size := len(l.data)
	if index < 0 || index >= size {
		var zero T
		return zero, NewIndexOutOfBoundError(index, size-1)
	}

	val := l.data[index]

	if index < size-1 {
		copy(l.data[index:], l.data[index+1:])
	}

	var zero T
	l.data[size-1] = zero
	l.data = l.data[:size-1]

	return val, nil
}

// Reverse inverts the order of the elements in the list in-place.
//
// Complexity: O(n).
// Note: This operation modifies the underlying slice using the optimized slices.Reverse.
func (l *ArrayList[T]) Reverse() {
	slices.Reverse(l.data)
}

// All returns a sequence that yields elements in order from index 0 to length-1.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *ArrayList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range l.data {
			if !yield(value) {
				return
			}
		}
	}
}

// Backward returns a sequence that yields elements in order from index length-1 to 0.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *ArrayList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for idx := l.Length() - 1; idx >= 0; idx-- {
			if !yield(l.data[idx]) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the index and value of each element.
//
// Complexity: O(n) for a full traversal, O(1) per step.
func (l *ArrayList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range l.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

// ToSlice exports the list elements into a native Go slice.
// It pre-allocates the slice based on the current list length for efficiency.
//
// Complexity: O(n).
func (l *ArrayList[T]) ToSlice() []T {
	if len(l.data) == 0 {
		return nil
	}

	slice := make([]T, len(l.data))
	copy(slice, l.data)
	return slice
}

// Clear removes all elements from the list.
//
// After calling Clear, the list will be empty and its length will be zero.
// This operation is typically more efficient than creating a new array list
// as it may reuse the underlying storage.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (l *ArrayList[T]) Clear() {
	clear(l.data)
	l.data = l.data[:0]
}

// MarshalJSON converts the list into a JSON array.
// It uses the internal serialization utility to ensure elements are
// encoded in their current logical order.
//
// Complexity: O(n).
func (l *ArrayList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(l)
}

// UnmarshalJSON populates the list from a JSON array.
// It clears any existing elements before appending the new ones from the JSON data.
//
// Note: This operation is destructive; it calls Clear() to remove all existing
// elements before appending the ones from the JSON data.
//
// Complexity: O(n + k) where k is the number of elements in the JSON.
func (l *ArrayList[T]) UnmarshalJSON(data []byte) error {
	return collection.Unmarshal(data, l.Clear, l.Append)
}

// Format implements the fmt.Formatter interface, allowing custom formatting
// with verbs like %v, %+v, and %#v.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *ArrayList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *ArrayList[T]) String() string {
	return fmt.Sprint(l)
}
