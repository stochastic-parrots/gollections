// Package sortedlist provides generic lists whose order is maintained by a comparator.
//
// Unlike package list, sorted lists do not expose positional mutation operations
// such as Insert, Set, or Reverse. Their order is derived from the comparator
// supplied at construction time and restored after every mutation.
//
// A sorted list is a good fit when data is built once or updated occasionally
// and then queried many times. The slice-backed implementation provides O(log N)
// lookup and O(1) indexed access, but single-element insertions and removals are
// O(N) because values may need to be shifted. For write-heavy workloads, prefer
// a heap, priority map, or a future tree/skip-list implementation depending on
// the access pattern.
//
// # Readonly Interface
//
// All sorted lists implement the [Readonly] interface:
//
//	type Readonly[T any] interface {
//		Get(idx int) (T, error)
//		Find(x T) (idx int, ok bool)
//		Contains(x T) bool
//		First() (T, bool)
//		Last() (T, bool)
//		Backward() iter.Seq[T]
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
//		Remove(x T) bool
//		Clear()
//		Readonly[T]
//		json.Unmarshaler
//	}
//
// # Implementations
//
// The package currently exposes one implementation:
//
//   - [ArraySortedList]: A slice-backed sorted list with O(log N) lookup and
//     O(N) single-element insertion/removal.
//
// Constructors can build sorted lists from empty capacity, slices, cloned
// slices, or iterators. Iterator constructors make it possible to sort any
// collection in this module that exposes All().
package sortedlist
