// Package sortedlist provides generic lists whose order is maintained by a comparator
// or by the natural order of cmp.Ordered values.
//
// Unlike package list, sorted lists do not expose arbitrary positional mutation
// operations such as Insert, Set, or Reverse. Replace is allowed only when the
// new value preserves sorted order. The list order is derived from either the
// comparator supplied at construction time or the natural order of the element
// type, and is preserved after every mutation.
//
// A sorted list is a good fit when data is built once or updated occasionally
// and then queried many times. The slice-backed implementation provides O(log N)
// lookup and O(1) indexed access, but single-element insertions and removals are
// O(N) because values may need to be shifted. For write-heavy workloads, prefer
// a heap, priority map, or a future tree/skip-list implementation depending on
// the access pattern.
//
// Use [OrderedArray] when T satisfies cmp.Ordered and natural order is enough.
// Its factory uses standard-library ordered sort and search paths without
// custom comparator calls. Use [Array] when values need a custom comparator.
// Custom comparators
// must define the same ordering for construction, lookup, replacement, and
// removal, and are called during search and sort operations. The relative order
// of values considered equivalent by the comparator is unspecified; include a
// tie-breaker in the comparator when that order matters.
//
// Bounds and range operations use list order. [Readonly.Range] is half-open:
// it yields values in [from, to), including values equivalent to from and
// excluding values equivalent to to.
//
// # Readonly Interface
//
// All sorted lists implement the [Readonly] interface:
//
//	type Readonly[T any] interface {
//		Get(idx int) (T, error)
//		LowerBound(x T) int
//		UpperBound(x T) int
//		EqualRange(x T) (start, end int)
//		Count(x T) int
//		Find(x T) (idx int, ok bool)
//		Ceiling(x T) (value T, idx int, ok bool)
//		Floor(x T) (value T, idx int, ok bool)
//		Higher(x T) (value T, idx int, ok bool)
//		Lower(x T) (value T, idx int, ok bool)
//		Contains(x T) bool
//		First() (T, bool)
//		Last() (T, bool)
//		Backward() iter.Seq[T]
//		Range(from, to T) iter.Seq[T]
//		ToSlice() []T
//		gollections.Collection[T]
//		fmt.Stringer
//		json.Marshaler
//	}
//
// # SortedList Interface
//
// Mutable sorted lists implement the [SortedList] interface:
//
//	type SortedList[T any] interface {
//		Add(xs ...T)
//		Replace(idx int, x T) error
//		Remove(x T) bool
//		Clear()
//		Readonly[T]
//		json.Unmarshaler
//	}
//
// # Implementations
//
// The package currently exposes two slice-backed implementations:
//
//   - [ArraySortedList]: A comparator-backed sorted list for custom ordering
//     with O(log N) lookup and O(N) single-element insertion/removal.
//   - [OrderedArraySortedList]: A sorted list for cmp.Ordered values that uses
//     standard-library ordered search and sort paths without custom comparator
//     calls.
//
// Factories build sorted lists from empty capacity, slices, cloned slices, or
// iterators. [ArrayFactory.From] and [OrderedArrayFactory.From] sort and retain
// the provided slice; Clone preserves it. FromSeq makes it possible to sort any
// collection in this module that exposes All(). Every factory method returns a
// concrete sorted-list type.
package sortedlist
