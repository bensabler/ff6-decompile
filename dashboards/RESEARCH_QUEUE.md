# Research Queue

## P0 — Active frontier

- [x] EXP-0020: `+$3A70` refuted as the varying state; question #30
  refined (scheduling interpretation).
- [ ] **EXP-0021 (question #30):** capture the attack-record index at
  `ROMCPU:$C22966` entry + scheduler-adjacent reads across
  identical-state trials; determine whether any RNG feeds selection.

## P1 — Semantic debt (next behavior units)

- [ ] Live MP spend/heal verification for `WRAM:+$3C08`/`+$3C30`
  (H-BATTLE-0004/0007; justifies the published Go names).
- [ ] `+$11AE`/`+$11AF` producers and meanings (question #27; `+$3B18,X`
  lead in EXP-0017/0019).
- [ ] Cross-check attack records against the local ROM via
  `attackdata.RecordAt` (Fire Beam entry: power 60, element bit 0).

## P1 — Phase-1 vertical proofs (rebalance targets)

- [ ] Graphics: one bounded FF6 target (menu font / battle HUD tiles)
  runtime→ROM→decoder→comparison.
- [ ] Audio: one short sound effect trigger→APU→sequence/sample→decode.

## Completed (see ACTIVITY_LOG and indexes/EXPERIMENTS.md)
Battle damage pipeline end to end (EXP-0001..0019): engine, display
chain, accumulator, element response, defense scaling, both base
formulas, attack-record format, variance localization.
