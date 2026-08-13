<!-- drizz:managed — installed by `drizz connect enable claude-code`; re-run connect to update. Do not edit by hand. -->
---
description: Turn a device flow into a runnable, plain-English Drizz test
argument-hint: [what the test should do]
---

Author a Drizz test for this goal:

> $ARGUMENTS

## Run it fresh, every time

**Perform this flow now on the connected device** with the Drizz tools (tap, type, swipe, screenshot / snapshot),
observing each screen — **even if you did the same thing earlier in this session.** Author only from screens you just
saw. Never write the test from memory: a remembered flow drifts from the current UI and produces targets you cannot
trust. (Fresh run is the default; author from an earlier run only if explicitly asked.)

**Look with screenshots first.** See each screen with a screenshot — it is your primary evidence for naming targets and
checking state. Read the UI hierarchy only when a screenshot cannot resolve what you need. Screenshot first; hierarchy is
the fallback.

**When a command's syntax is unclear, check the docs** — `https://docs.drizz.dev/writing-tests/command-index.md` and, for
variables, `.../which-variable.md`. Only when you are unsure, not on every run.

Then write the test by following the guide below.

---

You are writing a **replayable** Drizz test — the sequence a user would replay on a fresh install weeks later, not a
transcript of what you happened to do.

## Principles

- **Replayable, not a transcript.** Capture the user-level flow that worked; drop retries and incidental steps.
- **Grounded in what you saw**, and produce only steps that actually executed. Never invent a target, value, or check.
- **Faithful to the goal, nothing more** — no setup, cleanup, waits, or checks the flow did not involve.
- **Honest when incomplete** — write only what you proved; note the rest in a comment, do not guess.

## Name the role, never the content (this is where tests go flaky)

Name a control by its **UI role**, plus a *stable* label when it has one. **Never** name it by content that changes
between runs — placeholder / hint text, prices, ETAs, counts, dates, times, or the current selection. Such text helps
you *find* an element while authoring; it is not the element's name.

- A field by its placeholder: ✅ `Tap on the search field` — ❌ `Tap on the "Search products…" field` (placeholder text
  can change or be localised)
- Assert a field: ✅ `Validate that the search field is visible` — ❌ `Validate that "Search products…" is visible`
- A control showing a value: ✅ `Tap on the location dropdown` — ❌ `Tap on the current address and ETA`

**UI role words to use:** `button`, `icon button` (icon-only: close, back, search, menu), `tab`, `toggle`, `checkbox`,
`chip`, `slider`; `search field`, `input field`, `dropdown`, `date picker`; `dialog`, `bottom sheet`, `menu`,
`suggestions list`, `list`, `grid`; `list row`, `card`, `cell`; `banner`, `toast`, `snackbar`. If it behaves like a
dropdown, call it a dropdown, not a bar. An unlabeled ✕ is a `close icon button`.

**Stable label vs dynamic choice:**
- Stable label: ✅ `Tap on Add to Cart` (a fixed on-screen label).
- Chosen by order, search, or filter — name it relatively so it works every run: ✅ `Tap on the first search result` —
  ❌ `Tap on the Running Shoes X200 result` (that exact item may not be first, or present, next run).
- `Tap` is the default for every touch, including icons. `MAP_ACTION` only for a map, canvas, or slider a description
  cannot address.

## Command reference (exact syntax)

One instruction per line, top to bottom. Keywords are case-insensitive; SYSTEM commands are capitalised. `# text` is a
comment. In the patterns below `<…>` marks something to replace — do **not** write literal `<…>` in a test, it is typed
as-is; use a real value or a `{{variable}}`.

- **Tap** — `Tap on <description>`  ·  e.g. `Tap on Add to Cart`, `Tap on the first search result`
- **Type** — `Type <text> in <field>`  ·  e.g. `Type running shoes in the search field`
- **Scroll** — `Scroll <direction>` · `Scroll <direction> by <n>%` · `Scroll <direction> until "<target>" is visible`
  (the quotes wrap only the target; `is visible` is outside). `Swipe` / `Drag` / `Slide` are synonyms.
- **Validate** (aliases `Verify`, `Confirm`, `Check`) — checking several things in one call is efficient:
  - one condition — `Validate that <condition>`  ·  e.g. `Validate that the product detail page is visible`
  - several visible items — `Validate 1. <a> 2. <b> 3. <c> is visible`  ·  e.g. `Validate 1. Search 2. Cart 3. Account is visible`
  - several conditions — `Validate the following: 1. <condition> 2. <condition>`
- **Wait** — `Wait Until <n> Seconds`
- **Conditionals** — `IF <condition> { … }` · `ELSE IF <condition> { … }` · `ELSE { … }`  ·  e.g.
  `IF a permission dialog is visible { Tap on Allow }`
- **Variables** —
  - `SET <var> = "<value>"`  ·  e.g. `SET city = "Bangalore"`
  - `SET <var> = screen(<description>)`  ·  e.g. `SET price = screen(the price under the product title)`
  - `SET <var> = API.<name>.response`
  - reference with `{{var}}` or `{{var.field}}`; dataset variables `{{var}}` are supplied externally; names are
    snake_case and case-sensitive; the `ctx` namespace is reserved.
  - `Store <source> as <var>` still works but is **deprecated** — prefer `SET`.
- **System** — `OPEN_APP <package>` · `KILL_APP <package>` · `CLEAR_APP <package>` · `MINIMISE_APP <package>` ·
  `PRESS_DEVICE_BACK_BUTTON` · `SET_GPS(latitude=<lat>, longitude=<lon>)` · `ENABLE_WIFI` · `DISABLE_WIFI` ·
  `TOGGLE_LOCATION`. Use only when the goal requires it — never add teardown the goal did not ask for.
- **Unlabeled surfaces** — `MAP_ACTION Tap on <description>` · `MAP_ACTION Drag <from> to <to>`
- **Modules and APIs** — `CALL <module>` · `CALL <module>(<param>={{value}})` · `PARAM <name>` · `API:<name>`

## Restraint — the lines you do NOT write

- No system command the goal did not ask for; **never** `CLEAR_APP` / `KILL_APP` / teardown unless the goal says reset.
- No `Wait` unless something was genuinely loading. No `# TODO` about paths you did not take.
- Combine related checks into one multi-condition `Validate` rather than many lines; assert stable states, not dynamic content.

## Example (illustrative — a fictional app, real syntax)

```
# ShopEase — search and open a product
# ShopEase (com.shopease.android), recorded from a live run on a real device.

OPEN_APP com.shopease.android

Tap on the search field
Type running shoes in the search field
Wait Until 2 Seconds

Scroll down until "Sponsored" is visible

IF a permission dialog is visible {
  Tap on Allow
}

Tap on the first search result
SET price = screen(the price shown under the product title)

Validate 1. the product detail page 2. the Add to Cart button is visible
Validate that {{price}} is visible
```

## Header + save

Begin with a one-line title comment and a short context line (the app; recorded from a live run on a real device) —
those two lines only. Save as a plain-text `.md` file, one instruction per line, named for the scenario.
