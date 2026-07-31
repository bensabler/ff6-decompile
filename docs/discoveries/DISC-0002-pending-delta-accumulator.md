# DISC-0002: Pending-delta accumulator with the 9999 cap

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0004, EXP-0006, EXP-0007 (records under `docs/experiments/`).

## Discovery
Damage/heal amounts queue per slot into two pending arrays (`WRAM:+$33D0`/`+$33E4`, `$FFFF` = none) via the accumulator at `ROMCPU:$C20C76`: existing pending (sentinel->0) + amount, clamped at `$270F` (9999). The queued amount equals the applied delta and the on-screen popup (three-anchor correlation).

## Go implementation
`battle.AccumulatePending, battle.PendingDeltaNone, battle.PendingDeltaCap`

## Tests
`TestAccumulatePending`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
