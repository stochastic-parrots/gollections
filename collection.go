package gollections

import (
	"iter"
)

// Collection defines the basic operations for a generic data structure.
//
// Iterators returned by a Collection are lazy and observe the collection when
// iteration begins. Callers must not mutate the collection while an iterator is
// running unless the iterator explicitly documents destructive behavior.
type Collection[T any] interface {
	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool

	// Length returns the number of elements in the collection.
	Length() int

	// All returns an iterator for the collection.
	//
	// The collection must not be mutated while the iterator is running.
	All() iter.Seq[T]

	// Enumerate returns an iterator that yields both the index and the element.
	//
	// The collection must not be mutated while the iterator is running.
	Enumerate() iter.Seq2[int, T]
}
