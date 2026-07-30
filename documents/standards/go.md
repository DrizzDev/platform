# Go Standards

Status: Approved and mandatory

- `GO-001`: Use the repository-pinned Go toolchain and standard Go formatting.
- `GO-002`: Package, folder, and source base names MUST use one complete
  semantic word. Multiword concepts MUST use nested packages. Required Go
  filenames and suffixes are the exceptions listed in `STYLE-005`.
- `GO-003`: Exported identifiers use PascalCase; unexported identifiers use
  lowerCamelCase; established initialisms retain conventional casing.
- `GO-004`: Descriptive names are required. Single-letter variables are limited
  to narrow idiomatic scopes where meaning is obvious.
- `GO-005`: Constant and enumeration-like identifiers use normal Go naming:
  exported `PascalCase`, unexported `lowerCamelCase`, with initialisms kept
  intact. Internal semantic string values use `UPPER_SNAKE_CASE`. External
  values preserve their reviewed protocol vocabulary. Numeric and `iota`
  values have no casing.
- `GO-006`: Accept interfaces at the consumer boundary and return concrete
  types by default.
- `GO-007`: Test convenience alone is not sufficient reason for an interface.
  Tests may fake an interface required by a production application boundary.
- `GO-008`: `context.Context` is the first parameter for I/O and cancellable
  work and MUST NOT be stored in a struct. It carries cancellation, deadlines,
  and tracing, not domain data or optional parameters.
- `GO-009`: Errors are values. Preserve causes with wrapping and classify stable
  outcomes with typed errors or codes.
- `GO-010`: Avoid `any` in domain and application. At an opaque protocol
  boundary, validate and narrow it immediately.
- `GO-011`: Prefer useful zero values for technical types. Invariant-bearing
  domain values use validated construction or reject invalid values at the
  owning boundary. Mutable exported fields MUST NOT bypass invariants.
- `GO-012`: Goroutines have an owner, cancellation, bounded concurrency, and
  shutdown test. Related concurrent work SHOULD use structured cancellation.
- `GO-013`: Channels MUST have an explained owner and capacity. Unbounded
  goroutine creation is prohibited.
- `GO-014`: Production packages MUST NOT call `os.Exit` or panic outside the
  composition root.
- `GO-015`: Export does not automatically require a comment. Add concise Go
  documentation only for a public external contract or a technical abstraction
  that is not clear from its name and type. Follow the three-line hard limit in
  `STYLE-034`.
- `GO-016`: Same-package tests are reserved for necessary unexported invariants.
  Do not export production internals solely for tests.
- `GO-017`: Use raw SQL only inside persistence adapters and test its semantics.
  Hot-path SQL requires query-plan or benchmark evidence.
- `GO-018`: Run gofmt, goimports, vet, staticcheck, configured linters, tests,
  race checks for concurrency, and govulncheck.
- `GO-019`: Cancellation and deadline errors MUST propagate with their identity
  intact and MUST NOT be swallowed or converted into an unrelated failure.
- `GO-020`: Every project-owned `.go` file, including tests and generated
  project code, MUST contain no more than 500 physical lines.

## Owner-convention mapping

- Behavior with collaborators or lifecycle lives on a cohesive Go type.
- Package functions are limited to constructors and required entry points.
  Pure transformations use methods on their owning cohesive type.
- Entities, value objects, configuration, and boundary payloads use typed
  structs with validated constructors and narrow methods.
- Internal implementation details use unexported Go identifiers and narrow
  package boundaries.
- Typed parameter structs replace keyword-call requirements when argument
  meaning is not immediately obvious.
- Public methods orchestrate private steps and follow the line targets in
  `STYLE-024`.
- Test functions required by Go are framework adapters; their organization
  mirrors the production type or use case and covers concrete production,
  failure, and recovery behavior.
- A typed parameter struct is required when several arguments would otherwise
  be ambiguous or easy to transpose.
