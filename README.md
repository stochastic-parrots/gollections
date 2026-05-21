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

The module uses Go's standard iterator APIs and targets Go 1.25+.

## Collections

| Package | Structures | Use when you need |
| --- | --- | --- |
| `list` | `ArrayList`, `LinkedList` | Indexed, ordered sequences with forward/backward traversal |
| `deque` | `ArrayDeque`, `LinkedDeque` | Fast insertion and removal at both ends |
| `heap` | `BinaryHeap` | Priority queue behavior with min, max, or custom ordering |
| `prioritymap` | `BinaryHeapPriorityMap`, `PairingHeapPriorityMap` | A priority queue with key lookup and priority updates |

## Examples

### Lists

```go
package main

import (
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
)

func main() {
	l := list.NewArray[int](4)
	l.Append(10, 20, 30)
	_ = l.Insert(1, 15)

	fmt.Println(slices.Collect(l.All()))
	fmt.Println(slices.Collect(l.Backward()))
}
```

### Heaps

```go
package main

import (
	"fmt"

	"github.com/stochastic-parrots/gollections/heap"
)

func main() {
	h := heap.NewMinBinary[int](0)
	h.Push(10, 3, 7, 1)

	for !h.IsEmpty() {
		value, _ := h.Pop()
		fmt.Println(value)
	}
}
```

### Priority Maps

```go
package main

import (
	"fmt"

	"github.com/stochastic-parrots/gollections/prioritymap"
)

func main() {
	pm := prioritymap.MinPairingHeap[string, int](10)
	pm.Set("slow", 30)
	pm.Set("fast", 5)
	pm.Improve("slow", 2)

	for key, priority := range pm.Drain() {
		fmt.Println(key, priority)
	}
}
```

### Deques

```go
package main

import (
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/deque"
)

func main() {
	d := deque.NewArray[int](2)
	d.Append(2, 3)
	d.Prepend(0, 1)

	front, _ := d.Shift()
	back, _ := d.Pop()

	fmt.Println(front, back)
	fmt.Println(slices.Collect(d.All()))
}
```

## Design

- Generic APIs: no `interface{}` casting for stored values.
- Go iterators: collections expose `All` and `Enumerate` for `range` loops.
- Internal implementations: public packages expose stable constructors while
  concrete internals live under `internal`.
- Read-only views: packages such as `list`, `deque`, and `prioritymap` expose
  wrappers for sharing non-mutating access.
- JSON support: linear collections and heaps can marshal/unmarshal as arrays
  where the operation makes sense.

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
