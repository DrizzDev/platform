# Code Standards

Status: Approved and mandatory

The mandatory structure and naming policy is defined in
[style.md](style.md). Language profiles may add syntax rules but may not weaken
it.

- `CODE-001`: Names MUST follow the strict single-word and domain-vocabulary
  rules in `STYLE-001` through `STYLE-015`.
- `CODE-002`: Multiword concepts MUST use the nested structure required by
  `STYLE-003`, `STYLE-004`, and `STYLE-011`.
- `CODE-003`: Generic names such as `base`, `common`, `helper`, `manager`,
  `misc`, `processor`, `service`, `types`, `util`, and `wrapper` are prohibited
  unless the word is the exact domain or platform concept.
- `CODE-004`: `Impl` is prohibited. Concrete names MUST identify the strategy or
  provider.
- `CODE-005`: Established initialisms retain language-idiomatic casing.
- `CODE-006`: Magic values are prohibited. Constants remain with their owning
  module and describe meaning rather than merely repeating a literal. Internal
  semantic string values use `UPPER_SNAKE_CASE`; Go identifiers do not.
- `CODE-007`: Inputs and outputs crossing a boundary MUST use explicit typed
  schemas. Raw maps, tuples, or opaque payloads require a documented outermost
  boundary and immediate validation.
- `CODE-008`: Domain values SHOULD be immutable. Identifiers, digests, amounts,
  durations, revisions, and states SHOULD use defined types where ambiguity is
  possible.
- `CODE-009`: Public surface area MUST be minimal and deliberate.
- `CODE-010`: Do not write a comment unless the code contains the
  exceptional technical need defined by `STYLE-030`. Permitted comments follow
  the one-line preference and three-line hard limit in `STYLE-034`.
- `CODE-011`: AI filler, emojis in code artifacts, and TODOs without a tracking
  identifier are prohibited.
- `CODE-012`: Hidden side effects, mutable singletons, and implicit dependency
  lookup are prohibited.
- `CODE-013`: Time, randomness, environment, locale, identity, and other
  nondeterministic inputs MUST be injected when they affect behavior or tests.
- `CODE-014`: External input MUST be bounded, validated, and normalized before
  application or domain use.
- `CODE-015`: Provider responses MUST be treated as untrusted and validated.
- `CODE-016`: Expected failures MUST be typed or classifiable. Errors MUST
  preserve safe context and internal causes without string matching.
- `CODE-017`: Errors exposed across trust boundaries MUST NOT include secrets,
  content, stack traces, SQL, internal paths, or provider internals.
- `CODE-018`: Panic or fatal termination is reserved for unrecoverable startup
  programmer/configuration failure at the composition root.
- `CODE-019`: Resources MUST have a guaranteed release path.
- `CODE-020`: Reuse existing correct code; refactor recoverable code; rewrite
  only when evidence shows the structure is fundamentally unsuitable.

## Contracts and layering

- `CODE-032`: A call across a layer, module, process, storage, provider, or
  public boundary MUST use an explicit typed request, response, error, and state
  model appropriate to that boundary.
- `CODE-033`: Raw maps, `any`, provider SDK objects, database records, generated
  transport objects, and framework objects MUST NOT cross into another layer.
- `CODE-034`: Constants, enumerations, contracts, domain models, errors,
  repositories, use cases, adapters, handlers, and validators MUST remain in
  their owning layer and responsibility.
- `CODE-035`: A contract change MUST be evaluated for validation,
  compatibility, migration, failure mapping, and tests before implementation.

## Configuration

- `CODE-024`: Configuration MUST be typed, validated before dependent
  components start, and immutable after composition unless an explicit runtime
  policy contract says otherwise.
- `CODE-025`: Configuration sources, precedence, unknown-key behavior, and safe
  defaults MUST be explicit and tested.
- `CODE-031`: Invalid user or deployment configuration MUST return a normal
  actionable startup error. Panic is limited to impossible programmer-created
  composition states.
- `CODE-026`: Domain MUST NOT read configuration. Application receives typed
  values or policy facts through construction and use-case input.
- `CODE-027`: Secrets MUST be represented by credential references or secret
  ports, not ordinary configuration values.
- `CODE-028`: Diagnostics MAY expose effective non-secret configuration through
  an allowlist and mandatory redaction.

## Error contract

- `CODE-029`: Boundary failures MUST map from a transport-neutral error
  containing a stable code, category, retryability, recommended action, safe
  details, correlation identity, and optional retry time.
- `CODE-030`: Adapters MUST map that error into native CLI, MCP, IPC, or HTTP
  outcomes without exposing provider details.

## Documentation

- `CODE-021`: A public external contract documents purpose, ownership,
  cancellation, idempotency, and material side effects concisely. Internal
  exported Go identifiers do not require comments when naming and structure are
  self-explanatory.
- `CODE-022`: Documentation MUST change with public behavior, configuration,
  operational behavior, compatibility, or architecture.
- `CODE-023`: Generated code and generated contracts MUST be reproducible and
  MUST NOT be edited manually.
