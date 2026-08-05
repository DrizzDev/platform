# Stage 3 Authentication Evidence

Status: Draft — Gate 1 in progress

This document owns the current facts and approvals that
[the Stage 3 plan](../plans/authentication.md) Gate 1 must resolve before any
production code. It records only established or pending facts. A Drizz-specific
fact that requires access this agent does not hold is marked `PENDING` with its
owner; it is never invented. An unresolved fact that a slice depends on blocks
that slice.

## Evidence record schema

Every fact below carries:

- **Source** — where the fact was read or must be read.
- **Revision** — commit, tenant export, or dated inspection that fixes it.
- **Environment** — staging, production, or public standard.
- **Date** — recording date.
- **Owner** — the accountable party who confirms or supplies it.
- **Result** — `ESTABLISHED`, `PROPOSED`, or `PENDING`.

Drafting date for this pass: 2026-08-04.

## Resolution summary

| Fact | Result | Owner | Blocks |
| --- | --- | --- | --- |
| F1 Current Auth0 tenant state | RESOLVED (repo) | Platform owner | — |
| F2 Final Auth0 resource proposal | PENDING (inputs gathered) | Platform owner | Slice 4, 5, 8 |
| F3 Supported operating-system matrix | PROPOSED | Platform owner | Slice 3 |
| F4 Drizz Cloud organization endpoint and contract | RESOLVED (repo) | Platform owner | Slice 7 |
| F5 GitHub Actions federation facts | PARTIAL | Platform owner | Slice 8 |
| F6 Attended real-browser journey | PENDING | Platform owner | Gate 9 |
| F7 Cross-repository file manifest | PENDING | Platform owner | Slice 7, 8 |
| F8 Existing-flow baseline | PENDING | Platform owner | Slice 4, Gate 9 |
| C1 Windows credential blob limit | ESTABLISHED | Platform owner | Slice 3 |
| C2 Auth0 loopback port policy | RESOLVED (repo) | Platform owner | Slice 4 |
| C3 Pure-Go SQLite cross-process locking | PENDING | Platform owner | Slice 3 |
| D1–D4 Dependency proofs | DRAFT | Platform owner | Slice 3, 4 |

Evidence marked `RESOLVED (repo)` is grounded in the `infra` and `desktop-app`
repositories as read on 2026-08-04; a fact that lives only in the Auth0
dashboard or in continuous-integration secrets is still marked PENDING. Literal
client identifiers and secrets are referenced by `repository:file:line`, never
copied into this repository.

Slice 2 (typed contracts and architecture enforcement) depends on none of the
external facts above; it depends only on the threat model in
[threats](../threats/authentication.md). Slice 3 and later depend on the rows
that name them.

## F1: Current Auth0 tenant state

- **Source**: `infra` (Terraform) and `desktop-app` (`electron/auth.ts`), read
  2026-08-04; [authentication.md](../authentication.md) §3.
- **Environment**: the managed tenant serves prod, staging, and dev1.
- **Result**: RESOLVED for what the repositories hold; dashboard-only objects
  remain PENDING and are named below.

Authoritative tenant: `drizz-dev.eu.auth0.com` (EU region), instantiated per
environment through the single `modules/auth0` module. A second tenant,
`dev-kf5tbjokhczvagbn.eu.auth0.com`, is used by unrelated services and by the
desktop dev configuration; its role is unclear and its audience looks like a
placeholder.

One application is Terraform-managed: `web_app`, `app_type = "regular_web"`,
`client_secret_post`, grants `authorization_code`, `refresh_token`,
`client_credentials`, and passwordless one-time-password, with rotating and
expiring refresh tokens (`infra modules/auth0/main.tf`). The distributed
**desktop** and **toolbase/MCP** applications are referenced by client
identifier only and are managed in the Auth0 dashboard, not in Terraform, so
their exact grants, redirect registrations, and token settings are not in any
repository and are PENDING dashboard inspection.

