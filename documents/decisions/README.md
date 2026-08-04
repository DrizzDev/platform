# Architecture Decision Records

Status: Approved

Architecture Decision Records, or ADRs, preserve why the platform is built a
certain way. They prevent future contributors from treating an old choice as an
accident or repeating the same research.

This directory exists because the project requires an audit trail for material
technical decisions. It does not store meeting notes, routine implementation
choices, task history, or agent reasoning.

## Statuses

- `proposed`: written for review; implementation may not depend on it;
- `accepted`: approved and active;
- `superseded`: replaced by another ADR;
- `rejected`: considered and not selected;
- `deprecated`: still present but scheduled for removal.

Only the repository owner or an explicitly delegated technical owner changes a
proposal to accepted.

## When an ADR is required

Create an ADR only for a durable decision that would be expensive or risky to
reverse, including:

- architecture or module boundaries;
- language or framework;
- persistence and migrations;
- process and deployment topology;
- authentication and authorization;
- public contract and compatibility policy;
- a foundational third-party dependency;
- a deliberate exception to the engineering guide.

Do not create an ADR for a local refactor, a self-explanatory bug fix, a normal
library addition, a task plan, or a reversible implementation detail.

## Workflow

1. Copy [template.md](template.md).
2. Use the next four-digit number and a short lowercase title.
3. Mark it `proposed`.
4. Link evidence, experiments, and affected documents.
5. Obtain review from the actual owners of the affected boundaries.
6. Record the decision date and mark it `accepted`, `rejected`, or leave it
   `proposed`.
7. Update this index and the affected architecture documents.
8. If the decision changes later, add a new ADR and mark the old one
   `superseded`; do not rewrite history.

## Index

| ADR | Decision | Status |
| --- | --- | --- |
| [0001](0001-architecture.md) | Modular monolith with hexagonal module boundaries | Accepted |
| [0002](0002-language.md) | Go as the platform host | Accepted |
| [0003](0003-local.md) | Local execution with cloud authority | Accepted |
| [0004](0004-persistence.md) | SQLite journal and file artifact storage | Accepted |
| [0005](0005-interfaces.md) | One capability core with MCP, CLI, and desktop adapters | Accepted |
| [0006](0006-authentication.md) | OAuth flows by client type with cloud authorization | Accepted |
| [0007](0007-capture.md) | Native host integration with a Drizz execution record | Accepted |
