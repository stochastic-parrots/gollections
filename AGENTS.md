# Agent Instructions

Before changing this repository, read [.agents/project-standards.md](.agents/project-standards.md).
Use [.agents/agent-guide.md](.agents/agent-guide.md) as the fast task map for
what to inspect, edit, and verify.

This repository favors small, explicit Go data-structure implementations with
stable public APIs and concrete implementations hidden under `internal/...`.
Keep changes aligned with the project standards unless the user explicitly asks
for a different direction.

## Required Working Loop

1. Inspect the relevant package and the nearest existing pattern before editing.
2. Keep public API, documentation, tests, and benchmarks in sync when behavior
   changes.
3. Preserve user work in the git tree; never revert unrelated changes.
4. Run the narrowest useful verification first, then `go test ./...` when the
   change can affect shared behavior.
5. Summarize changed files, verification, and any remaining risk.
