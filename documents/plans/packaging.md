# Packaging and Distribution Plan

Status: Draft for owner review

## What this delivers

A person installs **one thing** — the Drizz application — with one standard command, and everything Drizz needs to drive a device is inside it. There is no second download, no separate runtime to install, and no manual wiring of a helper. This plan describes how the two programs Drizz is built from are combined into that single installable, how the combined application reaches a device at runtime, and how it is delivered, updated, and removed on each supported operating system. It is the first version shared with the team (roadmap Stage 7).

This is delivery and packaging work. It does not change any product behaviour, any command, or any architecture decision. It takes the parts that already work and makes them shippable as one artifact.

## Terms used in this plan

The **platform binary** is the Drizz program itself: a single compiled Go executable that carries the command line, the local MCP server, sign-in, the local record store, and the integration installer. It is what a person runs as `drizz`.

The **device helper** (also called the **bridge**) is a separate program, written in TypeScript, that performs the actual device mechanics — talking to `adb` on Android and to `simctl` and WebDriverAgent on iOS. Drizz reuses it rather than reimplementing device control in Go.

To **embed** a file means to bake its bytes into the platform binary at build time, so the finished executable literally contains a copy of the file. Go provides this through the standard `embed` package.

A **digest** is a SHA-256 hash of the helper's bytes. Drizz pins the expected digest and refuses to run a helper whose bytes do not match, so a tampered or stale helper cannot be executed.

**Notarization** is Apple's process of uploading a signed binary to Apple, which scans it and returns a ticket the binary carries, so macOS Gatekeeper runs it without a warning. **Code signing** proves a binary is from a known developer and unaltered; notarization is the extra Apple scan on top.

## What already exists

Two programs, each working on its own:

- The **platform binary** builds today for every supported target — macOS, Linux, and Windows on both `amd64` and `arm64` — through the existing crossbuild step. It contains the command line, the MCP server started by `drizz mcp`, sign-in, the local record store, and the installer that wires agent applications to Drizz.
- The **device helper** builds today as CommonJS JavaScript (`dist/serve.js`) and speaks a small stdio JSON-RPC protocol. It fully manages the *on-device* components it needs: on Android it installs Drizz's own instrumentation onto the phone (`adb install -r`) and runs an on-device gRPC server; on iOS it builds and code-signs WebDriverAgent through `xcodebuild` and provisions it. For the *host* toolchains it depends on — `adb` and the Android build-tools, or Xcode and its iOS platform SDK — it detects absence and returns a precise, actionable message telling the person exactly what to install, rather than installing those itself.

The seam between them already exists and is supervised: the platform binary launches the helper as a long-lived child process, keeps one running for the life of a connection, restarts it with backoff if it dies, and pins its digest before running it. Today the platform binary finds the helper through two environment variables — a path and a digest — which the helper resolver documents as being **for tests and continuous integration only**. Production is meant to resolve the helper from a copy carried inside the application. That carried-inside copy is what this plan adds.

Two things are therefore missing, and only these two:

1. The helper is not yet compiled into a self-contained executable, and is not yet embedded in the platform binary.
2. There is no signed, published installer for a person to download with one standard command.

## The packaging model

The finished product is a single platform binary per operating system and architecture, with the matching device helper embedded inside it. Installing Drizz installs the helper. The person never sees the helper as a separate file.

The helper is prepared for embedding in one step: it is compiled from TypeScript into a **single self-contained native executable** — one that bundles its own JavaScript runtime, so the person needs no separately installed Node. This produces one helper executable per target (for example `darwin/arm64`, `darwin/amd64`, `linux/amd64`, `windows/amd64`). The crossbuild step then embeds the helper that matches each platform binary it builds, so every published binary carries exactly the helper for its own operating system and architecture, and its pinned digest.

At first use, the platform binary makes the embedded helper real: it writes the helper to a protected per-user location, verifies the bytes against the pinned digest, marks it executable, and runs it — then reuses that same copy for the life of the connection. A later Drizz version ships a newer helper with a new digest; on first run it writes and verifies the new copy. The environment-variable override stays, unchanged, as the escape hatch for development and continuous integration.

## How the combined application reaches a device at runtime

The end-to-end path, from a person to a device, once Drizz is installed and connected to an agent:

