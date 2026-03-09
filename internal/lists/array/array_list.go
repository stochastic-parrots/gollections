package array

import (
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/formatters"
	"github.com/stochastic-parrots/gollections/internal/lists"
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
//
func (l *ArrayList[T]) Length() int {
	return len(l.data)
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
//
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
		return zero, lists.NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	return l.data[index], nil
}

// Set updates the value at the specified index.
//
// Complexity: O(1).
// Returns an IndexOutOfBounds error if the index is out of range.
func (l *ArrayList[T]) Set(index int, x T) error {
	if index < 0 || index >= len(l.data) {
		return lists.NewIndexOutOfBoundError(index, len(l.data)-1)
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

// Reverse inverts the order of the elements in the list in-place.
//
// Complexity: O(n).
// Note: This operation modifies the underlying slice using the optimized slices.Reverse.
func (l *ArrayList[T]) Reverse() {
	slices.Reverse(l.data)
}

// Format implements the fmt.Formatter interface for custom string formatting.
func (l *ArrayList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, cap(l.data))
}

// String returns a string representation of the list.
func (l *ArrayList[T]) String() string {
	return fmt.Sprint(l)
}
