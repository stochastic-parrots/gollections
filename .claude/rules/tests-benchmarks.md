---
paths:
  - "**/*_test.go"
  - "**/*_bench_test.go"
  - "internal/benchmarks/**"
---

# Tests And Benchmarks

- Read the `Tests` and `Benchmarks` sections of
  `.agents/project-standards.md`.
- Use `github.com/stretchr/testify/assert` by default; use `require` only when
  later checks depend on the current assertion.
- Internal implementation tests live in the same package and may inspect
  invariants, backing storage, freelists, and cleanup.
- Public package tests use external `_test` packages when validating public
  contracts.
- Iterator tests cover full iteration, partial iteration, and empty structures
  when those paths exist.
- Benchmarks use `b.ReportAllocs()` and keep setup outside measured sections
  unless setup is the measured operation.
- Use local `*_bench_test.go` files beside internal implementations for
  operation-level measurements.
- Use `internal/benchmarks` only for optional cross-implementation suites:
  contracts and optional pure Go baselines in `datastructs`, input helpers in
  `models`, algorithms that exercise those contracts in `algorithms`, and
  benchmark entrypoints in `suites`.
- Do not add behavior to benchmark adapters just to include an implementation
  in a suite.
