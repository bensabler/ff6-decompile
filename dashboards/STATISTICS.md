# Statistics

As of 2026-07-30 (post-EXP-0020 maintenance; regenerate via `/weekly-review`):

| Metric | Count |
|---|---|
| Documented sessions | 3 (SES-001..003; SES-003 reconstructed post-hoc) |
| Function records | 9 (FN-0001..0009; FN-0003..0009 Confirmed code via EXP-0001/0003/0006/0009/0010) |
| Experiment records | 26 completed (EXP-0001..0026; EXP-0008 negative-result, EXP-0012 and EXP-0020 refuted hypotheses) |
| Graphics asset records | 1 (GFX-0001, ROM provenance Confirmed) |
| Audio asset records | 1 (AUD-0001, trigger chain + ROM provenance Confirmed) |
| Content census entries | 43 across 10 of 12 domains ([COVERAGE.md](COVERAGE.md)) |
| ROM ownership | 20 regions; 10,753 bytes known (0.34%), 120 candidate ([ROM_REGIONS.md](../indexes/ROM_REGIONS.md)) |
| Variable records | 4 (VAR-0001..0004) |
| Structure records | 4 (ST-0001..0004) |
| Contradiction records | 1 resolved (CONTRA-0001, notation standards) |
| Go packages with tests | 10 of 15 (`cmd/ff6lab`, `internal/audit`, `internal/audio/brr`, `internal/census`, `internal/game/attackdata`, `internal/game/battle`, `internal/game/chardata`, `internal/graphics/tile2bpp`, `internal/graphics/tile4bpp`, `internal/platform/bgr555`) |
| Go packages missing required tests | 0 |
| Quality gates (2026-07-30) | gofmt clean; `go build`/`go vet`/`go test ./...`/`ff6lab audit` pass |
