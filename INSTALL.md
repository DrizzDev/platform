# Install Drizz

Drizz lets your AI assistant interact with a real device connected to your computer — seeing the screen and performing actions on your behalf.

## Before you start

Connect a device to your computer, and have the standard tools for its platform:

- **Android:** turn on USB debugging, and have `adb` (Android platform-tools) installed.
- **iOS:** have Xcode and its command-line tools installed.

## 1. Install (macOS / Linux)

```sh
brew install DrizzDev/tap/drizz
```

## 2. Connect your assistant

Run the command for the assistant you use:

```sh
drizz connect enable claude-code
drizz connect enable codex
drizz connect enable gemini
```

Use `claude-desktop` instead if that's what you use. Then **restart the assistant** so it loads the new tools.

## 3. Try it

Ask your assistant:

> Use the Drizz tools to list my connected devices and take a screenshot of the first one.

## Remove it later

```sh
drizz connect disable claude-code
```
