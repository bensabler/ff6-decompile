# Statistics

As of 2026-07-30 (post-EXP-0020 maintenance; regenerate via `/weekly-review`):

| Metric | Count |
|---|---|
| Documented sessions | 3 (SES-001..003; SES-003 reconstructed post-hoc) |
| Function records | 9 (FN-0001..0009; FN-0003..0009 Confirmed code via EXP-0001/0003/0006/0009/0010) |
| Experiment records | 21 completed (EXP-0001..0021; EXP-0008 negative-result, EXP-0012 and EXP-0020 refuted hypotheses) |
| Variable records | 4 (VAR-0001..0004) |
| Structure records | 4 (ST-0001..0004) |
| Contradiction records | 1 resolved (CONTRA-0001, notation standards) |
| Go packages with tests | 7 of 13 (`internal/audit`, `internal/audio/brr`, `internal/game/attackdata`, `internal/game/battle`, `internal/game/chardata`, `internal/graphics/tile4bpp`, `internal/platform/bgr555`) |
| Go packages missing required tests | 0 |
| Quality gates (2026-07-30) | gofmt clean; `go build`/`go vet`/`go test ./...`/`ff6lab audit` pass |
