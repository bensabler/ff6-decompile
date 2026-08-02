# DEMO-0001: New Game to Whelk Victory — playable Go demo

- **Status:** ACTIVE (program started 2026-08-02)
- **Kind:** Production milestone (deliverable is a runnable program, not a record)
- **Scenario record:** [SCN-0001](../scenarios/SCN-0001-opening-to-whelk.md)
- **Machine manifest:** `data/scenarios/opening-to-whelk.json` (19 beats B01–B19)
- **Readiness matrix:** [DEMO-0001-READINESS.md](DEMO-0001-READINESS.md)
- **Acceptance criteria:** [DEMO-0001-ACCEPTANCE.md](DEMO-0001-ACCEPTANCE.md)
- **Deviations register:** [DEMO-0001-DEVIATIONS.md](DEMO-0001-DEVIATIONS.md)
- **ROM revision:** SHA-256 `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
- **Branch:** `demo/new-game-to-whelk`

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
go run ./cmd/ff6demo
```

Headless frame capture, for automated comparison and CI:

```bash
go run ./cmd/ff6demo -headless -frames 120
```

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
| Unit 5 | **The demo runs.** Headless `cmd/ff6demo` renders real FF6 text | — |

## Exact next action

Unit 6 — the windowed host. Ebitengine v2 behind a `//go:build gui` tag,
confined to `internal/engine/ebitenhost`, plus ADR-0001.

The build tag is not a fallback, it is the design. Measured 2026-08-02:
ebiten v2.9.9 requires cgo on **both** Linux and macOS (`CGO_ENABLED=0` fails
on each; only Windows and js/wasm build without it). Keeping the host behind a
tag is what stops `go build ./...`, `go vet ./...` and `go test ./...` from
acquiring a system-library dependency, so CI keeps running the full demo test
surface on a bare container.

After that the loop turns to the battle vertical: port the Confirmed ATB layer
(EXP-0041…0047) into Go, then formation and monster decoders, then a battle
HUD scene that can use the font this milestone integrated.
