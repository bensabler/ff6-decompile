# Checkpoint 2026-08-01 — EXP-0046: the action-execution path (Unit 48)

## Current question

None open. EXP-0046 is **completed**; the writer EXP-0045 left unnamed is
identified, and an apparent contradiction was resolved by evidence.

## State

The ATB model's execution half is located. `ROMCPU:$C201BE`
(`INC $3219,X`) advances `+$3218` by `$0100`, guarded by `+$3AA0` bit 3,
reached from an action-execution path in `ROMCPU:$C207xx`–`$C209xx` that
runs outside `$C21124`'s ACTIVE/WAIT gate.

## Confirmed before this session

Config storage (EXP-0041); entry sampling (EXP-0042); the ATB layer
(EXP-0043); the pause matrix (EXP-0044); queued work resolving past the
gate (EXP-0045).

## Work completed

The pending-trigger from EXP-0045 arranged a gated interval with slot 7
pending. **All thirteen gated writes fell on a single frame** (60169,
122 frames after the trigger) — the deferred completion is a burst, not a
drain, consistent with EXP-0045's +78 and +119.

Gated writers: `$C20798`, `$C211BA`, `$C20974` (flags) and `$C201C1`
(accumulator).

**`+$3218`'s `$0100` is named.** `$C201C1`'s write went to `$7E3227`, the
high byte of slot 7's entry, with value 1. Decoding:

```asm
C201B7  BD A0 3A     LDA $3AA0,X
C201BA  89 08        BIT #$08        ; flag bit 3
C201BC  D0 08        BNE $C201C6
C201BE  FE 19 32     INC $3219,X     ; the +$0100
```

**An apparent contradiction, resolved.** `$C211BA` sits inside the
gauge-advance routine the gate skips, so a gated write there would
contradict EXP-0045. The stacks — identical across every capture, with
innermost return `$C208CB` — showed why:

```asm
C208C6  A9 50        LDA #$50
C208C8  20 B4 11     JSR $C211B4
C208CB  AD 04 34     LDA $3404
```

`$C211B4` is a **shared helper** (`ORA $3AA0,X / STA / RTS`) with at
least two entry points: the scheduler enters at `$C211B2` with `A=$20`,
`$C208C6` enters at `$C211B4` with `A=$50`. The same store PC therefore
appears in gated and un-gated traffic without the scheduler running.

**Correction propagated.** EXP-0043 attributed that `+$3AA0` store to the
scheduler's threshold path. It belongs to a shared helper, so a PC in
that range does not imply the scheduler ran. The memory map is amended.
No earlier conclusion depended on the stronger reading.

Also decoded: `$C20795` clears `+$3AA0` bit 7 and sets bit 6 of a new,
uncharacterised per-slot byte `+$3204,X`; `$C20974` sets `+$3AA0` bit 3 —
the bit `$C201BA` tests.

New **CEN-BATTLE-0013**.

## Last raw observation

`EXP46-ACCUM GATED-WRITE a=7E3227 v=1 pc=$C201C1 A=$0001 X=$000E
Y=$0000 SP=$15F3 PS=$35 DB=$7E frame=60169`

## Active emulator state

**None.** One headless instance; logs harvested before termination,
killed with `kill -9`, absence confirmed. `jobs -l` empty. **SRAM
verified still empty.**

## Breakpoints/watchers

None left installed. `mesen/probes/EXP-0046.lua` armed gate-tagged write
watches over `+$3AA0`/`+$3218` with `probelog` stack capture, plus the
pending-trigger; all die with the process.

## Evidence paths

`local_artifacts/experiments/EXP-0046/` — 5 artifacts with
`hashes.sha256`: `gated-flags.txt`, `gated-accum.txt`,
`bridge-commands.log` (with the ROM dumps in the transcript),
`bridge-events.log` (the stack captures), `experiment.json`.

## Files changed

- New: `docs/experiments/EXP-0046-action-queue-execution-path.md`,
  `mesen/probes/EXP-0046.lua`, this checkpoint.
- Updated: `docs/sessions/04_MEMORY_MAP.md` (correction to the
  `$C21193`–`$C211BA` entry; new entries for `$C201B1`, the `$C207xx`
  path, `+$3204`; `+$3AA0` bit picture), `manifests/experiments.json`,
  `manifests/content-census.json`, `docs/checkpoints/LATEST.md`.
- Regenerated: `indexes/`, `dashboards/COVERAGE.md`.

## Tests and quality gates

`gofmt -l .` clean; build, vet, `go test ./...` ok; `ff6lab audit` clean;
`ff6lab census sync` clean; restricted-extension scan clean. No Go source
changed.

## Git status

`main`, one coherent unit committed, worktree clean.

## Unresolved decisions

- **What invokes the execution path, and why the completion fired 122
  frames after the gate engaged.** The next unit.
- The caller chain above `$C208C6` (`$C208B1`, `$C21420`, `$C20EB6`) was
  captured but not decoded.
- Boundaries and purpose of the routines enclosing `$C20798`/`$C20974`.
- `+$3AA0` bits 0-2, 4, 5; and `+$3204`, new here.
- Whether several pending slots drain in one burst or several.
- Party slots at a gate engage; only enemy slots have been observed.

## Blockers

None.

This unit does **not** support, and no downstream record may assert:
that the 122-frame delay is characteristic (three delays observed: 78,
119, 122); that the burst is always one frame wide beyond the single
fully instrumented event; that `$C20798`/`$C20974` do what their few
decoded instructions suggest at routine scope; or that any PC in the
`$C211B4` region implies scheduler execution.

## Exact next action

**EXP-0047 — what invokes the action-execution path, and when.** Dump
the captured stack frames above `$C208C6` and exec-watch the entry point
across a gated interval to recover the invocation cadence. That is the
last structural piece before the ATB programme's action-lifecycle and
queue-model deliverables can be written down.

## Recommended next command

`/orchestrate`
