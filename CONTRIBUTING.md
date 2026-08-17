# Contributing to Drizz

Thanks for your interest in improving Drizz. This guide covers how to build,
test, and submit changes.

## Prerequisites

- The Go toolchain pinned in `go.mod` (the `go` command fetches it automatically).
- `make`.

## Build and verify

```sh
git clone git@github.com:DrizzDev/platform.git
cd platform
make build     # build the program
make verify    # the full check suite: format, lint, tests, race, cross-build, and more
```

`make verify` must pass before a change is merged; it runs the same checks as CI.

> The on-device helper is compiled into official release binaries only. A build
> from source embeds a placeholder, so device commands are inert unless you run an
> official binary. This does not affect building, testing, or contributing to the
> rest of the program.

## Standards

This codebase follows the engineering standards in
[`documents/standards/`](documents/standards/). In short:

- One primary responsibility per file; dependencies point toward the owner of policy.
- Boundaries use explicit, typed contracts; the hexagonal layering is enforced by
  the tests under `tests/architecture/`.
- Focused tests cover real behaviour and failure paths — not mocks for their own sake.

## Submitting changes

1. Open a topic branch.
2. Keep the change focused; add tests for new behaviour and failure paths.
3. Make sure `make verify` passes.
4. Open a pull request and fill in the template. A maintainer will review it.

## Reporting bugs and requesting features

Use GitHub Issues. For security vulnerabilities, do **not** open a public issue —
see [SECURITY.md](SECURITY.md).
