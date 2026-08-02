# FF6 Reconstruction Session 005 — DEMO-0001 content parity, units 10–18

- **Date:** 2026-08-02
- **Investigator:** Benjamin Sabler + Claude (autonomous, operator-directed)
- **ROM identity:** SHA-256
  `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
  re-verified at session start, unchanged.
- **Emulator:** **none.** No Mesen instance was opened at any point. Every
  runtime observation came from savestates frozen in earlier sessions.
- **Starting checkpoint:**
  [2026-08-02 — DEMO-0001A complete](../checkpoints/2026-08-02-demo0001a-shell-complete.md)
- **Ending checkpoint:**
  [2026-08-02 — Tactical pause](../checkpoints/2026-08-02-tactical-pause-evidence-harvested.md)
- **Branch:** `demo/whelk-content-parity`, cut from `main` at `297ba88`
- **HEAD:** `e6c0d1d`, 16 commits, pushed

## Objective

Advance DEMO-0001 from a demo with "a spine, not yet a game" toward one that
looks and behaves like the opening of FF6, by recovering and integrating
content dependencies along the SCN-0001 route.

## The session's central result

**The program's stated keystone was wrong, and the evidence to disprove it was
already on disk.**

Readiness X1 had asserted for ten units that FF6's compression format gated
maps, field sprites, battle backgrounds, enemy graphics and party sprites.
[EXP-0050](../experiments/EXP-0050-vram-rom-provenance-sweep.md) tested that
premise for the first time and refuted it: **47–52 % of a field scene's VRAM is
present in the ROM verbatim**, including all three contiguous BG tile blocks.
No decompressor is involved in placing them.

That finding was only reachable because of a second one: **every preserved
savestate already carried the runtime evidence.** `internal/mesenstate` had
parsed those files since it was written while exposing only WRAM and SRAM.
Exposing `ppu.vram`, `ppu.cgram`, `ppu.oamRam`, `spc.ram` and the PPU/DMA
registers converted maps, sprites, backgrounds and the audio driver from work
that needed an emulator session into work that needed analysis.

## Units completed

| Unit | Commit | Result |
|---|---|---|
| 10A | `b620635` | Records reconciled; route content matrix created |
| 10B | `c6f9844` | [ADR-0002](../decisions/ADR-0002-project-product-lanes.md), five product lanes |
| 11 | `6d36e42` | Palette defect fixed — the demo's text had been invisible |
| 12 | `94aa05f` | VRAM/CGRAM/OAM/ARAM and PPU/DMA state exposed offline |
| 13 | `59b7e2b` | **EXP-0050**; readiness X1 re-scoped |
| 14 | `54a5200` | `audit.CheckReadinessSummary` |
| 15 | `2673065` | **EXP-0051** — negative; EXP-0050 corrected |
| 16 | `7969d50` | DMA implementations; five missing commands created |
| 17 | `312f836` | **First authentic field graphics rendered**; header search stopped |
| 18 | `a4e899d` | **EXP-0053** — OAM read for the first time |
| — | `38a62a1`, `e6c0d1d` | Evidence harvest and tactical pause |

Checkpoints: `7ae1403`, `c861bc1`, `f322caf`, `8591aef`, `e6c0d1d`.

## Experiments

| ID | Question | Result |
|---|---|---|
| [EXP-0050](../experiments/EXP-0050-vram-rom-provenance-sweep.md) | Which captured VRAM regions are verbatim ROM copies? | **Confirmed** field tile blocks uncompressed; X1's blanket claim **refuted** |
| [EXP-0051](../experiments/EXP-0051-map-tileset-selector-search.md) | Are the block addresses stored as pointers? | **Negative** — four encodings, none found |
| [EXP-0052](../experiments/EXP-0052-map-id-and-33-byte-record.md) | What in WRAM distinguishes one field map from another? | **Partial** — map id and a 33-byte record table located; still not the header |
| [EXP-0053](../experiments/EXP-0053-oam-field-sprite-composition.md) | What does OAM hold on the Narshe exterior? | **Confirmed** composition; player sprite located at `ROMFILE:0x150180` |

No discovery records were promoted. Every result is either a located address
with an experiment behind it, or a negative — neither clears the discovery bar.

## Evidence

All four experiments have preserved directories under
`local_artifacts/experiments/`, each with a `README.md` naming its source
savestate and regeneration command, and a verified `hashes.sha256`. All
gitignored; every artifact is ROM-derived.

`EXP-0051/search-pointer-encodings.py` and `EXP-0052/wram-discriminator.py`
preserve analyses that were ad-hoc when run. **Both were replayed at session
end and reproduce their recorded numbers exactly.**

No new capture was taken. Sources were the frozen SCN-0001 milestone
savestates and `local_artifacts/experiments/EXP-0035/`.

## Code changes

78 files changed, +7 049 / −301 against `main`; 46 new tracked files.

New packages: `internal/romorigin` (VRAM→ROM provenance),
`internal/platform/snesdma` (DMA registers and trace parsing),
`internal/platform/snesoam` (object attribute table),
`internal/audit/readiness.go`, `internal/validate/framebuffer.go`.

Extended: `internal/mesenstate` (PPU/DMA state and four memory images),
`cmd/ff6lab state` (`vram|cgram|oam|aram`, `ppu`, `sprites`, `origin`),
`internal/extract` (`ASSET-GFX-0002`), `internal/content` (tileset loader,
`SubPalette`), `internal/game/scenes` (`FieldTiles`).

Instrumentation: `mesen/probes/dma-trace.lua` — written, **unexercised**.

Doctrine: five commands created (`recover-compression`, `recover-map`,
`recover-text`, `recover-event-opcode`, `capture-frame`) with playbooks; three
`reconstruct-*` aliased.

## Defects found and fixed

| Defect | Impact | Fix |
|---|---|---|
| **D0** — `PaletteBase` read as a brightness level | **32 % of drawn ink visible**; the demo's bright text was black on black since Unit 4 | `content.SubPalette` carries the units; visible ink 32 % → **63 %** |
| Readiness summary vs its own tables | "53 requirements" where the file carried 55, propagated into four checkpoints and three dashboards | Recounted; `CheckReadinessSummary` added |
| `MESEN_CAPABILITY_MATRIX.md` | VRAM/CGRAM/OAM/ARAM marked `Unknown` while wired and used — work deferred as "needs instrumentation" that existed | Rows corrected, three states distinguished |
| `mesenstate` DMA size | A 64 KB upload printed as `0` | `ByteCount`; zero means 65536 |
| EXP-0050's "shared with the mines" | An unmeasured claim in a canonical record | Corrected inline by EXP-0051 |
| EXP-0053 first pass | A 128-byte probe spanned a VRAM gap and reported the player's sprite absent from the ROM | Re-probed per tile; correction recorded |

**Three checks written this session later caught errors before commit**: the
palette assertion (Unit 17), and `CheckReadinessSummary` (Units 17 and 18).

## Readiness movement

| Status | Session start | Session end |
|---|---|---|
| Validated | 1 | 1 |
| Integrated | 13 | **15** |
| Researching | 2 | **3** |
| Unknown | 25 | 27 |
| Total rows | 55 | **57** |

**F1 → `Integrated` (tileset only)** and **F6 → `Researching`** are the first
and second Field rows ever to move. B19 and X4 were added after the content
matrix found requirements the route needs and no row owned.

## Tests and quality gates

Green at every commit: `gofmt -l .` clean; `go build`, `go vet`, `go test` on
both build variants; `ff6lab audit` (eleven checks, one new); `census validate`;
`archive verify` 9/9; restricted-file scan clean.

37 packages, 30 with tests, **18 fuzz targets** (four added). The
`ASSET-GFX-0002` differential asserts all 256 tiles against captured VRAM on
two independent savestates.

## What was not done

- No Mesen session. `probe dma-trace` is unexercised.
- No Ghidra correlation.
- The map header is unlocated after six approaches.
- No map is rendered — the field scene shows tiles in tile-number order,
  deviation **D7**, because there is no tilemap and inventing one would violate
  the acceptance criteria.
- Audio untouched; A2–A6 unchanged.

## Recurring failure mode

Recorded because it appeared **five times**: *two records of one fact are worth
nothing unless something compares them.* The `hud-font` extractor vs ROM-0016
(D1); the readiness summary vs its tables; `BlitOptions.PaletteBase`'s doc
comment vs its callers (D0); EXP-0050's unmeasured shared-blocks claim; and
EXP-0053's wrong probe geometry — where **a negative from the wrong geometry
looks exactly like a real one.**

## Next action

**Project-wide workflow and orchestration audit**, as directed. No research or
implementation unit before it.

When content work resumes: test the block-boundary alternative (free, would
invalidate all four pointer encodings); diff the sprite tile region between
milestones 02 and 04 (F6, no new capture); then run `probe dma-trace` over a map
load, which needs one operator session.
