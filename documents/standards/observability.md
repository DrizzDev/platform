# Observability Standards

Status: Approved and mandatory

- `OBS-001`: Domain MUST emit no telemetry. Application MAY return typed operational events. Outer adapters and infrastructure map approved events to logs, metrics, traces, and diagnostics.
- `OBS-002`: Structured events MUST use stable dotted lowercase names and typed, allowlisted attributes. Free-form content is prohibited.
- `OBS-003`: Log only meaningful boundaries, state transitions, external calls, retries, and failures. Do not log every private function or duplicate the same event across layers.
- `OBS-004`: Every event MUST have an operational consumer or diagnostic purpose. “May be useful” is insufficient.
- `OBS-005`: Correlation identifiers MUST propagate through every asynchronous, retry, process, and remote boundary. A synchronous in-process pure call does not introduce a new correlation requirement.
- `OBS-006`: Metrics MUST cover throughput, error rate, latency, saturation, backlog age, and resource pressure for critical paths.
- `OBS-007`: Metric labels MUST be bounded and low cardinality. User, request, execution, artifact, path, URL, and device identifiers do not belong in metric labels.
- `OBS-008`: Traces SHOULD cover latency, failure, retry, queue claim, synchronization, and provider boundaries, not per-frame or per-item hot loops.
- `OBS-009`: Telemetry export MUST be bounded, cancellable, non-blocking to product behavior, and disabled or no-op safely.
- `OBS-010`: Local diagnostic storage MUST be size-bounded, corruption-tolerant, content-allowlisted, and protected if persistent.
- `OBS-011`: Alerts MUST be actionable, deduplicated, owned, and linked to a runbook.
- `OBS-012`: Privacy canary tests MUST prove prohibited content never reaches any telemetry or crash sink.
- `OBS-013`: Production logs MUST be line-delimited JSON with stable `timestamp`, `level`, `message`, `service`, `version`, and `revision` fields.
- `OBS-014`: Logs MUST include `request_id`, `correlation_id`, and `session_id` when the owning boundary has those values. Unavailable identifiers are omitted, never invented or emitted as empty values.
- `OBS-015`: Correlation fields are attached through an explicit typed logging scope at the boundary that validates or creates them. Generic context values, positional strings, and implicit global lookup are prohibited.
- `OBS-016`: Every log sink MUST apply the approved redaction policy. Attribute names and content remain allowlisted because key-based redaction alone cannot prove that arbitrary values contain no sensitive data.
- `OBS-017`: Additional log context MUST use the optional `metadata` JSON object. Metadata accepts a bounded set of named scalar attributes. Generic maps, nested objects, bulk payloads, and unapproved content are prohibited. Metadata is not an escape hatch for screenshots, page source, device output, or other artifacts.
- `OBS-018`: Log timestamps MUST use UTC. Every record MUST include structured source information containing a repository-relative file path, line, and function. Released binaries MUST NOT emit absolute developer-machine paths.
- `OBS-019`: Error records that carry a cause MUST place the original Go `error` under the `error` key. A record whose cause may contain client, user, credential, or provider content MUST instead be reported as an approved code alone, with the cause withheld from every sink and retained only for control flow. Configured error reporting receives `ERROR` and higher records automatically and MUST flush within bounded shutdown. Logging MUST never terminate the process itself.
- `OBS-020`: The composition root MUST construct logging, reporting, tracing, and metrics once per process. Every CLI, MCP, server, desktop, and background flow receives those process-scoped providers through explicit construction. Packages MUST NOT initialize independent providers or read observability configuration directly.
- `OBS-021`: OpenTelemetry owns traces, spans, and metrics. Error-reporting adapters MAY correlate events with the active OpenTelemetry span but MUST NOT create a competing tracing pipeline.
- `OBS-022`: Automatic PII, cookies, headers, bodies, query parameters, machine identity, logs, and metrics in an error-reporting SDK MUST default to off. Additional collection requires an approved, allowlisted contract and privacy review.
- `OBS-023`: A vendor integration MUST implement the reporting sink boundary. Adding or replacing a vendor MUST NOT change logging callers, transports, or product behavior.
