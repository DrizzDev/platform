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
| [Arize coding harness tracing](https://github.com/Arize-ai/coding-harness-tracing) | Claude, Codex, Gemini, Cursor, and Copilot hook mapping and trace reconstruction |
| [OpenInference](https://github.com/Arize-ai/openinference) | Standard AI trace vocabulary and SDK instrumentation across supported languages |
| [Neatlogs](https://github.com/neatlogs/neatlogs) | Provider SDK wrapping, exposed-thinking capture, tool correlation, and OTLP export |

The recurring pattern is a transport-neutral behavior layer with thin MCP
registration. For local standard-input/output integrations, the MCP client
starts and supervises the configured process. Remote HTTP is a separate
deployment and authentication model.

No reference source is copied into this repository without a separate license,
security, maintenance, and fit review.

The agent-observability references do not provide the Drizz product boundary.
The primary Drizz journey uses official host integrations and a Drizz-owned
execution record. `coding-harness-tracing` and Neatlogs are not installed on
customer machines. OpenInference is deferred to a future SDK integration and
would remain behind an adapter.

OpenInference currently provides Go semantic conventions, shared
instrumentation utilities, and OpenAI and Anthropic Go SDK instrumentors. It
does not provide Go adapters for Claude Code, Codex, Gemini CLI, or their
desktop applications. Those external hosts require their official plugin,
hook, transcript, or structured-event surfaces regardless of language.

## Primary sources

- [Go release policy](https://go.dev/doc/devel/release)
- [MCP SDKs](https://modelcontextprotocol.io/docs/sdk)
- [Official Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Claude Code hooks](https://code.claude.com/docs/en/hooks)
- [Claude Code plugins](https://code.claude.com/docs/en/plugins-reference)
- [Codex hooks](https://developers.openai.com/codex/hooks)
- [Codex plugins](https://developers.openai.com/plugins/build/plugins)
- [Gemini CLI hooks](https://geminicli.com/docs/hooks/reference/)
- [OAuth for Native Apps](https://www.rfc-editor.org/rfc/rfc8252)
- [OAuth Security Best Current Practice](https://www.rfc-editor.org/rfc/rfc9700)
- [OAuth Device Authorization Grant](https://www.rfc-editor.org/rfc/rfc8628)
- [OAuth Resource Indicators](https://www.rfc-editor.org/rfc/rfc8707)
- [OAuth Protected Resource Metadata](https://www.rfc-editor.org/rfc/rfc9728)
- [OAuth Token Exchange](https://www.rfc-editor.org/rfc/rfc8693)
- [OpenID Connect Core](https://openid.net/specs/openid-connect-core-1_0-18.html)
- [MCP authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization)
- [Auth0 MCP authorization](https://auth0.com/ai/docs/mcp/intro/overview)
- [Auth0 MCP client registration](https://auth0.com/ai/docs/mcp/guides/registering-your-mcp-client-application)
- [Auth0 On-Behalf-Of Token Exchange](https://auth0.com/docs/secure/call-apis-on-users-behalf/on-behalf-of-token-exchange)
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