The existing desktop application is a public native client using Authorization
Code with PKCE `S256` and no client secret, requesting
`openid profile email offline_access`, with per-attempt state and nonce and full
ID-token verification (`desktop-app electron/auth.ts`). It stores refresh tokens
in the operating-system keychain and an encrypted local store. This confirms the
native public-client model the plan adopts.

Connections are not Terraform-managed. Google and email one-time-password are
the observed connections, inferred from the account-linking Action and the
passwordless grant, matching the desktop login surface. No Device Authorization
grant is enabled on any application in the tenant; the plan requires it, so it
is new work in F2.

The plan §4 references an application named Drizz MCP as evidence. The Platform
owner confirms it exists as a separate application created manually in the Auth0
dashboard, not through Terraform, which is why it is in neither repository. Per
the plan §4 it is evidence, not a template, and it is not the default or
approved path; F2 authors a fresh approved Platform native application
regardless of its configuration. It therefore constrains nothing and blocks no
slice. Its exact settings may be recorded here as optional reference for F2 if
the owner supplies them; they are not required to proceed.

## F2: Final Auth0 resource proposal

- **Owner**: Platform owner, authored through reviewed infrastructure code.
- **Result**: PENDING; inputs gathered in F1.

The proposal must define, additively and without altering shared defaults:

- a new Platform native application, distinct from `web_app`, with
  `token_endpoint_auth_method = "none"` and no secret; grants
  `authorization_code`, `refresh_token`, and `device_code` only; implicit,
  password, password-realm, and client credentials disabled; PKCE `S256`
  required. The tenant enables none of Device Authorization today, so that grant
  is new;
- an exact registered loopback redirect on `127.0.0.1` only, with the port
  policy resolved in C2; the plan does not register the `localhost` name that
  `web_app` currently also carries;
- rotating refresh tokens with the access, idle, and maximum lifetimes in
  authentication.md §8, which differ from the single lifetime set on `web_app`;
- a dedicated Platform API and audience distinct from the existing
  `app.drizz.dev`, `toolbase.drizz.dev`, and `auth.drizz.dev` audiences, with
  reviewed scopes, since the current resource server defines no scopes and
  authorization is carried as custom claims by a post-login Action;
- one workload application per continuous-integration integration (see F5);
- a decision on whether the shared post-login Actions apply to the new
  application (the RBAC Action already runs for the desktop client), with the
  regression from authentication.md §10.

## F3: Supported operating-system matrix

- **Source**: authentication.md §8; plan §3.
- **Result**: PROPOSED; owner confirms exact versions and Linux desktop
  environments.

| Platform | Vault | Notes |
| --- | --- | --- |
| macOS | Keychain | Confirm minimum supported release |
| Windows | Credential Manager | Per-user store; blob limit in C1 |
| Linux desktop | Secret Service | Confirm supported desktop environments and a headless fallback story |

If a supported platform lacks safe credential storage, the current process
model fails closed (authentication.md §8) and the ADR 0006 review trigger
applies.

## F4: Drizz Cloud organization endpoint and contract

- **Owner**: Platform owner.
- **Result**: RESOLVED from `drizz-tm-backend` source (2026-08-05).

The internal `/api/v1/rbac/internal/user-context` (x-internal-key) remains the
Auth0-Action path and is not used by the Platform CLI.

Platform-facing contract for Slice 7 (source: `drizz-tm-backend`
`app/api/v1/desktop_router.py`, `endpoints/organization.py`,
the organization resolver in `utils/auth.py`, `schemas/organization.py`):

- **Endpoint**: `GET /api/organizations/me` (desktop bearer router).
- **Auth**: `Authorization: Bearer <access_token>`, Auth0 RS256 verified via JWKS
  for the configured desktop domain; subject from `sub`, organization from the
  namespaced organization claim or an active-membership lookup.
- **200**: `Organization { id: int, name: string, created_at, updated_at,
  created_by, notification_config }`. Role, status, and permissions are not on
  this route; per-operation authorization stays at the cloud (`SEC-029`).
