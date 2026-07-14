// Package heap provides mutable generic priority queues backed by binary heaps.
//
// The value at the top is selected by a priority comparator. Pop, Replace, and
// Push maintain the heap property; Peek observes the current top value without
// removing it. Drain is destructive and yields values in priority order.
//
// # Ordering
//
// Use [Binary] when priority is defined by a custom function. The function must
// return true when its first argument has higher priority than its second.
//
// Use [OrderedBinary] for values satisfying cmp.Ordered. Pass [Min] to give
// smaller values priority or [Max] to give larger values priority. Both
// selectors return a reusable [BinaryFactory] whose methods return
// *[BinaryHeap].
// Binary heaps require a priority comparator, so the zero values of
// [BinaryFactory] and [BinaryHeap] are invalid.
//
// # Construction And Ownership
//
// [BinaryFactory.New] creates an empty heap with the requested capacity.
// [BinaryFactory.From] heapifies a slice in place and transfers ownership of its
// backing storage; the caller must not use the slice or aliases of its backing
// array afterward. [BinaryFactory.Clone] makes a shallow copy before heapifying
// and does not retain the source slice.
//
// # Complexity
//
// Peek is O(1). Pop and Replace are O(log N). From and Clone build a heap in
// O(N). Push is O(K log N) for K inserted values, with a bulk heapify path when
// rebuilding is cheaper. Drain is O(N log N) when fully consumed.
//
// All and Enumerate expose the internal heap representation and do not guarantee
// priority order. Use Drain when destructive priority-ordered traversal is required.
package heap