1. The person installs the platform binary with one command and signs in with `drizz login`.
2. The person runs `drizz connect enable <agent>`, which writes a Drizz entry into that agent's configuration — the command being the installed platform binary's own path, with the argument `mcp`. Optionally `drizz connect capture` adds the turn-event hooks.
3. The person restarts their agent application. On start it reads its configuration and launches `drizz mcp` as a child process, speaking MCP over standard input and output.
4. The agent asks for a device action, calling a Drizz MCP tool.
5. On the first such call, the platform binary extracts, verifies, and runs the embedded helper, then talks to it over stdio JSON-RPC.
6. The helper performs the action on the device through `adb` or `simctl`, provisioning any on-device component it needs.
7. Drizz writes a durable local record of the action, then returns the result to the agent.

Nothing in this path requires a network connection to Drizz's servers; uploading records to the cloud is a separate, later step.

## Why a local install is required at all

Drizz controls a device that is physically connected to the person's own computer. A remote, link-only MCP server — the kind a person adds by pasting a URL — cannot reach a device on someone's desk. Drizz must therefore run locally, beside `adb`, which is why the product is a locally installed binary rather than a hosted link. A future cloud-data capability such as reporting, which needs no local device, may instead be a hosted link-based MCP server; that is a later stage, not this one.

## Host prerequisites

The embedded helper removes every Drizz-specific installation step. It does not remove the platform toolchains a person needs to talk to their own devices:

- **Android:** `adb` and the Android build-tools must be present on the machine.
- **iOS:** Xcode, its command-line tools, and the iOS platform SDK must be present.

The helper already detects the absence of each and returns an exact instruction for installing it. This plan does not bundle those toolchains; it treats them as declared prerequisites with clear, actionable guidance, which is the correct behaviour — Xcode in particular cannot be installed silently. Bundling `adb` itself may be revisited later as a convenience, but it is not required for a working product.

## Distribution channels

Standard install experiences a person expects for a command-line tool, all produced by one release tool from the existing crossbuild:

- **Homebrew** (macOS and Linux): `brew install` from a Drizz tap — the default for command-line tools; it handles the binary, the path, and updates, and sidesteps the Gatekeeper prompt.
- **Shell installer**: `curl -fsSL https://get.drizz.dev | sh` — detects the operating system and architecture and places the right binary on the path.
- **Windows**: a Scoop or winget entry.

No `.dmg` or double-click bundle is used for the command line. Because the command line is delivered through Homebrew and the shell installer, **notarization is largely unnecessary for the first version** and is deferred until a channel needs it. Signing is set up so it can be enabled without reshaping the pipeline.

Lifecycle is defined explicitly: update, rollback, downgrade, and uninstall behaviour, including what happens to a person's local records and credentials on removal, and the version-compatibility rule between a desktop agent and the platform binary.

## The embed-and-run seam (design sketch)

The seam is small and lives entirely in the device infrastructure, beside the existing helper resolver. It is sketched here so it is ready to implement; it is not yet code in the tree.

The build embeds the matching helper and its digest:

```go
// internal/device/infrastructure/carrier/carrier.go  (sketch)
package carrier

import _ "embed"

//go:embed helper/drizz-bridge
var helper []byte

//go:embed helper/drizz-bridge.sha256
var digest string
```

Crossbuild writes `helper/drizz-bridge` (the compiled helper for the target) and its digest into that package before `go build`, so each platform binary embeds its own helper.

At first use the carrier makes the helper real, verifying before running:

```go
// Materialize writes the embedded helper to a protected per-user path once, verifies its bytes against
// the pinned digest, and returns the path to run. A digest mismatch is refused, never run.
func (carrier Carrier) Materialize() (string, error) {
    target := carrier.path()                       // protected per-user dir, e.g. under os.UserCacheDir()
    if carrier.verified(target) {                  // already extracted and matches the pinned digest
        return target, nil
    }
    sum := sha256.Sum256(helper)
    if hex.EncodeToString(sum[:]) != strings.TrimSpace(digest) {
        return "", errors.New("embedded device helper failed its integrity check")
    }
    // atomic publish: write to a temp file, chmod +x, rename over the target
    if failure := carrier.publish(target, helper); failure != nil {
        return "", failure
    }
    return target, nil
}
```

The existing resolver changes only in its fallback order: an explicit environment override still wins for development and continuous integration; otherwise the resolver calls `Materialize` and pins the digest it verified. The supervision, restart, and stdio JSON-RPC layers that already exist are untouched — they simply receive a path that now comes from inside the application instead of from an environment variable.

