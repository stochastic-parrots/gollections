#!/usr/bin/env bash
set -euo pipefail

fail() { echo "$*" >&2; exit 1; }

[[ $# == 3 ]] || fail "Usage: bash release.sh check|tag|publish PR_JSON REPOSITORY"
command=$1
pr=$(jq '.pull_request // .' "$2")
repository=$3
[[ "$command" == check || "$command" == tag || "$command" == publish ]] || fail "Unknown command: $command"

base=$(jq -er '.base.ref' <<< "$pr")
if [[ "$command" == check && "$base" == develop ]]; then
  exit 0
fi
[[ "$base" == main ]] || fail "Change PRs must target develop; releases must target main"
head=$(jq -er '.head.ref' <<< "$pr")
head_repository=$(jq -r '.head.repo.full_name // ""' <<< "$pr")
[[ "$head_repository" == "$repository" && ( "$head" == develop || "$head" == release/* ) ]] ||
  fail "Releases must come from develop or release/* in this repository"

jq -e '.title | test("[\\r\\n]") | not' <<< "$pr" > /dev/null || fail "Release title must be a single line"
title=$(jq -er '.title' <<< "$pr")
pattern='^release: (v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-([0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*))?)$'
[[ "$title" =~ $pattern ]] || fail "Release title must be release: vMAJOR.MINOR.PATCH[-prerelease]"
version=${BASH_REMATCH[1]}
major=${BASH_REMATCH[2]}
prerelease=${BASH_REMATCH[6]:-}
IFS='.' read -r -a identifiers <<< "$prerelease"
for identifier in "${identifiers[@]}"; do
  [[ ! "$identifier" =~ ^0[0-9]+$ ]] || fail "Numeric prerelease identifiers cannot have leading zeroes"
done

if [[ "$command" != check ]]; then
  jq -e '.merged == true' <<< "$pr" > /dev/null || fail "Only merged release PRs can be tagged"
  commit=$(jq -er '.merge_commit_sha' <<< "$pr")
  [[ "$commit" =~ ^[0-9a-f]{40}$ ]] || fail "Missing or invalid merge commit"
  git merge-base --is-ancestor "$commit" origin/main || fail "Release commit is not in main"
  read -r -a parents < <(git rev-list --parents -n 1 "$commit")
  [[ ${#parents[@]} == 3 ]] || fail "Releases must use a merge commit"
  go_mod=$(git show "$commit:go.mod")
else
  go_mod=$(cat go.mod)
fi

module=$(awk '$1 == "module" { gsub(/"/, "", $2); print $2; exit }' <<< "$go_mod")
expected="github.com/$repository"
if [[ "$major" != 0 && "$major" != 1 ]]; then
  expected+="/v$major"
fi
[[ "$module" == "$expected" ]] || fail "$version requires go.mod module $expected"
if [[ "$command" == check ]]; then
  echo "Valid release pull request: $version"
  exit 0
fi

tag_ref="refs/tags/$version"
if git show-ref --verify --quiet "$tag_ref"; then
  [[ $(git rev-parse "$tag_ref^{commit}") == "$commit" ]] ||
    fail "$version already points to a different commit"
  echo "$version already points to $commit; keeping the tag"
else
  status=$?
  [[ "$status" == 1 ]] || exit "$status"
  git tag -a "$version" "$commit" -m "Release $version" -m "$(jq -er '.html_url' <<< "$pr")"
  git push origin "$tag_ref"
  echo "Created annotated tag $version at $commit"
fi
[[ "$command" == publish ]] || exit 0

release_error=$(mktemp)
notes_file=$(mktemp)
trap 'rm -f -- "$release_error" "$notes_file"' EXIT
if existing=$(gh api "repos/$repository/releases/tags/$version" 2> "$release_error"); then
  jq -e '.draft == false' <<< "$existing" > /dev/null || fail "Existing release is still a draft"
  echo "GitHub Release $version already exists; keeping its notes and settings"
  exit 0
fi
[[ $(cat "$release_error") == *'HTTP 404'* ]] || fail "$(cat "$release_error")"

{
  printf 'Release PR: %s\n\n' "$(jq -er '.html_url' <<< "$pr")"
  jq -r '.body // ""' <<< "$pr" | awk '
    /^## Summary[[:space:]]*$/ { summary = 1; next }
    /^## / { summary = 0 }
    summary { print }
  '
} > "$notes_file"
flags=(--latest)
if [[ -n "$prerelease" ]]; then
  flags=(--prerelease --latest=false)
fi
gh release create "$version" --repo "$repository" --verify-tag \
  --title "$version" --notes-file "$notes_file" "${flags[@]}"
