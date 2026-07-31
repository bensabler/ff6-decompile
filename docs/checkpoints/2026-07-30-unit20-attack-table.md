# Checkpoint 2026-07-30 — Unit 20 complete (attack-data table); session wrap

## Current question
#30: the enemy-AI action-selection layer (the true RNG consumer).

## State
Session wrapping at the context boundary. Mesen running with bridge and
resident watches (dlog, vbase/vq, sseen/slist — code in EXP records);
`exp10-battle.mss` trial asset intact.

## Work completed (Unit 20)
EXP-0019: attack-record loader decoded (MVN from `ROMCPU:$C46AC0`,
14-byte stride, into `+$11A0`–`+$11AD`); fight staging decoded (per-slot
stat tables `+$3B18/+$3B2C/+$3B68/+$3B7C/+$3B90/+$3BA4`); variance
reinterpreted as AI action selection (#30). Go:
`internal/game/attackdata` — first ROM data-format package (record
type, bounds-checked reader, typed accessors for the six verified
fields, synthetic fixtures, fuzz). ST-0004 indexed.

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (6 packages).

## Git status
`main`, 4 commits ahead of origin after this commit. Not pushed.

## Blockers
None hard.

## Exact next action
Either: (a) #30 AI-layer watches during identical-state trials; or
(b) validation unit — read real records from the local ROM via
`attackdata.RecordAt` (ROM stays in local_artifacts/) and cross-check
the Fire Beam entry (power should read 60, element bit 0), upgrading
field labels; or (c) `ff6lab battle simulate` integration.

## Recommended next command
`/resume-session` in a fresh context.
