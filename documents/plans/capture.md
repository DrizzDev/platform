# Stage 5 — Capture, Records, and Synchronization Implementation Plan

Status: In progress. Companion architecture: [Agent integration and execution capture](../capture.md) and [ADR 0007](../decisions/0007-capture.md); persistence per [ADR 0004](../decisions/0004-persistence.md). Delivery order per the [roadmap](../roadmap.md) Stage 5.

## 1. Goal

Preserve the required capability and agent execution history and synchronize it to Drizz without adding a cloud round trip to each local device step. The ordered, durable local execution record is the product source of truth; every recorded value keeps its origin and fidelity, and the whole execution is navigable end to end.

## 2. Module

New bounded context `internal/capture` (parallel to `internal/identity` and `internal/device`). Layered domain → application → infrastructure, mirroring the existing modules. Persistence follows ADR 0004: SQLite (WAL) for metadata and the journal; large artifacts are content-addressed files with digest, size, state, and references in SQLite.

## 3. Delivery slices (roadmap Stage 5)

1. **Capture contract primitives** — the source/fidelity model, classification with per-category policy, consent, retention, limits, and correlation identity. Pure domain.
2. **Persistence proof** — the accepted SQLite + artifact design under concurrent processes, crash, migration, corruption, and disk pressure.
3. **Record intent + outcome** inside the shared application path below MCP and CLI, so every interface produces the same authoritative record.
4. **Idempotent synchronization** with restart and partial-upload recovery.
5. **Acknowledgement-based cleanup**, leases, budgets, and protection for active or un-synchronized evidence.
6. **Bounded pending agent context** and exact-versus-inferred correlation, in preparation for the supported host adapters.

## 4. Locked decisions (owner-approved)

- **Module** `capture`.
- **Fidelity** is a five-value model — `Exact | Summary | Inferred | Unavailable | Redacted`. It labels how truthfully a captured value reflects the source; an inferred value is never presentable as exact, and hidden model reasoning is never claimed.
- **Classification** is a fixed `category` enum (`Prompt, Response, Tool, Screen, Hierarchy, File, Log`); each category's `policy` (byte limit, retention, visibility, upload, redaction) lives in one canonical typed registry — no scattered literals. Unclassified data fails closed.
- **Consent, retention, and limits** are typed in Slice 1, not deferred.
- **Correlation is full end-to-end.** Beyond the six identity dimensions (session, turn, tool-call, MCP connection, execution, capability-call) the model carries a trace spine so any execution is walkable first-call to last through all nesting: a `trace` root shared by the whole execution, a per-hop `span`, a `parent` link (walk up to backtrace, down to descend), a monotonic `sequence` for ordered traversal, and an Exact/Inferred `mark` per link. It is a Drizz-owned tree, not an external trace format.

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

## 6. Slice 2b — artifact storage and deferrals

Slice 2b adds `internal/capture/infrastructure/artifact`: a content-addressed file store. `digest.Digest` (a validated SHA-256) names each object; `Put` streams to a bounded temp, hashes, and atomically publishes by rename; `Get` re-verifies the digest on read so silent corruption fails rather than returns bad bytes; `Sweep` removes temps orphaned by a crashed write.

- **Artifact size ceiling.** A single canonical default, `ceiling = 128 MiB`, large enough that no legitimate artifact (screenshots are bounded at 32 MiB) is rejected — set conservatively so we never lose data. It is **configurable in one place** via `Options.Ceiling`; the default applies only when unset.
- **DEFERRED TO SLICE 5 — do not miss.** Reference-aware artifact retention is NOT in 2b. Deleting an artifact once no journal entry references its digest, with leases, disk budgets, and protection for active or un-synchronized evidence, belongs to Slice 5 (acknowledgement-based cleanup). Slice 2b cleans up **orphaned temporary files only**; a published object is never deleted by 2b.

## 7. Store resilience — deletion, corruption, concurrency

The data directory is user-owned; a user or another agent can delete or damage it. We make the store resilient rather than trying to prevent deletion (which is not possible for user-owned data and is not the industry norm).

- **Corruption → quarantine and rebuild (done, capture sqlite).** On open, a corrupt or not-a-database file is detected by its SQLite result code, renamed aside as `<db>.corrupt-<nanos>` (with any `-wal`/`-shm`), and a fresh database is created — so a damaged file never blocks startup, and the bad file is preserved, never silently discarded. Identity's store should get the same (tracked follow-up; not on this branch).
- **Deletion → self-heal.** A missing directory or database is recreated on start; a missing artifact returns `Absent`; atomic temp+rename means no half-written file. Deleting the data dir loses history but never crashes or corrupts, and the binary itself lives elsewhere and is untouched.
- **Capture is observational (enforced in Slice 3).** A device or agent operation MUST proceed and return even if the journal or artifact write fails; the failure is logged, never propagated. This is the primary protection: even if the data dir is deleted mid-run, the device tool keeps working — only recording degrades.
- **Concurrency: no single-instance lock.** MCP and CLI share these stores at the same time by design. SQLite WAL + `busy_timeout` + one-writer, and content-addressed atomic-rename artifacts, already make multi-process access safe. A store-level single-instance lock is rejected — it would break legitimate coexistence to prevent a problem that cannot occur.

