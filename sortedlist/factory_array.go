package sortedlist

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// ArraySortedList is a slice-backed [SortedList] with custom comparator order.
// Its zero value is invalid; construct one with [Array].
type ArraySortedList[T any] = sortedlist.ArraySortedList[T]

var _ SortedList[any] = &sortedlist.ArraySortedList[any]{}

// ArrayFactory constructs slice-backed sorted lists with a fixed comparator.
//
// The comparator follows the same contract as cmp.Compare. Array sorted lists
// provide O(log N) lookup, O(1) indexed access, and O(N) single-value insertion
// or removal because values may need to be shifted.
//
// The zero value is invalid. Create a factory with [Array].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	From/Clone/FromSeq  O(N log N)
//	Add(x)/Remove(x)    O(N)
//	Lookup/Bounds       O(log N)
//	Range               O(log N + K)
//	Get                 O(1)
//	Clear               O(N)
type ArrayFactory[T any] struct {
	compare func(a, b T) int
}

// Array returns an array sorted-list factory using compare for every operation.
func Array[T any](compare func(a, b T) int) ArrayFactory[T] {
	return ArrayFactory[T]{compare: compare}
}

// New creates an empty array sorted list with the requested initial capacity.
func (factory ArrayFactory[T]) New(capacity int) *ArraySortedList[T] {
	return sortedlist.NewArraySortedList(capacity, factory.compare)
}

// From creates an array sorted list using data as its backing storage.
//
// WARNING: From sorts data in place and transfers ownership of its backing
// storage. The caller must not use data or aliases of its backing array after
// this call. Use [ArrayFactory.Clone] to preserve the source slice.
func (factory ArrayFactory[T]) From(data []T) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListFromSlice(data, factory.compare)
}

// Clone creates an array sorted list from a sorted copy of data.
//
// Clone does not modify or retain the provided slice.
func (factory ArrayFactory[T]) Clone(data []T) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListCloneSlice(data, factory.compare)
}

// FromSeq collects seq into new storage and creates an array sorted list.
func (factory ArrayFactory[T]) FromSeq(seq iter.Seq[T]) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListFromSeq(seq, factory.compare)
}
