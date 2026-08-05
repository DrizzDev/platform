# Dependencies

Status: Approved foundation baseline

This record captures the `DEL-002` evidence for direct runtime dependencies. It
is a living record, not a compliance certificate. Mechanically enforced fields:
license (`scripts/license`), vulnerabilities (`govulncheck`), and cross-platform
build (`scripts/crossbuild`, `CGO_ENABLED=0`), all inside `make verify`.
Accountable owner and upstream-maintenance review remain manual obligations.

The accountable owner for every dependency below is the Drizz Platform
maintainer; the "Upstream" column names the external project, not the Drizz
owner. All runtime dependencies are pure Go and build with `CGO_ENABLED=0` (no
native code). Measured baseline for size review: the development build
`CGO_ENABLED=0 go build -trimpath ./command/drizz` (`go1.26.5`, `darwin/arm64`)
produces a ~25 MiB binary (trimmed, not linker-stripped). Release stripping
(`-s -w`) is a packaging decision, not yet configured.

## Runtime dependencies

| Module | Version | Purpose | License | Upstream / maintenance | Replacement boundary |
| --- | --- | --- | --- | --- | --- |
| `github.com/modelcontextprotocol/go-sdk` | v1.7.0 | Official MCP server protocol | MIT / Apache-2.0 transition | MCP project with Google; active | `internal/transport/mcp` |
| `github.com/spf13/cobra` | v1.10.2 | CLI command routing | Apache-2.0 | spf13; active, widely adopted | `internal/transport/cli` |
| `github.com/getsentry/sentry-go` (+ `/otel`) | v0.48.0 | Error reporting sink; OTel span correlation | MIT | Sentry; active, vendor-maintained | `internal/platform/reporting/sentry` |
| `github.com/samber/slog-sentry/v2` | v2.11.0 | Deliver approved error-level records to the Sentry sink | MIT | samber; active | `internal/platform/reporting/sentry` |
| `go.opentelemetry.io/otel` (+ `sdk`, `metric`, `trace`, OTLP exporters) | v1.44.0 | Traces and metrics with OTLP export | Apache-2.0 | CNCF OpenTelemetry; active | `internal/platform/telemetry` |
| `github.com/zalando/go-keyring` | v0.2.8 | Operating-system credential store access (macOS Keychain, Windows Credential Manager, Linux Secret Service) | MIT | Zalando; active | `internal/identity/infrastructure/vault` |
| `modernc.org/sqlite` | v1.56.0 | Pure-Go SQLite for identity coordination (leases, fenced publication, cleanup) | BSD-3-Clause | modernc; active | `internal/identity/infrastructure/sqlite` |
| `github.com/hashicorp/go-retryablehttp` | v0.7.8 | HTTP retries with jittered backoff, restricted to idempotent requests | MPL-2.0 | HashiCorp; active, widely adopted | `internal/platform/transport` |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | v0.61.0 | HTTP client spans, metrics, and first-party trace-context propagation | Apache-2.0 | CNCF OpenTelemetry; active | `internal/platform/transport` |

`go-keyring` pulls two per-operating-system backends: `github.com/danieljoos/wincred`
(MIT, Windows) and `github.com/godbus/dbus/v5` (BSD-2-Clause, Linux); both are
permissive and verified by `scripts/license`. `modernc.org/sqlite` pulls the
`modernc.org` runtime family (`libc`, `memory`, `mathutil`, `gc`), all BSD-3-Clause
and permissive. Neither adds CGO — the build stays `CGO_ENABLED=0` and cross-builds
for all six targets, which is the reason a pure-Go SQLite was chosen over a
CGO driver (`DEL-001`, plan §10). Binary-size effect on the `drizz` command is
currently zero for both because the identity module is not yet linked into the
released binary; it is re-measured when the host wires it (the `modernc.org/sqlite`
contribution is expected to be several megabytes and will be recorded then).
`govulncheck` reports no known vulnerability. Necessity: the standard library has
no operating-system credential store and no embedded SQL engine, and the plan
requires the OS vault (`SEC-007`) and SQLite coordination (ADR 0004).

The resilient HTTP transport (`internal/platform/transport`) composes both new
direct dependencies with the standard library: per-attempt and dial timeouts, a
capped response body (`SEC-022`), `otelhttp` observability, and `go-retryablehttp`
retries that fire only on idempotent methods so a rotating refresh token is never
resubmitted after a lost response. `go-retryablehttp` pulls `github.com/hashicorp/go-cleanhttp`
(MPL-2.0) and `github.com/hashicorp/go-hclog` (MIT); `otelhttp` pulls
`github.com/felixge/httpsnoop` (MIT). Both vendors stay behind the transport
adapter (`DEL-004`, `framework_test.go`). Adopting `go-retryablehttp` and its
`go-cleanhttp` transitive brought **MPL-2.0** (weak, file-level copyleft) into the
tree for the first time; `scripts/license` was extended to accept MPL-2.0 by owner
decision (`DEL-006`), keeping every other license permissive as before.

Transitive dependencies are pinned by `go.sum`; their licenses are verified for
every module by `scripts/license`, and their vulnerabilities by `govulncheck`.
Each vendor stays behind its owning adapter package; no vendor type crosses into
the application or domain layers (`DEL-004`, enforced by
`tests/architecture/framework_test.go`).

## Coverage gaps (owner review)

- `scripts/license` verifies imported Go modules only. Development tools (`go.mod`
  `tool` block), GitHub Actions, pre-commit hooks, and the pinned Gitleaks and
  pre-commit versions are pinned reproducibly but their licenses are not yet
  machine-verified. All are permissive (MIT/Apache-2.0/BSD) by inspection.
- Upstream maintenance status is recorded by inspection, not automatically
  tracked.

## Tooling and automation

| Tool | Version | Purpose | Pinned in |
| --- | --- | --- | --- |
| Go toolchain | 1.26.5 | Build and test | `go.mod`, `.github/workflows/verify.yml` |
| `golangci-lint`, `staticcheck`, `goimports`, `govulncheck` | `go.mod` `tool` block | Static analysis, formatting, vulnerabilities | `go.mod` |
| `gitleaks` | v8.30.1 | Secret scanning | `scripts/secret` |
| `pre-commit` | 4.6.1 | Hook runner | `requirements.txt` |
| `pre-commit/pre-commit-hooks` | pinned commit SHA | File hygiene | `.pre-commit-config.yaml` |
| GitHub Actions (`checkout`, `setup-go`, `setup-python`) | pinned commit SHA | CI | `.github/workflows/verify.yml` |

## Policy

- Prefer the standard library and existing approved dependencies (`DEL-001`).
- Every new direct dependency records the fields above and passes
  `scripts/license` before it is accepted.
- Permissive licenses (MIT, Apache-2.0, BSD, ISC) are approved. MPL-2.0 (weak,
  file-level copyleft) is approved by owner decision for library use where the
  source is unmodified. Any other license (notably strong copyleft) requires
  legal review, an ADR, and an owner decision (`DEL-006`).
- Actions, hooks, and tools are pinned reproducibly (`DEL-020`).
