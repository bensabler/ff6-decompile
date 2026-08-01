# Latest Checkpoint

**[2026-08-01 — EXP-0036: scheduled route to the mines (partial)](2026-08-01-exp0036-scheduled-route.md)**

State: a 17-leg state-driven route controller walks power-on → mines
interior with no manual correction (model + tests in
`internal/scenario/route`, probe-sync guard against the Lua). Battle 5
= formation 84 {27,27,0,0}, ROM-verified; new pre-Whelk monster record
27. The `+$1EA5` falsifier **fired** — it reaches the mines value while
the party is still on the exterior, so it is not a simple map-id byte.
Milestone 05 **not claimed**: three acceptance runs required, and the
final encoding has not completed that set. Exact next action: three
power-on runs of `mesen/probes/EXP-0036.lua`, byte-compare milestone
WRAM, then promote milestone 05 only if all three land at (`$26`,`$1C`).
