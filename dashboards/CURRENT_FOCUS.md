# Current Focus

**State (2026-07-31): the opening-to-Whelk scenario program is
active.** Operator directive: reconstruct the complete opening
scenario — New Game selection through Whelk victory and the first
stable post-battle state — as an evidence-backed vertical slice.
Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`; machine
manifest: `data/scenarios/opening-to-whelk.json`. 19 timeline beats
(B01–B19), all PARTIAL at program start; per-beat gaps and links are
in the scenario record.

**Program baseline:** formation chain Confirmed end-to-end down to
`WRAM:+$11E0` (EXP-0028..0030: monster db $CF0000, loader $C22C30,
staging +$3F46, formation table $CF6200); spell db extracted
(EXP-0027); damage pipeline implemented + regression-tested
(EXP-0001..0019). Largest scenario blockers: no power-on route, event
engine unlocated (CEN-EVENT-0001), map/collision/encounter-zone
systems unlocated (CEN-WORLD-0001..0006), nothing observed from Whelk
introduction onward.

**Lab constraint:** unattended sessions run headless Mesen
(`--testrunner --timeout=7200`, `FF6_OUT` env, frame-scheduled input);
see BLOCKERS.md.

**Golden route status:** segments **1–2 done** (EXP-0031/0032),
milestones `00-new-game` → `03-first-scripted-battle`. The route is
one probe (`mesen/probes/EXP-0032.lua`) from power-on: start+a edge
toggling at the title (2500–4200), auto-run presentation, then held
Up + A every 240 frames (31000–46000) through the dialogue/walk
chain into the Narshe-gate battle, detected by the Confirmed
`+$3B18`/`$C22800-$C22FFF` battle-init signature. **WRAM
byte-identical across two power-on runs at all four milestones**;
screens identical except milestone 01 (CEN-QUIRK-0002 — PPU-phase
falsifier queued; treat WRAM as the assertion there).

**Scenario finding:** the first scripted battle is **formation 2,
single monster id 12** — live staged `+$3F44` bytes match ROM record
2 (`ROMFILE:0x0F621E`) exactly, independently re-verifying the
EXP-0030 formation table on a scripted (non-random) encounter.

Controlled lab variables (whole program): RamPowerOnState=AllZeros,
virgin SRAM — originals backed up under `local_artifacts/backups/`.

**Next exact action:** EXP-0033 — golden route segment 3: from
milestone 03, fight the Guard battle to victory (capturing victory /
reward processing) and continue to free field movement (milestone
`04-free-movement`), two-run determinism check.
