package prioritymap

import (
	"iter"
	"math/bits"

	"github.com/stochastic-parrots/gollections/constraint"
)

const radixBucketsCount = 65

// radixEntry holds a key-priority pair and the necessary pointers for the
// bucket list that contains it.
type radixEntry[K comparable, P constraint.Integer] struct {
	key      K
	priority P
	bucket   int
	next     *radixEntry[K, P]
	previous *radixEntry[K, P]
}

// RadixPriorityMap implements a priority-map on top of a radix heap.
// It maintains a map from keys to bucket entries, allowing O(1) lookups,
// updates, and removals by key under the radix heap monotonicity contract.
//
// K must be comparable. P is the priority type and must be an integer type.
// Priorities are ordered as a min-heap.
//
// The initial monotone lower bound is zero. Priorities must never be lower than
// the current lower bound, which is initially zero and is later advanced by Pop.
// Supplying a lower priority violates the radix heap contract and causes the
// heap ordering guarantees to be lost.
//
// Internally the structure keeps three representations:
//   - entries - a map[K]*radixEntry mapping each key to its bucket entry;
//   - buckets - linked lists grouped by the highest differing bit from last;
//   - free    - a freelist used to reuse pre-allocated entries.
//
// This makes radix priority maps particularly effective for monotone workloads
// such as Dijkstra with non-negative integer edge weights.
type RadixPriorityMap[K comparable, P constraint.Integer] struct {
	entries map[K]*radixEntry[K, P]
	buckets [radixBucketsCount]*radixEntry[K, P]
	free    *radixEntry[K, P]
	last    P
}

// NewRadixPriorityMap creates an empty radix priority map with pre-allocated
// entries and map capacity.
//
// Complexity: O(capacity).
func NewRadixPriorityMap[K comparable, P constraint.Integer](capacity int) *RadixPriorityMap[K, P] {
	pm := &RadixPriorityMap[K, P]{
		entries: make(map[K]*radixEntry[K, P], capacity),
	}
	storage := make([]radixEntry[K, P], capacity)
	for i := range storage {
		storage[i].next = pm.free
		pm.free = &storage[i]
	}
	return pm
}

// radixBucket returns the bucket index for priority relative to the current
// monotone lower bound.
//
// Complexity: O(1).
func radixBucket[P constraint.Integer](priority, last P) int {
	return bits.Len64(uint64(priority) ^ uint64(last))
}

// allocate retrieves an entry from the freelist or creates a new one.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) allocate(key K, priority P) *radixEntry[K, P] {
	if pm.free == nil {
		return &radixEntry[K, P]{key: key, priority: priority}
	}

	entry := pm.free
	pm.free = entry.next
	entry.key = key
	entry.priority = priority
	entry.next = nil
	return entry
}

// release clears an entry's data and returns it to the freelist.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) release(entry *radixEntry[K, P]) {
	var zeroK K
	var zeroP P
	entry.key = zeroK
	entry.priority = zeroP
	entry.bucket = 0
	entry.previous = nil
	entry.next = pm.free
	pm.free = entry
}

// attach inserts an entry into the bucket selected by its priority and the
// current monotone lower bound.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) attach(entry *radixEntry[K, P]) {
	bucket := radixBucket(entry.priority, pm.last)
	entry.bucket = bucket
	entry.previous = nil
	entry.next = pm.buckets[bucket]
	if entry.next != nil {
		entry.next.previous = entry
	}
	pm.buckets[bucket] = entry
}

// detach removes an entry from its current bucket list.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) detach(entry *radixEntry[K, P]) {
	if entry.previous == nil {
		pm.buckets[entry.bucket] = entry.next
	} else {
		entry.previous.next = entry.next
	}
	if entry.next != nil {
		entry.next.previous = entry.previous
	}
	entry.next = nil
	entry.previous = nil
}

// refillFirstBucket advances the monotone lower bound to the smallest priority
// in the first non-empty bucket and redistributes that bucket's entries.
//
// Complexity: O(radixBucketsCount + m), where m is the number of entries in
// the selected bucket.
func (pm *RadixPriorityMap[K, P]) refillFirstBucket() {
	if pm.buckets[0] != nil || len(pm.entries) == 0 {
		return
	}

	bucketIndex := 1
	for pm.buckets[bucketIndex] == nil {
		bucketIndex++
	}

	bucket := pm.buckets[bucketIndex]
	minimum := bucket.priority
	for entry := bucket.next; entry != nil; entry = entry.next {
		if entry.priority < minimum {
			minimum = entry.priority
		}
	}
	pm.last = minimum

	pm.buckets[bucketIndex] = nil
	for entry := bucket; entry != nil; {
		next := entry.next
		entry.next = nil
		entry.previous = nil
		pm.attach(entry)
		entry = next
	}
}

// Set inserts a key or replaces its priority.
//
// Priority must be at least the current monotone lower bound. The lower bound
// is initially zero and is advanced by Pop. Supplying a lower priority violates
// the radix heap contract and causes the heap ordering guarantees to be lost.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Set(key K, priority P) {
	if entry, ok := pm.entries[key]; ok {
		pm.detach(entry)
		entry.priority = priority
		pm.attach(entry)
		return
	}

	entry := pm.allocate(key, priority)
	pm.entries[key] = entry
	pm.attach(entry)
}

