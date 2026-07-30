# Current Focus

**State:** Autonomous overnight research session (2026-07-29). Units 1–4
complete: Session 003 documented; delta engine verified byte-exact
(EXP-0001); PerFrameBattleUpdate shown battle-scoped with a rate
correction and a second caller (EXP-0002); the `+$2E78` producer
identified — full battle HP chain now code-complete (EXP-0003). Mesen
running (post-"Annihilated" game-over screen).

**Next exact action:** Unit 9 — question #21 continuation: dump
`ROMCPU:$C26AE0`–`$C26B60` (deeper-return lead `$6B16` from the
EXP-0007 stacks) and decode toward the damage formula; fall back to a
full-DP-snapshot exec log at `$C20C76` if the lead is a dead end.
