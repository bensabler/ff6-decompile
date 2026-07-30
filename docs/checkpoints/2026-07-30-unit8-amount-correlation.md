# Checkpoint 2026-07-30 — Unit 8 complete (queued amount = displayed damage)

## Current question
Between units. Unit 9 per CURRENT_FOCUS: the `$6B16` deeper-return lead
toward the damage formula (question #21).

## State
Supervised daytime autonomous session. Mesen running with bridge; battle
from EXP-0007 likely still resolving or won (last seen: Repo Man alone,
VICKS 38/70). Injected watches resident: pfCount, dseen, cseen, eseen,
qseen (+ dlog) — all code preserved in experiment records.

## Confirmed before this session
See [previous checkpoint](2026-07-30-unit7-delta-accumulator.md).

## Work completed
EXP-0007: DP `$F0` = final per-hit damage — Confirmed on three anchors
(array arithmetic, HUD mid-values, captured popup "6"); 346-damage Fire
Beam kill observed end-to-end; X/Y = attacker/target slot×2 (strong,
5/5); `$F2` direction observation (`$20` vs `$00`) recorded. Question
#21 narrowed to the formula upstream of `$C20C28`; lead: deeper return
`$6B16`. Records updated: EXP-0007 + manifest, 02, 08, dashboards.

## Last raw observation
QUEUE lines (f0=9/4/2/6 → Y=$0004; f0=$015A → Y=$0008);
`+$3BF4` final read `00 00 00 00 26 00 … 23 00 …` (VICKS 38, entry 5
enemy 35, entry 4 dead); screenshots exp8-round1..6.png.

## Active emulator state
Mid/post-battle vs Repo Man. Savestates cp1/cp2/cp3 unchanged; `_11.mss`
volatile.

## Breakpoints/watchers
Bridge defaults + six injected globals.

## Evidence paths
`mesen/out/exp8.log`, QUEUE lines in `events.log`, `exp8-round*.png`.

## Files changed
EXP-0007 record + manifest; 02; 08; ACTIVITY_LOG; CURRENT_FOCUS; this
checkpoint; LATEST.md. (No Go changes this unit.)

## Tests and quality gates
No Go changes; last full run green at Unit 7 commit.

## Git status
`main`, 7 commits ahead of origin; committing Unit 8 now.

## Unresolved decisions
None.

## Blockers
None hard.

## Exact next action
Unit 9 / EXP-0008: dump `ROMCPU:$C26AE0`–`$C26B60`; hand-verify whether
`$6B16` is a return frame of the formula path; decode what feeds DP
`$F0`. Fallback: full DP `$F4`–`$FC` snapshot logging at `$C20C76`.

## Recommended next command
Continue autonomous session (Unit 9).
