# Security Policy

## Reporting a vulnerability

Please do **not** open a public issue for security vulnerabilities.

Report privately through GitHub's
[private vulnerability reporting](https://github.com/DrizzDev/platform/security/advisories/new),
or email **security@drizz.dev**.

Please include:

- a description of the issue and its impact,
- steps to reproduce or a proof of concept,
- the affected version (`drizz --version`) and your platform.

We aim to acknowledge a report within three business days and to keep you updated
as we investigate and prepare a fix. Please give us a reasonable window to release
a fix before any public disclosure.

## Data handling

When you connect an agent, Drizz records your prompts and responses for context. This
is disclosed at connect time, and you can turn it off (`--no-capture` when connecting,
or `drizz connect uncapture` later). Captured context is stored **locally**, under a
protected per-user location; it is not uploaded.

## Supported versions

Security fixes are provided for the latest released version.
