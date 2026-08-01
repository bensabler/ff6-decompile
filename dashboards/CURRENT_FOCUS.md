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
- **Event-flag system located** (CONTRA-0002): three bit arrays at
  `WRAM:+$1E80`/`+$1EA0`/`+$1EC0`, set via `ORA $C0BAFC,X` /
  `STA $1EA0,Y` with a flag-number decoder at `ROMCPU:$BAED`
  (Y = flag/8, X = flag&7). `+$1EA5` — twice mistaken for a map
  indicator — is byte 5 of the `+$1EA0` array (flags `$28`-`$2F`).
- Position bytes are **not field-meaningful during battle**.

## Controlled lab variables (whole program)

`RamPowerOnState=AllZeros`, virgin SRAM — originals backed up under
`local_artifacts/backups/`. Unattended sessions use headless Mesen
(`--testrunner`, `FF6_OUT`, frame-scheduled input); see BLOCKERS.md.

## Next exact action

**EXP-0037 — map the opening's event flags.** CONTRA-0002 removed
`+$1EA5` as a map lead and replaced it with something larger: the
event-flag arrays and their set/clear routines. The bounded next unit
write-watches `$1E80`/`$1EA0`/`$1EC0` across the scheduled route to
record **which flag numbers the opening sets and when** — that serves
B16 (treasure/interaction flags), B19 (post-battle state) and
persistence, and gives the event engine (CEN-EVENT-0001) a concrete
anchor it has never had.

The map header / tileset question (CEN-WORLD-0004) is **still open and
untouched** — `+$1EA5` was never going to answer it, so that lead is
now genuinely unstarted and needs its own discriminator.
