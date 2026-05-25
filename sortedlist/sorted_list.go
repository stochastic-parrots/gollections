package sortedlist

import (
	"encoding/json"
	"fmt"
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
)

// Readonly defines a non-mutable view of a comparator-sorted list.
//
// It provides access to elements by sorted position and allows efficient lookup
// by value without exposing mutation operations.
type Readonly[T any] interface {
	// Get returns the element at the specified sorted index.
	//
	// Returns an error if the index is out of bounds [0, Length).
	Get(idx int) (x T, err error)

	// Find locates the first index equivalent to x according to the list comparator.
	// It returns the index and true if found; otherwise, -1 and false.
	Find(x T) (idx int, ok bool)

	// Contains returns true if an equivalent value exists in the list.
	Contains(x T) bool

	// First returns the first element in comparator order.
	//
	// It returns the zero value of T and false if the list is empty.
	First() (T, bool)

	// Last returns the last element in comparator order.
	//
	// It returns the zero value of T and false if the list is empty.
	Last() (T, bool)

	// Backward returns an iterator that traverses the list in reverse sorted order.
	Backward() iter.Seq[T]

	// ToSlice exports the current elements of the collection into a native Go slice.
	//
	// This method creates a shallow copy of the underlying data. While the slice
	// structure itself is new, the elements (if they are pointers or reference types)
	// still point to the same memory addresses as the original collection.
	ToSlice() []T

	pkg.Collection[T]
	fmt.Stringer
	json.Marshaler
}

// SortedList defines a list whose order is derived from a comparator instead of
// explicit positional insertion.
type SortedList[T any] interface {
	Readonly[T]

	// Add inserts the given elements while preserving the sorted invariant.
	Add(xs ...T)

	// Remove deletes the first value equivalent to x according to the comparator.
	Remove(x T) bool

	// Clear removes all elements from the list, resetting it to an empty state.
	Clear()

	json.Unmarshaler
}

// AsReadonly returns a [Readonly] view of the provided [SortedList].
//
// The returned view is a wrapper that prevents type assertion back to
// the mutable interface, ensuring data safety for observers.
func AsReadonly[T any](list SortedList[T]) Readonly[T] {
	if list == nil {
		return nil
	}
	return readonly[T]{inner: list}
}

type readonly[T any] struct {
	inner SortedList[T]
}

func (w readonly[T]) Get(idx int) (x T, err error) { return w.inner.Get(idx) }

func (w readonly[T]) Find(x T) (idx int, ok bool) { return w.inner.Find(x) }

func (w readonly[T]) Contains(x T) bool { return w.inner.Contains(x) }

func (w readonly[T]) First() (T, bool) { return w.inner.First() }

func (w readonly[T]) Last() (T, bool) { return w.inner.Last() }

func (w readonly[T]) Backward() iter.Seq[T] { return w.inner.Backward() }

func (w readonly[T]) ToSlice() []T { return w.inner.ToSlice() }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

func (w readonly[T]) String() string { return w.inner.String() }

var _ Readonly[any] = (*readonly[any])(nil)
