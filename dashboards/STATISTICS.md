# Statistics

As of 2026-08-01 (post-audit refresh; regenerate via `/weekly-review`):

| Metric | Count |
|---|---|
| Documented sessions | 3 (SES-001..003; SES-003 reconstructed post-hoc) |
| Function records | 9 (FN-0001..0009; FN-0003..0009 Confirmed code via EXP-0001/0003/0006/0009/0010) |
| Experiment records | 36 completed (EXP-0001..0036; EXP-0008 negative-result; EXP-0012/0020 refuted hypotheses; EXP-0033/0035 completed-partial with carried work closed by EXP-0034/0036) |
| Graphics asset records | 1 (GFX-0001, ROM provenance Confirmed) |
| Audio asset records | 1 (AUD-0001, trigger chain + ROM provenance Confirmed) |
| Content census entries | 57 across 12 of 12 domains ([COVERAGE.md](COVERAGE.md)) |
| ROM ownership | 26 regions; 14,885 bytes known (0.47%), 120 candidate ([ROM_REGIONS.md](../indexes/ROM_REGIONS.md)) |
| Variable records | 4 (VAR-0001..0004) |
| Structure records | 4 (ST-0001..0004) |
| Contradiction records | 2 resolved (CONTRA-0001 notation; CONTRA-0002 $1EA5 event-flag byte), 0 open |
| Go packages with tests | 11 of 17 (adds `internal/scenario/route`; the 6 without tests are declaration-only or wiring packages) |
| Go packages missing required tests | 0 |
| Quality gates (2026-08-01) | gofmt clean; build/vet/test, `ff6lab audit`, `census validate`, and a clean tracked-only checkout build all pass |
