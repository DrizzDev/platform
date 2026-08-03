# Agent Instructions

These instructions apply to every coding agent working in this repository.
The detailed engineering rules are in
[documents/standards](documents/standards/README.md).

## 1. Understand the work

Before changing tracked files:

1. read the user request and identify the expected outcome;
2. inspect the current working tree;
3. read the files that own the affected behavior;
4. read the applicable standards and accepted decision records;
5. identify affected boundaries, state, failures, compatibility, security, and
   tests;
6. stop and ask when a material product or architecture decision is unresolved.

Never treat remembered context as repository evidence. Re-read relevant source
after context compaction, handoff, a user correction, or a material scope
change.

## 2. Design before editing

- Complete the design inventory required by
  [the agent protocol](documents/standards/agent.md) before each implementation
  slice. Do not begin a slice with an unresolved owner, boundary, dependency,
  state, failure, file, test, or verification decision.
- Apply SOLID, separation of concerns, high cohesion, low coupling, dependency
  inversion, and composition over inheritance wherever applicable.
- Use object-oriented design through cohesive Go types and methods where
  behavior, state, invariants, or lifecycle belong together.
- Use a design pattern only when it solves a named current problem.
- Preserve the repository's layered dependency direction.
- Give every module and type one clear responsibility.
- Use strict single-word names. Represent multi-word concepts through meaningful
  nesting.
- Follow normal Go visibility naming for identifiers.
- Use `UPPER_SNAKE_CASE` only for internal semantic string values.
- Preserve values defined by an external protocol or public contract.
- Keep public surfaces minimal and validate input at boundaries.
- Go has no keyword arguments. Use one typed input or options struct with keyed
  fields whenever a project-owned call needs multiple values; keep
  `context.Context` first when applicable.
- Do not introduce a dependency, abstraction, interface, or framework without a
  demonstrated requirement.
- Record a material, durable architecture decision only after its owner has
  reviewed it.
- Treat each implementation slice as a fresh repository review. Re-read the
  applicable standards and owning code instead of relying on conversation
  memory.

Before choosing a concrete design, answer:

1. Does this component need an interface or abstraction?
2. Is another implementation or provider reasonably expected?
3. Is a specific implementation leaking across a boundary?
4. Can the component be tested independently?
5. Does it violate SRP, OCP, or dependency inversion?
6. Do dependencies point toward the owning policy?
7. Could the module be extracted later without rewriting its consumers?

The answers may justify a concrete type. They may not be skipped.

## 3. Make the smallest complete change

- Stay within the requested scope.
- Preserve unrelated and pre-existing work.
- Keep domain policy out of transports and infrastructure.
- Keep constants, contracts, models, errors, application ports, use cases,
  adapters, transport handlers, and validators in their owning layer and
  responsibility.
- Use explicit strongly typed contracts at every boundary. Generic maps and
  implementation-specific objects must not cross layers.
- Keep secrets, authorization context, and customer data inside their approved
  boundaries.
- Do not add speculative features or future-facing scaffolding.
- No project-owned Go source file, including tests and generated project code,
  may exceed 500 physical lines.
- Do not write unnecessary comments. A justified technical comment should
  normally be one line and must not exceed three short lines.

## 4. Verify with evidence

1. run the narrowest relevant checks while editing;
2. run every available repository check required by the changed area;
3. test real integrations when the result depends on real provider behavior;
4. inspect the fresh complete diff and working-tree status;
5. review the diff against the completed design inventory;
6. distinguish scoped checks from repository-wide checks.

Never label a failure as pre-existing, flaky, or unrelated without reproduction
and baseline evidence. A required check that does not yet exist must be reported
as unavailable, not simulated through prose or configuration.

## 5. Report honestly

Report:

- the outcome;
- changed files;
- decisions made;
- checks run and their exact scope;
- unresolved risks or blocked checks.

Do not claim approval, completeness, compliance, or production readiness that
has not been demonstrated.