## Delivery slices

The stage is delivered as a sequence of self-contained slices, each built, reviewed, and proven before the next begins.

1. **Compile the helper for embedding.** Add a build step in the device-helper repository that produces a single self-contained native executable per target and writes its SHA-256 digest beside it. The compilation approach is decided by evidence in this slice: try `bun build --compile` first (it bundles its own runtime, so a person needs no separately installed Node), and fall back to Node's Single Executable Application only if one of the helper's dependencies — the gRPC client, process spawning, or the WebDriverAgent assets — will not bundle cleanly. The chosen approach and why is reported before the slice is committed. No platform-binary change yet; the output is validated by running the compiled helper directly over stdio.
2. **Carry the helper inside the platform binary.** Add the `carrier` package that embeds the helper and its digest, and its `Materialize` operation. Change only the helper resolver's fallback order: the environment override still wins; otherwise it calls `Materialize`. Supervision, restart, and the JSON-RPC seam are untouched.
3. **Wire crossbuild.** Before each target's `go build`, place the matching compiled helper and digest into the `carrier` package's embedded paths, so every published binary carries exactly its own helper.
4. **Standard distribution.** Add the release configuration (one tool) that produces the Homebrew tap, the shell installer, the Windows channel, archives, and checksums, with signing hooks in place; notarization left disabled until a channel requires it.
5. **Release qualification.** Install through `brew` and the shell installer on the supported operating-system matrix, connect Claude and Codex, invoke a capability, confirm capture — plus the live device pass — then qualify the first release.

The environment-variable override, the supervision layer, and the digest-pinning check already exist; slices 2 and 3 only change where the helper's bytes come from.

## Design inventory

| Field | Decision |
| --- | --- |
| Capability | One installable platform binary that carries the device helper inside it and makes it run on first use, delivered through standard install channels (Homebrew, shell installer, Windows). |
| Owner | Device infrastructure owns the embedded-helper carrier and the resolver fallback; the build and release scripts own compilation, embedding, and distribution. No domain or application code changes. |
| Layer | New `internal/device/infrastructure/carrier` (embed + materialize + verify); the existing `sidecar` resolver's fallback; `scripts/crossbuild`; a release configuration; a helper-compile step in the device-helper repository. |
| Contract | `Carrier.Materialize() (string, error)` returns the path of a digest-verified, executable helper, or refuses. The resolver returns the environment override when set, otherwise the materialized path, and pins the digest it verified. |
| Dependencies | Inward only. The carrier is infrastructure and uses only the standard library (`embed`, `crypto/sha256`, `os`, `path/filepath`). No new third-party dependency in the platform binary; the helper's runtime is bundled into the helper executable, not linked into Go. |
| State | The extracted helper lives at a protected per-user path (under the user cache directory); it is written once and reused, re-extracted only when the pinned digest changes. |
| Failures | A digest mismatch is refused and never run. An extraction or permission failure surfaces as the existing device-unprepared refusal, so a broken carrier degrades exactly like a missing helper — it never crashes the surface. |
| Files | device-helper repository: compile target + digest output. Platform: `internal/device/infrastructure/carrier/*`, `internal/device/infrastructure/sidecar/resolve.go` (fallback), `scripts/crossbuild`, a release configuration, `tests/architecture/*` if a framework confinement is needed. |
| Tests | Carrier: fresh extraction, idempotent reuse, digest-mismatch refusal, re-extraction when a tampered or partial copy is found, atomic publish under a crash between write and rename. Resolver: environment override wins; embedded path used otherwise. A release dry run that builds every target with its embedded helper. |
| Verification | `make verify` stays green, plus a release dry run (build all targets with the embedded helper and produce the install artifacts) a reviewer can run locally. |

Record material architecture exceptions in an approved ADR. This slice needs none: it adds one infrastructure package and build wiring, changes no boundary, and introduces no third-party dependency into the core.

## What is deliberately not in this plan

This plan does not change any device command, the MCP protocol, sign-in, or the local record format. It does not add cloud upload. It does not bundle `adb` or Xcode. It does not choose the signing identities or set up notarization — those are enabled in this stage only if a chosen channel requires them. It covers only the mechanism by which the helper is carried inside the platform binary and made to run, and the shape of the standard distribution that carries the result to a person.
