# Current Focus

**State:** Autonomous overnight research session (2026-07-29). Units 1–4
complete: Session 003 documented; delta engine verified byte-exact
(EXP-0001); PerFrameBattleUpdate shown battle-scoped with a rate
correction and a second caller (EXP-0002); the `+$2E78` producer
identified — full battle HP chain now code-complete (EXP-0003). Mesen
running (post-"Annihilated" game-over screen).

**Next exact action:** Unit 8 — question #21: exec-log
`ROMCPU:$C20C76` capturing DP `$F0`/`$F2`, Y, and stack during one
attack whose displayed damage is read from the screen; correlate the
queued amount with the display and walk the caller chain backward
toward the damage formula.
