# Contributing to Drizz

Thanks for your interest in improving Drizz. This guide covers how to build,
test, and submit changes.

## License and sign-off

Drizz is distributed under the Business Source License 1.1 (converting to
Apache-2.0 on 2030-08-15; see [LICENSE](LICENSE)). By contributing, you agree
that your contribution is licensed under those same terms.

Every commit must carry a Developer Certificate of Origin sign-off. Add one by
committing with `-s`:

```sh
git commit -s -m "Your message"
```

This appends a `Signed-off-by: Your Name <you@example.com>` line, which certifies
that you wrote the change or otherwise have the right to submit it under the
project license (see [DCO](DCO)). Pull requests whose commits are not signed off
will not pass the sign-off check.

## Prerequisites

- The Go toolchain pinned in `go.mod` (the `go` command fetches it automatically).
- `make`.
- For device work: Android `adb` / platform-tools, or Xcode and the iOS SDK.

## Build and verify

```sh
git clone git@github.com:DrizzDev/platform.git
cd platform
make build     # build the program
make verify    # the full check suite: format, lint, tests, race, cross-build, and more
```

`make verify` must pass before a change is merged; it runs the same checks as CI:
format, module, vet, staticcheck, lint, architecture, test, race, vulnerability,
license, secret, and smoke. Useful individual targets include `make test`,
`make race`, `make lint`, `make secret`, `make license`, and `make vulnerability`.

Install the git hooks once so the checks run automatically on commit and push:

```sh
make hook
```

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
- Documentation changes with public behavior, configuration, or architecture
  (`CODE-034`).
- Never commit secrets, real `.env` files, device serials, screenshots, or
  captured prompts and responses (`SEC-008`).

## Submitting changes

1. Open a topic branch from the default branch.
2. Keep the change focused; add tests for new behaviour and failure paths.
3. Sign off every commit (`git commit -s`).
4. Make sure `make verify` passes.
5. Open a pull request and fill in the template, describing the motivation, the
   change, and any compatibility or migration considerations. A maintainer will
   review it; CI must pass before merge.

External contributors work from a fork: fork the repository, push your branch to
your fork, and open a pull request against `main`.

## Reporting bugs and requesting features

Use GitHub Issues, with clear reproduction steps, your platform, and the output
of `drizz --version`. For security vulnerabilities, do **not** open a public
issue — see [SECURITY.md](SECURITY.md) and use private reporting instead.
