# DISC-0001: Unified ten-slot battle arrays

## Status
Confirmed; implemented and tested.

## Supporting experiments
EXP-0003, EXP-0004, EXP-0005 (records under `docs/experiments/`).

## Discovery
One $14-stride struct-of-arrays family holds all per-slot battle state: entries 0-3 party, 4-9 enemies; ~19 member arrays incl. current/max HP (`WRAM:+$3BF4`/`+$3C1C`), MP candidates (`+$3C08`/`+$3C30`), status words (`+$3EE4`/`+$3EF8`), element masks (`+$3BCC`/`+$3BE0`), defense pair (`+$3BB8`), battler stats (`+$3B18..+$3BA4`). The delta engine and death handler are slot-uniform (enemy entries observed live).

## Go implementation
`battle.BattleSlots, battle.BattleSlotCount, battle.PartySlotCount`

## Tests
`TestApplyHPDeltaEnemySlots, TestApplyHPDeltaIsolation`

## Confidence and residue
Behavior Confirmed byte-exact per the cited experiments; semantic labels
beyond the verified behavior (stat/element names, MP identity) remain at
their recorded hypothesis levels in `docs/sessions/` and the hypotheses
dashboard.
