# Technology Stack

Status: Approved foundation baseline

This document distinguishes implemented foundation choices from decisions that
still require product or integration proof. A library is not selected merely
because it appears here.

## 1. Selection principles

A dependency must solve a demonstrated requirement and be evaluated for:

- product fit;
- maintenance and release activity;
- security and license;
- cross-platform behavior;
- startup, memory, and binary impact;
- replacement cost;
- compatibility with the supported Go version.

Prefer the Go standard library when it provides a clear and maintainable
solution. Pin exact versions only when implementation begins.

## 2. Language recommendation

| Option | Product fit | Main concern | Recommendation |
| --- | --- | --- | --- |
| Go | Native executable, fast startup, process control, concurrency, cross-platform release, official MCP SDK | Existing systems need explicit process or API boundaries | Recommended for the local platform host |
| TypeScript | Strong MCP and desktop ecosystem | Heavier standalone runtime and packaging model | Keep in existing TypeScript products; do not choose as the new host without a blocking Go limitation |
| Python | Useful for existing Python behavior | Local environment, startup, source distribution, and packaging | Keep behind an approved boundary only when existing behavior must run locally |
| Rust | Native performance and memory safety | Higher team and integration cost for the same immediate outcome | Reconsider only if measured safety or performance requirements exceed Go |

Use the latest supported patch of Go 1.26 for initial proofs. The exact patch
must be pinned in the repository toolchain. Go 1.26 was released in February
2026 and remains supported under the Go release policy.

Official agent command hooks invoke an executable and exchange structured data,
so Claude and Codex adapters can use the installed Go application directly.
Python's broader model-SDK instrumentation ecosystem does not justify a Python
runtime for the primary external-agent journey.

## 3. Foundation choices

Implemented choices are active in the repository. Pending choices are approved
directions but do not authorize a library until their required proof passes.

| Concern | Choice | Status | Required evidence |
| --- | --- | --- | --- |
| Runtime | Go 1.26.x | Implemented | Build, startup, shutdown, process supervision, and supported-platform cross-compilation |
| MCP | Official `modelcontextprotocol/go-sdk` | Implemented foundation transport | Standard-input/output connection with each supported MCP client |
| CLI | Cobra | Implemented foundation entry point | Nested commands, stable machine output, cancellation, and binary impact |
| Logging | `log/slog` | Implemented | Line-delimited JSON, stable fields, correlation scope, redaction, and MCP standard-output isolation |
| Error reporting | Sentry Go SDK with the `slog-sentry` handler | Implemented | Approved code-only reporting, source attribution, asynchronous delivery, and bounded shutdown |
| Configuration | Typed Go structures with a small explicit loader | Implemented | Precedence, unknown-key rejection, secret exclusion, and deterministic tests |
| Construction | Manual constructors | Implemented | Clear lifecycle, no circular dependencies, and replacement at real boundaries |
| Authentication | Auth0 plus reviewed OAuth and OIDC primitives | Architecture approved; implementation pending | PKCE, Device Authorization, token validation, workload identity, and MCP OAuth compatibility |
| Credential storage | Operating-system credential stores behind an application-owned interface | Direction approved; implementation pending | macOS, Windows, and supported Linux behavior |
| Agent integration | Official MCP, plugin, hook, transcript, and structured-event surfaces | Architecture approved; implementation pending | Real Claude and Codex install, capture, compatibility, privacy, and removal journeys |
| Testing | Go `testing`, `httptest`, and focused comparison helpers | Implemented foundation gates | Unit, integration, process, cancellation, and race coverage |

The official MCP documentation currently lists the Go SDK as Tier 1. Its exact
release and supported protocol versions must be pinned and tested when the proof
is built.

### Logging

Use the standard `log/slog` API and JSON handler. The official MCP Go SDK
accepts `*slog.Logger` directly, so another logging API would require an adapter
and create two logging contracts. `slog` already provides structured
attributes, context-aware calls, immutable logger enrichment, stable handler
interfaces, and line-delimited JSON.

Use `HandlerOptions.ReplaceAttr` for the Drizz redaction and field policy. Do
not add a redaction package until it demonstrates stable maintenance, precise
control of Drizz fields, and a clear benefit over this standard-library
extension point. Reconsider Zerolog or Zap only if production measurements show
logging is a material CPU or allocation bottleneck.

The composition root constructs observability once per process and injects the
same providers into CLI, MCP, server, desktop, and background flows. No feature
creates its own logger, reporter, tracer, or meter.

Sentry uses `github.com/getsentry/sentry-go` v0.48.0 behind the reporting sink
boundary. `github.com/samber/slog-sentry/v2` v2.11.0 delivers only approved
error-level records created by that boundary. The production reporting surface
accepts a stable Drizz event code and no cause, message, or arbitrary attribute,
so provider and user content cannot reach diagnostics or Sentry through it.
The official `github.com/getsentry/sentry-go/otel` integration may correlate an
approved event with the active OpenTelemetry span without creating another
tracing pipeline. Both adapters remain replaceable; another reporting vendor
implements the same sink and is registered once in the reporting provider.
Logging callers and transports do not change.

