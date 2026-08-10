# Standards exceptions

Recorded per `DEL-010`. Each entry is a time-bounded deviation from an approved standard. An entry takes effect only once its owner marks it Accepted; until then it is Proposed and does not waive the standard. The enforcement point in code links back here.

## EXC-001: SEC-022 protocol limits are partially deferred

- **Status:** Proposed (pending owner approval).
- **Standard:** `SEC-022` (protocol boundaries declare and enforce limits for encoded and decoded body size, string/collection length, nesting, pagination, decompression ratio, execution time, concurrency, and output).
- **Owner:** Drizz Platform maintainer.
- **Scope:** the local MCP stdio transport (`internal/transport/mcp`). Encoded frame size and frame validity (single complete JSON value per newline-delimited frame) are enforced now. Decoded nesting, string/collection length, concurrency, execution time, and output size are deferred.
- **Reason:** no product tool surface exists yet, so there is no request payload, handler, or response whose decoded shape or execution can be bounded meaningfully. Bounding them now would be speculative.
- **Risk:** low. Without tools, a client can only drive protocol negotiation; oversized or malformed frames are already rejected before decoding, and the server holds no untrusted decoded payload.
- **Compensating control:** the 1 MiB frame limit and complete-frame validation in `connection.go` bound memory per message; fail-safe classification keeps client and transport errors out of reporting.
- **Review trigger / expiry:** re-open before the first MCP product tool is accepted (roadmap Stage 6). The deferred limits must be defined and tested as part of that tool's contract.
- **ADR:** not required; this defers implementation timing, not an architecture decision.
