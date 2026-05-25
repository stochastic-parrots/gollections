package sortedlist

import (
	"iter"

	constructor "github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// ArraySortedList is a slice-backed [SortedList].
type ArraySortedList[T any] = *constructor.ArraySortedList[T]

var _ SortedList[any] = &constructor.ArraySortedList[any]{}

// NewArray creates an empty ArraySortedList with a custom comparator.
//
// The comparator follows the same contract as [cmp.Compare]: it returns a
// negative value when a sorts before b, zero when they are equivalent, and a
// positive value when a sorts after b.
//
// ArraySortedList is optimized for read-heavy workloads: lookup is O(log N),
// indexed access is O(1), and traversal is cache-friendly. It is not ideal for
// frequent single-element updates because Add(x) and Remove(x) may shift O(N)
// elements to preserve sorted order.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Add(x)              O(N)
//	Add(xs...T)         O((N+K) log (N+K))
//	Remove(x)           O(N)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(index)          O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewArray[T any](capacity int, compare func(a, b T) int) ArraySortedList[T] {
	return constructor.NewArraySortedList(capacity, compare)
}

// ArrayFrom creates an ArraySortedList using the provided slice as storage.
//
// WARNING: This operation is in-place and will sort the original slice.
// Use [ArrayClone] when the original slice order must be preserved.
//
// This constructor is efficient when data already exists in a slice and the list
// will be queried many times after construction. Later single-element updates
// still have the same O(N) shifting cost as [NewArray].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(1)
//	Add(x)              O(N)
//	Remove(x)           O(N)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(index)          O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func ArrayFrom[T any](data []T, compare func(a, b T) int) ArraySortedList[T] {
	return constructor.NewArraySortedListFromSlice(data, compare)
}

// ArrayClone creates an ArraySortedList from a sorted clone of the provided slice.
//
// Use this constructor when the source slice must keep its original order. It is
// useful for read-heavy workflows that need a sorted view of existing data.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(N)
//	Remove(x)           O(N)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(index)          O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func ArrayClone[T any](data []T, compare func(a, b T) int) ArraySortedList[T] {
	return constructor.NewArraySortedListCloneSlice(data, compare)
}

// ArrayFromSeq creates an ArraySortedList from an iterator.
//
// This is useful for ordering any collection that exposes All(). The iterator is
// collected into new storage before sorting, so the source collection is never
// modified.
//
// Like the slice constructors, this is best for building a sorted view that will
// be queried many times and updated occasionally.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(N)
//	Remove(x)           O(N)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(index)          O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func ArrayFromSeq[T any](seq iter.Seq[T], compare func(a, b T) int) ArraySortedList[T] {
	return constructor.NewArraySortedListFromSeq(seq, compare)
}
