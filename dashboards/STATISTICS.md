# Statistics

As of 2026-08-02 (Unit 10 records reconciliation; regenerate via
`/weekly-review`). Counts are recounted from the tracked records, manifests and
indexes rather than edited incrementally.

| Metric | Count |
|---|---|
| Documented sessions | 3 (SES-001..003; SES-003 reconstructed post-hoc) |
| Function records | 9 (FN-0001..0009; FN-0003..0009 Confirmed code via EXP-0001/0003/0006/0009/0010) |
| Experiment records | 48 (EXP-0001..0047 and EXP-0049; **EXP-0048 is queued, no record exists**). EXP-0008/0020 negative results; EXP-0040 `blocked`, no Whelk victory |
| Discovery records | 8 (DISC-0001..0008, all Confirmed) |
| Correlation records | 1 ([CORR-0001](../docs/correlations/CORR-0001-C09B5C.md); mechanism Confirmed, semantic role Strong hypothesis) |
| Graphics asset records | 1 (GFX-0001, ROM provenance Confirmed; archive-vs-ROM differential passes on all 256 glyphs) |
| Audio asset records | 1 (AUD-0001, trigger chain + ROM provenance Confirmed) |
| Content census entries | 68 across 12 of 12 domains ([COVERAGE.md](COVERAGE.md)) |
| ROM ownership | 32 regions; 15,496 bytes known (0.49 %), 120 candidate, 3,130,112 unknown in 27 gaps ([ROM_REGIONS.md](../indexes/ROM_REGIONS.md)) |
| Variable records | 4 (VAR-0001..0004) |
| Structure records | 4 (ST-0001..0004) |
| Contradiction records | 2 resolved (CONTRA-0001 notation; CONTRA-0002 `$1EA5` event-flag byte), 0 open |
| Architecture decisions | 1 ([ADR-0001](../docs/decisions/ADR-0001-rendering-host.md) rendering host) |
| Generated assets | 8 in `manifests/assets.json`; `archive verify` 8/8. Categories `dialogue`, `maps`, `animations`, `scripts` are empty — what holds acceptance gate G4 at PARTIAL |
| Go packages | 34; **26 carry tests**; 14 fuzz targets |
| Go packages missing required tests | 0 |
| DEMO-0001 readiness | 57 rows — 1 Validated, 14 Integrated, 6 Implemented, 1 Evidence Ready, 1 Extractor Ready, 2 Researching, 2 Blocked, 1 Deferred, 29 Unknown ([matrix](../docs/demo/DEMO-0001-READINESS.md)) |
| DEMO-0001 acceptance | 1 of 17 progression steps; gates 4 PASS, 1 PARTIAL, 1 FAIL ([scorecard](../docs/demo/DEMO-0001-ACCEPTANCE.md)) |
| Route dependency pressure | compression (X1) blocks **8 of 19** beats, the largest single blocker ([content matrix](../docs/demo/DEMO-0001-CONTENT-MATRIX.md)) |
| Quality gates (2026-08-02) | gofmt clean; build/vet/test on both build variants, `ff6lab audit`, `census validate`, `archive verify` 8/8 all pass |

## Correction

The previous revision, dated 2026-08-01, reported 36 experiments, 57 census
entries, 26 ROM regions and 0.47 % ownership, and "11 of 17" Go packages with
tests. Each was stale. This revision recounts from `docs/experiments/`,
`docs/discoveries/`, `manifests/`, `indexes/ROM_REGIONS.md` and `go list`.
