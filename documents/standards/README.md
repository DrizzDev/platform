# Engineering Standards

Status: Approved and mandatory

These documents consolidate relevant engineering rules for this Go repository. They are deliberately separated by concern so contributors and coding agents can load only what a change requires.

They are the approved engineering baseline. They may be refined through the reviewed standards-change process as the project evolves.

## Reading order

For every tracked change:

1. read the [agent protocol](agent.md);
2. read [architecture](architecture.md), [code](code.md), [style](style.md), [Go](go.md), [testing](testing.md), and [delivery](delivery.md);
3. add [reliability](reliability.md) for I/O, state, concurrency, retry, synchronization, cleanup, or migration;
4. add [security](security.md) for authentication, authorization, external input, files, processes, credentials, or customer data;
5. add [observability](observability.md) for logs, metrics, traces, diagnostics, or telemetry.

## Rule strength

- `MUST` means the change cannot be accepted without correcting the violation or receiving an explicit owner-approved exception.
- `SHOULD` means deviation requires a written technical reason.
- `MAY` describes an allowed option.

## Authority

- The repository owner approves product scope and these standards.
- Approved product contracts define product behavior.
- Accepted decision records define durable architecture choices.
- These standards define implementation and delivery expectations.
- A roadmap defines sequence only; it does not approve product behavior or architecture.

When two authorities conflict, stop and resolve the conflict explicitly. Security, privacy, authorization, data integrity, compatibility, and destructive actions cannot be waived implicitly.

## Enforcement

Do not represent prose, draft configuration, or an agent review as automated enforcement. A mechanical rule is enforced only when an implemented repository command checks it and CI blocks violations.
