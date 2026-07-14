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
//
// # Construction And Ownership
//
// [BinaryFactory.New] creates an empty heap with the requested capacity.
// [BinaryFactory.From] heapifies a slice in place and transfers ownership of its
// backing storage. [BinaryFactory.Clone] makes a shallow copy before heapifying
// and does not retain the source slice.
//
// # JSON
//
// Heaps marshal as arrays in internal heap order, not Pop or Drain order. They
// do not implement json.Unmarshaler because JSON cannot preserve the comparator.
// Decode into []T and pass the values to the same factory used to define heap
// priority.
//
// # Complexity
//
// Peek is O(1). Pop and Replace are O(log N). From and Clone build a heap in
// O(N). Push is O(K log N) for K inserted values, with a bulk heapify path when
// rebuilding is cheaper. Drain is O(N log N) when fully consumed.
package heap
