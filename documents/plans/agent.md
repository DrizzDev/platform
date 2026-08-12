# Agent Integration Plan

Status: Draft for owner review

## What this delivers

This stage lets a person who is using an outside agent application — such as Claude Code or Codex — ask Drizz to perform actions on a device that is connected to their computer, and have every action recorded locally as it happens. The first release supports the full set of device commands that the device helper provides today — observing the screen, interacting with it, managing apps, and managing emulators — because the goal is to give an agent the complete device surface, not a narrow slice of it. Everything runs on the person's own machine; nothing is uploaded to Drizz's servers in this stage.

The goal for this stage is a complete, working local path: an agent asks for an action, Drizz carries it out on the connected device, and Drizz writes a durable local record of what was requested, what happened, and the result. Uploading those records to Drizz's cloud is a real part of the product, but it is deliberately a later step and is not part of this stage.

## Terms used in this plan

A few words appear throughout, so they are defined here in plain terms.

An **agent application** is a program the person already uses to work with a model, such as Claude Code or Codex; it owns the conversation and decides what to do, and it can call outside tools. The **Model Context Protocol**, written **MCP** below, is a single shared standard that every agent application uses to discover and call an outside tool; it is the same protocol for every agent, so there is one path to support rather than one per application. A **plugin** is a small installable package that tells an agent application how to reach Drizz. A **hook** is an optional message that some agent applications can send to notify Drizz that something happened on their side, such as the person typing a prompt; only some applications offer hooks, and they are the only agent-specific part of the work. The **integration manager**, referred to below as the installer, is the part of Drizz that finds an installed agent application on the machine and points it at Drizz without disturbing the person's other settings.

## What already exists

Earlier stages built the pieces this stage connects together, and each of them works on its own today. The device helper that Drizz reuses can already drive a real device across a broad set of commands — listing devices, reading the screen, interacting with it, managing apps, and managing emulators. Drizz's own neutral wrapper around that helper currently covers only a few of those commands, so widening it to the full set is part of this stage. The recording layer can write a durable, ordered local record of an execution and its captured data, survive a crash, and later hand that record to synchronization when synchronization is turned on. Drizz already runs as a single installed application with a command line and with a starting point that an agent can connect to over MCP, and it already has a native sign-in. None of these pieces are wired together into a single working flow yet, and doing that wiring is the substance of this stage.

## The order of work, and why

The work is ordered so that the central, riskiest part is proven first and each later part becomes small and predictable.

The central part is one action running the whole way through: an action is requested, the device-control layer performs it, and a local record is written. This central part can be exercised two ways without any installer at all — directly from the Drizz command line, and from a small test client that speaks MCP and calls the tool the way a real agent would. Because both of those callers exist without installing anything into Claude or Codex, the central part is fully provable on its own, first.

The installer is a genuinely separate piece whose only purpose is to connect a real Claude or Codex to Drizz, and it can only be proven against the real applications. Building it before the central part works would mean installing something that has no working action to call, so it is built second, once there is a real action for a real agent to reach. Because MCP is one shared standard, the same action path serves Claude, Codex, and any other MCP-capable agent; the only agent-specific work is reading each application's optional hook messages to capture extra context such as the person's prompt.

Uploading records to Drizz's cloud is left for last, as its own step, so that this stage can be completed and proven entirely on the person's machine.

The result is a simple sequence: build and prove the action path locally, then add the installer and the optional per-agent context, then add cloud upload in a later stage.

## The device commands in this release

This release exposes the full set of device commands the device helper provides today, organized into five families. Observing the screen reads what is on it and changes nothing. Interacting with the screen performs a normal, intended input. Managing apps and reading device information cover the surrounding context an agent needs to act sensibly. Managing emulators covers the emulator images a developer runs a device on.

Each command carries two names, and they always refer to the same single command defined in one place. Inside Drizz's own code a command keeps a short single-word identifier, and its recorded events use a dotted form such as a screen capture being recorded under a device-capture name. The name an agent or a person sees is spelled out in full so it is never ambiguous: the agent connection offers each command under a clear name such as `TakeScreenshot` or `PressBack`, and the command line offers the same command spelled with dashes, such as `take-screenshot` or `press-back`.

