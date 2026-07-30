# Checkpoint 2026-07-30 — Unit 18 complete (physical formula; RNG upstream)

## Current question
#29: the action-setup/hit-roll layer (the true RNG consumer).

## State
Autonomous session. Mesen running with bridge; variance loggers + dlog
resident; `exp10-battle.mss` trial asset intact.

## Work completed
EXP-0017: physical base path fully decoded byte-exact (vigor² shape,
×1.75 gate, party tail, `+$3C58,X` flags); complete read census shows
**no RNG anywhere in damage arithmetic** — misses arrive as power=0;
question #26 closed, #29 opened; `+$11AF` producer lead (`+$3B18,X`).
Go: `BaseAmountPhysical` + `PhysicalFlags` + table tests (incl. the
EXP-0016 base-7 golden closure).

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (5 packages).

## Git status
`main`, committing Unit 18 now (will be 2 ahead of origin).

## Exact next action
EXP-0018 (#29): `eval`-inject a `+$11A6` write watch (dlog pattern);
two identical-state trials from `exp10-battle.mss` with different
waits; the writer PC + its stack = the hit/setup layer; then dump it
and find the timing-varying read (the RNG state).

## Recommended next command
Continue autonomous session.
