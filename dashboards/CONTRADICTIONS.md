# Contradictions

Two records, both resolved; none open (2026-08-01).

| ID | Status | Summary |
|---|---|---|
| [CONTRA-0001](../docs/contradictions/CONTRA-0001-address-notation.md) | Resolved 2026-07-29 | Pre-V4 contextual address notation vs V4 domain prefixes — V4 supersedes for new/updated docs; historical docs keep declared notation. |
| [CONTRA-0002](../docs/contradictions/CONTRA-0002-1ea5-map-id-vs-event-flags.md) | Resolved 2026-08-01 | `WRAM:+$1EA5` "map id" (EXP-0035) vs "map-load target" (EXP-0036) — **both refuted**. Static decode of the writer shows a bit-set into an event-flag array at `+$1EA0`; the location correlation was incidental to story progress. Event-flag system located as a by-product (CEN-EVENT-0008). |

> Note: this dashboard read "none open" while CEN-WORLD-0007 carried an
> unresolved contradiction field. Contradictions recorded inside census
> entries must be surfaced here as well — checked at each resolution.
