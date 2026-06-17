# Gollections Project Standards

These standards describe how to extend and maintain `gollections`. They are
intended for coding agents and human contributors.

For a compact task map, start with [agent-guide.md](agent-guide.md). When a
change is non-trivial, use [task-checklists.md](task-checklists.md) alongside
the relevant sections below.

## Project Shape

- The module targets Go 1.24 and uses standard iterator APIs from `iter`.
- The root package exposes shared read-only capability contracts:
  - `Collection[T]` for linear collections.
  - `Map[K, V]` for key-value structures.
- Mutating operations belong in structure-specific subpackages such as `list`,
  `deque`, `heap`, and `prioritymap`.
- Public subpackages expose stable constructors, type aliases, interfaces,
  examples, and package documentation.
- Concrete implementations live under `internal/...`.
- Prefer small, focused packages over broad utility packages.

## Package Layout

Follow the file layout already used by the library:

- Root package:
  - `collection.go` and `map.go` contain read-only base interfaces.
  - `doc.go` contains module-level package documentation.
- Public collection package:
  - `<package>.go` contains the public interface, readonly interface, and
    `AsReadonly` wrapper.
  - `factory.go` exposes type aliases and constructors.
  - `doc.go` explains the package.
  - `examples_test.go` contains public examples.
  - `<package>_test.go` validates public contracts from the external test
    package when possible.
- Internal implementation package:
  - `<strategy>_<structure>.go` contains the concrete implementation.
  - `<strategy>_<structure>_test.go` contains same-package tests.
  - `<strategy>_<structure>_bench_test.go` contains local benchmarks.
- Shared internal helpers live under `internal/shared/...` only when multiple
  families actually need them.

When adding a new data structure family, create public and internal packages in
parallel. Do not expose an internal concrete type directly; expose a public type
alias and constructor from the public package.

## API Design

- Keep root interfaces observation-only. Do not add mutating methods to
  `gollections.Collection` or `gollections.Map`.
- Expose mutations through the domain interface where their semantics are clear,
  for example `list.List`, `deque.Deque`, `heap.Heap`, and
  `prioritymap.PriorityMap`.
- Empty reads and removals should be safe. Methods such as `Pop`, `Peek`, `Get`,
  and indexed reads return zero values plus `false` or an explicit error rather
  than panicking.
- Iterators follow this convention:
  - `All` yields values or key-value pairs without mutation.
  - `Enumerate` yields index-value pairs for ordered collections.
  - `Keys`, `Values`, and `All` for maps do not guarantee priority or sorted
    order unless the method explicitly says so.
  - `Drain` is destructive and yields in the structure's priority/order
    contract.
- `Clear` must leave the structure empty, preserve reusable capacity where
  practical, and zero references that could otherwise keep values alive.
- For ordered behavior, use `cmp.Ordered` or the local comparator helpers in
  `internal/comparator`.
- For numeric-only generic APIs, use the package-level constraints in
  `constraint`.
- Public packages should provide readonly views when the data structure can be
  shared safely for observation. The pattern is:
  - `Readonly[...]` interface.
  - Mutable interface embedding `Readonly[...]`.
  - `AsReadonly` returning `nil` for nil input.
  - private `readonly` wrapper with one-line forwarding methods.
  - compile-time assertion: `var _ Readonly[...] = (*readonly[...])(nil)`.
- Public factory files should include compile-time assertions that internal
  implementations satisfy the public interface.

## Naming

- Public package names are short, singular, and domain-specific: `list`,
  `deque`, `heap`, `prioritymap`.
- Public concrete type aliases use the data-structure name:
  `ArrayList`, `LinkedList`, `ArrayDeque`, `LinkedDeque`, `BinaryHeap`,
  `BinaryHeapPriorityMap`, `PairingHeapPriorityMap`,
  `RadixHeapPriorityMap`.
- Internal implementation types include their backing strategy:
  `ArrayList`, `DoubleLinkedList`, `RingBufferDeque`, `BinaryHeap`,
  `BinaryPriorityMap`, `PairingPriorityMap`, `RadixPriorityMap`.
- Constructors use predictable prefixes:
  - `New...` creates an empty structure.
  - `...From` may consume and reorder the provided slice in place.
  - `...Clone` clones input before building the structure.
  - `Min...` and `Max...` use natural ordering for `cmp.Ordered` types.
