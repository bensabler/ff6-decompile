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
- **Event-flag system located, decoded and inventoried**
  (CONTRA-0002 → EXP-0037/DISC-0008): three bit arrays at
  `WRAM:+$1E80`/`+$1EA0`/`+$1EC0` (decoder `ROMCPU:$BAED`, masks
  `$C0BAFC`/`$C0BB04`). The whole opening touches exactly **20
  flags** — 11 latched story flags, 4 transient, 5 engine working
  bits — deterministically across GUI and headless runs. The
  16-handler script-command family covers eight bases (five beyond
  the verified three, static-only), and every handler tail anchors
  the event interpreter at candidate `ROMCPU:$C09B5C`. Inventory:
  `data/scenarios/opening-event-flags.json`.
- Position bytes are **not field-meaningful during battle**.

## Controlled lab variables (whole program)

`RamPowerOnState=AllZeros`, virgin SRAM — originals backed up under
`local_artifacts/backups/`. Unattended sessions use headless Mesen
(`--testrunner`, `FF6_OUT`, frame-scheduled input); see BLOCKERS.md.

## Whelk reached (EXP-0039 breadth pass)

The scenario's furthest point to date. The mines-to-Whelk stretch is
short and linear: one corridor with a single turn, one dead-end nub, a
**bidirectional** exterior transition, a random-encounter zone drawing
formation 14, and two scripted beats bracketing the fight. Whelk is
**formation 432**, contact-triggered from `(2A,07)` after a scripted
beat at `(2A,09)`; its introduction/warning dialogue and **shell
counterattack** were observed, and the first attempt **ended in
defeat** (which captured the defeat flow, CEN-BATTLE-0007). A
pre-Whelk savestate is preserved.

Nothing in this stretch is gated behind an unsolved system — the only
obstacle to a victory is party HP management (entering at 26/19/56).

## EXP-0040 — Whelk victory attempted, not achieved (2026-08-01)

Two piloted attempts from the preserved pre-Whelk state, both under
`Bat.Mode = Wait`. **No victory**; milestone `10-whelk-victory` and B19
remain open. Confirmed along the way:

- **Whelk is two battle entities** — shell slot 4 = **50000/50000 HP**,
  head slot 5 = **1600/1600 HP** (DISC-0001 arrays, formation 432).
- **Head-only targeting works**: six measured hits (162-186 damage
  each) all reduced the head; the shell's 50000 never moved. No shell
  strike occurred, so EXP-0039's counter was not re-triggered.
- **Head/shell state is visually classifiable** at 4× upscale.
- **A field healing route exists** before contact — Tonic ×4 took the
  party from 26/19/56 to **76/77, 105/105, 106/107** at no turn cost.
- **MagiTek sets are character-specific and EXP-0039's list was
  incomplete**: Terra has eight abilities, Wedge and Vicks four.
- The **guard/Esper beat precedes Whelk** on a clean run
  (CEN-EVENT-0011 resolved).

## Blocked — no ATB model

Whelk execution is **deferred**. The project cannot yet model
ACTIVE/WAIT behavior, qualifying submenu pause states, timer domains,
or action-queue ordering, and EXP-0040 could not operate the battle
reliably without that. All of its head/shell transitions are
menu-pause-contaminated and **must not** be used as timing evidence.
See BLOCKERS.md.

## Next exact action

Start a new session, resume from the latest checkpoint, audit existing
battle infrastructure and ATB evidence, and **propose the first bounded
ATB experiment before operating Mesen**. Do not resume Whelk gameplay
first.
