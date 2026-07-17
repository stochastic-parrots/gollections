package sortedlist

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections"
)

// Order selects the direction of a type's natural order.
type Order bool

const (
	// Asc places smaller values before larger values.
	Asc Order = false

	// Desc places larger values before smaller values.
	Desc Order = true
)

// Readonly defines a non-mutable view of a sorted list.
//
// It provides access to elements by sorted position and allows efficient lookup
// by value without exposing mutation operations.
type Readonly[T any] interface {
	// Get returns the element at the specified sorted index.
	//
	// Returns an error if the index is out of bounds [0, Length).
	Get(idx int) (x T, err error)

	// LowerBound returns the first index whose value does not sort before x.
	//
	// It returns Length when every value sorts before x.
	LowerBound(x T) int

	// UpperBound returns the first index whose value sorts after x.
	//
	// It returns Length when no value sorts after x.
	UpperBound(x T) int

	// EqualRange returns the half-open index range containing values equivalent to x.
	EqualRange(x T) (start, end int)

	// Count returns the number of values equivalent to x.
	Count(x T) int

	// Find locates the first index equivalent to x according to the list order.
	// It returns the index and true if found; otherwise, -1 and false.
	Find(x T) (idx int, ok bool)

	// Ceiling returns the first value that does not sort before x.
	//
	// It returns a zero value, -1, and false when every value sorts before x.
	Ceiling(x T) (value T, idx int, ok bool)

	// Floor returns the last value that does not sort after x.
	//
	// It returns a zero value, -1, and false when every value sorts after x.
	Floor(x T) (value T, idx int, ok bool)

	// Higher returns the first value that sorts after x.
	//
	// It returns a zero value, -1, and false when no value sorts after x.
	Higher(x T) (value T, idx int, ok bool)

	// Lower returns the last value that sorts before x.
	//
	// It returns a zero value, -1, and false when no value sorts before x.
	Lower(x T) (value T, idx int, ok bool)

	// Contains returns true if an equivalent value exists in the list.
	Contains(x T) bool

	// First returns the first element in sorted order.
	//
	// It returns the zero value of T and false if the list is empty.
	First() (T, bool)

	// Last returns the last element in sorted order.
	//
	// It returns the zero value of T and false if the list is empty.
	Last() (T, bool)

	// Backward returns an iterator that traverses the list in reverse sorted order.
	Backward() iter.Seq[T]

	// Range returns an iterator that traverses values in the half-open range [from, to).
	//
	// Values are selected according to list order: yielded values do not sort
	// before from and do sort before to.
	Range(from, to T) iter.Seq[T]

	// ToSlice exports the current elements of the collection into a native Go slice.
	//
	// This method creates a shallow copy of the underlying data. While the slice
	// structure itself is new, the elements (if they are pointers or reference types)
	// still point to the same memory addresses as the original collection.
	ToSlice() []T

	gollections.Collection[T]
	fmt.Stringer
	json.Marshaler
}

// SortedList defines a list whose order is derived from a comparator or natural
// element order instead of explicit positional insertion.
type SortedList[T any] interface {
	Readonly[T]

	// Add inserts x while preserving the sorted invariant.
	Add(x T)

	// Adds inserts zero or more values while preserving the sorted invariant.
	Adds(xs ...T)

	// Replace changes the value at idx while preserving the sorted invariant.
	//
	// It returns an index error when idx is out of bounds and
	// [ErrOrderViolation] when x would break list order.
	Replace(idx int, x T) error

	// Remove deletes the first value equivalent to x according to the list order.
	Remove(x T) bool

	// Clear removes all elements and leaves the list ready for reuse with the
	// same ordering configuration.
	Clear()
}

// AsReadonly returns a [Readonly] view of the provided [SortedList].
//
// The returned view is a wrapper that prevents type assertion back to
// the mutable interface, ensuring data safety for observers.
func AsReadonly[T any](list SortedList[T]) *readonly[T] {
	if list == nil {
		return nil
	}
	return &readonly[T]{inner: list}
}

type readonly[T any] struct {
	inner SortedList[T]
}

func (w readonly[T]) Get(idx int) (x T, err error) { return w.inner.Get(idx) }

func (w readonly[T]) LowerBound(x T) int { return w.inner.LowerBound(x) }

func (w readonly[T]) UpperBound(x T) int { return w.inner.UpperBound(x) }

func (w readonly[T]) EqualRange(x T) (start, end int) { return w.inner.EqualRange(x) }

func (w readonly[T]) Count(x T) int { return w.inner.Count(x) }

func (w readonly[T]) Find(x T) (idx int, ok bool) { return w.inner.Find(x) }

func (w readonly[T]) Ceiling(x T) (value T, idx int, ok bool) { return w.inner.Ceiling(x) }

func (w readonly[T]) Floor(x T) (value T, idx int, ok bool) { return w.inner.Floor(x) }

func (w readonly[T]) Higher(x T) (value T, idx int, ok bool) { return w.inner.Higher(x) }

func (w readonly[T]) Lower(x T) (value T, idx int, ok bool) { return w.inner.Lower(x) }

func (w readonly[T]) Contains(x T) bool { return w.inner.Contains(x) }

func (w readonly[T]) First() (T, bool) { return w.inner.First() }

func (w readonly[T]) Last() (T, bool) { return w.inner.Last() }

func (w readonly[T]) Backward() iter.Seq[T] { return w.inner.Backward() }

func (w readonly[T]) Range(from, to T) iter.Seq[T] { return w.inner.Range(from, to) }

func (w readonly[T]) ToSlice() []T { return w.inner.ToSlice() }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

func (w readonly[T]) String() string { return w.inner.String() }

var _ Readonly[any] = (*readonly[any])(nil)
