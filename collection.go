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
	// Iteration order is defined by the concrete collection contract.
	// The collection must not be mutated while the iterator is running.
	All() iter.Seq[T]

	// Enumerate returns an iterator that yields both the zero-based iteration
	// position and the element.
	//
	// Values are yielded in the same order as All. For unordered collections,
	// the yielded position is not a stable index into the collection.
	// The collection must not be mutated while the iterator is running.
	Enumerate() iter.Seq2[int, T]
}
