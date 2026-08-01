# Latest Checkpoint

**[2026-08-01 — CONTRA-0002 resolved: `+$1EA5` is an event-flag byte](2026-08-01-contra0002-event-flags.md)**

State: the `+$1EA5` contradiction is resolved — **both** the "map id"
(EXP-0035) and "map-load target" (EXP-0036) readings are refuted. A
static decode of the writer shows `ORA $C0BAFC,X / STA $1EA0,Y`: the
byte is byte 5 of an **event-flag bit array** at `WRAM:+$1EA0`, and its
values accumulate bits rather than naming a map. The event-flag system
(three arrays, set/clear routines, mask tables, index decoder at
`$BAED`) is located as a by-product — CEN-EVENT-0008. Dependent code
renamed and unfrozen; no golden-route result changed. Note:
CEN-WORLD-0004 (map header / tileset path) is now genuinely unstarted.
Exact next action: EXP-0037 — write-watch all three flag arrays across
the scheduled route to inventory which flags the opening sets, and when.
