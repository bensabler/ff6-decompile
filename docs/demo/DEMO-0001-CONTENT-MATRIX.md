# DEMO-0001 route content dependency matrix

- **Milestone:** [DEMO-0001](DEMO-0001-new-game-to-whelk.md)
- **Scenario record:** [SCN-0001](../scenarios/SCN-0001-opening-to-whelk.md)
- **Readiness matrix:** [DEMO-0001-READINESS.md](DEMO-0001-READINESS.md)
- **Created:** 2026-08-02 (Unit 10 — records reconciliation)

## What this file is for

The readiness matrix answers *"what is the status of requirement F7?"*. It is
organized by subsystem, and it is the critical-path instrument that unit
selection reads.

This file answers the other question: *"what content does beat B10 need before
it can be played?"*. It is organized by the **19 SCN-0001 route beats**, and it
exists because the demo is built in route order even though research is done in
evidence order. A subsystem can be `Integrated` in readiness while the beat that
consumes it is still unplayable because a different dependency is missing.

**This file introduces no facts.** Every recovery state and every anchor here is
sourced from an experiment, a discovery, a census entry, or the readiness
matrix. Where readiness and this file could disagree, readiness wins and this
file is wrong and must be corrected.

## Recovery state vocabulary

Readiness tracks *demo* status. This file tracks *recovery pipeline* position,
which is a different axis: a dependency can be `Format Recovered` and still sit
at `Unknown` in readiness if no demo row consumes it yet.

| State | Meaning |
|---|---|
| `Unknown` | No location, no format, no records |
| `Located` | A ROM or ARAM address is identified, format not decoded |
| `Runtime Captured` | Observed in a preserved capture (savestate, VRAM/CGRAM/OAM/ARAM dump, log) |
| `Source Provenance Known` | The captured bytes are traced back to their ROM source |
| `Format Recovered` | Structure decoded and recorded with a falsifier |
| `Extractor Implemented` | Deterministic extractor exists in `internal/extract` |
| `Private Asset Generated` | Asset lands in `local_artifacts/archive/` with manifest provenance |
| `Loaded by Go` | A typed loader in `internal/content` reads it |
| `Integrated` | Reachable in the running demo |
| `Validated` | Compared against preserved emulator evidence, difference recorded |
| `Blocked` | Cannot advance until a named dependency resolves |

## Dependency register

Each distinct content dependency has an ID. Beats reference these IDs rather
than restating them, so a state change is edited in exactly one place.

### Compression

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-CMP-01 | FF6 compression format | X1 | none — zero records, zero code | `Unknown` | **The keystone.** Gates CD-MAP-02, CD-SPR-01/04/05/06, CD-BAT-04. Next: differential recovery against preserved VRAM |

### Maps and world

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-MAP-01 | Map header records | F1 | none; CEN-WORLD-0004 | `Unknown` | Location unknown |
| CD-MAP-02 | Tileset graphics | F1 | CEN-WORLD-0001/0002 observed only | `Unknown` | CD-CMP-01 |
| CD-MAP-03 | Tilemap layers | F1 | none | `Unknown` | CD-MAP-01 |
| CD-MAP-04 | Map palettes | F2 | `bgr555.Decode` implemented; no palette table located | `Unknown` | CD-MAP-01 |
| CD-MAP-05 | Collision data | F3 | player tile bytes `WRAM:+$00AF`/`+$00B0` (EXP-0035, CEN-WORLD-0007, strong hypothesis) | `Located` (position only) | Collision *data* and lookup unlocated |
| CD-MAP-06 | Map exits and transitions | F7 | milestone 05 reproducible 3×; `+$1EA5` lead refuted (CONTRA-0002) | `Runtime Captured` | Load path unlocated |
| CD-MAP-07 | Camera and scrolling | F5 | CEN-QUIRK-0002 (HDMA/PPU-phase nondeterminism at milestone 01) | `Unknown` | Scroll register driver unlocated |
| CD-MAP-08 | Layer ordering and priority | F1 | none | `Unknown` | CD-MAP-01 |

