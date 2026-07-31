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

**Next exact action (operator rebalance order):** Unit 3 — one
semantic-debt item: cross-check attack records against the local ROM
via `attackdata.RecordAt` (dump the `$C46AC0` table live, locate the
Fire Beam entry, verify power 60 / element bit 0). Live MP
verification stays queued: the only battle savestate (Magitek intro)
has no MP-consuming action, so it needs a new state first. Then
Unit 4 (graphics vertical: menu font / HUD tiles) and Unit 5 (audio
vertical: cursor SFX).
