# Latest Checkpoint

**[2026-08-01 — EXP-0035: route to the mines (partial)](2026-08-01-exp0035-route-to-mines.md)**

State: scheduled deterministic route still ends at milestone 04. The
remaining path to the mines interior is now mapped leg by leg
(11 legs, zigzag climb, gated by a fifth scripted battle) but was
walked interactively, so milestone 05 is **not** claimed. Found:
player tile bytes `WRAM:+$00AF`/`+$00B0` (Strong) and candidate
map-id `+$1EA5` (Tentative). Exact next action: EXP-0036 — encode the
leg table into a scheduled probe, capture `05-mines-entry`, two-run
determinism, and test the `+$1EA5` falsifier.
