package sortedlist

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"sort"

	ilist "github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// ArraySortedList is a slice-backed sorted list.
//
// Elements are kept in comparator order after every mutation. Duplicate values
// are allowed and inserted after existing equivalent values. This structure is
// optimized for read-heavy workloads; frequent insertions and removals may be
// costly because preserving sorted order requires shifting elements.
type ArraySortedList[T any] struct {
	data    []T
	compare func(a, b T) int
}

// NewArraySortedList creates an empty ArraySortedList with pre-allocated capacity.
func NewArraySortedList[T any](capacity int, compare func(a, b T) int) *ArraySortedList[T] {
	return &ArraySortedList[T]{
		data:    make([]T, 0, capacity),
		compare: compare,
	}
}

// NewArraySortedListFromSlice creates an ArraySortedList using the provided slice as storage.
//
// The input slice is sorted in place.
func NewArraySortedListFromSlice[T any](data []T, compare func(a, b T) int) *ArraySortedList[T] {
	slices.SortStableFunc(data, compare)
	return &ArraySortedList[T]{data: data, compare: compare}
}

// NewArraySortedListCloneSlice creates an ArraySortedList from a sorted clone of the provided slice.
func NewArraySortedListCloneSlice[T any](data []T, compare func(a, b T) int) *ArraySortedList[T] {
	return NewArraySortedListFromSlice(slices.Clone(data), compare)
}

// NewArraySortedListFromSeq creates an ArraySortedList from an iterator.
func NewArraySortedListFromSeq[T any](seq iter.Seq[T], compare func(a, b T) int) *ArraySortedList[T] {
	return NewArraySortedListFromSlice(slices.Collect(seq), compare)
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) Length() int {
	return len(l.data)
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) IsEmpty() bool {
	return len(l.data) == 0
}

// Get retrieves the value at the specified sorted index.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, ilist.NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	return l.data[index], nil
}

// Find locates the first index equivalent to x according to the comparator.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Find(x T) (idx int, ok bool) {
	if l.IsEmpty() {
		return -1, false
	}

	idx = sort.Search(len(l.data), func(i int) bool {
		return l.compare(l.data[i], x) >= 0
	})

	if idx < len(l.data) && l.compare(l.data[idx], x) == 0 {
		return idx, true
	}

	return -1, false
}

// Contains returns true if an equivalent value exists in the list.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Contains(x T) bool {
	_, ok := l.Find(x)
	return ok
}

// First returns the first value in comparator order.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) First() (T, bool) {
	if l.IsEmpty() {
		var zero T
		return zero, false
	}

	return l.data[0], true
}

// Last returns the last value in comparator order.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) Last() (T, bool) {
	if l.IsEmpty() {
		var zero T
		return zero, false
	}

	return l.data[len(l.data)-1], true
}

// Add inserts one or more values while preserving the sorted invariant.
//
// Complexity: O(N) for a single value, O((N+K) log (N+K)) for multiple values.
func (l *ArraySortedList[T]) Add(xs ...T) {
	if len(xs) == 0 {
		return
	}

	if len(xs) == 1 {
		l.insert(xs[0])
		return
	}

	l.data = append(l.data, xs...)
	slices.SortStableFunc(l.data, l.compare)
}

// Remove removes the first value equivalent to x according to the comparator.
//
// Complexity: O(N).
func (l *ArraySortedList[T]) Remove(x T) bool {
	idx, ok := l.Find(x)
	if !ok {
		return false
	}

	copy(l.data[idx:], l.data[idx+1:])

	var zero T
	l.data[len(l.data)-1] = zero
	l.data = l.data[:len(l.data)-1]

	return true
}

// All returns a sequence that yields elements in sorted order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (l *ArraySortedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range l.data {
			if !yield(value) {
				return
			}
		}
	}
}

// Backward returns a sequence that yields elements in reverse sorted order.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (l *ArraySortedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for idx := len(l.data) - 1; idx >= 0; idx-- {
			if !yield(l.data[idx]) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the sorted index and value of each element.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (l *ArraySortedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range l.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

// ToSlice exports the sorted elements into a native Go slice.
//
// Complexity: O(N).
func (l *ArraySortedList[T]) ToSlice() []T {
	if len(l.data) == 0 {
		return nil
	}

	slice := make([]T, len(l.data))
	copy(slice, l.data)
	return slice
}

// Clear removes all elements from the list.
//
// Complexity: O(N).
func (l *ArraySortedList[T]) Clear() {
	clear(l.data)
	l.data = l.data[:0]
}

// MarshalJSON converts the list into a JSON array in sorted order.
//
// Complexity: O(N).
func (l *ArraySortedList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(l)
}

// UnmarshalJSON populates the list from a JSON array and restores sorted order.
//
// The input order is not preserved; values are sorted according to the list
// comparator before replacing the current contents.
//
// Complexity: O(M + N log N), where M is the current length and N is the number
// of decoded values.
func (l *ArraySortedList[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	slices.SortStableFunc(values, l.compare)
	l.Clear()
	l.data = values
	return nil
}

// Format implements fmt.Formatter.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *ArraySortedList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, l, l.Length())
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *ArraySortedList[T]) String() string {
	return fmt.Sprint(l)
}

func (l *ArraySortedList[T]) insert(x T) {
	idx := sort.Search(len(l.data), func(i int) bool {
		return l.compare(l.data[i], x) > 0
	})

	var zero T
	l.data = append(l.data, zero)
	copy(l.data[idx+1:], l.data[idx:])
	l.data[idx] = x
}
