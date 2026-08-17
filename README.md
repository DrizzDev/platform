# Drizz

**Let AI coding assistants — and your command line — drive real mobile devices, with everything recorded on your own machine.**

## Overview

Drizz lets an AI coding assistant such as Claude Code, Codex, or Gemini — or any tool that speaks the Model Context Protocol — perform real actions on a phone or emulator connected to your computer, and it keeps a private, timestamped record of everything as it happens. You install one program, connect it to your assistant with a single command, and the assistant can then observe the screen, tap and swipe, type text, install and launch apps, read device information, and manage emulators. It works the same on Android and iOS. Everything is recorded and kept on your own machine.

## How it works

Drizz runs entirely on your own computer as a single, self-contained program. When your assistant asks to do something, Drizz carries the request out on the connected device through a device helper that it installs, verifies, and manages for you — you never download or configure that helper separately, because it is carried inside the one program you install. Every action is written to a durable local record describing what was requested, what actually happened, and the result. The identical behavior is available two ways — to an assistant over the Model Context Protocol, and directly from your command line — so a person and an assistant always get the same results, described the same way, from the same source.

## Requirements

Drizz itself needs nothing beyond the single installed program. To talk to your own devices, you need the standard platform tools already used for mobile development, which Drizz detects and clearly guides you to install if they are missing:

- **macOS, Linux, or Windows** for the Drizz program itself.
- **Android:** the Android Debug Bridge (`adb`) and platform-tools available on your `PATH`.
- **iOS:** Xcode with its command-line tools and iOS platform SDK. Drizz builds and provisions the on-device components it needs automatically.

## Install

**Homebrew (macOS and Linux):**

```sh
brew install DrizzDev/tap/drizz
```

**Shell installer (macOS and Linux):**

```sh
curl -fsSL https://get.drizz.dev | sh
```

**Windows:** download the archive from the [latest release](https://github.com/DrizzDev/platform/releases) and place `drizz.exe` on your `PATH`.

## Connect your assistant

Connecting takes one command per assistant. Drizz finds the assistant already installed on your machine, adds itself to that assistant's configuration without disturbing any of your other settings, and always keeps a backup of the original so the change can be undone cleanly. By default it also records the prompts and replies around each action as local context, so a recorded action can later be understood alongside the conversation that led to it; pass `--no-capture` to connect without that. After connecting, restart the assistant so it loads the new tools.

```sh
drizz login                         # sign in once

drizz connect enable claude-code    # connect one assistant (also records context)
drizz connect enable                # connect every detected assistant
drizz connect enable codex --no-capture   # connect without recording context

drizz connect list                  # see which assistants are connected
```

Supported assistant names are `claude-code`, `claude-desktop`, `codex`, and `gemini`. **Restart the assistant after connecting** — assistants load their tools at startup, so the Drizz tools appear in the next session.

## Use it

Once an assistant is connected, ask it in plain language. A good first prompt:

> Using the Drizz tools, list the connected devices, take a screenshot of the first one and describe what's on screen, then tell me which app is in the foreground and how much free disk it has.

Everything is also available directly from the command line, which is useful for scripting and for trying commands yourself:

```sh
drizz list-devices
drizz take-screenshot <serial>
drizz tap <serial> <x> <y>
drizz launch-app <serial> <package>
```

## What Drizz can do

Drizz exposes one consistent set of device actions to both your assistant and your command line: observing the screen (screenshots, the on-screen element tree, screen size), interacting with it (tap, swipe, pinch, type, clear, back, home, buttons, location), managing applications (install, launch, terminate, clear data, list installed and running apps, read the foreground app and current link, read free disk space), and managing emulators (list images, boot, pause, resume). Every action is recorded locally as it runs, and a capability that a particular device cannot perform is reported clearly rather than failing obscurely.

## Recording and privacy

Recording is local-first and on by default. Each action produces a durable record on your own machine, and connecting an assistant also captures the prompts and responses around those actions as context, unless you opt out. This context capture can be turned off at connect time with `--no-capture`, or removed later with `drizz connect uncapture <assistant>`. No device content, screen image, typed text, or prompt is ever placed into Drizz's own operational logs or diagnostics. Every record and any captured context is stored only on your own computer.

## Manage or remove the connection

Removing Drizz from an assistant is as clean as adding it — only Drizz's own entries are touched, and every other setting you had is preserved. Disconnecting removes the tools, and uncapturing removes the context recording; each keeps a backup of the file it changed. After changing a connection, restart the assistant so it reflects the new state.

```sh
drizz connect disable claude-code    # remove the tools (or: with no name, from every assistant)
drizz connect uncapture claude-code  # stop recording prompts and responses
drizz logout                         # sign out
```

To remove the program itself, run `brew uninstall drizz`, or delete the binary you installed.

## Build from source

Drizz is a single program built with the pinned Go toolchain, plus a device helper that is compiled and carried inside it for release builds.

```sh
git clone git@github.com:DrizzDev/platform.git
cd platform
make build        # build the program
make verify       # run the full check suite (format, lint, tests, race, cross-build, and more)
```

A release compiles the device helper for every supported target, carries it inside the program, and publishes the signed archives, checksums, Homebrew cask, and shell installer.

## Security

Drizz is local-first by design. It signs you in through the system browser using a standard native flow and stores durable credentials in your operating system's secure credential store; short-lived material is never written to ordinary files. Product-owned writes stay under a protected per-user location and are checked against tampering, and the carried device helper is integrity-verified before it is ever run. Sensitive content — screen images, typed text, and your prompts — is excluded from all operational telemetry by policy, proven by tests. Authorization for any future cloud access is always evaluated against current trusted state at the point of use.

## License

Drizz is licensed under the Business Source License 1.1, converting to the Apache License 2.0 on 2030-08-15. See [LICENSE](LICENSE).
