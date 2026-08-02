# AUDIT-0002 — methodology review (Phase 3)

Generator replay against frozen snapshots, and a tested matcher fixture suite.
AUDIT-0001 shipped its matcher with **zero tests** after a mid-run correctness
fix.

## Replay method

Two detached worktrees, never the evolving audit branch:

- **`581ddbc`** — the tree as it was when the generators actually ran, since
  AUDIT-0001's own outputs did not yet exist.
- **`93f7d03`** — the closure tree, to quantify self-contamination.

## Reproducibility classification

| Artifact | Classification | Evidence |
|---|---|---|
| `inventory.json` | **deterministically generated** | byte-identical vs `581ddbc` |
| `reference-graph.json` | **deterministically generated** | byte-identical |
| `historical-usage.json` | **deterministically generated** | byte-identical |
| `broken-references.json` | **deterministically generated** | byte-identical |
| `command-capabilities.json` | **manually transcribed** | byte-identical, but `classify_capabilities.py` makes **0 repository reads** — it is hand-authored constants in a Python file |
| `agent-smoke-tests.json` | **manually transcribed, aggregate-only** | byte-identical, **0 repository reads**; raw agent transcripts were never preserved |
| `AUDIT-0001-baseline.json` | **partly generated, partly hand-authored** | derives from the four deterministic inputs plus the two hand-authored ones |

**Reproducing identically proves nothing about accuracy when the generator
contains constants rather than measurements.** Two of the six "generated"
artifacts are transcription, and one of those (`agent-smoke-tests.json`) is the
sole record of evidence that no longer exists.

## The matcher corpus omits repository-root files

`build_inventory.py::corpus()` walks ten subdirectories: `.claude`, `docs`,
`dashboards`, `indexes`, `manifests`, `schemas`, `internal`, `cmd`, `mesen`,
`scripts`. It adds `CLAUDE.md` and `.claude/README.md` explicitly as
documentation entry points — but no other root-level file.

`PACKAGE_MANIFEST.json` is at the root and lists **24 of the 31 reported
orphans** by full path. The orphan count is therefore wrong: **7, not 31.**

## Self-contamination: the audit destroyed its own metric

| Replay tree | Zero-inbound resources |
|---|---|
| `581ddbc` (true input) | **31** — reproduces AUDIT-0001 exactly |
| `93f7d03` (closure) | **0** |

`AUDIT-0001-baseline.json` alone un-orphans all 31, because it enumerates every
resource id. Four AUDIT-0001 markdown records, `LATEST.md`, the checkpoint and
`SESSION_006.md` contribute the rest.

**Any future run of this method returns 0 orphans regardless of the real
routing situation.** This is a permanent, self-inflicted measurement failure,
and it is the strongest argument in this audit for separating the corpus a
validator scans from the reports the validator produces.

## Matcher fixture suite

`AUDIT-0002-matcher-fixtures.json` — **63 cases, 0 failures.** All seven
resource types × nine match conditions (full path, filename, stable id, bare
prose, alias spelling, near-match-must-not-count, self-reference exclusion,
manual entry point, auto-loaded rule), plus a pinned regression for the
skill-stem bug: `.claude/skills/example-skill/SKILL.md` must yield
`example-skill`, never `SKILL`.

**A passing suite is not automatically a vindication.** Two conditions encode
matcher behaviour that is arguably wrong, and both were chased into the real
tree rather than ratified:

1. **Bare-name prose is invisible** for shared contracts, playbooks, rules and
   templates — the matcher only sees `NAME.md`. Checked against the tree: this
   did not create false orphans, because the real references are full paths.
2. **Auto-loaded rules are invisible.** The matcher has no concept of
   `paths:` frontmatter, so a rule that is correctly never named reads as an
   orphan. This *did* create a false finding — see the errata.

## Broken references, decomposed

AUDIT-0001's "zero broken references" measured two of eight categories.

| Category | Measured? | Result |
|---|---|---|
| Broken markdown links | yes | 0 |
| Broken backtick paths | yes | 0 |
| Missing named commands | no | not_tested |
| Missing named skills | no | not_tested |
| Missing named agents | no | not_tested |
| Missing named playbooks | no | not_tested |
| Invalid declared output paths | no | not_tested — `/correlate-static-runtime` never names `docs/correlations/`, its actual output location |
| Contract/backend gaps | no | covered separately in finding verification |

## File-count taxonomy

AUDIT-0001's Phase 1 inventory reported "sessions 15", counting every `.md`
under `docs/sessions/`. Only **5** are session records; the other 10 are
topic documents (`00_PROJECT_GOALS.md` … `08_OPEN_QUESTIONS.md`) and a README.
Canonical record counts and Markdown file counts are not the same measurement.
