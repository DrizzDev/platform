# Stage 3 Authentication Threat Model

Status: Owner-accepted; Slice 2 unblocked (F4, F5 evidence still pending)

This document owns the Stage 3 authentication threats required by
[the plan](../plans/authentication.md) §9. It has the eight named rows the plan
requires. Each row records assets, entry points, control, detection, test,
residual risk, and owner. An unowned high or critical residual risk blocks
Slice 2.

The Platform owner accepts accountability for every row and its residual risk as
of 2026-08-04. No row is high or critical after controls, so the plan §9 blocker
is cleared. Two rows still carry a pending evidence fact (T1 on F5, T5 on F4)
that must resolve before those rows reach `LOW`; this constrains Slices 7 and 8,
not Slice 2.

## Method and scope

Scope is the identity module and its boundaries: system browser and loopback,
Device Authorization, the operating-system vault, SQLite coordination, Drizz
Cloud, and GitHub Actions federation. Out of scope: remote HTTP MCP
authorization, Dynamic Client Registration, and the machine-credential fallback
(all excluded by the plan §3).

Residual-risk scale after controls: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.
Controls marked pending depend on an evidence fact in
[evidence](../evidence/authentication.md); their residual risk cannot drop below
`MEDIUM` until that fact resolves.

## Summary

| Row | Threat | Residual | Owner |
| --- | --- | --- | --- |
| T1 | Replay | MEDIUM (pending F5) | Platform owner |
| T2 | Downgrade | LOW | Platform owner |
| T3 | Credential theft | LOW | Platform owner |
| T4 | Tampering | LOW | Platform owner |
| T5 | Tenant crossover | MEDIUM (pending F4) | Platform owner |
| T6 | Provider compromise | MEDIUM | Platform owner |
| T7 | Support misuse | LOW | Platform owner |
| T8 | Dependency integrity | LOW | Platform owner |

No row is `HIGH` or `CRITICAL` after controls. The Platform owner has accepted
every row (see the acceptance note above); two rows also carry a pending
evidence fact that constrains Slices 7 and 8.

## T1: Replay

- **Assets**: authorization code, ID token, device code, GitHub OIDC assertion.
- **Entry points**: loopback callback, token endpoint, device polling, workload
  exchange.
- **Control**: PKCE `S256` binds the code to a per-attempt verifier; nonce binds
  the ID token to the attempt; state binds the callback; the callback consumes
  exactly one request; the workload exchange requires a unique token identifier
  and atomic one-time exchange or a bounded replay store.
- **Detection**: reject a second callback, a mismatched nonce or state, and a
  repeated token identifier; rely on Auth0 code single-use.
- **Test**: replayed code, replayed ID token, replayed assertion, and duplicate
  identifier are rejected (Slices 4, 8).
- **Residual**: `MEDIUM` until evidence F5 confirms atomic exchange or a replay
  store owner and retention; `LOW` thereafter.
- **Owner**: Platform owner (accepted).

## T2: Downgrade

- **Assets**: flow integrity, transport confidentiality.
- **Entry points**: authorization request, redirect, provider transport, device
  path.
- **Control**: implicit, password, password-realm, and Client Credentials
  disabled on the native application; PKCE always `S256`, never plain; redirect
  restricted to a `127.0.0.1` loopback; Device Authorization has no browser
  fallback; all provider traffic is over TLS.
- **Detection**: reject a non-`S256` challenge, a non-loopback redirect, and a
  non-TLS endpoint.
- **Test**: forced plain PKCE, forced implicit response, and a non-loopback
  redirect are rejected (Slices 4, 5).
- **Residual**: `LOW`.
- **Owner**: Platform owner (accepted).

## T3: Credential theft

- **Assets**: refresh token, access token, provider subject.
- **Entry points**: on-disk storage, logs, telemetry, environment, process
  arguments, another local process, MCP traffic, model context.
- **Control**: refresh tokens live only in the operating-system vault as
  immutable versioned entries; access tokens live in process memory only; a
  fenced epoch prevents last-writer account switching; diagnostics and reporting
  are code-only, so no token or cause reaches logs, traces, or Sentry; no token
  enters environment, arguments, MCP messages, or model context.
- **Detection**: privacy canaries across every prohibited boundary (plan §12).
- **Test**: canary absence in arguments, environment, standard output and error,
  logs, metrics, traces, Sentry, crash output, SQLite, paths, support bundles,
  prompts, model context, hooks, and MCP messages (Slices 3, 4, 6).
