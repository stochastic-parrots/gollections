# gollections

Generic collection data structures for Go, built around type safety, predictable
APIs, and idiomatic iteration with `iter.Seq`.

`gollections` provides focused implementations for common data-structure needs
that are either missing from the standard library or awkward to use with
generics, such as lists, deques, heaps, and priority maps.

## Install

```bash
go get github.com/stochastic-parrots/gollections
```

The module uses Go's standard iterator APIs and targets Go 1.24+.

## Collections

| Package | Structures | Use when you need |
| --- | --- | --- |
| `list` | `ArrayList`, `LinkedList` | Indexed, ordered sequences with forward/backward traversal |
| `sortedlist` | `ArraySortedList`, `OrderedArraySortedList` | Sorted sequences that allow duplicate values |
| `deque` | `ArrayDeque`, `LinkedDeque` | Fast insertion and removal at both ends |
| `heap` | `BinaryHeap` | Priority queue behavior with min, max, or custom ordering |
| `prioritymap` | `BinaryHeapPriorityMap`, `PairingHeapPriorityMap`, `RadixHeapPriorityMap` | Keyed priority queues, including monotone integer workloads |

Package-level examples live alongside each public package and are rendered by
Go documentation tools.

## Design

- Generic APIs: no `interface{}` casting for stored values.
- Read-only root interfaces: `gollections.Collection` and `gollections.Map`
  are observation contracts. They expose inspection and iteration, while
  mutating operations live in structure-specific interfaces such as
  `list.List`, `sortedlist.SortedList`, `deque.Deque`, `heap.Heap`, and
  `prioritymap.PriorityMap`.
- Go iterators: collections expose `All` and `Enumerate` for `range` loops.
- Reusable clearing: `Clear` removes stored references, preserves construction
  configuration, and leaves the structure ready for reuse. It does not guarantee
  storage shrinkage or return memory to the Go runtime. Slice-backed structures
  retain capacity where practical, freelists retain at most their construction
  limit, and linked structures without a freelist may release their nodes.
- Internal implementations: public packages expose stable concrete factories while
  concrete internals live under `internal`.
- Read-only views: packages such as `list`, `sortedlist`, `deque`, and
  `prioritymap` expose wrappers for sharing non-mutating access without
  allowing type assertion back to the mutable interface. These wrappers are
  capability restrictions, not snapshots or concurrency synchronization.
- Concurrency: collections are not safe for concurrent use unless explicitly
  documented otherwise. Callers must synchronize shared access when any
  goroutine may mutate the collection.
- JSON support: lists and deques marshal and unmarshal as arrays in logical
  order. Sorted lists and heaps marshal as arrays but intentionally do not
  implement `json.Unmarshaler`; decode into a slice, then construct them through
  a factory so ordering and comparator settings remain explicit.

## Construction

Public packages select an implementation through a reusable typed factory.
Factory methods return concrete collection types without interface dispatch:

```go
items := list.Array[string]().New(16)
queue := deque.Linked[int]().From([]int{1, 2, 3})
scores := sortedlist.OrderedArray[int]().Clone([]int{3, 1, 2})
pending := prioritymap.OrderedBinaryHeap[string, int](prioritymap.Min).New(32)
```

Sorted lists and heaps intentionally do not implement `json.Unmarshaler`.
Decode into a slice and use the selected factory to establish their invariants:

```go
var values []int
if err := json.Unmarshal(data, &values); err != nil {
	return err
}
scores := sortedlist.Array(cmp.Compare[int]).From(values)
```

Array-backed `From` methods transfer ownership of the provided slice and may
reorder it. The caller must not use the slice or aliases of its backing array
afterward. Use `Clone` when the source must remain available. Linked structures
always copy values into nodes, so their `From` methods do not retain the source.

## Choosing a list

- `ArrayList`: use as the default indexed sequence when O(1) random access,
  O(1) indexed writes, append-heavy growth, and memory locality matter.
- `LinkedList`: use when boundary insertions/removals and O(1) reversal matter
  more than random indexed access.

Lists preserve explicit positional order. Use `sortedlist` when the collection
should stay ordered by value instead of by insertion position.

## Choosing a sorted list

- `OrderedArraySortedList`: use by default when element values satisfy
  `cmp.Ordered`; it uses standard-library ordered sort and binary search paths
  without custom comparator calls.
- `ArraySortedList`: use when values need a custom comparator, such as sorting
  structs by one field or using descending order. This flexibility means each
  search and sort comparison calls the comparator.

Both sorted-list implementations are best for data that is built once or
updated occasionally and queried many times. Lookups and bounds are O(log N),
indexed access is O(1), and single-element insertion/removal shifts O(N) values.
For custom comparators, values that compare equal may appear in any relative
order; add a tie-breaker to the comparator when that order matters.

## Choosing a deque

- `ArrayDeque`: use as the default double-ended queue when amortized O(1)
  operations at both ends and cache-friendly traversal are useful.
- `LinkedDeque`: use when growth is unpredictable or when avoiding backing
  array reallocations matters more than memory locality.

Deques are the right fit for queue, stack, sliding-window, worklist, and small
scheduling-buffer patterns.

## Choosing a priority map

- `BinaryHeapPriorityMap`: predictable O(log N) updates with compact,
  contiguous storage.
- `PairingHeapPriorityMap`: a strong general-purpose choice for workloads with
  frequent priority improvements.
- `RadixHeapPriorityMap`: optimized for unsigned integer priorities where
  popped priorities never decrease, such as Dijkstra with non-negative integer
  edge weights. After `Pop` returns `p`, callers must never insert or update an
  entry to a priority below `p`; this precondition is intentionally unchecked.

## Development

Run the full test suite:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

The project aims for complete coverage on internal data-structure
implementations. To inspect coverage for a package:

```bash
go test -cover ./internal/list
go test -cover ./internal/sortedlist
go test -cover ./internal/heap
go test -cover ./internal/prioritymap
go test -cover ./internal/deque
```

## Status

The available packages are `list`, `sortedlist`, `deque`, `heap`, and
`prioritymap`. `queue`, `stack`, and set families are planned as focused APIs on
top of the same collection foundations.
