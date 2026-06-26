package list

import (
	"encoding/json"
	"fmt"
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
)

// Readonly defines a non-mutable view of an indexed collection.
//
// It provides access to elements by their position and the ability to iterate
// through them, but forbids any operation that modifies the list's state.
type Readonly[T any] interface {
	// Get returns the element at the specified index.
	//
	// Returns an error if the index is out of bounds [0, Length).
	Get(idx int) (x T, err error)

	// Find locates the index of an element using a linear search.
	// It returns the index and true if found; otherwise, -1 and false.
	Find(x T, cmp func(a, b T) int) (idx int, ok bool)

	// Contains returns true if the element exists in the list according
	// to the provided comparator. This is typically an O(N) operation.
	Contains(x T, cmp func(a, b T) int) bool

	// Backward returns an iterator that traverses the list in reverse order.
	//
	// This provides a reversed view without modifying the underlying
	// structure.
	Backward() iter.Seq[T]

	// ToSlice exports the current elements of the collection into a native Go slice.
	//
	// This method creates a shallow copy of the underlying data. While the slice
	// structure itself is new, the elements (if they are pointers or reference types)
	// still point to the same memory addresses as the original collection.
	//
	// Performance Note: This is an O(n) operation as it allocates a new slice
	// and copies each element.
	ToSlice() []T

	pkg.Collection[T]
	fmt.Stringer
	json.Marshaler
}

// List defines the operations for a collection where the order is determined
// explicitly by the user (insertion order).
type List[T any] interface {
	Readonly[T]

	// Append adds the given elements to the end of the list.
	Append(xs ...T)

	// Insert places an element at a specific index, shifting subsequent
	// elements to the right.
	Insert(idx int, x T) error

	// Set replaces the element at the specified index with a new value.
	Set(idx int, x T) error

	// Remove deletes the element at the specified position and returns it.
	Remove(idx int) (x T, err error)

	// Reverse reorders the elements in the list in-place.
	Reverse()

	// Clear removes all elements from the list, resetting it to an empty state.
	Clear()

	json.Unmarshaler
}

// AsReadonly returns a [Readonly] view of the provided [List].
//
// The returned view is a wrapper that prevents type assertion back to
// the mutable interface, ensuring data safety for observers.
func AsReadonly[T any](list List[T]) *readonly[T] {
	if list == nil {
		return nil
	}
	return &readonly[T]{inner: list}
}

type readonly[T any] struct {
	inner List[T]
}

func (w readonly[T]) Get(idx int) (x T, err error) { return w.inner.Get(idx) }

func (w readonly[T]) Find(x T, cmp func(a, b T) int) (idx int, ok bool) { return w.inner.Find(x, cmp) }

func (w readonly[T]) Contains(x T, cmp func(a, b T) int) bool { return w.inner.Contains(x, cmp) }

func (w readonly[T]) Backward() iter.Seq[T] { return w.inner.Backward() }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) ToSlice() []T { return w.inner.ToSlice() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

func (w readonly[T]) String() string { return w.inner.String() }

var _ Readonly[any] = (*readonly[any])(nil)
