# Statistics

As of 2026-07-29 (V4 migration; regenerate via `/weekly-review`):

| Metric | Count |
|---|---|
| Documented sessions | 2 (SES-001, SES-002) + 1 undocumented (SES-003) |
| Function records | 2 documented (FN-0001, FN-0002); 3 pending documentation (FN-0003..0005) |
| Variable records | 4 (VAR-0001..0004) |
| Structure records | 2 (ST-0001, ST-0002) |
| Contradiction records | 1 resolved (CONTRA-0001, notation standards) |
| Go packages with tests | 4 of 12 (`internal/game/chardata`, `internal/audio/brr`, `internal/graphics/tile4bpp`, `internal/platform/bgr555`) |
| Go packages missing required tests | 1 (`internal/game/battle` — implements behavior without tests) |
| Quality gates (2026-07-29) | gofmt clean; `go build`/`go vet`/`go test ./...` pass |