## 8. Slice 5 — acknowledgement-based cleanup, retention, leases, budgets

A synchronized entry (`Synced`, set when the courier's upload is acknowledged) is the only entry a reclaim pass ever removes; un-synchronized evidence is the hard floor and is never deleted.

- **`internal/capture/application/janitor`.** `Janitor.Run` is drained by the host loop (no timer). Each pass: age out `Synced` records past their category's retention (`catalogue`, measured off `stamped`); then sweep artifacts no remaining entry references, sparing any within a grace window that closes the `Put`-then-`Append` race. Selection (row delete) is one transaction; physical file deletion is idempotent and re-derived each pass, so an interrupted pass resumes. No cross-process lock — deletes are idempotent and reclaim is conservative.
- **Budget with hysteresis.** `Ceiling` (high-water) and `Relief` (low-water). Over the ceiling, reclaim escalates to `Synced` records even within retention — safe, the cloud holds them. When only un-synchronized evidence remains, the artifact store's write gate is raised (`Restrict`; `Put` returns a typed `Saturated`) and the user is alerted through a `notifier` port plus a pressure metric; new writes degrade rather than deleting evidence. The gate lifts (`Admit`) once reclaim frees below the relief mark.
- **Active-execution lease.** The recorder refreshes a durable lease on its trace each `Record` (no timer; swallowed on failure). A reclaim pass skips a live-leased trace; a crashed holder's lease lapses and protection lifts.
- **Store surface added.** journal: `Settled`, `Discard`, `Digests`, `Lease`, `Leases`. artifact: `Prune(keep predicate)`, `Footprint`, `Restrict`, `Admit`, and `Put` honouring the gate. The `janitor` ports pass only domain and standard types (a `keep func(digest, time) bool` predicate keeps grace and reference policy in the application while the walk stays in the store), so no port points outward. `lease` table folded into the single `0001` baseline.
- **Deferred wiring.** The composition root must adapt the identity session flow to the courier's credential provider (Slice 4b) and construct the notifier and the janitor loop. No `cmd` wiring exists yet.

## 9. Slice 6 — correlation matching and the bounded pending window

Host hooks and MCP see different parts of one execution; the capture layer associates a host observation with the Drizz capability call it belongs to, then activates it into that execution. Unmatched host data expires locally and is never synchronized.

- **Domain `affinity`.** `bearings.Shares` (any of the six dimensions present and equal), and `Signal.Match` — exact on a shared identifier, otherwise inferred when an earlier observation is within the call's window and agrees on normalized input or process, otherwise no match. Pure; time is passed in, never read. `identifier.Same` and `digest.Same` express the value-object equality the match needs. The mark is carried honestly: inferred is never relabelled exact.
- **Ephemeral `pending` table** (folded into `0001`, never synchronized). Bounded by a time-to-live and a capacity. Its columns derive from the same generic `layout` the journal uses, so its insert, select, and scan cannot drift.
- **Application `lobby`.** `Observe` holds an observation and expires the stale; `Activate` matches a `Call` against the live window, returns the claimed observations tagged exact or inferred, and evicts them so each activates at most once. The lobby matches and evicts; the caller records the claims into the activated execution, so the lobby does not depend on recording.
- **Recording** carries the mark: `Note` gains a mark that defaults to exact, so a direct capability record is unchanged and an inferred match records as inferred.
- **Host adapters deferred (Stage 6).** No host adapter exists yet, so real per-adapter and real-client fixtures ride with each integration. This slice ships the host-neutral engine with synthetic fixtures for exact, inferred, missing, duplicate, delayed, and out-of-order events.

Stage 5 is closed end to end: `tests/capture` drives one observation through the whole stack — held in the window, claimed by a call, recorded, synchronized to a fake cloud, then reclaimed once acknowledged. The composition root wiring (the identity-backed credential provider, the notifier, and the observe/activate/loop schedule) lands with the device tool in Stage 6.

## 10. Verification

Platform general gate: `make verify`. Later slices add the ADR 0004 persistence qualification (crash, corruption, disk-full, duplicate-process, migration, cleanup, power-loss) on every supported operating system, and idempotent synchronization recovery tests.
