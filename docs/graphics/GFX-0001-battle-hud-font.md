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
**ROMFILE `0x046FC0` + 16×tileindex (`ROMCPU:$C4:6FC0`), raw copy.**
Established by byte-identity: 15 distinctive glyph tiles each found
exactly once in the 3 MiB ROM, all implying base `0x046FC0`; the
contiguous identical run covers tiles `$FF-$1FF` (4112 bytes,
`ROMFILE 0x047FB0-0x048FBF`). The load path (which code/DMA copies it
into VRAM, and when) is **unresolved** — queued follow-up.

## Compression/format
Stored uncompressed in ROM, SNES 2bpp planar (row = plane0,plane1
byte pair). Decoder: `internal/graphics/tile2bpp`.

## Reconstruction recipe
1. Read 16×N bytes at ROMFILE `0x046FC0 + 16×index` (index `$FF-$1FF`
   verified; lower indices carry no runtime identity claim).
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
