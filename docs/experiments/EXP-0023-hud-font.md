# EXP-0023: Battle HUD font — runtime capture, ROM provenance, decoder proof

- **Status:** completed (2026-07-30)
- **Question (Unit 4, graphics vertical):** where do the battle HUD
  text tiles live at runtime (VRAM chr base, layer, bit depth), and
  does a contiguous ROM region hold those tile bytes **verbatim**
  (raw-copy provenance)? Deliver the full vertical: runtime state →
  ROM source → Go decoder with golden tests → byte/pixel comparison.
- **Starting state:** headless bridge session; `exp10-battle.mss`
  reloaded and advanced a fixed +120 frames with no input (EXP-0021
  established schedule-determinism, so the capture frame is
  reproducible). Bridge `read` extended with `vram|cgram|oam|aram|rom`
  memory types (tracked change; restart required).
- **Method:**
  1. Capture the battle PPU state (`state` dump: bgMode, per-layer
     chrAddress/tilemapAddress).
  2. Dump CGRAM (512 B), the candidate text-layer tilemap (2 KB), and
     the text-layer chr region (4 KB+) via the extended `read`.
  3. Identify HUD glyph tiles from the tilemap's populated rows
     (digit/letter tile indices referenced by the HUD).
  4. Search the **local ROM file** (SHA-256
     `0f51b4fc…d8d5e2`, never committed) for several distinct
     captured tile byte spans; require consistent relative offsets for
     a base identification.
  5. Implement/extend the Go tile decoder for the observed bit depth
     with golden vectors taken from the capture; compare decoded ROM
     bytes against decoded VRAM bytes pixel-for-pixel.
- **Expected outcomes:**
  - *Raw-copy provenance:* all searched spans found at one consistent
    ROM base → font ROM address Confirmed (byte identity), decoder
    proof closes the vertical.
  - *Transformed provenance:* spans absent from ROM → the font is
    stored compressed/transformed; record the negative result, keep
    the decoder + VRAM golden tests, and queue the transform hunt.
- **Falsifying outcome (for raw-copy):** any searched span missing, or
  inconsistent relative offsets across spans.
- **Limitations:** provenance by byte identity, not by a traced DMA
  load — the load-path trace (who copies font→VRAM, when) is a
  follow-up unit. Bit depth and layer are read from captured PPU
  state, not assumed.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0023/`
  (exp23-ppu-state.txt, exp23-cgram.hex, exp23-tilemap.hex,
  exp23-chr.hex, rom_046FC0_8192.hex, decoded-grids.txt,
  hashes.sha256 — SHA-256s in the asset record).
- **Result:** **raw-copy provenance Confirmed; the vertical closes.**
  Full detail in [GFX-0001](../graphics/GFX-0001-battle-hud-font.md).
  - Capture landed atomically at frame 210102 (anchor+120 exactly).
  - Battle PPU: mode 1, BG3 (2bpp text layer) chr word `$5000`,
    tilemap word `$7800`; HUD text uses BG palette 0
    (`$0000/$0000/$5294/$7FFF`).
  - Tilemap references 37 distinct tiles, all in `$180-$1FF`.
  - ROM search: 15 distinctive glyph tiles → each exactly one hit in
    the 3 MiB ROM, unanimous base **ROMFILE `0x046FC0`** = chr byte
    base. Per-tile comparison over all 512 tiles: one contiguous
    identical run, tiles `$FF-$1FF` (257 tiles), containing every
    HUD-referenced tile. Tiles `$00-$FE` differ (runtime-composed;
    Tentative: dynamic text compose area).
  - Go proof: new `internal/graphics/tile2bpp` (tests + fuzz) and
    `ff6lab tiles decode2bpp`; decoded VRAM grids == decoded ROM
    grids for sampled tiles `$196/$19E/$1BF` (archived).
- **Confidence:** ROM provenance (tiles `$FF-$1FF` ↔
  `ROMFILE 0x047FB0-0x048FBF`) — **Confirmed**. 2bpp format —
  Confirmed. `mainScreenLayers=0` at capture — recorded oddity,
  mid-frame sampling suspected (Tentative). Load path (who copies
  ROM→VRAM, when) — Unknown. Glyph semantics — Unknown (not
  assigned).
- **Next action:** follow-ups queued — trace the ROM→VRAM load path
  (DMA watch at battle entry); name glyphs empirically (render known
  WRAM values against tilemap digits); investigate the `$00-$FE`
  compose region.
