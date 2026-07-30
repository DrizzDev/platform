# Security and Privacy Standards

Status: Approved and mandatory

- `SEC-001`: Every external input and provider response MUST be treated as
  untrusted, bounded, validated, normalized, and mapped to internal types.
- `SEC-002`: Authentication establishes identity. Authorization independently
  evaluates action, resource, tenant or organization, and current policy.
- `SEC-003`: Authorization MUST run at the trusted boundary before use-case
  execution and again at cloud data access.
- `SEC-004`: A client-supplied identity, tenant, organization, permission,
  audit field, or trusted deadline MUST NOT override trusted context.
- `SEC-005`: Use established cryptographic protocols and standard
  implementations. Custom cryptography is prohibited.
- `SEC-006`: Distributed clients MUST NOT contain universal or recoverable
  client secrets.
- `SEC-007`: Credentials and encryption keys MUST use the approved operating
  system or workload secret store and least privilege.
- `SEC-008`: Secrets MUST NOT appear in source, ordinary configuration, process
  arguments, telemetry, crash data, test fixtures, repository artifacts, or
  support bundles. Credentials may persist only as encrypted values in an
  approved credential store. User-authorized product artifacts that may contain
  sensitive content require classification, encryption, authorization,
  retention, and redaction policy and never enter telemetry or generic
  diagnostics.
- `SEC-009`: Sensitive user content MUST be excluded from telemetry by default.
  An allowlist, not a denylist, controls telemetry attributes.
- `SEC-010`: Data MUST have classified source, purpose, sensitivity, retention,
  upload, processing, and visibility before persistence.
- `SEC-011`: Unclassified or unsupported data fails closed.
- `SEC-012`: Product-owned writes and deletions MUST remain under a resolved
  product-owned root. Reads outside that root require an explicit user grant,
  trusted platform handle, or approved system location. Every path operation
  MUST resist traversal, symlink escape, race substitution, and unintended
  recursive scope according to its authority.
- `SEC-013`: Subprocess execution MUST use an explicit executable and argument
  list, environment allowlist, working directory, deadline, output bound, and
  shutdown policy. User input MUST NOT construct shell commands.
- `SEC-014`: Downloads, updates, plugins, and provider artifacts MUST be
  integrity-verified and authenticated before execution.
- `SEC-015`: Local IPC MUST authenticate the current user or supervising process
  and defend against cross-user access and DNS rebinding where applicable.
- `SEC-016`: Raw artifacts and organization data require explicit purpose,
  scope, expiry, byte limit, authorization, and audit.
- `SEC-017`: Deletion MUST prevent resurrection from retry, cache, backup,
  derived data, or synchronization.
- `SEC-018`: Every new trust boundary MUST have a threat model covering replay,
  downgrade, credential theft, tampering, tenant crossover, provider
  compromise, support misuse, and dependency integrity.
- `SEC-019`: Release artifacts MUST be signed and accompanied by checksums,
  dependency inventory, SBOM, and provenance.
- `SEC-020`: Security claims MUST be supported by tests or current primary
  evidence. Absence of evidence MUST NOT be reported as protection.
- `SEC-021`: Secret or digest equality MUST use constant-time comparison when
  timing can reveal security-sensitive information.
- `SEC-022`: Protocol boundaries MUST declare and enforce limits for encoded and
  decoded body size, string and collection length, nesting, pagination,
  decompression ratio, execution time, concurrency, and output.
