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

**Golden route status:** segment 1 **done** (EXP-0031): deterministic
power-on route — title press (start+a edge toggling, frames
2500–4200), milestone `00-new-game` at frame 5200, auto-run to the
input-waiting Narshe-entry dialogue asserted at frame 30000; two runs
byte-identical (WRAM + screen). Controlled lab variables:
RamPowerOnState=AllZeros, virgin SRAM (original backed up). New
census: CEN-EVENT-0004 (title/attract flow), CEN-QUIRK-0001 (boot
uninitialized reads).

**Next exact action:** EXP-0032 — golden route segment 2: extend the
schedule from the Narshe-entry stall through the scripted approach
and the first scripted battle (milestones `01-opening-cinematic`,
`02-narshe-entry`, `03-first-scripted-battle`), A-press cadence per
stalled dialogue, with per-milestone WRAM assertions and a two-run
determinism check.
