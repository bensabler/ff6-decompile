# DEMO-0001 readiness matrix

- **Milestone:** [DEMO-0001](DEMO-0001-new-game-to-whelk.md)
- **Updated:** 2026-08-02 (Unit 7 — ATB layer in Go)

Every player-visible requirement of the acceptance run appears here exactly
once. This file is the demo's critical-path instrument: unit selection reads it,
and every completed unit updates it.

## Status vocabulary

| Status | Meaning |
|---|---|
| `Unknown` | No ROM location, no format, no records. Nothing to implement |
| `Researching` | A bounded unit is actively establishing the evidence |
| `Evidence Ready` | Format/behavior Confirmed and recorded; no extractor or Go code yet |
| `Extractor Ready` | Deterministic extractor exists; asset lands in the local archive with provenance |
| `Implemented` | Go code exists with tests, but is not wired into the demo executable |
| `Integrated` | Reachable in the running demo |
| `Validated` | Compared against Mesen with recorded evidence |
| `Blocked` | Cannot advance until a named dependency resolves |
| `Deferred` | Deliberately out of the current milestone |

Status reflects the **demo**, not the research. A subsystem can be Confirmed in
`docs/discoveries/` and still sit at `Evidence Ready` here.

## Engine and platform

| # | Requirement | Subsystem | Known evidence | Missing evidence | Asset / data | Go integration point | Validation | Depends on | Status | Next action |
|---|---|---|---|---|---|---|---|---|---|---|
| E1 | Executable launches, stable window | — | n/a (project-authored) | — | — | `cmd/ff6demo`, `internal/engine/ebitenhost` | launched and ran stably; ADR-0001 | E2, E3 | `Integrated` | — |
| E2 | Deterministic fixed-step loop, frame counter | — | n/a | — | — | `internal/engine.Machine` | determinism test runs the same script twice with sleeps and extra Renders interleaved | — | `Integrated` | — |
| E3 | 256×224 indexed framebuffer, integer scaling | SNES PPU output stage | CGRAM is 256 BGR555 entries | — | — | `internal/graphics/framebuf`, `ebitenhost.integerScale` | blit/clip/flip tests + fuzz; scale table test incl. degenerate sizes | — | `Integrated` | — |
| E4 | Input: SNES pad + deterministic script source | SNES joypad | `$4218` bit order | gamepad (as opposed to keyboard) unmapped | — | `internal/platform/snespad`, `ebitenhost` | edge tests; every button proven reachable from the keyboard; `-input` drives a reproducible run | E2 | `Integrated` | gamepad when the route needs it |
| E5 | Headless frame capture (no display) | — | n/a | — | — | `internal/engine/headless` | 60 committed frame-hash goldens from a synthetic font | E2, E3 | `Integrated` | — |
| E6 | Archive location + hash-verified asset loading | — | `manifests/assets.json`, `archive verify` 8/8 | — | `local_artifacts/archive/` | `internal/content` | fixture-archive tests cover all three sentinel errors | — | `Integrated` | — |
| E7 | ROMCPU↔ROMFILE address translation | HiROM mapping | `offset = cpu − 0xC00000` Confirmed 18/18 vs Mesen (CORR-0001) | mirror/lower windows unverified for this ROM | — | `internal/platform/snesaddr` | 20 table cases + 21 constant assertions + round-trip fuzz | — | `Implemented` | — |

## Text and UI

| # | Requirement | Subsystem | Known evidence | Missing evidence | Asset / data | Go integration point | Validation | Depends on | Status | Next action |
|---|---|---|---|---|---|---|---|---|---|---|
| T1 | Fixed-width font tiles | Battle HUD font | GFX-0001, EXP-0023: 257 2bpp tiles, `ROMFILE:0x047FB0-0x048FBF` (ROM-0016), uncompressed | load path (which code/DMA copies it) | `hud-font-sheet.png` | `internal/content.LoadHUDFont` | **archive-vs-ROM differential passes on all 256 glyphs** | E6, E7 | `Validated` | — |
| T2 | Glyph → character mapping | Battle HUD font | **EXP-0049 Confirmed**: `VRAMtile = 0x100 + encodedByte`, verified by 8 direct decodes *and* the whole EXP-0023 HUD tilemap decoding to coherent game text over 37 tiles. 64 characters named; `$BF` = `?` added to `textenc` | 47 non-blank tiles carry unidentified glyphs | `data/graphics/hud-font-glyphs.json` (generated, hashes only) | `internal/game/hudfont` | relation + tracked-table consistency tests | T1 | `Integrated` | identify the 47 when a capture renders them |
| T3 | Text drawing to framebuffer | — | `textenc.Encode`/`EncodeFixed` added (EXP-0049) | proportional/dialogue text is a different system (T4) | — | `content.Font.DrawString` | synthetic-font frame goldens; unverified runes never draw a ROM tile | T1, T2, E3 | `Integrated` | — |
| T4 | Variable-width dialogue font | Dialogue rendering | CEN-GFX-0004 `OBSERVED` only | ROM location, format, widths table | — | — | — | — | `Unknown` | research unit not yet scheduled |
| T5 | Dialogue window / border tiles | Menu graphics | CEN-MENU-0003 `OBSERVED` | ROM location | — | — | — | — | `Unknown` | research unit not yet scheduled |
| T6 | Menu windows, cursor | Menu graphics | EXP-0026 captured menu tilemap/CHR/CGRAM | ROM location of the tile source | — | — | — | — | `Unknown` | research unit not yet scheduled |

