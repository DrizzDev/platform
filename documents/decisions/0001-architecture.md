# ADR 0001: Modular monolith with hexagonal module boundaries

- Status: Accepted
- Date: 2026-07-31
- Related: [Platform architecture](../architecture.md)

## Context

The platform must run local capabilities, support multiple interfaces, integrate
existing providers, and grow into more Drizz product areas. It is installed as
one user-facing product but must remain separable later.

## Problem

Choose an architecture that avoids both a tightly coupled binary and premature
local microservices.

## Options considered

### Horizontal layered monolith

Simple initially, but unrelated product domains accumulate in shared domain,
service, and repository directories.

### Local microservices

Strong deployment separation, but adds networking, process supervision,
version skew, installation complexity, and failure modes without an independent
scaling need.

### Modular monolith with hexagonal modules

One deployment with explicit product modules, inward dependencies, and
replaceable inbound and outbound adapters.

## Decision

Select a modular monolith. Apply hexagonal boundaries inside each real product
module and clean inward dependency direction across layers.

## Consequences

The platform keeps one installation and can execute local paths in process.
Architecture tests and application interfaces are mandatory because process
boundaries will not enforce separation. A module may be extracted later by
replacing an in-process adapter with a remote adapter.

## Validation

Prove the structure with one vertical authentication path and one vertical
device path before expanding the module set.

## Review trigger

Revisit when a module has a measured independent scaling, security, deployment,
or team ownership requirement.
