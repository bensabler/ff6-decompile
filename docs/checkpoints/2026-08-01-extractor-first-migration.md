# Checkpoint 2026-08-01 — extractor-first archival architecture (Unit 38)

## Current question
None open. The architecture migration is complete; the next research
unit is EXP-0037 (opening event-flag inventory).

## State
The project is now **extractor-first**. The public repository carries
the reconstruction; substantial assets are regenerated locally from a
user-supplied verified ROM and described by a tracked manifest. A fresh
clone plus the ROM reproduces the archive byte-for-byte.

## Work completed

### Policy correction
`docs/legal/ASSET_POLICY.md` rewritten. Names, labels, IDs, statistics,
mechanical data, table contents and complete structured gameplay
databases are **permitted in Git**. Substantial graphics, audio,
dialogue, maps, scripts and binaries are generated locally. The policy
explicitly may not be used to hide ordinary reconstruction data behind
placeholders, and that rule is now mechanically enforced.
`.claude/skills/_shared/LEGAL_BOUNDARY.md` and `CLAUDE.md` updated to
match (CLAUDE.md gains an Asset policy section and two new gates).

### Placeholder migration
All **54 spell names restored** to `data/census/spells.json`. They were
the only `see local extraction (asset policy)` placeholders in the
tracked tree — a repo-wide sweep of `data/ manifests/ docs/ indexes/
dashboards/` confirms none remain.

Names are decoded from the ROM by the new extractor and were
cross-checked against EXP-0027's independent local extraction: 54/54
agree on content. **One rendering discrepancy resolved with evidence** —
EXP-0027's script rendered byte `$FE` as `.` ("Fire.2"), while EXP-0026
Confirmed `$FE` as a *narrow space*. The bytes are identical; the space
rendering is correct, so the inventory reads "Fire 2". Recorded in the
textenc tests so it cannot silently regress.

### New packages
- `internal/rom` — loads and identity-checks the ROM; refuses wrong
  size or hash outright; `FF6_ROM` / `-rom` resolution.
- `internal/game/textenc` — the EXP-0026 fixed-width text encoding.
  Unmapped bytes decode to visible `\xNN` escapes; nothing is guessed.
- `internal/extract` — extractor registry, archive writer with
  overwrite safety, asset manifest, verification, inventory.

### Extractors (8 assets, deterministic)
`spells` (54 name+data records), `espers` (27 names), `hud-font`
(257-tile 2bpp sheet → PNG using the EXP-0023 live palette),
`sfx-pack` (288-byte BRR pack + sample 5 decoded to PCM via
`internal/audio/brr`), `raw-tables` (monster, formation, formation-flag
slices). `dialogue` and `maps` categories are registered but report
honestly that no confirmed source is located yet (CEN-EVENT-0002,
CEN-WORLD-0004) rather than emitting fabricated output.

### Commands
`ff6lab extract all|<category>... [-rom p] [-force]` and
`ff6lab archive verify|inventory`. Extraction verifies the ROM hash,
refuses unsupported revisions, writes only under the archive root,
refuses to touch `local_artifacts/experiments/`, and never silently
overwrites a differing file. `archive verify` compares three ways:
current extractor output, tracked manifest hash, on-disk bytes — and
reports missing / changed / drifted / unknown separately.

### Tests
- `TestNoPlaceholdersInPermittedNameFields` — walks all tracked
  data/manifest JSON; fails on placeholder markers in name fields.
  Allows honest `unknown`/empty, forbids policy excuses.
- `TestSpellInventoryIsComplete` — 54 records, non-empty, distinct.
- `TestExtractionIsDeterministic` — four passes, byte-identical.
- `TestAssetsMatchTrackedManifest` — manifest describes reality.
- `TestManifestEntriesAreComplete` — full provenance per entry
  (runs without a ROM).
- `textenc` table-driven decode tests including the `$FE` case.

ROM-dependent tests skip cleanly when `FF6_ROM` is unset, so CI without
a ROM still runs the whole non-ROM suite.

## What remains uncertain
- `dialogue` and `maps` have no extractors because their sources are
  unlocated — a real research gap, not an architectural one.
- The monster/formation raw slices are archived as bytes because their
  record layouts are only partly decoded; per-record inventories should
  replace them as decoding advances.
- Music sequences, sprites, portraits, backgrounds and animation frames
  have no extractors yet (sources unlocated).

## Tests and quality gates
gofmt clean; build/vet clean; `go test ./...` green (13 packages);
`ff6lab audit` clean; `census validate` clean; `archive verify` clean
(8/8); restricted-extension scan clean; unsupported-ROM refusal
verified; clean tracked-only checkout verified.

## Git status
main; migration committed as one unit and pushed.

## Exact next action
EXP-0037: write-watch the event-flag arrays (`WRAM:+$1E80`, `+$1EA0`,
`+$1EC0`) across the scheduled golden route and inventory which flag
numbers the opening sets and when (CEN-EVENT-0008). Unchanged by this
migration; the archive now gives that unit a place to put extracted
script material once the event dispatcher is located.
