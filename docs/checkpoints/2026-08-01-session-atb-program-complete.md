# Checkpoint 2026-08-01 — session close: the ATB research program (Units 43–49)

Session-level checkpoint spanning ten commits. Per-unit checkpoints
remain authoritative for their own detail; this one exists so a resuming
engineer sees the arc without reading seven files.

## Current question

None open. The ATB research program is **complete for its purpose** and
the blocker EXP-0040 raised is **discharged**. One narrow follow-up is
queued (EXP-0048) and it blocks nothing.

## State

The project can now interpret, reproduce and stage battle timing under
documented configurations — which is exactly the completion criterion the
program was given. Whelk is no longer blocked by an absent model.

Nothing is in flight. No emulator running, no instrumentation resident,
worktree clean, `main` in sync with `origin/main`.

## Confirmed before this session

Golden route power-on → milestone `06-random-encounter`; event-flag
system (DISC-0008); unified battle arrays (DISC-0001); formation table;
Whelk reached, entered, not defeated (EXP-0039/0040).

## Work completed

Ten commits: three infrastructure, seven experiments.

| Commit | Unit | Result |
|---|---|---|
| `8a24f17` | workflow prep | battle-config fingerprint became a required, audited record field |
| `dfdf896` | tooling | `ff6lab state` reads work/save RAM out of preserved `.mss` files |
| `72c8985` | EXP-0041 | configuration located and bit-decoded across three non-contiguous bytes |
| `1ca0375` | EXP-0042 | configuration sampled **once at battle entry**; staging rule established |
| `2cf84f6` | EXP-0043 | **the ATB layer located** — gauges, increments, flags, tick counter, gate |
| `a14a801` | tooling | `/battle-baseline`, once real addresses existed to read |
| `ff4e663` | EXP-0044 | **ACTIVE/WAIT pause matrix**; blocker discharged |
| `cee354c` | EXP-0045 | queued work resolves past the gate; `+$3AA0` bit 6 = pending |
| `adfee7d` | EXP-0046 | completion write `$C201BE`; shared-helper correction to EXP-0043 |
| `22b52af` | EXP-0047 | execution path is **periodic and ungated**; two refutations |

The model itself is summarised in `dashboards/CURRENT_FOCUS.md` and
carried in CEN-MENU-0007 and CEN-BATTLE-0010..0013. It is not repeated
here.

### Methodological notes worth carrying forward

- **Two predictions were made and confirmed.** EXP-0042 predicted
  `$3A8F`/`$3A90` from a static decode before running the second
  configuration, and both matched exactly. EXP-0045 predicted that
  engaging the gate while a slot was pending would produce a deferred
  completion, and it fired on a different slot at a different delay.
  Those are the two strongest results of the session.
- **One-variable comparison by patching.** EXP-0044 flipped ACTIVE/WAIT
  by writing `+$3A8F` in place, inside a single savestate, then validated
  the patch against genuinely configured WAIT. That avoided two route
  runs with divergent RNG.
- **Retrospective mining beat re-running.** EXP-0041's trial 0 extracted
  the `+$1D4E` Cursor candidate from EXP-0040's preserved savestate pair
  *before Mesen was ever launched*.
- **A `JSR` at `return − 3` is necessary but not sufficient** to confirm a
  stack frame. `$C20EB6` is not a call frame; `$C20016` holds a plausible
  `JSR` and never executes. Confirm frames by **execution**. This cost
  one unit's reconstruction and is recorded in CEN-BATTLE-0013.
- **A correction was propagated, not buried.** EXP-0043 attributed a
  `+$3AA0` store to the scheduler; EXP-0046 showed `$C211B4` is a shared
  helper with at least two entry points. The memory map is amended and no
  earlier conclusion depended on the stronger reading.

## Last raw observation

`EXEC C201BE gate=1 pc=$C201BE A=$0001 X=$000E Y=$0000 SP=$15F3 PS=$37
DB=$7E frame=76208` — the completion write executing while the
ACTIVE/WAIT gate was shut.

## Active emulator state

**None.** Every unit launched one headless `--testrunner` instance,
harvested logs before termination, killed with `kill -9`, and confirmed
absence with `pgrep`. `jobs -l` empty. **SRAM directory verified still
empty** — virgin boot preserved across the whole session.

## Breakpoints/watchers

None resident. All watches lived in per-unit probes
(`mesen/probes/EXP-0041..0047.lua`, none loaded now) plus `eval`-installed
callbacks, all of which die with the process. `mesen/probes/common.lua`
gained `battleconfig`, `watchreads` and `watchdump` as tracked shared
helpers.

