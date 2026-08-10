# ADR 0003: Local execution with cloud authority

- Status: Accepted
- Date: 2026-07-31
- Related: [Platform architecture](../architecture.md)

## Context

Device observation and action happen on a user's Android or iOS device. Routing each step through the cloud adds unacceptable latency. Drizz still needs durable history and must protect organization data.

## Decision

Execute device work and selected deterministic authoring locally. Record durable local facts and synchronize them asynchronously. Keep Drizz Cloud authoritative for organization membership, permissions, cloud resources, and cross-surface history.

## Consequences

Local use remains fast and can tolerate temporary network loss. The platform must implement a durable journal, idempotent synchronization, bounded local storage, and honest offline authorization behavior. Local cached context never grants access to cloud data.

## Validation

Measure a real device action locally and verify record synchronization across restart, network loss, duplicate delivery, and partial artifact upload.

## Review trigger

Revisit a capability when its required data or security authority exists only in the cloud, or when local distribution is not legally or technically safe.
