package pkg

import (
	"fmt"
)

// Collection defines the basic operations for a generic data structure.
type Collection[T any] interface {
	// IsEmpty Checks if the collection is empty.
	IsEmpty() bool

	// Length Returns the collection length (number of elements).
	Length() int

	// Reverse the collection in place.
	Reverse()

	// Iterator Returns an iterator for the collection.
	Iterator() Iterator[T]

	fmt.Stringer
}
