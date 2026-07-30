# Checkpoint 2026-07-29 — Overnight Unit 2 complete (EXP-0001 verified)

## Current question
Between units. Next: H-BATTLE-0003 (is `ROMCPU:$C101FB` battle-only?).

## State
Autonomous overnight session. **Mesen 2.1.1 running at the title screen**
with `mesen/bridge.lua` active (exec log at `ROMCPU:$C10DF3`, write watch
`WRAM:+$2E78`–`+$2E7F`).

## Work completed
EXP-0001: dumped `ROMCPU:$C212F0`–`$C2141F` (304 bytes,
`mesen/out/rom_C212F0_304.hex`, SHA-256 `2800f34b…d5d56a`); verified every
battle.go disassembly claim byte-exact; full annotated listing promoted to
02; 04/05/06/08 updated; FN-0003..0006 Confirmed (code); new unknowns
(`+$11A2` bit7 selector, dispatch tail, fetch gates `+$3A3C`,
`+$3A81/+$3A82`) recorded as questions 13–16; SESSION_003 dispatch-address
inference corrected ($C21300, not $C212FF); ROM mapping Confirmed HiROM
via Mesen header parse; Session 003 evidence archived
(`mesen/out/session003/`, hashes verified).

## Last raw observation
ROM dump bytes (static, deterministic); Mesen log-window header parse
(HiROM/FastROM/Map Mode $31).

## Active emulator state
Title screen, no savestate loaded. Bridge cmd/resp protocol live
(`mesen/out/cmd.txt` → `resp.txt`, 10-frame polling).

## Breakpoints/watchers
Bridge defaults only. No eval-injected watches at this checkpoint.

## Evidence paths
`mesen/out/rom_C212F0_304.hex` (`2800f34b…`); archives under
`mesen/out/session003/`; EXP record
`docs/experiments/EXP-0001-c2-delta-engine-dump.md`.

## Files changed
EXP-0001 record + manifest; `02/04/05/06/08`; `SESSION_003.md`
(correction note); `indexes/FUNCTIONS.md`; `docs/research/ROM_IDENTITY.md`;
`internal/game/battle/battle.go` (provenance comment);
dashboards (5 files).

## Tests and quality gates
To be run at commit (documentation + comment-only Go change).

## Git status
`main`, committing Unit 2 now; no push (overnight rules).

## Unresolved decisions
None for the operator.

## Blockers
None hard. M4 gate lifted (engine code Confirmed).

## Exact next action
EXP-0002 per CURRENT_FOCUS: eval-inject counting exec watch at
`ROMCPU:$C101FB` (preserve eval text in the record), sample title vs
field (`checkpoint3-mines.mss`) vs battle (`checkpoint1.mss`).

## Recommended next command
Continue autonomous session (Unit 3).
