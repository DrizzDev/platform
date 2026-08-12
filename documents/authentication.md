# Authentication and Authorization

Status: Approved architecture

## 1. Purpose

This document defines authentication and authorization for the new Drizz Platform, including its CLI, local MCP server, future remote MCP endpoint, desktop integration, and automation.

The existing sign-in flows used by Drizz web, playground, desktop, load balancers, and backend services are not changed by this decision. The Platform may add Auth0 applications, APIs, grants, and backend support, but it must not break or silently alter an existing client.

Stage 3 of the [Delivery Roadmap](roadmap.md) owns delivery order. The [Stage 3 Authentication Implementation Plan](plans/authentication.md) owns the implementation slices, inventories, and evidence. This document owns the security and protocol contract; it does not create a competing delivery plan.

## 2. Core model

Drizz uses one identity system, Auth0, with the standard OAuth flow appropriate to each client.

| Use case | Authentication flow | Credential owner |
| --- | --- | --- |
| Person using the installed CLI or host | Authorization Code with PKCE in the system browser | Installed Drizz host |
| Person on SSH or another headless terminal | Device Authorization Grant | Installed Drizz host |
| Local stdio MCP started by Claude, Codex, an IDE, or another client | Reuse the installed user's Drizz session | Installed Drizz host |
| Person connecting a hosted agent to a future remote MCP endpoint | MCP OAuth 2.1 with Authorization Code and PKCE | MCP client |
| CI/CD or another unattended workload | Workload identity federation, with a dedicated machine client as fallback | Workload secret or identity provider |

Authentication proves an identity. Drizz Cloud independently decides whether that identity may perform an action on an organization or resource. A token, local session, selected organization, or client request never replaces that cloud authorization decision.

## 3. Auth0 boundaries

The Platform uses separate Auth0 objects for separate trust models.

As inspected on 3 August 2026, the existing staging application named Drizz MCP is research evidence, not a production template. Its enabled implicit grant, non-rotating refresh tokens, unbounded refresh lifetime, and broad delegated API policy do not meet this architecture. Staging may be rebuilt to the approved configuration; production must not copy that configuration unchanged.

### Platform Native Application

The distributed CLI and installed host are a public native client. They cannot keep a client secret.

Approved grants:

- Authorization Code;
- Refresh Token;
- Device Code.

Implicit, password, password-realm, and Client Credentials grants are disabled for this application. Authorization Code requests always use PKCE with `S256`. Browser callbacks use an exact registered loopback URI on `127.0.0.1` and an ephemeral port where Auth0 configuration permits it. The implementation also validates state, issuer, audience, signature, expiry, and nonce where an ID token is used.

### Platform API

The Platform has a dedicated Auth0 API and audience. Platform access tokens are not Auth0 Management API tokens and are not accepted by unrelated Drizz APIs. The API enables per-application authorization and exposes only reviewed Platform scopes. Public scope names are approved with the capability contract, not invented by an interface adapter.

External Platform clients send only a Platform-audience token to the Platform cloud boundary. That boundary never passes the token to an API with a different audience. It authorizes the operation, then uses the existing service contract through either a narrowly scoped service identity or Auth0 On-Behalf-Of token exchange when downstream user identity is required. The choice is made per existing service without changing its current login flow.

### Workload Applications

Each CI system or autonomous integration receives its own workload identity. Workload identity federation using the CI provider's signed OIDC identity is preferred because it avoids a long-lived Drizz secret. Where federation is not available, a confidential Auth0 machine-to-machine application may use Client Credentials. It receives only the scopes and organization access required for that integration.

A human refresh token is never copied into CI, and one machine credential is never shared across organizations or unrelated integrations.

### Remote MCP Clients

A future remote HTTP MCP endpoint is an OAuth resource server separate from the local stdio server. Auth0 remains the authorization server.

Supported client registration follows this order:

1. pre-register supported clients where Drizz and the client have an existing integration;
2. support Client ID Metadata Documents for standards-based clients;
3. use Dynamic Client Registration only as a compatibility fallback.

Dynamic Client Registration is not enabled until tenant-wide default grants, connections, consent, redirect validation, rate limits, and access to every existing API have been reviewed. A dynamically registered third-party client must receive no unrelated API access by default.

## 4. Interactive local sign-in

