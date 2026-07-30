# Current Focus

**State:** Autonomous overnight research session (2026-07-29). Units 1–4
complete: Session 003 documented; delta engine verified byte-exact
(EXP-0001); PerFrameBattleUpdate shown battle-scoped with a rate
correction and a second caller (EXP-0002); the `+$2E78` producer
identified — full battle HP chain now code-complete (EXP-0003). Mesen
running (post-"Annihilated" game-over screen).

**Next exact action:** Unit 7 — question #13's formula layer: dump
`ROMCPU:$C20C60`–`$C20CE0` (around the pending-delta setter `$C20C9B`)
and the caller-return neighborhoods from the EXP-0004 stack
(`$C20430`–`$C20440`, `$C20C20`–`$C20C30`); decode how the delta value
is computed and where it comes from.
