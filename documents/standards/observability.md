# Observability Standards

Status: Approved and mandatory

- `OBS-001`: Domain MUST emit no telemetry. Application MAY return typed
  operational events. Outer adapters and infrastructure map approved events to
  logs, metrics, traces, and diagnostics.
- `OBS-002`: Structured events MUST use stable dotted lowercase names and typed,
  allowlisted attributes. Free-form content is prohibited.
- `OBS-003`: Log only meaningful boundaries, state transitions, external calls,
  retries, and failures. Do not log every private function or duplicate the same
  event across layers.
- `OBS-004`: Every event MUST have an operational consumer or diagnostic
  purpose. “May be useful” is insufficient.
- `OBS-005`: Correlation identifiers MUST propagate through every asynchronous,
  retry, process, and remote boundary. A synchronous in-process pure call does
  not introduce a new correlation requirement.
- `OBS-006`: Metrics MUST cover throughput, error rate, latency, saturation,
  backlog age, and resource pressure for critical paths.
- `OBS-007`: Metric labels MUST be bounded and low cardinality. User, request,
  execution, artifact, path, URL, and device identifiers do not belong in metric
  labels.
- `OBS-008`: Traces SHOULD cover latency, failure, retry, queue claim,
  synchronization, and provider boundaries, not per-frame or per-item hot loops.
- `OBS-009`: Telemetry export MUST be bounded, cancellable, non-blocking to
  product behavior, and disabled or no-op safely.
- `OBS-010`: Local diagnostic storage MUST be size-bounded,
  corruption-tolerant, content-allowlisted, and protected if persistent.
- `OBS-011`: Alerts MUST be actionable, deduplicated, owned, and linked to a
  runbook.
- `OBS-012`: Privacy canary tests MUST prove prohibited content never reaches
  any telemetry or crash sink.
