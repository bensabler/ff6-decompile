# Checkpoint 2026-08-01 — EXP-0044: the ACTIVE/WAIT pause matrix (Unit 46)

## Current question

None open in the unit. EXP-0044 is **completed** and the ATB blocker
raised by EXP-0040 is **discharged for its original purpose**.

## State

A usable ATB model exists. Across EXP-0041..0044 the project now has:
configuration storage and encoding, battle-entry sampling, gauges,
increments, the tick counter, and the exact ACTIVE/WAIT pause condition.

**Whelk is no longer blocked by an absent model.** Its intervals can now
be scoped: only time inside the ability list and target selection was
paused.

## Confirmed before this session

Config storage (EXP-0041, CEN-MENU-0007); battle-entry sampling and the
staging rule (EXP-0042, CEN-BATTLE-0010); the ATB layer and the gate
instruction (EXP-0043, CEN-BATTLE-0011).

## Work completed

**`WRAM:+$2F41` is the battle submenu flag.** Resting `$00`, cleared
per-frame at `ROMCPU:$C17A92` (`STZ`), raised at `ROMCPU:$C17C01` (`INC`)
when a qualifying submenu opens — it fired **exactly twice** for two
submenu opens. A second clear path exists at `ROMCPU:$C14434`.

`ROMCPU:$C21124` ANDs it with the Wait flag `+$3A8F` and skips the entire
per-frame battle update.

### The matrix

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| **Main command window open** | advances | **advances** |
| Ability list open | advances* | **paused** |
| Target selection | advances* | **paused** |
| Action resolving / animation | advances* | **advances** |
| Item / Magic / Row / Defend, dialogue, damage display, victory, defeat | not sampled | not sampled |

\* structurally implied — the gate is an `AND`, so `+$3A8F` = `$00` makes
it zero regardless of `$2F41`. Directly verified for the ability-list row.

All four located domains — tick `+$3A3E`, gauges `+$3AB4`, flags
`+$3AA0`, accumulator `+$3218` — froze and resumed **together** in every
trial. No independent clock was found among them, which is itself a
result: the master programme specifically warned to look for one.

**The pause is narrower than the folk model.** Two states that intuition
expects to pause do not: the main battle command window, which is on
screen for most of a WAIT battle, and action resolution. Only a
qualifying submenu raises the flag.

### Method

The mode was flipped **in place** by patching `+$3A8F`, which EXP-0042
showed is what battle timing actually runs on. That made
ACTIVE-versus-WAIT a genuine one-variable comparison inside a single
savestate — same submenu, `$2F41` = `01` throughout, frozen at Wait and
resuming at Active — instead of two route runs with divergent RNG.

The patch was then validated against **genuinely configured WAIT**
(`in-battle-formation14.mss`), which froze identically. Falsifier 2 did
not fire, so every patched trial stands. Full runtime-patch record is in
the experiment record.

Preparatory commit: `/battle-baseline` plus `battleconfig()` in
`probes/common.lua` — deferred at EXP-0041 precisely because it could
then only have re-read the Config screen from pixels.

## Last raw observation

Trial 9, genuinely configured WAIT, MagiTek list open:
`gate2F41=01 wait3A8F=01 spd3A90=CF`, tick frozen at `$17C7` across three
samples.

## Active emulator state

**None.** One headless `--testrunner` instance; logs harvested before
termination, killed with `kill -9`, absence confirmed by `pgrep`.
`jobs -l` empty. **SRAM directory verified still empty.**

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0044.lua` armed a write watch on
`WRAM:+$2F41` plus an all-domain sampler and an in-place mode patch; all
die with the process.

## Evidence paths

`local_artifacts/experiments/EXP-0044/` — **10 artifacts, 124 KB**, with
`hashes.sha256`: 6 screenshots, `gate-writers.txt`,
`bridge-commands.log`, `bridge-events.log`, `experiment.json`.

Starting state: `local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss`,
cross-check `in-battle-formation14.mss`. No route replay needed.

## Files changed

- New: `docs/experiments/EXP-0044-active-wait-pause-matrix.md`,
  `mesen/probes/EXP-0044.lua`, `.claude/commands/battle-baseline.md`,
  this checkpoint.
- Updated: `mesen/probes/common.lua` (+`battleconfig`),
  `docs/sessions/04_MEMORY_MAP.md`, `manifests/experiments.json`,
  `manifests/content-census.json`, `dashboards/` (CURRENT_FOCUS,
  BLOCKERS, RESEARCH_QUEUE, ACTIVITY_LOG), `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.

Stale claim removed: CURRENT_FOCUS's "Blocked — no ATB model" section is
now struck and pointed at the discharge note, per the dashboard rule that
no dashboard may retain a claim contradicted by newer evidence.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test ./...` ok;
`ff6lab audit` clean; `ff6lab census sync` clean (65 entries);
restricted-extension scan clean. No Go source changed this unit.

## Git status

`main`. Two commits: the `/battle-baseline` preparation, then this unit.
Worktree clean.

## Unresolved decisions

- **Six matrix rows are `not sampled`** — Item, Magic, Row, Defend,
  dialogue, damage display, victory and defeat presentations.
- **The settling transient.** Just after the gate engages, one sample
  pair showed slot 6's flag byte clearing and `+$3218` advancing by
  `$0100` before everything froze for the following 920 frames. At this
  sampling resolution that cannot be distinguished from work completing
  on the last un-gated frame.
- Whether any **unlocated** domain (animation, AI script, status,
  boss-state timers) advances during a paused interval. "Everything
  pauses together" means everything the project can currently see.
- `ROMCPU:$C14434`, the second clear path, fired once and is unexplained.
- Battle types other than random encounters — **Whelk is a boss with its
  own script**, the one caveat for any Whelk re-run.

## Blockers

**The ATB blocker is discharged.** Configuration, entry sampling, gauges,
increments, tick counter and the pause condition are all located and
verified. What remains — the increment formula, the action queue, status
modifiers, the unsampled matrix rows — is ordinary queued work, not a
blocker.

EXP-0040's head/shell transitions remain **unusable as natural timing**,
but they are now *scoped* rather than dismissed: intervals inside the
ability list and target selection were paused; intervals at the command
window and during action resolution were not.

This unit does **not** support, and no downstream record may assert:
that any untested menu state pauses or does not pause; that unlocated
timer domains obey this gate; that queued work does or does not resolve
past the gate; or that a boss battle behaves like a random encounter.

## Exact next action

**EXP-0045 — finish the matrix and settle the transient.** Walk every
remaining battle menu and presentation state sampling `+$2F41`, and
frame-step the gate transition using an `emu.addEventCallback` on
`endFrame` rather than bridge round-trips, which are orders of magnitude
too coarse for a one-frame question.

Then the ATB programme's remaining questions — increment formula and
threshold, action queue and readiness arbitration, status modifiers —
are all non-blocking, and the Whelk decision (reinterpret EXP-0040's
scoped captures, or re-run the fight with the model in hand) becomes an
ordinary orchestration call.

## Recommended next command

`/orchestrate`
