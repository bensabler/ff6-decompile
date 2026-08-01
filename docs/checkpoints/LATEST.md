# Latest Checkpoint

**[2026-08-01 — EXP-0036: scheduled route to the mines](2026-08-01-exp0036-scheduled-route.md)**

State: milestone `05-mines-entry` **established** — three power-on runs
of a 17-leg state-driven route controller all reach (`$26`,`$1C`)
inside the mines with byte-identical WRAM. The golden route now spans
power-on → milestones 00–05. Battle 5 = formation 84 {27,27,0,0},
ROM-verified. The `+$1EA5` falsifier fired: it reaches the mines value
before the transition is visible, so it is not a simple map-id byte.
Exact next action: EXP-0037 — write-watch `ROMCPU:$C0B5B6` across the
reproducible transition to settle `+$1EA5` and locate the map header /
tileset load path (CEN-WORLD-0004/0007).
