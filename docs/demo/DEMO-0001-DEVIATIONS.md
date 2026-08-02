# DEMO-0001 deviations register

- **Milestone:** [DEMO-0001](DEMO-0001-new-game-to-whelk.md)
- **Updated:** 2026-08-02 (Unit 11 — palette correctness)

Every way the Go demo differs from the original, and every piece of provisional
scaffolding, gets a row here. The register exists so that differences are
**visible**, not so they are permitted.

Rules:

- A deviation is recorded when it is introduced, not when it is noticed later.
- Every row names an **exact replacement requirement** — the specific evidence
  or implementation that retires it.
- A deviation is never retired by weakening a test.
- Scaffolding without a row here is a defect.

## Severity

| Level | Meaning |
|---|---|
| `Cosmetic` | Visible difference with no effect on behavior or progression |
| `Behavioral` | The demo does something the original does not, or vice versa |
| `Structural` | An unverified model stands in for a system whose real behavior is unknown |
| `Parity-blocking` | Prevents an acceptance step from ever passing honestly |

## Open

| ID | Severity | Deviation | Introduced | Exact replacement requirement | Status |
|---|---|---|---|---|---|
| D2 | `Cosmetic` | The battle scene labels enemies by **monster record id** ("MONSTER 19") instead of by name. The monster name table is unlocated (readiness B9), so no name is extractable | 2026-08-02, Unit 9 | Locate the monster name table, add an extractor, and render the extracted names. Until then the label must stay obviously non-diegetic — an invented FF6-sounding name would be indistinguishable from recovered data, which is the exact failure this register exists to prevent | **Open** |
| D3 | `Structural` | Enemy ATB increments use the constant `$9C` that EXP-0043 **measured** for formation 14's enemies, rather than a computed value. `ROMCPU:$C209F0` is only partly decoded — it reads per-slot Speed and adds `$14` before an enemy-only multiply by `+$3A90` | 2026-08-02, Unit 9 | Decode `ROMCPU:$C209F0` fully and compute the increment from the monster record's Speed field and the sampled `+$3A90`. `internal/game/atb` already takes increments as an input for this reason, so the replacement is local to the scene | **Open** |
| D4 | `Structural` | The battle scene renders **no party side**. Character initialisation is unlocated (readiness F14), so there is no table to read party names, HP, or MP from; EXP-0040's party HP values are observations of one playthrough, not a data source | 2026-08-02, Unit 9 | Locate the new-game character initialisation and the party record layout, then render slots 0-3 from extracted data. The scene states "PARTY SIDE NOT IMPLEMENTED / NO CHARACTER DATA SOURCE" on screen rather than drawing invented rows | **Open** |
| D5 | `Cosmetic` | Gauge fill quantisation is a **reading of two samples**, not a decoded rule. EXP-0049 saw a full bar and a nearly empty one, from which "filled cells use `$F8`, the boundary cell uses `$F1`, the rest `$F0`" is inferred | 2026-08-02, Unit 9 | Capture intermediate gauge states and decode how the engine picks a segment tile. CEN-GFX-0005 carries the open question | **Open** |
| D6 | `Cosmetic` | Every colour on screen comes from `framebuf.GrayPalette`, a **project-authored** two-sub-palette ramp, not from an extracted FF6 palette. No extractor produces palette data, because no palette table has been located. Entries 1-3 coincide with the BG 2bpp palette 0 that EXP-0023 measured (BGR555 `$0000/$5294/$7FFF`), but black / mid-grey / white is the obvious three-level ramp and the coincidence is not evidence | 2026-08-02, Unit 11 | Locate the palette tables (readiness **F2**), add an extractor, and load real CGRAM data through `internal/content`. The sub-palette *names* (`content.SubPalettePrimary`, `SubPaletteDim`) are chosen to survive that change — only the colours behind them move | **Open** |

## Retired

| ID | Severity | Deviation | Introduced | How it was retired | Retired |
|---|---|---|---|---|---|
| D0 | `Parity-blocking` | **Text rendered black on black.** `framebuf.BlitTile` computes `dst = ink + PaletteBase`, so `PaletteBase` is a sub-palette base; both scenes passed `white = 3` / `gray = 2`, reading the field as a brightness level. The HUD font's ink values are 1-3, so bright strings landed on palette entries 4-6 while `GrayPalette` defined only 0-3. Measured on the real font: **32 % of drawn ink resolved to a visible colour**. The frame goldens could not see it — `Sum256` hashes indices and is deliberately palette-independent, so a frame can hash correctly and be blank on screen | Pre-existing (Unit 4, `2febc83`, 2026-08-02); identified 2026-08-02 during Unit 10 planning | The units were the defect, so the fix is a type: `content.SubPalette` is a **slot number**, `Base()` converts it, and `TextOptions.Palette` takes it. Expressed as slots the old values address entries 8-15, far outside what the palette defines, and `TestScenesDrawOnlyDefinedPaletteEntries` fails — verified by re-applying them. `GrayPalette` gained a second sub-palette so secondary text reads as subordinate rather than absent. Visible ink went from 32 % to **63 %** on the real font, with the ink mask byte-identical (6406 pixels before and after on the synthetic-font boot scene), so the change is index remapping only | 2026-08-02, Unit 11 |
| D1 | `Parity-blocking` | The `hud-font` extractor read `ROMFILE:0x046FC0` as a block start, but that address is only the arithmetic anchor for VRAM tile `$000`. The real 257-tile block is `ROMFILE:0x047FB0-0x048FBF` (ROM-0016). 255 of 257 tiles in `hud-font-sheet.png` were attack-table bytes rendered as tiles. `manifests/rom-regions.json` and `internal/extract/extractors.go` contradicted each other, and nothing compared them | Pre-existing (extractor v1.0.0, 2026-07-30); identified 2026-08-02 during DEMO-0001 planning | Extractor v2.0.0 derives `hudFontBase` from the anchor relation (`hudFontAnchor + hudFontFirstVRAMTile*16`) so the block start stays checkable rather than magic. Archive regenerated: exactly one asset changed, `archive verify` 8/8 clean, `rom_source` now `ROMFILE:0x047FB0-0x048FBF`. Regression tests `TestHUDFontMatchesROMLedger` (asserts the extractor span against ROM-0016) and `TestHUDFontAnchorRelation` (pins the affine relation) added — neither needs a ROM. GFX-0001 and CEN-GFX-0001 synchronized. Visually confirmed: the sheet now renders a legible font | 2026-08-02, Unit 1 |

## Explicitly out of scope (not deviations)

These are boundary decisions, recorded so they are not later mistaken for gaps:

| Item | Rationale |
|---|---|
| The frozen-Esper reaction and Terra's awakening | Out of SCN-0001's boundary; only the few frames proving the scenario ended are in scope |
| Event opcodes outside the DEMO-0001 route | The supported subset must be documented and tested; unsupported opcodes may remain unimplemented |
| Battle types other than those the route exercises | Random, scripted, and the Whelk boss battle are in scope; pincer/back attacks are not |
| Save/load to SRAM | The acceptance run starts at New Game and does not save |
