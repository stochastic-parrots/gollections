// Package sortedlist provides generic lists whose order is maintained by a
// comparator.
//
// Unlike package list, sorted lists do not expose arbitrary positional mutation
// operations such as Insert, Set, or Reverse. Replace is allowed only when the
// new value preserves sorted order. The list order is derived from the
// comparator selected at construction time and is preserved after every
// mutation. [NewOrderedArray] supplies that comparator for a cmp.Ordered type in
// ascending or descending order.
//
// A sorted list is a good fit when data is built once or updated occasionally
// and then queried many times. The slice-backed implementation provides O(log N)
// lookup and O(1) indexed access, but single-element insertions and removals are
// O(N) because values may need to be shifted. For write-heavy workloads, prefer
// a heap, priority map, or a future tree/skip-list implementation depending on
// the access pattern.
//
// Use [NewOrderedArray], [OrderedArrayFrom], [OrderedArrayClone], or
// [OrderedArrayFromSeq] when T satisfies cmp.Ordered and natural order is enough.
// Use [NewArray], [ArrayFrom], [ArrayClone], or [ArrayFromSeq] for a custom
// comparator. The comparator must remain stable and define the same ordering
// for construction, lookup, replacement, and removal. The relative order of
// equivalent values is unspecified; include a tie-breaker if it matters.
//
// [ArraySortedList] requires a comparator; its zero value is invalid.
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
//		Add(xs ...T)
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
// Direct constructors build sorted lists from empty capacity, slices, cloned
// slices, or iterators. [ArrayFrom] and [OrderedArrayFrom] sort the supplied
// slice in place and transfer ownership of its backing array. The caller must
// stop using the slice and every alias afterward. Clone constructors preserve
// the source through an independent shallow copy. FromSeq constructors consume
// their iterators once. Every constructor returns an [ArraySortedList].
//
// # JSON
//
// Sorted lists marshal as arrays in sorted order. They do not implement
// json.Unmarshaler; decode into []T and pass the values to the same constructor family so
// the natural or custom ordering strategy remains explicit.
package sortedlist
