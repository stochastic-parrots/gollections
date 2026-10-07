// Package heap provides mutable generic priority queues backed by binary heaps.
//
// The value at the top is selected by a priority comparator. Pop, Replace,
// and Push maintain the heap property; Peek observes the current top
// value without removing it. Drain is destructive and yields values in
// priority order.
//
// # Ordering
//
// Use [NewBinary], [BinaryFrom], [BinaryClone], or [BinaryFromSeq] when
// priority is defined by a custom predicate. It returns true when its first
// argument has higher priority and must define a stable strict weak ordering.
//
// Use the corresponding OrderedBinary constructors for cmp.Ordered values.
// Pass [Min] to give smaller values priority or [Max] for larger values.
// [BinaryHeap] requires a comparator, so its zero value is invalid.
//
// # Construction And Ownership
//
// [NewBinary] creates an empty heap with the requested capacity. [BinaryFrom]
// heapifies a slice in place and transfers ownership of its backing array;
// the caller must stop using the slice and all aliases afterward. [BinaryClone]
// makes an independent shallow copy before heapifying. [BinaryFromSeq] consumes
// an iterator once into new storage and heapifies in O(N) time. OrderedBinary
// variants follow the same ownership and complexity contracts.
//
// # JSON
//
// BinaryHeap marshals as an array in internal heap order, not Pop or Drain
// order. The Heap interface does not require JSON or formatting methods.
// BinaryHeap does not implement json.Unmarshaler because JSON cannot preserve
// the comparator. Decode into []T and pass the values to the same constructor
// family used to define heap priority.
//
// # Complexity
//
// Peek is O(1). Pop and Replace are O(log N). From, Clone, and FromSeq build a
// heap in O(N). Push of one value is O(log N). Push of multiple values is
// O(N + len(xs)) when the heap is empty or len(xs) exceeds both N and 64;
// otherwise it is O(len(xs) log (N + len(xs))). Drain is O(N log N) when
// fully consumed.
//
// All and Enumerate expose the internal heap representation and do not guarantee
// priority order. Use Drain when destructive priority-ordered traversal is required.
package heap