### Events

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-EVT-01 | Event interpreter dispatch | F8 | CORR-0001: `ROMCPU:$C09B5C` advances `WRAM:+$00E5..+$00E7` by `A & $FF`, Confirmed 24/24 | `Format Recovered` (advance only) | Dispatch predecessor at `ROMCPU:$C09B59` unresolved |
| CD-EVT-02 | Opcode jump table | F8 | `ROMCPU:$C098C4`; first 64 entries plausible bank-`$C0` pointers spanning `$C09C44-$C0A336` | `Located` (static-only, never exercised) | Never dereferenced under observation |
| CD-EVT-03 | Opcode command lengths | F8 | `A & $FF` per CORR-0001 — strong hypothesis, value never seen dereferenced | `Located` | Needs an opcode→length table |
| CD-EVT-04 | Event flags | F9 | DISC-0008: three arrays, decoder `ROMCPU:$C0BAED`; opening touches exactly 20 flags | `Loaded by Go` (`internal/game/eventflags`) | Not wired into progression; flag *meanings* Unknown by policy |
| CD-EVT-05 | Scripted movement | F4 | route traversal reproducible (EXP-0036) | `Runtime Captured` | Movement routine unlocated |
| CD-EVT-06 | Scripted battle invocation | F13 | EXP-0034: exactly 4 scripted battles, formations 2, 1, 2, 41 | `Runtime Captured` | Invocation opcode unknown (CEN-EVENT-0005) |
| CD-EVT-07 | Event trigger placement | F1, F8 | contact triggers observed at named tiles (EXP-0036/0039) | `Runtime Captured` | Placement data lives in map records — CD-MAP-01 |

### Text and dialogue

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-TXT-01 | Dialogue corpus and pointer tables | F10 | lead only: banks `$CA`/`$CC` near `ROMCPU:$CC9A55` (CORR-0001 delta chain) | `Unknown` | Storage, pointers, encoding all unlocated |
| CD-TXT-02 | Dialogue control codes | F10 | none | `Unknown` | CD-TXT-01 |
| CD-TXT-03 | Variable-width dialogue font | T4 | CEN-GFX-0004 observed only | `Unknown` | ROM location, widths table |
| CD-TXT-04 | Dialogue window and border tiles | T5 | CEN-MENU-0003 observed | `Unknown` | ROM location |
| CD-TXT-05 | Pagination and advancement | F11 | none | `Unknown` | CD-TXT-01, CD-TXT-03 |
| CD-TXT-06 | Fixed-width HUD font tiles | T1 | GFX-0001, EXP-0023, ROM-0016 `ROMFILE:0x047FB0-0x048FBF` | `Validated` | — |
| CD-TXT-07 | Fixed-width text encoding | T2, T3 | EXP-0049: `VRAM tile = $100 + encoded byte` | `Integrated` | 47 tiles deliberately unidentified |

### Sprites and animation

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-SPR-01 | Field sprite tiles | F6 | CEN-GFX-0002 observed, no OAM capture | `Unknown` | CD-CMP-01 |
| CD-SPR-02 | OAM / metasprite layout | F6 | none | `Unknown` | No OAM has been read out of any capture yet |
| CD-SPR-03 | Walking frames and animation timing | F4, F6 | none | `Unknown` | CD-SPR-01 |
| CD-SPR-04 | Party battle sprites | B16 | none | `Unknown` | CD-CMP-01 |
| CD-SPR-05 | Enemy graphics | B15 | CEN-MONSTER-0003 (green guard, observed) | `Unknown` | CD-CMP-01 |
| CD-SPR-06 | Whelk head and shell graphics | B15 | head/shell visually classifiable at 4× (EXP-0040) | `Unknown` | CD-CMP-01, CD-BAT-03 |

### Battle

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-BAT-01 | Formation records | B7 | EXP-0030, ROM-0024 `ROMCPU:$CF6200`, 15-byte, verified 7× | `Integrated` | Tail `+$08..+$0E` undecoded |
| CD-BAT-02 | Monster stat records | B8 | EXP-0028/0029, ROM-0022 `ROMCPU:$CF0000`, 32-byte | `Integrated` | 20+ of 32 bytes unmapped |
| CD-BAT-03 | Monster names | B9 | none | `Unknown` | Name table unlocated — deviation D2 |
| CD-BAT-04 | Battle backgrounds | B14 | CEN-GFX-0003 observed | `Unknown` | CD-CMP-01 |
| CD-BAT-05 | Battle HUD composition | B17 | EXP-0049 decoded the layout (CEN-BATTLE-0014) and the five gauge tiles (CEN-GFX-0005) | `Integrated` (enemy side) | Party side needs CD-INI-02 — deviation D4 |
| CD-BAT-06 | Command menus and target selection | B11 | EXP-0040: Magitek sets are character-specific | `Unknown` | Command record table unlocated |
| CD-BAT-07 | Damage numerals | B17 | none | `Unknown` | Numeral tiles and placement unlocated |
| CD-BAT-08 | Action animations and effects | **no readiness row** | none | `Unknown` | Gap surfaced by this matrix — see below |
| CD-BAT-09 | Enemy AI | B12 | none | `Unknown` | AI script format unlocated |
| CD-BAT-10 | Victory, rewards, battle exit | B18 | EXP-0033: 32 EXP / 96 GP, writeback `ROMCPU:$C2496E`/`$C24979` | `Format Recovered` (writeback only) | Reward fields in monster records unmapped |
| CD-BAT-11 | ATB gauges, increments, ACTIVE/WAIT | B4, B5, B6 | EXP-0041…0047 | `Integrated` | Increment formula (`ROMCPU:$C209F0`) — deviation D3 |
| CD-BAT-12 | Encounter triggering, zones, packs | F12 | EXP-0038 reproducible at frame 51 307, formation 14; fixed-tile triggering refuted (EXP-0039) | `Runtime Captured` | `WRAM:+$11E0` producer unknown |
| CD-BAT-13 | Whelk formation and monster ids | B13 | formation 432; shell slot 4 = 50000 HP, head slot 5 = 1600 HP | `Blocked` | Monster-id extension undecoded; leading word `+$00/+$01` refuted |

