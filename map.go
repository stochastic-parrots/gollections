package gollections

import "iter"

// Map defines the base operations for key-value structures.
type Map[K comparable, V any] interface {
	// Get retrieves the current value for a key.
	Get(key K) (value V, ok bool)

	// Contains returns true if the key exists in the map.
	Contains(key K) bool

	// Keys returns an iterator for all keys in the collection.
	Keys() iter.Seq[K]

	// Values returns an iterator for all values (priorities/data).
	Values() iter.Seq[V]

	// All returns an iterator for key-value pairs.
	All() iter.Seq2[K, V]

	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool

	// Length Returns the collection length (number of elements).
	Length() int
}
