# DEMO-0001: New Game to Whelk Victory — playable Go demo

- **Status:** ACTIVE (program started 2026-08-02)
- **Kind:** Production milestone (deliverable is a runnable program, not a record)
- **Scenario record:** [SCN-0001](../scenarios/SCN-0001-opening-to-whelk.md)
- **Machine manifest:** `data/scenarios/opening-to-whelk.json` (19 beats B01–B19)
- **Readiness matrix:** [DEMO-0001-READINESS.md](DEMO-0001-READINESS.md) (by subsystem)
- **Route content matrix:** [DEMO-0001-CONTENT-MATRIX.md](DEMO-0001-CONTENT-MATRIX.md) (by beat)
- **Acceptance criteria:** [DEMO-0001-ACCEPTANCE.md](DEMO-0001-ACCEPTANCE.md)
- **Deviations register:** [DEMO-0001-DEVIATIONS.md](DEMO-0001-DEVIATIONS.md)
- **ROM revision:** SHA-256 `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
- **Branch:** `demo/whelk-content-parity` (foundation merged to `main` at `297ba88`, tagged `demo-0001-foundation-v0.1`)

## Deliverable

A standalone Go program that launches from a documented command, begins at New
Game, and plays continuously through the defeat of Whelk to the first stable
post-battle state.

The acceptance run must not require Mesen or Ghidra to be running. Those remain
research and validation tools only.

Boundary is inherited verbatim from SCN-0001:

```text
Start:
New Game is selected from a fresh boot or reproducible equivalent state.

