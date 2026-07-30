# Platform Architecture

Status: Approved baseline

## 1. Purpose

Drizz Platform provides one local foundation through which MCP clients, the
command line, and the Drizz desktop application can use approved Drizz
capabilities.

External agents such as Claude, Codex, and Gemini interpret the user's goal and
plan the work. Drizz does not add another planner. Drizz authenticates the user,
executes approved operations, interacts with local devices, records outcomes,
and synchronizes required data with Drizz services.

## 2. Known product constraints

- Device observation and actions should run locally when the device is local.
- MCP and CLI must reuse the same product behavior.
- The desktop application must be able to use the same behavior without a
  second implementation.
- Drizz services remain authoritative for cloud data and organization access.
- Local work must be recorded and synchronized without making every device
  action wait for a cloud round trip.
- Public tool names and schemas are not yet approved.

## 3. System context

The following diagram shows runtime communication. Its arrows mean “invokes” or
“exchanges data with”; they do not represent source-code dependencies.

```mermaid
flowchart TB
    user["User"]

    subgraph clients["User-facing clients"]
        agent["External agent<br/>Claude, Codex, Gemini, IDE"]
        terminal["Terminal or automation"]
        desktop["Drizz desktop"]
    end

    subgraph local["Local computer"]
        mcp["Drizz MCP interface"]
        cli["Drizz CLI interface"]
        host["Drizz application"]
        device["Android or iOS device integration"]
        records["Local records and artifacts"]
    end

    cloud["Drizz services"]

    user --> agent
    user --> terminal
    user --> desktop
    agent --> mcp
    terminal --> cli
    desktop --> host
    mcp --> host
    cli --> host
    host --> device
    host --> records
    host <--> cloud
```

The desktop-to-host protocol is intentionally unspecified until the current
desktop application is inspected and the lifecycle requirements are known.

## 4. Local and cloud responsibility

### Local responsibility

- start the requested interface;
- obtain and protect the user's authenticated session;
- validate local requests;
- interact with local devices;
- run approved local deterministic behavior;
- record local execution facts and artifacts;
- retry and resume synchronization;
- bound disk, memory, process, and concurrency usage.

### Cloud responsibility

- validate cloud credentials;
- enforce organization membership and resource authorization;
- own durable cloud resources and shared history;
- accept synchronized records and artifacts idempotently;
- expose cloud-backed product operations.

Local cached identity can support explicitly approved offline behavior. It never
grants access to cloud data and never replaces cloud authorization.

## 5. Architectural shape

The approved shape is a modular monolith with inward dependencies. A single
installable application is appropriate for local execution, while module and
port boundaries preserve the option to extract a capability later.

The following diagram shows source-code dependency direction. Every arrow means
“depends on.”

```mermaid
flowchart TB
    subgraph outer["Outer adapters"]
        mcpAdapter["MCP adapter"]
        cliAdapter["CLI adapter"]
        desktopAdapter["Desktop adapter"]
        deviceAdapter["Device adapter"]
        cloudAdapter["Cloud adapter"]
        storageAdapter["Storage adapter"]
    end

    subgraph core["Application core"]
        application["Application use cases and ports"]
        domain["Domain rules and values"]
    end

    mcpAdapter --> application
    cliAdapter --> application
    desktopAdapter --> application
    deviceAdapter --> application
    cloudAdapter --> application
    storageAdapter --> application
    application --> domain
```

The domain does not import adapters, SDKs, databases, transports, or cloud
clients. Application code owns the interfaces that outer adapters implement.

## 6. Capability reuse

A product capability is an application use case. MCP tools, CLI commands, and
desktop actions are interface-specific mappings to that use case.

For example, a future approved device tap operation would be implemented once:

```mermaid
flowchart TB
    subgraph interfaces["Interfaces"]
        mcpTool["MCP tool"]
        cliCommand["CLI command"]
        desktopAction["Desktop action"]
    end

    tap["Tap use case"]
    devicePort["Device port"]
    deviceAdapter["Device adapter"]
    platform["Android or iOS integration"]

    mcpTool --> tap
    cliCommand --> tap
    desktopAction --> tap
    tap --> devicePort
    deviceAdapter -.->|implements| devicePort
    deviceAdapter --> platform
```

This diagram deliberately separates the use case from the interface name.
Inputs, outputs, authorization, recording, and failure behavior must be agreed
before the operation becomes public.

## 7. Process model

### MCP

For a local standard-input/output integration, the MCP client starts
`drizz mcp` when it needs the configured server. The client owns that child
process lifecycle. Drizz writes protocol messages to standard output and
diagnostics to standard error.

### CLI

`drizz <command>` starts a normal command-line process, invokes the same
application use case, returns a stable exit result, and exits.

### Desktop

The desktop application should use the same released Drizz implementation. The
choice between a supervised child process and another local boundary remains
open until the desktop source and multi-client lifecycle are audited. The
desktop must not reimplement capabilities.

### Remote interfaces

A remote MCP or HTTP deployment is a separate product mode. It has different
authentication, authorization, availability, and local-device constraints and
is not implied by the local design.

## 8. Candidate product modules

These are candidate ownership boundaries, not approved public tools or a
required folder tree.

| Module | Responsibility |
| --- | --- |
| Identity | Local sign-in lifecycle and authenticated context |
| Device | Device discovery, connection, observation, and actions |
| Execution | Operation lifecycle, cancellation, and outcome |
| Artifact | Local artifact identity, integrity, and retention state |
| Sync | Durable transfer, retry, resume, and reconciliation |
| Authoring | Approved deterministic authoring behavior |

Test plan, debugging, application, and report modules should be designed when
their product requirements and existing cloud contracts are inspected.

## 9. Local records and synchronization

The requirements are:

- an operation and its synchronization intent cannot silently diverge;
- retries must not create duplicate cloud effects;
- restart, sleep, loss of network, and process termination must be recoverable;
- credentials must not be stored with normal execution data;
- local artifacts must have bounded retention and safe cleanup;
- in-use or unsynchronized evidence must not be deleted.

SQLite metadata plus filesystem artifacts is the accepted architectural
direction. Its driver, schema, locking, crash, migration, multi-process,
cleanup, and supported-operating-system behavior must be proven before the
implementation is accepted.

## 10. Security boundaries

- Native sign-in should use the system browser and PKCE if supported by the
  current Drizz identity service.
- Credentials belong in the operating-system credential store.
- The cloud rechecks authorization for every cloud resource.
- Interface input cannot supply trusted user or organization identity.
- Device data, screenshots, page structure, logs, and artifacts are potentially
  sensitive.
- Subprocesses require explicit executable, arguments, environment, deadline,
  output, and shutdown limits.
- Local files must remain within approved Drizz-owned locations unless the user
  explicitly grants another path.

These requirements do not define the final authentication contract. That
requires an audit of the existing Drizz identity and authorization flow.

## 11. Open decisions

- approved first capability set and public vocabulary;
- current authentication endpoints, token model, scopes, and organization
  context;
- desktop process and update lifecycle;
- device integration boundary for Android and iOS;
- which deterministic authoring behavior may be distributed locally;
- local persistence and cleanup design;
- synchronization contract with existing Drizz services;
- supported operating systems and installation channels;
- whether any remote MCP or HTTP mode is a product requirement.

No open decision in this section should be treated as an implementation choice.
