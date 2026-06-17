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
3. Explain behavioral contracts, empty results, destructive operations, and
   performance tradeoffs near the API that owns them.
4. Keep README updates focused on when to choose a structure.
5. Run tests for examples when examples change.

## Benchmark Change

1. Keep local benchmarks beside internal implementations.
2. Keep cross-implementation suites under `internal/benchmarks/suites`.
3. Use `b.ReportAllocs()`.
4. Reset or rebuild setup outside the measured section unless setup is the
   operation being measured.

## Commit Preparation

1. Review `git status --short`.
2. Inspect the final diff.
3. Run the appropriate verification.
4. Use Conventional Commits with concise scopes, such as
   `feat(priority map): add radix heap priority map`.
