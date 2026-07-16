package datastructs

import (
	"fmt"
	"slices"
)

// StdSortedList adapts a sorted slice to the benchmark SortedList contract.
type StdSortedList struct {
	data []int
}

// NewStdSortedList creates an empty sorted-slice baseline.
func NewStdSortedList(capacity int) *StdSortedList {
	return &StdSortedList{data: make([]int, 0, capacity)}
}

// Add inserts x while preserving sorted order.
func (l *StdSortedList) Add(x int) {
	idx, _ := slices.BinarySearch(l.data, x)
	l.data = append(l.data, 0)
	copy(l.data[idx+1:], l.data[idx:])
	l.data[idx] = x
}

// Adds inserts all values while preserving sorted order.
func (l *StdSortedList) Adds(xs ...int) {
	switch len(xs) {
	case 0:
		return
	case 1:
		l.Add(xs[0])
	default:
		l.data = append(l.data, xs...)
		slices.Sort(l.data)
	}
}

// Remove deletes the first matching value.
func (l *StdSortedList) Remove(x int) bool {
	idx, ok := slices.BinarySearch(l.data, x)
	if !ok {
		return false
	}

	copy(l.data[idx:], l.data[idx+1:])
	l.data[len(l.data)-1] = 0
	l.data = l.data[:len(l.data)-1]
	return true
}

// Get returns the value at idx.
func (l *StdSortedList) Get(idx int) (int, error) {
	if idx < 0 || idx >= len(l.data) {
		return 0, fmt.Errorf("index %d out of bounds", idx)
	}
	return l.data[idx], nil
}

// LowerBound returns the first index whose value is not less than x.
func (l *StdSortedList) LowerBound(x int) int {
	idx, _ := slices.BinarySearch(l.data, x)
	return idx
}

// Contains returns true when x exists in the list.
func (l *StdSortedList) Contains(x int) bool {
	_, ok := slices.BinarySearch(l.data, x)
	return ok
}

// Length returns the number of values in the list.
func (l *StdSortedList) Length() int {
	return len(l.data)
}

// Clear removes all values while preserving capacity.
func (l *StdSortedList) Clear() {
	clear(l.data)
	l.data = l.data[:0]
}
