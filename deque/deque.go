package deque

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections"
)

// Readonly defines a non-mutable view of a double-ended queue.
type Readonly[T any] interface {
	// Front returns the element at the beginning of the deque without removing it.
	//
	// It returns the zero value of T and false if the deque is empty.
	Front() (T, bool)

	// Back returns the element at the end of the deque without removing it.
	//
	// It returns the zero value of T and false if the deque is empty.
	Back() (T, bool)

	// ToSlice exports the current elements of the collection into a native Go slice.
	//
	// This method creates a shallow copy of the underlying data. While the slice
	// structure itself is new, the elements (if they are pointers or reference types)
	// still point to the same memory addresses as the original collection.
	//
	// Performance Note: This is an O(n) operation as it allocates a new slice
	// and copies each element.
	ToSlice() []T

	gollections.Collection[T]
	fmt.Stringer
	json.Marshaler
}

// Deque defines the operations for a double-ended queue, typically implemented
// via a ring buffer.
//
// It supports highly efficient O(1) insertions and removals at both ends.
type Deque[T any] interface {
	// Prepend adds the given elements to the beginning of the deque.
	//
	// If multiple elements are provided, they are inserted such that their
	// relative order is preserved at the front of the deque.
	Prepend(xs ...T)

	// Append adds the given elements to the end of the deque.
	Append(xs ...T)

	// Shift removes and returns the element from the beginning of the deque.
	//
	// It returns the zero value of T and false if the deque is empty.
	Shift() (T, bool)

	// Pop removes and returns the element from the end of the deque.
	//
	// It returns the zero value of T and false if the deque is empty.
	Pop() (T, bool)

	// Clear removes all elements from the deque, resetting it to an empty state.
	Clear()

	Readonly[T]
	json.Unmarshaler
}

// AsReadonly returns a [Readonly] view of the provided [Deque].
//
// The returned view is a wrapper that prevents type assertion back to
// the mutable interface, ensuring data safety for observers.
func AsReadonly[T any](d Deque[T]) Readonly[T] {
	if d == nil {
		return nil
	}
	return readonly[T]{inner: d}
}

type readonly[T any] struct {
	inner Deque[T]
}

func (w readonly[T]) Front() (T, bool) { return w.inner.Front() }

func (w readonly[T]) Back() (T, bool) { return w.inner.Back() }

func (w readonly[T]) ToSlice() []T { return w.inner.ToSlice() }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

func (w readonly[T]) String() string { return w.inner.String() }

var _ Readonly[any] = (*readonly[any])(nil)