## Evidence paths

| Unit | Directory | Artifacts |
|---|---|---|
| EXP-0041 | `local_artifacts/experiments/EXP-0041/` | 17, 616 KB |
| EXP-0042 | `.../EXP-0042/` | 8, 304 KB (two **live battle savestates**, reused by every later unit) |
| EXP-0043 | `.../EXP-0043/` | 8, 400 KB |
| EXP-0044 | `.../EXP-0044/` | 10, 124 KB |
| EXP-0045 | `.../EXP-0045/` | 12, 304 KB (four per-frame trace logs) |
| EXP-0046 | `.../EXP-0046/` | 5, 28 KB |
| EXP-0047 | `.../EXP-0047/` | 4, 24 KB |

Each frozen with `hashes.sha256`. Nothing restricted is tracked.

## Files changed

New tracked code: `internal/mesenstate/` (+tests), `cmd/ff6lab/state.go`
(+tests), `.claude/commands/battle-baseline.md`,
`mesen/probes/EXP-0041..0047.lua`, seven experiment records, seven
per-unit checkpoints, this checkpoint.

Updated: `internal/audit/audit.go` (+`CheckBattleExperimentConfig`),
`internal/project/project.go`, `schemas/experiment.schema.json`,
`.claude/templates/EXPERIMENT.md`, two skills,
`.claude/skills/_shared/ADDRESS_SPACES.md`,
`docs/research/ADDRESS_NOTATION.md` (SRAM prefix),
`docs/sessions/04_MEMORY_MAP.md` (≈20 entries), `manifests/*`,
`dashboards/*`.

## Tests and quality gates

Run before every commit and green at each: `gofmt -l .` clean;
`go build ./...`; `go vet ./...`; `go test ./...` (15 packages);
`ff6lab audit` clean; `ff6lab census sync` clean (66 entries);
`ff6lab archive verify` with `FF6_ROM` — 8/8 ok; restricted-extension
scan of `git ls-files` clean.

The audit earned its keep twice: it rejected a battle experiment missing
its configuration fingerprint (the check added this session), and caught
a dangling `LATEST.md` link written before its checkpoint existed.

## Git status

`main`, in sync with `origin/main`, worktree clean. Ten commits pushed.

## Unresolved decisions

- **The invoker of the action-execution path** — two candidates refuted.
- The invocation period beyond one run (gaps 121/35/105/122).
- `$C223ED`, `$C2083F`, `$C213D3` — named call targets, none decoded.
- The exact increment formula, threshold, overflow and gauge reset.
- `+$3AA0` bits 0-2, 4, 5; `+$3204`.
- Queue ordering and arbitration when several slots are ready.
- Status modifiers (Haste, Slow, Stop) — never touched.
- **Battle types other than random encounters.** Every ATB unit ran on
  formation 14. Whelk is a boss with its own script.
- Whether party slots behave like enemy slots at a gate engage.
- Six matrix rows for a normal-party battle (Magic, Row, Defend) and the
  defeat presentation with gate instrumentation.

## Blockers

**None.** The ATB blocker is discharged.

Standing evidence constraint, unchanged: EXP-0040's head/shell
transitions may not be used as *natural* Whelk timing. They can now be
**scoped** — intervals inside the ability list and target selection were
paused; intervals at the command window and during action resolution were
not — but scoping is not the same as reinterpreting, and no record has
done the latter.

No downstream record may assert: that the invocation period is constant;
that these findings hold in scripted, pincer, back or boss battles; that
party slots behave like enemy slots; or that any address is on a call
path merely because a `JSR` sits at `return − 3`.

## Exact next action

**EXP-0048 — name the invoker of the action-execution path.** Stack
archaeology has failed twice, so change instrument: exec-watch outward
from the confirmed `ROMCPU:$C2141D`, or trace a single invocation with
Mesen's trace facility. Start from
`local_artifacts/experiments/EXP-0042/in-battle-active-speed6.mss`, patch
`+$3A8F` to `$01`, and use the pending-trigger from
`mesen/probes/EXP-0047.lua`.

It blocks nothing. If the operator would rather bank the ATB program and
return to the scenario, the alternative next action is the **Whelk
decision** — reinterpret EXP-0040's scoped captures, or re-run the fight
with the model in hand — which is now an ordinary orchestration call
rather than a blocked one.

## Recommended next command

`/orchestrate`
