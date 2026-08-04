# ADR 0006: OAuth flows by client type with cloud authorization

- Status: Accepted
- Date: 2026-08-03
- Related: [Authentication and authorization](../authentication.md),
  [Platform architecture](../architecture.md)

## Context

The Platform serves people, local MCP clients, hosted MCP clients, desktop
surfaces, and unattended workloads. These clients have different abilities to
protect credentials and interact with a browser. Existing Drizz applications
already use Auth0 and must continue to work unchanged.

## Decision

Use Auth0 as the common identity system with a separate Platform native
application, Platform API, workload applications, and remote-client
registrations.

Use Authorization Code with PKCE as the default installed-user flow, Device
Authorization as the explicit headless fallback, workload identity federation
as the preferred CI flow, and Client Credentials as the machine fallback.
Local stdio MCP reuses the installed Drizz session without exposing credentials
to the MCP client. A future remote HTTP MCP endpoint follows the MCP OAuth 2.1
authorization specification.

Drizz Cloud remains authoritative for organization and resource authorization.
No existing application flow is replaced by this decision.

## Consequences

Users receive familiar browser sign-in and agents do not handle Drizz tokens.
CI does not depend on a person's session. Remote MCP clients can use standard
discovery and delegated authorization. Drizz must operate several isolated
Auth0 client types, validate tokens for their exact audience, protect renewable
credentials, and regression-test shared Auth0 tenant changes.

## Validation

Complete the verification matrix in the authentication document, including
real supported clients, operating-system credential stores, CI federation,
cross-organization denial, and regression tests for existing Drizz login flows.

## Review trigger

Revisit if Auth0 cannot support a required MCP registration or resource-binding
standard, a supported platform lacks safe credential storage, or a new client
type cannot use one of the approved flows.
