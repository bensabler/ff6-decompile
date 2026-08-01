# Checkpoint 2026-08-01 — EXP-0042: battle-entry configuration sampling (Unit 44)

## Current question

None open in the unit. EXP-0042 is **completed**: the sampling question
is answered and both derived values were predicted before being observed.

## State

The ATB program's **staging rule is established**. `Bat.Mode` and
`Bat.Speed` are sampled **once at battle entry** into battle-local cells,
so every ACTIVE/WAIT and Battle Speed condition must be set before entry
— or injected directly at those cells, which is a far better controlled
handle than driving menus.

The hard ATB blocker **remains open**: no timer domain, pause condition,
or queue semantics is known. What changed is that there is now a concrete
way in, `WRAM:+$3A90`, whose consumer is unlocated.

Whelk was not resumed and its savestates were not reloaded.

## Confirmed before this session

Golden route power-on → milestone `06-random-encounter`; unified battle
arrays (DISC-0001); configuration storage and bit layout (EXP-0041,
CEN-MENU-0007).

## Work completed

**Answer: mixed, split by setting.** At battle entry, `ROMCPU:$C22472`
reads `+$1D4D` and `+$1D4E` once each and **decomposes** them:

| Setting | Destination | Value |
|---|---|---|
| Bat.Mode (`+$1D4D` bit 3) | `WRAM:+$3A8F` | `INC` iff Wait — `01` = Wait, `00` = Active |
| Bat.Speed (`+$1D4D` bits 0-2) | `WRAM:+$3A90` | `255 − 24 × speed`; Fast `$FF`, default `$CF`, Slow `$87` |
| Cmd.Set (`+$1D4D` bit 7) | `WRAM:+$2F2E` | cleared when Window |
| Gauge (`+$1D4E` bit 7) | `WRAM:+$2021` | cleared when Off |
| `+$1D4E` bits 0-2 | `WRAM:+$2F34` | at `$C10FF7` (bits no Config setting touches) |

Neither Bat.Mode nor Bat.Speed is re-read for timing during the battle.
`Msg.Speed` and `Cursor` **are** read live from the persistent bytes:
`$C198AC` extracts Msg.Speed and indexes a delay table at
`ROMCPU:$C19872`; `$C159D6` tests Cursor and clears the `$5C`-byte
cursor-memory block at `WRAM:+$890F` when the setting is Reset — the
mechanism behind EXP-0040's `Cursor = Memory` observation.

Method note: the `$C22472` arithmetic was decoded statically from a ROM
dump, then used to **predict** `$3A8F` and `$3A90` for a second,
differently-configured run. Both matched exactly (`00`/`$87` for Active +
Bat.Speed 6). That is what makes the transform Confirmed rather than
inferred.

Instrumentation added: `watchreads` and `watchdump` in
`mesen/probes/common.lua` (the shared helper had a write watch but no
read counterpart, and no dump path at all), plus
`mesen/probes/EXP-0042.lua`.

New **CEN-BATTLE-0010**. Memory map extended with `+$3A8F`, `+$3A90`,
`+$2F2E`/`+$2021`/`+$2F34`, `+$890F`, and the three decoded ROM regions.

## Last raw observation

Run 2, frame 58 058, battle entry at 55 922, `+$1D4D` = `$25`
(Active, Bat.Speed 6): `$C22475` and `$C22493` each **count=1** at the
entry frame; `WRAM:+$3A8F`/`+$3A90` read `00 87`, matching the prediction
`255 − 24 × 5 = 135 = $87`.

## Active emulator state

**None.** One headless `--testrunner` instance; logs harvested before
termination, killed with `kill -9`, absence confirmed by `pgrep`.
`jobs -l` empty. **SRAM directory verified still empty** — virgin boot
preserved.

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0042.lua` armed a read watch on
`WRAM:+$1D4D`-`+$1D4E` and a battle-entry write detector on
`+$3B18`-`+$3BB7`; both die with the process.

## Evidence paths

`local_artifacts/experiments/EXP-0042/` — **8 artifacts, 304 KB**, with
`hashes.sha256`: 2 screenshots, 2 savestates (`in-battle-formation14.mss`
and `in-battle-active-speed6.mss`, both **live battle states** worth
reusing), `run2-read-table.txt`, `bridge-commands.log`,
`bridge-events.log`, `experiment.json`.

Starting state: `local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss`.
Both runs reached formation 14, reproducing EXP-0038's mines encounter.

## Files changed

- New: `docs/experiments/EXP-0042-battle-entry-config-sampling.md`,
  `mesen/probes/EXP-0042.lua`, this checkpoint.
- Updated: `mesen/probes/common.lua` (+`watchreads`, +`watchdump`),
  `docs/sessions/04_MEMORY_MAP.md`, `manifests/experiments.json`,
  `manifests/content-census.json`, `dashboards/` (CURRENT_FOCUS,
  BLOCKERS, RESEARCH_QUEUE, ACTIVITY_LOG), `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/EXPERIMENTS.md`, `indexes/CONTENT_CENSUS.md`,
  `indexes/ROM_REGIONS.md`, `dashboards/COVERAGE.md`.

## Tests and quality gates

`gofmt -l .` clean; `go build ./...`, `go vet ./...`, `go test ./...` ok;
`ff6lab audit` clean; `ff6lab census sync` clean (63 entries);
restricted-extension scan clean. No Go source changed this unit.

## Git status

`main`, one coherent unit committed, worktree clean.

## Unresolved decisions

- **The consumer of `WRAM:+$3A90` is unlocated.** It is the sharpest
  current lead into ATB rate, and finding it is the next unit.
- The consumer of `+$3A8F`, and of `+$2F2E`/`+$2021`/`+$2F34`.
- Whether `$3A8F` can be incremented more than once — only `00` and `01`
  were observed.
- Whether battle types other than random encounters sample configuration
  differently. Both runs were formation-14 random encounters.
- Owner of `+$1D4E` bits 0-3 and `+$1D54` bits 0-6 (carried from
  EXP-0041); `$C10FF7` now gives bits 0-2 a known consumer.

## Blockers

**HARD BLOCKER — no ATB model. Still open.** Timer domains, pause
conditions and action-queue ordering remain Unknown. EXP-0041 and
EXP-0042 established the configuration layer and the staging rule; no
gauge, tick, or queue has been observed.

**Whelk gameplay must not resume before that research.** The evidence
constraint is unchanged: EXP-0040's head/shell timing is
menu-pause-contaminated, and ACTIVE and WAIT are separate experimental
conditions.

This unit does **not** support, and no downstream record may assert:
that `$3A90` is an ATB rate, increment, or threshold; that larger
`$3A90` means faster anything; that `$3A8F` gates any particular clock;
or that configuration sampling behaves the same way in scripted, pincer,
back or boss battles.

## Exact next action

**EXP-0043 — locate the consumer of `WRAM:+$3A90`, and the ATB gauges.**
Read-watch `+$3A8F`–`+$3A90` across a battle and follow the reading
routine to whatever it advances. Expect convergence with the eight
undumped callees of `ROMCPU:$C101FB` (open question #6); this should also
resolve open question #18, the sub-1.0 per-frame dispatch rate. Start
from `local_artifacts/experiments/EXP-0042/in-battle-formation14.mss`, a
live battle preserved by this unit, so no route replay is needed.

## Recommended next command

`/orchestrate`
