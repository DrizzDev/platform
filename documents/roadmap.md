# Delivery Roadmap

Status: Approved delivery baseline

Current position: the foundation review and application skeleton are complete. Authentication and authorization are the next delivery stage.

This roadmap describes delivery order. It does not approve public tools, contracts, architecture decisions, dependencies, staffing, or dates.

## 1. First release outcome

A user can install Drizz, authenticate, connect it to a supported MCP client, use approved local device capabilities through an external agent, and see the required execution record synchronized with Drizz.

CLI and desktop reuse the same application behavior. External agents own planning.

## 2. Planning rules

- Build vertical outcomes, not empty layers.
- Resolve current-system facts before designing replacement contracts.
- Approve public operation names, inputs, outputs, and risk behavior before exposing them.
- Test local behavior on real supported devices.
- Treat installation, authentication, recovery, synchronization, and cleanup as product behavior.
- Do not assign dates or story points until the team, capacity, and proof results are known.

## 3. Delivery sequence

| Stage | Outcome | Depends on |
| --- | --- | --- |
| 1 | Reviewed product and integration foundation | None |
| 2 | Installable application skeleton with MCP and CLI entry points | Stage 1 |
| 3 | Reusable authentication and authorization context | Stage 2 |
| 4 | One real Android device journey | Stage 3 |
| 5 | Durable execution capture, records, and synchronization | Stage 4 |
| 6 | Claude and Codex integration with the approved first capability set | Stage 5 |
| 7 | Approved deterministic authoring behavior | Stage 6 |
| 8 | iOS support through the same capability path | Stages 4–7 |
| 9 | Desktop integration and supported release | Stages 3–8 |
| 10 | Additional Drizz product capabilities | Stable released foundation |

Stage 3 implements [Authentication and Authorization](authentication.md) and [ADR 0006](decisions/0006-authentication.md) through the [Stage 3 Authentication Implementation Plan](plans/authentication.md). Stages 5 and 6 implement [Agent Integration and Execution Capture](capture.md) and [ADR 0007](decisions/0007-capture.md). The sections below are the implementation plans for those approved architectures; they do not create parallel plans elsewhere.

## 4. Stage 1: Foundation review

### Goal

Replace assumptions with evidence from the existing Drizz systems and approve only the decisions needed for the first implementation.

### Work

- Confirm the first user journey and non-goals.
- Confirm that the approved Platform authentication design can be added without changing existing Drizz login flows.
- Audit Android and iOS device integrations, supported operations, lifecycle, errors, concurrency, and packaging.
- Audit how the current desktop application can start, update, and communicate with the Drizz application.
- Classify existing authoring behavior into agent-owned planning, reusable deterministic behavior, cloud-owned behavior, and excluded behavior.
- Prove that the official Go MCP SDK connects over standard input/output to Claude, Codex, and a generic client.
- Review the official plugin, hook, and structured-event surfaces for Claude and Codex and separate them from SDK instrumentation.
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

The repository owner approves the first journey and the technical owners agree that no unresolved issue blocks the application skeleton, authentication, or first device proof.

## 5. Stage 2: Application skeleton

### Goal

Create the smallest installable application that proves lifecycle and interface composition without duplicating product behavior.

### Work

- Pin the approved Go toolchain and module.
- Add one composition entry point.
- Implement signal handling, cancellation, shutdown, and stable exit behavior.
- Implement typed configuration with explicit precedence.
- Implement structured diagnostics without contaminating MCP standard output.
- Add minimal MCP and CLI adapters over one application use case.
- Add build metadata and version reporting.
- Establish formatting, static analysis, unit testing, race testing, dependency review, secret scanning, file-size enforcement, architecture-boundary checks, and build checks.
- Produce unsigned development artifacts for the approved platform matrix.

### Exit

A clean checkout runs one verification command and builds the development application. MCP and CLI invoke the same non-destructive application behavior.

## 6. Stage 3: Authentication and authorization

### Goal

Provide one authenticated context reusable by MCP, CLI, desktop, and future approved interfaces.

### Work

- Create the isolated Auth0 Platform Native Application and Platform API.
- Implement system-browser Authorization Code with PKCE as the default native sign-in flow.
- Implement Device Authorization as the explicit headless fallback.
- Store durable credentials in the operating-system credential store.
- Keep short-lived session material in memory where practical.
- Reuse the same identity service from CLI, local MCP, and the desktop test boundary without exposing credentials to an agent.
- Derive trusted identity and organization context outside interface input.
- Revalidate cloud authorization at the cloud boundary.
- Add workload identity federation for the first approved CI provider and an isolated Client Credentials fallback only if required.
- Define offline behavior explicitly.
- Handle expired, revoked, missing, corrupted, and organization-changed credentials.
- Test browser interruption, callback failure, restart, concurrent login, logout, refresh rotation, revocation, and credential-store denial.
- Regression-test every existing Drizz login flow affected by shared Auth0 Actions or tenant policy.
- Verify the same authenticated context through MCP, CLI, and the desktop test boundary.

