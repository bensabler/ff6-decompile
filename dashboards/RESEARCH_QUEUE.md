# Research Queue

## P0 — Active frontier: SCN-0001 opening-to-Whelk program

Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`.
Ordered next units (reorder when evidence exposes a better chain):

- [x] **EXP-0031: golden route segment 1** — done: deterministic
  power-on → New Game → Narshe-entry stall; milestone `00-new-game`;
  two runs byte-identical
  ([record](../docs/experiments/EXP-0031-golden-route-newgame.md)).
- [x] **EXP-0032: golden route segment 2** — done: milestones
  `01`–`03`; first scripted battle = formation 2 / monster 12
  (Confirmed live+static)
  ([record](../docs/experiments/EXP-0032-golden-route-seg2.md)).
- [ ] **EXP-0033: golden route segment 3** — Guard battle to victory
  (victory/reward processing) → `04-free-movement`.
- [ ] Golden route segments 4..N → milestones `05`–`07-pre-whelk`,
  then Whelk branches A/B/C (SCN-0001 route plan).
- [ ] Milestone-01 PPU/HDMA falsifier (CEN-QUIRK-0002): capture PPU
  registers + frame buffer at frame 15000 across two runs.
- [ ] Monster record 12 extraction + on-screen cross-check (B14;
  `ROMFILE:0x0F0180`).
- [ ] Scripted-battle invocation opcode (CEN-EVENT-0005): exec-watch
  the frames before battle detection.
- [ ] New-game initialization capture (B01; CEN-SAVE-0001).
- [ ] Event dispatcher location (B03/B05; CEN-EVENT-0001) — largest
  single blocker.
- [ ] Map/transition/collision data (B06/B08–B11; CEN-WORLD-0001..0004).
- [ ] `+$11E0` encounter-roll producer (B12; CEN-WORLD-0006).
- [ ] Encounter packs + opening monster set, names, rewards, AI
  (B13/B14).
- [ ] Magitek command set (B15; CEN-MAGIC-0001, Fire Beam record).
- [ ] Whelk formation/records/AI/state machine (B17/B18).
- [ ] Opening graphics + audio families (SCN-0001 domain lists).

## P0 — Prior frontier (closed)

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
- [x] EXP-0024 / AUD-0001: battle confirm SFX — trigger chain Confirmed
  (press → `$21`@`$2140` from `$C117CC` → DSP voice 7 → SRCN 5 →
  `ARAM:$48D8`); SFX pack `ARAM:$4800-$491F` byte-identical to
  `ROMFILE:0x051EC9` (288 B); `brr.Decode` + `ff6lab brr info` proof.
  Follow-ups queued: SPC dispatch, `$E4/id/$18` background protocol,
  remaining pack cues.

## Completed (see ACTIVITY_LOG and indexes/EXPERIMENTS.md)
Battle damage pipeline end to end (EXP-0001..0019): engine, display
chain, accumulator, element response, defense scaling, both base
formulas, attack-record format, variance localization.
