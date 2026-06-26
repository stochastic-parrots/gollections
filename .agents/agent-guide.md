# Agent Guide

Use this file to orient quickly, then rely on
[project-standards.md](project-standards.md) for the full rules.

## Repository Snapshot

- `gollections` is a Go 1.24 generic collections library.
- Root package interfaces are read-only contracts: `Collection[T]` and
  `Map[K, V]`.
- Public mutable APIs live in focused packages such as `list`, `deque`, `heap`,
  and `prioritymap`.
- Concrete implementations live under `internal/...`; public packages expose
  interfaces, type aliases, constructors, examples, and docs.
- The codebase favors direct data-structure code over broad helper layers.

## Before Editing

1. Read the relevant sections of [project-standards.md](project-standards.md).
2. Inspect the closest existing package or implementation pattern.
3. Check `git status --short` and preserve unrelated user changes.
4. Decide whether the change affects public API, docs, tests, benchmarks, or
   README guidance.

## Task Routing

- Public API or new package: read `Project Shape`, `Package Layout`,
  `API Design`, `Naming`, and `Documentation`.
- Internal implementation: read `Implementation Style`, `Code Style`, and the
  related public API section.
- Tests: read `Tests`, then mirror the implementation method order where
  practical.
- Benchmarks: read `Benchmarks`; use local `*_bench_test.go` files beside
  internal implementations for operation-level measurements. Use
  `internal/benchmarks` only for optional cross-implementation suites:
  structure contracts, implementation lists, adapters, and optional pure Go
  `stdlib_*.go` baselines go in `datastructs`; benchmark input helpers go in
  `models`; algorithms that exercise `datastructs` contracts go in
  `algorithms`; benchmark entrypoints and timer control go in `suites`. Model
  suites after the existing `heap` and `prioritymap` suites. Adapters normalize
  APIs only; they must not invent behavior.
- Docs or examples: read `Documentation`; exported identifiers need comments
  that start with the identifier. AI-written comments should explain contracts,
  tradeoffs, invariants, mutation, errors, or complexity, and should not merely
  restate obvious code.
- Commits or tags: read `Commits And Tags`.

## Default Engineering Rules

- Preserve existing public API names unless the user asks for a breaking
  change.
- Keep empty reads and removals safe: return zero values plus `false` or a
  package-owned error.
- Keep root interfaces observation-only; mutations belong in structure-specific
  packages.
- Clear removed slots, nodes, and freelist entries so references do not leak.
- Use `gofmt` on changed Go files.
- Prefer focused tests with `assert`; use `require` only when later assertions
  depend on the current one.
- For benchmark suites, use a `get<Family>Suite(...)` helper and check that the
  benchmark abstraction itself does not introduce hot-loop allocations.

## Verification

- Run targeted package tests while iterating.
- Run `go test ./...` before finishing changes that affect shared behavior.
- For coverage-sensitive internal changes, run the affected package with
  `-coverprofile` and inspect `go tool cover -func`.
- Use `make lint` when changing public APIs, docs, or broad implementation
  behavior.

## When To Ask

Ask the user before making an irreversible public API change, changing
documented semantics, removing behavior, or choosing between two incompatible
design directions. Otherwise make the smallest conservative change that matches
the existing codebase.