The commands, by family, with the name an agent sees:

| Family | Command | What it does |
| --- | --- | --- |
| Observe the screen | `TakeScreenshot` | Read the current screen and return an image |
| | `TakeSnapshot` | Return the screen image together with its on-screen element tree |
| | `GetUIHierarchy` | Return the on-screen element tree on its own |
| | `GetScreenSize` | Return the screen's width and height |
| Interact with the screen | `Tap` | Press a chosen point |
| | `Swipe` | Drag from one point to another |
| | `Pinch` | Pinch the screen to zoom |
| | `PressButton` | Press a named hardware or remote button |
| | `PressBack` | Press the back button |
| | `PressHome` | Press the home button |
| | `TypeText` | Type text into the focused field |
| | `ClearText` | Clear the focused text field |
| | `SetLocation` | Set the device's reported location |
| Manage apps and context | `InstallApp` | Install an app package |
| | `LaunchApp` | Launch an app |
| | `TerminateApp` | Stop a running app |
| | `ClearAppData` | Clear an app's stored data |
| | `ListInstalledApps` | List the apps installed on the device |
| | `ListRunningApps` | List the apps currently running |
| | `GetForegroundApp` | Return the app currently in the foreground |
| | `GetCurrentURL` | Return the current link open in the active app |
| Device information | `ListDevices` | List the connected devices |
| | `GetFreeDiskSpace` | Return the device's free disk space |
| Manage emulators | `ListEmulatorImages` | List the available emulator images |
| | `BootEmulator` | Start an emulator from an image |
| | `PauseEmulator` | Pause a running emulator |
| | `ResumeEmulator` | Resume a paused emulator |

Two things the helper can do are deliberately left out of this release. Screen video recording captures a heavier, different kind of data and is added later on its own. The helper's internal housekeeping — preparing its on-device engine, readiness checks, warm-ups, and teardown — is not exposed as commands at all, because that is machinery Drizz manages on the person's behalf, not something an agent asks for.

## The pieces of work

The stage is delivered as a sequence of self-contained pieces, each built, reviewed, and proven before the next begins.

The **first piece** defines a single, transport-neutral registry of the supported commands — each command's outward name, its short internal name, the information it needs, and what it returns — in one place, so that the command line and the agent connection both offer exactly the same commands described the same way. The registry is deliberately built to hold commands from more than one source: the device commands in this release, and the cloud commands planned for a later version, which drop into the same registry and the same shape without disturbing what is already there. Adding a command, now or later, is one entry in this one place, never wiring scattered across the application.

The **second piece** exposes every registry entry two ways: as commands the person can run directly on the Drizz command line, and as tools an agent can call over MCP. Both front doors render the same registry, and both routes lead into the same underlying command, so there is one behavior with two front doors rather than two separate implementations.

The **third piece** records every action locally as it runs, reusing the recording layer built earlier. Each requested action produces a durable local record of what was asked, what the device-control layer did, and the result, written before the result is returned.

The **fourth piece** proves the first three end to end without any installer, using the command line directly and a small test client that speaks MCP, on a real connected device. This is the point at which the central path is demonstrably complete and trustworthy.

The **fifth piece** adds the installer and the per-agent context. The installer detects an installed Claude, points it at Drizz safely, and adds Claude's optional hook messages so Drizz can capture the surrounding context of an action; this is then repeated for Codex. Each is proven against the real application, including installing, using, updating, disabling, and cleanly removing the integration.

The **sixth piece**, deferred to a later stage and named here only for completeness, uploads the local records to Drizz's cloud.

## The current piece in detail: the full command set, running for real

A working skeleton is already built and committed: the single list of commands, both front doors rendering from it, and the on-demand device runtime that lets the application start normally and answer plainly that device support is not set up on this computer yet. That skeleton proved the shape from end to end for a first pair of commands. This piece makes the commands actually reach a device and be recorded, and widens the surface to the full command set. It is built in a deliberate order: the whole path is first made to work and proven end to end for a single command, then the remaining commands are added onto that already-proven path. Because every request and result crosses the internal link as a strongly-typed, checked shape, adding each further command is a small, low-risk addition — one more typed entry on each side — rather than new plumbing. It is written for approval before any of it is built, and it has four connected parts.

