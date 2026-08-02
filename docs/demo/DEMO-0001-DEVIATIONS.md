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
| D1 | `Parity-blocking` | The `hud-font` extractor reads `ROMFILE:0x046FC0` as a block start, but that address is only the arithmetic anchor for tile index 0. The real 257-tile block is `ROMFILE:0x047FB0-0x048FBF` (ROM-0016). 255 of 257 tiles in `hud-font-sheet.png` are attack-table bytes rendered as tiles. `manifests/rom-regions.json` and `internal/extract/extractors.go` currently contradict each other | Pre-existing; identified 2026-08-02 during DEMO-0001 planning | Set `hudFontBase = 0x047FB0`, bump the extractor version, regenerate the archive, and add a skip-guarded differential test decoding `ROMFILE:0x047FB0` with `tile2bpp` and comparing tile-for-tile against the loaded PNG | **Open** — Unit 1 |

## Retired

_None yet._

## Explicitly out of scope (not deviations)

These are boundary decisions, recorded so they are not later mistaken for gaps:

| Item | Rationale |
|---|---|
| The frozen-Esper reaction and Terra's awakening | Out of SCN-0001's boundary; only the few frames proving the scenario ended are in scope |
| Event opcodes outside the DEMO-0001 route | The supported subset must be documented and tested; unsupported opcodes may remain unimplemented |
| Battle types other than those the route exercises | Random, scripted, and the Whelk boss battle are in scope; pincer/back attacks are not |
| Save/load to SRAM | The acceptance run starts at New Game and does not save |
