# Checkpoint 2026-07-30 — EXP-0021: action layer deterministic; battle path pauses

## Current question
Question #30 resolved for the tested window. Battle path paused at a
natural boundary (operator rebalance). Next: Unit 3 semantic-debt item
(attack-record ROM cross-check), then graphics (Unit 4) and audio
(Unit 5) verticals.

## State
Headless Mesen testrunner live (`--testrunner --timeout=7200`,
`FF6_OUT` env, bridge v2, probe EXP-0021 armed, PID under `caffeinate
-i`). The locked display breaks GUI Mesen (RenderTimer −6661 / script
window never opens) — headless is the unattended-session mode now;
recipe in memory and BLOCKERS. Trial savestate `mesen/out/
exp10-battle.mss` intact. Evidence frozen at
`local_artifacts/experiments/EXP-0021/` (hashes.sha256 written).

## Work completed
EXP-0021: two frame-exact trials (first press anchor+72 vs +270) from
the same savestate produced **byte-identical matched-ordinal action
content** — all 8 record loads index 238; enemy powers 13/0/19/0 with
the miss at the same ordinal; only cluster timing (+2/+69 frames) and
two outcome-neutral press residues (`+$3A71` $04-vs-$00, POP3 carry)
varied. EXP-0020's scheduling interpretation upgraded to Confirmed for
this window; GUI-era variance attributed to harness wall-clock jitter
(Strong hypothesis); the "free RNG" framing retired. Probe with
frame-scheduled trial driver tracked at `mesen/probes/EXP-0021.lua`.
Earlier in session: manifest-completeness repair (EXP-0019 entry +
audit reverse check, commit 6ac8108).

## Tests and quality gates
To run at commit: gofmt/build/vet/test/cli/audit (all passed at the
6ac8108 commit; no Go changes since except none — probe is Lua).

## Git status
`main`, 7 ahead of origin after this commit. Not pushed.

## Blockers
None hard. Soft: GUI/testrunner input-latching parity unverified; MP
verification needs a savestate with an MP-consuming action.

## Exact next action
Unit 3 (bounded semantic-debt item): via the live bridge, `read cpu
C46AC0 3584` (256 records × 14), archive the dump under
`local_artifacts/experiments/EXP-0022/`, scan with
`attackdata.RecordAt` for power-60 entries, and test the Fire Beam
candidate (element bit 0, power 60 per EXP-0011/0015). Write
EXP-0022 record + manifest entry first (plan-before-Mesen).

## Then
Unit 4 graphics vertical (menu font / HUD tiles: runtime → DMA → ROM →
decoder → comparison) and Unit 5 audio vertical (cursor SFX), then
final project-state synchronization.
