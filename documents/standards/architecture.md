# Architecture Standards

Status: Approved and mandatory

## Design principles

- `ARC-001`: Architecture and implementation MUST apply SOLID principles where
  they fit the language and problem: single responsibility, open/closed,
  substitutability, interface segregation, and dependency inversion.
- `ARC-002`: Object-oriented design MUST be used where behavior, state,
  invariants, or lifecycle form a cohesive responsibility. Go applies these
  principles through types, methods, interfaces, and composition rather than
  class hierarchies.
- `ARC-003`: Composition MUST be preferred over inheritance. Inheritance or
  embedding used for behavioral reuse requires a demonstrated semantic reason
  and MUST preserve substitutability.
- `ARC-004`: A design pattern MAY be introduced only when it solves a named,
  current design problem. Pattern-driven scaffolding and abstractions added
  solely for possible future use are prohibited.
- `ARC-005`: Every component MUST have high cohesion: its behavior, state, and
  dependencies serve one clearly stated responsibility.
- `ARC-006`: Coupling between modules and layers MUST be explicit, narrow, and
  directed through owned contracts.
- `ARC-007`: Separation of concerns is mandatory. Product policy, orchestration,
  transport mapping, persistence, provider integration, and process lifecycle
  MUST NOT be mixed.
- `ARC-008`: Every material design decision MUST be evaluated for
  extensibility, maintainability, testability, and long-term evolution.

## Dependency direction

- `ARC-009`: Dependencies MUST point inward. Transport and infrastructure
  adapters depend on application-owned ports; application depends on domain.
- `ARC-010`: Domain contains invariants, entities, value objects, policies, and
  deterministic decisions only.
- `ARC-011`: Domain MUST NOT perform I/O, read environment or current time,
  access runtime or provider frameworks, log, emit telemetry, or use global
  mutable state.
- `ARC-012`: Application owns use cases, orchestration, transaction intent,
  cancellation, and ports. It MUST NOT interpret provider APIs or invent domain
  policy.
- `ARC-013`: Adapters validate and translate external representations.
  Infrastructure implements persistence, network, file, process, queue,
  credential, and provider ports.
- `ARC-014`: Composition and process lifecycle remain in the composition root.
  Business behavior MUST NOT be placed there.
- `ARC-015`: No module may bypass another module's application interface,
  import its internal implementation, or access its storage directly.
- `ARC-016`: Business behavior MUST be implemented once and reused by MCP, CLI,
  desktop, jobs, and future approved adapters.

## Extensibility

- `ARC-017`: Before implementation, the owner MUST evaluate whether a component
  needs an interface, abstraction, alternate implementation, replacement
  boundary, or independent lifecycle.
- `ARC-018`: The evaluation MUST ask whether the component is coupled to a
  provider, can be tested independently, violates SRP or OCP, points dependencies
  inward, and can be extracted later without rewriting its consumers.
- `ARC-019`: An interface MUST be owned by its consumer and introduced when a
  real boundary, multiple behaviors, external provider, test seam required by
  production architecture, or accepted replacement requirement exists.
- `ARC-020`: Test convenience alone is insufficient reason for an interface.
  Tests may fake an interface required by production architecture.
- `ARC-021`: Concrete providers MUST be replaceable without changing domain or
  application behavior.
- `ARC-022`: Shared code requires a precise owner, stable meaning, and at least
  two demonstrated consumers. Generic shared dumping grounds are prohibited.
- `ARC-023`: A service extraction MUST solve a measured deployment, security,
  scaling, reliability, ownership, or isolation requirement.
- `ARC-024`: A module intended for possible extraction MUST expose an explicit
  application contract and MUST NOT share internal storage or provider types
  with another module.

## Contracts

- `ARC-025`: Every layer and module boundary MUST use explicit, strongly typed
  input, output, error, and state contracts.
- `ARC-026`: Generic maps, untyped payloads, loosely typed containers, and
  implementation-specific objects MUST NOT cross a boundary.
- `ARC-027`: Raw external data is allowed only at the outermost decoding
  boundary and MUST be validated and converted immediately.
- `ARC-028`: Provider, transport, generated, persistence, and framework types
  MUST NOT leak into domain, application, or public product contracts.
- `ARC-029`: A boundary contract MUST define ownership, validation,
  cancellation, failure, compatibility, and material side effects.

## Structure and responsibility

- `ARC-030`: Every module, package, directory, file, type, and public interface
  MUST have one clear responsibility.
- `ARC-031`: Modules MUST use strict layered nesting. Related sub-concepts
  belong below their owner instead of beside it in a flat directory.
- `ARC-032`: Constants and enumeration definitions remain with the behavior
  that owns their meaning. Boundary DTOs remain in their boundary. Domain
  models and errors remain in domain. Application ports, repositories, and use
  cases remain in application. Provider implementations remain in adapters or
  infrastructure. Transport handlers and transport validation remain in the
  owning transport.
- `ARC-033`: A file or directory MUST NOT mix constants, contracts, domain
  policy, orchestration, persistence, provider integration, transport handling,
  and generic utility behavior merely for convenience.
- `ARC-034`: Cross-cutting technical primitives require a precise technical
  owner. Generic `common`, `helper`, `service`, `types`, or `utils` packages are
  prohibited.

## Complexity

- `ARC-035`: A type, method, function, or file SHOULD have one cohesive reason
  to change.
- `ARC-036`: A name requiring “and” to describe its responsibility is evidence
  that the responsibility SHOULD be split.
- `ARC-037`: Size limits are architecture guards, not permission to accumulate
  unrelated behavior up to the limit.
- `ARC-038`: A corrective refactor SHOULD precede new behavior when the touched
  boundary already violates these standards. The refactor MUST remain scoped
  and behavior-preserving.

## Go ownership

- `ARC-039`: Go package co-location does not waive separation of concerns. A
  source file contains one primary concept and only implementation inseparable
  from that concept. Contracts, errors, configuration, lifecycle,
  orchestration, persistence, provider integration, transport behavior,
  validation, and observability MUST use separate files or narrower packages
  when they can change independently.
- `ARC-040`: Structure is owner-first and layer-second. Global directories for
  constants, schemas, models, repositories, or services are prohibited.
  Constants remain with their semantic owner, boundary models with their
  boundary, persistence ports with their use cases, and implementations with
  their adapter or infrastructure owner.
