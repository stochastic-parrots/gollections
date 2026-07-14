package sortedlist

import (
	"cmp"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// OrderedArraySortedList is a slice-backed [SortedList] for naturally ordered values.
type OrderedArraySortedList[T cmp.Ordered] = sortedlist.OrderedArraySortedList[T]

var _ SortedList[int] = &sortedlist.OrderedArraySortedList[int]{}

// OrderedArrayFactory constructs slice-backed sorted lists in natural order.
//
// Ordered array sorted lists avoid custom comparator calls. They provide
// O(log N) lookup, O(1) indexed access, and O(N) single-value insertion or
// removal because values may need to be shifted.
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
type OrderedArrayFactory[T cmp.Ordered] struct{}

// OrderedArray returns an array sorted-list factory using natural order.
func OrderedArray[T cmp.Ordered]() OrderedArrayFactory[T] {
	return OrderedArrayFactory[T]{}
}

// New creates an empty ordered array sorted list with the requested capacity.
func (OrderedArrayFactory[T]) New(capacity int) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedList[T](capacity)
}

// From creates an ordered array sorted list using data as its backing storage.
//
// WARNING: From sorts data in place and transfers its backing storage to the
// returned list. Use [OrderedArrayFactory.Clone] to preserve the source slice.
func (OrderedArrayFactory[T]) From(data []T) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListFromSlice(data)
}

// Clone creates an ordered array sorted list from a sorted copy of data.
//
// Clone does not modify or retain the provided slice.
func (OrderedArrayFactory[T]) Clone(data []T) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListCloneSlice(data)
}

// FromSeq collects seq into new storage and creates an ordered array sorted list.
func (OrderedArrayFactory[T]) FromSeq(seq iter.Seq[T]) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListFromSeq(seq)
}
