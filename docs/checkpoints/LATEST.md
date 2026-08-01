# Latest Checkpoint

**[2026-08-01 — EXP-0037: opening event-flag inventory](2026-08-01-exp0037-event-flag-inventory.md)**

State: the opening's event-flag behavior is inventoried and
deterministic — 20 flags (11 latched story, 4 transient, 5 working
bits), 162 value-changing writes, byte-identical across one GUI and
two headless runs, final WRAM = milestone 05 (five byte-identical
runs). Every writer statically decoded; 16-handler family over eight
bases; event interpreter anchored at candidate `ROMCPU:$C09B5C`;
GUI/testrunner parity verified for this schedule. New: DISC-0008,
`internal/game/eventflags`, ROM-0027..0032,
`data/scenarios/opening-event-flags.json`. Exact next action:
EXP-0038 — golden route segment 5, mines traversal to milestone
`06-random-encounter` with encounter-trigger context capture
(CEN-WORLD-0006).
