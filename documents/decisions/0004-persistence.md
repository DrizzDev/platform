# ADR 0004: SQLite journal and file artifact storage

- Status: Accepted
- Date: 2026-07-31
- Related: [Technology stack](../stack.md)

## Context

The local platform needs crash-safe execution records, offline work, synchronization progress, leases, and bounded storage without requiring users to install a database.

## Options considered

Plain files make related atomic updates, querying, and leases difficult. PostgreSQL is operationally inappropriate for a desktop installation. SQLite provides transactions and recovery in an embedded database. Large blobs are more efficiently streamed and cleaned as files.

## Decision

Use SQLite in WAL mode for metadata and durable work. Store large artifacts as content-addressed files with digest, size, state, and references in SQLite. Store credentials only in the operating-system credential vault.

## Consequences

The system gains atomic local state and straightforward installation. It must manage one-writer contention, migrations, WAL checkpoints, filesystem permissions, artifact consistency, leases, and disk pressure carefully.

## Validation

Run crash, corruption, disk-full, duplicate-process, migration, cleanup, and power-loss-oriented tests on every supported operating system.

## Review trigger

Revisit if measured write volume or multi-process concurrency exceeds SQLite's validated limits.
