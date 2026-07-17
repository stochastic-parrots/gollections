package sortedlist

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// OrderedArraySortedList is a slice-backed sorted list for naturally ordered values.
//
// Elements are kept in natural order after every mutation. Equivalent values
// are allowed, but their relative order is not part of the list contract. This
// structure is optimized for read-heavy workloads; frequent insertions and
// removals may be costly because preserving sorted order requires shifting
// elements.
type OrderedArraySortedList[T cmp.Ordered] struct {
	data []T
}

// NewOrderedArraySortedList creates an empty OrderedArraySortedList with pre-allocated capacity.
func NewOrderedArraySortedList[T cmp.Ordered](capacity int) *OrderedArraySortedList[T] {
	return &OrderedArraySortedList[T]{
		data: make([]T, 0, capacity),
	}
}

// NewOrderedArraySortedListFromSlice creates an OrderedArraySortedList using the provided slice as storage.
//
// The input slice is sorted in place.
func NewOrderedArraySortedListFromSlice[T cmp.Ordered](data []T) *OrderedArraySortedList[T] {
	slices.Sort(data)
	return &OrderedArraySortedList[T]{data: data}
}

// NewOrderedArraySortedListCloneSlice creates an OrderedArraySortedList from a sorted clone of the provided slice.
func NewOrderedArraySortedListCloneSlice[T cmp.Ordered](data []T) *OrderedArraySortedList[T] {
	return NewOrderedArraySortedListFromSlice(slices.Clone(data))
}

// NewOrderedArraySortedListFromSeq creates an OrderedArraySortedList from an iterator.
func NewOrderedArraySortedListFromSeq[T cmp.Ordered](seq iter.Seq[T]) *OrderedArraySortedList[T] {
	return NewOrderedArraySortedListFromSlice(slices.Collect(seq))
}

// Length returns the current number of elements in the list.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) Length() int {
	return len(l.data)
}

// IsEmpty returns true if the list contains no elements.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) IsEmpty() bool {
	return len(l.data) == 0
}

// Get retrieves the value at the specified sorted index.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) Get(idx int) (T, error) {
	if idx < 0 || idx >= len(l.data) {
		var zero T
		return zero, list.NewIndexOutOfBoundError(idx, len(l.data)-1)
	}

	return l.data[idx], nil
}

// LowerBound returns the first index whose value does not sort before x.
//
// It returns Length when every value sorts before x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) LowerBound(x T) int {
	idx, _ := slices.BinarySearch(l.data, x)
	return idx
}

// UpperBound returns the first index whose value sorts after x.
//
// It returns Length when no value sorts after x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) UpperBound(x T) int {
	lo, hi := 0, len(l.data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if cmp.Less(x, l.data[mid]) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

// EqualRange returns the half-open index range containing values equivalent to x.
//
// The returned range is [start, end), and is empty when start == end.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) EqualRange(x T) (start, end int) {
	data := l.data
	start, found := slices.BinarySearch(data, x)
	if !found {
		return start, start
	}

	lo, hi := start+1, len(data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if cmp.Less(x, data[mid]) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	return start, lo
}

// Count returns the number of values equivalent to x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Count(x T) int {
	start, end := l.EqualRange(x)
	return end - start
}

// Find locates the first index equivalent to x according to natural order.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Find(x T) (idx int, ok bool) {
	idx, ok = slices.BinarySearch(l.data, x)
	if !ok {
		return -1, false
	}
	return idx, true
}

// Ceiling returns the first value that does not sort before x.
//
// It returns a zero value, -1, and false when every value sorts before x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Ceiling(x T) (value T, idx int, ok bool) {
	idx = l.LowerBound(x)
	if idx >= len(l.data) {
		return value, -1, false
	}

	return l.data[idx], idx, true
}

// Floor returns the last value that does not sort after x.
//
// It returns a zero value, -1, and false when every value sorts after x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Floor(x T) (value T, idx int, ok bool) {
	idx = l.UpperBound(x) - 1
	if idx < 0 {
		return value, -1, false
	}

	return l.data[idx], idx, true
}

// Higher returns the first value that sorts after x.
//
// It returns a zero value, -1, and false when no value sorts after x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Higher(x T) (value T, idx int, ok bool) {
	idx = l.UpperBound(x)
	if idx >= len(l.data) {
		return value, -1, false
	}

	return l.data[idx], idx, true
}

