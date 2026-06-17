# Agent Guides

This directory contains shared working instructions for AI coding agents such
as Codex, Claude Code, and IDE assistants.

Start with [agent-guide.md](agent-guide.md) for a compact task map, then read
the relevant sections of [project-standards.md](project-standards.md) before
editing. Root-level `AGENTS.md` and `CLAUDE.md` point here so agent tools can
discover the same guidance from their usual entry points.

## Files

- [agent-guide.md](agent-guide.md): short operational guide for agents.
- [task-checklists.md](task-checklists.md): repeatable checklists for common
  coding, test, documentation, benchmark, and release tasks.
- [project-standards.md](project-standards.md): canonical conventions for
  architecture, naming, documentation, tests, benchmarks, and commit style.

Keep these files factual, concise, and free of personal preferences. If a rule
only matters to Claude Code path loading, put the short reminder in
`.claude/rules/` and keep the canonical version here.
