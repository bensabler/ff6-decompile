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

**Next exact action:** final project-state synchronization (operator
plan item 6), then choose the next bounded unit from the queue. All
five rebalance units are done — the session delivered EXP-0021
(determinism), EXP-0022 (record cross-check), EXP-0023/GFX-0001
(graphics vertical, M2), EXP-0024/AUD-0001 (audio vertical, M3).
High-value queued follow-ups: HUD font load-path trace, SPC dispatch
trace, Fire Beam index disambiguation, MP-capable savestate, GUI/
testrunner parity check.
