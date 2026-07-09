package sortedlist

import (
	"cmp"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// ArraySortedList is a slice-backed [SortedList] with custom comparator order.
type ArraySortedList[T any] = sortedlist.ArraySortedList[T]

// OrderedArraySortedList is a slice-backed [SortedList] for naturally ordered values.
type OrderedArraySortedList[T cmp.Ordered] = sortedlist.OrderedArraySortedList[T]

// SkipList is a probabilistic skip-list-backed [SortedList] with custom comparator order.
type SkipList[T any] = sortedlist.SkipList[T]

// OrderedSkipList is a probabilistic skip-list-backed [SortedList] for naturally ordered values.
type OrderedSkipList[T cmp.Ordered] = sortedlist.OrderedSkipList[T]

var _ SortedList[any] = &sortedlist.ArraySortedList[any]{}
var _ SortedList[int] = &sortedlist.OrderedArraySortedList[int]{}
var _ SortedList[any] = &sortedlist.SkipList[any]{}
var _ SortedList[int] = &sortedlist.OrderedSkipList[int]{}

// NewArray creates an empty ArraySortedList with a custom comparator.
//
// The comparator follows the same contract as [cmp.Compare]: it returns a
// negative value when a sorts before b, zero when they are equivalent, and a
// positive value when a sorts after b.
// It must define the same strict weak ordering used for every later lookup and
// mutation on the list. Values for which compare returns zero may appear in any
// relative order; include a tie-breaker when that order matters.
// Each search and sort comparison calls compare; prefer [NewOrderedArray] when
// T satisfies cmp.Ordered and natural order is enough.
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

// NewSkip creates an empty SkipList with a custom comparator.
//
// The comparator follows the same contract as [cmp.Compare]: it returns a
// negative value when a sorts before b, zero when they are equivalent, and a
// positive value when a sorts after b.
// It must define the same strict weak ordering used for every later lookup and
// mutation on the list. Values for which compare returns zero may appear in any
// relative order; include a tie-breaker when that order matters.
//
// SkipList is optimized for mixed lookup and update workloads: lookup,
// insertion, removal, replacement by sorted index, and indexed access are
// expected O(log N). It keeps skip-list nodes and level links in internal
// arenas, but [ArraySortedList] remains the better fit when data is built once
// and queried many times with few updates.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Add(xs...T)         O(K log (N+K))
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewSkip[T any](compare func(a, b T) int) *SkipList[T] {
	return sortedlist.NewSkipList(compare)
}

// NewSkipFrom creates a SkipList from the provided slice.
//
// WARNING: This operation is in-place and will sort the original slice.
// Use [NewSkipClone] when the original slice order must be preserved.
//
// This constructor is useful when data already exists in a slice but later
// single-element updates should avoid the O(N) shifting cost of [ArraySortedList].
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewSkipFrom[T any](data []T, compare func(a, b T) int) *SkipList[T] {
	return sortedlist.NewSkipListFromSlice(data, compare)
}

// NewSkipClone creates a SkipList from a sorted clone of the provided slice.
//
// Use this constructor when the source slice must keep its original order.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewSkipClone[T any](data []T, compare func(a, b T) int) *SkipList[T] {
	return sortedlist.NewSkipListCloneSlice(data, compare)
}

// NewSkipFromSeq creates a SkipList from an iterator.
//
// The iterator is collected into new storage before sorting, so the source
// collection is never modified.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewSkipFromSeq[T any](seq iter.Seq[T], compare func(a, b T) int) *SkipList[T] {
	return sortedlist.NewSkipListFromSeq(seq, compare)
}

// NewOrderedArray creates an empty OrderedArraySortedList with natural ordering.
//
// OrderedArraySortedList is optimized for read-heavy workloads where T already
// supports Go's ordered operations. Lookup uses the standard library's ordered
// binary search path without custom comparator calls.
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

// NewOrderedSkip creates an empty OrderedSkipList with natural ordering.
//
// OrderedSkipList is optimized for mixed lookup and update workloads where T
// already supports Go's ordered operations. It avoids custom comparator calls
// while keeping expected O(log N) lookup, insertion, removal, replacement by
// sorted index, and indexed access.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Add(xs...T)         O(K log (N+K))
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedSkip[T cmp.Ordered]() *OrderedSkipList[T] {
	return sortedlist.NewOrderedSkipList[T]()
}

// NewOrderedSkipFrom creates an OrderedSkipList from the provided slice.
//
// WARNING: This operation is in-place and will sort the original slice.
// Use [NewOrderedSkipClone] when the original slice order must be preserved.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedSkipFrom[T cmp.Ordered](data []T) *OrderedSkipList[T] {
	return sortedlist.NewOrderedSkipListFromSlice(data)
}

// NewOrderedSkipClone creates an OrderedSkipList from a sorted clone of the provided slice.
//
// Use this constructor when the source slice must keep its original order.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedSkipClone[T cmp.Ordered](data []T) *OrderedSkipList[T] {
	return sortedlist.NewOrderedSkipListCloneSlice(data)
}

// NewOrderedSkipFromSeq creates an OrderedSkipList from an iterator.
//
// The iterator is collected into new storage before sorting, so the source
// collection is never modified.
//
// Performance Summary (Expected Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Construction        O(N log N)
//	Extra Space         O(N)
//	Add(x)              O(log N)
//	Remove(x)           O(log N)
//	Replace(idx, x)     O(log N)
//	LowerBound(x)       O(log N)
//	UpperBound(x)       O(log N)
//	EqualRange(x)       O(log N)
//	Count(x)            O(log N)
//	Ceiling/Floor(x)    O(log N)
//	Higher/Lower(x)     O(log N)
//	Range(from, to)     O(log N + K)
//	Find(x)             O(log N)
//	Contains(x)         O(log N)
//	Get(idx)            O(log N)
//	First()/Last()      O(1)
//	Clear()             O(N)
func NewOrderedSkipFromSeq[T cmp.Ordered](seq iter.Seq[T]) *OrderedSkipList[T] {
	return sortedlist.NewOrderedSkipListFromSeq(seq)
}