- Boolean return variables should be named `ok` when matching Go lookup style.
- Prefer `key`, `value`, `priority`, `idx`, `entry`, `node`, and `current`
  over abbreviated or domain-opaque names.
- Use receiver names that match the structure family:
  - `list` implementations use `list` or a concise local name when clearer.
  - `deque` implementations use `deque`.
  - `heap` implementations use `heap`.
  - priority maps use `pm`.
  - readonly wrappers use `w`.
- Generic type parameter names are short and meaningful:
  - `T` for element values.
  - `K` for keys.
  - `V` for generic map values.
  - `P` for priorities.
- Use `xs ...T` for variadic element inputs, matching `Append(xs ...T)` and
  `Push(xs ...T)`.
- Use `idx` for indexes, not `index`, in method signatures and tests.
- Constructors should preserve the user-facing vocabulary:
  - `NewArray`, `NewLinked`, `NewBinary`, `NewMinBinary`, `NewMaxBinary`.
  - `BinaryFrom` means in-place construction from a slice.
  - `BinaryClone` means clone-before-build.
  - `NewBinaryHeap`, `NewPairingHeap`, and `NewRadixHeap` are priority map
    factories even though the internal type name includes `PriorityMap`.
- Test names use `Test<Type>_<Method>` or `TestNew<Type>`. Subtests use compact
  PascalCase names such as `FullIteration`, `PartialIteration`, `EmptyMap`,
  `WorsePriority`, or `Reuse`.

## Implementation Style

- Keep implementations direct and data-structure oriented. This repo favors
  explicit slice/node/map manipulation over clever helper abstractions.
- Internal structs should expose fields only inside `internal` packages. Tests
  may inspect those fields directly to verify invariants.
- Prefer pointer receivers for mutable data structures.
- Use named returns when they make zero-value failure paths obvious, especially
  for `Get`, `Peek`, and `Pop`:

  ```go
  func (pm *BinaryPriorityMap[K, P]) Peek() (key K, priority P, ok bool)
  ```

- Empty-state methods return zero values and `false`. They should not panic.
- Index-based methods return package-owned errors for invalid indexes.
- Mutating removal paths must clear discarded slots or nodes before releasing
  storage so references do not leak.
- Backing storage should be reused after `Clear` when the structure owns the
  storage. Preserve capacity unless there is a strong reason not to.
- Freelist-backed structures must reset all pointer fields before returning a
  node/entry to the freelist.
- Iterator functions should be written as closures returning early when
  `yield` returns false:

  ```go
  return func(yield func(T) bool) {
      for _, value := range values {
          if !yield(value) {
              return
          }
      }
  }
  ```

- Destructive iterators such as `Drain` should stop immediately on either empty
  structure or `yield == false`.
- Use simple `for range` loops over integer counts when targeting Go 1.24:
  `for i := range capacity`.
- Avoid reflection in collection implementations. Type switches are acceptable
  in constrained benchmark helpers where they are isolated from public hot
  paths.

## Documentation

- Public Go documentation is written in English.
- Every exported package, type, interface, and function should have a comment
  that starts with the exported identifier.
- Package docs live in `doc.go`.
- Public constructors should include a performance table when they define a
  data-structure choice. Use this style:

  ```go
  // Performance Summary (Time Complexity):
  //
  //	Operation           Time Complexity
  //	-----------------   ---------------
  //	Push(xs... T)       O(K log N)
  //	Pop()               O(log N)
  //	Peek()              O(1)
  ```

- Internal data-structure methods use shorter documentation with `Complexity:`
  blocks when the implementation is non-trivial.
- Document behavioral contracts close to the API that depends on them. For
  example, radix priority maps must document their monotonicity requirement near
  the type and mutation methods.
- Examples live in `examples_test.go`, use `Example...` names, and include
  `// Output:` blocks when deterministic.
- README changes should explain when to choose a structure, not only list that
  it exists.
- Interface method comments should describe return values and empty/error
  behavior, not implementation details.
- Constructor comments should explain when to choose that implementation.
- Use "highest priority" for comparator-based priority maps because the
  comparator defines what "highest" means. Use "minimum priority" only for
  min-only implementations such as radix heaps.
- When a method is intentionally destructive, say so in the first paragraph.
- Warnings should be explicit and near the constructor or method, for example
  `WARNING: This operation is In-Place and WILL modify the original slice
  order.`

