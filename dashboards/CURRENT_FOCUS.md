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

**Next exact action:** EXP-0031 — golden route segment 1: reproducible
frame-scheduled route from power-on through New Game selection;
milestone `00-new-game` savestate + input transcript + determinism
check (two identical runs).
