# Milestones

| ID | Milestone | Status |
|---|---|---|
| M0 | Reproducible lab | **Complete** (2026-07-29) — ROM identity, capability matrix, checkpoint, git repository published to `github.com/bensabler/ff6-decompile`; remaining nuance: re-verify Mesen version string on next launch ([BLOCKERS.md](BLOCKERS.md)) |
| M1 | First confirmed Go behavior | **Complete** (2026-07-29) — `chardata.CopyCharacterFields`, [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| M2 | First provenance-complete sprite/tile asset | Not started |
| M3 | First provenance-complete audio sample/cue | Not started |
| M4 | Party battle record | In progress — unified 10-slot arrays, battler stat tables, and the full damage pipeline decoded/implemented (EXP-0003..0019); remaining: MP verification, records' consumer, status semantics |
| M5 | Battle subsystem vertical slice | In progress — damage subsystem byte-exact formula-to-HUD; #30 resolved for the tested window (EXP-0021: action layer deterministic given the input frame schedule); remaining: press-coupled state semantics (`+$3A71`), GUI/testrunner parity |
| M6 | Graphics subsystem vertical slice | Not started |
| M7 | Audio subsystem vertical slice | Not started |
| M8 | Public clean release | Not started |
