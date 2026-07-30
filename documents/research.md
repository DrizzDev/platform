# Research

Status: Approved supporting evidence

Reference implementations inform design choices but do not define Drizz
behavior, architecture, naming, or security policy.

## Reference implementations

| Repository | Relevant area |
| --- | --- |
| [GitHub MCP Server](https://github.com/github/github-mcp-server) | Go MCP composition, local and remote transports, tool grouping, release packaging |
| [Terraform MCP Server](https://github.com/hashicorp/terraform-mcp-server) | Go MCP composition, tool controls, telemetry |
| [Docker MCP Gateway](https://github.com/docker/mcp-gateway) | MCP aggregation, configuration, credentials, local state |
| [Playwright MCP](https://github.com/microsoft/playwright-mcp) | Stateful local provider behavior and MCP client integration |
| [Sentry MCP](https://github.com/getsentry/sentry-mcp) | Shared behavior across local and remote transports |
| [MCP reference servers](https://github.com/modelcontextprotocol/servers) | Protocol examples, schemas, roots, and filesystem safety |

The recurring pattern is a transport-neutral behavior layer with thin MCP
registration. For local standard-input/output integrations, the MCP client
starts and supervises the configured process. Remote HTTP is a separate
deployment and authentication model.

No reference source is copied into this repository without a separate license,
security, maintenance, and fit review.

## Primary sources

- [Go release policy](https://go.dev/doc/devel/release)
- [MCP SDKs](https://modelcontextprotocol.io/docs/sdk)
- [Official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- [OAuth for Native Apps](https://www.rfc-editor.org/rfc/rfc8252)
- [OAuth Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [OpenID Connect Core](https://openid.net/specs/openid-connect-core-1_0-18.html)
- [SQLite write-ahead logging](https://www.sqlite.org/wal.html)
- [OpenTelemetry for Go](https://opentelemetry.io/docs/languages/go/)

## Evidence requirements

A design claim derived from source code must record:

- the public repository URL;
- the exact commit or release;
- the files inspected;
- the observed behavior;
- the limitation of the comparison.

Personal filesystem paths, private workstation layout, and uncommitted local
clones must never appear as shared evidence.
