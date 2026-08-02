# Current Focus

**Active program:** DEMO-0001 — produce a **playable Go demo** of New
Game through Whelk victory. Started 2026-08-02. The foundation (units
0–9) is merged to `main` at `297ba88`, tagged `demo-0001-foundation-v0.1`;
work continues on branch `demo/whelk-content-parity`. The deliverable is a
runnable program; research, extraction, and documentation are supporting
activities. Records: `docs/demo/DEMO-0001-new-game-to-whelk.md`,
`DEMO-0001-READINESS.md` (the critical-path instrument, by subsystem),
`DEMO-0001-CONTENT-MATRIX.md` (the route view, by beat),
`DEMO-0001-ACCEPTANCE.md`, `DEMO-0001-DEVIATIONS.md`.

**Evidence program:** SCN-0001 — reconstruct the opening scenario, New
Game through Whelk victory, as an evidence-backed vertical slice.
Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`;
machine manifest: `data/scenarios/opening-to-whelk.json`
(19 beats B01–B19, per-beat status and links live there).
DEMO-0001 **consumes** SCN-0001; it does not replace it.

## DEMO-0001 build order — evidence-led, not route-ordered

Decided 2026-08-02. The milestone ladder reads in route order (shell →
opening → first map → events → battle), but the project's evidence is
distributed almost inversely: a field room needed the map system and,
it was assumed, a compression format — both at zero records — while a
battle needs systems already Confirmed and partly in Go. EXP-0050 has
since removed the compression half of that assumption.

1. **Technical shell** — executable, deterministic loop, input, indexed
   framebuffer, archive-backed asset loading, headless frame capture.
2. **Battle vertical** — ATB model into Go, formation/monster decoders,
   battle HUD scene, then a playable encounter.
3. **Field/event vertical** — gated behind research into the map system,
   the dialogue corpus, and the event opcode table.

Readiness at program start was **0 of 55 requirements Integrated** — 5
Implemented, 6 Evidence Ready, 2 Extractor Ready, 3 Researching, 2
Blocked, 1 Deferred, **36 Unknown**. ROM ownership 0.49 %. (Unit 0's own
summary said "53 / 33 Unknown / 7 Evidence Ready"; those figures were
wrong and propagated. Unit 10 recounted from the table — see the
correction note in `DEMO-0001-READINESS.md`.)

After units 0–14: **14 Integrated, 1 Validated, 6 Implemented**, 29
Unknown, of 57 rows. Every Integrated row is engine, text, or
battle-table plumbing — **no Field or Audio row has moved**. The demo has
a spine, not yet a game.

The route content matrix measures the same fact by dependency pressure,
re-derived after EXP-0050: **map headers (F1) lead at 6 of 19 beats**,
then field sprites (F6) and music sequences (A3) at 5 each, event
dispatch (F8) at 4, the dialogue corpus (F10) at 3.

Compression is **withdrawn** from that table. It was ranked first at 8
beats when the matrix was created, on readiness X1's prose rather than on
evidence; EXP-0050 tested the premise and refuted it for the map tile
graphics on this route. What X1 gates is Unknown and is recorded as
Unknown.

**Defect D1 (parity-blocking) — found at program start, fixed in Unit 1.**
The `hud-font` extractor read `ROMFILE:0x046FC0` as a block start, but
that address is only the back-projected anchor for VRAM tile `$000`. The
real block is `ROMFILE:0x047FB0-0x048FBF`, which
`manifests/rom-regions.json` ROM-0016 recorded correctly the whole time —
the extractor and the ledger were two records of one fact that disagreed,
and nothing compared them. **255 of 257 tiles in the shipped
`hud-font-sheet.png` were attack-table bytes rendered as tiles**, from
2026-07-30 until 2026-08-02.

Now asserted two ways: `TestHUDFontMatchesROMLedger` compares the
extractor's span to ROM-0016 with no ROM needed, and a skip-guarded
archive-vs-ROM differential compares all 256 glyphs. The general lesson is
recorded in `ARCHITECTURE.md`: a hand-computed offset needs something that
asserts it against its other record.

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

## Event interpreter — CORR-0001 (2026-08-01)

The first static/runtime correlation landed on `ROMCPU:$C09B5C`, the shared
tail every event-flag handler jumps to. Record:
`docs/correlations/CORR-0001-C09B5C.md`.

Confirmed across 24 observations (12 hits × 2 runs, byte-identical logs):
entry state E = native, **M = 1**, **X = 0**, PBR = `$C0`, DBR = `$00`,
**DPR = `$0000`**, SP = `$15FD`. DPR was measured, so the direct-page
operands resolve to `WRAM:+$00E3` and `WRAM:+$00E5`-`+$00E7`. The 24-bit
value there increases by **exactly `A & $FF`** (deltas varied 1-7, refuting
a constant), control reaches `ROMCPU:$C09A6D` in the same frame every time,
and `WRAM:+$00E3` counts down once per frame until zero (twice, at exactly
30 frames). `$C09B5C` is a **genuine shared routine entry**, entered from
multiple command handlers — not only the flag family.

Still **Strong hypothesis, not Confirmed** — the value was never observed
being dereferenced: event-script pointer, command length, script frame-wait.
The **immediate predecessor is unresolved**.

## Next exact action

**PAUSED at an atomic boundary.** No unit is in flight. The next session
performs a **project-wide workflow and orchestration audit**; no research
or implementation unit should begin before it.

Branch `demo/whelk-content-parity`, HEAD `38a62a1`, worktree clean, no
emulator running, no background processes. Evidence for EXP-0050 through
EXP-0053 is harvested into `local_artifacts/experiments/` with READMEs,
replay scripts and verified hashes.

**The map-descriptor line is stopped by choice**, under the
three-attempts rule — six meaningfully different approaches across
EXP-0051 and EXP-0052 failed to find what selects the field tile blocks.
The instrument is wrong: this is a question for a disassembler or a
trace, not a search. What the line *did* establish is recorded and
hashed: the tile block addresses, a map identifier (`$39` Narshe / `$17`
mines at `WRAM:+$1305`/`+$13E2`/`+$1F80`), and a 33-byte record table in
bank `$ED` carrying that id at `record[28]` on all three captures.

When content work resumes, in priority order:

1. **Test the block-boundary alternative** — nearly free, needs no
   instrument, and would invalidate all four pointer encodings if it
   holds. `ff6lab state origin` anchors a matched run where the *image*
   begins, so a ROM block starting before `VRAM:$0000`'s source is
   reported from its midpoint and its true start has never been searched
   for. Probe the ROM immediately before `0x208460` and `0x20DFA0`.
2. **Diff the sprite tile region between milestones 02 and 04** — both
   are the Narshe exterior, so if the player faces differently the
   walking-frame set and the sheet's row stride fall out with no new
   capture (F6).
3. Needing one operator session: **run `probe dma-trace` over a map
   load**. Written in Unit 16, still **unexercised**. One run would name
   the source address and the routine that set it up.

Carrying forward, five instances deep now: **two records of one fact are
worth nothing unless something compares them.** The `hud-font` extractor
vs ROM-0016 (D1). The readiness summary vs its own tables (Unit 10).
`BlitOptions.PaletteBase`'s doc comment vs its callers (D0, Unit 11).
EXP-0050's "shared with the mines interior", never measured (EXP-0051).
And EXP-0053's first pass, where a 128-byte probe spanned a VRAM gap and
reported the player's sprite as absent from the ROM — a negative from
the wrong probe geometry looks exactly like a real one.

### Research actions, now sequenced behind the demo's critical path

These remain the correct next research units, and each is named in the
readiness matrix against the demo requirement it unblocks. They are no
longer the *first* action, because implementation-ready work exists.

**Name the immediate predecessor at `ROMCPU:$C09B59`.** Exec-watch the
dispatcher's `JMP ($002A)` alongside `$C09B5C`, capturing DP `$2A`/`$2B`
and `$EA` at dispatch, and correlate each entry with the opcode that
reached it. Closes CORR-0001's one gap, turns "A is a command length" into
a measured opcode→length table, and decodes the candidate opcode table at
`ROMCPU:$C098C4` as a by-product. Also the **deferred blocker** for any
demonstration claiming to read or execute event script.

Still open and unchanged, non-blocking: **EXP-0048 — name the invoker of
the execution path.** Stack archaeology has failed twice (`$C20EB6` is not
a call frame; `$C20016` holds a plausible `JSR` yet never executes), so
switch instrument: exec-watch outward from the confirmed `$C2141D`, or
trace a single invocation.

**Method note:** a `JSR` at `return − 3` is necessary but not sufficient
to confirm a stack frame. Confirm by execution.

Then the programme's remaining questions — increment formula and
threshold, queue ordering and arbitration, status modifiers, other battle
types — are all non-blocking, and the Whelk decision (reinterpret
EXP-0040's scoped captures, or re-run with the model in hand) is an
ordinary orchestration call.
