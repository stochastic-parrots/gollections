# gollections

[![CI](https://github.com/stochastic-parrots/gollections/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/stochastic-parrots/gollections/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/stochastic-parrots/gollections/branch/main/graph/badge.svg)](https://app.codecov.io/gh/stochastic-parrots/gollections)

Generic collection data structures for Go, built around type safety, predictable
APIs, and idiomatic iteration with `iter.Seq`.

`gollections` provides focused implementations for common data-structure needs
that are either missing from the standard library or awkward to use with
generics, such as lists, deques, sets, heaps, and priority maps.

## Install

```bash
go get github.com/stochastic-parrots/gollections
```

The module uses Go's standard iterator APIs and targets Go 1.24+.

## Collections

| Package | Structures | Use when you need |
| --- | --- | --- |
| `list` | `ArrayList`, `LinkedList` | Indexed, ordered sequences with forward/backward traversal |
| `sortedlist` | `ArraySortedList` | Sorted sequences that allow duplicate values |
| `deque` | `ArrayDeque`, `LinkedDeque` | Fast insertion and removal at both ends |
| `set` | `HashSet`, `KeyedHashSet` | Unique values using Go equality or a derived identity key |
| `heap` | `BinaryHeap` | Priority queue behavior with min, max, or custom ordering |
| `prioritymap` | `BinaryHeapPriorityMap`, `PairingHeapPriorityMap`, `RadixHeapPriorityMap` | Keyed priority queues, including monotone integer workloads |
| `disjointset` | `IntsRangeByRank`, `IntsRangeBySize` | Track connectivity among all integers in a fixed inclusive range |

Package-level examples live alongside each public package and are rendered by
Go documentation tools.

## Design

- Generic APIs: no `interface{}` casting for stored values.
- Read-only root interfaces: `gollections.Collection` and `gollections.Map` are
  observation contracts. They expose inspection and iteration, while mutating
  operations live in structure-specific interfaces such as
  `list.List`, `sortedlist.SortedList`, `deque.Deque`, `set.Set`, `heap.Heap`,
  `prioritymap.PriorityMap`, and `disjointset.DisjointSet`.
- Go iterators: collections expose `All` and `Enumerate` for `range` loops.
- Reusable clearing: `Clear` removes stored references, preserves construction
  configuration, and leaves the structure ready for reuse. It does not guarantee
  storage shrinkage or return memory to the Go runtime. Slice-backed structures
  retain capacity where practical, freelists retain at most their construction
  limit, and linked structures without a freelist may release their nodes.
  Fixed-range disjoint sets use `Reset` to restore singleton sets instead.
- Internal implementations: public packages expose direct constructors and concrete aliases while
  concrete internals live under `internal`.
- Read-only views: packages such as `list`, `sortedlist`, `deque`, `set`,
  `prioritymap`, and `disjointset` expose wrappers for sharing observation without
  allowing type assertion back to the mutable interface. These wrappers are
  capability restrictions, not snapshots or concurrency synchronization.
- Concurrency: collections are not safe for concurrent use unless explicitly
  documented otherwise. Callers must synchronize shared access when any
  goroutine may mutate the collection.
- JSON support: concrete lists, deques, sorted lists, heaps, and sets marshal as
  arrays without external construction input, using each package's traversal
  order. Concrete lists and deques also unmarshal because array order completely
  defines their logical state. Structural interfaces do not require JSON or
  formatting capabilities.
  Sorted lists, heaps, and sets require slice decoding followed by the matching
  constructor so ordering, priority, and set identity remain explicit.
  Disjoint sets have no JSON representation because an array of elements does
  not encode the partition.

## Formatting and JSON capabilities

Collection interfaces describe structural operations. Concrete types retain
their optional formatting and JSON methods, and standard `fmt` and
`encoding/json` calls discover those methods on the underlying concrete value.

This is a breaking change for direct optional method calls and assignments
through structural interfaces:

| Structural interface | Capability no longer required |
| --- | --- |
| `list.Readonly`, `deque.Readonly`, `sortedlist.Readonly`, `heap.Heap` | `fmt.Stringer`, `json.Marshaler` |
| `set.Readonly` | `json.Marshaler` |
| `list.List`, `deque.Deque` | `json.Unmarshaler`, plus inherited readonly capabilities |

Constructors return concrete types, so direct calls such as `items.String()`
remain available. For built-in collections stored behind a structural interface,
use standard formatting and JSON functions without type assertions:

```go
var items list.List[int] = list.ArrayFrom([]int{1, 2})
fmt.Println(items)
data, err := json.Marshal(items)
if err != nil {
    return err
}
fmt.Println(string(data))
```

Use a checked capability assertion when you need to call an optional method
directly through a structural interface or pass it to a function that requires
that capability.

Readonly wrappers keep their formatting and JSON methods by calling `fmt.Sprint`
and `json.Marshal` on the wrapped collection. Built-in collections retain their
five-element string display limit and JSON arrays, including `[]` when empty.
External implementations follow the standard packages' behavior; their elements
are not rendered automatically through collection iteration. `json.Marshal`
validates and compacts custom JSON, escapes HTML characters, and wraps errors
returned by `MarshalJSON` in `json.MarshalerError`. Views remain live and do not
expose mutation or `json.Unmarshaler`.

## Construction

Public packages expose direct constructors returning concrete collection pointers:

```go
items := list.NewArray[string](16)
queue := deque.LinkedFrom([]int{1, 2, 3})
scores := sortedlist.OrderedArrayClone(sortedlist.Asc, []int{3, 1, 2})
pending := prioritymap.NewOrderedBinaryHeap[string, int](prioritymap.Min, 32)
visited := set.NewHashSet[string](32)
groups := disjointset.NewIntsRangeBySize(0, 99)
```

Sorted lists, heaps, and sets do not implement `json.Unmarshaler`. Decode
into a slice and use a constructor that establishes the collection invariant:

```go
var values []int
if err := json.Unmarshal(data, &values); err != nil {
    return err
}
scores := sortedlist.ArrayFrom(cmp.Compare[int], values)
```

Array-backed `From` functions transfer ownership of the provided slice and may
reorder it. Stop using the slice and all aliases of its backing array after
construction. `Clone` functions make independent shallow copies of slice
storage. Linked `From` functions copy values into nodes. Set `From` functions
copy values into map storage. `FromSeq` functions consume an iterator once.

### Migration from factories

| Package | Previous factory calls | Direct constructors |
| --- | --- | --- |
| `list` | `Array[T]().New/Clone/From/FromSeq` | `NewArray[T]` / `ArrayClone` / `ArrayFrom` / `ArrayFromSeq` |
| `list` | `Linked[T]().New/From/FromSeq` | `NewLinked[T]` / `LinkedFrom` / `LinkedFromSeq` |
| `deque` | `Array[T]().New/Clone/From/FromSeq` | `NewArray[T]` / `ArrayClone` / `ArrayFrom` / `ArrayFromSeq` |
| `deque` | `Linked[T]().New/From/FromSeq` | `NewLinked[T]` / `LinkedFrom` / `LinkedFromSeq` |
| `sortedlist` | `Array(cmp).New/Clone/From/FromSeq` | `NewArray(cmp, cap)` / `ArrayClone(cmp, xs)` / `ArrayFrom(cmp, xs)` / `ArrayFromSeq(cmp, seq)` |
| `sortedlist` | `OrderedArray[T](order).New/Clone/From/FromSeq` | `NewOrderedArray[T](order, cap)` / `OrderedArrayClone(order, xs)` / `OrderedArrayFrom(order, xs)` / `OrderedArrayFromSeq(order, seq)` |
| `heap` | `Binary(hasPriority).New/Clone/From` | `NewBinary(hasPriority, cap)` / `BinaryClone(hasPriority, xs)` / `BinaryFrom(hasPriority, xs)` |
| `heap` | `OrderedBinary[T](order).New/Clone/From` | `NewOrderedBinary[T](order, cap)` / `OrderedBinaryClone(order, xs)` / `OrderedBinaryFrom(order, xs)` |
| `prioritymap` | `BinaryHeap[K, P](hasPriority).New` | `NewBinaryHeap[K, P](hasPriority, cap)` |
| `prioritymap` | `OrderedBinaryHeap[K, P](order).New` | `NewOrderedBinaryHeap[K, P](order, cap)` |
| `prioritymap` | `PairingHeap[K, P](hasPriority).New` | `NewPairingHeap[K, P](hasPriority, cap)` |
| `prioritymap` | `OrderedPairingHeap[K, P](order).New` | `NewOrderedPairingHeap[K, P](order, cap)` |
| `prioritymap` | `RadixHeap[K, P]().New` | `NewRadixHeap[K, P](cap)` |
| `set` | `HashSetOf[T]().New/From/FromSeq` | `NewHashSet[T]` / `HashSetFrom` / `HashSetFromSeq` |
| `set` | `HashSetBy(keyOf).New/From/FromSeq` | `NewKeyedHashSet(keyOf, cap)` / `KeyedHashSetFrom(keyOf, xs)` / `KeyedHashSetFromSeq(keyOf, seq)` |

`heap.BinaryFromSeq(hasPriority, seq)` and
`heap.OrderedBinaryFromSeq(order, seq)` are new O(N) constructors.

## Mutations

`Append`, `Prepend`, `Push`, and `Add` accept zero, one, or many values through
one variadic method. Calling them without values leaves the collection unchanged:

```go
items.Append()
items.Append("one")
items.Append("two", "three")
xs := []string{"four", "five"}
items.Append(xs...)
```

Batch calls use the same singular name as single-value calls. Custom interfaces
and method values must use the variadic signature. Set `Add` and `Remove` return
counts; use `> 0` when a caller needs a boolean.

Variadic calls through collection interfaces can allocate a temporary argument
slice, even for one value. Calling the concrete type directly avoids interface
dispatch overhead. Benchmark the call form used by your workload; reusing an
argument slice can avoid repeated temporary allocations.

## Choosing a list

- `ArrayList`: use as the default indexed sequence when O(1) random access,
  O(1) indexed writes, append-heavy growth, and memory locality matter.
- `LinkedList`: use when boundary insertions/removals and O(1) reversal matter
  more than random indexed access.

Lists preserve explicit positional order. Use `sortedlist` when the collection
should stay ordered by value instead of by insertion position.

## Choosing a sorted list

`ArraySortedList` is the single slice-backed implementation. Choose its
ordering through one of two constructor families:

- `NewOrderedArray`, `OrderedArrayFrom`, or `OrderedArrayClone` when element values
  satisfy `cmp.Ordered` and natural order is enough.
- `NewArray`, `ArrayFrom`, or `ArrayClone` for custom ordering, such as sorting structs by one or more
  fields.

Both constructor families build `ArraySortedList` values, so they
have the same performance characteristics. Sorted lists are best for data that
is built once or updated occasionally and queried many times. Lookups and bounds
are O(log N), indexed access is O(1), and single-element insertion/removal
shifts O(N) values. Values that compare equal may appear in any relative order;
add a tie-breaker to the comparator when that order matters.

## Choosing a deque

- `ArrayDeque`: use as the default double-ended queue when amortized O(1)
  operations at both ends and cache-friendly traversal are useful.
- `LinkedDeque`: use when growth is unpredictable or when avoiding backing
  array reallocations matters more than memory locality.

Deques are the right fit for queue, stack, sliding-window, worklist, and small
scheduling-buffer patterns.

## Choosing a set

- `HashSet`: use when values are comparable and Go equality defines membership.
- `KeyedHashSet`: use for arbitrary values, or when a derived comparable key such
  as an ID defines membership. The first value inserted for a key remains its
  representative, and the derived key must remain stable while the value is in
  the set.

Set identities must have reflexive equality. Floating-point NaN cannot be used
directly as a `HashSet` value or `KeyedHashSet` key because NaN is not equal to
itself. Use `NewKeyedHashSet` with a canonical comparable key when NaN membership is
required. When `HashSet` uses an interface type, its dynamic values must also be
comparable, matching native Go map requirements.

`Add(xs ...T)` and `Remove(xs ...T)` return the number of values whose membership
changed. Both accept zero, one, or many values; empty input returns zero. Use
`Add(x) > 0` or `Remove(x) > 0` when a boolean is needed. Iteration and JSON
array order are unspecified.

Both implementations provide pure set algebra (`Clone`, `Union`,
`Intersection`, `Difference`, and `SymmetricDifference`) as receiver methods.
The package-level `Equal`, `IsSubset`, `IsProperSubset`, `IsSuperset`,
`IsProperSuperset`, and `IsDisjoint` functions compare `Readonly` sets using Go
equality. Their `...By` counterparts take a shared identity function for both
operands. `Equal` compares each operand with the first, while `IsDisjoint`
requires every pair of operands to be disjoint. `KeyedHashSet` methods use the
receiver's identity function.
Destructive `UnionWith`, `IntersectWith`,
`DifferenceWith`, and `SymmetricDifferenceWith` operations reuse the receiver
and report how many memberships changed, which avoids allocating a second full
set when mutation is appropriate. Keep and pass the set pointers returned by
the constructors; initialized set structs must not be copied.

## Choosing a disjoint set

Use `disjointset` when every integer in an inclusive range is known up front
and you need to merge groups or test connectivity, such as components of a
graph whose vertices have contiguous integer IDs. `IntsRangeByRank` exposes
the rank heuristic; `IntsRangeBySize` reports how many integers belong to a
group. Both keep the range fixed, so `Reset` restores singleton groups rather
than removing elements.

With path compression and union by rank or size, `Find`, `Connected`, `Union`,
and the corresponding `Rank` or `Size` query take O(α(N)) amortized time, where
N is the number of integers and α is the inverse Ackermann function. A single
operation can take O(log N) time; `Reset` takes O(N).

## Choosing a priority map

- `BinaryHeapPriorityMap`: predictable O(log N) updates with compact,
  contiguous storage.
- `PairingHeapPriorityMap`: useful when priority improvements are frequent; each
  improvement does constant immediate heap-link work. Pop consolidates children
  later; use the conservative O(log N) amortized bound for mixed workloads.
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
go test -cover ./internal/set
go test -cover ./internal/disjointset
```

CI builds and tests every package with the minimum supported Go release and the
latest stable release. The published coverage report measures library packages;
the cross-implementation benchmark harness under `internal/benchmarks` remains
part of the build and test suite but is excluded from the coverage percentage.
Codecov requires at least 95% project coverage and 100% coverage for changed
lines.

## Status

The available packages are `list`, `sortedlist`, `deque`, `set`, `heap`, and
`prioritymap`, plus `disjointset` for fixed integer ranges. `queue` and
`stack` families are planned as focused APIs on top of the same collection
foundations.
