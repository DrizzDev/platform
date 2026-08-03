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
- `DEL-006`: Approved licenses are permissive licenses recorded by repository
  policy. Any other license requires legal review, an ADR, an owner, and a
  documented distribution impact before the dependency is added.

## Review

- `DEL-007`: A change MUST identify product behavior, non-goals, module, layer,
  trust boundaries, state, failure, compatibility, tests, and documentation.
- `DEL-008`: Public contract, generated code, migration, dependency, security,
  and authorization changes require explicit review.
- `DEL-009`: A change MUST be independently reviewable and MUST NOT mix
  unrelated refactoring or formatting.
- `DEL-010`: A standards exception requires owner, reason, scope, risk,
  compensating control, expiry or review trigger, and ADR when architectural.
- `DEL-011`: Review MUST evaluate SOLID, separation of concerns, cohesion,
  coupling, composition, dependency inversion, and correct layer ownership.
- `DEL-012`: Review MUST confirm that interface and abstraction needs,
  alternative implementations, provider coupling, independent testability, SRP,
  OCP, dependency direction, and future extraction were evaluated.
- `DEL-013`: Review MUST reject generic or implementation-specific data crossing
  a boundary.
- `DEL-014`: Review MUST reject a project-owned Go source file exceeding 500
  physical lines and MUST treat a file approaching the limit as a cohesion
  smell.
- `DEL-015`: Before product implementation begins, repository verification and
  CI MUST count physical lines in every project-owned Go source file and block
  any file above 500.
- `DEL-016`: Architecture rules that can be checked mechanically, including
  forbidden dependency direction and boundary type leakage, MUST become blocking
  repository checks before affected production code is accepted.
- `DEL-017`: Every implementation slice MUST include the design inventory
  defined by `standards/agent.md`. Review MUST compare the final diff with that
  inventory and reject silent changes in ownership, layers, contracts,
  dependencies, state, failures, files, tests, or verification.

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

- `DEL-018`: One repository-owned aggregate command MUST represent the merge
  gate.
- `DEL-019`: The merge gate MUST include every implemented repository check
  applicable to the change, including format, build, architecture, unit,
  integration, contract, concurrency, static analysis, vulnerability, license,
  secret, generated-code, and migration checks where they exist. A skipped
  required gate needs an explicit reason and owner.
- `DEL-020`: CI actions and tools MUST be pinned reproducibly.
- `DEL-021`: Auto-fixes MUST be deterministic and semantics-preserving.
- `DEL-022`: A clean checkout MUST reproduce generated code and release
  build inputs and unsigned artifacts. Signed or notarized output is not
  required to be byte-identical when its trusted timestamp or signature format
  is intentionally nondeterministic; provenance MUST bind it to the reproducible
  unsigned input.
- `DEL-023`: Fast deterministic checks run before a local commit. The complete
  merge gate runs before a local push and on every pull request. Expensive
  checks may be omitted from the commit hook only when the push hook and pull
  request gate run them.
- `DEL-024`: Pull request verification MUST be a required branch or repository
  ruleset check. A workflow file alone runs checks but does not prevent a merge
  when repository protection is absent.
- `DEL-025`: Repository verification MUST mechanically enforce every applicable
  measurable standard, including forbidden names, dependency direction, file
  size, complexity, argument limits, test-package isolation, domain
  mutability, context propagation, framework isolation, and direct operating
  system access. Rules that cannot be checked mechanically remain explicit
  review obligations.

## Completion

- `DEL-026`: Work is done only when approved behavior and non-goals are met,
  required gates pass, documentation and ADRs are current, and risks are owned.
- `DEL-027`: A phase or feature is not complete with only interfaces, mocks, or
  unit tests when its value depends on a real integration.
- `DEL-028`: Release qualification MUST include install, upgrade, rollback,
  restore, migration, resource, security, compatibility, and support evidence
  on the supported matrix.
- `DEL-029`: Claims in handoff MUST distinguish verified fact, inference,
  proposal, and unresolved risk.
- `DEL-030`: Verification claims record exact revision, working-tree
  fingerprint, command, environment, and result. Scoped checks are labelled
  scoped. Claims of pre-existing, flaky, or unrelated failures require
  reproduction plus clean-baseline or parent-revision comparison.
