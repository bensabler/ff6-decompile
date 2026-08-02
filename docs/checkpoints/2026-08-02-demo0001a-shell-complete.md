# Checkpoint 2026-08-02 — DEMO-0001A complete, the demo runs

## Current question

None open. Eight bounded units closed cleanly. The next unit is chosen and
unblocked: decode formations and monster records from the archived tables.

## State

The project has shifted into **playable-vertical-slice production**.
`DEMO-0001 — New Game to Whelk Victory` is the active program; SCN-0001
remains the evidence program that DEMO-0001 consumes.

**The demo runs**, windowed and headless, rendering real extracted FF6 data.
Branch `demo/new-game-to-whelk`, 8 commits ahead of `main`. Worktree clean.

## Confirmed before this session

The ATB program (EXP-0041…0047), the battle damage pipeline (DISC-0001…0007),
event flags (DISC-0008), formations `ROMCPU:$CF6200`, monster records
`ROMCPU:$CF0000`, the HUD font block (GFX-0001), and CORR-0001's pointer
advance at `ROMCPU:$C09B5C`.

## Work completed

| Unit | Commit | Result |
|---|---|---|
| 0 | `969b5dd` | DEMO-0001 records: milestone, readiness matrix (53 requirements), acceptance scorecard, deviations register |
| 1 | `cdc55e8` | **Defect fixed.** HUD font extractor read `ROMFILE:0x046FC0`, the arithmetic anchor, not the block start `0x047FB0`. 255 of 257 archived tiles were attack-table bytes |
| 2 | `887f676` | `internal/platform/snesaddr`; 21 existing offsets asserted through it |
| 3 | `aabad71` | **EXP-0049** — the text encoding indexes the font block directly |
| 4 | `2febc83` | Engine core: framebuffer, scene machine, content layer, enforced import boundaries |
| 5 | `5a8ffd8` | **The demo runs.** Headless `cmd/ff6demo` |
| 6 | `a3c9797` | Windowed host (Ebitengine, confined + build-tagged), ADR-0001 |
| 7 | `726f20e` | ATB layer ported to Go |

### Two findings worth carrying forward

**1. The extractor and the ROM ledger disagreed, and nothing compared them.**
`manifests/rom-regions.json` ROM-0016 recorded the correct font span the whole
time, and even noted the anchor had "no data overlap". The extractor, the asset
manifest and the archive were mutually consistent and all wrong about the ROM.
Two independent records of one fact are worth nothing unless something asserts
them against each other — now `TestHUDFontMatchesROMLedger` does, and the
archive-vs-ROM differential covers the general case.

**2. EXP-0049 resolved a mapping the census had carried as Unknown**, from
tracked evidence, with no emulator run. `VRAM tile = $100 + encoded byte`.
The discriminating trial was decoding the whole EXP-0023 HUD tilemap through
the relation: it yields "Were-Rat", "Repo Man", "WEDGE", "VICKS", "MagiTek",
"Item", "Row", "Tek" and HP digits across all 37 referenced tiles. Secondary
results: `$BF` = `?`; five gauge tiles identified (CEN-GFX-0005); the HUD text
layout registered (CEN-BATTLE-0014); 47 tiles recorded as **deliberately
unidentified**.

## Last raw observation

`local_artifacts/demo-frames/frame-000060.png` — the boot scene rendering the
corrected font. Verified by eye and by per-row ink measurement.

## Active emulator state

**None.** No emulator was run this session. Every finding is static analysis of
evidence EXP-0023 preserved, or implementation.

## Breakpoints/watchers

None. No resident instrumentation.

## Evidence paths

- `local_artifacts/experiments/EXP-0023/rom_046FC0_8192.hex` and
  `exp23-tilemap.hex` — EXP-0049's inputs, frozen in that directory's
  `hashes.sha256`
- `local_artifacts/archive/` — regenerated; `archive verify` 8/8 clean
- `local_artifacts/demo-frames/` — demo output (gitignored)

## Files changed

New packages: `internal/platform/snesaddr`, `internal/platform/snespad`,
`internal/graphics/framebuf`, `internal/engine`, `internal/engine/headless`,
`internal/engine/ebitenhost`, `internal/content`, `internal/assetmanifest`,
`internal/game/hudfont`, `internal/game/atb`, `internal/game/scenes`,
`cmd/ff6demo`.

New records: `docs/demo/` (4), `docs/decisions/ADR-0001-rendering-host.md`,
`docs/experiments/EXP-0049-hud-font-glyph-mapping.md`,
`data/graphics/hud-font-glyphs.json`.

Modified: `internal/extract` (font fix + manifest split), `internal/game/textenc`
(`Encode`, `$BF`), `internal/game/attackdata` (`TableCPUAddr`), `internal/audit`
(import boundaries, rendered-asset guard), `ARCHITECTURE.md`, `README.md`,
`CLAUDE.md`, `Makefile`, `.github/workflows/ci.yml`, dashboards, census.

## Tests and quality gates

All green at the last commit:

```text
gofmt -l .            clean
go build ./...        ok        go vet ./...      ok
go test ./...         25 packages, no failures
go build -tags gui ./cmd/ff6demo   ok
go vet -tags gui ./...             ok
go test -tags gui ./...            no failures
go mod verify         all modules verified
ff6lab audit          clean       ff6lab census validate   clean
ff6lab archive verify 8/8 ok
restricted-file scan  clean (now covers rendered images too)
```

Fuzz targets clean at 5s each: `snesaddr` (×2), `framebuf` (×2), `textenc`,
`atb`.

## Git status

Branch `demo/new-game-to-whelk`, worktree clean, 8 commits ahead of `main`.
**Not pushed** — no push was requested.

## Unresolved decisions

- Whether to merge `demo/new-game-to-whelk` into `main` or keep it long-lived.
- Whether the boot scene stays once a real title screen becomes implementable.

## Blockers

Nothing blocks the next unit. The standing blockers are unchanged and all
belong to the field half: the map system, the dialogue corpus, the event
opcode table, sprite formats, audio sequences, and FF6's compression format
are all at zero records. `dashboards/BLOCKERS.md` is accurate.

One correction to a plan assumption, recorded because it changed the design:
**Ebitengine v2.9.9 requires cgo on both macOS and Linux** (`CGO_ENABLED=0`
fails on each; only Windows and js/wasm build without it). The `gui` build tag
is therefore the design, not a fallback.

## Exact next action

**Unit 8 — decode formations and monster records into Go**, feeding a battle
scene from real data.

Both tables are already extracted to the archive and both formats are
Confirmed:

- `local_artifacts/archive/raw/formations.bin` — `ROMCPU:$CF6200` + 15×id,
  verified against seven independent encounters (EXP-0030). Unknown: bytes
  `+$00/01` and `+$08..+$0E`, and the **monster-id high-bit extension** that
  Whelk (formation 432) needs.
- `local_artifacts/archive/raw/monsters.bin` — `ROMCPU:$CF0000`, 32-byte
  stride; `+$08` HP, `+$0A` MP, `+$01` battle power (EXP-0028/0029).

New package `internal/content/battledata`, reading through the archive (never
the ROM — the import boundary forbids it). Then Unit 9: a battle HUD scene
consuming the font, the ATB layer, and these tables.

## Recommended next command

```text
/resume-session
```

Then proceed directly to Unit 8. `docs/demo/DEMO-0001-READINESS.md` is the
critical-path instrument; rows B7 and B8 are the target.
