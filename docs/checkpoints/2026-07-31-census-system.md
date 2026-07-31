# Checkpoint 2026-07-31 — Content census and coverage system established

## What was created
- `docs/research/CONTENT_TAXONOMY.md`: 12-domain taxonomy; two
  independent status ladders (reconstruction vs runtime), never
  collapsed.
- `schemas/content-census.schema.json`, `schemas/rom-regions.schema.json`.
- `manifests/content-census.json` (43 entries) and
  `manifests/rom-regions.json` (20 regions) — seeded strictly from
  existing evidence; every discovery now has census ownership.
- `internal/census`: load/validate (duplicate ids, enum validity,
  evidence/implementation/test ladder requirements, cross-links,
  undeclared ROM overlaps, generated-file staleness), coverage
  summaries, ROM gap analysis, and generators for
  `indexes/CONTENT_CENSUS.md`, `indexes/ROM_REGIONS.md`,
  `dashboards/COVERAGE.md`.
- `ff6lab census validate|sync`, `coverage summary|gaps|domain`,
  `rom gaps`; census validation joined `ff6lab audit`.
- `data/census/`: 19 data-family inventories (spells populated; the
  rest hold descriptions + bounded next steps; no fabricated
  records).
- Workflow: census-observer skill, `/census-observations`,
  `/register-system`, `/update-coverage`; stopping rules and
  CLAUDE.md now require the observation pass before closing an
  experiment; the orchestrator playbook gained depth/breadth
  prioritization rules (max three consecutive same-subsystem
  experiments before a coverage review).

## What was actually discovered
- **EXP-0025** (screenshot sweep): all four opening states render
  headlessly; the field state is controllable; Terra renders as
  `?????` pre-naming; 25 systems registered at honest statuses.
- **EXP-0026** (magic census): fixed-font text encoding derived from
  the menu tilemap (A-Z=$80+n, a-z=$9A+n, '-'=$C4); spell name table
  at `ROMFILE:0x26F567` (7-byte `[icon][name]`, Fire=0, Cure=45);
  the `$C46AC0` table **is** the global spell database with record
  byte 5 = MP cost (Strong hypothesis, one on-screen anchor);
  per-character spell availability at `WRAM:+$1A6E` (unique 128 KB
  search hit). Preliminary inventory in `data/census/spells.json`
  (structure + 2 observed records).

## What remains unknown
- ROM: 3,134,855 bytes (99.66%) unowned — largest gap 2.8 MB from
  `0x051FE9`. ITEM and QUIRK domains have no entries (nothing
  observed). Spell table: end boundary, help-text table, targeting/
  animation/sound linkage, menu grid ordering, field castability
  rule. Fire's cost needs a second on-screen anchor.

## Verification
gofmt clean; build/vet pass; `go test ./...` 10 packages;
`ff6lab audit` clean (census checks included); clean tracked-only
clone green (run at commit 5).

## Git status
`main`, 5 census commits ahead of the previous push point after this
commit. Not pushed (push only on request).

## Exact next breadth target
Verify the spell name-table boundary + a second MP-cost anchor
(battle Magic menu navigation), then bulk-extract all 54 spell
records to `local_artifacts/` and mirror ids/numbers into
`data/census/spells.json` — the first EXTRACTED_COMPLETE candidate.
