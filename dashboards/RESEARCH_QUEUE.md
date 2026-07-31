# Research Queue

## P0 — Active frontier

- [x] EXP-0020: `+$3A70` refuted as the varying state; question #30
  refined (scheduling interpretation).
- [x] **EXP-0021 (question #30, resolved for this window):**
  matched-ordinal action content is identical across frame-exact
  schedules — indices all 238, powers 13/0/19/0, the miss at the same
  ordinal; only timing (+2/+69-frame shifts) and two outcome-neutral
  press residues (`+$3A71`, a carry bit) vary. GUI-era variance
  attributed to harness wall-clock jitter (Strong hypothesis). The
  "free RNG" framing is retired; battle path pauses here
  (operator rebalance to graphics/audio verticals).

## P1 — Semantic debt (next behavior units)

- [ ] Live MP spend/heal verification for `WRAM:+$3C08`/`+$3C30`
  (H-BATTLE-0004/0007; justifies the published Go names).
- [ ] `+$11AE`/`+$11AF` producers and meanings (question #27; `+$3B18,X`
  lead in EXP-0017/0019).
- [x] EXP-0022: attack records cross-checked against the local ROM via
  `ff6lab attackdata scan` — record 238 power=0 + physical flag
  (Confirmed, converges with EXP-0017/0018/0021); Fire Beam candidates
  narrowed to indices 5 and 131 (Tentative — value coincidence).
- [ ] Disambiguate the Fire Beam record index: capture `+$11A0`–`+$11AD`
  (or the MVN X source) during an actual Magitek Fire Beam confirm
  (needs a menu-navigation press script).

## P1 — Phase-1 vertical proofs (rebalance targets)

- [x] EXP-0023 / GFX-0001: battle HUD font — ROM provenance Confirmed
  (raw copy at ROMFILE 0x046FC0+16·N, tiles $FF-$1FF), `tile2bpp`
  decoder + `ff6lab tiles decode2bpp` proof. Follow-ups queued: load
  path, glyph semantics, the $00-$FE compose region.
- [ ] Audio: one short sound effect trigger→APU→sequence/sample→decode.

## Completed (see ACTIVITY_LOG and indexes/EXPERIMENTS.md)
Battle damage pipeline end to end (EXP-0001..0019): engine, display
chain, accumulator, element response, defense scaling, both base
formulas, attack-record format, variance localization.
