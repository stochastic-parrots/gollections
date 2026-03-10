package heap

import "iter"

// entry holds a key–priority pair stored in the underlying heap.
type entry[K comparable, T any] struct {
	key   K
	value T
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
type BinaryPriorityMap[K comparable, V any] struct {
	indexes map[K]int
	data    []entry[K, V]
	less    func(a, b V) bool
}

// NewBinaryPriorityMap creates an empty priority map with the given
// initial capacity.  The cmp function should return true if a has
// higher priority than b (for a min‑heap use a < b).
//
// Complexity: O(1).
func NewBinaryPriorityMap[K comparable, V any](
	capacity int,
	less func(V, V) bool) *BinaryPriorityMap[K, V] {

	data := make([]entry[K, V], 0, capacity)
	indexes := make(map[K]int, capacity)
	return &BinaryPriorityMap[K, V]{indexes, data, less}
}

// fixup restores the heap invariant by moving the element at idx upward.
// Used after an insertion or priority decrease.
func (m *BinaryPriorityMap[K, V]) fixup(idx int) {
	i := idx
	target := m.data[i]

	for i > 0 {
		parent := (i - 1) / 2
		if !m.less(target.value, m.data[parent].value) {
			break
		}
		m.data[i] = m.data[parent]
		m.indexes[m.data[i].key] = i
		i = parent
	}
	m.data[i] = target
	m.indexes[target.key] = i
}

// fixdown restores the heap invariant by moving the element at idx downward.
// Used after removal or priority increase.
func (m *BinaryPriorityMap[K, V]) fixdown(idx int) {
	n := len(m.data)
	i := idx
	target := m.data[i]

	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		right := left + 1
		smallest := left
		if right < n && m.less(m.data[right].value, m.data[left].value) {
			smallest = right
		}
		if !m.less(m.data[smallest].value, target.value) {
			break
		}

		m.data[i] = m.data[smallest]
		m.indexes[m.data[i].key] = i
		i = smallest
	}
	m.data[i] = target
	m.indexes[target.key] = i
}

// fix chooses between fixup and fixdown based on the element’s relation
// to its parent. It returns the final index of the element (not used yet
// but kept for symmetry with BinaryHeap.fix).
func (m *BinaryPriorityMap[K, V]) fix(idx int) int {
	if idx == 0 {
		m.fixdown(idx)
		return idx
	}
	parent := (idx - 1) / 2
	if m.less(m.data[idx].value, m.data[parent].value) {
		m.fixup(idx)
	} else {
		m.fixdown(idx)
	}
	return idx
}

// Get returns the priority associated with key.  The zero value of V and
// false are returned if the key does not exist.
//
// Complexity: O(1).
func (m *BinaryPriorityMap[K, V]) Get(key K) (V, bool) {
	if idx, ok := m.indexes[key]; ok {
		return m.data[idx].value, true
	}
	var zero V
	return zero, false
}

// Set inserts a new key with the given priority or updates an existing
// key’s priority.  The heap is adjusted accordingly.
//
// Complexity: O(log n).
func (m *BinaryPriorityMap[K, V]) Set(key K, value V) {
	if idx, ok := m.indexes[key]; ok {
		m.data[idx].value = value
		parent := (idx - 1) / 2
		if idx > 0 && m.less(value, m.data[parent].value) {
			m.fixup(idx)
		} else {
			m.fixdown(idx)
		}

		return
	}
	e := entry[K, V]{key: key, value: value}
	m.data = append(m.data, e)
	idx := len(m.data) - 1
	m.indexes[key] = idx
	m.fix(idx)
}

// Remove deletes the entry for key if present, returning true if an entry
// was removed.  The heap and index map are kept consistent.
//
// Complexity: O(log n).
func (m *BinaryPriorityMap[K, V]) Remove(key K) bool {
	if idx, exists := m.indexes[key]; exists {
		last := len(m.data) - 1
		m.data[idx], m.data[last] = m.data[last], m.data[idx]
		m.indexes[m.data[idx].key] = idx
		m.data = m.data[:last]
		delete(m.indexes, key)
		if idx < len(m.data) {
			m.fix(idx)
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
func (m *BinaryPriorityMap[K, V]) Pop() (K, V, bool) {
	n := len(m.data)
	if n == 0 {
		var zk K
		var zv V
		return zk, zv, false
	}

	e := m.data[0]
	delete(m.indexes, e.key)

	n--
	if n > 0 {
		m.data[0] = m.data[n]
		m.indexes[m.data[0].key] = 0
		m.data = m.data[:n]
		m.fixdown(0)
	} else {
		m.data = m.data[:0]
	}
	return e.key, e.value, true
}

// Peek returns the highest‑priority key‑value pair without removing it.
// If empty returns zero-key, zero-value, false.
//
// Complexity: O(1).
func (m *BinaryPriorityMap[K, V]) Peek() (K, V, bool) {
	if len(m.data) != 0 {
		entry := m.data[0]
		return entry.key, entry.value, true
	}

	var zK K
	var zv V
	return zK, zv, false
}

// Keys returns an iterator for all keys in the collection.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (m *BinaryPriorityMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range m.indexes {
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
func (m *BinaryPriorityMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, idx := range m.indexes {
			if !yield(m.data[idx].value) {
				return
			}
		}
	}
}

// All returns an iterator for key-value pairs.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (m *BinaryPriorityMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, idx := range m.indexes {
			if !yield(key, m.data[idx].value) {
				return
			}
		}
	}
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1).
func (m *BinaryPriorityMap[K, V]) IsEmpty() bool {
	return len(m.indexes) == 0
}

// Length returns the current number of elements in the heap.
//
// Complexity: O(1).
func (m *BinaryPriorityMap[K, V]) Length() int {
	return len(m.indexes)
}
