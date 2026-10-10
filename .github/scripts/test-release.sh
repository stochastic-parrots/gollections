#!/usr/bin/env bash
set -euo pipefail

script="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)/release.sh"
temporary=$(mktemp -d)
trap 'rm -rf -- "$temporary"' EXIT
repository=stochastic-parrots/gollections

expect_failure() {
  if "$@" > "$temporary/failure.log" 2>&1; then
    echo "Expected command to fail: $*" >&2
    exit 1
  fi
}

git init --bare --initial-branch=main "$temporary/origin.git" > /dev/null
git init --initial-branch=main "$temporary/work" > /dev/null
cd "$temporary/work"
git config user.name 'Release test'
git config user.email release@example.com
git remote add origin "$temporary/origin.git"
printf 'module github.com/%s\n\ngo 1.24\n' "$repository" > go.mod
git add go.mod
git commit -qm 'Initial module'
git switch -qc release/test
echo 'Release change' > release.txt
git add release.txt
git commit -qm 'Release change'
git switch -q main
git merge -q --no-ff release/test -m 'Merge release'
git push -q origin main
commit=$(git rev-parse HEAD)

jq -n --arg repository "$repository" --arg commit "$commit" '{
  title: "release: v1.0.0-beta.2",
  base: {ref: "main"},
  head: {ref: "release/test", repo: {full_name: $repository}},
  merged: true, merge_commit_sha: $commit,
  body: "## Summary\n\nRelease changes.\n\n## Testing\n\nCI passed.",
  html_url: ("https://github.com/" + $repository + "/pull/54")
}' > "$temporary/pr.json"
pr_file="$temporary/pr.json"
invalid_file="$temporary/invalid.json"

# Change PRs need no release title; release branches must be from this repository.
jq '.base.ref = "develop" | .title = "ci: improve releases" | {pull_request: .}' "$pr_file" > "$invalid_file"
bash "$script" check "$invalid_file" "$repository"
for expression in '.base.ref = "other"' '.head.ref = "feature/test"' '.head.repo.full_name = "other/fork"'; do
  jq "$expression" "$pr_file" > "$invalid_file"
  expect_failure bash "$script" check "$invalid_file" "$repository"
done

# Canonical titles and the module major version are checked before publication.
for title in 'chore(release): v1.0.0' 'release: v01.0.0' 'release: v1.0.0-beta.02' 'release: v1.0.0+meta' 'release: v2.0.0-beta.1' $'release: v1.0.0\n'; do
  jq --arg title "$title" '.title = $title' "$pr_file" > "$invalid_file"
  expect_failure bash "$script" check "$invalid_file" "$repository"
done
bash "$script" check "$pr_file" "$repository"

# Policy runs trusted code even when the proposed script accepts every PR.
mkdir -p "$temporary/trusted" "$temporary/proposed/.github/scripts"
cp "$script" "$temporary/trusted/release.sh"
printf 'touch "%s"\nexit 0\n' "$temporary/untrusted-executed" > "$temporary/proposed/.github/scripts/release.sh"
printf 'module github.com/%s/v2\n' "$repository" > "$temporary/proposed/go.mod"
jq '.head.ref = "feature/bypass"' "$pr_file" > "$invalid_file"
awk '
  /^  pr-target:/ { policy = 1; next }
  policy && /^  [a-zA-Z0-9_-]+:/ { exit }
  policy && /^        run: \|/ { body = 1; next }
  body && /^          / { sub(/^          /, ""); print; next }
  body { exit }
' "${script%/scripts/release.sh}/workflows/ci.yml" > "$temporary/branch-step.sh"
[[ -s "$temporary/branch-step.sh" ]]
(
  cd "$temporary/proposed"
  export PR_BASE=main PR_HEAD=feature/bypass PR_HEAD_REPOSITORY="$repository" REPOSITORY="$repository"
  export GITHUB_EVENT_PATH="$invalid_file" GITHUB_REPOSITORY="$repository"
  expect_failure bash "$temporary/branch-step.sh"
)
expect_failure bash "$temporary/trusted/release.sh" check "$invalid_file" "$repository" "$temporary/proposed/go.mod"
expect_failure bash "$temporary/trusted/release.sh" check "$pr_file" "$repository" "$temporary/proposed/go.mod"
jq '.title = "release: v2.0.0-beta.1"' "$pr_file" > "$invalid_file"
bash "$temporary/trusted/release.sh" check "$invalid_file" "$repository" "$temporary/proposed/go.mod"
[[ ! -e "$temporary/untrusted-executed" ]]
expect_failure bash "$script" tag "$pr_file" "$repository" "$temporary/proposed/go.mod"

