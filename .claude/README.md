# Claude Code Setup

Claude Code reads the root `CLAUDE.md`, which imports `AGENTS.md` as the shared
agent entrypoint. This directory holds Claude-specific helpers that should not
be duplicated in the root files.

## Layout

- `rules/`: path-scoped reminders loaded by Claude Code when matching files are
  in scope.

Keep canonical project rules in `.agents/project-standards.md`. Use this
directory for Claude Code loading behavior only.
