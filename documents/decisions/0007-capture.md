# ADR 0007: Native host integration with a Drizz execution record

- Status: Accepted
- Date: 2026-08-03
- Related: [Agent integration and execution capture](../capture.md),
  [Platform architecture](../architecture.md)

## Context

The primary user runs Claude, Codex, or another external agent application. The
agent owns conversation and planning while Drizz exposes local capabilities and
must retain sufficient execution evidence for product history, debugging,
evaluation, and future datasets.

External projects demonstrate SDK instrumentation and coding-harness tracing,
but their language, storage, export destination, and data model do not own the
Drizz product boundary.

## Options considered

### MCP capture only

This reliably records Drizz tool execution but cannot provide host session,
prompt, response, other tool, or provider-exposed reasoning context.

### OpenInference instrumentation as the primary integration

This is useful when a developer controls the process that calls a model or
agent SDK. It does not instrument an arbitrary external CLI or desktop
application and would make an observability model define Drizz product data.

### Python coding-harness dependency

Existing projects provide useful host mappings, but shipping them would add a
Python environment, a second installer and updater, external-backend
assumptions, and state that Drizz does not own.

### Native host integration with a Drizz-owned record

Official MCP, plugin, hook, transcript, and structured-event surfaces provide
the available host data. Thin Go adapters translate that data into a typed
Drizz record shared with authoritative capability instrumentation.

## Decision

Use official host integrations for primary external agent applications. Package
MCP configuration and lifecycle hooks in supported plugins, and implement their
receivers as replaceable Go adapters over a Drizz-owned capture contract.

Keep the ordered local Drizz execution record as the canonical product and
future dataset representation. Record only provider-exposed reasoning with its
actual representation and mark missing or inferred information honestly.

Do not add OpenInference to the first-release runtime or roadmap. A future SDK
integration may accept or export OpenInference through a separate adapter after
approval. Do not install `coding-harness-tracing`, Neatlogs, or a Python sidecar
on customer machines; use them only as reviewed research evidence.

## Consequences

The primary user receives a native installation without a Python runtime or
observability account. Drizz controls consent, durability, artifact identity,
authorization, synchronization, retention, and dataset evolution. Each agent
application still requires a small versioned adapter and real compatibility
tests, and surfaces without lifecycle hooks provide only partial context.

OpenInference can be added or replaced later without changing domain or
application contracts. Drizz must maintain an explicit capability matrix and
must not claim access to private reasoning that a provider does not expose.

## Validation

Complete the real-client verification matrix in the capture document for Claude
and Codex. Prove plugin installation, hook capture, MCP execution, correlation,
privacy scoping, durable recording, synchronization, update, and removal on each
supported surface and operating system.

## Review trigger

Revisit when a primary supported agent removes or materially changes its hook
or plugin surface, when a required host can expose data only through a standard
trace protocol, or when an approved custom-SDK journey makes OpenInference a
current product requirement.
