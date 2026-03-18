package prioritymap

import (
	"iter"
)

// node holds a key–priority pair and the necessary pointers for the
// multi-way pairing heap structure.
type node[K comparable, P any] struct {
	key      K
	priority P

	child    *node[K, P]
	next     *node[K, P]
	previous *node[K, P]
}

// PairingPriorityMap implements a priority‑map on top of a pairing heap.
// It maintains a map from keys to node pointers, allowing O(1) amortized
// insertion and priority updates (decrease-key), and O(log n) amortized
// removal and extraction of the highest priority element.
//
// K must be comparable. P is the priority type; ordering is determined
// by the hasPriority function passed to the constructor.
//
// Internally, the structure uses a multi-way tree (pairing heap) where
// each node is indexed in a map for instant access. It also includes a
// freelist to reuse node allocations and mitigate GC pressure.
type PairingPriorityMap[K comparable, P any] struct {
	hasPriority func(P, P) bool
	root        *node[K, P]
	indexes     map[K]*node[K, P]
	freelist    *node[K, P]
}

// NewPairingPriorityMap creates an empty pairing priority map.
// The hasPriority function should return true if a has higher priority than b.
//
// Complexity: O(1).
func NewPairingPriorityMap[K comparable, P any](hasPriority func(P, P) bool) *PairingPriorityMap[K, P] {
	return &PairingPriorityMap[K, P]{
		hasPriority: hasPriority,
		indexes:     make(map[K]*node[K, P]),
	}
}

// NewPairingPriorityMapWithCapacity creates an empty pairing priority map with
// pre-allocated nodes and map capacity. This is recommended for high-performance
// scenarios to reduce allocations during initial bursts.
//
// Complexity: O(capacity).
func NewPairingPriorityMapWithCapacity[K comparable, P any](
	capacity int,
	hasPriority func(P, P) bool) *PairingPriorityMap[K, P] {

	pm := new(PairingPriorityMap[K, P])
	indexes := make(map[K]*node[K, P], capacity)
	storage := make([]node[K, P], capacity)

	for i := range capacity {
		n := &storage[i]
		n.next = pm.freelist
		pm.freelist = n
	}

	pm.hasPriority = hasPriority
	pm.root = nil
	pm.indexes = indexes
	return pm
}

// allocate retrieves a node from the freelist or creates a new one.
func (pm *PairingPriorityMap[K, P]) allocate(key K, priority P) *node[K, P] {
	var n *node[K, P]
	if pm.freelist == nil {
		n = &node[K, P]{}
	} else {
		n = pm.freelist
		pm.freelist = pm.freelist.next
	}

	n.key = key
	n.priority = priority
	n.child, n.previous, n.next = nil, nil, nil
	return n
}

// deallocate clears a node's data and returns it to the freelist.
func (pm *PairingPriorityMap[K, P]) deallocate(address *node[K, P]) {
	var zK K
	var zP P
	address.key, address.priority = zK, zP
	address.child, address.previous, address.next = nil, nil, nil

	address.next = pm.freelist
	pm.freelist = address
}

// merge unites two sub-heaps. The node with higher priority becomes the
// parent of the other.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) merge(first, second *node[K, P]) *node[K, P] {
	if first == nil {
		return second
	}

	if second == nil {
		return first
	}

	if pm.hasPriority(second.priority, first.priority) {
		first, second = second, first
	}

	second.next = first.child
	if first.child != nil {
		first.child.previous = second
	}

	first.child = second
	second.previous = first
	return first
}

// cut removes a node from its current position in the tree, effectively
// isolating it as a standalone sub-heap.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) cut(n *node[K, P]) {
	if n == pm.root || n.previous == nil {
		return
	}

	if n.previous.child == n {
		n.previous.child = n.next
	} else {
		n.previous.next = n.next
	}

	if n.next != nil {
		n.next.previous = n.previous
	}

	n.next = nil
	n.previous = nil
}

