# Task Checklists

Use these checklists for repeatable agent workflows. They supplement
[project-standards.md](project-standards.md); they do not replace it.

## Code Change

1. Read the relevant package docs, factory file, public interface, and internal
   implementation.
2. Find the closest existing test pattern before adding new tests.
3. Keep implementation, public constructors, examples, and README guidance in
   sync when behavior changes.
4. Run `gofmt` on changed Go files.
5. Run targeted tests, then broaden to `go test ./...` when needed.

## New Data Structure Family

1. Add public and internal packages in parallel.
2. Expose public interfaces, readonly views, type aliases, constructors,
   package docs, examples, and external package tests.
3. Keep concrete types under `internal/...`.
4. Add internal same-package tests for invariants and cleanup.
5. Add benchmarks beside implementations and suite wrappers when comparisons
   are meaningful.

## Bug Fix

1. Reproduce or explain the failing behavior with a focused test.
2. Fix the smallest implementation path that owns the invariant.
3. Add regression coverage named for the behavior, not the bug number.
4. Verify empty, partial, and cleanup paths when the bug touches iteration,
   removal, or storage reuse.

## Documentation Change

1. Keep public Go docs in English.
2. Start exported comments with the exported identifier.
3. Inspect nearby package docs, public interfaces, factory comments, and
   implementation comments before adding new wording.
4. Explain behavioral contracts, empty results, destructive operations,
   invariants, ownership, error behavior, complexity, and performance tradeoffs
   near the API that owns them.
5. Remove comments that merely repeat names or obvious code.
6. Keep README updates focused on when to choose a structure.
7. Run tests for examples when examples change.

## Benchmark Change

1. Keep local benchmarks beside internal implementations.
2. Use `internal/benchmarks` only for optional cross-implementation benchmark
   suites.
3. Place each benchmark concern in the right subpackage:
   - structure contracts, implementation lists, adapters, and optional pure Go
     `stdlib_*.go` baselines in `internal/benchmarks/datastructs`;
   - benchmark input helpers such as random slices and graphs in
     `internal/benchmarks/models`;
   - algorithms that exercise a `datastructs` contract in
     `internal/benchmarks/algorithms`;
   - benchmark entrypoints, suite selection, timer control, and result sinks in
     `internal/benchmarks/suites`.
4. Follow the existing `heap` and `prioritymap` suite shape with a
   `get<Family>Suite(...)` helper returning
   `datastructs.Implementations[...]`.
5. Keep suite files focused on data selection, implementation selection, timer
   control, and workload invocation.
6. Do not put thin method wrappers such as Build, Add, Contains, Remove, Find,
   or iteration in `algorithms`; use local implementation benchmarks for
   operation-level measurements.
7. Confirm that adapters only normalize APIs and preserve each real
   implementation's semantics. Do not invent behavior just to include an
   implementation in a suite.
8. Use `b.ReportAllocs()` and verify that interfaces or variadic calls do not
   create artificial per-operation allocations.
9. Reset or rebuild setup outside the measured section unless setup is the
   operation being measured.
10. Run each new benchmark with a short `-benchtime` before the full suite to
   validate practicality and allocation behavior.

## Commit Preparation

1. Review `git status --short`.
2. Inspect the final diff.
3. Run the appropriate verification.
4. Use Conventional Commits with concise scopes, such as
   `feat(priority map): add radix heap priority map`.
