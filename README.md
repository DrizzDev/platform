# Drizz Platform

Drizz Platform is the local execution and integration foundation that allows AI
agents, developer tools, and Drizz applications to use the same Drizz
capabilities.

This repository contains the approved design and engineering baseline plus the
shared Go foundation for the Drizz CLI, local MCP server, and future interfaces.
Product capabilities are added only after their contracts and boundaries are
approved.

## Intended ownership

- One installable Drizz application.
- Shared application behavior for MCP, CLI, desktop, and future approved
  interfaces.
- Fast local device interaction.
- Durable local execution records and synchronization with Drizz services.
- Explicit integration boundaries for existing Drizz systems.

External agents own planning. Drizz executes approved capabilities, enforces its
boundaries, and records the resulting work.

## Documents

1. [Architecture](documents/architecture.md)
2. [Technology stack](documents/stack.md)
3. [Engineering guide](documents/engineering.md)
4. [Delivery roadmap](documents/roadmap.md)
5. [Decision records](documents/decisions/README.md)
6. [Dependencies](documents/dependencies.md)
7. [Standards exceptions](documents/exceptions.md)
8. [Research](documents/research.md)
9. [Agent instructions](AGENTS.md)

## Current status

| Area | Status |
| --- | --- |
| Shared runtime foundation | Implemented |
| Product capabilities | Not started |
| Architecture | Approved baseline |
| Technology stack | Approved recommendation baseline |
| Engineering standards | Approved and mandatory |
| Delivery roadmap | Approved baseline |
| Decision records | Accepted |

## Commands

The foundation exposes two commands:

- `drizz version` prints the released application identity.
- `drizz mcp` runs the local MCP server over standard input and output. An MCP
  client starts and stops this process; protocol messages use standard output
  and diagnostics use standard error. Incoming messages are size-bounded.

The CLI and MCP adapters consume one application-owned value, the released
identity, rather than reading build information independently. This proves the
reuse boundary; product capabilities and public MCP tools are not yet defined.
`drizz version` performs no configuration, telemetry, or reporting setup; those
are constructed only when `drizz mcp` runs.

## Development

The repository pins Go and all development tools in `go.mod`.

```bash
python3 -m pip install -r requirements.txt
make hook
make fix
make verify
```

The pinned Python dependency installs pre-commit. `make hook` installs the
commit and push hooks. Commits run the fast deterministic checks. Pushes and
pull requests run the complete `make verify` merge gate, including the scripts
under `scripts`.

Repository administrators must configure the GitHub `verify` job as a required
status check. The workflow runs on every pull request, but GitHub branch or
repository rules are what prevent a failed check from being merged.
