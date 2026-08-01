# Checkpoint 2026-07-31 — EXP-0033: golden route segment 3, partial (Unit 33)

## Current question
SCN-0001 program. Segment 3 closed as **partial**. Next: EXP-0034
(segment 3b) — a state-driven free-movement detector.

## State
Golden route: power-on → New Game → opening presentation → Narshe
entry → first scripted battle → **battle won + rewards** (frame
32 706). Milestones 00–03 exist; **milestone 04 does not** — the
captured post-victory state was demoted into the experiment evidence
directory because it landed inside a *second* scripted battle rather
than free movement. `04-free-movement/` is deliberately empty.

Lab controls unchanged (AllZeros RAM, virgin SRAM; originals in
`local_artifacts/backups/`). No Mesen running.

## Work completed
EXP-0033 (record: docs/experiments/EXP-0033-golden-route-seg3.md):
- Battle confirm cadence (A, 10 polls every 90 frames from battle+60)
  wins the fight; both enemy slots reach 0 HP at frame 32 706.
- **Enemy slot identification:** slots 0-2 are the party; **slots 6-7
  are the enemies** (HP 40 / MP 15 each = monster record 0's
  `+$08`/`+$0A`). v1's detector watched slot 0 (a party member) and
  never fired — negative result preserved.
- Victory sequence: `Got 32 Exp. point(s)` → `Got 96 GP` →
  transition → second scripted battle. Reward windows are
  input-waiting.
- **Battle→field writeback captured:** `$C2496E` → `+$1609` (55, the
  on-screen post-battle HP) and `$C24979` → `+$160D` (24, MP) — the
  same offsets EXP-0027 located as field current HP/MP. Five further
  post-victory writers logged with Unknown meaning.
- Census: +CEN-BATTLE-0008 (victory/reward processing + writeback),
  +CEN-EVENT-0006 (multiple scripted battles precede free movement).

## Correction issued this session
EXP-0032's first draft misread formation record **byte 1** (`$0C`) as
a monster id and claimed "single monster id 12". Formation 2's id
bytes are record bytes 2-7 = `FF FF 00 00 FF FF` → **two entries of
monster id 0**, verified live (slots 6/7, HP 40 / MP 15 vs record 0).
Corrected in EXP-0032 (explicit correction section), EXP-0033,
CEN-EVENT-0005, the scenario record, and the scenario manifest. The
error and its evidence trail are preserved, not overwritten.

## What remains uncertain
- Milestone 04 undefined: how many scripted battles precede free
  movement is unknown (CEN-EVENT-0006).
- **Segment 3's two-run determinism check was not performed** —
  carried forward to EXP-0034.
- Reward/EXP/GP storage offsets and the level-up path are Unknown;
  only the HP/MP writeback is a Strong hypothesis.
- Milestone-01 capture instability (CEN-QUIRK-0002) still open.

## Tests and quality gates
census validate/sync clean, indexes regenerated, audit clean;
gofmt/build/vet/test at commit.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
`mesen/probes/EXP-0033.lua` (tracked schedule);
`local_artifacts/experiments/EXP-0033/` (three iterations' logs,
beat + reward screenshots, demoted post-victory state,
hashes.sha256).

## Exact next action
EXP-0034 (write record first): modify the EXP-0033 probe so that the
battle-init watch **re-arms after each battle end**, count battles,
and continue the walk+A cadence until no battle re-arms for a fixed
window (candidate M = 1200 frames); capture milestone
`04-free-movement` at that point, logging each battle's `+$11E0`
formation id along the way (this also answers CEN-EVENT-0006's
count question). Then run two full power-on runs and byte-compare
milestone 04 WRAM to clear the deferred determinism check.