### Audio

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-AUD-01 | Music sequences | A3 | `docs/audio/SEQUENCES.md` is a stub | `Unknown` | Format, location, track-id mapping |
| CD-AUD-02 | SPC700 driver and upload protocol | A4 | AUD-0001: `$21` → `$2140` from `ROMCPU:$C117CC` | `Located` (one command) | Driver dispatch unanalyzed |
| CD-AUD-03 | Sample directory, SRCN → sample | A2 | AUD-0001: one cue → SRCN 5 → `ARAM:$48D8`; pack ROM-0017 extracted | `Private Asset Generated` (one pack) | Directory table unlocated |
| CD-AUD-04 | Cue selection (music + SFX triggers) | A5 | one cue Confirmed; `$E4/id/$18` triples uninterpreted | `Located` | Command encoding unknown |
| CD-AUD-05 | Audio output and mixing | A6 | `internal/audio/brr` implemented and fuzz-tested | `Loaded by Go` (decoder only) | No mixer designed |

### Initialization and persistence

| ID | Dependency | Readiness | Anchor | State | Blocker / next action |
|---|---|---|---|---|---|
| CD-INI-01 | New Game initialization | F14 | milestone 00 deterministic (EXP-0031) | `Runtime Captured` | Init write set uncaptured |
| CD-INI-02 | Character and party records | F14 | field block ~`WRAM:+$1600` located (EXP-0027, strong hypothesis) | `Located` | Name store and initializer unlocated |

### Engine (satisfied)

`CD-ENG-01` framebuffer, `CD-ENG-02` scene machine, `CD-ENG-03` input,
`CD-ENG-04` archive loading, `CD-ENG-05` headless capture — all `Integrated`
(readiness E1–E6). Every beat depends on these; they are not repeated per beat.

## Beat → dependency map

Beat anchors are the concrete facts SCN-0001 has established. "Playable when"
names the dependencies that must reach `Integrated` for the beat to run in the
demo.

