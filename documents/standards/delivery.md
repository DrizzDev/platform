# Delivery Standards

Status: Approved and mandatory

## Dependencies

- `DEL-001`: Prefer the standard library and existing approved dependencies.
- `DEL-002`: Every direct dependency MUST record purpose, owner, pinned version,
  license, maintenance, vulnerabilities, transitive and binary-size effect,
  cross-platform or CGO effect, and replacement boundary.
- `DEL-003`: Overlapping frameworks or persistence paradigms are prohibited
  without an ADR and migration plan.
- `DEL-004`: Vendor SDKs MUST remain in adapters and MUST NOT enter domain or
  application.
- `DEL-005`: Copied reference code requires license review, preserved notices,
  security review, and explicit provenance.
- `DEL-019`: Approved licenses are permissive licenses recorded by repository
  policy. Any other license requires legal review, an ADR, an owner, and a
  documented distribution impact before the dependency is added.

## Review

- `DEL-006`: A change MUST identify product behavior, non-goals, module, layer,
  trust boundaries, state, failure, compatibility, tests, and documentation.
- `DEL-007`: Public contract, generated code, migration, dependency, security,
  and authorization changes require explicit review.
- `DEL-008`: A change MUST be independently reviewable and MUST NOT mix
  unrelated refactoring or formatting.
- `DEL-009`: A standards exception requires owner, reason, scope, risk,
  compensating control, expiry or review trigger, and ADR when architectural.
- `DEL-021`: Review MUST evaluate SOLID, separation of concerns, cohesion,
  coupling, composition, dependency inversion, and correct layer ownership.
- `DEL-022`: Review MUST confirm that interface and abstraction needs,
  alternative implementations, provider coupling, independent testability, SRP,
  OCP, dependency direction, and future extraction were evaluated.
- `DEL-023`: Review MUST reject generic or implementation-specific data crossing
  a boundary.
- `DEL-024`: Review MUST reject a project-owned Go source file exceeding 500
  physical lines and MUST treat a file approaching the limit as a cohesion
  smell.
- `DEL-025`: Before product implementation begins, repository verification and
  CI MUST count physical lines in every project-owned Go source file and block
  any file above 500.
- `DEL-026`: Architecture rules that can be checked mechanically, including
  forbidden dependency direction and boundary type leakage, MUST become blocking
  repository checks before affected production code is accepted.

### Required code review checklist

- [ ] Does the change preserve one responsibility per component?
- [ ] Are cohesion high and coupling explicit and narrow?
- [ ] Do dependencies point toward the layer that owns the policy?
- [ ] Was the need for an interface or abstraction evaluated?
- [ ] Can a second implementation be added without rewriting consumers?
- [ ] Is the component independently testable through production boundaries?
- [ ] Does the design preserve SRP, OCP, substitution, and interface
      segregation?
- [ ] Is composition used instead of inheritance or behavioral embedding?
- [ ] Does every concern live in its owning layer, directory, and file?
- [ ] Are all boundary contracts explicit and strongly typed?
- [ ] Can the module be extracted later without exposing storage or provider
      internals?
- [ ] Is every project-owned Go source file within 500 physical lines?
- [ ] Does each pattern, framework, dependency, and abstraction solve a named
      current problem?

## Verification

- `DEL-010`: One repository-owned aggregate command MUST represent the merge
  gate.
- `DEL-011`: The merge gate MUST include every implemented repository check
  applicable to the change, including format, build, architecture, unit,
  integration, contract, concurrency, static analysis, vulnerability, license,
  secret, generated-code, and migration checks where they exist. A skipped
  required gate needs an explicit reason and owner.
- `DEL-012`: CI actions and tools MUST be pinned reproducibly.
- `DEL-013`: Auto-fixes MUST be deterministic and semantics-preserving.
- `DEL-014`: A clean checkout MUST reproduce generated code and release
  build inputs and unsigned artifacts. Signed or notarized output is not
  required to be byte-identical when its trusted timestamp or signature format
  is intentionally nondeterministic; provenance MUST bind it to the reproducible
  unsigned input.

## Completion

- `DEL-015`: Work is done only when approved behavior and non-goals are met,
  required gates pass, documentation and ADRs are current, and risks are owned.
- `DEL-016`: A phase or feature is not complete with only interfaces, mocks, or
  unit tests when its value depends on a real integration.
- `DEL-017`: Release qualification MUST include install, upgrade, rollback,
  restore, migration, resource, security, compatibility, and support evidence
  on the supported matrix.
- `DEL-018`: Claims in handoff MUST distinguish verified fact, inference,
  proposal, and unresolved risk.
- `DEL-020`: Verification claims record exact revision, working-tree
  fingerprint, command, environment, and result. Scoped checks are labelled
  scoped. Claims of pre-existing, flaky, or unrelated failures require
  reproduction plus clean-baseline or parent-revision comparison.