// Update changes the priority of an existing key.
//
// It returns true if the key was found and updated. If the key does not exist,
// it performs no operation and returns false.
//
// Priority must satisfy the same monotonicity requirement as Set.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Update(key K, priority P) bool {
	entry, ok := pm.entries[key]
	if !ok {
		return false
	}

	pm.detach(entry)
	entry.priority = priority
	pm.attach(entry)
	return true
}

// Improve ensures that the key has at most the given priority.
//
// If the key does not exist, it is inserted with the provided priority.
// If the key already exists, its priority is updated only if the new priority
// is smaller than the current one.
//
// It returns true if the map was modified (either by insertion or update).
//
// Priority must still be at least the current monotone lower bound.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Improve(key K, priority P) bool {
	entry, ok := pm.entries[key]
	if !ok {
		pm.Set(key, priority)
		return true
	}
	if priority >= entry.priority {
		return false
	}

	pm.detach(entry)
	entry.priority = priority
	pm.attach(entry)
	return true
}

// Remove deletes the entry for key if present, returning true if an entry
// was removed.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Remove(key K) bool {
	entry, ok := pm.entries[key]
	if !ok {
		return false
	}

	pm.detach(entry)
	delete(pm.entries, key)
	pm.release(entry)
	return true
}

// Pop removes and returns the key-priority pair with the minimum priority.
// If empty, it returns zero values and false.
//
// Complexity: O(radixBucketsCount + m), where m is the number of entries
// redistributed from the selected bucket. Amortized over a full drain, each
// entry is redistributed at most once per bucket level.
func (pm *RadixPriorityMap[K, P]) Pop() (key K, priority P, ok bool) {
	if len(pm.entries) == 0 {
		return key, priority, false
	}

	pm.refillFirstBucket()
	entry := pm.buckets[0]
	key, priority = entry.key, entry.priority
	pm.detach(entry)
	delete(pm.entries, key)
	pm.release(entry)
	return key, priority, true
}

// Peek returns the minimum-priority key-value pair without removing it.
// If empty returns zero-key, zero-value, false.
//
// Peek does not advance the monotone lower bound.
//
// Complexity: O(radixBucketsCount + m), where m is the number of entries in
// the first non-empty bucket.
func (pm *RadixPriorityMap[K, P]) Peek() (key K, priority P, ok bool) {
	if len(pm.entries) == 0 {
		return key, priority, false
	}

	if entry := pm.buckets[0]; entry != nil {
		return entry.key, entry.priority, true
	}

	bucketIndex := 1
	for pm.buckets[bucketIndex] == nil {
		bucketIndex++
	}

	minimum := pm.buckets[bucketIndex]
	for entry := minimum.next; entry != nil; entry = entry.next {
		if entry.priority < minimum.priority {
			minimum = entry
		}
	}

	return minimum.key, minimum.priority, true
}

// Get returns the priority associated with key.
//
// The zero value of P and false are returned if the key does not exist.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Get(key K) (priority P, ok bool) {
	entry, ok := pm.entries[key]
	if !ok {
		return priority, false
	}
	return entry.priority, true
}

// Contains returns true if the key exists in the map.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Contains(key K) bool {
	_, ok := pm.entries[key]
	return ok
}

// Keys returns an iterator for all keys in the collection.
//
// Complexity: O(n) for a full traversal, O(1) per step.
// Note: This does not guarantee priority order; use Drain for priority-ordered traversal.
func (pm *RadixPriorityMap[K, P]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range pm.entries {
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
func (pm *RadixPriorityMap[K, P]) Values() iter.Seq[P] {
	return func(yield func(P) bool) {
		for _, entry := range pm.entries {
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
func (pm *RadixPriorityMap[K, P]) All() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for key, entry := range pm.entries {
			if !yield(key, entry.priority) {
				return
			}
		}
	}
}

// Drain returns a destructive iterator that removes and yields elements
// in priority order (minimum priority first).
//
// Since this is a destructive operation, the map will be empty after a
// full traversal. If the iteration is stopped early (e.g., via break),
// the map will retain only the remaining elements.
//
// Complexity: O(n * radixBucketsCount) for a full traversal.
func (pm *RadixPriorityMap[K, P]) Drain() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for {
			key, priority, ok := pm.Pop()
			if !ok || !yield(key, priority) {
				return
			}
		}
	}
}

// Clear removes all elements from the priority map.
//
// After calling Clear, the map will be empty, its length will be zero, and the
// monotone lower bound will be reset to zero. This operation is typically more
// efficient than creating a new map as it may reuse the pre-allocated entries.
//
// Complexity: O(n) to zero out elements (avoiding memory leaks).
func (pm *RadixPriorityMap[K, P]) Clear() {
	for i := range pm.buckets {
		pm.buckets[i] = nil
	}
	for _, entry := range pm.entries {
		pm.release(entry)
	}
	clear(pm.entries)
	pm.last = 0
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) IsEmpty() bool {
	return len(pm.entries) == 0
}

// Length returns the current number of elements in the heap.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) Length() int {
	return len(pm.entries)
}

// LastPriority returns the monotone lower bound established by Pop.
//
// Complexity: O(1).
func (pm *RadixPriorityMap[K, P]) LastPriority() P {
	return pm.last
}
