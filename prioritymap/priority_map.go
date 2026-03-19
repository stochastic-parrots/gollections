package prioritymap

import (
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
)

// Readonly defines a non-mutable view of a PriorityMap.
//
// It provides access to elements by key and the ability to peek at the
// highest priority element, but forbids any operation that modifies
// the map's state (like Set, Pop, or Clear).
//
// This is particularly useful for exposing the map's state to observers
// while guaranteeing data integrity.
type Readonly[K comparable, P any] interface {
	// Get returns the priority associated with key.
	//
	// If the key exists, it returns the priority and true.
	//
	// Otherwise, it returns the zero value and false.
	Get(key K) (priority P, ok bool)

	// Peek returns the key and priority of the element with the highest priority
	// without removing it.
	//
	// It returns ok as false if the map is empty.
	Peek() (key K, priority P, ok bool)

	pkg.Map[K, P]
}

// PriorityMap defines a structure that combines a map (key-based access) with a heap (priority ordering).
// It allows efficient insertion, update, and removal of key-priority pairs, as well as extraction of the minimum priority element.
// K must be comparable for map operations, and V is the priority type (requires a comparator in implementações).
type PriorityMap[K comparable, P any] interface {
	// Set inserts or updates the priority for the given key.
	//
	// If the key already exists, its priority is updated to the new value.
	Set(key K, priority P)

	// Update changes the priority of an existing key.
	//
	// If the key exists, its priority is updated to the new value and it returns true.
	// If the key does not exist, it performs no operation and returns false.
	Update(key K, priority P) (ok bool)

	// Remove deletes the key-priority pair from the map.
	//
	// It returns true if the key was found and removed.
	Remove(key K) (ok bool)

	// Pop removes and returns the key and priority with the highest priority.
	//
	// It returns ok as false if the map is empty.
	Pop() (key K, priority P, ok bool)

	// Drain returns a destructive iterator that removes and yields elements
	// in priority order (highest priority first).
	//
	// The map is emptied as the iterator progresses. If the iteration
	// is stopped early (e.g., via break), only the yielded elements
	// are removed, and the remaining ones stay in the map.
	//
	// For a non-destructive iteration, use [pkg.Map.All].
	Drain() iter.Seq2[K, P]

	// Clear removes all elements from the priority map.
	//
	// After calling Clear, the map will be empty and its length will be zero.
	// This operation is typically more efficient than creating a new map
	// as it may reuse the underlying storage.
	Clear()

	Readonly[K, P]
}

// AsReadonly returns a [Readonly] view of the provided [PriorityMap].
//
// The returned view is a wrapper that prevents type assertion back to
// the mutable interface, ensuring data safety for observers.
func AsReadonly[K comparable, P any](pm PriorityMap[K, P]) Readonly[K, P] {
	if pm == nil {
		return nil
	}
	return readonly[K, P]{inner: pm}
}

type readonly[K comparable, P any] struct {
	inner PriorityMap[K, P]
}

func (w readonly[K, P]) Get(key K) (P, bool) { return w.inner.Get(key) }

func (w readonly[K, P]) Peek() (K, P, bool) { return w.inner.Peek() }

func (w readonly[K, P]) Contains(key K) bool { return w.inner.Contains(key) }

func (w readonly[K, P]) Keys() iter.Seq[K] { return w.inner.Keys() }

func (w readonly[K, P]) Values() iter.Seq[P] { return w.inner.Values() }

func (w readonly[K, P]) All() iter.Seq2[K, P] { return w.inner.All() }

func (w readonly[K, P]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[K, P]) Length() int { return w.inner.Length() }

var _ Readonly[string, int] = (*readonly[string, int])(nil)
