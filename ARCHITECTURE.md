# Architecture

## Repository layers

```text
cmd/ff6lab
    CLI wiring and subcommands

internal/platform
    SNES address spaces, ROM mapping, endian/bit utilities

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
