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
- [x] **EXP-0039: mines-to-Whelk breadth reconnaissance** — done:
  **Whelk reached, entered and observed** (formation 432,
  contact-triggered from `(2A,07)`); corridor branches/dead ends
  mapped; mines transition found **bidirectional**; fixed-tile
  encounter triggering **refuted**; Magitek ability list captured;
  first attempt lost to the shell counter, capturing the defeat flow
  ([record](../docs/experiments/EXP-0039-mines-to-whelk-breadth-recon.md)).
- [~] **EXP-0040: Whelk victory attempt, branch A — STOPPED, not
  achieved.** Two attempts from the preserved pre-Whelk state, both
  under `Bat.Mode = Wait`. Confirmed that Whelk is **two battle
  entities** (shell slot 4 = 50000 HP, head slot 5 = 1600 HP), that
  striking the visibly extended head damages the head and never the
  shell (six measured deltas, 162-186 each), that head/shell state is
  visually classifiable, and that a **field healing route** exists
  (Tonic ×4) taking the party to 76/105/106. **No victory**; milestone
  `10-whelk-victory` and B19 remain open. Stopped by operator directive
  on a **methodological blocker** — see below
  ([record](../docs/experiments/EXP-0040-whelk-victory.md)).
- [~] **P0 — bounded ATB research program (blocks all further Whelk
  execution).** Started 2026-08-01. The project still has no model of
  ACTIVE/WAIT semantics, qualifying submenu pause states, timer domains,
  or action-queue ordering. **Whelk gameplay must not resume before this
  program produces a usable model.** Progress below.
- [x] **EXP-0041: battle configuration storage** — done: all nine Config
  settings located and bit-decoded across `WRAM:+$1D4D`, `+$1D4E` and
  `+$1D54` (the block is not contiguous); both speed fields swept to
  their clamps; the Config screen's selected/unselected attribute
  convention settled (`$20`/`$28`), independently confirming EXP-0040's
  `Bat.Mode = Wait` correction; configuration shown not SRAM-backed
  before a save. Prerequisites landed with it: an audited
  battle-configuration fingerprint field, and `ff6lab state` for reading
  preserved `.mss` files without an emulator
  ([record](../docs/experiments/EXP-0041-battle-config-storage.md)).
- [x] **EXP-0042: battle-entry configuration sampling** — done, answered
  **mixed**. `Bat.Mode` and `Bat.Speed` are read **once** at battle entry
  by `ROMCPU:$C22472` and decomposed into `WRAM:+$3A8F` (Wait flag) and
  `WRAM:+$3A90` (`255 − 24 × speed`), never re-read for timing;
  `Msg.Speed` and `Cursor` *are* read live by `$C198AC`/`$C159D6`
  (presentation only). Both derived values were predicted from the
  disassembly and matched. **Staging rule established:** ACTIVE/WAIT and
  Battle Speed must be set before battle entry, or injected at
  `+$3A8F`/`+$3A90`
  ([record](../docs/experiments/EXP-0042-battle-entry-config-sampling.md)).
- [x] **EXP-0043: ATB gauges and the `+$3A90` consumer** — done. **Gauge
  array `WRAM:+$3AB4`** (10 entries, **stride 2**, 16-bit), advanced by
  `$3AC8,X >> 1` at `ROMCPU:$C21195`; per-slot increment `WRAM:+$3AC8`
  built at `ROMCPU:$C209E0` from Speed at `+$3B19,X`; scheduler flags
  `WRAM:+$3AA0`; battle tick counter `WRAM:+$3A3E`. `+$3A8F`'s consumer
  is **`ROMCPU:$C21124`** — `LDA $2F41 / AND $3A8F / BNE` gates the whole
  per-frame update, so that is where ACTIVE and WAIT diverge.
  **Battle Speed scales enemy gauges only** (party increments identical
  at Bat.Speed 3 and 6; enemy 240 vs 156)
  ([record](../docs/experiments/EXP-0043-atb-gauges-and-speed-consumer.md)).
- [x] **EXP-0044: the ACTIVE/WAIT pause matrix** — done. `WRAM:+$2F41`
  is the **battle submenu flag** (cleared per-frame at `ROMCPU:$C17A92`,
  raised at `ROMCPU:$C17C01` on submenu open); `ROMCPU:$C21124` ANDs it
  with the Wait flag and skips the whole per-frame battle update. All
  four located domains freeze and resume **together**. The pause is
  **narrower than assumed**: the main command window does **not** pause
  under WAIT, and neither do action animations — only the ability list
  and target selection did. Verified by an in-place one-variable flip of
  `+$3A8F` and cross-checked against genuinely configured WAIT
  ([record](../docs/experiments/EXP-0044-active-wait-pause-matrix.md)).
- [ ] **P0 NEXT — EXP-0045: finish the matrix, settle the transient.**
  Six matrix rows read `not sampled` (Item, Magic, Row, Defend, dialogue,
  damage display, victory, defeat) — walk each sampling `+$2F41`. Then
  frame-step the gate transition with an `endFrame` callback to decide
  whether the settling transient is queued work resolving past the gate
  or the last un-gated frame; bridge round-trips are far too coarse.
- [ ] **ATB program, remaining (none blocking):** the exact increment
  formula, threshold and gauge reset; `+$3AA0` bit semantics; the action
  queue and readiness arbitration; status modifiers (Haste/Slow/Stop);
  battle types other than random encounters — Whelk is a boss with its
  own script, the one caveat for any Whelk re-run.
- [ ] **Whelk (B17/B18/B19, milestone `10-whelk-victory`) — unblocked.**
  The ATB model exists. Decide whether to reinterpret EXP-0040's scoped
  captures or re-run the fight with the model in hand.
- [ ] Whether battle types other than random encounters sample
  configuration differently (scripted, pincer, back, boss) — cheap to
  re-check at the first scripted battle rather than as its own unit.
- [ ] Whelk formation monster-id extension: reading formation 432's
  record bytes 2-7 as ids yields record 0 (the opening guard), so the
  id field must carry a high-bit/extension not yet decoded (FF6
  exceeds 256 monsters). Blocks Whelk's monster records (B18/B14).
- [x] Post-Whelk guard/Esper beat (CEN-EVENT-0011) — **resolved by
  EXP-0040**: the "We won't hand over the Esper!!" beat fires at
  `(2A,07)` **before** the Whelk battle on a clean, never-defeated run,
  reproduced on both attempts. It is normal pre-Whelk progression, not
  a post-defeat artifact.
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
