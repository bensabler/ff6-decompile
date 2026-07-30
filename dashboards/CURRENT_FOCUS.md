# Current Focus

**State:** Autonomous overnight research session (2026-07-29). Units 1–4
complete: Session 003 documented; delta engine verified byte-exact
(EXP-0001); PerFrameBattleUpdate shown battle-scoped with a rate
correction and a second caller (EXP-0002); the `+$2E78` producer
identified — full battle HP chain now code-complete (EXP-0003). Mesen
running (post-"Annihilated" game-over screen).

**Next exact action:** Unit 6 — test H-BATTLE-0008 (enemy slots 4–9 in
the same arrays): reload `checkpoint3-mines.mss`, watch
`WRAM:+$3BFC`–`+$3C07` (candidate enemy-HP entries) through one
encounter while VICKS attacks; delta-engine stores with Y≥8 would
confirm. Pair with a dump of `$C20C60`–`$C20CE0` (question #13 formula
layer) if the session budget allows a seventh unit.
