# GFX-0001: Battle HUD fixed 8×8 text tile block

## Class
2bpp planar BG tile block (mode-1 BG3 text layer), 257 tiles.

## Runtime scenario/frame
`exp10-battle.mss` reloaded, captured atomically at anchor+120
(frame 210102, schedule-deterministic per EXP-0021, no input).
Capture probe: `mesen/probes/EXP-0023.lua` (headless bridge session).

## Local paths and hashes
`local_artifacts/experiments/EXP-0023/` — exp23-chr.hex
(`09e3787d…0fc5e0`), exp23-tilemap.hex (`7a0fb4a7…dfc2a4`),
exp23-cgram.hex (`3b09256b…8f2b0d60`), exp23-ppu-state.txt
(`9248025f…3d8d75`), rom_046FC0_8192.hex (`c90220fa…2235ef`),
decoded-grids.txt (`2caef05a…40d3c7`), hashes.sha256.

## VRAM ranges
BG3 chr word base `$5000` (byte `$A000`); dumped `$A000-$BFFF`
(512 tiles). The asset is the identical-to-ROM run: tiles
`$FF-$1FF` (byte `$AFF0-$BFFF`). Tiles `$00-$FE` are
runtime-composed (differ from ROM; interpretation: dynamic text
compose area — Tentative). Tilemap word base `$7800` (byte `$F000`),
32×32; every referenced tile index (37 distinct, all in `$180-$1FF`)
lies inside the identical run.

## CGRAM ranges
HUD text entries use BG-2bpp palette 0: BGR555
`$0000, $0000, $5294, $7FFF` (words 0–3). Other palettes recorded in
exp23-cgram.hex.

## OAM/layout
Not applicable (BG layer asset); OAM not captured in this unit.

## PPU/layer state
Mode 1, `mode1Bg3Priority=true`, brightness 15, layers 0/1 chr word
`$2000` (tilemaps `$6800`/`$7000`). `mainScreenLayers` read 0 at the
capture instant — mid-frame sampling artifact suspected (Tentative);
screen visibly renders (brightness 15, forcedBlank false).

## Source/loading provenance
**Block: `ROMFILE 0x047FB0-0x048FBF`** (4112 bytes, 257 tiles), addressed
by the affine relation **`ROMFILE = 0x046FC0 + 16×VRAMtileindex`**
(`ROMCPU:$C4:6FC0`), raw copy. Established by byte-identity: 15
distinctive glyph tiles each found exactly once in the 3 MiB ROM, all
implying anchor `0x046FC0`; the contiguous identical run covers tiles
`$FF-$1FF`. Ledger region: **ROM-0016**. The load path (which code/DMA
copies it into VRAM, and when) is **unresolved** — queued follow-up.

> **`0x046FC0` is an anchor, not a block start.** It back-projects to
> VRAM tile `$000`, which does not exist in the copy, and lands inside
> the attack-data region (`0x046AC0`), with which the font shares no
> bytes. The block begins at the first tile that does exist, VRAM `$FF`:
> `0x046FC0 + 16×0xFF = 0x047FB0`.
>
> This distinction is load-bearing. From 2026-07-30 to 2026-08-02 the
> `hud-font` extractor read 257 tiles starting at the anchor, so 255 of
> them were attack-table bytes rendered as tiles, while
> `manifests/rom-regions.json` ROM-0016 recorded the correct span the
> whole time. Nothing compared the two records. Fixed in extractor
> version 2.0.0; regression test
> `internal/extract.TestHUDFontMatchesROMLedger` now asserts the
> extractor's span against the ledger's.

## Compression/format
Stored uncompressed in ROM, SNES 2bpp planar (row = plane0,plane1
byte pair). Decoder: `internal/graphics/tile2bpp`.

## Reconstruction recipe
1. Read 4112 bytes at ROMFILE `0x047FB0` — equivalently, 16 bytes at
   `0x046FC0 + 16×index` for each index `$FF`…`$1FF`. **Start at index
   `$FF`, not 0**: lower indices are not font data and carry no runtime
   identity claim. Sheet index `i` corresponds to VRAM tile `$FF + i`.
2. Decode with `tile2bpp.Decode` (or `ff6lab tiles decode2bpp`).
3. Apply BG palette 0 via `internal/platform/bgr555` for display
   colors.

## Validation
Go-decoded VRAM tiles `$196/$19E/$1BF` == Go-decoded ROM tiles at the
same indices (grids archived in decoded-grids.txt); full-run byte
identity checked over all 512 tiles (257 identical, exactly tiles
`$FF-$1FF`). EXP-0023 record: falsifier (missing/inconsistent spans)
not met.

## Confidence
Raw-copy ROM provenance for tiles `$FF-$1FF` — **Confirmed**
(byte-exact, single-hit search consensus ×15 + contiguous-run check).
2bpp format — Confirmed (mode-1 BG3 + successful decode).
Glyph semantics (which tile is which character) — Unknown
(deliberately not assigned). Load path — Unknown.
