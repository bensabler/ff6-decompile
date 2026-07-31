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

**Next exact action (operator rebalance order):** Unit 4 — graphics
vertical proof (menu font / battle HUD tiles: runtime → DMA → ROM →
decoder → comparison), then Unit 5 — audio vertical proof (cursor
SFX). Unit 3 done: EXP-0022 cross-checked the attack-record table
against the live ROM (record 238 power 0 + physical flag Confirmed;
Fire Beam candidates 5/131 Tentative). Live MP verification stays
queued behind a savestate with an MP-consuming action (BLOCKERS).
