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

1. Create a focused branch using the prefixes described in
   [Branch And Pull Request Flow](#branch-and-pull-request-flow), followed by a
   short description.
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

Use these prefixes for change branches targeting `develop`, including branches
from forks:

| Prefix | Purpose | Example |
| --- | --- | --- |
| `feature/` | New capabilities | `feature/hashset` |
| `bugfix/` | Bug fixes | `bugfix/comparator-contracts` |
| `refactor/` | Internal changes that preserve behavior | `refactor/direct-constructors` |
| `docs/` | Documentation changes | `docs/complexity-guidance` |

The [branch-naming ruleset](.github/rulesets/branch-naming.json) blocks creation
of other branch names in this repository, allowing `main`, `develop`, and
`release/*` as well. The required **Pull request target** check validates the
source and target combination, including contributions from forks.

- Open all change pull requests against `develop` and merge them using squash.
  This keeps each change in a single commit.
- When a batch of changes is ready, create a `release/<version>` branch from
  the current `develop` and manually open its promotion pull request against
  `main`. Later changes to `develop` are not automatically included in that
  release.
- Pull requests to `main` must come from a `release/*` branch in this repository
  and use merge commits to preserve the release history. Direct `develop` to
  `main` pull requests are rejected.
- Backport any release fixes to `develop` through `bugfix/*` pull requests.
  Start new change branches from `develop`.
- Use a Conventional Commit style PR title with a concise scope, for example
  `feat(set): add set algebra operations` or
  `refactor(api)!: replace collection factories with direct constructors`.
- A CI check rejects pull requests with any other base/head combination.
  Pull request edits, including changes to the target branch, rerun CI so this
  check evaluates the current branch combination.

CI runs for pull requests and pushes to `main`. A push to `develop` does not
start a duplicate CI run or create a release pull request. Maintainers choose
when to create each release branch and its pull request.

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
behavior. It installs the latest golangci-lint v2 release before each run and
runs gopls hint diagnostics. CI and the pre-commit lint hook also follow the
latest linter release. Local lint needs network access to check for updates;
a new release may introduce diagnostics that need configuration or source
changes. CI uses Go 1.24 for lint and tests both Go 1.24 and stable Go.

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
