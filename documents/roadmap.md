# Delivery Roadmap

Status: Approved delivery baseline

This roadmap describes delivery order. It does not approve public tools,
contracts, architecture decisions, dependencies, staffing, or dates.

## 1. First release outcome

A user can install Drizz, authenticate, connect it to a supported MCP client,
use approved local device capabilities through an external agent, and see the
required execution record synchronized with Drizz.

CLI and desktop reuse the same application behavior. External agents own
planning.

## 2. Planning rules

- Build vertical outcomes, not empty layers.
- Resolve current-system facts before designing replacement contracts.
- Approve public operation names, inputs, outputs, and risk behavior before
  exposing them.
- Test local behavior on real supported devices.
- Treat installation, authentication, recovery, synchronization, and cleanup as
  product behavior.
- Do not assign dates or story points until the team, capacity, and proof results
  are known.

## 3. Delivery sequence

| Stage | Outcome | Depends on |
| --- | --- | --- |
| 1 | Reviewed product and integration foundation | None |
| 2 | Installable application skeleton with MCP and CLI entry points | Stage 1 |
| 3 | Reusable authentication and authorization context | Stage 2 |
| 4 | One real local device journey | Stage 3 |
| 5 | Durable execution records and synchronization | Stage 4 |
| 6 | Approved first capability set across MCP and CLI | Stage 5 |
| 7 | Approved deterministic authoring behavior | Stage 6 |
| 8 | Desktop integration and supported release | Stages 3–7 |
| 9 | Additional Drizz product capabilities | Stable released foundation |

## 4. Stage 1: Foundation review

### Goal

Replace assumptions with evidence from the existing Drizz systems and approve
only the decisions needed for the first implementation.

### Work

- Confirm the first user journey and non-goals.
- Audit the existing authentication and authorization implementation.
- Audit Android and iOS device integrations, supported operations, lifecycle,
  errors, concurrency, and packaging.
- Audit how the current desktop application can start, update, and communicate
  with the Drizz application.
- Classify existing authoring behavior into agent-owned planning, reusable
  deterministic behavior, cloud-owned behavior, and excluded behavior.
- Prove that the official Go MCP SDK connects over standard input/output to
  Claude, Codex, and a generic client.
- Prove that Go can supervise every required existing local provider.
- Define supported operating systems and architectures.
- Review the security and privacy boundaries.
- Approve or reject the architecture and initial stack recommendations.

### Evidence

- source-backed audit notes with repository revision;
- working MCP connection proof;
- working provider-supervision proof;
- approved first journey;
- reviewed open-decision list;
- accepted or revised foundational decision records.

### Exit

The repository owner approves the first journey and the technical owners agree
that no unresolved issue blocks the application skeleton, authentication, or
first device proof.

## 5. Stage 2: Application skeleton

### Goal

Create the smallest installable application that proves lifecycle and interface
composition without duplicating product behavior.

### Work

- Pin the approved Go toolchain and module.
- Add one composition entry point.
- Implement signal handling, cancellation, shutdown, and stable exit behavior.
- Implement typed configuration with explicit precedence.
- Implement structured diagnostics without contaminating MCP standard output.
- Add minimal MCP and CLI adapters over one application use case.
- Add build metadata and version reporting.
- Establish formatting, static analysis, unit testing, race testing, dependency
  review, secret scanning, file-size enforcement, architecture-boundary checks,
  and build checks.
- Produce unsigned development artifacts for the approved platform matrix.

### Exit

A clean checkout runs one verification command and builds the development
application. MCP and CLI invoke the same non-destructive application behavior.

## 6. Stage 3: Authentication and authorization

### Goal

Provide one authenticated context reusable by MCP, CLI, desktop, and future
approved interfaces.

### Work

- Document the current Drizz identity endpoints, token types, claims, scopes,
  organization selection, refresh, revocation, expiry, and errors.
- Define local login, logout, status, and recovery journeys.
- Implement the approved native sign-in flow.
- Store durable credentials in the operating-system credential store.
- Keep short-lived session material in memory where practical.
- Derive trusted identity and organization context outside interface input.
- Revalidate cloud authorization at the cloud boundary.
- Define offline behavior explicitly.
- Handle expired, revoked, missing, corrupted, and organization-changed
  credentials.
- Test browser interruption, callback failure, restart, concurrent login,
  logout, and credential-store denial.
- Verify the same authenticated context through MCP, CLI, and the desktop test
  boundary.

### Exit

One user can sign in and out safely on every supported operating system.
Interface input cannot override trusted identity or organization context.

## 7. Stage 4: First local device journey

### Goal

Execute one approved observation and one approved action on real Android and iOS
targets through the shared application core.

### Work