- **Failure**: `403` "access has been revoked" → `AUTHORIZATION_FORBIDDEN`; `403`
  "organization not found, log in again" and `401` (expired, session terminated,
  no subject) → `AUTHENTICATION_REQUIRED`; `5xx` or malformed → `AUTHENTICATION_UNAVAILABLE`.
- Membership status is `active` or `pending` only; removal is non-active, not a
  distinct "suspended" state.

Both deployment facts are now resolved and proven live (2026-08-05):

- **Base URL** (`DRIZZ_CLOUD`): `https://stag.drizz.dev/api/tm/desktop`. The
  desktop bearer routes mount under `/desktop` (`app/api/v1/router.py`) behind
  the `/api/tm` ingress, so the client's `/api/organizations/me` resolves to
  `https://stag.drizz.dev/api/tm/desktop/api/organizations/me`.
- **Desktop-token JWKS domain**: `drizz-dev.eu.auth0.com`, set by
  `DRIZZ_AUTH0_DESKTOP_DOMAIN` in `infra` (`stag/05-compute/main.tf`,
  `prod/05-compute/main.tf`) — the same tenant the CLI uses, so its tokens
  validate. The backend uses signature-only validation for desktop tokens.
- **Live result**: `drizz login` with `DRIZZ_CLOUD` set resolved and displayed
  the caller's organization end to end.

No organization identity is persisted locally in Stage 3.

## F5: GitHub Actions federation facts

- **Environment**: public standard plus Drizz-specific policy.
- **Result**: PARTIAL.

Established public facts:

- issuer `https://token.actions.githubusercontent.com`;
- discovery and JWKS are published by that issuer;
- the OIDC token carries repository, reference, and workflow claims and a unique
  token identifier (`jti`).

Pending Drizz-specific facts (owner: Identity and Cloud):

- the exact audience Drizz requires;
- the exact subject claim shape Drizz accepts;
- assertion lifetime;
- ingress path (protected file or socket, never ordinary configuration);
- exchange authentication and whether the exchange provider offers atomic
  one-time exchange or a bounded replay store is required, including its owner
  and retention.

## F6: Attended real-browser journey

- **Owner**: Identity owner.
- **Result**: PENDING. An attended system-browser sign-in through every
  approved real connection, including multi-factor where configured, is run
  before release. No password, seed, cookie, or browser profile enters
  continuous integration or evidence.

## F7: Cross-repository file manifest

- **Owner**: Cloud and Infrastructure owners.
- **Result**: PENDING. The exact external files for the cloud organization
  resolver, the workload exchange, and any replay store are recorded here
  before any external edit (plan §11 Gate 1 and Slices 7–8).

## F8: Existing-flow baseline

- **Owner**: Web, playground, desktop, backend product owners.
- **Result**: PENDING. Capture the current staging login journeys before any
  shared Auth0 change so regression is measurable (authentication.md §10).

## C1: Windows credential blob limit

- **Source**: Windows Credential Manager public contract
  (`CRED_MAX_CREDENTIAL_BLOB_SIZE`).
- **Environment**: public standard.
- **Result**: ESTABLISHED.

A single Windows credential blob is limited to 2560 bytes. The plan §9 bounds
`credential.Record` at 16 KiB serialized, which is not storable as one Windows
blob.

Resolved in Slice 3: the vault stores one blob per version and caps the encoded
record at 2048 bytes (`infrastructure/vault`), under the 2560 limit. A realistic
record (256-byte refresh token, full-length provider subject, UUID session,
32-character client identity, issuer URL) was measured at 651 bytes encoded, so
the single-blob approach holds with wide margin; no chunking is needed. Oversized
records fail closed rather than truncate.

The real macOS Keychain round-trip is verified: the `go-keyring` adapter writes,
reads back byte-identical, and deletes the record through the actual Keychain
(`tests/identity/vault_test.go`, run with `-tags keyring`). The Linux Secret
Service and Windows Credential Manager round-trips remain for their protected
identity CI jobs; `go-keyring` returns `ErrSetDataTooBig` if a Windows blob is
oversized, and the 2048 cap stays under that limit.

## C2: Auth0 loopback port policy

