# Latest Checkpoint

**[2026-08-01 — EXP-0044: the ACTIVE/WAIT pause matrix](2026-08-01-exp0044-active-wait-pause-matrix.md)**

State: **the ATB blocker raised by EXP-0040 is discharged.** Across
EXP-0041..0044 the project now has configuration storage and encoding,
battle-entry sampling, gauges, increments, the tick counter, and the
exact ACTIVE/WAIT pause condition.

`WRAM:+$2F41` is the **battle submenu flag** — resting `$00`, cleared
per-frame at `ROMCPU:$C17A92`, raised at `ROMCPU:$C17C01` when a
qualifying submenu opens. `ROMCPU:$C21124` ANDs it with the Wait flag
`+$3A8F` and skips the entire per-frame battle update.

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| **Main command window open** | advances | **advances** |
| Ability list open | advances* | **paused** |
| Target selection | advances* | **paused** |
| Action resolving / animation | advances* | **advances** |
| Item / Magic / Row / Defend, dialogue, damage display, victory, defeat | not sampled | not sampled |

\* structurally implied — the gate is an `AND`; verified directly for the
ability-list row.

All four located domains (`+$3A3E`, `+$3AB4`, `+$3AA0`, `+$3218`) froze
and resumed **together**; no independent clock was found among them.
**The pause is narrower than the folk model**: the command window is on
screen for most of a WAIT battle and does not pause, and neither does
action resolution.

Method: the mode was flipped **in place** by patching `+$3A8F`, making
ACTIVE-versus-WAIT a one-variable comparison inside a single savestate;
the patch was then validated against genuinely configured WAIT, which
froze identically.

**Whelk is no longer blocked by an absent model.** EXP-0040's timing is
now *scoped* rather than dismissed: only intervals inside the ability
list and target selection were paused.

Carried forward honestly: six matrix rows are `not sampled`; a settling
transient just after the gate engages is unresolved at this sampling
resolution; unlocated domains (animation, AI script, status, boss-state)
were not covered; and Whelk is a boss with its own script, untested here.

10 evidence artifacts with hashes; no background processes; SRAM virgin.
All gates clean (gofmt/build/vet/test, `ff6lab audit`, census sync 65
entries, restricted-file scan).

Exact next action: **EXP-0045 — finish the matrix and settle the
transient.** Walk the remaining menu and presentation states sampling
`+$2F41`, and frame-step the gate transition with an `endFrame` callback
rather than bridge round-trips.
