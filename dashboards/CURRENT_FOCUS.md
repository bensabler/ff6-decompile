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

## ATB research program — complete (EXP-0041..0047, 2026-08-01)

The blocker EXP-0040 raised is **discharged**. Seven bounded units built
a usable ATB model; the records carry the evidence, this is the summary.

### Configuration

| Byte | Contents | Default |
|---|---|---|
| `WRAM:+$1D4D` | bits 0-2 Bat.Speed (0-5 → 1-6), bit 3 **Bat.Mode** (1 = Wait), bits 4-6 Msg.Speed, bit 7 Cmd.Set | `$2A` |
| `WRAM:+$1D4E` | bit 4 Reequip, 5 Sound, 6 Cursor, 7 Gauge | `$00` |
| `WRAM:+$1D54` | bit 7 Controller | `$00` |

Not contiguous. Cleared bit = the left-hand screen option. The Config
screen marks the **selected** option with tile attribute `$20` and the
unselected with `$28` — the inverse of the intuitive reading, and exactly
why EXP-0040 misread `Bat.Mode`. Not SRAM-backed before a save.
(EXP-0041, CEN-MENU-0007.)

### Battle entry

`ROMCPU:$C22472` reads the two config bytes **once each** and decomposes
them: Bat.Mode → `+$3A8F`, Bat.Speed → `+$3A90` (`255 − 24 × speed`),
Cmd.Set → `+$2F2E`, Gauge → `+$2021`, `+$1D4E` bits 0-2 → `+$2F34`.
Neither timing setting is re-read during the battle; `Msg.Speed` and
`Cursor` **are**, by `$C198AC` and `$C159D6` (presentation only).

**Staging rule:** ACTIVE/WAIT and Battle Speed must be set *before*
battle entry, or injected at `+$3A8F`/`+$3A90`. (EXP-0042,
CEN-BATTLE-0010.)

### The ATB layer

| Address | Role |
|---|---|
| `WRAM:+$3AB4` | **ATB gauges** — 10 entries, **stride 2**, 16-bit |
| `WRAM:+$3AC8` | per-slot increment; gauge += `$3AC8,X >> 1` at `$C21195` |
| `WRAM:+$3AA0` | scheduler flags — bit 3 gates the `+$3218` increment, **bit 6 = pending action**, bit 7 cleared on completion |
| `WRAM:+$3A3E` | 16-bit battle tick counter |
| `WRAM:+$3218` | second accumulator; `+$0100` on action completion |

**Battle Speed scales enemy gauges only** — the `CPX #$08 / BCC` at
`$C209F6` skips the `$3A90` multiply for party slots, and party
increments were byte-identical at Bat.Speed 3 and 6 while enemy
increments went 240 → 156. Watch the stride: the ATB family is **2**, not
the `$14` of the HP/stat family. (EXP-0043, CEN-BATTLE-0011.)

### ACTIVE/WAIT

`ROMCPU:$C21124` — `LDA $2F41 / AND $3A8F / BNE` — skips the entire
per-frame battle update. `+$2F41` is the **battle submenu flag**, cleared
per-frame at `$C17A92` and raised at `$C17C01`.

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| **Main command window open** | advances | **advances** |
| Ability list / Item / target selection | advances* | **paused** |
| Action resolving, victory presentation | advances* | **advances** |
| Magic / Row / Defend | unreachable from a Magitek battle | |

\* structurally implied — the gate is an `AND`; verified directly for the
ability-list row.

**The pause is narrower than the folk model**: the command window is on
screen for most of a WAIT battle and does not pause, and neither does
action resolution. (EXP-0044/0045, CEN-BATTLE-0012.)

### Action execution

The execution path (`ROMCPU:$C207xx`-`$C209xx`, completion write
`$C201BE` `INC $3219,X`) is **periodic and ungated** — it fires roughly
every 100-120 frames regardless of the gate, sweeping the battle slots.
So work pending when the gate engages still completes, and the
78/119/122-frame delays are just the wait for the next invocation.

This **vindicates EXP-0040**: "queued actions resolved out of issue
order" while menus were open was real system behaviour, not operator
error. (EXP-0046/0047, CEN-BATTLE-0013.)

### Consequence for Whelk

**No longer blocked by an absent model.** EXP-0040's timing can be
*scoped*: intervals inside the ability list and target selection were
paused; intervals at the command window and during action resolution were
not. Whelk is a boss with its own script and untested here — the one
caveat for any re-run.

## Next exact action

**EXP-0048 — name the invoker of the execution path.** Stack archaeology
has failed twice (`$C20EB6` is not a call frame; `$C20016` holds a
plausible `JSR` yet never executes), so switch instrument: exec-watch
outward from the confirmed `$C2141D`, or trace a single invocation.
Narrow question — one routine, ~every 100 frames, confirmed sites to walk
from.

**Method note:** a `JSR` at `return − 3` is necessary but not sufficient
to confirm a stack frame. Confirm by execution.

Then the programme's remaining questions — increment formula and
threshold, queue ordering and arbitration, status modifiers, other battle
types — are all non-blocking, and the Whelk decision (reinterpret
EXP-0040's scoped captures, or re-run with the model in hand) is an
ordinary orchestration call.