The first part is the helper program itself, which Drizz builds and compiles here rather than waits on from elsewhere. Drizz reuses the existing device library unchanged and adds a thin serving layer on top of it — the part that receives a request, calls the matching device operation, and returns the result — then compiles the two together into the single standalone helper program that Drizz carries and runs. Alongside it, Drizz's own neutral wrapper around the helper and the internal link that carries a request to it gain the commands, starting with the first one and then the rest, so every command in the release has a real path to the device.

The second part turns the single list of commands into a general registry — one typed entry per command, each carrying its outward name, its short internal name, the information it needs, and what it returns. Both front doors render every entry, so no command is ever wired separately in the command line, the agent connection, or the recording; adding one is a single entry in this one place. The registry is built to hold more than device commands, so the cloud commands planned for a later version drop in the same way and are offered and recorded the same way, reporting that the person is not signed in exactly as device commands report that support is not set up.

The third part connects the device helper and the local recorder into the shared path, so that each command runs on the device and produces a durable local record before its result is returned. How the helper is delivered was decided in the earlier device stage and is only restated here: the helper is carried inside the Drizz program itself, so installing Drizz installs the helper; the first time a device command runs, Drizz places the helper in a protected location, checks that it is exactly the trusted version, starts it, and reuses it. The person never installs, locates, or configures the helper separately. Because the helper program is built and compiled here, the path is proven against the real helper on a real device; where a test must run without a device, a scripted stand-in that speaks the same internal protocol stands in so the surrounding logic is still exercised.

Because the helper is a single long-running program, it is started once and kept running for the life of the connection, then shut down cleanly at the end, rather than being started and stopped for every individual command. This matters most for the agent connection, where an agent may ask for many commands in a row: restarting the helper each time would be slow and wasteful, while keeping one running lets each command return quickly. If the helper stops unexpectedly in the middle of a connection, Drizz notices and starts a fresh one automatically, pausing briefly between attempts so a repeatedly failing helper is given up on cleanly rather than retried forever; only the commands that were in flight at that moment fail, each with a clear message the caller can retry, and the connection itself keeps working. This recovery is already part of the helper-supervision layer built in the earlier device stage, so this piece only has to keep one helper for the life of the connection and shut it down cleanly at the end.

Each record is kept in the same private, per-user location on the machine that Drizz already uses for its other local data, and a recording that is still in progress is protected so that routine cleanup never removes it while it is being written. The length of time an in-progress record stays protected is a single, named setting rather than a number scattered through the code. Screenshots are treated as sensitive: the only place a captured image is written is the local record store, never a log line, a diagnostic, or any usage measurement, and this piece adds a test whose sole purpose is to guarantee that the captured bytes never appear on any of those channels.

The fourth part proves the whole path, first for the single starting command and then for the rest. The primary proof runs against the real compiled helper on a real connected device, both from the command line and through a real agent client, so that a requested command is shown to reach the device, carry out its work, and leave a durable local record — and, for a screenshot, produce a real captured image. The scripted stand-in is used for the tests that must run without a device, so the surrounding logic is still exercised everywhere. Where the real-device run cannot be performed in the current environment, that is reported honestly as not yet run rather than reported as passing.

## How each piece is proven

Every piece is proven with focused automated tests for the code it changes, and the action path is additionally proven on a real connected device, because a screenshot or a tap can only be trusted when it has run against a real device. The end-to-end proof for the central path uses the command line and a real MCP client rather than only test doubles. The installer and per-agent pieces are proven against the actual Claude and Codex applications, because a claim of compatibility with a real application cannot be demonstrated with a stand-in. Where a runtime genuinely cannot be exercised in the current environment, that is reported honestly as unavailable rather than reported as passing.

## What is deliberately not in this release

This stage does not upload anything to Drizz's cloud; that is a later step. It does not add screen video recording or any of the helper's internal housekeeping as commands. It does not change how the person signs in, and it never places any Drizz credential inside a plugin, a hook message, a record, or a diagnostic. It does not add any external tracing system, an extra runtime such as a separate scripting environment, or any dependency drawn from outside research; any such addition would need its own separate approval in the stage that owns it.
