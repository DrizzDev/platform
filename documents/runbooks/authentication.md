# Stage 3 Authentication — Release Qualification (Gate 9)

Status: In progress — macOS proven; cross-platform, signing, and supply-chain evidence outstanding.

This runbook tracks Gate 9 qualification for the Stage 3 identity work. Slices 2 through 7 (typed contracts, vault and SQLite, browser login, device login, session refresh, logout, and organization authorization) are complete, green through `make verify`, and proven end to end against real Auth0 and staging.

## Owner decision (2026-08-05)

Slice 8 (workload authentication) is **deferred**: it is not the current priority. Gate 9 proceeds now for the completed user-authentication vertical, starting with macOS, which is the only supported vault available to the owner. Linux and Windows qualification and Slice 8 are tracked as outstanding work below.

## Done

| Check | Evidence | Command |
| --- | --- | --- |
| Full merge gate | Green | `make verify` |
| Real macOS Keychain round-trip | Pass | `go test -tags keyring ./tests/identity -run TestKeyring` |
| Cross-process fenced publication | Pass | `go test -tags fencing ./tests/identity -run TestContention` |
| MCP process behaviour | Pass | `go test ./tests/process` (in `make verify`) |
| Browser login, device login, logout | Proven live | attended run on macOS |
| Organization authorization | Proven live against staging | `drizz login` with `DRIZZ_CLOUD` set returns the resolved organization |
| Secret hygiene, license, vulnerability | Pass | `make verify` (`secret`, `license`, `vulnerability`) |
| Cross-compilation, all six targets | Pass | `scripts/crossbuild` (`CGO_ENABLED=0`) |

## Outstanding

### TODO(slice-8): Workload authentication (GitHub Actions OIDC)
Deferred by owner decision; not primary. When resumed:

- Resolve evidence **F5** first (audience, subject-claim shape, assertion lifetime, exchange endpoint, and replay policy) from `infra` and the cloud setup, exactly as F4 was resolved for Slice 7.
- Build `application/workload` (`Flow.Exchange`), `infrastructure/github` (assertion adapter over the Actions OIDC request), the exchange adapter, and a replay guard if the exchange provider lacks atomic one-time exchange.
- Reuse the transport, failure contract, and `grant.Credential`.
- Client-side effort is medium; the end-to-end size depends on whether the cloud exchange endpoint and replay store already exist.

### TODO(cross-platform): Linux and Windows vault qualification
No Linux or Windows machine is available to the owner. The real-vault tests exist and run per operating system in continuous integration:

- Linux Secret Service: `go test -tags keyring ./tests/identity` on a protected Linux runner with an isolated Secret Service.
- Windows Credential Manager: same, on a protected Windows runner with an isolated account.
- Decide the runner strategy (protected CI matrix) before release. macOS is the only surface proven locally today.

### TODO(release-engineering): Signing and supply chain
None of this tooling exists yet; it is release-engineering, not application code.

- Code signing: macOS notarization and Windows signing (certificates and pipeline).
- Software bill of materials and build provenance binding signed output to the reproducible unsigned input.
- Release artifact checksums.

### TODO(lifecycle): Install, upgrade, rollback, uninstall
- A clean-machine install, an upgrade, a rollback, and an uninstall on the supported matrix. A manual macOS install and uninstall smoke (build, login, logout, confirm the Keychain entry and SQLite pointer are cleared) has been run; the packaged lifecycle is not yet built.

### TODO(review): Independent security review
- A formal independent security review before release (DEL-028).

### TODO(infra): Codify the Auth0 Platform application
- The Platform native application and API were created by hand in the Auth0 dashboard. Codify them in the `infra` repository with import blocks so they are not dashboard-only. Set `DRIZZ_CLOUD` in the real staging and production configuration (only a local value exists today).

## Regression

The new Platform application is additive: the shared post-login Actions are gated to the existing web and desktop clients and skip the new client, so existing sign-ins are unaffected (evidence `documents/evidence/authentication.md` F1 and the Auth0 Action review). A formal regression run of the existing web and desktop login journeys remains part of release sign-off.
