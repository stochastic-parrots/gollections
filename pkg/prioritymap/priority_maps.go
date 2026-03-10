package prioritymap

import "github.com/stochastic-parrots/gollections/pkg"

// PriorityMap defines a structure that combines a map (key-based access) with a heap (priority ordering).
// It allows efficient insertion, update, and removal of key-priority pairs, as well as extraction of the minimum priority element.
// K must be comparable for map operations, and V is the priority type (requires a comparator in implementações).
type PriorityMap[K comparable, V any] interface {
	// Set inserts or updates the value for a given key.
	// If the key exists, its priority is updated (decrease-key operation).
	// Parameters:
	//   key      The key to insert or update.
	//   value    The new value.
	Set(key K, value V)

	// Get retrieves the current priority for a key.
	// Returns:
	//   V     The value.
	//   bool  True if the key exists, false otherwise.
	Get(key K) (V, bool)

	// Remove deletes a key-priority pair from the map.
	// Returns:
	//   bool  True if the key was removed, false if it didn't exist.
	Remove(key K) bool

	// Pop removes and returns the key-priority pair with the priority.
	// Returns:
	//   K     The key.
	//   V     The value.
	//   bool  True if a pair was removed, false if the map is empty.
	Pop() (K, V, bool)

	// Peek returns the key-priority pair with the priority without removing it.
	// Returns:
	//   K     The key.
	//   V     The value.
	//   bool  True if the map is not empty, false otherwise.
	Peek() (K, V, bool)

	pkg.Map[K, V]
}
