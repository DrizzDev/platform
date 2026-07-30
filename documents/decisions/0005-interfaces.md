# ADR 0005: One capability core with multiple inbound adapters

- Status: Accepted
- Date: 2026-07-31
- Related: [Platform architecture](../architecture.md)

## Context

Users may invoke Drizz from Claude, Codex, another MCP client, shell automation,
an IDE, or the Drizz desktop application. Building behavior inside each
interface would create drift.

## Decision

Define product capabilities as application use cases. MCP, CLI, desktop IPC,
and future HTTP handlers are thin inbound adapters over those use cases. Use the
official Go MCP SDK. Run local MCP over standard input/output and let the MCP
client start and stop `drizz mcp`.

## Consequences

Behavior, authorization, recording, and provider selection remain consistent.
Each transport still needs contract mapping and conformance tests. A persistent
desktop host is added only if the desktop journey requires shared live state.

## Validation

Invoke the same device proof through CLI, MCP, and a desktop test harness and
verify equivalent application outcomes and records.

## Review trigger

Revisit if a transport has product semantics that cannot be represented by a
shared application use case.