The normal journey is:

1. The user runs the Drizz sign-in command or starts a Drizz surface that needs authentication.
2. Drizz creates a PKCE verifier, challenge, state, and nonce, then starts a loopback callback listener.
3. Drizz opens the system browser at Auth0.
4. Auth0 performs the existing enabled identity experience, such as Google, passwordless email, or another approved connection.
5. Auth0 redirects the authorization code to the loopback listener.
6. Drizz validates the response and exchanges the code with the PKCE verifier.
7. Drizz keeps the short-lived access token in memory and stores renewable credentials in the operating-system credential store.
8. Drizz obtains the user's allowed organization context from a trusted Drizz cloud boundary.

The CLI never asks the user to paste their Auth0 password or a bearer token. The browser remains responsible for credentials, MFA, enterprise connection, and consent.

If the machine cannot complete a browser callback, the user explicitly selects the Device Authorization flow. Drizz displays the verification URI and user code, polls within the server-provided interval, respects expiry and slowdown responses, and can be cancelled. Device Authorization is a fallback, not the default desktop login.

## 5. Local MCP behavior

For local MCP, the external MCP client starts `drizz mcp` over standard input and output when Drizz tools are needed. This process does not run the HTTP MCP OAuth flow.

The process asks the application-owned identity service for the current local Drizz session. That service reads renewable credentials from the current operating-system user's credential store and refreshes them when required. The MCP client, agent, model, prompt, tool input, process arguments, environment variables, and MCP protocol messages never receive the credential.

If no usable session exists, the tool returns a stable authentication-required result with the normal Drizz sign-in action. It does not open repeated browsers inside tool calls. All local Drizz surfaces reuse the same session lifecycle, but they do not duplicate token-handling code.

This follows the MCP guidance that stdio should not use the HTTP authorization protocol. Drizz deliberately uses the operating-system credential store instead of placing a bearer token in an environment variable.

### Agent plugins and hooks

The Claude, Codex, and other approved local integrations package MCP and hook configuration but contain no Drizz credential. A hook invokes the installed Drizz application with a structured host event. The application obtains any credential needed for synchronization only through the same identity service used by CLI, MCP, and desktop flows.

Installing a Drizz integration does not replace, read, or modify the user's Claude, Codex, ChatGPT, Gemini, or IDE login. Agent authentication and Drizz authentication remain separate. Hook input, plugin metadata, process arguments, environment variables, model context, and MCP messages never carry the Drizz access or refresh token.

The integration installer resolves and verifies the installed Drizz executable, merges supported host configuration without overwriting unrelated settings, and uses the host's native review or trust flow. The agent-integration and capture contract is defined in [Agent Integration and Execution Capture](capture.md).

## 6. Remote MCP behavior

The remote MCP model applies only when Drizz offers a hosted HTTP MCP endpoint. It does not make a hosted agent capable of directly controlling a device on the user's computer.

The remote endpoint:

- publishes OAuth Protected Resource Metadata;
- uses Auth0 authorization-server or OpenID Connect discovery;
- requires Authorization Code with PKCE for delegated user access;
- requires the OAuth `resource` parameter for the canonical MCP resource;
- validates issuer, signature, expiry, intended audience, and required scope;
- receives bearer tokens only in the `Authorization` header over TLS;
- returns `401` for missing, invalid, or expired credentials;
- returns `403` with an insufficient-scope challenge when identity is valid but the granted scope is insufficient;
- never accepts or forwards a token issued for another resource.

A future hosted agent can reach a local device only through a separately designed installed-host relay. Such a relay must use an authenticated outbound connection, explicit device ownership, bounded commands, and its own threat model. It is not part of local MCP or approved by this document.

## 7. Organization authorization

Auth0 authenticates the person or workload. Drizz Cloud is authoritative for:

- current organization membership;
- role and entitlement;
- requested action;
- resource ownership and state;
- suspension, revocation, and policy changes.

The client may send an organization or resource identifier as request data, but the cloud resolves it under the authenticated subject and verifies access at the point of use. Organization identifiers, roles, and permissions from MCP arguments, CLI flags, local files, or cached state are never trusted authority.

Local capability policy may restrict device operations before execution. It cannot grant access to cloud-owned data. Offline behavior is limited to explicitly approved local operations and locally owned work; it never creates permanent offline access to cloud resources.