- **Source**: `infra/prod/target.tfvars:8-13`; `desktop-app electron/auth.ts`.
- **Result**: RESOLVED for the current tenant; the Platform native
  application's own registration is authored in F2.

The tenant registers a fixed loopback callback, `http://127.0.0.1:8490/callback`
and `http://localhost:8490/callback`, in production only, on the confidential
`web_app` client. No ephemeral or wildcard loopback port is registered. The
desktop code attempts an operating-system-assigned port when 8490 is busy, but
the tenant does not register that port, so an arbitrary ephemeral port is not
supported by the current configuration.

Decision for the plan: the Platform native application registers a fixed
`127.0.0.1` loopback port, not an ephemeral one, and does not register the
`localhost` name (the plan binds `127.0.0.1` only to resist DNS rebinding). If
the chosen fixed port is busy at run time, the plan fails safely rather than
falling back to an unregistered port. The exact port and its registration are
part of the F2 proposal.

## C3: Pure-Go SQLite cross-process locking

- **Owner**: Identity owner.
- **Result**: PENDING (proven in Slice 3). The fenced publication and lease
  correctness in the plan §8 rely on `modernc.org/sqlite` honoring cross-process
  locking equivalently to C SQLite in WAL mode. Slice 3 proves this with the
  multi-process, fencing, and uncertain-refresh tests before it is relied upon.

## Security findings during Gate 1

While reading the infrastructure repository for F1, a committed plaintext OAuth
client secret and other third-party keys were found in git-tracked Terraform,
alongside continuous-integration state files in the tree. This is unrelated to
the Platform identity module but was raised to the Platform owner for a separate
remediation in that repository: rotate the affected credentials, remove the
values from tracked files and history, and inject them through the secret store
as the managed web application already does. No secret value was copied into
this repository.

## Dependency evidence

Draft `DEL-002` records for the approved stack (plan §10). Each is completed and
promoted to [dependencies.md](../dependencies.md) when the dependency is first
added, in Slice 3 or 4. Binary-size and performance figures are measured at
that point.

| Module | Purpose | License | Upstream | Replacement boundary | Added in |
| --- | --- | --- | --- | --- | --- |
| `golang.org/x/oauth2` | OAuth Authorization Code, PKCE, Device Code, token exchange | BSD-3-Clause | Go team; active | `infrastructure/auth0` | Slice 4 |
| `github.com/coreos/go-oidc/v3` | OIDC discovery and ID-token validation | Apache-2.0 | CoreOS/Red Hat; active | `infrastructure/auth0` | Slice 4 |
| `github.com/zalando/go-keyring` | macOS, Windows, Linux vault access | MIT | Zalando; active | `infrastructure/vault` | Slice 3 |
| `modernc.org/sqlite` | Pure-Go SQLite, keeps `CGO_ENABLED=0` | BSD-3-Clause | modernc; active | `infrastructure/sqlite` | Slice 3 |

Each record must still cover, at add time: necessity versus the standard
library, transitive risk, `govulncheck` result, supported cross-builds,
measured binary-size delta, and hot-path performance where relevant. No Auth0
Management SDK, ORM, migration framework, authentication framework, or
automatically refreshing general-purpose token source is approved (plan §10).

There is no Auth0 SDK for the authentication flows in Go; the standard native
and headless flows use `x/oauth2` with `go-oidc`, which this stack adopts.

## Gate exit

The [threat model](../threats/authentication.md) is owner accepted with no
unowned high or critical risk, so **Slice 2 may proceed**. C1 and C2 are
resolved from the repositories; F1 is resolved for what the repositories hold.

Still open before the slices that name them: F2 authored (Slice 4, 5, 8); the
dashboard-only desktop, toolbase, and Drizz MCP application details (F1 tail);
F4 Platform-facing organization contract (Slice 7); F5 GitHub facts (Slice 8);
F6, F7, F8 (Slice 4 and Gate 9); C3 proven (Slice 3); and the dependency records
accepted at add time (Slice 3, 4). None of these block Slice 2.
