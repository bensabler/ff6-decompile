# DISC-0007: Attack/spell record format (14 bytes)

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0018, EXP-0019 (records under `docs/experiments/`).

## Discovery
Battle actions load a 14-byte record from `ROMCPU:$C46AC0 + 14*index` (`ROMFILE:0x046AC0`) via MVN into `WRAM:+$11A0-+$11AD`; verified fields: +1 element, +2 flags (bit0 physical formula, bit7 MP dispatch), +3 bit7 MP-retarget, +4 mode, +6 power, +10 abort bits.

## Go implementation
`attackdata.Record, attackdata.Decode, attackdata.RecordAt, attackdata.TableFileOffset`

## Tests
`TestDecodeAndAccessors, TestRecordAt, FuzzDecode (attackdata)`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
