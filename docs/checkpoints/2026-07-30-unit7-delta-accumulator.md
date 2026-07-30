# Checkpoint 2026-07-30 — Unit 7 complete (pending-delta accumulator decoded)

## Current question
Between units. Unit 8 per CURRENT_FOCUS: question #21 (damage formula →
DP `$F0` correlation experiment).

## State
Supervised daytime autonomous session. Mesen running with bridge; game in
post-victory field state. Injected watches resident (pfCount, dseen,
cseen, eseen, dlog — code preserved in EXP-0002/0004/0005 records).

## Confirmed before this session
See [previous checkpoint](2026-07-30-unit6-enemy-slots.md).

## Work completed
EXP-0006: PendingDeltaAccumulate (`ROMCPU:$C20C76`) decoded byte-exact —
`$FFFF`-sentinel init, accumulate DP `$F0`, **cap `$270F` (9999)**, store
`$C20C98`; dual-array retarget (`Y += $14`) from DP `$F2`/carry; caller
`JSR $0C2D` at `$C20C28` with gate block (`+$11A2` bit 0,
`+$3A82`&`+$3A83`, `+$3EE4,X` bit 1). EXP-0004 stack misparse corrected
in both records. Go: `battle.AccumulatePending` + constants
(`PendingDeltaNone`, `PendingDeltaCap`) + 8-case table test. FN-0008
indexed; 02/04/08 updated; question #21 opened.

## Last raw observation
ROM dumps `rom_C20C60_128.hex` (`95aa6214…`), `rom_C20420_48.hex`
(`bd420802…`), `rom_C20C10_48.hex` (`c4c4c562…`).

## Active emulator state
Field, post-victory. Savestates cp1/cp2/cp3 unchanged; `_11.mss`
volatile.

## Breakpoints/watchers
Bridge defaults + five injected globals (see Unit 6 checkpoint).

## Evidence paths
The three new `.hex` artifacts; EXP-0006 record; prior artifacts per
earlier checkpoints.

## Files changed
EXP-0006 record + manifest; EXP-0004 correction; 02/04/08;
indexes/FUNCTIONS.md; ACTIVITY_LOG; CURRENT_FOCUS;
`internal/game/battle/battle.go` + tests; this checkpoint; LATEST.md.

## Tests and quality gates
gofmt clean; `go build`/`go vet` pass; `go test ./...` pass (5 packages,
including the new AccumulatePending table). Run 2026-07-30.

## Git status
`main`, 6 commits ahead of origin; committing Unit 7 now.

## Unresolved decisions
None.

## Blockers
None hard.

## Exact next action
Unit 8 / EXP-0007: exec-log `ROMCPU:$C20C76` with DP `$F0`/`$F2`, Y, and
stack during one attack; read the displayed damage from a screenshot;
correlate; walk callers backward toward the formula.

## Recommended next command
Continue autonomous session (Unit 8).
