# Current Focus

**State:** Autonomous resume session (2026-07-30 late). The battle
damage/AI investigation is **paused at a natural boundary**: EXP-0021
resolved question #30 for the tested window — action content is
deterministic given the input frame schedule (all record loads index
238, powers 13/0/19/0, miss at the same matched ordinal across
frame-exact trials); GUI-era variance attributed to harness wall-clock
jitter (Strong hypothesis). The lab now runs headless
(`--testrunner --timeout=7200`, `FF6_OUT` env, frame-scheduled input)
because the locked display breaks GUI Mesen.

**State:** the content census system is live (2026-07-31): taxonomy,
dual status ladders, ROM ownership ledger, coverage tooling in
`ff6lab`, and the observe-register-return workflow wired into the
stopping rules. EXP-0025/0026 registered the opening sequence and
censused the magic system (spell name table located, $C46AC0 shown to
be the spell database, availability array at +$1A6E). 43 census
entries, 10 of 12 domains; ROM known 0.34%.

**Next exact action:** EXP-0029 done (loader chain decoded; formation
ids stage at WRAM:+$3F46). Next: **EXP-0030 — write-watch
+$3F46/+$3F52 at an encounter entry** to capture the ROM
formation-record reader (one hop to the formation table). Rotation
alternates: HUD font load-path (GFX), SPC dispatch (AUDIO).
