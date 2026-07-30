# Technology Stack

Status: Approved recommendation baseline

This document distinguishes current recommendations from decisions that require
proof. A library is not selected merely because it appears here.

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

## 3. Initial foundation

These choices are approved as the initial proof targets. A dependency is added
to production only after its required proof passes.

| Concern | Recommendation | Required proof |
| --- | --- | --- |
| Runtime | Go 1.26.x | Build, startup, shutdown, process supervision, and supported-platform cross-compilation |
| MCP | Official `modelcontextprotocol/go-sdk` | Standard-input/output connection with Claude, Codex, and a generic MCP client |
| CLI | Cobra | Nested command UX, completion, stable machine output, cancellation, and binary impact |
| Logging | `log/slog` | Structured standard-error output with redaction and MCP standard-output isolation |
| Configuration | Typed Go structures with a small explicit loader | Precedence, unknown-key rejection, secret exclusion, and deterministic tests |
| Construction | Manual constructors | Clear lifecycle, no circular dependencies, and test replacement at real boundaries |
| Authentication primitives | `golang.org/x/oauth2` where compatible | Compatibility with the existing Drizz identity service and native-app flow |
| Credential storage | Operating-system credential stores behind an application-owned interface | macOS, Windows, and supported Linux behavior |
| Testing | Go `testing`, `httptest`, and focused comparison helpers | Unit, integration, process, cancellation, and race coverage |

The official MCP documentation currently lists the Go SDK as Tier 1. Its exact
release and supported protocol versions must be pinned and tested when the proof
is built.

## 4. Decisions that require product or architecture evidence

### Local persistence

SQLite is the approved persistence direction for execution metadata,
synchronization state, and leases. Large artifacts remain files with integrity
metadata.

Do not select an ORM, driver, or migration library until the data model and
access patterns exist. The proof must cover locking, crash recovery, migration,
multi-process behavior, disk-full behavior, cleanup, and supported operating
systems.

### Desktop integration

Do not select Protobuf, local HTTP, sockets, or another IPC framework until the
desktop lifecycle is audited. First determine whether the desktop needs a
persistent shared process, a per-request process, streaming events, or multiple
concurrent clients.

### Cloud communication

Use the existing Drizz service contracts where suitable. Do not add Chi,
OpenAPI generation, an HTTP server, a queue client, or a cloud database library
until an approved capability requires it.

### Background work

Synchronization needs bounded, cancellable, restart-safe work. Whether that is
a SQLite-backed journal, a simpler persisted state machine, or an existing
Drizz mechanism depends on the persistence and cloud contract proofs. No
generic job or scheduler framework is selected.

### Observability

Use structured logs from the beginning. Add OpenTelemetry only after trace or
metric export requirements, privacy policy, and destination are known. Domain
and application behavior must not depend on a telemetry vendor.

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
- a TUI;
- a plugin runtime;
- a remote MCP deployment;
- containers or Kubernetes for local device execution.

Deferral means “no demonstrated requirement yet,” not permanent rejection.

## 6. Acceptance before implementation

Before a recommendation becomes an accepted dependency:

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
