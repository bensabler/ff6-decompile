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

**Next exact action:** EXP-0027 done (spell database extracted: 54
records Confirmed, cost field Confirmed, esper name table found,
field character-record block located). Per the prioritization rules,
three consecutive magic-adjacent units have now run (EXP-0025/26/27)
— **review `ff6lab coverage summary` and rotate domains.** Top
candidates: monster stat-record source trace (battle-init writes into
the +$3B18 family — unlocks the monster database the way EXP-0026
unlocked spells), the HUD font load-path trace (completes GFX-0001),
or the SPC dispatch trace (completes AUD-0001).
