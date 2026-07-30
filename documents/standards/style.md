# Structure and Naming Standards

Status: Approved and mandatory

These rules preserve the repository conventions explicitly chosen by the owner.
They are not suggestions and are not replaced by a language community's default
style.

## Structure

- `STYLE-001`: Every project-owned directory and package name MUST be one
  complete lowercase English word.
- `STYLE-002`: Underscores, hyphens, camel case, abbreviations, and joined words
  are prohibited in directory and package names. Approved industry protocol
  terms listed in `STYLE-013`, such as `mcp`, `cli`, and `http`, remain valid
  single-word names.
- `STYLE-003`: A concept needing multiple words MUST be represented through
  nested directories. Nesting expresses ownership and relationship.
- `STYLE-004`: Repository structure MUST remain layered and nested at directory,
  package, file, type, entity, and field level. A flat collection of unrelated
  concepts is prohibited.
- `STYLE-005`: Every project-owned source filename MUST have one complete
  semantic word as its base name. Language-required suffixes such as `_test.go`,
  operating-system build suffixes, generated filenames, `README.md`, `go.mod`,
  and migration identifiers are explicit exceptions.
- `STYLE-006`: A file contains one primary responsibility. A primary type and
  its tightly bound private implementation may share one file.
- `STYLE-007`: Tests MUST mirror the production module and directory structure.
## Names

- `STYLE-008`: Every project-owned identifier MUST use one precise domain word
  whenever the language grammar permits it.
- `STYLE-009`: Types use a single precise noun. Functions and methods use a
  single present-tense verb. Variables and fields use a single precise noun.
- `STYLE-010`: Context already expressed by the package, receiver, or enclosing
  model MUST NOT be repeated in an identifier. Use `device.Session`, not
  `device.DeviceSession`.
- `STYLE-011`: Relationships MUST be expressed through nesting instead of
  flattened multiword fields. Use `ocr.enabled`, not `ocr_enabled` or
  `ocrEnabled`.
- `STYLE-012`: Names MUST come from approved product and domain vocabulary.
  Decorative, metaphorical, fashionable, or thesaurus-generated names are
  prohibited.
- `STYLE-013`: Abbreviations are prohibited except approved industry protocol
  terms such as API, CLI, HTTP, ID, JSON, MCP, SQL, TLS, URL, and XML.
- `STYLE-014`: Vague names and suffixes such as `Base`, `Common`, `Controller`,
  `Engine`, `Handler`, `Helper`, `Impl`, `Manager`, `Misc`, `Processor`,
  `Service`, `Types`, `Util`, `Utils`, `Worker`, and `Wrapper` are prohibited
  unless the word is the exact platform or product concept.
- `STYLE-015`: Names MUST be short because responsibility is narrow, not because
  meaning is abbreviated.

## Constants and enumerations

- `STYLE-016`: Constant identifiers MUST follow Go naming and visibility:
  exported `PascalCase`, unexported `lowerCamelCase`, with approved initialisms
  kept intact.
- `STYLE-017`: Enumeration-like constant identifiers MUST follow the same Go
  naming and visibility rule. Identifier casing and stored value casing are
  separate concerns.
- `STYLE-018`: Internal semantic enumeration values MUST use
  `UPPER_SNAKE_CASE`. This rule applies to the stored string value, not the Go
  variable or constant identifier. Numeric and `iota` values have no casing.
- `STYLE-019`: An enumeration value emitted verbatim into an external API,
  persisted format, rendered output, or protocol MUST preserve that contract's
  literal vocabulary.
- `STYLE-020`: Hardcoded semantic values are prohibited. Constants and
  enumerations remain beside the module that owns their meaning; generic
  constant packages are prohibited.
- `STYLE-021`: Units such as `ms`, `s`, `bytes`, or `mb` MUST NOT appear as
  identifier prefixes or suffixes. Use a typed duration, size, rate, or amount
  and document its meaning.

## Behavior and methods

- `STYLE-022`: Behavior with collaborators, state, or lifecycle MUST live on a
  cohesive type. Standalone functions are limited to language-required entry
  points and constructors. Pure transformations belong to the owning cohesive
  type as static, class, or receiver methods according to the language.
- `STYLE-023`: Public methods are concise orchestrators that read as the
  operation's sequence. Distinct steps belong in precisely named private
  methods.
- `STYLE-024`: A public method SHOULD remain under 20 logical lines and a
  private method under 30 logical lines. Exceeding the target requires a
  cohesion review, not cosmetic splitting.
- `STYLE-025`: Visibility defaults to private. Public and module-visible
  surfaces require a real consumer.
- `STYLE-026`: Composition, OOP principles, SOLID, SRP, OCP, LSP, and hexagonal
  boundaries are mandatory. A language without classes applies these through
  types, interfaces, composition, and package boundaries.

## Models and documentation

- `STYLE-027`: Boundary models, entities, value objects, and configuration MUST
  be strongly typed. Raw dictionaries, maps, tuples, pairs, or `any` values are
  prohibited after the outermost decoding boundary.
- `STYLE-028`: Models SHOULD be immutable after valid construction.
- `STYLE-029`: Go entities, value objects, configuration, and boundary payloads
  MUST use explicit typed structs and validated construction. Raw maps and
  `any` values are confined to immediate decoding boundaries.
- `STYLE-030`: Comments are exceptional, not routine. Add one
  only for a genuinely technical abstraction, invariant, concurrency rule,
  security constraint, compatibility requirement, protocol contract, or
  non-obvious tradeoff that clear code cannot express.
- `STYLE-031`: A code comment MUST be short, crisp, technically
  precise, and well written. It MUST NOT narrate code, repeat a name, contain
  filler, or explain an obvious implementation.
- `STYLE-034`: A code comment SHOULD be one line and MUST NOT
  exceed three short lines. Longer explanation belongs in an ADR, design
  document, contract document, or test.
- `STYLE-035`: Do not add comments merely because a type, function, method,
  field, or package is exported. Prefer a precise name and small responsibility.
  Public external contracts may require concise generated-documentation text.

## Tests

- `STYLE-032`: Tests MUST cover concrete production behavior, edge cases,
  failure, and recovery. Vague, random, or artificially convenient cases are
  prohibited.
- `STYLE-033`: Test grouping follows the production type or use case. Languages
  that require function entry points may use them only as the test framework
  adapter around clearly grouped behavior.

## File size

- `STYLE-036`: Every project-owned Go source file, including tests and generated
  project code, MUST NOT exceed 500 physical lines. Approaching 500 lines
  requires a responsibility and cohesion review; splitting MUST follow ownership
  rather than arbitrary line ranges.
