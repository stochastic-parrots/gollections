package set

import (
	"encoding/json"
	"iter"

	"github.com/stochastic-parrots/gollections"
)

// Readonly defines a non-mutable view of a set.
type Readonly[T any] interface {
	gollections.Collection[T]
	json.Marshaler

	// Contains returns true if the set contains a value with the same identity
	// as x according to the concrete implementation.
	Contains(x T) bool
}

// Set defines a mutable collection of values that are unique according to its
// concrete implementation's identity policy.
type Set[T any] interface {
	Readonly[T]

	// Add inserts x if the set does not already contain a value with the same
	// identity. Existing values are preserved. It returns true if x was added.
	Add(x T) (added bool)

	// Adds inserts every value in xs that does not already exist in the set.
	// When multiple values have the same identity, the first one is preserved.
	// It returns the number of values added.
	Adds(xs ...T) (added int)

	// Remove deletes the value with the same identity as x.
	// It returns true if a value was removed.
	Remove(x T) bool

	// Removes deletes every value in xs whose identity belongs to the set.
	// It returns the number of values removed.
	Removes(xs ...T) (removed int)

	// Clear removes all values and leaves the set ready for reuse with the same
	// identity policy.
	Clear()
}

// AsReadonly returns a [Readonly] view of the provided [Set].
//
// The returned view observes the same underlying set while preventing type
// assertion back to the mutable interface. It is not a snapshot and does not
// provide synchronization.
func AsReadonly[T any](set Set[T]) *readonly[T] {
	if set == nil {
		return nil
	}
	return &readonly[T]{inner: set}
}

type readonly[T any] struct {
	inner Set[T]
}

func (w readonly[T]) Contains(x T) bool { return w.inner.Contains(x) }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

var _ Readonly[any] = (*readonly[any])(nil)
