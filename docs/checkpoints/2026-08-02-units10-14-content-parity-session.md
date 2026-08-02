# Checkpoint 2026-08-02 — Units 10-14, the compression premise tested and refuted

## Current question

None open. The next unit is chosen, unblocked, and needs no emulator.

## State

Branch `demo/whelk-content-parity`, six commits ahead of `main` at `297ba88`.
Worktree clean, not pushed. No emulator running, no resident instrumentation,
no background processes.

The session's result is not the code it wrote. It is that **the program's
stated keystone was wrong, and the evidence to show that was already on disk.**

## Confirmed before this session

The merged DEMO-0001 foundation (units 0-9), re-verified: ROM SHA-256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2` unchanged,
all gates green on both build variants.

## Work completed

| Unit | Commit | Result |
|---|---|---|
| 10A | `b620635` | Records reconciled; route content matrix created |
| 10B | `c6f9844` | ADR-0002, five product lanes |
| 11 | `6d36e42` | **The demo's text was invisible.** Fixed |
| 12 | `94aa05f` | VRAM/CGRAM/OAM/ARAM and PPU/DMA state readable offline |
| 13 | `59b7e2b` | **EXP-0050**: compression is not the keystone |
| 14 | `54a5200` | The readiness summary is now asserted against its own tables |

### The finding that redirects the program

Readiness X1 said FF6's compression format gates F1, F6, B14, B15 and B16. The
route content matrix, built at the start of this session, counted the beats
those rows touch and ranked compression first at 8 of 19. **Neither number was
measured.** Both were inherited from the absence of records.

EXP-0050 tested it. `internal/romorigin` walks a captured memory image, probes
the ROM for each 32-byte window, and extends any match to its full run;
`ff6lab state origin` runs it over a preserved savestate. No emulator, no new
capture, one command.

Verbatim coverage of the full 64 KB VRAM image:

| Milestone | Scene | Verbatim |
|---|---|---|
| 04-free-movement | Narshe exterior | **51 %** |
| 05-mines-entry | Mines interior | **52 %** |
| 02-narshe-entry | Narshe entry | **47 %** |
| 03-first-scripted-battle | Battle | **38 %** |
| 01-opening-cinematic | Opening, BG mode 7 | **0 %** |

The field scenes resolve into three contiguous **uncompressed** spans carrying
20 KB of BG tile data — `ROMFILE:0x208460` (8192 B), `0x223000` (8192 B),
`0x224F00` (4096 B) — plus ~18 short runs of 128-171 bytes from bank `$E6`,
consistent with animated tiles. The first two blocks are **shared between the
Narshe exterior and the mines interior**. No decompressor places any of it.

Compression was **withdrawn** from the matrix's pressure table rather than
re-ranked. 48-62 % of each image is unmatched, but a verbatim search is
defeated by bit-plane reordering, by runtime composition, by tilemaps built in
WRAM and by a single changed byte. Ranking it on the unmatched fraction would
repeat the original error in the other direction. What X1 gates is Unknown, and
is now recorded as Unknown.

**F1 gains what it never had: anchors.** Its next action is to follow
`0x208460` and `0x223000` to the header that selects them. Map headers now lead
the pressure table at 6 beats.

By-products: the HUD font's destination is `VRAM:$AFF0`, half of the load path
T1 lists as missing. Battle sprite regions `VRAM:$6180-$7829` map to short
verbatim runs in banks `$D8`/`$E9`. The battle BG chr region `$4000-$47FF` is
**not** verbatim while `$4800-$4FFF` is, which makes battle backgrounds the
strongest remaining compression candidate on the route.

### The enabling unit

Every preserved savestate carries `ppu.vram`, `ppu.cgram`, `ppu.oamRam`,
`spc.ram` (64 KB ARAM), the full PPU register set and per-channel DMA state.
`internal/mesenstate` had parsed those files since it was written and exposed
only WRAM and SRAM. Unit 12 exposed the rest. EXP-0050 exists because Unit 12
did.

A caveat is recorded in the code, the CLI output and the census because it is
easy to get wrong: **DMA channel state is configuration at an instant, not a
transfer log.** Live DMA tracing does not exist in this project — no Lua probe
and no Go path reads a DMA register during execution, despite `/trace-dma`, the
`dma-tracer` skill, the `dma-researcher` agent and `TRACE_DMA.md` all existing.
Measured falsifier: channel 0 in the mines state is set up to write `$2118`
from `$7EC180`, but VRAM matches WRAM at `$0C180` on only 18 % of bytes with a
49-byte leading run. That setup is not the transfer that produced the VRAM.

### The defect that was on screen the whole time

`framebuf.BlitTile` computes `dst = ink + PaletteBase`, so `PaletteBase` is a
sub-palette base. Both scenes passed `white = 3` / `gray = 2`, reading it as a
brightness level. HUD font ink values are 1-3, so bright strings landed on
entries 4-6 while `GrayPalette` defined only 0-3.

Measured on the real font and archive: **32 % of drawn ink resolved to a
visible colour.** Two thirds of what the demo drew was invisible, since Unit 4.

The goldens could not see it. `Sum256` hashes indices and is deliberately
palette-independent — the right design for a composition golden, and exactly
why a frame can hash correctly and be blank on screen. The doc comment on
`BlitOptions.PaletteBase` had said "4- or 16-color sub-palette" the whole time;
the callers and the contract disagreed and nothing compared them.

The units were the defect, so the fix is a type. `content.SubPalette` is a slot
number. Expressed as slots the old values address entries 8-15, outside
anything the palette defines, so the same mistake now fails a test — verified
by re-applying it. Visible ink **32 % → 63 %**, with the ink mask byte-identical
before and after (6406 pixels), so the change is index remapping only.

### Three records that disagreed with themselves

All the same failure mode as retired deviation D1.

1. **The readiness summary vs its own tables.** Unit 0 claimed 53 requirements,
   33 Unknown, 7 Evidence Ready; the file carried 55, 36 and 6, and has since
   `969b5dd`. Recounted, and **the comparison is now a check** —
   `audit.CheckReadinessSummary`, which caught a status token this very session
   had borrowed from the content matrix's vocabulary within a minute of being
   written.
2. **`MESEN_CAPABILITY_MATRIX.md` understated the project's own instrument.**
   VRAM, CGRAM, OAM, ARAM and DSP access were all recorded `Unknown`, although
   `bridge.lua:238` has had them wired since it was written and EXP-0023/0024
   had used four of them. Graphics and audio work was being deferred as "needs
   new instrumentation" when the instrument existed.
3. **`STATISTICS.md`** was stale on every count.

## Last raw observation

`ff6lab state origin local_artifacts/scenarios/SCN-0001/04-free-movement/run1-04.mss vram`
— three contiguous spans at `0x208460`, `0x223000`, `0x224F00`, 51 % coverage
in 33 spans.

## Active emulator state

None.

## Breakpoints/watchers

None.

## Evidence paths

No new evidence was captured. EXP-0050 read the existing frozen savestates
under `local_artifacts/scenarios/SCN-0001/`, unmodified, each covered by its
directory's `hashes.sha256`. Reproducible with one command and the operator's
ROM.

## Files changed

Six commits. New: `internal/romorigin`, `internal/mesenstate/hardware.go`,
`internal/audit/readiness.go`, `docs/demo/DEMO-0001-CONTENT-MATRIX.md`,
`docs/decisions/ADR-0002-project-product-lanes.md`,
`docs/experiments/EXP-0050-vram-rom-provenance-sweep.md`,
`reconstruction/snes/README.md`, `specs/README.md`. Modified: the palette path
(`framebuf`, `content`, both scenes), `ff6lab state`, and the DEMO-0001 records,
dashboards, manifests and indexes.

## Tests and quality gates

`gofmt -l .` clean. `go build`, `go vet`, `go test` pass on both build variants.
`ff6lab audit` clean (now eleven checks). `census validate` clean.
`archive verify` 8/8. Restricted-file scan clean. `romorigin`'s fuzz target ran
5.8 M executions with no failures.

## Git status

Branch `demo/whelk-content-parity`, six commits ahead of `main`, worktree
clean, **not pushed**.

## Unresolved decisions

None.

## Blockers

None for the next unit. Route dependency pressure, re-derived: map headers
(F1) 6 beats, field sprites (F6) 5, music sequences (A3) 5, event dispatch (F8)
4, dialogue corpus (F10) 3, Whelk ids (B13) 2. Compression withdrawn pending
evidence.

## Exact next action

**Unit 15 — find the map header record that selects `ROMFILE:0x208460` and
`0x223000`.**

F1 has anchors for the first time. Two field maps share those two blocks and
differ in a third, so the selecting record is a small table of block ids or
pointers, and two maps that agree on two of three entries is a strong
discriminator. Search the ROM for pointer or id patterns that would produce the
observed triples, and cross-check against the mines interior's third block.

Falsifier: no candidate record reproduces both maps' block sets.

Cheaper unit available if that stalls: **read OAM from milestone 04**
(CEN-GFX-0008). Every savestate carries `ppu.oamRam`, `ff6lab state oam` now
reaches it, and no experiment has ever looked at it. F6 depends on it.

## Recommended next command

`/investigate-variable` or direct static work on the ROM; then
`/run-quality-gates` and `/checkpoint`.
