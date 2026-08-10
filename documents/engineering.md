# Engineering Guide

Status: Approved

This guide explains how the approved engineering rules are organized. The standards evolve through the reviewed change process below.

## Source of truth

The canonical rules are in [standards](standards/README.md). The root [agent instructions](../AGENTS.md) provide the short, ordered workflow that coding agents must keep in context.

| Area | Document |
| --- | --- |
| Agent workflow | [Agent protocol](standards/agent.md) |
| Architecture and boundaries | [Architecture standards](standards/architecture.md) |
| Code and configuration | [Code standards](standards/code.md) |
| Naming and structure | [Style standards](standards/style.md) |
| Go conventions | [Go standards](standards/go.md) |
| Reliability and state | [Reliability standards](standards/reliability.md) |
| Security and privacy | [Security standards](standards/security.md) |
| Logs, metrics, and traces | [Observability standards](standards/observability.md) |
| Testing | [Testing standards](standards/testing.md) |
| Delivery and dependencies | [Delivery standards](standards/delivery.md) |

The approved product authentication contract is in [Authentication and Authorization](authentication.md). The security standards govern its implementation; they do not replace that architecture. The complete Stage 3 slice inventories and gates are in the [Authentication Implementation Plan](plans/authentication.md).

## Mandatory design baseline

Every implementation applies:

- SOLID principles where applicable;
- object-oriented design through cohesive Go types and methods where it improves ownership of behavior, state, invariants, or lifecycle;
- composition over inheritance;
- dependency inversion and inward dependency direction;
- separation of concerns;
- high cohesion and low coupling;
- strongly typed contracts at every boundary;
- typed input or options structs with keyed Go literals instead of order-dependent multi-value calls;
- design patterns only for a named current problem;
- explicit evaluation of future substitution and module extraction;
- a hard maximum of 500 physical lines for every project-owned Go source file.

Extensibility does not mean creating an interface or pattern for every type. It means evaluating the need before implementation and choosing the smallest design that preserves the required evolution boundary.

## Approval and enforcement

The standards are approved repository policy.

Prose tells contributors what is expected. It does not prove compliance. Mechanical rules must become real repository checks and blocking CI gates before affected production code is accepted. Until a check exists, reviews must verify the applicable rule manually and must not describe it as machine-enforced.

## Changing a standard

A standards change must:

1. identify the problem with the current rule;
2. show the proposed wording;
3. explain its effect on existing and future code;
4. resolve conflicts with other rules;
5. receive repository-owner approval.

Historical review conversations and reviewer scores do not belong in the canonical standards.
