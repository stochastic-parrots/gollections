package lists

import (
	"github.com/stochastic-parrots/gollections/pkg"
)

// List defines the operations for an ordered, indexable collection of elements.
// It extends the basic functionalities provided by the gollections.Collection interface.
type List[T any] interface {
	// Append adds the given elements (xs) to the end of the list.
	Append(...T)

	// Get returns the element at the specified index.
	// Returns an error if the index is out of bounds.
	//
	// Parameters:
	//   index  The zero-based index of the element to retrieve.
	//
	// Errors:
	//   ErrIndexOutOfBound: Returned if 'index' is negative or greater than or equal to Length().
	Get(index int) (T, error)

	// Set replaces the element at the specified index with the new value (x).
	// Returns an error if the index is out of bounds.
	//
	// Parameters:
	//   index  The zero-based index of the element to modify.
	//   x      The new value to set at the given index.
	//
	// Errors:
	//   - [ErrIndexOutOfBound](IndexOutOfBound): Returned if 'index' is negative or greater than or equal to Length().
	Set(index int, x T) error

	// Reverse the list in place.
	Reverse()

	pkg.Collection[T]
}
