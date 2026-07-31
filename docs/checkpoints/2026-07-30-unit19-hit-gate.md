# Checkpoint 2026-07-30 — Unit 19 complete (hit/act gate localized); session wrap

## Current question
#29 endgame: the RNG state read by the gate between `$C2297D` (clear)
and `$C229D4` (populate).

## State
Session wrapping at the context boundary. **Mesen running** with bridge;
resident injected watches: dlog, vbase/vq (EXP-0016), sseen/slist
(EXP-0018) — all code preserved in experiment records.
`exp10-battle.mss` trial asset intact (the key instrument).

## Work completed this stretch (Units 18–19)
EXP-0017: physical base formula decoded byte-exact (vigor² shape);
damage arithmetic proven RNG-free; Go `BaseAmountPhysical`.
EXP-0018: hit/act gate localized between block-clear and power-populate;
populate values deterministic; miss = cleared-but-never-populated;
caller bracket `$C2319D`.

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (5 packages). Run at
Unit 18 commit; no Go changes in Unit 19.

## Git status
`main`, 3 commits ahead of origin after this commit (b58a72f, 807a209,
this). Not pushed.

## Blockers
None hard.

## Exact next action
EXP-0019: `read cpu C22950 160` and `read cpu C23190 48`; decode the
clear/populate routine; enumerate its reads; the timing-varying one is
the RNG state; verify via write watch during identical-state trials
from `exp10-battle.mss`.

## Recommended next command
`/resume-session` in a fresh context; continue with EXP-0019.
