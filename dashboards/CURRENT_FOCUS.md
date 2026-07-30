# Current Focus

**State:** Autonomous overnight research session (2026-07-29). Units 1–4
complete: Session 003 documented; delta engine verified byte-exact
(EXP-0001); PerFrameBattleUpdate shown battle-scoped with a rate
correction and a second caller (EXP-0002); the `+$2E78` producer
identified — full battle HP chain now code-complete (EXP-0003). Mesen
running (post-"Annihilated" game-over screen).

**Next exact action:** Unit 5 — orchestrator pick from: (a) question #19,
exec-watch `ROMCPU:$C25D26` for its caller/trigger during one attack;
(b) question #13, write-watch the pending-delta arrays `+$33E4`/`+$33D0`
to find the damage-formula layer; (c) question #9, live MP observation to
settle H-BATTLE-0004. All need a fresh battle from
`checkpoint3-mines.mss` (reload required — party was annihilated).
