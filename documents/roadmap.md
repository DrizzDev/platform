# Delivery Roadmap

Status: Approved delivery baseline (revised)

Current position: foundation, application skeleton, and authentication are complete. The device journey, execution capture, and agent integration with the full first capability set are built and proven green. The next delivery stage is packaging the result into one installable binary and shipping the first version to the team.

This roadmap describes delivery order. It does not approve public tools, contracts, architecture decisions, dependencies, staffing, or dates.

## 1. First release outcome

A user can install Drizz with one standard command, authenticate, connect it to a supported MCP client, and use the approved local device capabilities — on Android and iOS — through an external agent, with every action recorded locally. Synchronizing those records to Drizz's cloud is real product behaviour but is deliberately a later step; the first release is complete and useful without it.

The command line and the desktop reuse the same application behaviour. External agents own planning.

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
| 4 | First real device journey — Android and iOS through one neutral path | Stage 3 |
| 5 | Durable execution capture and local records; cloud synchronization | Stage 4 |
| 6 | Agent integration with the approved first capability set, on both platforms | Stage 5 |
| 7 | Packaged, installable, supported release — the first shippable version | Stage 6 |
| 8 | Additional product capabilities, added incrementally | Stage 7 |
| 9 | Deterministic authoring | Stage 8 |

Two revisions from the original baseline are folded in here and explained below:

- **iOS is not a separate stage.** The Device port is platform-neutral and the device helper already drives Android and iOS; every capability routes by the device's own platform. iOS is therefore delivered inside Stages 4 and 6, not as a later stage of its own. The only iOS-specific work is threading the platform through the emulator/simulator commands and testing on Apple hardware — both part of the device work, not a new architecture.
- **Packaging is pulled forward and is the first shippable version.** Turning the two programs Drizz is built from into one installable binary, delivered through standard channels, is what lets the team actually use it. It becomes Stage 7 — the first release — rather than being bundled with far-later work. **Deterministic authoring is moved to the end** (Stage 9): it is added after the team is using the released product, alongside other new capabilities, one at a time.

Stage 3 implements [Authentication and Authorization](authentication.md) and [ADR 0006](decisions/0006-authentication.md) through the [Stage 3 Authentication Implementation Plan](plans/authentication.md). Stages 5 and 6 implement [Agent Integration and Execution Capture](capture.md) and [ADR 0007](decisions/0007-capture.md). Stage 7 follows the [Packaging and Distribution Plan](plans/packaging.md). The sections below are the delivery record and plan for those architectures; they do not create parallel plans elsewhere.

## 4. Stage 1: Foundation review — complete

Evidence-backed audit of the existing Drizz device and desktop systems, a proven MCP stdio connection to Claude and Codex, proven Go supervision of the required local providers, the approved first journey, and accepted foundational decision records. Complete.

## 5. Stage 2: Application skeleton — complete

The smallest installable application: one composition entry point, signal handling and clean shutdown, typed configuration, structured diagnostics that never contaminate MCP standard output, minimal MCP and CLI adapters over one use case, build metadata, and the full verification pipeline (format, lint, unit, race, dependency, secret, file-size, architecture-boundary, and build checks). Complete.

## 6. Stage 3: Authentication and authorization — complete

One authenticated context reusable by MCP, CLI, and the desktop boundary: system-browser Authorization Code with PKCE by default, Device Authorization as the headless fallback, durable credentials in the operating-system credential store, trusted identity and organization derived outside interface input and revalidated at the cloud boundary. Merged and proven live. The isolated CI workload identity is deferred until a workload actually needs it.

## 7. Stage 4: First device journey — Android and iOS through one neutral path

### Goal

Execute the approved observations and actions on a real device through the shared application core, on **both** Android and iOS, behind one platform-neutral Device port.

### Work

- Define device identity, session ownership, discovery, connect, disconnect, and reconnect.
- Define the platform-neutral device contract without hiding platform behaviour; every request routes by the device's own platform (Android, iOS device, iOS simulator).
- Integrate the existing device helper through an adapter over stdio JSON-RPC; supervise it, bound its subprocesses, output, memory, and concurrency, and pin its digest.
- The device helper manages the on-device components itself — installing Drizz's instrumentation on Android, building and provisioning WebDriverAgent on iOS — and detects and guides on missing host toolchains (`adb`, Xcode) rather than failing opaquely.
- Thread the platform through the emulator and simulator commands so device management is not Android-only.
- Capture only the evidence the product journey requires.
- Test disconnect, stale device, locked device, application and provider crash, cancellation, and restart on real Android and iOS targets.

### Status

The neutral port, the helper adapter, supervision, digest pinning, and the read and interaction commands are built and proven on a real Android device. iOS routes through the same commands today; the remaining device work is the simulator/emulator platform threading and the live pass on Apple hardware.

