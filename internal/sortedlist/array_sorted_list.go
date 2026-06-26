package sortedlist

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

// ArraySortedList is a slice-backed sorted list.
//
// Elements are kept in comparator order after every mutation. Equivalent values
// are allowed, but their relative order is not part of the list contract. This
// structure is optimized for read-heavy workloads; frequent insertions and
// removals may be costly because preserving sorted order requires shifting
// elements.
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
	slices.SortFunc(data, compare)
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
func (l *ArraySortedList[T]) Get(idx int) (T, error) {
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
func (l *ArraySortedList[T]) LowerBound(x T) int {
	data := l.data
	compare := l.compare

	lo, hi := 0, len(data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if compare(data[mid], x) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo
}

// UpperBound returns the first index whose value sorts after x.
//
// It returns Length when no value sorts after x.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) UpperBound(x T) int {
	data := l.data
	compare := l.compare

	lo, hi := 0, len(data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if compare(data[mid], x) <= 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo
}

// EqualRange returns the half-open index range containing values equivalent to x.
//
// The returned range is [start, end), and is empty when start == end.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) EqualRange(x T) (start, end int) {
	data := l.data
	compare := l.compare

	lo, hi := 0, len(data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if compare(data[mid], x) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	start = lo
	if start >= len(data) || compare(data[start], x) != 0 {
		return start, start
	}

	lo, hi = start+1, len(data)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if compare(data[mid], x) <= 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return start, lo
}

// Count returns the number of values equivalent to x.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Count(x T) int {
	start, end := l.EqualRange(x)
	return end - start
}

// Find locates the first index equivalent to x according to the comparator.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Find(x T) (idx int, ok bool) {
	idx = l.LowerBound(x)

	if idx < len(l.data) && l.compare(l.data[idx], x) == 0 {
		return idx, true
	}

	return -1, false
}

// Ceiling returns the first value that does not sort before x.
//
// It returns a zero value, -1, and false when every value sorts before x.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Ceiling(x T) (value T, idx int, ok bool) {
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
func (l *ArraySortedList[T]) Floor(x T) (value T, idx int, ok bool) {
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
func (l *ArraySortedList[T]) Higher(x T) (value T, idx int, ok bool) {
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
func (l *ArraySortedList[T]) Lower(x T) (value T, idx int, ok bool) {
	idx = l.LowerBound(x) - 1
	if idx < 0 {
		return value, -1, false
	}

	return l.data[idx], idx, true
}

// Contains returns true if an equivalent value exists in the list.
//
// Complexity: O(log N).
func (l *ArraySortedList[T]) Contains(x T) bool {
	idx := l.LowerBound(x)
	return idx < len(l.data) && l.compare(l.data[idx], x) == 0
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
		idx := l.UpperBound(xs[0])

		var zero T
		l.data = append(l.data, zero)
		copy(l.data[idx+1:], l.data[idx:])
		l.data[idx] = xs[0]
		return
	}

	l.data = append(l.data, xs...)
	slices.SortFunc(l.data, l.compare)
}

// Replace changes the value at idx while preserving the sorted invariant.
//
// It returns an index error when idx is out of bounds and ErrOrderViolation
// when x would sort before the previous value or after the next value.
//
// Complexity: O(1).
func (l *ArraySortedList[T]) Replace(idx int, x T) error {
	if idx < 0 || idx >= len(l.data) {
		return list.NewIndexOutOfBoundError(idx, len(l.data)-1)
	}
	if idx > 0 && l.compare(l.data[idx-1], x) > 0 {
		return ErrOrderViolation
	}
	if idx+1 < len(l.data) && l.compare(x, l.data[idx+1]) > 0 {
		return ErrOrderViolation
	}

	l.data[idx] = x
	return nil
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

// Range returns a sequence that yields values in the half-open range [from, to).
//
// Values are selected according to comparator order: yielded values do not sort
// before from and do sort before to.
//
// Complexity: O(log N + K), where K is the number of yielded values.
func (l *ArraySortedList[T]) Range(from, to T) iter.Seq[T] {
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
// Existing backing storage is reused when it has enough capacity.
//
// Complexity: O(M + N log N), where M is the current length and N is the number
// of decoded values.
func (l *ArraySortedList[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	slices.SortFunc(values, l.compare)
	l.Clear()
	if len(values) > cap(l.data) {
		l.data = values
		return nil
	}

	l.data = append(l.data, values...)
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
