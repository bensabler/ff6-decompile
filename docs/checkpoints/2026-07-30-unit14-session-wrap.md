# Checkpoint 2026-07-30 — Unit 14 complete; daytime session wrap

## Current question
Question #23: what computes the base amount at `WRAM:+$11B0` (battle
power × level/stat layer) — the innermost undecoded layer.

## State
Supervised daytime autonomous session ending at a clean boundary.
**Mesen 2.1.1 still running** with `mesen/bridge.lua` and these
eval-injected globals (exact code preserved in the named records):
`pfCount` (EXP-0002), `dseen`/`dlog` (EXP-0004), `eseen` (EXP-0005),
`qseen` (EXP-0007). Game idle in a mines field/post-victory state.
Kill with `kill -9` (Mesen ignores SIGTERM) or keep for the next unit.

## Confirmed this session (Units 6–14; see ACTIVITY_LOG for the ledger)
- Enemy slots 4–9 in the unified 10-entry arrays; slot-uniform engine.
- Pending-delta accumulator with the 9999 cap; queue amount = displayed
  damage (three anchors).
- Target loop (`$C23442` region) → base-amount routines (variant A
  defense/halvings over precomputed `+$11B0`; variant B fraction-of-HP)
  → elemental-modifier block → accumulator → delta engine → display.
- Twelve-member `$14`-stride array family (05 master table); HP→MP
  retargeting via `+$14` index arithmetic.
- Helpers: hardware multiply/divide; defense scaling exactly
  `(amount×(255−def))/256+1`; ×1.5 chain (randomness hypothesis
  refuted — variance source still unlocated).

## Go state
`internal/game/battle`: `BattleSlots` (10 slots), `ApplyHPDelta`,
`ApplyMPDelta`, `AccumulatePending`, `ElementResponse`/
`ApplyElementResponse`, `Scale256`, `ApplyDefense`, `ChainBoost` — all
byte-exact-verified behavior, table tests + `FuzzScale256Composition`.

## Tests and quality gates
gofmt clean; `go build`/`go vet`/`go test ./...` pass (5 packages);
new fuzz target passes (8s run). Verified at wrap.

## Git status
`main`, 20 commits ahead of `origin/main` after this unit's commit.
No pushes during autonomous work (operator reviews and pushes).

## Evidence paths
All EXP-0001..0013 records under `docs/experiments/`; raw dumps and
logs under `mesen/out/` (gitignored; hashes in the records);
Session 003 archive under `mesen/out/session003/`.

## Unresolved decisions
None requiring the operator beyond reviewing/pushing the commit stack.

## Blockers
None hard. Live-experiment questions (#23, #14, #19b, #2…) need a
battle; ROM-dump questions (#15, #16, #17, #20, #6, #7) are
state-independent.

## Exact next action
EXP-0014 (question #23): inject a write watch on `WRAM:+$11B0` (both
bytes) with `dlog` stack capture; reload `checkpoint3-mines.mss`; one
encounter; correlate captured writer PCs and values with the queued
`$F0` amounts; then dump the writer's routine.

## Recommended next command
`/resume-session` in a fresh context, then continue with EXP-0014.
