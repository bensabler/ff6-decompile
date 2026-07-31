# DISC-0005: Standard base-damage calculation

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0014, EXP-0015 (records under `docs/experiments/`).

## Discovery
The standard path computes base = power*4 + (power*`+$11AE`*`+$11AF`)>>5 at `ROMCPU:$C22B69-$C22B9C` (24-bit product via the wrapper's natural memory layout, shifted by the `$C20DD1` helper). Numerically closed live: 60/28/4 -> 450 -> defense -> 346-observed Fire Beam.

## Go implementation
`battle.BaseAmountStandard`

## Tests
`TestBaseAmountStandard, TestDamageChainGolden`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
