# DISC-0004: Defense scaling

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0011, EXP-0013 (records under `docs/experiments/`).

## Discovery
Variant-A base amounts scale by floor(amount*(255-def)/256)+1 (defense `$FF` = skip), built on the hardware 8x8 multiply at `ROMCPU:$C24781` composed by the `$C247B7` wrapper - algebraically exact, no carry loss. Defense comes from the 16-bit pair at `WRAM:+$3BB8,Y`.

## Go implementation
`battle.Scale256, battle.ApplyDefense, battle.NoDefense`

## Tests
`TestScale256, FuzzScale256Composition, TestApplyDefense, TestDamageChainGolden`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
