# Agent Integration Plan

Status: Draft for owner review

## What this delivers

This stage lets a person who is using an outside agent application — such as Claude Code or Codex — ask Drizz to perform actions on a device that is connected to their computer, and have every action recorded locally as it happens. The first release supports two actions, taking a screenshot of the device's screen and tapping the screen at a chosen point, because those are the two actions the device-control layer already supports. Everything runs on the person's own machine; nothing is uploaded to Drizz's servers in this stage.

The goal for this stage is a complete, working local path: an agent asks for an action, Drizz carries it out on the connected device, and Drizz writes a durable local record of what was requested, what happened, and the result. Uploading those records to Drizz's cloud is a real part of the product, but it is deliberately a later step and is not part of this stage.

## Terms used in this plan

A few words appear throughout, so they are defined here in plain terms.

An **agent application** is a program the person already uses to work with a model, such as Claude Code or Codex; it owns the conversation and decides what to do, and it can call outside tools. The **Model Context Protocol**, written **MCP** below, is a single shared standard that every agent application uses to discover and call an outside tool; it is the same protocol for every agent, so there is one path to support rather than one per application. A **plugin** is a small installable package that tells an agent application how to reach Drizz. A **hook** is an optional message that some agent applications can send to notify Drizz that something happened on their side, such as the person typing a prompt; only some applications offer hooks, and they are the only agent-specific part of the work. The **integration manager**, referred to below as the installer, is the part of Drizz that finds an installed agent application on the machine and points it at Drizz without disturbing the person's other settings.

## What already exists

Earlier stages built the pieces this stage connects together, and each of them works on its own today. The device-control layer can list connected devices and can perform the screenshot and tap actions on a real device. The recording layer can write a durable, ordered local record of an execution and its captured data, survive a crash, and later hand that record to synchronization when synchronization is turned on. Drizz already runs as a single installed application with a command line and with a starting point that an agent can connect to over MCP, and it already has a native sign-in. None of these pieces are wired together into a single working flow yet, and doing that wiring is the substance of this stage.

## The order of work, and why

The work is ordered so that the central, riskiest part is proven first and each later part becomes small and predictable.

The central part is one action running the whole way through: an action is requested, the device-control layer performs it, and a local record is written. This central part can be exercised two ways without any installer at all — directly from the Drizz command line, and from a small test client that speaks MCP and calls the tool the way a real agent would. Because both of those callers exist without installing anything into Claude or Codex, the central part is fully provable on its own, first.

The installer is a genuinely separate piece whose only purpose is to connect a real Claude or Codex to Drizz, and it can only be proven against the real applications. Building it before the central part works would mean installing something that has no working action to call, so it is built second, once there is a real action for a real agent to reach. Because MCP is one shared standard, the same action path serves Claude, Codex, and any other MCP-capable agent; the only agent-specific work is reading each application's optional hook messages to capture extra context such as the person's prompt.

Uploading records to Drizz's cloud is left for last, as its own step, so that this stage can be completed and proven entirely on the person's machine.

The result is a simple sequence: build and prove the action path locally, then add the installer and the optional per-agent context, then add cloud upload in a later stage.

## The actions in this first release

The first release exposes exactly the actions the device-control layer supports today, and no more. Listing connected devices lets the person and the agent choose which device an action applies to. Taking a screenshot reads the current screen and returns an image; it only reads, and changes nothing on the device. Tapping the screen presses a chosen point; this is a normal, intended action. New actions such as swiping, typing text, or pressing the back button are not part of this release because the device-control layer does not implement them yet; they are added later, each as its own small unit of work, once the path they travel is already proven.

## The pieces of work

The stage is delivered as a sequence of self-contained pieces, each built, reviewed, and proven before the next begins.

The **first piece** defines a single, transport-neutral list of the supported actions — their names, the information each one needs, and what each one returns — in one place, so that the command line and the agent connection both offer exactly the same actions described the same way.

The **second piece** exposes those actions two ways: as commands the person can run directly on the Drizz command line, and as tools an agent can call over MCP. Both routes lead into the same underlying action, so there is one behavior with two front doors rather than two separate implementations.

The **third piece** records every action locally as it runs, reusing the recording layer built earlier. Each requested action produces a durable local record of what was asked, what the device-control layer did, and the result, written before the result is returned.

The **fourth piece** proves the first three end to end without any installer, using the command line directly and a small test client that speaks MCP, on a real connected device. This is the point at which the central path is demonstrably complete and trustworthy.

The **fifth piece** adds the installer and the per-agent context. The installer detects an installed Claude, points it at Drizz safely, and adds Claude's optional hook messages so Drizz can capture the surrounding context of an action; this is then repeated for Codex. Each is proven against the real application, including installing, using, updating, disabling, and cleanly removing the integration.

The **sixth piece**, deferred to a later stage and named here only for completeness, uploads the local records to Drizz's cloud.

## How each piece is proven

Every piece is proven with focused automated tests for the code it changes, and the action path is additionally proven on a real connected device, because a screenshot or a tap can only be trusted when it has run against a real device. The end-to-end proof for the central path uses the command line and a real MCP client rather than only test doubles. The installer and per-agent pieces are proven against the actual Claude and Codex applications, because a claim of compatibility with a real application cannot be demonstrated with a stand-in. Where a runtime genuinely cannot be exercised in the current environment, that is reported honestly as unavailable rather than reported as passing.

## What is deliberately not in this release

This stage does not upload anything to Drizz's cloud; that is a later step. It does not add any new device action beyond the two the device-control layer already supports. It does not change how the person signs in, and it never places any Drizz credential inside a plugin, a hook message, a record, or a diagnostic. It does not add any external tracing system, an extra runtime such as a separate scripting environment, or any dependency drawn from outside research; any such addition would need its own separate approval in the stage that owns it.
