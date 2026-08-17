# Reliability Standards

Status: Approved and mandatory

## Context and lifecycle

- `REL-001`: Every cancellable operation MUST propagate caller cancellation. Remote, blocking, or potentially unbounded operations MUST have a caller deadline or a documented lifecycle-owned termination condition. Long-lived subscriptions require explicit ownership and shutdown, not an artificial request deadline.
- `REL-002`: A request or job path MUST NOT replace caller context with a background context.
- `REL-003`: Every goroutine, task, worker, subprocess, timer, and subscription MUST have one owner, bounded lifetime, cancellation, cleanup, and test.
- `REL-004`: Concurrency, queue, channel, pool, cache, retry, memory, disk, and output sizes MUST be bounded and validated.
- `REL-005`: Locks MUST NOT be held across remote or slow I/O.
- `REL-006`: In-process channels coordinate work; they are not durable queues.

## State and transactions

- `REL-007`: State machines MUST define allowed transitions, terminal states, ownership, and recovery.
- `REL-008`: Transactions MUST be short and MUST NOT contain remote I/O.
- `REL-009`: Durable state and the journal or outbox describing its remote work MUST commit atomically.
- `REL-010`: A crash MUST NOT advance state beyond durable evidence.
- `REL-011`: Multi-process coordination MUST use database constraints, atomic compare-and-swap, locks, or renewable leases appropriate to the store.

## Retry and idempotency

- `REL-012`: Every retryable mutation MUST have one stable idempotency identity and exactly one logical effect under at-least-once delivery.
- `REL-013`: Retry policy MUST declare deadline, attempt or time ceiling, exponential backoff with jitter, server retry time, eligibility, and cancellation.
- `REL-014`: Validation, authorization, integrity, and incompatible-schema failures MUST NOT be blindly retried.
- `REL-015`: Lost acknowledgements MUST reconcile remote state before resending large or nontrivial bodies.
- `REL-016`: Partial success MUST remain partial and be recoverable per item.
- `REL-017`: Back-pressure MUST produce a typed, observable outcome and MUST NOT silently discard durable work.

## Storage and artifacts

- `REL-018`: Large artifacts MUST be streamed and MUST NOT be loaded entirely into memory.
- `REL-019`: Artifact writes MUST use bounded temporary storage, digest verification, atomic publication, explicit lifecycle state, and guaranteed cleanup.
- `REL-020`: Cleanup spanning stores MUST use a durable, idempotent, crash-consistent state machine. Selection and lifecycle-state changes are transactional in their authoritative store. Physical deletion is retryable and reconciled after interruption. Cleanup respects retention, leases, references, and unacknowledged recoverability.
- `REL-021`: Corruption MUST isolate the affected item and preserve bounded diagnostic evidence; it MUST NOT be silently ignored.
- `REL-022`: Disk pressure MUST follow an explicit priority and hysteresis policy. Irrecoverable loss MUST be recorded honestly.

## Migration and compatibility

- `REL-023`: Prefer additive compatibility and the sequence extend, migrate, observe, deprecate, remove.
- `REL-024`: Schema migrations MUST use expand, migrate, and contract when multiple versions may overlap.
- `REL-025`: Migrations are immutable after release and MUST be tested from an empty store and every supported prior version.
- `REL-026`: Every migration MUST be tested for interruption and supported-version upgrade. Low-space, large-data, contention, and performance qualification are required when the migration can exercise those conditions. Recovery specifies verified rollback or a documented forward-only application and data recovery procedure.
- `REL-027`: Rollback MUST preserve data even when schema rollback is unsafe.
- `REL-028`: API, wire protocol, payload schema, persistence schema, policy, provider, and artifact protocol versions evolve independently and MUST NOT be collapsed into one application version.

## Performance

- `REL-029`: Every material operation MUST declare latency class, resource bounds, concurrency, deadline, and retry budget.
- `REL-030`: Hot paths MUST be identified by measurement. Avoid unnecessary allocation, copying, repeated work, and quadratic behavior.
- `REL-031`: Optimization MUST be supported by a reproducible benchmark or profile and checked for regression.
