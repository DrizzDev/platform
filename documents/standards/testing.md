# Testing Standards

Status: Approved and mandatory

- `TEST-001`: Test observable behavior, contracts, invariants, and failure
  recovery rather than implementation call order.
- `TEST-002`: Domain tests MUST be deterministic and infrastructure-free.
- `TEST-003`: Application tests SHOULD use small hand-written fakes. Generated
  mocks require a large or externally owned interface that justifies them.
- `TEST-004`: Infrastructure and adapters MUST have contract and integration
  tests against real stores, protocol fixtures, or production-shaped
  simulators.
- `TEST-005`: End-to-end tests MUST use normal process composition and public
  boundaries.
- `TEST-006`: Every defect fix MUST include a regression test that fails for the
  original defect.
- `TEST-007`: Tests MUST inject clock, randomness, identifiers, retry, and
  provider behavior. Sleeps and developer machine state are prohibited as
  coordination.
- `TEST-008`: Concurrent code MUST pass race and leak checks appropriate to the
  language.
- `TEST-009`: External parsers, schemas, state machines, and migration readers
  SHOULD have fuzz or property tests.
- `TEST-010`: Public contracts MUST keep minimum, previous, current,
  next-additive, unknown, malformed, oversized, duplicate, replayed, partial,
  and privacy-canary fixtures when the contract supports the corresponding
  versioning, parsing, replay, cardinality, or sensitive-data risk. `N/A`
  categories are recorded with reasons in the contract inventory.
- `TEST-011`: Persistence MUST test empty install, supported upgrades,
  interrupted migration, corruption, low storage, large backlog, contention,
  rollback, and query plans for hot paths.
- `TEST-012`: Durable workflows MUST inject failure at every transition and
  verify crash, kill, restart, duplicate, reorder, timeout, partial response,
  network loss, and disk pressure.
- `TEST-013`: Physical or environment-specific behavior MUST NOT ship using only
  unit, emulator, or simulator evidence.
- `TEST-014`: A skipped, disabled, or flaky test in a merge gate is prohibited.
  Quarantine requires an owner, tracked issue, expiry, and separate visible
  failing signal.
- `TEST-015`: Verification evidence MUST record command, revision, environment,
  result, and artifacts needed to reproduce the claim.
