# Checkpoint 2026-07-30 — Unit 17 complete (variance confirmed; wrap)

## Current question
#26/#28: decode the enemy/physical base path and find the RNG state it
reads.

## State
Session wrapping. **Mesen running** with bridge; injected watches
vbRef/vqRef (variance loggers) + dlog resident; local battle savestate
`mesen/out/exp10-battle.mss` available for identical-state trials —
this is the reusable experimental asset for the RNG hunt.

## Work completed
EXP-0016: timing-dependent variance Confirmed via identical-state
trials (miss-vs-hit first action; varying bases 7/0/12; same base
queuing 4 vs 2). Miss-path writer `$C22C02` found. `createSavestate`
queued-callback technique established and recorded.

## Tests and quality gates
No Go changes this unit; last full run green at Unit 16.

## Git status
`main`; committing this unit now; not pushed.

## Evidence paths
`exp10.log`, `exp10-battle.mss` (local), VAR-capture eval dumps in the
transcript.

## Blockers
None hard.

## Exact next action
EXP-0017: hand-decode `$C22B9D`–`$C22C10` from `rom_C22B40_192.hex`
(+ small supplemental dump if the miss path extends past `$C22C10`);
identify every memory read in the path; the timing-varying one is the
RNG state; verify with a write watch during identical-state trials
from `exp10-battle.mss`.

## Recommended next command
`/resume-session` in a fresh context; continue with EXP-0017.
