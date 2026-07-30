# Checkpoint 2026-07-29 — Overnight Unit 5 complete (delta producers located)

## Current question
Between units. Unit 6 per CURRENT_FOCUS: H-BATTLE-0008 enemy-slot test.

## State
Autonomous overnight session. Mesen running; battle from EXP-0004 was in
progress at last observation (VICKS 25/70 vs unknown enemies) and has been
idling — assume it resolved (likely another annihilation); reload
`checkpoint3-mines.mss` before further live work. Injected watches still
armed: `pfCount` ($C101FB exec), `dseen`/`dRef` ($33D0–$33EB write),
`cseen`/`cRef` ($C25D26 exec), `dlog` helper — exact code preserved in
EXP-0002/0004 records.

## Confirmed before this session
See [previous checkpoint](2026-07-29-unit4-2e78-producer.md).

## Work completed
EXP-0004: pending-delta writer PCs identified (setter `ROMCPU:$C20C9B`
×12; sweepers `$C2638E` ×48 / `$C26391` ×120; init `$C22408`); arrays
transient (`$FFFF` between events); **10-entry array-family discovery**
(`$14` stride: `+$3BF4/+$3C08/+$3C1C/+$3C30`; entry-9 write at `+$33E2`);
refresh trigger localized (13 event-driven calls; steady stacks top at the
`JSR $069B` @ `ROMCPU:$C21409`). H-BATTLE-0008 opened. Records updated:
EXP-0004, 04, 05, 08 (#13/#19), hypotheses, dashboards.

## Last raw observation
`dseen: C2638E=48 C26391=120 C20C9B=12 C22408=28 cseen=13`;
`WRAM:+$3BF4` = `00 00 00 00 19 00 00 00` (VICKS 25);
`+$33D0`–`+$33EB` all `$FF` between events.

## Active emulator state
As above; savestates unchanged (cp1 guard event / cp2 mines entrance /
cp3 mine tunnels). `_11.mss` volatile (auto-save) — never cite.

## Breakpoints/watchers
Bridge defaults + the four injected globals listed above.

## Evidence paths
`mesen/out/exp5.log`, labeled lines in `mesen/out/events.log`,
`exp5-battle.png`; prior units' artifacts per earlier checkpoints.

## Files changed
EXP-0004 record + manifest; 04/05/08; OPEN_HYPOTHESES (H-BATTLE-0008);
ACTIVITY_LOG; CURRENT_FOCUS; this checkpoint; LATEST.md.

## Tests and quality gates
No Go changes this unit; gates last run green at Unit 3 commit (and Unit 4
had docs only). Will re-run at any future Go-touching unit.

## Git status
`main`, 4 commits ahead of origin (no push per overnight rules);
committing Unit 5 now.

## Unresolved decisions
None for the operator. (Flag for morning: battle.go models the party view
of 10-wide arrays; revisit once H-BATTLE-0008 is settled.)

## Blockers
None hard.

## Exact next action
Unit 6 / EXP-0005: write-watch `WRAM:+$3BFC`–`+$3C07` (enemy-HP
candidates) through one encounter in which VICKS attacks (press A to
target); look for delta-engine stores with Y≥8.

## Recommended next command
Continue autonomous session (Unit 6).