| Beat | Route anchor (established) | Dependencies | Playable when |
|---|---|---|---|
| B01 New Game initialization | milestone `00-new-game`, frame 5200, two runs byte-identical (EXP-0031) | CD-INI-01, CD-INI-02, CD-TXT-06/07 | CD-INI-01 |
| B02 Opening cinematic | auto-runs without input to first input-waiting dialogue (EXP-0031) | CD-MAP-01…04, CD-MAP-07, CD-EVT-01/02/03, CD-TXT-01/03, CD-SPR-01/03, CD-AUD-01/04 | CD-CMP-01, CD-EVT-01, CD-TXT-01 |
| B03 Snowfield march | no capture | CD-MAP-01…04, CD-EVT-05, CD-SPR-01/03, CD-AUD-01 | CD-CMP-01, CD-EVT-05 |
| B04 Character and party init | `?????`/WEDGE/VICKS observed (EXP-0025) | CD-INI-02, CD-TXT-06/07 | CD-INI-02 |
| B05 Opening dialogue, scripted movement | dialogue beat mid-sequence (EXP-0025) | CD-TXT-01…05, CD-EVT-01/05 | CD-TXT-01, CD-TXT-03 |
| B06 Entry into Narshe | milestone `02-narshe-entry`, deterministic (EXP-0032) | CD-MAP-01…06, CD-EVT-01 | CD-CMP-01, CD-MAP-06 |
| B07 Guard encounters, scripted battles | frame 31 557; four battles, formations 2, 1, 2, 41 (EXP-0034) | CD-EVT-06, CD-BAT-01/02/04/05/06/07/10, CD-SPR-04/05, CD-AUD-01/04 | CD-EVT-06, CD-SPR-04/05 |
| B08 Controllable field movement | milestone `04-free-movement`, frame 46 375 (EXP-0034) | CD-MAP-05, CD-SPR-01/02/03, CD-MAP-07 | CD-MAP-05, CD-SPR-01 |
| B09 Narshe exterior traversal | route mapped leg by leg; fifth scripted battle at tile ($1E,$27) (EXP-0035) | CD-MAP-01…08, CD-EVT-06/07 | B08 dependencies + CD-EVT-07 |
| B10 Mine entry and transition | milestone `05-mines-entry` (`$26`,`$1C`), three byte-identical runs; flag `EVF-1EA0-$2B` at frame 50 699 (EXP-0036/0037) | CD-MAP-06, CD-EVT-04/07 | CD-MAP-06 |
| B11 Mine traversal | corridor `(26,1C)`→`(26,0B)`→`(28,0B)`→`(28,09)`→`(2A,09)` (EXP-0038) | CD-MAP-01…05, CD-EVT-07 | CD-CMP-01, CD-MAP-05 |
| B12 Random encounter trigger | frame 51 307, formation 14; sets no event flags (EXP-0038) | CD-BAT-12 | CD-BAT-12 |
| B13 Encounter packs | mines pack yields formations 44 and 14 (EXP-0030/0038) | CD-BAT-12, CD-BAT-01 | CD-BAT-12 |
| B14 Monsters before Whelk | records 0, 19, 25, 27, 77 | CD-BAT-02/03/09, CD-SPR-05 | CD-BAT-03, CD-SPR-05 |
| B15 Accessible menus and commands | field menu + submenus censused (EXP-0026) | CD-BAT-06, CD-TXT-04/06/07 | CD-BAT-06, CD-TXT-04 |
| B16 Chests, save points, NPCs | golden route performs no interaction | CD-MAP-01, CD-EVT-04/07, CD-SPR-01 | Deferred by readiness F15 |
| B17 Whelk introduction | scripted beat at `(2A,09)`; contact-triggered pushing north from `(2A,07)` (EXP-0039) | CD-EVT-07, CD-TXT-01/03/05, CD-AUD-01/04 | CD-TXT-01, CD-EVT-07 |
| B18 Whelk battle | formation 432; shell slot 4 = 50000 HP, head slot 5 = 1600 HP; six head hits 162–186; shell counter confirmed behaviorally (EXP-0039/0040) | CD-BAT-13, CD-BAT-04/05/06/07/09/10, CD-SPR-04/06, CD-BAT-11, CD-AUD-01/04 | CD-BAT-13, CD-SPR-06, CD-BAT-09 |
| B19 Stable post-Whelk state | **not reached.** Milestone `10-whelk-victory` not established | CD-BAT-10, CD-EVT-04 | B18 |

## Dependency pressure

How many of the 19 beats each unresolved dependency blocks. This is the
prioritization signal the matrix exists to produce.

| Dependency | Beats blocked | Readiness |
|---|---|---|
| **CD-CMP-01 compression** | B02, B03, B06, B07, B09, B11, B14, B18 — **8** | X1 |
| CD-MAP-01 map headers | B02, B03, B06, B09, B11, B16 — 6 | F1 |
| CD-TXT-01 dialogue corpus | B02, B05, B17 — 3 | F10 |
| CD-EVT-01 event dispatch | B02, B03, B05, B06 — 4 | F8 |
| CD-SPR-01 field sprites | B02, B03, B08, B09, B16 — 5 | F6 |
| CD-AUD-01 music sequences | B02, B03, B07, B17, B18 — 5 | A3 |
| CD-BAT-13 Whelk ids | B18, B19 — 2 | B13 |

Compression leads on both counts — beats blocked and dependencies gated. That
agrees with readiness X1's own assessment and is the reason it is scheduled
first.

## Gaps this matrix surfaces

Recorded here when found, and carried into the readiness matrix so the
critical-path instrument does not lose them:

1. **CD-BAT-08, action animations and effects, has no readiness row.** The
   readiness matrix covers battle backgrounds (B14), enemy graphics (B15),
   party sprites (B16) and the HUD (B17), but nothing owns attack, spell, or
   damage-effect animation. B18's "Supports a functional Whelk battle" cannot
   pass without it.
2. **Transition effects between scenes** (fades, battle-entry wipes) are named
   by the route but owned by no requirement.
3. **CD-SPR-02, OAM layout, has never been read out of any capture.** Every
   preserved savestate contains `ppu.oamRam`; none has been examined.

Items 1 and 2 are added to readiness as new rows in the same unit that created
this file. Item 3 is an action, not a gap, and is scheduled.

## Maintenance

- A dependency changes state here only when an experiment, discovery, or
  implementation record supports it.
- Readiness is authoritative for demo status; this file is authoritative for
  route composition. Neither restates the other.
- Every completed unit that moves a dependency updates the register row, the
  affected beat rows, and the dependency-pressure table.
