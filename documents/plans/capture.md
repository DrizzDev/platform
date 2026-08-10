# Stage 5 — Capture, Records, and Synchronization Implementation Plan

Status: In progress. Companion architecture: [Agent integration and execution
capture](../capture.md) and [ADR 0007](../decisions/0007-capture.md); persistence
per [ADR 0004](../decisions/0004-persistence.md). Delivery order per the
[roadmap](../roadmap.md) Stage 5.

## 1. Goal

Preserve the required capability and agent execution history and synchronize it to
Drizz without adding a cloud round trip to each local device step. The ordered,
durable local execution record is the product source of truth; every recorded
value keeps its origin and fidelity, and the whole execution is navigable end to
end.

## 2. Module

New bounded context `internal/capture` (parallel to `internal/identity` and
`internal/device`). Layered domain → application → infrastructure, mirroring the
existing modules. Persistence follows ADR 0004: SQLite (WAL) for metadata and the
journal; large artifacts are content-addressed files with digest, size, state,
and references in SQLite.

## 3. Delivery slices (roadmap Stage 5)

1. **Capture contract primitives** — the source/fidelity model, classification
   with per-category policy, consent, retention, limits, and correlation identity.
   Pure domain. (This slice — inventory in §5.)
2. **Persistence proof** — the accepted SQLite + artifact design under concurrent
   processes, crash, migration, corruption, and disk pressure.
3. **Record intent + outcome** inside the shared application path below MCP and
   CLI, so every interface produces the same authoritative record.
4. **Idempotent synchronization** with restart and partial-upload recovery.
5. **Acknowledgement-based cleanup**, leases, budgets, and protection for active
   or un-synchronized evidence.
6. **Bounded pending agent context** and exact-versus-inferred correlation, in
   preparation for the supported host adapters.

## 4. Locked decisions (owner-approved)

- **Module** `capture`.
- **Fidelity** is a five-value model — `Exact | Summary | Inferred | Unavailable |
  Redacted` (`capture.md` §8). It labels how truthfully a captured value reflects
  the source; an inferred value is never presentable as exact, and hidden model
  reasoning is never claimed (§9).
- **Classification** is a fixed `category` enum (`Prompt, Response, Tool, Screen,
  Hierarchy, File, Log`); each category's `policy` (byte limit, retention,
  visibility, upload, redaction) lives in one canonical typed registry — no
  scattered literals. Unclassified data fails closed (SEC-011).
- **Consent, retention, and limits** are typed in Slice 1, not deferred.
- **Correlation is full end-to-end.** Beyond the six identity dimensions (session,
  turn, tool-call, MCP connection, execution, capability-call, §10) the model
  carries a trace spine so any execution is walkable first-call to last through
  all nesting: a `trace` root shared by the whole execution, a per-hop `span`, a
  `parent` link (walk up to backtrace, down to descend), a monotonic `sequence`
  for ordered traversal, and an Exact/Inferred `mark` per link. It is a
  Drizz-owned tree, not an external trace format.

## 5. Slice 1 — capture contract primitives (design inventory)

| Field | Decision |
| --- | --- |
| Capability | Define + approve the foundational capture-contract value objects: source/fidelity, classification with per-category policy, consent, and end-to-end correlation identity. No persistence. |
| Owner | `internal/capture/domain`. |
| Layer | domain only. |
| Contract | `fidelity` (Exact\|Summary\|Inferred\|Unavailable\|Redacted); `origin` (Capability = authoritative Drizz execution \| Host = agent-side observation); `category` (Prompt, Response, Tool, Screen, Hierarchy, File, Log); `policy` per category (limit, retention, visibility, upload, redaction) from one canonical registry; `consent` (approved categories); `mark` (Exact\|Inferred); `trace` + `span` identifiers; `correlation` composing the trace spine (trace, span, parent, sequence) + the six identity dimensions + mark. Immutable, self-validating, typed `Input`. |
| Dependencies | domain + stdlib only; independent of identity/device; points inward. |
| State | none — immutable value objects. |
| Failures | construction-time validation only: invalid category, non-positive limit/retention, empty correlation, a non-root span missing its parent, an inferred link missing the Inferred mark, and a category whose policy forbids upload can never be marked syncable. |
| Files | `internal/capture/domain/{fidelity,origin,category,policy,consent,mark,trace,span,correlation}.go` + tests. One concept per package, single-word, ≤500 lines. |
| Tests | per-VO validation + invariants (inferred-not-exact; no-upload-forbids-sync; retention/limit bounds; parent required off-root; backtrace walks root→leaf). |
| Verification | scoped `go test` + architecture gates; `make verify` is the merge gate. |

## 6. Verification

Platform general gate: `make verify`. Later slices add the ADR 0004 persistence
qualification (crash, corruption, disk-full, duplicate-process, migration,
cleanup, power-loss) on every supported operating system, and idempotent
synchronization recovery tests.