// combine merges a list of siblings into a single heap using a two-pass
// strategy (left-to-right pairs, then right-to-left reduction).
//
// Complexity: O(log n) amortized.
func (pm *PairingPriorityMap[K, P]) combine(first *node[K, P]) *node[K, P] {
	if first == nil {
		return nil
	}

	if first.next == nil {
		return first
	}

	var head *node[K, P]
	current := first
	for current != nil {
		a := current
		b := current.next
		if b != nil {
			nextGroup := b.next
			a.next, a.previous = nil, nil
			b.next, b.previous = nil, nil

			res := pm.merge(a, b)
			res.next = head
			head = res
			current = nextGroup
		} else {
			a.next, a.previous = head, nil
			head = a
			current = nil
		}
	}

	var root *node[K, P]
	current = head
	for current != nil {
		next := current.next
		current.next = nil
		root = pm.merge(current, root)
		current = next
	}
	return root
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) IsEmpty() bool {
	return len(pm.indexes) == 0
}

// Length returns the current number of elements in the heap.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) Length() int {
	return len(pm.indexes)
}

// Get returns the priority associated with key.
//
// The zero value of P and false are returned if the key does not exist.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) Get(key K) (P, bool) {
	if node, ok := pm.indexes[key]; ok {
		return node.priority, true
	}
	var zero P
	return zero, false
}

// Set inserts a new key with the given priority or updates an existing
// key’s priority. If the priority is improved, it is an O(1) operation.
// If it is worsened, it is O(log n).
//
// Complexity: O(1) amortized for insertions and priority improvements.
func (pm *PairingPriorityMap[K, P]) Set(key K, priority P) {
	n, exists := pm.indexes[key]
	if !exists {
		newNode := pm.allocate(key, priority)
		pm.indexes[key] = newNode
		pm.root = pm.merge(pm.root, newNode)
		return
	}

	old := n.priority
	n.priority = priority
	if pm.hasPriority(priority, old) || (!pm.hasPriority(old, priority)) {
		if n != pm.root {
			pm.cut(n)
			pm.root = pm.merge(pm.root, n)
		}
		return
	}

	if n == pm.root {
		children := n.child
		n.child = nil
		pm.root = n
		if children != nil {
			pm.root = pm.merge(pm.root, pm.combine(children))
		}
	} else {
		pm.cut(n)
		children := n.child
		n.child = nil
		pm.root = pm.merge(pm.root, n)
		if children != nil {
			pm.root = pm.merge(pm.root, pm.combine(children))
		}
	}
}

// Remove deletes the entry for key if present, returning true if an entry
// was removed.
//
// Complexity: O(log n) amortized.
func (pm *PairingPriorityMap[K, P]) Remove(key K) bool {
	node, ok := pm.indexes[key]
	if !ok {
		return false
	}

	if node == pm.root {
		pm.Pop()
		return true
	}

	pm.cut(node)
	delete(pm.indexes, key)

	subtree := pm.combine(node.child)
	pm.root = pm.merge(pm.root, subtree)
	pm.deallocate(node)

	return true
}

// Pop removes and returns the key‑priority pair with the highest
// priority. If empty, it returns zero values and false.
//
// Complexity: O(log n) amortized.
func (pm *PairingPriorityMap[K, P]) Pop() (K, P, bool) {
	if pm.root == nil {
		var zK K
		var zP P
		return zK, zP, false
	}

	root := pm.root
	children := root.child
	key, priority := root.key, root.priority

	delete(pm.indexes, key)
	root.child = nil

	pm.root = pm.combine(children)
	if pm.root != nil {
		pm.root.previous = nil
	}

	pm.deallocate(root)
	return key, priority, true
}

// Peek returns the highest‑priority key‑value pair without removing it.
// If empty returns zero-key, zero-value, false.
//
// Complexity: O(1).
func (pm *PairingPriorityMap[K, P]) Peek() (K, P, bool) {
	if pm.root == nil {
		var zK K
		var zP P
		return zK, zP, false
	}
	return pm.root.key, pm.root.priority, true
}

// Keys returns an iterator for all keys in the collection.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (pm *PairingPriorityMap[K, P]) Keys() iter.Seq[K] {
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
func (pm *PairingPriorityMap[K, P]) Values() iter.Seq[P] {
	return func(yield func(P) bool) {
		for _, entry := range pm.indexes {
			if !yield(entry.priority) {
				return
			}
		}
	}
}

// All returns an iterator for key-value pairs.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (m *PairingPriorityMap[K, P]) All() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for key, entry := range m.indexes {
			if !yield(key, entry.priority) {
				return
			}
		}
	}
}