## 8. Credential lifecycle

The starting policy is:

- access token lifetime: 15 minutes;
- rotating refresh tokens for public native clients;
- refresh-token idle lifetime: 90 days;
- refresh-token maximum lifetime: 1 year;
- immediate local removal on logout, followed by server revocation where supported;
- session recovery through a new login when refresh is expired, revoked, reused, corrupted, or denied.

The operating-system credential store is Keychain on macOS, Credential Manager on Windows, and an approved Secret Service implementation on supported Linux desktops. If secure persistent storage is unavailable, the current CLI and local MCP process model fails closed because a memory-only login would disappear when the login process exits. A future persistent desktop-owned process may propose a memory-only session through a separate reviewed design. Drizz never falls back to a plaintext token file.

Access tokens, refresh tokens, authorization codes, device codes, client secrets, and token-bearing URLs are excluded from logs, telemetry, crash data, support bundles, command arguments, project files, local execution records, and model context. CI credentials live only in the CI provider's protected identity or secret store.

## 9. Failure behavior

| Condition | Required behavior |
| --- | --- |
| No local session | Return authentication required and the sign-in action |
| Access token expired | Refresh once through the shared identity service |
| Refresh denied, expired, or reused | Remove unusable local state and require sign-in |
| User cancels browser or device flow | Stop cleanly without storing partial credentials |
| Callback state, issuer, audience, nonce, or signature is invalid | Reject the response and store nothing |
| Organization access changed | Reject the cloud operation and refresh organization context |
| Workload identity is invalid | Fail without falling back to a human session |
| Valid identity lacks permission | Return a stable forbidden result without hiding it as authentication failure |
| Credential store is unavailable | Report a safe actionable error without exposing credential material |

Authentication retries are bounded. A rejected operation is not executed and is not silently retried under another identity or organization.

## 10. Compatibility with existing Drizz products

This decision adds a Platform authentication boundary. It does not replace, rewrite, or reconfigure existing application login behavior.

Before production rollout, Auth0 and backend changes are tested against the existing web, playground, desktop, load-balancer, and service flows. New Platform applications and audiences use distinct identifiers and callbacks. Any shared Auth0 Action or tenant policy change requires explicit regression testing for every affected existing application.

## 11. Verification

The authentication foundation is complete only when tests prove:

- PKCE login, cancellation, callback validation, refresh rotation, restart, logout, and revocation on every supported operating system;
- the Device Authorization fallback on a real headless terminal;
- one local session reused by CLI, local MCP, and the desktop test boundary;
- credentials absent from MCP traffic, process arguments, environment, logs, telemetry, crash data, and local records;
- expired, revoked, replayed, wrong-issuer, wrong-audience, and wrong-resource tokens are rejected;
- cross-organization requests fail at the cloud boundary;
- CI federation and the machine-client fallback cannot exceed their assigned scopes or organization;
- remote MCP discovery, registration, authorization, scope challenge, and token validation work with each supported hosted client before that endpoint is listed as supported;
- existing Drizz login journeys continue to pass unchanged.

## 12. Standards and sources

- [OAuth 2.0 for Native Apps](https://www.rfc-editor.org/rfc/rfc8252)
- [OAuth 2.0 Device Authorization Grant](https://www.rfc-editor.org/rfc/rfc8628)
- [OAuth 2.0 Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [OAuth 2.0 Resource Indicators](https://www.rfc-editor.org/rfc/rfc8707)
- [OAuth 2.0 Protected Resource Metadata](https://www.rfc-editor.org/rfc/rfc9728)
- [OAuth 2.0 Token Exchange](https://www.rfc-editor.org/rfc/rfc8693)
- [MCP authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization)
- [Auth0 MCP authorization](https://auth0.com/ai/docs/mcp/intro/overview)
- [Auth0 MCP client registration](https://auth0.com/ai/docs/mcp/guides/registering-your-mcp-client-application)
- [Auth0 Dynamic Client Registration](https://auth0.com/docs/get-started/applications/dynamic-client-registration)
- [Auth0 On-Behalf-Of Token Exchange](https://auth0.com/docs/secure/call-apis-on-users-behalf/on-behalf-of-token-exchange)