## 8. Stage 5: Capture, local records, and synchronization

### Goal

Preserve the required capability and agent execution history locally, then synchronize it to Drizz without adding a cloud round trip to each local device step.

### Work and status

- The typed capture contract, source and fidelity model, classification, consent, retention, and limits are **approved and built**.
- The local stores — an ordered SQLite journal and a content-addressed artifact store — are **built and qualified** under concurrent processes, crash, migration, corruption, and disk pressure.
- Capability intent and outcome are **recorded inside the shared application path** below MCP and CLI, so every interface produces the same authoritative facts. Host-side agent context (prompts, responses) is recorded through the hook path as an inferred host observation, never relabelled as an exact Drizz event.
- **Cloud synchronization is the remaining piece and may land after the first release.** It targets Drizz's own backend API — authenticated with the user's token, organization-authorized — never a direct client-to-object-store write. Large artifacts use short-lived backend-issued upload URLs; storage stays server-controlled and the client never holds cloud credentials. Synchronization is idempotent by identity, resumable after restart and partial upload, and reclaims local disk only after the cloud acknowledges.

### Exit

Local records are durable and correct today. Once synchronization is enabled, the first device journey appears in the required Drizz surface, repeated delivery does not duplicate cloud effects, unrelated agent sessions are not uploaded, and local storage stays bounded without deleting required evidence.

## 9. Stage 6: Agent integration and the first capability set — built

### Goal

Expose the approved first-release capabilities consistently through the command line and the supported Claude, Codex, and other MCP clients, on both platforms, with native installation and optional host-context capture.

### Status

Built and proven green. One transport-neutral catalog is the single source for the full device command set; the agent connection (MCP) and the command line each render from it. A single `drizz connect` command wires Drizz into an agent's configuration without disturbing other settings, across every supported agent, and a separate opt-in step registers turn-event hooks that record prompt and response context. Every flow emits traces, metrics, and structured logs, with privacy canaries proving device content and prompts never reach telemetry.

The remaining Stage 6 work is run alongside Stage 7: the live install→invoke→capture pass against the actual Claude and Codex applications, and the live device pass on the few state-changing and emulator commands.

## 10. Stage 7: Packaged, installable, supported release — the first shippable version

### Goal

Deliver Drizz as **one installable binary** that carries the device helper inside it, through standard channels, so a person installs it with one command and connects an agent without any manual setup. This is the first version shared with the team.

### Work

- Compile the device helper into a single self-contained native executable per operating system and architecture, and embed the matching helper — and its pinned digest — into the platform binary at build time.
- On first use, extract the embedded helper to a protected per-user location, verify its digest, and run it; keep the environment-variable override for development and continuous integration.
- Produce the standard install experiences with one release tool: a Homebrew tap (`brew install`), a shell installer (`curl … | sh`), and the equivalent on Windows — plus archives and checksums.
- Sign binaries, and notarize where a channel requires it; for the command line delivered through Homebrew and the shell installer, this is largely unnecessary and is deferred until a channel needs it.
- Define install, update, rollback, downgrade, and uninstall behaviour, including what happens to local records and credentials on removal.
- Run the supported operating-system and client matrix, including the live Claude and Codex install→invoke→capture journey and the live device pass, and qualify the first release.

### Exit

A person installs Drizz with one standard command, connects Claude or Codex, uses the full device capability set on Android or iOS without manually starting anything, and removes the integration cleanly. This is the version the team uses.

## 11. Stage 8: Additional product capabilities

New product areas — reporting, debugging aids, application management, and others — are added as separate vertical outcomes, one at a time, on the released foundation. A capability that reads cloud-held data rather than a local device may be delivered as a hosted, link-based MCP server instead of the local binary. Each capability reuses the same authentication and application boundaries, enforces organization and resource authorization in Drizz services, reuses execution, artifact, synchronization, and cleanup behaviour where its semantics match, defines its own approved product contract, and works through every interface its requirement names.

## 12. Stage 9: Deterministic authoring

Reuse only approved deterministic authoring behaviour while the external agent continues to plan: identify the minimum authoring inputs and outputs, separate planning and model-specific control from reusable deterministic behaviour, decide whether approved behaviour is ported or invoked through a supervised local process, define author/validate/save journeys without exposing internal vocabulary, verify output through multiple external agents, and synchronize the authored result into the approved Drizz surface. This is delivered after the team is using the released product.

## 13. Open planning inputs

Dates, staffing, parallel delivery, and estimates remain open until the responsible team is known, the first release journey is qualified on the supported matrix, and proof results expose the real integration and packaging cost. Any schedule created before those inputs is a guess and must be labelled as one.
