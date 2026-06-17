---
paths:
  - "**/*.md"
  - "**/doc.go"
  - "**/examples_test.go"
---

# Documentation Rules

- Read the `Documentation` section of `.agents/project-standards.md`.
- Public Go documentation is written in English.
- Exported package, type, interface, and function comments start with the
  exported identifier.
- Constructor docs explain when to choose that implementation and include
  performance tables when they define a data-structure choice.
- Destructive methods say they are destructive in the first paragraph.
- README updates should explain when to choose a structure, not only list that
  it exists.
- Examples live in `examples_test.go`, use `Example...` names, and include
  `// Output:` blocks when deterministic.
