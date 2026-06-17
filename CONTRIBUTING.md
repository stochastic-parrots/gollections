# Contributing

Thanks for helping improve `gollections`. This project welcomes focused
contributions that preserve the library's API clarity, documentation quality,
and data-structure invariants.

## Ways to Contribute

- Report bugs with clear reproduction steps and expected behavior.
- Suggest focused features or API improvements.
- Improve package documentation, examples, or README guidance.
- Add missing tests or benchmarks for existing behavior.
- Submit small implementation fixes or new data-structure work that follows the
  existing package patterns.

## Before Contributing

1. Search existing issues and pull requests to avoid duplicate work.
2. Open an issue before investing in large changes, public API changes,
   documented behavior changes, or new data-structure families.
3. Read the README and the documentation for the package you plan to change.
4. Inspect the nearest existing implementation, test, documentation, or
   benchmark pattern before editing.
5. Check [.agents/project-standards.md](.agents/project-standards.md) when the
   change touches project conventions such as public APIs, package layout,
   tests, benchmarks, or documentation.

## Development Workflow

1. Create a focused branch for the change.
2. Keep the diff scoped to one bug, feature, documentation update, test gap, or
   benchmark gap.
3. Follow the existing package layout: public APIs live in focused packages,
   while concrete implementations stay under `internal/...`.
4. Keep public APIs, documentation, tests, benchmarks, and README guidance in
   sync when behavior changes.
5. Run `gofmt` on changed Go files.
6. Run the narrowest useful verification while iterating, then broaden checks
   before opening a PR.

## Testing

For most code changes, run:

```bash
go test ./...
```

For package-specific iteration, run the affected package first, for example:

```bash
go test ./internal/prioritymap
```

For coverage-sensitive internal changes, also inspect coverage:

```bash
go test ./internal/prioritymap -coverprofile=/tmp/internal_prioritymap.cover -covermode=count
go tool cover -func=/tmp/internal_prioritymap.cover
```

For race-sensitive changes, run:

```bash
go test -race ./...
```

Use `make lint` when changing public APIs, docs, or broad implementation
behavior.

## Pull Requests

Before opening a PR:

- Confirm the diff is focused and does not include unrelated local changes.
- Explain what changed and why.
- Add or update tests for behavior changes.
- Update docs, examples, or README guidance when public behavior changes.
- Add or update benchmarks when performance-sensitive behavior changes.
- Include verification commands and results in the PR description.
- Call out risk, migration impact, or follow-up work when relevant.

## Commit Messages

Commit messages should follow Conventional Commits with concise scopes, such
as:

```text
feat(priority map): add radix heap priority map
test(priority map): cover radix heap priority map
doc(priority map): document radix heap priority map
perf(priority map): add radix heap benchmarks
fix(priority map): clear stale priority state
```

## Optional AI Assistance

AI tools may be used, but contributors are responsible for the final patch. Do
not commit secrets, personal tokens, local prompts, or private tool output, and
verify generated changes the same way as handwritten changes.
