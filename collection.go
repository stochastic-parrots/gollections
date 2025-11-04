package gollections

import "fmt"

// Collection defines the basic operations for a generic data structure.
type Collection[T any] interface {
	// Checks if the collection is empty.
	IsEmpty() bool

	// Returns the collection length (number of elements).
	Length() int

	// Reverse the collection in place.
	Reverse()

	// Returns an iterator for the collection.
	Iterator() Iterator[T]

	fmt.Stringer
}
