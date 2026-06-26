package sortedlist

import (
	"cmp"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// ArraySortedList is a slice-backed [SortedList].
type ArraySortedList[T any] = sortedlist.ArraySortedList[T]

// OrderedArraySortedList is a slice-backed [SortedList] for naturally ordered values.
type OrderedArraySortedList[T cmp.Ordered] = sortedlist.OrderedArraySortedList[T]

var _ SortedList[any] = &sortedlist.ArraySortedList[any]{}
var _ SortedList[int] = &sortedlist.OrderedArraySortedList[int]{}

// NewArray creates an empty ArraySortedList with a custom comparator.
//
// The comparator follows the same contract as [cmp.Compare]: it returns a
// negative value when a sorts before b, zero when they are equivalent, and a
// positive value when a sorts after b.
// It must define the same strict weak ordering used for every later lookup and
// mutation on the list. Values for which compare returns zero may appear in any
// relative order; include a tie-breaker when that order matters.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewArray[T any](capacity int, compare func(a, b T) int) *ArraySortedList[T] {
	return sortedlist.NewArraySortedList(capacity, compare)
}

// NewArrayFrom creates an ArraySortedList using the provided slice as storage.
//
// WARNING: This operation is in-place and will sort the original slice.
// Use [NewArrayClone] when the original slice order must be preserved.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewArrayFrom[T any](data []T, compare func(a, b T) int) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListFromSlice(data, compare)
}

// NewArrayClone creates an ArraySortedList from a sorted clone of the provided slice.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewArrayClone[T any](data []T, compare func(a, b T) int) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListCloneSlice(data, compare)
}

// NewArrayFromSeq creates an ArraySortedList from an iterator.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewArrayFromSeq[T any](seq iter.Seq[T], compare func(a, b T) int) *ArraySortedList[T] {
	return sortedlist.NewArraySortedListFromSeq(seq, compare)
}

// NewOrderedArray creates an empty OrderedArraySortedList with natural ordering.
//
// OrderedArraySortedList is optimized for read-heavy workloads where T already
// supports Go's ordered operations. Lookup uses the standard library's ordered
// binary search path without a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Add(x)              O(N)
//	Add(xs...T)         O((N+K) log (N+K))
//	Remove(x)           O(N)
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedArray[T cmp.Ordered](capacity int) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedList[T](capacity)
}

// NewOrderedArrayFrom creates an OrderedArraySortedList using the provided slice as storage.
//
// WARNING: This operation is in-place and will sort the original slice.
// Use [NewOrderedArrayClone] when the original slice order must be preserved.
//
// This constructor is efficient when data already exists in a slice and the list
// will be queried many times after construction. Later single-element updates
// still have the same O(N) shifting cost as [NewOrderedArray].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(1)
//	Add(x)              O(N)
//	Remove(x)           O(N)
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedArrayFrom[T cmp.Ordered](data []T) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListFromSlice(data)
}

// NewOrderedArrayClone creates an OrderedArraySortedList from a sorted clone of the provided slice.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedArrayClone[T cmp.Ordered](data []T) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListCloneSlice(data)
}

// NewOrderedArrayFromSeq creates an OrderedArraySortedList from an iterator.
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
//	Replace(idx, x)     O(1)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(1)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedArrayFromSeq[T cmp.Ordered](seq iter.Seq[T]) *OrderedArraySortedList[T] {
	return sortedlist.NewOrderedArraySortedListFromSeq(seq)
}
