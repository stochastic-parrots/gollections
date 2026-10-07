// Package gollections provides a suite of high-performance, generic data structures for Go.
//
// The library leverages Go generics for type safety and standard iterators for
// idiomatic data traversal. Instead of a "one-size-fits-all" approach,
// gollections provides specialized implementations optimized for specific memory
// and performance profiles.
//
// # Core Interfaces
//
// Most data structures in this module implement one of the two base interfaces,
// ensuring a predictable API across different implementations:
//
//   - [Collection]: Foundation for value collections like lists, heaps, and sets.
//   - [Map]: Base operations for key-value based structures.
//
// Root-level interfaces are intentionally read-only capability contracts. They
// describe how callers can inspect or iterate over a structure without changing
// its state. Mutating operations such as Add, Append, Push, Set, Remove, Pop,
// or Clear are exposed only by the structure-specific interfaces in each
// subpackage, where their semantics are precise and unambiguous.
//
// # Clearing And Reuse
//
// Clear removes all elements, releases stored references, preserves the
// structure's construction configuration, and leaves it ready for reuse.
// Resource retention is implementation-specific: slice-backed structures may
// preserve capacity, bounded freelists retain at most their configured limit,
// and linked structures without a freelist may release their nodes. Clear does
// not guarantee that storage is shrunk or memory is returned to the Go runtime.
// Fixed-range disjoint sets use Reset to restore singleton sets while retaining
// all values in their range.
//
// # JSON
//
// Concrete lists, deques, sorted lists, heaps, and sets marshal without external
// construction input; each package documents its JSON array's traversal order.
// Concrete lists and deques also implement json.Unmarshaler because the array
// completely defines their logical element order. Sorted lists, heaps, and sets
// require explicit ordering, priority, or identity configuration; decode into a
// slice and use the matching constructor family. Disjoint sets have no JSON
// representation because an array of range values does not encode the partition.
//
// Structural interfaces do not require formatting or JSON capabilities. Use a
// capability assertion for direct calls to String, MarshalJSON, or UnmarshalJSON
// through an interface. Standard fmt and encoding/json calls discover the
// capabilities of the underlying concrete value. Readonly wrappers retain their
// formatting and JSON methods through fmt.Sprint and json.Marshal on the wrapped
// collection, following those packages' formatting and encoding behavior.
//
// # Concurrency
//
// Unless a type explicitly documents otherwise, collections in this module are
// not safe for concurrent use. Callers must synchronize access when at least one
// goroutine may mutate a shared collection. Readonly interfaces and views limit
// the operations available through an API; they do not provide synchronization
// or a snapshot of the underlying collection. Disjoint-set queries such as Find
// and Connected may compress paths and require synchronization on shared access.
//
// # Subpackages
//
// The library is organized into specialized subpackages. Refer to each package
// documentation for detailed time complexity tables:
//
//   - [github.com/stochastic-parrots/gollections/list]:
//     Indexed sequences like ArrayList and LinkedList.
//
//   - [github.com/stochastic-parrots/gollections/sortedlist]:
//     Sorted sequences like ArraySortedList.
//
//   - [github.com/stochastic-parrots/gollections/deque]:
//     Double-ended queues backed by circular arrays or linked nodes.
//
//   - [github.com/stochastic-parrots/gollections/heap]:
//     Priority-based ordering (Min/Max Heaps).
//
//   - [github.com/stochastic-parrots/gollections/prioritymap]:
//     A hybrid structure combining Map lookups with Heap ordering.
//
//   - [github.com/stochastic-parrots/gollections/set]:
//     Unique comparable values or arbitrary values identified by derived keys.
//
//   - [github.com/stochastic-parrots/gollections/disjointset]:
//     Union-find over fixed integer ranges, using union by rank or size.
//
// # Design Principles
//
//   - Type Safety: Full generic support ensures compile-time type checking.
//
//   - Read-only Roots: Shared root interfaces expose observation and iteration;
//     mutation belongs to the concrete collection family.
//
//   - Performance: Focused on O(1) and O(log N) operations where possible.
//
//   - Idiomatic Go: Full support for range over iterators.
//
// # Documentation Characteristics
//
// Each data structure includes detailed performance documentation in its subpackage.
// Refer to the package constructors (for example, [list.NewArray]) for
// construction and ownership contracts.
package gollections