## Tests

- Use `github.com/stretchr/testify/assert` by default. Use `require` only when a
  failed assertion would make the rest of the test invalid.
- Internal implementation tests live in the same package so they can inspect
  internal invariants, freelists, backing slices, and node links.
- Public package tests use the external `_test` package when validating public
  contracts.
- Prefer deterministic, focused tests over randomized tests unless randomness
  is essential to the behavior being validated.
- Keep tests organized by method:
  - constructor tests first,
  - public operations next,
  - iterators and destructive iterators,
  - integrity and cleanup tests near the end.
- Iterator tests should cover full iteration, partial iteration, and empty
  structures when the method supports those paths.
- Cleanup tests should verify both user-visible state and reference cleanup when
  the implementation stores pointers.
- The project aims for complete coverage on core internal data-structure
  packages: `internal/list`, `internal/deque`, `internal/heap`, and
  `internal/prioritymap`.
- Internal test files should mirror the implementation method order as much as
  practical. This makes coverage gaps easy to map back to code.
- Prefer direct assertions over helper-heavy test DSLs.
- Use helpers only when they remove repeated invariant checks or collection
  boilerplate. Mark helpers with `t.Helper()` when they accept `*testing.T`.
- Iterator tests should avoid assuming map iteration order. Use counts or
  `assert.ElementsMatch` for unordered output.
- Priority order tests should assert exact order only for deterministic
  priorities and deterministic tie behavior. Avoid relying on map iteration
  order for ties.
- For same-package internal tests, it is acceptable to inspect backing slices,
  capacities, freelists, node links, indexes, and zeroed discarded slots.
- Regression tests should name the behavior, not the bug number.

## Benchmarks

- Structure-specific benchmarks live next to internal implementations, for
  example `internal/prioritymap/*_bench_test.go`.
- Cross-implementation benchmark suites live in `internal/benchmarks/suites`.
- Benchmark wrapper types live under `internal/benchmarks/datastructs`.
- Benchmark models and algorithms should be generic only when that enables a
  meaningful comparison, such as comparing float and integer priority maps.
- Use `b.ReportAllocs()` for suite benchmarks.
- Reset or rebuild structures outside the measured section when setup is not
  part of the operation under test.

## Code Style

- Run `gofmt` on all changed Go files.
- Prefer simple generic code over reflection or `interface{}` casts.
- Keep allocations predictable. Reuse capacity when the implementation already
  owns backing storage.
- Avoid exposing `internal` implementation details through public interfaces.
- Add abstractions only when multiple structures actually need them.
- Keep error values centralized inside the package that owns the operation.
- Do not introduce panics for normal collection states such as empty reads.
- Preserve existing public API names unless the user explicitly asks for a
  breaking change.
- Imports are grouped as standard library first, then blank line, then project
  imports.
- Use package aliases only when they remove ambiguity, such as `pkg` for the
  root `gollections` package from a subpackage or `constructor` for an internal
  implementation package.
- Keep one-line forwarding methods in readonly wrappers on one line.
- Prefer small methods that each preserve one invariant over large methods with
  many modes.
- Do not add comments that simply repeat the code. Add comments for public API
  contracts, complexity, ownership, mutation, and non-obvious invariants.
- Keep public examples small and deterministic. They should demonstrate the
  intended API shape, not exhaustively test behavior.

## Verification

- For normal changes, run:

  ```bash
  go test ./...
  ```

- For coverage-sensitive internal changes, run the affected package with
  coverage, for example:

  ```bash
  go test ./internal/prioritymap -coverprofile=/tmp/internal_prioritymap.cover -covermode=count
  go tool cover -func=/tmp/internal_prioritymap.cover
  ```

- For race-sensitive changes, run:

  ```bash
  go test -race ./...
  ```

- Use `make lint` when changing public APIs, docs, or broad implementation
  behavior.

## Commits And Tags

- Commit messages follow Conventional Commits with concise scopes:
  - `feat(priority map): add radix heap priority map`
  - `test(priority map): cover radix heap priority map`
  - `doc(priority map): document radix heap priority map`
  - `perf(priority map): add radix heap benchmarks`
  - `fix(priority map): clear stale priority state`
- Keep commits reviewable and grouped by intent: feature, tests, docs,
  benchmarks, and cleanup should be separate when practical.
- Tags are annotated and currently follow `v0.0.N-alpha`.
