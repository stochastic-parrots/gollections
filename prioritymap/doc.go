// Package prioritymap provides mutable maps whose entries are also ordered by
// priority.
//
// A priority map combines key-based lookup with priority-queue operations.
// Entries can be inserted, updated, improved, removed by key, or popped in
// priority order. The comparator associated with an implementation defines what
// "highest priority" means.
//
// # Implementations
//
// [BinaryHeapPriorityMap] provides predictable O(log N) insertion, update,
// removal, and Pop, with O(1) key lookup and Peek. Select it with [BinaryHeap]
// or [OrderedBinaryHeap].
//
// [PairingHeapPriorityMap] provides O(1) amortized priority improvements and is
// intended for workloads with frequent decrease-key or increase-key operations.
// Other priority mutations and removals are O(log N) amortized. Select it with
// [PairingHeap] or [OrderedPairingHeap].
//
// [RadixHeapPriorityMap] is a monotone min-priority map for non-negative integer
// priorities. It provides O(1) Set, Update, Improve, Get, and Remove, with O(W)
// amortized Pop for a W-bit priority. Select it with [RadixHeap].
//
// # Ordering
//
// Binary and pairing heap factories accept either a custom priority comparator
// or the natural order of a cmp.Ordered priority type. Pass [Min] to an ordered
// factory when smaller values have higher priority, or [Max] when larger values
// have higher priority.
//
// The radix heap is always min-priority and monotone. Its lower bound starts at
// zero and advances when Pop returns a value. Subsequent insertions and updates
// must not use priorities smaller than that bound. This precondition is not
// checked at runtime; violating it invalidates ordering guarantees. Peek does
// not advance the lower bound.
//
// # Views And Iteration
//
// [PriorityMap] exposes mutation, while [Readonly] provides key lookup, Peek,
// and map iteration. [AsReadonly] prevents callers from recovering the mutable
// map through a type assertion. Drain is destructive and removes entries as
// they are yielded; stopping iteration early leaves unvisited entries in the
// map.
//
// Keys, Values, and All do not guarantee priority order. Separate Keys and
// Values iterations are not positionally related. Use All to retain key-priority
// association and Drain for destructive priority-ordered traversal.
//
// Priority maps require initialized indexes and, except for radix heaps, a
// priority comparator. Their concrete zero values are invalid and must be
// constructed through the corresponding factory selector.
package prioritymap
