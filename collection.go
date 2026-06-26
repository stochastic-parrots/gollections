package gollections

import (
	"iter"
)

// Collection defines the basic operations for a generic data structure.
type Collection[T any] interface {
	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool

	// Length returns the number of elements in the collection.
	Length() int

	// All returns an iterator for the collection.
	All() iter.Seq[T]

	// Enumerate returns an iterator that yields both the index and the element.
	Enumerate() iter.Seq2[int, T]
}
