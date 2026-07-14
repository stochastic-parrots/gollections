package gollections

import "iter"

// Map defines the base operations for key-value structures.
//
// Iterators returned by a Map are lazy and observe the map when iteration
// begins. Callers must not mutate the map while an iterator is running unless
// the iterator explicitly documents destructive behavior.
type Map[K comparable, V any] interface {
	// Get retrieves the current value for a key.
	Get(key K) (value V, ok bool)

	// Contains returns true if the key exists in the map.
	Contains(key K) bool

	// Keys returns an iterator for all keys in the collection.
	//
	// The map must not be mutated while the iterator is running.
	Keys() iter.Seq[K]

	// Values returns an iterator for all values (priorities/data).
	//
	// The map must not be mutated while the iterator is running.
	Values() iter.Seq[V]

	// All returns an iterator for key-value pairs.
	//
	// The map must not be mutated while the iterator is running.
	All() iter.Seq2[K, V]

	// IsEmpty returns true if the collection is empty.
	IsEmpty() bool

	// Length returns the number of key-value pairs in the map.
	Length() int
}
