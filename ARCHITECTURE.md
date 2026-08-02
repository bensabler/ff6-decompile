# Architecture

## Repository layers

```text
cmd/ff6lab
    CLI wiring and subcommands

internal/rom
    ROM image loading and revision verification

internal/extract
    deterministic extractors and the asset manifest

internal/census
    content census, ROM-ownership ledger, coverage generation

internal/audit
    repository integrity checks behind `ff6lab audit`

internal/mesenstate
    Mesen savestate reading, no emulator required

internal/scenario
    reproducible emulator route models

internal/platform
    SNES address spaces, ROM mapping, endian/bit utilities
    - platform/snesaddr: ROMCPU <-> ROMFILE (HiROM), per-window confidence
    - platform/bgr555:   SNES 15-bit color

internal/research
    evidence, experiment, discovery, provenance, catalog, validation models

internal/graphics
    planar tiles, palettes, tilemaps, OAM, sprites, backgrounds, animation

internal/audio
    BRR, SPC700, DSP, sequences, instruments, rendering and comparison

internal/game
    verified FF6-domain structures and behavior

internal/validate
    hashes, image comparisons, event comparisons, audio comparisons

internal/project
    configuration, paths, workspace, indexes, migrations
```

## Dependency direction

```text
cmd
 ↓
project / research / game / graphics / audio / validate
 ↓
platform
```

`platform` must not import game-specific packages. `research` stores evidence metadata but does not contain emulator automation. `game` represents verified behavior and may use narrow platform types.

Because `platform` may not import `rom`, `snesaddr.ImageSize` duplicates `rom.Size` deliberately; a test asserts the two agree.

## Address arithmetic

Every ROM offset used to be a hand-computed constant with its mapping in a comment. Two independent records of one fact can silently diverge — and did, for the battle HUD font block (DEMO-0001 deviation D1, fixed 2026-08-02).

`internal/platform/snesaddr` holds the arithmetic. Existing offsets stay constants rather than becoming derived values: deriving them would trade a compile-time constant for a runtime variable and an init-time panic path. Instead each constant is paired with the CPU address its evidence record cites, and a table-driven test asserts the pair through `snesaddr`. New offsets should follow that pattern.

The mapping windows do not carry equal evidence. Banks `$C0-$FF` are Confirmed (18/18 Mesen ROM captures, CORR-0001); the lower and mirror windows are Strong hypothesis. `snesaddr.Window.Confirmed()` lets a caller require the evidenced window.

## Command strategy

The primary executable is `ff6lab`.

Planned command groups:

```text
ff6lab project
ff6lab rom
ff6lab evidence
ff6lab experiment
ff6lab asset
ff6lab graphics
ff6lab audio
ff6lab validate
ff6lab report
```

Command packages contain parsing and dependency wiring. Domain logic belongs in `internal/`.

## Data boundary

The repository contains metadata, source, schemas, synthetic fixtures, and hashes. Local commercial data stays in `local_artifacts/`.

## Stability policy

All implementation packages remain internal until a genuine externally useful and stable API emerges. This follows Go's recommendation to use `internal` for supporting packages that should remain refactorable.
