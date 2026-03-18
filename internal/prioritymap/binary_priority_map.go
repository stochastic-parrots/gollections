package prioritymap

import "iter"

// entry holds a key–priority pair stored in the underlying heap.
type entry[K comparable, P any] struct {
	key      K
	priority P
}

// BinaryPriorityMap implements a priority‑map on top of a binary heap.
// It maintains a map from keys to positions, allowing O(log n) insertion,
// update and removal by key, and O(1) lookups.
//
// K must be comparable (so it can be used as a map key). V is the
// priority type; ordering among values is determined by the comparator
// passed to NewBinaryPriorityHeap.
//
// Internally the structure keeps two parallel representations:
//   - data      – a slice forming the binary heap, sorted by priority;
//   - indexes   – a map[K]int mapping each key to its index in the slice.
//
// All heap movements update the map via the swap/fix helpers below.
type BinaryPriorityMap[K comparable, P any] struct {
	indexes     map[K]int
	data        []entry[K, P]
	hasPriority func(P, P) bool
}

// NewBinaryPriorityMap creates an empty priority map with the given
// initial capacity.  The cmp function should return true if a has
// higher priority than b (for a min‑heap use a < b).
//
// Complexity: O(1).
func NewBinaryPriorityMap[K comparable, P any](
	capacity int,
	hasPriority func(P, P) bool) *BinaryPriorityMap[K, P] {

	data := make([]entry[K, P], 0, capacity)
	indexes := make(map[K]int, capacity)
	return &BinaryPriorityMap[K, P]{indexes, data, hasPriority}
}

// fixup restores the heap invariant by moving the element at idx upward.
// Used after an insertion or priority decrease.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) fixup(idx int) {
	i := idx
	target := pm.data[i]

	for i > 0 {
		parent := (i - 1) / 2
		if !pm.hasPriority(target.priority, pm.data[parent].priority) {
			break
		}
		pm.data[i] = pm.data[parent]
		pm.indexes[pm.data[i].key] = i
		i = parent
	}
	pm.data[i] = target
	pm.indexes[target.key] = i
}

// fixdown restores the heap invariant by moving the element at idx downward.
// Used after removal or priority increase.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) fixdown(idx int) {
	n := len(pm.data)
	i := idx
	target := pm.data[i]

	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		right := left + 1
		smallest := left
		if right < n && pm.hasPriority(pm.data[right].priority, pm.data[left].priority) {
			smallest = right
		}
		if !pm.hasPriority(pm.data[smallest].priority, target.priority) {
			break
		}

		pm.data[i] = pm.data[smallest]
		pm.indexes[pm.data[i].key] = i
		i = smallest
	}
	pm.data[i] = target
	pm.indexes[target.key] = i
}

// fix chooses between fixup and fixdown based on the element’s relation
// to its parent. It returns the final index of the element (not used yet
// but kept for symmetry with BinaryHeap.fix).
func (pm *BinaryPriorityMap[K, P]) fix(idx int) int {
	if idx == 0 {
		pm.fixdown(idx)
		return idx
	}
	parent := (idx - 1) / 2
	if pm.hasPriority(pm.data[idx].priority, pm.data[parent].priority) {
		pm.fixup(idx)
	} else {
		pm.fixdown(idx)
	}
	return idx
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1).
func (pm *BinaryPriorityMap[K, P]) IsEmpty() bool {
	return len(pm.indexes) == 0
}

// Length returns the current number of elements in the heap.
//
// Complexity: O(1).
func (pm *BinaryPriorityMap[K, P]) Length() int {
	return len(pm.indexes)
}

// Get returns the priority associated with key.  The zero value of V and
// false are returned if the key does not exist.
//
// Complexity: O(1).
func (pm *BinaryPriorityMap[K, P]) Get(key K) (priority P, ok bool) {
	if idx, ok := pm.indexes[key]; ok {
		return pm.data[idx].priority, true
	}
	var zP P
	return zP, false
}

// Set inserts a new key with the given priority or updates an existing
// key’s priority.  The heap is adjusted accordingly.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) Set(key K, priority P) {
	if idx, ok := pm.indexes[key]; ok {
		pm.data[idx].priority = priority
		parent := (idx - 1) / 2
		if idx > 0 && pm.hasPriority(priority, pm.data[parent].priority) {
			pm.fixup(idx)
		} else {
			pm.fixdown(idx)
		}
		return
	}
	e := entry[K, P]{key: key, priority: priority}
	pm.data = append(pm.data, e)
	idx := len(pm.data) - 1
	pm.indexes[key] = idx
	pm.fix(idx)
}