T1 was qualified at Unit 0 (the extractor read the wrong ROM range). Unit 1
retired that deviation; see DEVIATIONS D1.

## Battle

| # | Requirement | Subsystem | Known evidence | Missing evidence | Asset / data | Go integration point | Validation | Depends on | Status | Next action |
|---|---|---|---|---|---|---|---|---|---|---|
| B1 | Damage / healing arithmetic | Battle formulas | DISC-0002…0006, byte-exact, golden chain test | — | — | `internal/game/battle` | existing tests + differential | — | `Implemented` | wire into a battle scene |
| B2 | Attack/spell records | Attack data | DISC-0007: 14-byte records, `ROMCPU:$C46AC0` | 6 of 14 bytes unmapped; true table length | `raw/spells.json` | `internal/game/attackdata` | fuzz + ROM cross-check | E7 | `Implemented` | — |
| B3 | Ten-slot battle arrays | Battle state | DISC-0001: stride `$14`, 10 slots, ~19 arrays | — | — | `internal/game/battle.BattleSlots` | existing tests | — | `Implemented` | — |
| B4 | ATB gauges and increments | ATB | EXP-0043: gauges `WRAM:+$3AB4` stride **2**, increments `+$3AC8`, `gauge += $3AC8,X >> 1` at `ROMCPU:$C21195` | increment **formula** (`ROMCPU:$C209F0` partly decoded), readiness threshold (`$322C` undecoded), gauge reset | — | `internal/game/atb` | measured $9C→$4E advance and the observed $00B6→$0004 wrap, both reproduced; fuzz over the 16-bit add | — | `Implemented` | integrate in Unit 9 |
| B5 | ACTIVE/WAIT behavior | ATB | EXP-0044/0045: gate `ROMCPU:$C21124`, submenu flag `WRAM:+$2F41`; pause narrower than folk model | 6 of 10 pause-matrix rows unsampled | — | `internal/game/atb.Paused`, `engine.LiveScene` | the gate is an AND of bits; gated frames freeze the tick counter with the gauges | B4 | `Implemented` | integrate in Unit 9 |
| B6 | Battle config sampling at entry | Battle init | EXP-0041/0042: `WRAM:+$1D4D` bits, sampled once at `ROMCPU:$C22472` → `+$3A8F`/`+$3A90`; speed scales **enemy only** | whether other battle types sample differently | — | `internal/game/atb.DecodeConfig`, `SampleEntry` | decodes the two recorded fingerprints ($2A default, $25 from EXP-0047) | — | `Implemented` | integrate in Unit 9 |
| B7 | Formations | Formation table | EXP-0030: `ROMCPU:$CF6200` + 15×id, verified 7× | bytes `+$00/01`, `+$08..+$0E`; **monster-id high-bit extension** | `raw/formations.bin` | new `internal/content/formations` | byte-for-byte vs staged `+$3F44` | E6, E7 | `Evidence Ready` | Unit 8 |
| B8 | Monster stat records | Monster DB | EXP-0028/0029: `ROMCPU:$CF0000`, 32-byte stride; `+$08` HP, `+$0A` MP, `+$01` power | 20+ of 32 bytes; table length; **name table unlocated** | `raw/monsters.bin` | new `internal/content/monsters` | vs battle-entry WRAM | E6, E7 | `Evidence Ready` | Unit 8 |
| B9 | Monster names | Monster DB | — | **name table unlocated** | — | — | — | — | `Unknown` | research unit (queued P0 #9) |
| B10 | Action queue, ordering, arbitration | Battle scheduler | EXP-0046/0047: execution path periodic and ungated, ~100–120 frames | invoker unnamed (EXP-0048); queue ordering rules | — | — | — | B4 | `Researching` | EXP-0048 |
| B11 | Command menus, target selection | Battle UI | EXP-0040: Magitek sets are character-specific (Terra 8, Wedge/Vicks 4) | command record table unlocated; Fire Beam index ambiguous | — | — | — | T1, T3 | `Unknown` | research unit not yet scheduled |
| B12 | Enemy AI | Monster AI | — | AI script format and location unlocated | — | — | — | — | `Unknown` | research unit not yet scheduled |
| B13 | Whelk head/shell behavior | Boss AI | EXP-0039/0040: formation 432; shell slot 4 = 50000 HP, head slot 5 = 1600 HP; head-only targeting works (6 hits, 162–186 dmg); shell counterattack observed | Whelk's monster ids (blocked on B7 extension); natural head/shell timing (EXP-0040 evidence is operator-contaminated) | — | — | — | B7, B8, B12 | `Blocked` | needs B7 monster-id extension |
| B14 | Battle backgrounds | Battle graphics | CEN-GFX-0003 `OBSERVED` | ROM location, format | — | — | — | — | `Unknown` | research unit not yet scheduled |
| B15 | Enemy graphics incl. Whelk | Battle graphics | CEN-MONSTER-0003 `OBSERVED` (green guard only) | ROM location, format, compression | — | — | — | compression | `Unknown` | research unit not yet scheduled |
| B16 | Party battle sprites | Battle graphics | — | ROM location, format | — | — | — | compression | `Unknown` | research unit not yet scheduled |
| B17 | Battle HUD (names, HP, gauges, damage numerals) | Battle UI | HUD font + glyph mapping Confirmed; EXP-0049 decoded the layout (CEN-BATTLE-0014) and identified the five gauge tiles (CEN-GFX-0005) | the compose routine; per-field cell addresses; what drives gauge fill | — | new battle scene | frame goldens | T1, T3, B4 | `Evidence Ready` | Unit 9 |
| B18 | Victory detection, rewards, battle exit | Battle flow | EXP-0033: battle 1 gives 32 EXP / 96 GP; writeback `ROMCPU:$C2496E`/`$C24979` → `WRAM:+$1609`/`+$160D` | reward field layout in monster records | — | — | — | B8 | `Evidence Ready` | Unit 8+ |

## Field, event, and world

| # | Requirement | Subsystem | Known evidence | Missing evidence | Asset / data | Go integration point | Validation | Depends on | Status | Next action |
|---|---|---|---|---|---|---|---|---|---|---|
| F1 | Map headers, tilesets, tilemaps | Map system | **none** | everything: location, structure, format, compression | — | — | — | compression | `Unknown` | research unit — blocks B06, B08–B13, B16 |
| F2 | Map palettes | Map system | `bgr555.Decode` implemented | palette table location, bank selection | — | — | — | F1 | `Unknown` | research unit not yet scheduled |
| F3 | Collision, map boundaries, exits | Field engine | player tile bytes `WRAM:+$00AF`/`+$00B0` (EXP-0035, CEN-WORLD-0007, strong hypothesis) | collision data location and lookup | — | — | — | F1 | `Unknown` | research unit not yet scheduled |
| F4 | Player/NPC movement, directional animation | Field engine | route traversal reproducible (EXP-0036) | movement routine, sprite format, animation timing | — | — | — | F1, F6 | `Unknown` | research unit not yet scheduled |
| F5 | Camera / scrolling | Field engine | CEN-QUIRK-0002: HDMA/PPU-phase nondeterminism at milestone 01 | scroll registers driver | — | — | — | F1 | `Unknown` | research unit not yet scheduled |
| F6 | Field sprites, walking frames | Field graphics | CEN-GFX-0002 `OBSERVED`, no OAM capture | ROM location, format, OAM layout | — | — | — | compression | `Unknown` | research unit not yet scheduled |
| F7 | Map transitions | Map system | milestone 05 reproducible 3× byte-identical | header/tileset load path unlocated; `+$1EA5` lead **refuted** (CONTRA-0002) | — | — | — | F1 | `Unknown` | research unit — queued P0 #6 |
| F8 | Event interpreter (command subset) | Event system | CORR-0001: `ROMCPU:$C09B5C` advances 24-bit `WRAM:+$00E5..+$00E7` by `A & $FF`, Confirmed 24/24; `+$00E3` per-frame wait | value never observed dereferenced; **dispatch predecessor unresolved**; opcode table `ROMCPU:$C098C4` static-only | — | — | — | — | `Researching` | CORR-0002 at `ROMCPU:$C09B59` |
| F9 | Event flags | Event system | DISC-0008: three arrays, decoder `ROMCPU:$C0BAED`, masks `$C0BAFC`/`$C0BB04`; opening touches exactly 20 flags | flag *meanings* deliberately unassigned; 4 further array bases static-only | `data/scenarios/opening-event-flags.json` | `internal/game/eventflags` | existing tests | — | `Implemented` | wire into progression |
| F10 | Dialogue text corpus | Dialogue | **none** | storage, pointer tables, encoding, control codes | — | — | — | F8 | `Unknown` | research unit — lead: banks `$CA`/`$CC` near `ROMCPU:$CC9A55` |
| F11 | Dialogue windows, pagination, advancement | Dialogue | — | window layout, line-break rules | — | — | — | F10, T4, T5 | `Unknown` | research unit not yet scheduled |
| F12 | Encounter triggering, zones, packs | Encounters | EXP-0038: reproducible at frame 51 307, formation 14; fixed-tile triggering **refuted** (EXP-0039) | `WRAM:+$11E0` producer; zone→pack structures | — | — | — | F1 | `Unknown` | research unit — queued P0 #15 |
| F13 | Scripted battle invocation | Event system | EXP-0034: exactly 4 scripted battles, formations 2,1,2,41 | invocation opcode unknown (CEN-EVENT-0005) | — | — | — | F8 | `Unknown` | research unit not yet scheduled |
| F14 | New Game initialization | Save/init | milestone 00 deterministic (EXP-0031) | init write set uncaptured | — | — | — | — | `Unknown` | research unit — queued P0 #11 |
| F15 | Chests, save points, NPCs, objects | Field objects | B16: nothing yet | everything | — | — | — | F1 | `Deferred` | not required by the route until B16 |

## Audio

| # | Requirement | Subsystem | Known evidence | Missing evidence | Asset / data | Go integration point | Validation | Depends on | Status | Next action |
|---|---|---|---|---|---|---|---|---|---|---|
| A1 | BRR sample decode | Audio | `internal/audio/brr` full S-DSP semantics, fuzz-tested | — | — | `internal/audio/brr` | existing tests | — | `Implemented` | needs a mixer to be audible |
| A2 | Sample directory, SRCN → sample | Audio | AUD-0001: one cue → SRCN 5 → `ARAM:$48D8` | directory table location | `audio/sfx-pack.brr` | — | — | — | `Unknown` | research unit not yet scheduled |
| A3 | Music sequences | Audio | `docs/audio/SEQUENCES.md` is a stub | format, location, track id mapping, instrument tables | — | — | — | — | `Unknown` | research unit not yet scheduled |
| A4 | SPC700 driver, upload protocol | Audio | AUD-0001: `$21` → `$2140` from `ROMCPU:$C117CC` | driver dispatch entirely unanalyzed | — | — | — | — | `Unknown` | research unit not yet scheduled |
| A5 | Cue selection (music + SFX triggers) | Audio | one cue Confirmed; `$E4/id/$18` triples and `$28` heartbeat uninterpreted | command encoding | — | — | — | A4 | `Unknown` | research unit not yet scheduled |
| A6 | Audio output / mixing | — | n/a | — | — | not yet designed | cue-order comparison | A1–A5 | `Unknown` | after A3/A4 |

## Cross-cutting

| # | Requirement | Known evidence | Missing evidence | Status | Next action |
|---|---|---|---|---|---|
| X1 | FF6 compression format | **none** — never investigated, zero records, zero code | identification, algorithm, decompressor | `Unknown` | **Highest-leverage unscheduled research.** Gates F1, F6, B14, B15, B16 |
| X2 | Reproducible one-command asset generation | `ff6lab extract all` + `archive verify` 8/8 | dialogue/maps/animations/scripts categories are empty | `Extractor Ready` | extend per family as formats land |
| X3 | Whelk victory observed at all | **never achieved.** EXP-0040 lost twice; ATB blocker since discharged | a completed Branch-A run | `Blocked` | re-run with the ATB model in hand |

## Summary at Unit 0

| Status | Unit 0 | Now (Unit 7) |
|---|---|---|
| `Validated` | 0 | **1** |
| `Integrated` | 0 | **7** |
| `Implemented` | 5 | 10 |
| `Evidence Ready` | 7 | 5 |
| `Extractor Ready` | 2 | 1 |
| `Researching` | 3 | 2 |
| `Blocked` | 2 | 2 |
| `Deferred` | 1 | 1 |
| `Unknown` | 33 | 25 |

**The demo now runs, in a window.** Units 5-6 moved the first eight rows into the executable
and one to Validated — T1, whose archive-vs-ROM differential passes on all
256 glyphs.

Every Integrated row is engine or text plumbing. No row of the Battle, Field,
or Audio sections has moved, and none will until either the ATB port (Unit 7)
or the research that unblocks the field half lands. That is the honest shape
of the progress: the demo has a spine, not yet a game.
