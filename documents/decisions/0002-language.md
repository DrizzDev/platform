# ADR 0002: Go as the platform host

- Status: Accepted
- Date: 2026-07-31
- Related: [Technology stack](../stack.md)

## Context

The public host must be easy to install, start quickly, manage local processes
and devices, support CLI and MCP, and integrate with a TypeScript desktop and
selected Python authoring code.

## Options considered

### TypeScript

Strong MCP and desktop ecosystem, but requires a runtime or heavier bundling and
does not improve the selected Python integration.

### Python

Strong Fathom compatibility, but distributing the private implementation and a
reliable cross-platform environment is harder.

### Go

Produces a native binary, has strong process and concurrency primitives, a
mature CLI ecosystem, and an official Tier 1 MCP SDK.

## Decision

Use Go 1.26 for the platform host. Keep existing TypeScript and Python systems
behind narrow, supervised provider boundaries. Do not translate those systems
without a measured reason.

## Consequences

The main host is easy to package and can reuse one application core. Provider
IPC and lifecycle need explicit design. Go cannot make shipped Python source
private; selected Fathom distribution therefore requires a separate product and
security decision.

## Validation

Build foundation proofs for MCP startup, CLI startup, desktop supervision,
device provider invocation, and one selected authoring operation.

## Review trigger

Revisit only if a mandatory protocol, platform API, or distribution requirement
cannot be supported safely from Go or a bounded provider.