- Define device identity and session ownership.
- Define discovery, connect, disconnect, and reconnect behavior.
- Normalize Android and iOS observations without hiding platform differences.
- Select one low-risk observation and one representative action.
- Define input, output, cancellation, timeout, and error behavior.
- Build the application-owned device interface.
- Integrate the existing Android and iOS implementations through adapters.
- Bound subprocesses, output, memory, concurrency, and device ownership.
- Capture only the evidence required by the product journey.
- Test disconnect, stale device, locked device, application crash, provider
  crash, cancellation, and restart on real targets.

### Exit

The approved journey works locally on real supported Android and iOS targets
with deterministic errors and recovery.

## 8. Stage 5: Records and synchronization

### Goal

Preserve the required execution history and synchronize it without adding a
cloud round trip to each local device step.

### Work

- Define which execution facts, tool calls, assets, failures, and agent-provided
  context Drizz is allowed and required to store.
- Define data classification, redaction, consent, and retention.
- Define the cloud synchronization contract and idempotency behavior.
- Select local metadata and artifact storage after a real persistence proof.
- Commit execution state and synchronization intent safely.
- Resume after process termination, machine restart, sleep, and network loss.
- Handle duplicate requests, partial upload, server rejection, and version
  mismatch.
- Verify cloud acknowledgement before cleanup eligibility.
- Protect active, referenced, failed, and unsynchronized artifacts.
- Enforce disk budgets and observable cleanup.
- Test corruption, disk full, clock changes, competing processes, and upgrade.

### Exit

The first device journey appears in the required Drizz surface after
synchronization. Repeated delivery does not duplicate cloud effects, and local
storage remains bounded without deleting required evidence.

## 9. Stage 6: First MCP and CLI capability set

### Goal

Expose the approved first-release capabilities consistently through standard
MCP clients and CLI automation.

### Work

- Product-review every operation name, description, input, output, limit, risk,
  and destructive effect.
- Define one transport-neutral capability catalog.
- Map approved capabilities to MCP tools.
- Map the same capabilities to CLI commands.
- Provide stable machine-readable CLI output and exit behavior.
- Define confirmation and explicit enablement for destructive operations.
- Add client configuration helpers that never overwrite user configuration
  without consent.
- Validate install, configure, authenticate, invoke, cancel, restart, upgrade,
  and uninstall with Claude.
- Repeat the supported journey with Codex and a generic MCP client.
- Verify that client-specific behavior does not enter the application core.

### Exit

Every approved first-release operation has equivalent authorization, execution,
recording, cancellation, and failure behavior through MCP and CLI.

## 10. Stage 7: Deterministic authoring

### Goal

Reuse only approved deterministic authoring behavior while the external agent
continues to plan.

### Work

- Identify the minimum authoring inputs and outputs.
- Separate planning, hidden reasoning, and model-specific control from reusable
  deterministic behavior.
- Review source distribution, intellectual property, assets, dependencies, and
  packaging.
- Decide whether approved behavior is ported, invoked through a supervised
  local process, or replaced through a smaller contract.
- Define author, validate, and save journeys without exposing internal
  implementation vocabulary.
- Define session memory, recovery, loop protection, and cleanup only where
  deterministic product behavior requires them.
- Verify output through multiple external agents.
- Synchronize the authored result into the approved Drizz surface.

### Exit

Approved authoring works locally without Drizz owning planning and without
shipping unapproved proprietary implementation.

## 11. Stage 8: Desktop and release

### Goal

Deliver one supported Drizz implementation through standalone and desktop
journeys.

### Work

- Approve how the desktop locates, starts, stops, updates, and verifies Drizz.
- Prevent desktop and standalone processes from corrupting shared state.
- Define version compatibility between desktop and Drizz.
- Build signed artifacts for the supported platform matrix.
- Add notarization where required.
- Publish approved package-manager channels.
- Define update, rollback, downgrade, and data-compatibility behavior.
- Verify install, upgrade, downgrade, uninstall, and retained-user-data policy.
- Complete privacy-safe diagnostics and support guidance.
- Run the full first-release journey on every supported client and platform.

### Exit

A user installs Drizz through an approved channel, connects a supported client,
authenticates, completes the first journey, sees synchronized history, upgrades,
and uninstalls without credential leakage or data corruption.

## 12. Stage 9: Additional capabilities

Test plans, debugging, application management, reports, and other product areas
are added as separate vertical outcomes.

Each capability must:

- use the same authentication and application boundaries;
- enforce organization and resource authorization in Drizz services;
- reuse execution, artifact, synchronization, and cleanup behavior where its
  semantics match;
- define its own approved product contract;
- avoid exposing another module's storage or internal provider vocabulary;
- work through every interface that its product requirement names.

## 13. Open planning inputs

Dates, staffing, parallel delivery, and estimates remain open until:

- the responsible team is known;
- the first journey is approved;
- the authentication, device, desktop, and provider audits are complete;
- supported platforms are chosen;
- proof results expose the real integration cost.

Any schedule created before those inputs is a guess and must be labeled as one.