End:
Whelk has been defeated, battle rewards and victory processing have
completed, and the game has reached the first stable post-battle state
before the frozen Esper interaction begins.
```

## Relationship to SCN-0001

SCN-0001 is the **evidence** program: what the original does, established by
controlled experiment. DEMO-0001 is the **production** program: a Go
reconstruction that reproduces it. They share the boundary, the 19 beats, and
the ROM revision.

DEMO-0001 consumes SCN-0001; it does not replace it. A beat is not demo-complete
until SCN-0001 carries the evidence *and* the readiness matrix carries an
`Integrated` or `Validated` row.

## Honest starting position (2026-08-02)

This is recorded so no later reader mistakes scaffolding for parity.

The project's confirmed knowledge is distributed almost inversely to the demo's
route order. At program start:

| Area | State |
|---|---|
| Battle damage pipeline | Confirmed, in Go, tested (DISC-0001…0007) |
| ATB gauges, increments, ACTIVE/WAIT gate | Confirmed, documented, **not yet in Go** (EXP-0041…0047) |
| Formations `ROMCPU:$CF6200`, monsters `ROMCPU:$CF0000`, attacks `ROMCPU:$C46AC0` | Confirmed |
| Event-flag arrays | Confirmed, in Go (DISC-0008) |
| Battle HUD font `ROMFILE:0x047FB0` | Confirmed (GFX-0001) |
| Map headers, tilesets, tilemaps, collision, exits, map palettes | **Unknown — no ROM location** |
| Dialogue text, pointer tables, control codes, variable-width font | **Unknown** |
| Event bytecode opcode semantics | Candidate table `ROMCPU:$C098C4`, static-only, never exercised |
| Field / battle / enemy / Whelk sprites, battle backgrounds | **Unknown** |
| Music sequences, SPC driver, sample directory | **Unknown** (one 288-byte SFX pack extracted) |
| FF6 compression format | **Never investigated** — zero records, zero code |
| Whelk victory | **Never observed.** EXP-0040 attempted twice and lost |

ROM ownership at program start: **0.49 % known** (`indexes/ROM_REGIONS.md`).

DEMO-0001 is therefore a multi-session program. Progress is measured by the
readiness matrix, not by elapsed units.

## Build order — evidence-led, not route-ordered

The natural reading of the milestone ladder is route order: shell → opening →
first map → events → battle. The project's evidence does not support that order.
A field room requires the map system and a compression format, both at zero; a
battle requires systems that are already Confirmed and partly implemented.

Adopted order, decided 2026-08-02:

1. **DEMO-0001A** — technical shell. Executable, deterministic loop, input
   abstraction, framebuffer, archive-backed asset loading, headless frame
   capture. Renders real FF6 text.
2. **Battle vertical** — ATB model into Go, formation/monster decoders, battle
   HUD scene, then a playable encounter.
3. **Field/event vertical** — gated behind research into the map system,
   dialogue text, and the event opcode table.

Every milestone builds the same integrated executable forward. Route order is
filled in as research lands, not skipped.

## Executable

```bash
go run -tags gui ./cmd/ff6demo
```

Headless frame capture, for automated comparison and CI — no display needed,
and the authoritative mode for any frame-exact claim:

```bash
go run ./cmd/ff6demo -headless -frames 120 -capture-last
```

The `gui` tag is required because Ebitengine needs cgo on macOS and Linux;
keeping it optional is what lets every gate run without system libraries
([ADR-0001](../decisions/ADR-0001-rendering-host.md)).

The demo binary **structurally cannot read a ROM**: `internal/rom` is forbidden
— **transitively** — from `cmd/ff6demo`, `internal/content`, and
`internal/engine`, and `ff6lab audit` walks the import graph to enforce it.
Transitivity is not a detail: `internal/content` reached `internal/rom` through
`internal/extract` until the manifest model was split into
`internal/assetmanifest`, and a direct-import check called that clean.

All game data is loaded from the locally generated archive under
`local_artifacts/archive/`, which the operator produces from their own verified
ROM:

```bash
export FF6_ROM=/path/to/your/verified.sfc
go run ./cmd/ff6lab extract all
```

## Unit log

| Unit | Result | Commit |
|---|---|---|
| Unit 0 | DEMO-0001 program records established | `969b5dd` |
| Unit 1 | HUD font extractor read the wrong ROM range; fixed, D1 retired | `cdc55e8` |
| Unit 2 | `internal/platform/snesaddr`; existing offsets asserted through it | `887f676` |
| Unit 3 | EXP-0049 — the encoding indexes the font block directly | `aabad71` |
| Unit 4 | Engine core: framebuffer, scene machine, content layer | `2febc83` |
| Unit 5 | **The demo runs.** Headless `cmd/ff6demo` renders real FF6 text | `5a8ffd8` |
| Unit 6 | Windowed host (Ebitengine, confined + build-tagged) and ADR-0001 | `a3c9797` |
| Unit 7 | ATB layer ported to Go (`internal/game/atb`) | `726f20e` |
| Unit 8 | Formation and monster decoders (`internal/game/battledata`) | `a4bf261` |
| Unit 9 | **The demo runs a battle screen.** Battle HUD scene from extracted data; D2–D5 recorded | `a5ed352` |
| — | Foundation merged to `main`, tagged `demo-0001-foundation-v0.1` | `297ba88` |
| Unit 10 | Records reconciled; [route content matrix](DEMO-0001-CONTENT-MATRIX.md) created; readiness recount corrected | — |

## Exact next action

**DEMO-0001A and the battle vertical's first pass are complete.** The shell
runs windowed and headless, every engine row is Integrated, and Unit 9 put a
battle screen on that shell driven by real formations, real monster records and
a live ATB layer.

The program is now on branch `demo/whelk-content-parity`, working the third
phase: **field/event content recovery**. Two things changed the shape of that
phase, both established during Unit 10 planning:

1. **Compression is the keystone, measured.** The
   [route content matrix](DEMO-0001-CONTENT-MATRIX.md) counts dependency
   pressure across the 19 beats: X1 blocks **8 of 19**, more than any other
   single dependency, and it still has zero records and zero code.
2. **The evidence for most of the missing families is already captured.** Every
   preserved Mesen savestate under `local_artifacts/experiments/` carries
   `ppu.vram` (64 KiB), `ppu.cgram`, `ppu.oamRam`, `spc.ram` (64 KiB ARAM) and
   the full PPU/DMA register set. `internal/mesenstate` parses those files
   already but exposes only WRAM and SRAM. Exposing the rest turns maps,
   sprites, backgrounds and the audio driver into work that needs **no live
   emulator session** — the decompressed output is already frozen and hashed.

Unit sequence from here:

- **Unit 11** — fix the palette defect that renders the demo's text black on
  black, and add the assertion that would have caught it.
- **Unit 12** — expose VRAM/CGRAM/OAM/ARAM and the PPU/DMA scalars from
  preserved savestates through `internal/mesenstate` and `ff6lab state`.
- **Unit 13** — **X1, the FF6 compression format**, recovered offline by
  differential reproduction of captured VRAM from a ROM span. Falsifier:
  byte-exact reproduction, or a recorded negative result.
- **Unit 14** — F1/F2, the first authentic field scene.
