# ADR 0002: Go as the platform host

- Status: Accepted
- Date: 2026-07-31
- Related: [Technology stack](../stack.md)

## Context

The public host must be easy to install, start quickly, manage local processes and devices, support CLI and MCP, integrate with a TypeScript desktop, and call any approved local authoring provider.

## Options considered

### TypeScript

Strong MCP and desktop ecosystem, but requires a runtime or heavier bundling.

### Python

Strong compatibility with existing Fathom code, but distributing private source and a reliable cross-platform environment is harder.

### Go

Produces a native binary, has strong process and concurrency primitives, a mature CLI ecosystem, and an official Tier 1 MCP SDK.

## Decision

Use Go 1.26 for the platform host. Keep existing systems behind narrow, supervised provider boundaries when they must run locally. Do not translate or ship those systems without a measured and approved reason.

Supported agent hooks invoke an external executable and exchange structured data; they do not require the hook receiver to use the host's implementation language. The installed Go application therefore implements Claude and Codex hook receivers without a Python runtime.

## Consequences

The main host is easy to package and can reuse one application core. Provider communication and lifecycle need explicit design. Go cannot make shipped source private; deterministic authoring therefore requires a separate product, security, and distribution decision.

## Validation

Build foundation proofs for MCP startup, CLI startup, desktop supervision, device provider invocation, and one selected authoring operation.

## Review trigger

Revisit only if a mandatory protocol, platform API, or distribution requirement cannot be supported safely from Go or a bounded provider.
