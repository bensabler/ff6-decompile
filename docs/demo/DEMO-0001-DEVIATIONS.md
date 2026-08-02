# DEMO-0001 deviations register

- **Milestone:** [DEMO-0001](DEMO-0001-new-game-to-whelk.md)
- **Updated:** 2026-08-02 (Unit 0 — program start)

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

## Retired

| ID | Severity | Deviation | Introduced | How it was retired | Retired |
|---|---|---|---|---|---|
| D1 | `Parity-blocking` | The `hud-font` extractor read `ROMFILE:0x046FC0` as a block start, but that address is only the arithmetic anchor for VRAM tile `$000`. The real 257-tile block is `ROMFILE:0x047FB0-0x048FBF` (ROM-0016). 255 of 257 tiles in `hud-font-sheet.png` were attack-table bytes rendered as tiles. `manifests/rom-regions.json` and `internal/extract/extractors.go` contradicted each other, and nothing compared them | Pre-existing (extractor v1.0.0, 2026-07-30); identified 2026-08-02 during DEMO-0001 planning | Extractor v2.0.0 derives `hudFontBase` from the anchor relation (`hudFontAnchor + hudFontFirstVRAMTile*16`) so the block start stays checkable rather than magic. Archive regenerated: exactly one asset changed, `archive verify` 8/8 clean, `rom_source` now `ROMFILE:0x047FB0-0x048FBF`. Regression tests `TestHUDFontMatchesROMLedger` (asserts the extractor span against ROM-0016) and `TestHUDFontAnchorRelation` (pins the affine relation) added — neither needs a ROM. GFX-0001 and CEN-GFX-0001 synchronized. Visually confirmed: the sheet now renders a legible font | 2026-08-02, Unit 1 |

## Explicitly out of scope (not deviations)

These are boundary decisions, recorded so they are not later mistaken for gaps:

| Item | Rationale |
|---|---|
| The frozen-Esper reaction and Terra's awakening | Out of SCN-0001's boundary; only the few frames proving the scenario ended are in scope |
| Event opcodes outside the DEMO-0001 route | The supported subset must be documented and tested; unsupported opcodes may remain unimplemented |
| Battle types other than those the route exercises | Random, scripted, and the Whelk boss battle are in scope; pincer/back attacks are not |
| Save/load to SRAM | The acceptance run starts at New Game and does not save |
