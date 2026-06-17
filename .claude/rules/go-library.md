---
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
---

# Go Library Rules

- Read `.agents/project-standards.md` before changing Go code.
- Keep root interfaces read-only; mutating APIs belong in structure-specific
  packages.
- Keep concrete implementations under `internal/...`; expose public type
  aliases and constructors from public packages.
- Empty reads and removals return zero values plus `false` or package-owned
  errors.
- Iterators stop immediately when `yield` returns `false`.
- `Clear` and removal paths must zero discarded references and preserve
  reusable capacity where practical.
- Use `cmp.Ordered`, `constraint`, or `internal/comparator` according to the
  existing package pattern.
- Run `gofmt` on changed Go files.
