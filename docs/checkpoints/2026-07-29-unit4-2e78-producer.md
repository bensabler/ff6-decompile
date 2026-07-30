# Checkpoint 2026-07-29 — Overnight Units 3–4 complete (battle HP chain closed)

## Current question
Between units. Unit 5 candidates in CURRENT_FOCUS (questions #19, #13, #9).

## State
Autonomous overnight session. Mesen 2.1.1 running with bridge; game at the
"Annihilated" game-over aftermath of the EXP-0002 encounter — reload
`checkpoint3-mines.mss` before any further live work. Injected `pfCount`
watch (EXP-0002) still armed.

## Confirmed before this session
See [previous checkpoint](2026-07-29-unit2-exp0001-verified.md).

## Work completed
- EXP-0002: PerFrameBattleUpdate battle-scoped across all tested contexts
  (≈175k non-battle frames, zero fires); rate corrected to phase-dependent
  ≈0.23–0.87/frame; second one-shot caller `JSR ROMCPU:$C11090` at battle
  entry; auto-save hazard recorded (destroyed `_11.mss` battle state).
- EXP-0003: PartyDisplaySourceRefresh `ROMCPU:$C25D26` decoded byte-exact —
  copies `+$3BF4/+$3C1C/+$3C08/+$3C30/+$3EE4/+$3EF8` →
  `+$2E78/+$2E80/+$2E88/+$2E90/+$2E98/+$2EA0`. Open question #1 answered;
  H-BATTLE-0002/0006 resolved Confirmed; H-BATTLE-0004 upgraded.
- FN-0007 added; 02/03/04/05/06/08, hypotheses, indexes, dashboards updated.

## Last raw observation
EXP-0003 ROM dump `mesen/out/rom_C25CC0_176.hex` (SHA-256 `76861e12…a5dc7`);
EXP-0002 `writers` snapshot: `$C25D33` count=168 lastFrame=214340;
game-over screenshot `mesen/out/exp4-battle-after.png` (Were-Rat, Repo Man).

## Active emulator state
Game-over aftermath; do not trust current WRAM for party state. Savestates:
`checkpoint1.mss` (Narshe guard event), `checkpoint2.mss` (mines entrance
field), `checkpoint3-mines.mss` (mine tunnels field — encounter-capable;
party: slots 0–1 dead, VICKS 46/70). `_11.mss` is the volatile auto-save —
never cite it.

## Breakpoints/watchers
Bridge defaults + eval-injected `pfCount` exec counter at `ROMCPU:$C101FB`
(cumulative; arm-frame ~34250 tonight).

## Evidence paths
`mesen/out/exp2.log`, `exp4.log`, `rom_C212F0_304.hex` (`2800f34b…`),
`rom_C25CC0_176.hex` (`76861e12…`), screenshots `exp2-*.png`,
`exp4-battle-after.png`; archives `mesen/out/session003/`.

## Files changed
EXP-0002/0003 records + manifest; 02/03/04/05/06/08; SESSION_003
correction; capability matrix (auto-save hazard); indexes; dashboards;
checkpoints.

## Tests and quality gates
gofmt clean, build/vet pass, 5 test packages pass (run at Unit 3 commit;
no Go changes in Unit 4).

## Git status
`main`; Unit 4 commit next; no push (overnight rules).

## Unresolved decisions
None for the operator.

## Blockers
None hard. Live MP observation limited by the surviving party (no MP
users among WEDGE/VICKS Magitek? — verify; Terra's Magitek has MP-costing
attacks per game knowledge, but per project rules treat as unknown until
observed).

## Exact next action
Reload `checkpoint3-mines.mss`; walk into an encounter; run question #19's
exec watch at `ROMCPU:$C25D26` (stack capture) during one enemy attack
resolution — one experiment record (EXP-0004) before touching Mesen.

## Recommended next command
Continue autonomous session (Unit 5).
