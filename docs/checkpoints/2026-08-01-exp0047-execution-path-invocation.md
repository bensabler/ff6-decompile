# Checkpoint 2026-08-01 — EXP-0047: the execution path is periodic and ungated (Unit 49)

## Current question

None open. EXP-0047 is **completed**: the delay question is answered, and
two reconstructions from EXP-0046 were refuted in the process.

## State

The action-execution path is **periodic and ungated**. The 78/119/122-frame
completion delays are the wait for its next invocation, not a countdown
and not a gate effect. Its invoker remains unnamed.

## Confirmed before this session

Config storage (EXP-0041); entry sampling (EXP-0042); the ATB layer
(EXP-0043); the pause matrix (EXP-0044); queued work resolving past the
gate (EXP-0045); the completion write `ROMCPU:$C201BE` (EXP-0046).

## Work completed

**The cadence answer.** Exec watches on the four verified sites, gate
state recorded at each hit:

```
EXEC C201A0 gate=0  frame=75825
EXEC C2141D gate=0  frame=75946
EXEC C208AE gate=0  frame=75946   (x3, X = $0010/$000E/$000C)
EXEC C201BE gate=0  frame=75946
EXEC C2141D gate=0  frame=75981
EXEC C201A0 gate=0  frame=76086
EXEC C2141D gate=1  frame=76208
EXEC C201BE gate=1  frame=76208
```

`$C201BE` executed at **both** gate states — the path is not gated at
all. Firing frames give gaps of **121, 35, 105, 122** emulator frames,
bracketing the completion delays measured earlier (78, 119, 122).
`$C208AE`'s three hits on one frame carry the stride-2 indices for slots
8, 7 and 6, so the path sweeps the slots per invocation.

This **reframes EXP-0045**: the gate stops the scheduler, and the
execution path was never behind the gate to begin with. Nothing in
EXP-0045's measurements changes — only the explanation improves.

**Two refutations.**

- `$C20EB6` is **not a call frame** — `return − 3` holds `BA` (`TSX`).
  Falsifier 3 fired; the chain truncates there.
- `$C20016` is **not the entry point** despite holding
  `20 ED 23` (`JSR $C223ED`) at `return − 3`. An exec watch recorded
  **zero** executions across a gated interval in which a completion
  demonstrably occurred. Falsifier 1 fired.

**Method note, recorded for reuse:** a `JSR` at `return − 3` is necessary
but **not sufficient** to confirm a stack frame. Raw stack windows carry
routine-pushed data and `JSL` frames are three bytes, so pair-reading
misleads. Confirm frames by **execution**.

## Last raw observation

`EXEC C201BE gate=1 pc=$C201BE A=$0001 X=$000E Y=$0000 SP=$15F3 PS=$37
DB=$7E frame=76208`

## Active emulator state

**None.** One headless instance; logs harvested before termination,
killed with `kill -9`, absence confirmed. `jobs -l` empty. **SRAM
verified still empty.**

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0047.lua` armed `execwatch` with
gate tagging plus the pending-trigger; four additional exec callbacks
were installed via `eval`. All die with the process.

## Evidence paths

`local_artifacts/experiments/EXP-0047/` — 4 artifacts with
`hashes.sha256`: `bridge-commands.log` (ROM dumps in the transcript),
`bridge-events.log` (the tagged exec captures), `exp47-entry.log`,
`experiment.json`.

## Files changed

- New: `docs/experiments/EXP-0047-execution-path-invocation.md`,
  `mesen/probes/EXP-0047.lua`, this checkpoint.
- Updated: `manifests/experiments.json`, `manifests/content-census.json`
  (CEN-BATTLE-0013 extended with the cadence finding and the method
  note), `dashboards/RESEARCH_QUEUE.md`, `dashboards/ACTIVITY_LOG.md`,
  `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/`, `dashboards/COVERAGE.md`.

`mesen/probes/common.lua` gains nothing this unit; `execwatch` lives in
the EXP-0047 probe because its gate-tagging is unit-specific.

## Tests and quality gates

`gofmt -l .` clean; build, vet, `go test ./...` ok; `ff6lab audit` clean;
`ff6lab census sync` clean; restricted-extension scan clean. No Go source
changed.

The audit caught a dangling link when `LATEST.md` was written before this
checkpoint existed — working as intended, fixed by writing the file.

## Git status

`main`, one coherent unit committed, worktree clean.

## Unresolved decisions

- **The invoker of the execution path.** Two candidates refuted; needs a
  different instrument.
- The invocation period on more than one run — gaps of 121/35/105/122
  from a single run, one markedly shorter.
- `$C223ED`, `$C2083F`, `$C213D3` — named call targets, none decoded.
- `+$3AA0` bits 0-2, 4, 5; `+$3204`.
- Whether several pending slots drain in one invocation.

## Blockers

None.

This unit does **not** support, and no downstream record may assert:
that the invocation period is a constant; that `$C208AE` sweeps exactly
ten slots (three X values seen); that any un-executed address is on the
path merely because a `JSR` sits at `return − 3`; or that these findings
hold outside a random encounter.

## Exact next action

**EXP-0048 — name the invoker, with the right instrument.** Stack
archaeology has failed twice on this path. Exec-watch outward from the
confirmed `$C2141D`, or trace a single invocation with Mesen's trace
facility. The question is narrow: one routine, firing about every 100
frames, with confirmed sites to walk from.

## Recommended next command

`/orchestrate`
