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

1. Create a focused branch for the change, named with its kind and a short
   description, such as `feature/hashset`, `fix/comparator-contracts`,
   `refactor/direct-constructors`, or `docs/complexity-guidance`.
2. Keep the diff scoped to one bug, feature, documentation update, test gap, or
   benchmark gap.
3. Follow the existing package layout: public APIs live in focused packages,
   while concrete implementations stay under `internal/...`.
4. Keep public APIs, documentation, tests, benchmarks, and README guidance in
   sync when behavior changes.
5. Run `gofmt` on changed Go files.
6. Run the narrowest useful verification while iterating, then broaden checks
   before opening a PR.

## Branch And Pull Request Flow

- Open all change pull requests against `develop` and merge them using squash.
  This keeps each feature or fix in a single commit.
- Pull requests to `main` must come from `develop` or a `release/*` branch in
  this repository and use merge commits to preserve the release history.
  GitHub Actions opens or updates the `develop` promotion after pushes to
  `develop`; open `release/*` promotions manually.
- Keep `develop` current with `main` after release promotions. Start a new
  change branch from `develop`.
- Use a Conventional Commit style PR title with a concise scope, for example
  `feat(set): add set algebra operations` or
  `refactor(api)!: replace collection factories with direct constructors`.
- A CI check rejects pull requests with any other base/head combination.

The release promotion workflow needs GitHub Actions' **Allow GitHub Actions to
create and approve pull requests** repository setting. It grants the workflow
only `contents: read` and `pull-requests: write`.

CI runs for pull requests and pushes to `main`. A push to `develop` updates the
release pull request without starting a second CI run for the same change.
Pull requests opened by the workflow's `GITHUB_TOKEN` may require a maintainer
to approve their workflow runs in GitHub before the checks start.

For CLI merges, use `gh pr merge --squash` when the target is `develop` and
`gh pr merge --merge` when the target is `main`.

## Dependencies

Production Go files in every package may import only the standard library and
packages from this module. This also applies to the benchmark harness.
Third-party testing libraries are allowed only in `*_test.go` files; the
current allowlist contains `github.com/stretchr/testify/`. Approve any additional
testing library by updating the test allowlist in `.golangci.yml`.

Packages under `internal/...`, including their tests, may import other packages
from this module only under `internal/...`.

The `depguard` rules in `.golangci.yml` enforce both boundaries through
`make lint` locally and the **Lint** check in CI.

## Testing

For most code changes, run:

```bash
go test ./...
```

For package-specific iteration, run the affected package first, for example:

```bash
go test ./internal/prioritymap
go test ./internal/set
```

For coverage-sensitive internal changes, also inspect coverage:

```bash
go test ./internal/prioritymap -coverprofile=/tmp/internal_prioritymap.cover -covermode=count
go tool cover -func=/tmp/internal_prioritymap.cover

go test ./internal/set -coverprofile=/tmp/internal_set.cover -covermode=count
go tool cover -func=/tmp/internal_set.cover
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
