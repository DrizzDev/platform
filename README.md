# Drizz Platform

Drizz Platform is the local execution and integration foundation that allows AI
agents, developer tools, and Drizz applications to use the same Drizz
capabilities.

This repository currently contains the approved design and engineering
baseline. Product implementation has not started. Proof-gated integrations and
dependencies still require their documented validation before implementation.

## Intended ownership

- One installable Drizz application.
- Shared application behavior for MCP, CLI, desktop, and future approved
  interfaces.
- Fast local device interaction.
- Durable local execution records and synchronization with Drizz services.
- Explicit integration boundaries for existing Drizz systems.

External agents own planning. Drizz executes approved capabilities, enforces its
boundaries, and records the resulting work.

## Documents

1. [Architecture](documents/architecture.md)
2. [Technology stack](documents/stack.md)
3. [Engineering guide](documents/engineering.md)
4. [Delivery roadmap](documents/roadmap.md)
5. [Decision records](documents/decisions/README.md)
6. [Research](documents/research.md)
7. [Agent instructions](AGENTS.md)

## Current status

| Area | Status |
| --- | --- |
| Product implementation | Not started |
| Architecture | Approved baseline |
| Technology stack | Approved recommendation baseline |
| Engineering standards | Approved and mandatory |
| Delivery roadmap | Approved baseline |
| Decision records | Accepted |

Implementation should begin only after the repository owner has reviewed the
foundation documents relevant to the first deliverable.