### Delivery slices

The exact slice order, design inventory, planned ownership, contracts, state, failures, files, tests, evidence, rollout, and rollback are defined only in the [Stage 3 Authentication Implementation Plan](plans/authentication.md).

Stage 3 prepares the shared identity application for CLI, MCP capability adapters, and the desktop boundary. It does not add a public MCP tool merely to test authentication. The first authenticated public capability provides the real MCP invocation proof in Stage 6.

### Exit

One user can sign in and out safely on every supported operating system. The shared identity application is composed for future MCP capability adapters without exposing credentials to the MCP client or model. An approved CI workload can authenticate without a human token. Interface input cannot override trusted identity or organization context. Every completion gate in the Stage 3 implementation plan passes.

## 7. Stage 4: First local device journey

### Goal

Execute one approved observation and one approved action on a real Android target through the shared application core.

### Work

- Define device identity and session ownership.
- Define discovery, connect, disconnect, and reconnect behavior.
- Define the platform-neutral device contract without hiding Android behavior.
- Select one low-risk observation and one representative action.
- Define input, output, cancellation, timeout, and error behavior.
- Build the application-owned device interface.
- Integrate the existing Android implementation through an adapter.
- Bound subprocesses, output, memory, concurrency, and device ownership.
- Capture only the evidence required by the product journey.
- Test disconnect, stale device, locked device, application crash, provider crash, cancellation, and restart on real targets.

### Exit

The approved journey works locally on a real supported Android target with deterministic errors and recovery.

## 8. Stage 5: Capture, records, and synchronization

### Goal

Preserve the required capability and agent execution history and synchronize it without adding a cloud round trip to each local device step.

### Work

- Define which execution facts, tool calls, assets, failures, and agent-provided context Drizz is allowed and required to store.
- Define one typed Drizz capture contract for agent events and authoritative capability events. It must record source, fidelity, ordering, and correlation without adopting a provider or observability schema.
- Define bounded pending capture so unrelated Claude or Codex conversations are not retained or synchronized when they never invoke Drizz.
- Define data classification, redaction, consent, and retention.
- Define the cloud synchronization contract and idempotency behavior.
- Select local metadata and artifact storage after a real persistence proof.
- Commit execution state and synchronization intent safely.
- Resume after process termination, machine restart, sleep, and network loss.
- Handle duplicate requests, partial upload, server rejection, and version mismatch.
- Verify cloud acknowledgement before cleanup eligibility.
- Protect active, referenced, failed, and unsynchronized artifacts.
- Enforce disk budgets and observable cleanup.
- Instrument the shared capability core below MCP and CLI so every interface produces the same authoritative execution facts.
- Prove exact and inferred correlation behavior without representing an inferred host relationship as exact.
- Test corruption, disk full, clock changes, competing processes, and upgrade.

### Delivery slices

1. Approve the typed capture contract, source and fidelity model, classification, consent, retention, and limits.
2. Prove the accepted SQLite and artifact design under concurrent processes, crash, migration, corruption, and disk pressure.
3. Record capability intent and outcome inside the shared application path below MCP and CLI.
4. Implement idempotent event and artifact synchronization with restart and partial-upload recovery.
5. Implement acknowledgement-based cleanup, leases, budgets, and protection for active or unsynchronized evidence.
6. Prove bounded pending agent context and exact versus inferred correlation in preparation for supported host adapters.

### Exit

The first device journey appears in the required Drizz surface after synchronization. Its source and capture fidelity are explicit, repeated delivery does not duplicate cloud effects, unrelated agent sessions are not uploaded, and local storage remains bounded without deleting required evidence.

## 9. Stage 6: First agent integration and capability set

### Goal

Expose the approved first-release capabilities consistently through CLI and the supported Claude and Codex MCP clients, with native installation and available host context capture.

### Work

