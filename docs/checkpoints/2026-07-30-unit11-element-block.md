# Checkpoint 2026-07-30 — Unit 11 complete (elemental-modifier block decoded)

## Current question
Question #22: the base-amount formulas at `ROMCPU:$C20C9E` / `$C20D87`
(the innermost damage computation).

## State
Supervised daytime session. Mesen live with bridge + injected watches
(pfCount, dseen, cseen, eseen, qseen, dlog — code in EXP records).

## Work completed
EXP-0010: elemental-modifier block (`$C20B83`–`$C20C2C`) decoded
byte-exact (joins the EXP-0006 dump): `$11A6` gate; base-amount JSR
variants by `+$11A4` bit 7; polarity flips; element block vs `+$11A1`
with battle-wide nullify (`+$3EC8`) and per-target masks in two more
`$14`-stride family arrays (`+$3BCC` flip|zero, `+$3BE0` double|halve)
applied first-match-wins with an `$8000` doubling guard. Go:
`battle.ElementResponse` + `ApplyElementResponse` (behavior-derived
names) + 14-case table test. FN-0009 indexed; 02/04/08, hypotheses log,
statistics updated.

## Last raw observation
`rom_C20B83_141.hex` (SHA-256 `3a9034d3…c0a39f`).

## Active emulator state
Field, post-victory (from EXP-0007's battle). cp1/cp2/cp3 unchanged.

## Tests and quality gates
gofmt clean; build/vet pass; `go test ./...` pass (5 packages; battle
suite now 4 test functions + fuzz-adjacent tables). Run 2026-07-30.

## Git status
`main`, 15 commits ahead of origin; committing Unit 11 now.

## Blockers
None hard.

## Exact next action
Unit 12 / EXP-0011: dump `ROMCPU:$C20C9E`–`$C20D86` and `$C20D87`–
`$C20E00` (both base-amount routines); decode toward stats/multipliers/
randomness; correlate with live pre-element `$F0` values if needed.

## Recommended next command
Continue autonomous session (Unit 12).
