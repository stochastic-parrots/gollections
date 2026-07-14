package gollections

import "iter"

// Map defines the base operations for key-value structures.
type Map[K comparable, V any] interface {
	// Get retrieves the current value for a key.
	Get(key K) (value V, ok bool)

	// Contains returns true if the key exists in the map.
	Contains(key K) bool

	// Keys returns an iterator for all keys in unspecified order.
	Keys() iter.Seq[K]

	// Values returns an iterator for all values in unspecified order.
	//
	// A separate Keys iteration is not guaranteed to use the same order. Use All
	// when each value must remain associated with its key.
	Values() iter.Seq[V]

	// All returns an iterator for key-value pairs in unspecified order.
	All() iter.Seq2[K, V]

	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool

	// Length returns the number of key-value pairs in the map.
	Length() int
}
