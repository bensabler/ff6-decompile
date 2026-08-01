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

**Segment 3 (EXP-0033) — partial.** The first battle is won on
schedule (frame 32 706; rewards **32 EXP / 96 GP**) and victory
processing is captured: `$C2496E`/`$C24979` write the post-battle HP
(55) and MP (24) back into the field character block at `+$1609`/
`+$160D` (CEN-BATTLE-0008). **Milestone `04-free-movement` was not
reached** — a second scripted battle follows the first
(CEN-EVENT-0006), so a fixed post-victory offset cannot define it.
The two-run determinism check for segment 3 is deferred.

**Correction carried this session:** EXP-0032's first draft misread
formation record byte 1 as a monster id. Formation 2 is **two
monsters of record id 0** (live slots 6/7, HP 40 / MP 15 matching
record 0's `+$08`/`+$0A`), not one monster 12. Records, census, and
the scenario manifest are corrected with the correction preserved.

**Segment 3b (EXP-0034) — done.** The opening runs **exactly four
scripted battles** before player control — formations **2, 1, 2, 41**
(entries 31 557 / 34 953 / 36 828 / 39 500) — every staged record
matching the ROM formation table byte-for-byte. **Milestone
`04-free-movement` established at frame 46 375, byte-identical across
two power-on runs** (WRAM, screenshot, and every battle frame), which
also clears EXP-0033's deferred determinism check. New monster
reachable pre-Whelk: record **25** (27 HP / 5 MP), alongside record
**0** (40 HP / 15 MP) — CEN-MONSTER-0004.

**Segment 4 (EXP-0035) — partial.** The route from milestone 04 to
the mines is **mapped leg by leg** and the mines interior was
reached, but the scheduled encoding and determinism check are not
done, so **milestone 05 is deliberately not claimed**. Two findings
beyond routing: **player tile position at `WRAM:+$00AF`/`+$00B0`**
(Strong hypothesis — clean per-axis single-step behaviour), and a
**candidate map-id byte `+$1EA5`** ($00 opening / $05 Narshe
exterior / $0D mines across eight states — Tentative, with one
unexplained tension at milestone 02). A **fifth scripted battle**
gates the route (CEN-EVENT-0007); the climb is a zigzag, so the
up-only cadence used by earlier segments cannot walk it.

**Segment 4 (EXP-0036) — route works, milestone still unclaimed.** The
17-leg **state-driven route controller** walks the whole stretch from
milestone 04 into the mines interior with no manual correction:
position targets (`WRAM:+$00AF`/`+$00B0`, overshoot-tolerant), battle
edges, and elapsed settles, with per-leg timeouts that name the
earliest divergent leg. Tracked model + tests in
`internal/scenario/route`, with a probe-sync test that stops the Lua
and Go encodings drifting apart.

Findings: the guard trigger is at (`$1E`,`$25`) and is
**contact-triggered** (EXP-0035's condensed table had dropped an
intermediate `up` step — corrected in place); **battle 5 = formation
84**, ROM-verified at `0x0F66EC`, monsters {27, 27, 0, 0}, adding
**record 27** (115 HP / 30 MP, HP anchored to a live word) to the
pre-Whelk set; position bytes are **not field-meaningful during
battle**.

**`+$1EA5` falsifier fired.** It reaches the mines value `$0D` during
the shaft dialogue **while the party is still visibly on the
exterior** and has not moved, so it is **not** a simple current-map
byte — the evidence-safe reading is a map-load target / event-state
value written by `ROMCPU:$C0B5B6`. Transition detection was moved to
the player-position jump. Confidence deliberately **not** promoted.

**Milestone 05 is NOT claimed:** the acceptance criteria require three
scheduled power-on runs reaching (`$26`,`$1C`) in the mines, and the
final 17-leg encoding has not yet completed that set.

**Next exact action:** run the 17-leg encoding three times from
power-on, byte-compare milestone-05 WRAM across them, and create
milestone 05 only if all three land at (`$26`,`$1C`) inside the mines.