- Product-review every operation name, description, input, output, limit, risk, and destructive effect.
- Define one transport-neutral capability catalog.
- Map approved capabilities to MCP tools.
- Map the same capabilities to CLI commands.
- Build the Claude plugin with MCP configuration and supported lifecycle hooks.
- Build the Codex plugin with MCP configuration and supported lifecycle hooks.
- Implement separate typed Go host adapters for Claude and Codex.
- Resolve and verify the installed Drizz executable and merge host configuration without overwriting unrelated user settings.
- Add the native review, trust, consent, disablement, and removal journeys.
- Publish a tested capability matrix per agent product surface and version.
- Capture only provider-exposed reasoning representations and preserve their source; never claim private chain-of-thought.
- Provide stable machine-readable CLI output and exit behavior.
- Define confirmation and explicit enablement for destructive operations.
- Add client configuration helpers that never overwrite user configuration without consent.
- Validate install, configure, authenticate, invoke, cancel, restart, upgrade, disable, and uninstall with Claude.
- Repeat the supported journey with Codex and a generic MCP client, including honest partial capture where a surface exposes only MCP activity.
- Verify that client-specific behavior does not enter the application core.
- Do not add OpenInference, a Python sidecar, `coding-harness-tracing`, Neatlogs, or an external observability account to this journey.

### Delivery slices

1. Implement the integration manager for detection, exact-path configuration, consent, native trust, verification, update, disablement, and removal.
2. Implement and verify the Claude plugin and typed host adapter against the supported CLI and IDE versions.
3. Implement and verify the Codex plugin and typed host adapter against the supported CLI and desktop versions.
4. Complete the generic MCP client journey with explicitly partial agent context.
5. Run the supported operating-system and client matrix, publish the capture capability matrix, and qualify the first agent-integrated release.

### Exit

Every approved first-release operation has equivalent authorization, execution, recording, cancellation, and failure behavior through MCP and CLI. A non-technical user can connect Claude or Codex through the supported installation journey, approve the integration, use Drizz without manually starting a process, and remove the integration cleanly.

## 10. Stage 7: Deterministic authoring

### Goal

Reuse only approved deterministic authoring behavior while the external agent continues to plan.

### Work

- Identify the minimum authoring inputs and outputs.
- Separate planning, hidden reasoning, and model-specific control from reusable deterministic behavior.
- Review source distribution, intellectual property, assets, dependencies, and packaging.
- Decide whether approved behavior is ported, invoked through a supervised local process, or replaced through a smaller contract.
- Define author, validate, and save journeys without exposing internal implementation vocabulary.
- Define session memory, recovery, loop protection, and cleanup only where deterministic product behavior requires them.
- Verify output through multiple external agents.
- Synchronize the authored result into the approved Drizz surface.

### Exit

Approved authoring works locally without Drizz owning planning and without shipping unapproved proprietary implementation.

## 11. Stage 8: iOS support

### Goal

Extend the proven device, execution, synchronization, and authoring journey to the supported iOS simulator and device matrix.

### Work

- Map shared and iOS-specific Device Bridge behavior.
- Implement supported observations and actions through the existing device contract.
- Preserve platform-specific capabilities and failures where normalization would hide useful behavior.
- Run the approved journey on the supported simulator and physical-device matrix.
- Verify that MCP, CLI, execution, synchronization, and authoring behavior do not fork by interface.

### Exit

The approved iOS journey passes through the same product path and meets the same security, reliability, and evidence requirements as Android.

## 12. Stage 9: Desktop and release

### Goal

Deliver one supported Drizz implementation through standalone and desktop journeys.

### Work

- Approve how the desktop locates, starts, stops, updates, and verifies Drizz.
- Package the local MCP integration for each approved desktop agent surface. A Claude Desktop Extension is MCP-only unless that product exposes a separately approved lifecycle interface.
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

A user installs Drizz through an approved channel, connects a supported client, authenticates, completes the first journey, sees synchronized history, upgrades, and uninstalls without credential leakage or data corruption.

## 13. Stage 10: Additional capabilities

Test plans, debugging, application management, reports, and other product areas are added as separate vertical outcomes.

Each capability must:

- use the same authentication and application boundaries;
- enforce organization and resource authorization in Drizz services;
- reuse execution, artifact, synchronization, and cleanup behavior where its semantics match;
- define its own approved product contract;
- avoid exposing another module's storage or internal provider vocabulary;
- work through every interface that its product requirement names.

OpenInference or another standard trace input is added only when an approved custom agent or model SDK journey requires it. It is not a prerequisite for Claude, Codex, Gemini, or another external agent application and has no current delivery stage.

## 14. Open planning inputs

Dates, staffing, parallel delivery, and estimates remain open until:

- the responsible team is known;
- the first journey is approved;
- the remaining device, desktop, and provider audits are complete;
- supported platforms are chosen;
- proof results expose the real integration cost.

Any schedule created before those inputs is a guess and must be labeled as one.
