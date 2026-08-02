# Milestones

| ID | Milestone | Status |
|---|---|---|
| M0 | Reproducible lab | **Complete** (2026-07-29) — ROM identity, capability matrix, checkpoint, git repository published to `github.com/bensabler/ff6-decompile`; remaining nuance: re-verify Mesen version string on next launch ([BLOCKERS.md](BLOCKERS.md)) |
| M1 | First confirmed Go behavior | **Complete** (2026-07-29) — `chardata.CopyCharacterFields`, [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| M2 | First provenance-complete sprite/tile asset | **Complete** (2026-07-30) — battle HUD font tile block, [GFX-0001](../docs/graphics/GFX-0001-battle-hud-font.md) (ROM provenance Confirmed by byte identity; load-path trace still queued) |
| M3 | First provenance-complete audio sample/cue | **Complete** (2026-07-30) — battle confirm SFX, [AUD-0001](../docs/audio/AUD-0001-battle-confirm-sfx.md) (trigger chain + ROM provenance Confirmed; SPC driver dispatch still queued) |
| M4 | Party battle record | In progress — unified 10-slot arrays, battler stat tables, and the full damage pipeline decoded/implemented (EXP-0003..0019); remaining: MP verification, records' consumer, status semantics |
| M5 | Battle subsystem vertical slice | In progress — damage subsystem byte-exact formula-to-HUD; #30 resolved for the tested window (EXP-0021: action layer deterministic given the input frame schedule); battle configuration located and decoded (EXP-0041); remaining: the **ATB program** (timer domains, ACTIVE/WAIT, action queue), press-coupled state semantics (`+$3A71`) |
| M6 | Graphics subsystem vertical slice | In progress — first vertical proof complete (EXP-0023: HUD font runtime→ROM→decoder→comparison); next: load-path trace, glyph semantics, a compressed asset |
| M7 | Audio subsystem vertical slice | In progress — first vertical proof complete (EXP-0024: confirm SFX trigger→port→voice→sample→ROM); next: SPC driver dispatch, sequence decode, music cue |
| M8 | Public clean release | Not started |
| D1 | DEMO-0001 playable Go demo, New Game to Whelk victory | In progress (started 2026-08-02) — production milestone; the deliverable is a runnable program, not a record ([DEMO-0001](../docs/demo/DEMO-0001-new-game-to-whelk.md)). Readiness at start: **0 of 53 requirements Integrated**, 5 Implemented, 7 Evidence Ready, 33 Unknown ([matrix](../docs/demo/DEMO-0001-READINESS.md)). Build order is **evidence-led** (shell → battle vertical → field/event vertical), not route order, because the map system, dialogue corpus, event opcodes, sprites, audio sequences, and FF6's compression format are all at zero |
| S1 | SCN-0001 opening-to-Whelk golden route | In progress — milestones **00-06** established from power-on with byte-identical WRAM across repeated runs ([SCN-0001](../docs/scenarios/SCN-0001-opening-to-whelk.md)); remaining: 07-pre-whelk through 10-whelk-victory and branches A/B/C, all **blocked behind the ATB program** ([BLOCKERS](BLOCKERS.md)) |
