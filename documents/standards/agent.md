# Agent Protocol

Status: Approved and mandatory

This protocol keeps agent work ordered and recoverable during long sessions.
The root [agent instructions](../../AGENTS.md) are the active entry point.

## Before work

1. Confirm the requested outcome and non-goals.
2. Inspect the working tree and preserve unrelated changes.
3. Read the owning code and documents.
4. Read the standards relevant to the change.
5. Identify module, layer, boundaries, state, failure modes, compatibility,
   security, tests, and documentation impact.
6. Evaluate interface need, alternative implementations, implementation
   coupling, independent testability, SRP, OCP, dependency direction, and future
   module extraction.
7. Resolve material uncertainty before editing.

## During work

1. Make one coherent, reviewable change at a time.
2. Re-check the target files before each new implementation slice.
3. Keep behavior in its owning layer and module.
4. Keep every boundary strongly typed.
5. Avoid speculative abstractions and unrelated cleanup.
6. Keep every project-owned Go source file within the 500-line hard limit.
7. Run focused checks after each meaningful change.

After compaction, handoff, correction, or a material scope change, rebuild
context from the repository. Do not continue from narrative memory alone.

## Before completion

1. Re-read the request and applicable standards.
2. Inspect fresh status and the complete diff.
3. Run all available checks required by the changed area.
4. Verify failure, recovery, security, compatibility, and resource behavior
   where applicable.
5. Report the result, exact check scope, risks, and unavailable checks.

## Evidence rules

- A scoped test is not a full-suite result.
- A mock-only integration test is not proof of real integration.
- A missing checker is unavailable evidence, not a passing check.
- “Pre-existing,” “flaky,” and “unrelated” require reproduction and baseline
  evidence.
- A document or agent review is not machine enforcement.
