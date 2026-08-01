# Checkpoint 2026-08-01 — EXP-0038: mines encounter, milestone 06 (Unit 40)

## Current question
None open in the unit. EXP-0038 met its stopping condition. The next
unit is EXP-0039, a **breadth-first** reconnaissance pass (deliberate
mode change away from deep subsystem work).

## State
SCN-0001 now reaches **milestone `06-random-encounter`**. Two scheduled
power-on runs of the 26-leg route reproduce the mines random encounter
at **frame 51 307** — formation **14** (three of monster record 19),
interrupting leg 19 near tile `(26,0B)` — with **byte-identical
milestone WRAM** (`c6e69ad7…`) and identical event-flag timelines
(162 value-changing writes, frame+addr+value+PC).

## Work completed
- **EXP-0038** (record + `mesen/probes/EXP-0038.lua` +
  `MinesEncounterRoute()` in `internal/scenario/route`, guarded by the
  probe-sync test, which now maps each probe to its own model).
- Mines entry corridor mapped by GUI recon and walked on schedule:
  `(26,1C)→(26,0B)→(28,0B)→(28,09)`; a scripted event at ~`(2A,09)`
  (dialogue, party splits) registered as **CEN-EVENT-0009** and
  deliberately not entered.
- **Seventh** independent verification of EXP-0030's formation table
  (live staging = `ROMFILE:0x0F62D2` byte-for-byte).
- Mines interior yields **two** formations (14, 44) — a B13 datum.
- **Neither traversal nor a random encounter writes any event flag**
  in the three verified arrays.
- Census (CEN-WORLD-0002/0006, CEN-MONSTER-0004, CEN-EVENT-0009),
  scenario record + manifest (B11/B12/B13), indexes, dashboards synced.
- Two probe defects found and fixed: `shot()` and `mkstate()` were
  dropped when EXP-0037 was derived from EXP-0036 and were nil globals.
  Artifact writing is now `pcall`-guarded so no artifact error can
  strand a run again.

## What remains uncertain
- The `+$11E0` producer / encounter mechanism (step vs position vs
  zone) — CEN-WORLD-0006, untouched by design.
- Whether the mines zone yields more than formations 14 and 44 (B13).
- The `(2A,09)` scripted event's identity and whether it is the
  pre-Whelk beat (CEN-EVENT-0009, B17 candidate).
- Mines branches, dead ends, objects and collision — no sweep yet.
- **Evidence gap:** no milestone-06 *savestate* was captured (the
  `mkstate` defect). The milestone rests on byte-identical WRAM, the
  project's stated assertion channel (CEN-QUIRK-0002). A re-run purely
  for the convenience artifact was not performed; the probe is fixed,
  so the next unit needing that state captures it for free.

## Active instrumentation and evidence
**No background processes running** — verified with `jobs -l` and
`pgrep` after an audit that found and cleaned three: an orphaned
`tail` from the EXP-0037 GUI monitor (STALE), the EXP-0038 run1 Mesen
(HUNG on the `shot()` defect, evidence harvested first), and run1's
monitor (redundant). Evidence under
`local_artifacts/experiments/EXP-0038/` (recon transcript + 4
screenshots, two runs × flags/snapshots/events, testrunner stdout,
`hashes.sha256`) and
`local_artifacts/scenarios/SCN-0001/06-random-encounter/` (two WRAM
dumps, one screenshot, `hashes.sha256`).

## Tests and quality gates
gofmt clean; build/vet clean; `go test ./...` green (14 packages);
`ff6lab audit` clean; census clean; `archive verify` 8/8 clean;
restricted-extension scan clean; probe-sync test guards all three
probe encodings (EXP-0036/0037 → MinesRoute, EXP-0038 →
MinesEncounterRoute).

## Git status
main; one coherent unit committed and pushed.

## Exact next action
**EXP-0039 — mines-to-Whelk breadth reconnaissance**, in
**operator-visible GUI Mesen**, starting from
`local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss`
(the verified state with a savestate artifact; the corridor to the
encounter is known and short). Operating rule: **explore widely,
register briefly, continue forward** — sweep the mines branches and
dead ends, register every visible system/content family in the census
with one bounded future question each, and advance toward Whelk;
if reached, preserve a pre-Whelk savestate, fight head-only, and
capture the first stable post-battle state. Strict anti-depth rules
apply (record was written before the pass: see the Bounds section).
