// Package sortedlist provides generic lists whose order is maintained by a
// comparator.
//
// Unlike package list, sorted lists do not expose arbitrary positional mutation
// operations such as Insert, Set, or Reverse. Replace is allowed only when the
// new value preserves sorted order. The list order is derived from the
// comparator selected at construction time and is preserved after every
// mutation. [OrderedArray] supplies that comparator for a cmp.Ordered type in
// ascending or descending order.
//
// A sorted list is a good fit when data is built once or updated occasionally
// and then queried many times. The slice-backed implementation provides O(log N)
// lookup and O(1) indexed access, but single-element insertions and removals are
// O(N) because values may need to be shifted. For write-heavy workloads, prefer
// a heap, priority map, or a future tree/skip-list implementation depending on
// the access pattern.
//
// Use [OrderedArray] when T satisfies cmp.Ordered and ascending or descending
// natural order is enough. Use [Array] when values need a custom comparator.
// Both selectors return [ArrayFactory] and use the same implementation. The
// comparator must define the same ordering for construction, lookup,
// replacement, and removal. The relative order of values considered equivalent
// by the comparator is unspecified; include a tie-breaker when that order
// matters.
//
// [ArraySortedList] requires a comparator and must be constructed through
// [Array] or [OrderedArray]; its zero value is invalid.
//
// Bounds and range operations use list order. [Readonly.Range] is half-open:
// it yields values in [from, to), including values equivalent to from and
// excluding values equivalent to to.
// All traverses values in sorted order, Enumerate uses the same order, and
// Backward traverses it in reverse.
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
// A readonly view observes the same underlying sorted list. It restricts
// mutation through that interface but is not a snapshot or concurrency
// mechanism.
//
// # SortedList Interface
//
// Mutable sorted lists implement the [SortedList] interface:
//
//	type SortedList[T any] interface {
//		Add(x T)
//		Adds(xs ...T)
//		Replace(idx int, x T) error
//		Remove(x T) bool
//		Clear()
//		Readonly[T]
//	}
//
// # Implementations
//
// The package exposes one slice-backed implementation:
//
//   - [ArraySortedList]: A comparator-backed sorted list with O(log N) lookup
//     and O(N) single-element insertion/removal.
//
// Factories build sorted lists from empty capacity, slices, cloned slices, or
// iterators. [ArrayFactory.From] sorts and retains the provided slice,
// transferring ownership of its backing array. The caller must not use the
// slice or its aliases afterward. Clone preserves the source. FromSeq makes it
// possible to sort any collection in this module that exposes All(). Every
// factory method returns an [ArraySortedList].
//
// # JSON
//
// Sorted lists marshal as arrays in sorted order. They do not implement
// json.Unmarshaler; decode into []T and pass the values to the same factory so
// the natural or custom ordering strategy remains explicit.
package sortedlist