// Update changes the priority of an existing key.
//
// It returns true if the key was found and updated. If the key does not exist,
// it performs no operation and returns false.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) Update(key K, priority P) (ok bool) {
	if idx, exists := pm.indexes[key]; exists {
		pm.data[idx].priority = priority
		parent := (idx - 1) / 2
		if idx > 0 && pm.hasPriority(priority, pm.data[parent].priority) {
			pm.fixup(idx)
		} else {
			pm.fixdown(idx)
		}
		return true
	}
	return false
}

// Remove deletes the entry for key if present, returning true if an entry
// was removed.  The heap and index map are kept consistent.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) Remove(key K) bool {
	if idx, exists := pm.indexes[key]; exists {
		last := len(pm.data) - 1
		pm.data[idx], pm.data[last] = pm.data[last], pm.data[idx]
		pm.indexes[pm.data[idx].key] = idx
		pm.data = pm.data[:last]
		delete(pm.indexes, key)
		if idx < len(pm.data) {
			pm.fix(idx)
		}
		return true
	}

	return false
}

// Pop removes and returns the key‑priority pair with the highest
// priority (the root of the heap).  If the map is empty it returns
// zero-key, zero-value and false.
//
// Complexity: O(log n).
func (pm *BinaryPriorityMap[K, P]) Pop() (key K, priority P, ok bool) {
	n := len(pm.data)
	if n == 0 {
		var zK K
		var zP P
		return zK, zP, false
	}

	e := pm.data[0]
	delete(pm.indexes, e.key)

	n--
	if n > 0 {
		pm.data[0] = pm.data[n]
		pm.indexes[pm.data[0].key] = 0
		pm.data = pm.data[:n]
		pm.fixdown(0)
	} else {
		pm.data = pm.data[:0]
	}
	return e.key, e.priority, true
}

// Peek returns the highest‑priority key‑value pair without removing it.
// If empty returns zero-key, zero-value, false.
//
// Complexity: O(1).
func (pm *BinaryPriorityMap[K, P]) Peek() (key K, priority P, ok bool) {
	if len(pm.data) != 0 {
		entry := pm.data[0]
		return entry.key, entry.priority, true
	}

	var zK K
	var zP P
	return zK, zP, false
}

// Contains returns true if the key exists in the map.
//
// Complexity: O(1).
func (pm *BinaryPriorityMap[K, P]) Contains(key K) bool {
	_, exists := pm.indexes[key]
	return exists
}

// Keys returns an iterator for all keys in the collection.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (pm *BinaryPriorityMap[K, P]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range pm.indexes {
			if !yield(key) {
				return
			}
		}
	}
}

// Values returns an iterator for all values (priorities/data).
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (pm *BinaryPriorityMap[K, P]) Values() iter.Seq[P] {
	return func(yield func(P) bool) {
		for _, idx := range pm.indexes {
			if !yield(pm.data[idx].priority) {
				return
			}
		}
	}
}

// All returns an iterator for key-value pairs.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (pm *BinaryPriorityMap[K, P]) All() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for key, idx := range pm.indexes {
			if !yield(key, pm.data[idx].priority) {
				return
			}
		}
	}
}

// Drain returns a destructive iterator that removes and yields elements
// in priority order (highest priority first).
//
// Since this is a destructive operation, the map will be empty after a
// full traversal. If the iteration is stopped early (e.g., via break),
// the map will retain only the remaining elements.
//
// Complexity: O(n log n) for a full traversal.
func (pm *BinaryPriorityMap[K, P]) Drain() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for {
			key, priority, ok := pm.Pop()
			if !ok {
				break
			}
			if !yield(key, priority) {
				break
			}
		}
	}
}

// Clear removes all elements from the priority map.
//
// After calling Clear, the map will be empty and its length will be zero.
// This operation is typically more efficient than creating a new map
// as it may reuse the underlying storage.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (pm *BinaryPriorityMap[K, P]) Clear() {
	clear(pm.indexes)
	clear(pm.data)
	pm.data = pm.data[:0]
}
