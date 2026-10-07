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
}

// Deque defines the operations for a double-ended queue.
type Deque[T any] interface {
	// Prepend adds the given elements to the beginning of the deque.
	//
	// The relative order of the elements is preserved at the front of the deque.
	// With no elements, it leaves the deque unchanged.
	Prepend(xs ...T)

	// Append adds the given elements to the end of the deque in input order.
	// With no elements, it leaves the deque unchanged.
	Append(xs ...T)

	// Shift removes and returns the element from the beginning of the deque.
	//
	// It returns the zero value of T and false if the deque is empty.
	Shift() (T, bool)

	// Pop removes and returns the element from the end of the deque.
	//
	// It returns the zero value of T and false if the deque is empty.
	Pop() (T, bool)

	// Clear removes all elements and leaves the deque ready for reuse.
	// Storage retention depends on the concrete deque strategy.
	Clear()

	Readonly[T]
}

// AsReadonly returns a [Readonly] view of the provided [Deque].
//
// The returned view observes the same underlying collection while preventing
// type assertion back to the mutable interface. It is not a snapshot and does
// not provide synchronization.
//
// The wrapper implements fmt.Stringer and json.Marshaler through fmt.Sprint
// and json.Marshal on the wrapped collection, following those packages'
// formatting, JSON validation, and error wrapping behavior.
func AsReadonly[T any](d Deque[T]) *readonly[T] {
	if d == nil {
		return nil
	}
	return &readonly[T]{inner: d}
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

func (w readonly[T]) Len() int { return w.inner.Len() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return json.Marshal(w.inner) }

func (w readonly[T]) String() string { return fmt.Sprint(w.inner) }

var _ Readonly[any] = (*readonly[any])(nil)
