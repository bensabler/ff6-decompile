# Current Focus

**Program:** SCN-0001 — reconstruct the opening scenario, New Game
through Whelk victory, as an evidence-backed vertical slice.
Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`;
machine manifest: `data/scenarios/opening-to-whelk.json`
(19 beats B01–B19, per-beat status and links live there).

## Golden route — power-on through milestone 05

One tracked probe (`mesen/probes/EXP-0036.lua`) drives power-on →
mines interior, reproducibly:

| Milestone | Frame | Established by |
|---|---|---|
| `00-new-game` | 5 200 | EXP-0031 |
| `01-opening-cinematic` | 15 000 | EXP-0032 |
| `02-narshe-entry` | 30 000 | EXP-0032 |
| `03-first-scripted-battle` | 31 677 | EXP-0032 |
| `04-free-movement` | 46 375 | EXP-0034 |
| `05-mines-entry` | 51 578 | EXP-0036 |

Milestone WRAM is byte-identical across repeated power-on runs at
every milestone (two runs each for 00–04; **three** for 05). The
route's second half is a **state-driven 17-leg controller** —
position targets on `WRAM:+$00AF`/`+$00B0`, battle edges, elapsed
settles, per-leg timeouts that name the earliest divergent leg.
Tracked model and tests: `internal/scenario/route`, including a
probe-sync test so the Lua and Go encodings cannot drift apart.

**Assertion channel is WRAM, not screenshots** — CEN-QUIRK-0002 has
now been seen at two milestones (01 and 05): identical WRAM,
non-byte-stable frame capture.

## Scenario facts established so far

- The opening runs **five scripted battles** before the mines:
  formations **2, 1, 2, 41** (before free movement) then **84** (en
  route to the shaft). Every staged `+$3F44` record matches the ROM
  formation table `$CF6200 + id×15` byte-for-byte — six independent
  verifications of EXP-0030's table.
- Pre-Whelk monsters identified so far: records **0** (40 HP/15 MP),
  **25** (27 HP/5 MP), **27** (115 HP/30 MP) from the scripted
  battles, plus **19** and **77** from the mines random encounter.
  Monster *names* remain unlocated — the name table is not found, so
  on-screen labels are observational only.
- Battle 1 rewards **32 EXP / 96 GP**, with post-battle HP/MP written
  back into the field character block (`$C2496E`/`$C24979` →
  `+$1609`/`+$160D`).
- `WRAM:+$1EA5` is **not** a simple map-id byte: it reaches the mines
  value before the transition is visible or the party moves. Recorded
  as a map-load-target / event-state candidate written by
  `ROMCPU:$C0B5B6`.
- Position bytes are **not field-meaningful during battle**.

## Controlled lab variables (whole program)

`RamPowerOnState=AllZeros`, virgin SRAM — originals backed up under
`local_artifacts/backups/`. Unattended sessions use headless Mesen
(`--testrunner`, `FF6_OUT`, frame-scheduled input); see BLOCKERS.md.

## Next exact action

**EXP-0037 — the map system.** The mines transition is now
reproducible on demand, which makes this the highest-leverage target:
write-watch `ROMCPU:$C0B5B6` across the transition to settle
`+$1EA5`'s meaning, and follow it to the map header / tileset load
path (CEN-WORLD-0004, CEN-WORLD-0007). It unlocks beats B06 and
B08–B11 and B16 — the widest remaining PARTIAL rows in the matrix —
and the event engine (CEN-EVENT-0001) is the only comparable blocker.
