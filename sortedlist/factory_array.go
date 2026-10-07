package sortedlist

import (
	"cmp"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// ArraySortedList is a slice-backed SortedList. Its zero value is invalid;
// use NewArray or NewOrderedArray.
type ArraySortedList[T any] = sortedlist.ArraySortedList[T]

var _ SortedList[int] = &sortedlist.ArraySortedList[int]{}

func checkCompare[T any](compare func(T, T) int) {
	if compare == nil {
		panic("sortedlist: nil comparator")
	}
}

func orderedCompare[T cmp.Ordered](order Order) func(T, T) int {
	if order == Desc {
		return func(a, b T) int { return cmp.Compare(b, a) }
	}
	return cmp.Compare[T]
}

// NewArray creates an empty sorted list using a stable cmp.Compare-compatible
// comparator. Choose it for O(log N) lookup and O(1) indexed access when
// single-value insertion or removal at O(N) cost is acceptable.
// The comparator must remain stable for the lifetime of the list.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(capacity).
func NewArray[T any](compare func(T, T) int, capacity int) *ArraySortedList[T] {
	checkCompare(compare)
	return sortedlist.NewArraySortedList(capacity, compare)
}

// ArrayFrom sorts data in place and transfers ownership of its backing array.
// The caller must stop using data and every alias afterward. Use ArrayClone to preserve it.
// The comparator must be non-nil and provide a stable cmp.Compare-compatible order.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data) log len(data)).
func ArrayFrom[T any](compare func(T, T) int, data []T) *ArraySortedList[T] {
	checkCompare(compare)
	return sortedlist.NewArraySortedListFromSlice(data, compare)
}

// ArrayClone sorts an independent shallow copy of data and does not modify or retain its backing array.
// The comparator must be non-nil and provide a stable cmp.Compare-compatible order.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data) log len(data)).
func ArrayClone[T any](compare func(T, T) int, data []T) *ArraySortedList[T] {
	checkCompare(compare)
	return sortedlist.NewArraySortedListCloneSlice(data, compare)
}

// ArrayFromSeq consumes seq once, then sorts new storage using compare.
// The comparator must be non-nil and provide a stable cmp.Compare-compatible order.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(N log N), where N is the number of values yielded by seq.
func ArrayFromSeq[T any](compare func(T, T) int, seq iter.Seq[T]) *ArraySortedList[T] {
	checkCompare(compare)
	return sortedlist.NewArraySortedListFromSeq(seq, compare)
}

// NewOrderedArray creates an empty sorted list using T's natural order in the selected direction.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(capacity).
func NewOrderedArray[T cmp.Ordered](order Order, capacity int) *ArraySortedList[T] {
	return NewArray(orderedCompare[T](order), capacity)
}

// OrderedArrayFrom sorts data in place using T's natural order and transfers its backing array.
// The caller must stop using data and every alias afterward. Use OrderedArrayClone to preserve it.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data) log len(data)).
func OrderedArrayFrom[T cmp.Ordered](order Order, data []T) *ArraySortedList[T] {
	return ArrayFrom(orderedCompare[T](order), data)
}

// OrderedArrayClone sorts an independent shallow copy of data using T's natural order.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data) log len(data)).
func OrderedArrayClone[T cmp.Ordered](order Order, data []T) *ArraySortedList[T] {
	return ArrayClone(orderedCompare[T](order), data)
}

// OrderedArrayFromSeq consumes seq once, then sorts the values using T's natural order.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Get(idx)          O(1)
//	LowerBound(x)     O(log N)
//	UpperBound(x)     O(log N)
//	EqualRange(x)     O(log N)
//	Count(x)          O(log N)
//	Find(x)           O(log N)
//	Ceiling(x)        O(log N)
//	Floor(x)          O(log N)
//	Higher(x)         O(log N)
//	Lower(x)          O(log N)
//	Contains(x)       O(log N)
//	First()           O(1)
//	Last()            O(1)
//	Backward()        O(N)
//	Range(from, to)   O(log N + yielded values)
//	ToSlice()         O(N)
//	Add(xs... T)      O((N + len(xs)) log (N + len(xs)))
//	Remove(x)         O(N)
//	Replace(idx, x)   O(1)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(N log N), where N is the number of values yielded by seq.
func OrderedArrayFromSeq[T cmp.Ordered](order Order, seq iter.Seq[T]) *ArraySortedList[T] {
	return ArrayFromSeq(orderedCompare[T](order), seq)
}
