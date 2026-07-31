# DISC-0006: Physical base-damage calculation

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0016, EXP-0017 (records under `docs/experiments/`).

## Discovery
The `+$11A2`-bit-0 path computes the vigor-squared shape: t = power (x4 enemy attackers, x1.75 gated by DP `$B2` bit 14) + `+$11AE`, then ((t*`+$11AF`)&0xFFFF)*`+$11AF`/256; party attackers add t*1.5+power and `+$3C58,X` flags (halve, x0.75). Deterministic - no RNG anywhere in damage arithmetic.

## Go implementation
`battle.BaseAmountPhysical, battle.PhysicalFlags`

## Tests
`TestBaseAmountPhysical`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
