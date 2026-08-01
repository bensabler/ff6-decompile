# Checkpoint 2026-08-01 — EXP-0045: queued work resolves past the gate (Unit 47)

## Current question

None open. EXP-0045 is **completed**; both parts answered, part A with a
prediction test.

## State

The ATB pause model is now precise. The scheduler stops dead under the
gate, but **work already pending when the gate engages still completes**.
That was EXP-0044's one unresolved signal, and it is now resolved.

## Confirmed before this session

Config storage (EXP-0041); entry sampling (EXP-0042); the ATB layer and
gate instruction (EXP-0043); the pause matrix (EXP-0044).

## Work completed

**Part A — the transient.** Per-frame `endFrame` tracing, buffered in
memory and flushed on demand so file I/O could not distort the timing.

| Trace | Pending at gate engage | Gated-frame changes |
|---|---|---|
| 1 | slot 6 (`+$3AA0` bit 6 set) | 1, at +78 frames |
| 2 | **none** | **0** across 438 gated frames |
| 3 | slot 8 — **predicted** | 1, at +119 frames |

In traces 1 and 3 the change was identical in shape: that slot's
`+$3AA0` bit 6 cleared and its `+$3218` entry advanced by exactly
`$0100`, while the tick counter and every gauge stayed frozen. Trace 1
had 77 clean gated frames before it and 349 after, so it is neither a
boundary artifact nor an ongoing clock.

Trace 3 was a genuine prediction test. Traces 1 and 2 differ in exactly
one respect — whether any slot was pending — which predicts that
engaging the gate while a slot is pending will produce a completion. It
did, on a **different slot** at a **different delay**. Arming had to move
into Lua: bridge round-trips are hundreds of frames apart, which is
precisely why EXP-0044 could not settle this.

**`+$3AA0` bit 6 is the pending-action marker** — a semantics question
EXP-0043 and EXP-0044 both left open.

**Part B — matrix rows.** Item is a qualifying submenu (`$2F41` = `01`,
370 frames fully frozen). The victory presentation is **not** (`$2F41` =
`00` across 429 in-battle frames, tick advancing on 195). Magic, Row and
Defend are **unreachable** from a Magitek battle and are left unsampled
rather than guessed.

**Consequence for EXP-0040.** It reported that "queued actions resolved
out of issue order" while menus were open and treated that as operator
failure against an unmodelled system. It was real system behaviour. That
observation is vindicated and now has a mechanism.

## Last raw observation

Trace 3, f=82651, gate set: `fl=0101,0101,0141` → `0101,0101,0101` and
`ac=...,0083` → `...,0183`, with `tick=0D9A` and gauges unchanged.

## Active emulator state

**None.** One headless instance; logs harvested before termination,
killed with `kill -9`, absence confirmed. `jobs -l` empty. **SRAM
verified still empty.**

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0045.lua` armed a per-frame
`endFrame` trace and a `+$2F41` write watch; an additional
pending-triggered arm was installed via `eval` for trace 3. All die with
the process.

## Evidence paths

`local_artifacts/experiments/EXP-0045/` — **12 artifacts, 304 KB**, with
`hashes.sha256`: four per-frame trace logs
(`exp45-transition{,2,3}.log`, `exp45-item.log`, `exp45-presentation.log`),
4 screenshots, `bridge-commands.log`, `bridge-events.log`,
`experiment.json`.

## Files changed

- New: `docs/experiments/EXP-0045-gate-transition-and-matrix-completion.md`,
  `mesen/probes/EXP-0045.lua`, this checkpoint.
- Updated: `docs/sessions/04_MEMORY_MAP.md` (`+$3AA0` bit 6, `+$3218`),
  `manifests/experiments.json`, `manifests/content-census.json`
  (CEN-BATTLE-0012 extended), `dashboards/` (CURRENT_FOCUS,
  RESEARCH_QUEUE, ACTIVITY_LOG), `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/`, `dashboards/COVERAGE.md`.

## Tests and quality gates

`gofmt -l .` clean; build, vet, `go test ./...` ok; `ff6lab audit` clean;
`ff6lab census sync` clean; restricted-extension scan clean. No Go source
changed.

## Git status

`main`, one coherent unit committed, worktree clean.

## Unresolved decisions

- **What drives a pending action to completion while the gate is shut.**
  The next unit.
- Whether several pending slots all drain during one gated interval —
  never observed; no trace had two pending at once.
- Whether `+$3218` always advances by exactly `$0100`; its low byte was
  unchanged both times.
- Magic / Row / Defend submenus — need a normal-party battle.
- The defeat presentation, with gate instrumentation.
- Party slots pending at a gate engage; only enemy slots were seen.

## Blockers

None. The ATB blocker was discharged at EXP-0044 and nothing here
reinstates it.

This unit does **not** support, and no downstream record may assert:
that the completion delay is a fixed value (78 and 119 frames observed);
that `+$3218` always steps by `$0100`; that multiple pending slots behave
like one; that party slots behave like enemy slots; or that any
un-instrumented presentation state pauses or does not pause.

## Exact next action

**EXP-0046 — the action-queue execution path.** Read-watch `+$3AA0` and
`+$3218` through a gated interval arranged as trace 3 arranged it, and
capture the writing PC. That routine completes a pending action while
`ROMCPU:$C21124` is shut, so it sits outside the per-frame scheduler —
the queue/execution half of the ATB model.

## Recommended next command

`/orchestrate`
