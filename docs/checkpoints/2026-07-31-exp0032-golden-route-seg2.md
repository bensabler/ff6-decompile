# Checkpoint 2026-07-31 — EXP-0032: golden route segment 2 (Unit 32)

## Current question
SCN-0001 program. Segment 2 closed; next is segment 3 (EXP-0033:
Guard battle victory → free field movement).

## State
The golden route now runs power-on → first scripted battle as a
single tracked probe (`mesen/probes/EXP-0032.lua`). Milestones 00–03
exist under `local_artifacts/scenarios/SCN-0001/` with WRAM dumps,
screenshots, savestates, and hashes. Lab controls unchanged
(AllZeros RAM, virgin SRAM; originals in `local_artifacts/backups/`).

## Work completed
EXP-0032 (record: docs/experiments/EXP-0032-golden-route-seg2.md):
- Three schedule iterations, all preserved as evidence: A-metronome
  only (no battle in 60k frames), walk-only (stalls at the guard
  box), and **walk + A metronome** (reaches battle). The beat is a
  chain of input-waiting dialogue boxes separated by player-controlled
  walking — established behaviorally by that discrimination.
- Milestones `01-opening-cinematic` (frame 15000),
  `02-narshe-entry` (30000), `03-first-scripted-battle`
  (battle-detect + 120). Battle-init detection is state-driven
  (first `+$3B18` family write from `$C22800-$C22FFF`), firing at
  frame 31557 in both runs.
- **First scripted battle identified: formation id 2, single monster
  id 12.** Live `+$11E0`=$0002 and the 16-byte staged `+$3F44` record
  match ROM formation record 2 (`ROMFILE:0x0F621E`) byte-for-byte —
  a second, independent verification of the EXP-0030 formation table,
  this time on a scripted rather than random encounter.
- Determinism: WRAM byte-identical across two power-on runs at all
  three milestones (01 `011588bc…`, 02 `0f4369d5…`, 03 `24302078…`).
- Census: +CEN-EVENT-0005 (scripted battle invocation, LOCATED /
  NORMAL_PATH_VERIFIED), +CEN-QUIRK-0002 (milestone-01 capture
  instability).

## What remains uncertain
Milestone 01's **screenshot** is not byte-stable across runs while
its WRAM is (11536 vs 11339 B; savestate sizes also differ).
Interpretation (Tentative): per-scanline/HDMA state outside WRAM
sampled at an unfixed phase; genuine emulation nondeterminism not
excluded. Falsifier queued in RESEARCH_QUEUE. Until resolved, only
WRAM is a valid assertion channel at milestone 01.

## Tests and quality gates
census validate/sync clean, indexes regenerated, audit clean;
gofmt/build/vet/test at commit.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
No Mesen running. `local_artifacts/experiments/EXP-0032/` (three
iterations' logs, beat screenshots, recon shots),
`local_artifacts/scenarios/SCN-0001/0{1,2,3}-*/` (dumps, screens,
states, hashes.sha256).

## Exact next action
EXP-0033 (write record first): extend the probe from milestone 03 —
fight the Guard battle to victory with a scheduled confirm cadence
(the EXP-0021 frame-exact A pattern is the model), capture the
victory/reward processing window and the first stable post-battle
field state as milestone `04-free-movement`, then run the two-run
determinism check. Battle-end detection should be state-driven like
battle-init (candidate: the post-battle transition writes or the
enemy HP array reaching zero).
