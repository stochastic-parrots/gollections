// Package prioritymap provides high-performance, generic priority map implementations.
//
// The package offers a sophisticated data structure that combines the O(1) average
// lookup of a Hash Map with the O(log n) ordering of a Priority Queue. This allows
// for efficient "decrease-key" operations, which are essential for algorithms
// like Dijkstra's or A*.
//
// # PriorityMap Interface
//
// All priority maps implement the [PriorityMap] interface, providing a unified API
// for key-value management with priority ordering:
//
//	type PriorityMap[K comparable, V any] interface {
//		Set(key K, value V)
//		Get(key K) (V, bool)
//		Remove(key K) bool
//		Pop() (K, V, bool)
//		Peek() (K, V, bool)
//		pkg.Map[K, V]
//	}
//
// # Why this package?
//
//   - Efficient Updates: Unlike a standard heap, you can update the priority
//     of an existing key in O(log n) time without searching the entire structure.
//
//   - Type Safety: Fully leverages Go Generics (1.18+) to ensure keys and
//     priorities are strictly typed, eliminating interface{} casting.
//
//   - Dual Nature: Implements both the [pkg.Map] and priority queue
//     behaviors, making it a versatile tool for scheduling and graph traversal.
//
//   - Go Idiomatic: Designed to feel like a native Go collection, integrating
//     seamlessly with the broader 'gollections' ecosystem.
//
// # Core Concepts
//
// The [PriorityMap] is particularly powerful when you need to track the "best"
// element while frequently changing the scores of other elements.
//
// Example of a basic workflow:
//
//	pm := prioritymap.NewMin[string, int]()
//	pm.Set("task1", 10)
//	pm.Set("task2", 5)
//
//	// Update task1 to a higher priority (lower value)
//	pm.Set("task1", 2)
//
//	key, val, _ := pm.Pop() // Returns "task1", 2
//
// When using [Set], if the key already exists, the implementation internally
// performs a "fix" or "re-heapify" operation to maintain the correct order
// based on the new value.
package prioritymap
