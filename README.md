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
  `list.List`, `deque.Deque`, `heap.Heap`, and `prioritymap.PriorityMap`.
- Go iterators: collections expose `All` and `Enumerate` for `range` loops.
- Internal implementations: public packages expose stable constructors while
  concrete internals live under `internal`.
- Read-only views: packages such as `list`, `deque`, and `prioritymap` expose
  wrappers for sharing non-mutating access without allowing type assertion back
  to the mutable interface.
- JSON support: linear collections and heaps can marshal/unmarshal as arrays
  where the operation makes sense.

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
go test -cover ./internal/heap
go test -cover ./internal/prioritymap
go test -cover ./internal/deque
```

## Status

The available packages are `list`, `deque`, `heap`, and `prioritymap`. `queue`
and `stack` are planned as thin, focused APIs on top of the same collection
foundations.
