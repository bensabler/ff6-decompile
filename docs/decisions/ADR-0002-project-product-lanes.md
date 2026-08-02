# ADR-0002: Project product lanes

- **Status:** Accepted
- **Date:** 2026-08-02
- **Context:** [DEMO-0001](../demo/DEMO-0001-new-game-to-whelk.md), Unit 10
- **Supersedes:** nothing

## Context

This repository has been described as one thing — "FF6 reverse engineering" —
but it has been producing at least four different kinds of artifact for some
time, with different audiences, different correctness criteria, and different
legal constraints:

- research records whose criterion is **falsifiability**;
- Go code whose criterion is **it runs and its tests pass**;
- generated assets whose criterion is **byte-identical regeneration from the
  operator's own ROM**;
- and, prospectively, 65816 source whose criterion would be **byte-matching
  assembly of the original ROM**.

Conflating them causes concrete mistakes. Two have already happened. Deviation
D1 shipped because a *generated asset* was consistent with the *code that
generated it* while both disagreed with the *research record* — three artifacts
from three lanes, and nothing compared across the boundary. Unit 10 found the
same shape twice more in the records themselves.

The immediate trigger is that a matching SNES reconstruction is now a plausible
future direction, and it must not be allowed to compete with DEMO-0001 for
attention or to leak assembly-level concerns into the Go port.

## Decision

Recognize **five logical products** in one repository. Name them, say what each
owns, and say what its correctness criterion is. Do not split repositories, and
do not move any existing code or record to establish them.

### Lane 1 — Evidence and reverse-engineering laboratory

The falsifiable record of what the original does, and the tools that produce it.

`docs/` · `dashboards/` · `indexes/` · `manifests/` · `schemas/` ·
`mesen/` · `.claude/` · `cmd/ff6lab` · `internal/audit` · `internal/census` ·
`internal/research` · `internal/rom` · `internal/extract` ·
`internal/mesenstate` · `internal/scenario` ·
`local_artifacts/{experiments,scenarios,static-analysis}`

**Criterion:** every claim carries an observation, an interpretation kept
separate from it, alternatives, a falsifier, and a confidence level. Negative
results are results.

**May** read the ROM, drive the emulator, and hold addresses.

### Lane 2 — Native Go port

The runnable program: a reconstruction that plays, not a description of one.

`cmd/ff6demo` · `internal/engine` · `internal/content` · `internal/game` ·
`internal/graphics` · `internal/platform` · `internal/audio` ·
`internal/assetmanifest`

**Criterion:** it runs, its tests pass, and each behavior traces back to a
Confirmed record in Lane 1. Deviations are registered, never hidden.

**May not** read the ROM — not even transitively. This is the one lane boundary
that is already **mechanically enforced**: `audit.CheckImportBoundaries` walks
the import graph and fails the build if `cmd/ff6demo`, `internal/content` or
`internal/engine` can reach `internal/rom` by any path. That check is what makes
these lanes real rather than aspirational, and it has already caught a
transitive violation that a direct-import check called clean.

### Lane 3 — Matching SNES reconstruction (not started)

65816 and SPC700 source that assembles to the original ROM.

`reconstruction/snes/` — a README only, at this revision.

**Criterion:** byte-matching output. That is a categorically stricter bar than
Lane 2's, and the reason it is a separate lane: Go code may reproduce *behavior*
by any correct means, while matching assembly must reproduce the *encoding*.

**This lane must not delay DEMO-0001.** It is created now so that assembly-level
concerns have somewhere to go other than into the Go port.

### Lane 4 — Shared semantic specifications

Format and behavior descriptions that both Lane 2 and Lane 3 implement, and
that neither owns.

`specs/` — a README only, at this revision.

**Criterion:** implementation-independent, and precise enough that two
implementations from it agree byte-for-byte.

Today these specifications exist, but embedded: inside `docs/discoveries/`
prose and inside Go doc comments and constants. That is adequate while there is
one implementation. It stops being adequate the moment Lane 3 has any content,
because a spec that lives inside one implementation cannot arbitrate between
two.

### Lane 5 — Private generated asset archive

The commercial bytes, regenerated locally and never distributed.

`local_artifacts/archive/`, produced by `ff6lab extract`, described by the
tracked `manifests/assets.json`.

**Criterion:** byte-identical regeneration from the operator's verified ROM.
`ff6lab archive verify` checks extractor output against the manifest hash
*and* against the bytes on disk, three ways.

**Never tracked.** Lane 5 is the only lane whose output is entirely absent from
Git; what Git carries is Lane 1's description of it.

## What this ADR deliberately does not do

- **No repository split.** The lanes share one history, one module, one CI
  configuration and one review surface. Splitting would make the cross-lane
  comparisons that catch D1-class defects harder, which is precisely backwards.
- **No code moves.** Not one file changes location. A lane is a statement about
  ownership and criteria, not a directory reshuffle. Moving code to match a
  diagram would cost review history and buy nothing checkable.
- **No new enforcement.** The one boundary worth enforcing is already enforced.
  Lane 3 and Lane 4 have no content to constrain yet, and inventing rules for
  empty directories produces rules nobody has tested.

## Consequences

**Accepted:**

- Two directories exist that contain only a README. That is deliberate: an
  empty, named place with a stated criterion is what stops assembly and
  specification concerns from being filed into the Go port by default.
- Lane 4's separation is currently notional. Until specifications are extracted
  from `docs/discoveries/` and Go doc comments, Lane 2 remains their de facto
  home, and that is recorded here as a known debt rather than described as done.

**Gained:**

- A vocabulary for cross-lane disagreement, which is where this project's real
  defects have been found. "The extractor and the ledger disagreed" is a
  Lane 1/Lane 5 comparison failure; naming that makes the class searchable.
- A place to put a matching reconstruction that does not compete with the
  playable demo for the same directory or the same standard of correctness.

## Revisiting

Reconsider if: Lane 3 acquires enough content that its build needs tooling the
Go module cannot host; Lane 4's specifications are extracted and want their own
validation; or the repository's size makes a split cheaper than the cross-lane
checks it would cost.
