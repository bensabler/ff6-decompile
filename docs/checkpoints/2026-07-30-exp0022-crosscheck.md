# Checkpoint 2026-07-30 — EXP-0022: attack-record ROM cross-check (Unit 3)

## Current question
None open in this unit. Next per operator rebalance: Unit 4 graphics
vertical, Unit 5 audio vertical, then final synchronization.

## State
Headless Mesen testrunner still live (bridge v2, EXP-0021 probe armed,
`--timeout=7200`). Evidence frozen at
`local_artifacts/experiments/EXP-0022/` (table dump, scan output,
hashes). EXP-0021 archive frozen earlier tonight.

## Work completed
EXP-0022: 256 records dumped live from `ROMCPU:$C46AC0`, decoded via
the new `ff6lab attackdata scan` subcommand (hex parser table-tested).
Record 238: power 0 ✓ (EXP-0018 v=0 cross-check) and physical-formula
flag set ✓ (EXP-0017 convergence) — Confirmed. Fire Beam signature
(power 60, element bit 0, standard formula) → candidates 5 and 131
(Tentative, value coincidence). CLI help synchronized. MP verification
recorded as savestate-blocked.

## Tests and quality gates
Run at commit: gofmt clean, build/vet pass, `go test ./...` (8
packages incl. new cmd tests), `ff6lab audit` clean.

## Git status
`main`, 8 ahead of origin after this commit. Not pushed.

## Blockers
None hard. Soft: MP savestate; GUI/testrunner parity; Mesen version
string.

## Exact next action
Unit 4 (graphics vertical): from the live battle savestate, capture
the HUD/menu-font layer — PPU mode + BG chr/tilemap bases, dump VRAM
tile region + CGRAM palette via the bridge, locate the ROM source of
the font/HUD tiles (DMA provenance or byte search), decode with
`internal/graphics/tile4bpp` (or 2bpp as evidence dictates), and
compare decoder output against the VRAM bytes. Plan record
EXP-0023 first (question, falsifier, evidence paths).