- **Residual**: `LOW`. A locked or absent vault fails closed with
  `SECURE_STORAGE_UNAVAILABLE` and retains no copied credential.
- **Owner**: Platform owner (accepted).

## T4: Tampering

- **Assets**: callback response, provider response, loopback listener.
- **Entry points**: loopback socket, provider transport, a malicious local
  caller.
- **Control**: validate method, path, size, state, issuer, audience, signature,
  expiry, nonce, and schema; bind the listener to `127.0.0.1` only, not
  `0.0.0.0` or a name, which resists DNS rebinding; consume one `GET`; reply
  with at most one kilobyte of fixed Drizz-authored content under
  `Content-Security-Policy: default-src 'none'`, `Cache-Control: no-store`,
  `Referrer-Policy: no-referrer`, and `X-Content-Type-Options: nosniff`; reflect
  no provider or query value; a bad method, state, or schema returns a fixed
  rejection and closes the listener; three stray paths close the listener.
- **Detection**: rejection counters and a fixed rejection response.
- **Test**: tampered signature, wrong issuer or audience, oversized or malformed
  callback, and stray-path flooding are rejected without reflection (Slice 4).
- **Residual**: `LOW`.
- **Owner**: Platform owner (accepted).

## T5: Tenant crossover

- **Assets**: organization data, Platform audience.
- **Entry points**: token audience, cloud request, client-supplied organization
  identifier.
- **Control**: validate the exact Platform audience; Drizz Cloud is the sole
  authority per operation; an organization or resource identifier from a client,
  flag, file, or cache is request data, never authority; a wrong audience or a
  changed membership returns a stable forbidden or authentication-required
  result.
- **Detection**: audience-mismatch rejection and cloud denial evidence.
- **Test**: wrong audience, cross-organization request, and changed membership
  are denied at the cloud boundary (Slice 7).
- **Residual**: `MEDIUM` until evidence F4 fixes the cloud contract; `LOW`
  thereafter.
- **Owner**: Platform owner (accepted).

## T6: Provider compromise

- **Assets**: trust in Auth0 and the GitHub OIDC issuer.
- **Entry points**: discovery, JWKS, token, device, revocation, exchange
  endpoints.
- **Control**: fail closed when a provider is unavailable; bound every response
  size and every connect and total timeout; pin one issuer with a bounded key
  cache and a maximum freshness window; permit one controlled key refresh after
  an unknown key identifier; never blindly retry token exchange, refresh,
  publication, or revocation.
- **Detection**: bounded-timeout and size-limit rejections; freshness expiry.
- **Test**: unavailable provider, oversized discovery or JWKS, and an unknown
  key identifier follow the bounded policy (Slices 4, 6, 8).
- **Residual**: `MEDIUM`. Inherent trust in the identity provider remains and is
  accepted with these bounds.
- **Owner**: Platform owner (accepted).

## T7: Support misuse

- **Assets**: provider subject, organization identity.
- **Entry points**: support bundles, diagnostics, crash output.
- **Control**: the provider subject stays inside the vault record and never
  enters SQLite, paths, output, telemetry, or support data; organization
  identity and roles are not persisted locally in Stage 3; failures are
  code-only.
- **Detection**: privacy canaries in support and crash paths.
- **Test**: subject and organization canary absence in support bundles and crash
  output (Slices 3, 6).
- **Residual**: `LOW`.
- **Owner**: Platform owner (accepted).

## T8: Dependency integrity

- **Assets**: the four approved dependencies and their transitive graph.
- **Entry points**: module download, build, release artifact.
- **Control**: pinned modules in `go.sum`; `govulncheck`, license, and secret
  checks in the merge gate; pure-Go build with `CGO_ENABLED=0`; signatures,
  checksums, dependency inventory, software bill of materials, and provenance at
  release (Gate 9); a `DEL-002` proof per dependency.
- **Detection**: vulnerability and license gate failures; release verification.
- **Test**: gate failure on a known vulnerable or unapproved-license module;
  clean-machine release verification (Slice 3 onward, Gate 9).
- **Residual**: `LOW`.
- **Owner**: Platform owner (accepted).

## Gate exit

The Platform owner has accepted every row, and no row is `HIGH` or `CRITICAL`
after controls, so **Slice 2 is unblocked**. T1 and T5 additionally need
evidence F5 and F4 resolved to reach `LOW`, which constrains Slices 7 and 8.
Slices 3 and later still wait on the evidence rows that name them.