// Lower returns the last value that sorts before x.
//
// It returns a zero value, -1, and false when no value sorts before x.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Lower(x T) (value T, idx int, ok bool) {
	idx = l.LowerBound(x) - 1
	if idx < 0 {
		return value, -1, false
	}

	return l.data[idx], idx, true
}

// Contains returns true if an equivalent value exists in the list.
//
// Complexity: O(log N).
func (l *OrderedArraySortedList[T]) Contains(x T) bool {
	_, found := slices.BinarySearch(l.data, x)
	return found
}

// First returns the first value in natural order.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) First() (T, bool) {
	if l.IsEmpty() {
		var zero T
		return zero, false
	}

	return l.data[0], true
}

// Last returns the last value in natural order.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) Last() (T, bool) {
	if l.IsEmpty() {
		var zero T
		return zero, false
	}

	return l.data[len(l.data)-1], true
}

// Add inserts x while preserving the sorted invariant.
//
// Complexity: O(N).
func (l *OrderedArraySortedList[T]) Add(x T) {
	idx := l.UpperBound(x)

	var zero T
	l.data = append(l.data, zero)
	copy(l.data[idx+1:], l.data[idx:])
	l.data[idx] = x
}

// Adds inserts zero or more values while preserving the sorted invariant.
//
// Complexity: O(N) for one value, or
// O((N + len(xs)) log (N + len(xs))) for multiple values.
func (l *OrderedArraySortedList[T]) Adds(xs ...T) {
	switch len(xs) {
	case 0:
		return
	case 1:
		l.Add(xs[0])
		return
	}

	l.data = append(l.data, xs...)
	slices.Sort(l.data)
}

// Replace changes the value at idx while preserving the sorted invariant.
//
// It returns an index error when idx is out of bounds and ErrOrderViolation
// when x would sort before the previous value or after the next value.
//
// Complexity: O(1).
func (l *OrderedArraySortedList[T]) Replace(idx int, x T) error {
	if idx < 0 || idx >= len(l.data) {
		return list.NewIndexOutOfBoundError(idx, len(l.data)-1)
	}
	if idx > 0 && cmp.Less(x, l.data[idx-1]) {
		return ErrOrderViolation
	}
	if idx+1 < len(l.data) && cmp.Less(l.data[idx+1], x) {
		return ErrOrderViolation
	}

	l.data[idx] = x
	return nil
}

// Remove removes the first value equivalent to x according to natural order.
//
// Complexity: O(N).
func (l *OrderedArraySortedList[T]) Remove(x T) bool {
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
func (l *OrderedArraySortedList[T]) All() iter.Seq[T] {
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
func (l *OrderedArraySortedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range slices.Backward(l.data) {
			if !yield(value) {
				return
			}
		}
	}
}

// Range returns a sequence that yields values in the half-open range [from, to).
//
// Values are selected according to natural order: yielded values do not sort
// before from and do sort before to.
//
// Complexity: O(log N + number of yielded values).
func (l *OrderedArraySortedList[T]) Range(from, to T) iter.Seq[T] {
	return func(yield func(T) bool) {
		start := l.LowerBound(from)
		end := l.LowerBound(to)
		if end < start {
			return
		}

		for idx := start; idx < end; idx++ {
			if !yield(l.data[idx]) {
				return
			}
		}
	}
}

// Enumerate returns a sequence that yields the sorted index and value of each element.
//
// Complexity: O(N) for a full traversal, O(1) per step.
func (l *OrderedArraySortedList[T]) Enumerate() iter.Seq2[int, T] {
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
func (l *OrderedArraySortedList[T]) ToSlice() []T {
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
func (l *OrderedArraySortedList[T]) Clear() {
	clear(l.data)
	l.data = l.data[:0]
}

// MarshalJSON converts the list into a JSON array in sorted order.
//
// Complexity: O(N).
func (l *OrderedArraySortedList[T]) MarshalJSON() ([]byte, error) {
	return collection.Marshal(l)
}

// Format implements fmt.Formatter.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *OrderedArraySortedList[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, l, cap(l.data))
}

// String returns a string representation of the list.
//
// Complexity: O(1) as it respects a fixed display limit.
func (l *OrderedArraySortedList[T]) String() string {
	return fmt.Sprint(l)
}
