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
  Open release pull requests manually with the title
  `release: vMAJOR.MINOR.PATCH[-prerelease]`, for example `release: v1.0.0-beta.2`.
- Keep `develop` current with `main` after release promotions. Start a new
  change branch from `develop`.
- Use a Conventional Commit style PR title with a concise scope, for example
  `feat(set): add set algebra operations` or
  `refactor(api)!: replace collection factories with direct constructors`.
- The **Pull request target** check validates the branches, release title, and
  compatibility between the version and the module path in `go.mod`.

CI runs for pull requests and pushes to `main`. Editing a PR title or base also
reruns CI so release validation uses the current metadata. The workflows do not
create pull requests or change their titles and descriptions.

For CLI merges, use `gh pr merge --squash` when the target is `develop` and
`gh pr merge --merge` when the target is `main`.

## Release Tags

After a release PR merges into `main`, **Publish release** creates an annotated
tag using the version in its title and publishes a GitHub Release. The tag
points to that PR's merge commit, even if `main` has advanced before the
workflow starts. This publishes a Go module version; review the version and
release contents before merging.

Use canonical versions such as `v1.0.0-beta.2` or `v1.0.0`, without build
metadata. Go requires a `/vN` module path suffix for major versions 2 and above,
including prereleases. The current module path has no suffix, so releases use
major version 0 or 1. A v2 release requires migrating the module path and imports
to `/v2` first.

The workflow runs the automation from `main`, validates the merged PR and the
module at its merge commit, and uses `GITHUB_TOKEN` with `contents: write` to
push the tag. It does not need permission to create or approve pull requests.
An existing tag at the same commit makes a retry a no-op; an existing tag at a
different commit fails without overwriting it. If a run fails, rerun it or
dispatch **Publish release** from `main` with the already-merged PR number.
Release notes use the PR's `Summary` section and link back to the PR. Versions
with a prerelease suffix are marked **Pre-release**, without **Latest**; stable
versions are marked **Latest**. GitHub does not allow prereleases to be marked
as latest. Retries preserve an existing published release's notes and settings.
The tag push and release publication use `GITHUB_TOKEN` and do not trigger
additional workflows.

The automation uses Bash, Git, `jq`, and GitHub CLI, already available on the
GitHub-hosted Ubuntu runner. Test validation and tagging locally with temporary
Git repositories:

```bash
bash .github/scripts/test-release.sh
```

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
CI enforces at least 95% statement coverage across library packages, excluding
`internal/benchmarks`. Codecov checks 100% coverage of changed lines. Branch
rules require both Go test jobs, lint, coverage upload, PR validation, and the
Codecov patch check; `main` also requires the managed Code Quality analysis.
The pre-commit test hook and `make test` keep the race detector enabled.

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
