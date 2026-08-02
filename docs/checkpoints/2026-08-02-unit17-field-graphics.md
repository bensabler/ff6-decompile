# Checkpoint 2026-08-02 — Unit 17, first authentic field graphics

(preceding: [Units 15-16, a negative result and the missing implementations](2026-08-02-units15-16-negative-result-and-tooling.md))

State: branch `demo/whelk-content-parity`, eleven commits ahead of `main` at
`297ba88`. Worktree clean. No emulator running, no background processes.

**Readiness F1 is Integrated (tileset only) — the first Field row ever to
move.** `ASSET-GFX-0002` is the Narshe field BG tileset: 256 tiles, 4bpp,
uncompressed, `ROMFILE:0x208460`, extracted from the operator's ROM and drawn
by `scenes.FieldTiles`. **All 256 tiles are proven byte-identical to captured
VRAM on two independent savestates**, as a skip-guarded differential test.

The scene draws the block in tile-number order and labels itself on screen,
because there is no tilemap and no header. An invented arrangement would look
far more like FF6 and is exactly the "route hardcoded as one bespoke cinematic"
the acceptance criteria prohibit. Deviation **D7** records it.

**The map-header search is stopped, cleanly.** Six meaningfully different
approaches across EXP-0051 and EXP-0052 have failed. The instrument is wrong:
this is a question for a disassembler or a trace, not a search.

**EXP-0052 did find real structure.** A map identifier — `$39` Narshe, `$17`
mines — at `WRAM:+$1305`, `+$13E2`, `+$1F80`, found by the discriminator
`02 == 04 != 05`. And a 33-byte record copied verbatim from ROM bank `$ED`
into `WRAM:+$0520`, with three captures on one lattice (`0x2D9173`, +20
records, +22 records) and `record[28]` carrying that same id. It is **not** the
header: milestones 02 and 04 share an id and a tileset but load different
records, so whatever indexes that table is not the identifier.

**The validation repair landed.** `internal/validate` now separates
`CompareIndexed` (composition — what `Sum256` hashes) from `CompareResolved`
(what a viewer sees, resolving each frame through its own palette and reporting
the visible fraction of drawn ink). `DecodeCGRAM` turns a captured CGRAM image
into a `Palette`, the only true FF6 palette available until F2. The tests state
why the distinction matters: two frames with different indices can look
identical, and a frame whose ink lands on undefined entries hashes fine and
renders blank — which is D0, and it passed every golden for seven units.

Two checks written earlier this session fired inside this unit and caught real
errors before commit: the Unit 11 palette assertion caught the field scene
drawing above `GrayPaletteDefined`, and `CheckReadinessSummary` caught the
summary drifting from its own tables.

Gates: gofmt clean; build/vet/test on both variants; `ff6lab audit` clean;
census clean; `archive verify` 9/9; restricted scan clean.

Exact next action: **read OAM from milestone 04** (CEN-GFX-0008). Every
preserved savestate carries `ppu.oamRam`, `ff6lab state oam` reaches it, and no
experiment has ever looked at it. Readiness F6 (field sprites) depends on it,
and it is the cheapest unstarted graphics unit in the project.

Needing one operator session, and it would close the header question in a
single observation: **run `probe dma-trace` over a map load.**