Sentry is disabled when `DRIZZ_SENTRY_DSN` is absent. `DRIZZ_SENTRY_SAMPLE_RATE`
controls error-event sampling and defaults to `1`. Every Drizz-owned setting uses
the `DRIZZ_` prefix so inherited third-party `SENTRY_*` or `OTEL_*` variables
never change Drizz behavior. The adapter receives only approved `ERROR` and
higher records, limits breadcrumbs and event depth, and flushes once
during bounded shutdown. Automatic collection of user data, cookies, headers,
HTTP bodies, query parameters, machine identity, Sentry-native logs, and
Sentry-native metrics is explicitly disabled. OpenTelemetry remains the single
owner of traces, spans, and metrics; Sentry-native tracing is disabled.

The selected Sentry Go SDK does not expose a native profile-sampling option.
Profiling is therefore deferred until a supported mechanism, destination,
sampling budget, runtime cost, and privacy contract are proven. It MUST NOT be
silently enabled through a future dependency update.

## 4. Decisions that require product or architecture evidence

### Local persistence

SQLite is the approved persistence direction for execution metadata,
synchronization state, and leases. Large artifacts remain files with integrity
metadata.

The final Stage 3 authentication plan defines the first concrete access pattern:
non-secret identity coordination with fenced credential publication and durable
cleanup. Its pure-Go SQLite driver remains unapproved until the Stage 3
dependency proof covers transactions, crash recovery, migration, multi-process
behavior, disk-full behavior, cleanup, size, and every supported operating
system. Execution, synchronization, and artifact access patterns remain open
and do not inherit that driver choice automatically. No ORM or migration
framework is selected.

### Desktop integration

Do not select Protobuf, local HTTP, sockets, or another IPC framework until the
desktop lifecycle is audited. First determine whether the desktop needs a
persistent shared process, a per-request process, streaming events, or multiple
concurrent clients.

### Cloud communication

Use the existing Drizz service contracts where suitable. Do not add Chi,
OpenAPI generation, an HTTP server, a queue client, or a cloud database library
until an approved capability requires it.

### Authentication

Auth0 is the approved identity provider. The protocol and security contract is
defined in [Authentication and Authorization](authentication.md); a Go library
is not approved until a proof covers PKCE, Device Authorization, discovery,
token validation, refresh rotation, cancellation, and supported-platform
credential storage. `golang.org/x/oauth2` is a candidate primitive, not an
authorization framework and not a substitute for explicit validation.

### Background work

Synchronization needs bounded, cancellable, restart-safe work. Whether that is
a SQLite-backed journal, a simpler persisted state machine, or an existing
Drizz mechanism depends on the persistence and cloud contract proofs. No
generic job or scheduler framework is selected.

### Observability

Use structured logs and OpenTelemetry interfaces from the beginning. Export is
disabled by default and is enabled only when its destination and privacy policy
are configured. Domain and application behavior must not depend on a telemetry
vendor.

OpenInference is not selected for the primary Claude, Codex, or other external
agent application journey. Those surfaces are integrated through their official
MCP, plugin, hook, transcript, and structured-event contracts. A future SDK
trace integration may evaluate OpenInference behind a Drizz-owned adapter; it
cannot define the product record. See
[Agent Integration and Execution Capture](capture.md).

### Packaging

GoReleaser is a candidate for signed multi-platform releases and package-manager
artifacts. It is selected only after the supported operating-system,
architecture, signing, notarization, update, and distribution requirements are
approved.

## 5. Explicitly deferred

The current requirements do not justify selecting:

- an HTTP router;
- an ORM;
- a migration framework;
- a local message broker;
- Redis or another cache service;
- Kafka, NATS, or RabbitMQ;
- a scheduler framework;
- a runtime dependency-injection container;
- a public SDK;
- OpenInference ingestion or export;
- a Python agent-instrumentation sidecar;
- a TUI;
- a Drizz-owned plugin runtime;
- a remote MCP deployment;
- containers or Kubernetes for local device execution.

Deferral means “no demonstrated requirement yet,” not permanent rejection.

## 6. Dependency acceptance

Before a new recommendation becomes an accepted dependency:

1. record the exact requirement it solves;
2. run the required proof against real supported systems;
3. review license, maintenance, vulnerabilities, transitive dependencies, and
   release cadence;
4. measure relevant startup, memory, binary, or runtime cost;
5. document the replacement boundary;
6. obtain repository-owner approval for foundational choices.

## 7. Sources

- [Go release history](https://go.dev/doc/devel/release)
- [MCP SDK support tiers](https://modelcontextprotocol.io/docs/sdk)
- [Official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Cobra](https://github.com/spf13/cobra)
- [OAuth for Native Apps](https://www.rfc-editor.org/rfc/rfc8252)
- [OAuth Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [OAuth Device Authorization Grant](https://www.rfc-editor.org/rfc/rfc8628)
- [OAuth Resource Indicators](https://www.rfc-editor.org/rfc/rfc8707)
- [OAuth Protected Resource Metadata](https://www.rfc-editor.org/rfc/rfc9728)
- [OAuth Token Exchange](https://www.rfc-editor.org/rfc/rfc8693)
- [MCP authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization)