jq '.merged = false' "$pr_file" > "$invalid_file"
expect_failure bash "$script" tag "$invalid_file" "$repository"
jq --arg commit "$(git rev-parse HEAD^1)" '.merge_commit_sha = $commit' "$pr_file" > "$invalid_file"
expect_failure bash "$script" tag "$invalid_file" "$repository"
[[ -z $(git tag --list) ]]

# An advanced main and a changed worktree module must not change the tagged snapshot.
echo 'Later change' > later.txt
git add later.txt
git commit -qm 'Later change'
git push -q origin main
printf 'module github.com/%s/v2\n' "$repository" > go.mod
bash "$script" tag "$pr_file" "$repository"
[[ $(git --git-dir="$temporary/origin.git" cat-file -t v1.0.0-beta.2) == tag ]]
[[ $(git --git-dir="$temporary/origin.git" rev-parse 'v1.0.0-beta.2^{commit}') == "$commit" ]]
original=$(git rev-parse v1.0.0-beta.2)
bash "$script" tag "$pr_file" "$repository"
[[ $(git rev-parse v1.0.0-beta.2) == "$original" ]]

# Existing versions at other commits are never overwritten.
git tag -a v1.0.0-beta.3 HEAD -m 'Existing version'
git push -q origin refs/tags/v1.0.0-beta.3
jq '.title = "release: v1.0.0-beta.3"' "$pr_file" > "$invalid_file"
expect_failure bash "$script" tag "$invalid_file" "$repository"
[[ $(git --git-dir="$temporary/origin.git" rev-parse 'v1.0.0-beta.3^{commit}') == $(git rev-parse HEAD) ]]

# Mock only the GitHub API transport; tag creation still uses real Git repositories.
mkdir "$temporary/bin"
cat > "$temporary/bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ "$1" == api ]]; then
  if [[ -n ${RELEASE_TEST_ERROR:-} ]]; then
    echo "gh: Forbidden (HTTP $RELEASE_TEST_ERROR)" >&2
    exit 1
  fi
  if [[ ! -f "$RELEASE_TEST_DIR/release.json" ]]; then
    echo 'gh: Not Found (HTTP 404)' >&2
    exit 1
  fi
  cat "$RELEASE_TEST_DIR/release.json"
elif [[ "$1" == release && "$2" == create ]]; then
  printf '%s\n' "$@" > "$RELEASE_TEST_DIR/arguments"
  while [[ $# -gt 0 ]]; do
    if [[ "$1" == --notes-file ]]; then
      cp "$2" "$RELEASE_TEST_DIR/notes"
      break
    fi
    shift
  done
  echo '{"draft":false}' > "$RELEASE_TEST_DIR/release.json"
else
  exit 2
fi
MOCK
chmod +x "$temporary/bin/gh"
export PATH="$temporary/bin:$PATH"
export RELEASE_TEST_DIR="$temporary"
bash "$script" publish "$pr_file" "$repository"
[[ $(cat "$temporary/arguments") == *--prerelease* ]]
[[ $(cat "$temporary/arguments") == *--latest=false* ]]
[[ $(cat "$temporary/arguments") == *--verify-tag* ]]
[[ $(cat "$temporary/notes") == *'Release changes.'* ]]
[[ $(cat "$temporary/notes") != *'CI passed.'* ]]
echo 'Manually edited release notes' > "$temporary/notes"
bash "$script" publish "$pr_file" "$repository"
[[ $(cat "$temporary/notes") == 'Manually edited release notes' ]]
RELEASE_TEST_ERROR=403 expect_failure bash "$script" publish "$pr_file" "$repository"
[[ $(cat "$temporary/notes") == 'Manually edited release notes' ]]

rm "$temporary/release.json"
jq '.title = "release: v1.0.0"' "$pr_file" > "$invalid_file"
bash "$script" publish "$invalid_file" "$repository"
[[ $(cat "$temporary/arguments") == *--latest* ]]
[[ $(cat "$temporary/arguments") != *--latest=false* ]]
[[ $(cat "$temporary/arguments") != *--prerelease* ]]
echo 'Release validation and Git integration checks passed'
