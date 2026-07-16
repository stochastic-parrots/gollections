package list

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/list"
)

// ArrayList is a slice-backed [List]. Its zero value is ready for use.
type ArrayList[T any] = list.ArrayList[T]

var _ List[any] = &list.ArrayList[any]{}

// ArrayFactory constructs slice-backed lists.
//
// Array lists provide O(1) indexed access and writes, amortized O(1) append,
// and O(N) insertion or removal at an arbitrary position.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	From(data)          O(1)
//	Clone/FromSeq       O(N)
//	Append(xs...T)      O(len(xs)) amortized
//	Insert/Remove       O(N)
//	Get/Set             O(1)
//	Find/Contains       O(N)
//	Reverse/Clear       O(N)
type ArrayFactory[T any] struct{}

// Array returns a factory for slice-backed lists.
func Array[T any]() ArrayFactory[T] {
	return ArrayFactory[T]{}
}

// New creates an empty array list with space preallocated for capacity values.
func (ArrayFactory[T]) New(capacity int) *ArrayList[T] {
	return list.NewArrayList[T](capacity)
}

// From creates an array list using data as its backing storage.
//
// The caller transfers ownership of data to the returned list and must not use
// the slice afterward. Use [ArrayFactory.Clone] to preserve the source slice.
func (ArrayFactory[T]) From(data []T) *ArrayList[T] {
	return list.NewArrayListFromSlice(data)
}

// Clone creates an array list from a shallow copy of data.
//
// Clone does not modify or retain the provided slice.
func (ArrayFactory[T]) Clone(data []T) *ArrayList[T] {
	return list.NewArrayListCloneSlice(data)
}

// FromSeq collects seq into a new array list.
func (ArrayFactory[T]) FromSeq(seq iter.Seq[T]) *ArrayList[T] {
	return list.NewArrayListFromSeq(seq)
}
