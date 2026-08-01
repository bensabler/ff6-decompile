# Research Queue

## P0 — Active frontier: SCN-0001 opening-to-Whelk program

Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`.
Ordered next units (reorder when evidence exposes a better chain):

- [x] **EXP-0031: golden route segment 1** — done: deterministic
  power-on → New Game → Narshe-entry stall; milestone `00-new-game`;
  two runs byte-identical
  ([record](../docs/experiments/EXP-0031-golden-route-newgame.md)).
- [x] **EXP-0032: golden route segment 2** — done: milestones
  `01`–`03`; first scripted battle = formation 2 = two of monster
  record 0 (corrected in unit 33; Confirmed live+static)
  ([record](../docs/experiments/EXP-0032-golden-route-seg2.md)).
- [x] **EXP-0033: golden route segment 3 (partial)** — battle won on
  schedule, rewards 32 EXP / 96 GP, battle→field HP/MP writeback
  captured; free movement NOT reached (more scripted battles follow)
  ([record](../docs/experiments/EXP-0033-golden-route-seg3.md)).
- [x] **EXP-0034: segment 3b** — done: four scripted battles
  (formations 2/1/2/41), milestone `04-free-movement` at frame 46 375
  byte-identical across two runs; deferred check cleared
  ([record](../docs/experiments/EXP-0034-golden-route-seg3b.md)).
- [x] **EXP-0035: segment 4 (partial)** — route to the mines mapped
  leg by leg and walked; player tile bytes `+$00AF`/`+$00B0` and
  candidate map-id `+$1EA5` found (reading later refuted — CONTRA-0002);
  milestone 05 not claimed at the time
  ([record](../docs/experiments/EXP-0035-golden-route-seg4.md)).
- [x] **EXP-0036: segment 4 scheduled** — done: 17-leg state-driven
  route controller; battle 5 = formation 84 {27,27,0,0}; `+$1EA5`
  falsifier fired (not a simple map id); **milestone 05 established**
  over three byte-identical power-on runs
  ([record](../docs/experiments/EXP-0036-scheduled-route-to-mines.md)).
- [x] **CONTRA-0002:** `+$1EA5` map-id vs map-load-target — **both
  refuted**; it is an event-flag byte. Event-flag system located
  ([record](../docs/contradictions/CONTRA-0002-1ea5-map-id-vs-event-flags.md)).
- [x] **EXP-0037: opening event flags** — done: all writes to
  `$1E80`/`$1EA0`/`$1EC0` inventoried across the route (20 flags:
  11 latched story, 4 transient, 5 working bits), deterministic across
  one GUI + two headless runs; every writer PC statically decoded;
  16-handler family over eight bases found; event interpreter anchored
  at candidate `$C09B5C` (CEN-EVENT-0001); GUI/testrunner parity
  verified for this schedule; implemented as `internal/game/eventflags`
  ([record](../docs/experiments/EXP-0037-opening-event-flags.md),
  [DISC-0008](../docs/discoveries/DISC-0008-event-flag-system.md),
  inventory `data/scenarios/opening-event-flags.json`).
- [x] **EXP-0038: golden route segment 5** — done: milestone
  `06-random-encounter` established over two byte-identical scheduled
  runs (frame 51 307, formation 14 = three of monster record 19, leg
  19 near tile `(26,0B)`); corridor mapped; seventh verification of
  the formation table; traversal and encounters write no event flags
  ([record](../docs/experiments/EXP-0038-mines-route-to-encounter.md)).
- [ ] **EXP-0039 (next, breadth-first): mines-to-Whelk
  reconnaissance.** Visible GUI pass — sweep mines branches/dead ends,
  register every visible system and content family, advance toward
  Whelk and attempt a first visible victory. Explore widely, register
  briefly, continue forward
  ([record](../docs/experiments/EXP-0039-mines-to-whelk-breadth-recon.md)).
- [ ] `(2A,09)` mines scripted event (CEN-EVENT-0009) — candidate
  pre-Whelk beat (B17); observed once, uninvestigated.
- [ ] **Map header / tileset load path** (CEN-WORLD-0004) — still open
  and now genuinely unstarted; `+$1EA5` was a false lead. Needs its own
  discriminator (candidate: VRAM/DMA watch across the mines transition,
  which is reproducible on demand).
- [ ] Golden route segments 5..N → milestones `06-random-encounter`
  and `07-pre-whelk`, then Whelk branches A/B/C (SCN-0001 route plan).
- [ ] Milestone-01 PPU/HDMA falsifier (CEN-QUIRK-0002): capture PPU
  registers + frame buffer at frame 15000 across two runs.
- [ ] Pre-Whelk monster record extraction (B14): records 0, 19, 25,
  27, 77 — full 32-byte field map beyond the Confirmed +$08/+$0A, plus
  the monster NAME table, which is still unlocated (CEN-MONSTER-0004).
- [ ] Scripted-battle invocation opcode (CEN-EVENT-0005): exec-watch
  the frames before battle detection.
- [ ] New-game initialization capture (B01; CEN-SAVE-0001).
- [ ] Event dispatcher location (B03/B05; CEN-EVENT-0001) — largest
  single blocker; **now anchored**: every decoded event-command
  handler ends `JMP $C09B5C` with A = plausible command length
  (EXP-0037). Next: static-decode `$C09B5C` and its callers to find
  the dispatch loop and opcode table.
- [ ] Extended flag arrays (CEN-EVENT-0008 follow-up): write-watch the
  statically-decoded additional bases (`$1EE0`-`$1F5F`, `$1DC9`) on a
  future route unit; capture a save event to settle SRAM backing
  (CEN-SAVE-0001).
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
