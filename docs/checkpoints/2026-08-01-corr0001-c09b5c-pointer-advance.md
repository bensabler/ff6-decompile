# Checkpoint 2026-08-01 — CORR-0001: pointer advancement at ROMCPU:$C09B5C

First static/runtime correlation unit. Closes the Ghidra/Mesen workflow's
pilot; opens no new investigation.

## Current question

What does `ROMCPU:$C09B5C` actually do, in what 65816 state is it entered,
and what is its real relationship with `ROMCPU:$C09A6D`? Answered for the
mechanism; the semantic layer remains hypothesis.

## State

CORR-0001 is complete and atomic. Record:
`docs/correlations/CORR-0001-C09B5C.md`. The unit was **bounded to
`$C09B5C`** and did not follow the dispatcher, the opcode table, or the
script data it exposed.

## Confirmed before this session

- ROM-0027 / DISC-0008 / EXP-0037: the 16 event-flag command handlers at
  `ROMCPU:$C0B593-$C0B6D2`, each ending `LDA #$02 / JMP $9B5C`.
- CEN-EVENT-0001 registered `$C09B5C` as a **candidate** interpreter-advance.
- `/bootstrap-ghidra`: Ghidra 12.1.2, SNES Loader 1.3.0, language
  `65816:LE:24:snes`, `SNES ROM Loader` import, 48 HiROM bank blocks, ROM
  SHA-256 matching `ROM_IDENTITY.md` on both the file and inside the Ghidra
  database.

## Work completed

### Confirmed processor state at `ROMCPU:$C09B5C`

Read from the live CPU at the entry PC on **all 24 observations** (12 hits ×
2 runs), invariant across every one:

| Field | Value |
|---|---|
| E (emulation) | `false` — native mode |
| M (accumulator width) | **1 — 8-bit A** |
| X (index width) | **0 — 16-bit X/Y** |
| PBR | `$C0` |
| DBR | `$00` |
| DPR | **`$0000`** |
| SP | `$15FD` |

No processor-state assumption remains at the entry PC. P itself varied
(`$20/$21/$60/$61`); the M and X bits did not.

Because DPR was **measured** at `$0000`, the direct-page operands resolve to
`WRAM:+$00E3` and `WRAM:+$00E5`-`+$00E7`. Every prior direct-page read in
this project hardcoded D = 0; that assumption is now measured at this site
rather than inherited.

### Measured 24-bit increment at `WRAM:+$00E5`-`+$00E7`

The 24-bit little-endian value increases by **exactly `A & $FF`** on 24 of 24
observations. Deltas varied across **1, 2, 3, 4, 5, 7**, which refutes any
constant-delta or A-independent reading outright — that was the
pre-registered falsifier and it was not met.

Observed sequence (identical in both runs): `$CA5E92`→`$CA5E93`, then
`$CC9A55` through `$CC9A7B` in steps of 1, 3, 2, 4, 1, 4, 5, 2, 7, 7, 2.

### Measured wait-counter behaviour at `WRAM:+$00E3`

At `ROMCPU:$C09A6D`: when `+$00E3` was `$00`, the bypass path `$C09A75` ran.
When it was `$1E` (30), the wait path `$C09A71` executed **once per frame for
exactly 30 consecutive frames**, then the bypass ran on the next frame.
Observed twice — frames 18950-18979 and 18981-19010 — with the counter
reaching zero both times.

### Confirmed transfer to `ROMCPU:$C09A6D`

Control reached `$C09A6D` in the **same frame**, on 24 of 24 observations.
No exceptions.

### `$C09B5C` is a genuine shared routine entry

- **Entry, not an internal block:** the instruction immediately before it is
  `JMP ($002A)` at `$C09B59` — unconditional, so `$C09B5C` cannot be reached
  by fallthrough. Every observed entry was by explicit transfer.
- **Shared, not the flag handlers' tail:** 9 of 12 hits had no flag-family
  handler executing, and the observed A values of 1, 3, 4, 5 and 7 are
  unreachable from handlers that always load `#$02`. This **refutes the
  narrower reading** implied by ROM-0027's note.

### Two independently reproduced runs

`01-opening-cinematic/run1-01.mss` and `run2-01.mss` — the two independently
produced states from the EXP-0032 determinism pair, frame 15,001, no inputs
injected. The two logs differ on **exactly two lines**: the savestate byte
size (148796 vs 148697) and its path. All 12 hit records, 12 continuation
records, 60 wait-path lines and 37 bypass lines are byte-identical. This
meets the log-level reproducibility standard EXP-0037 set.

A third run is preserved as a **null result**: frames 15,001-16,801 produced
zero hits, so the routine does not run continuously — it runs in bursts while
the script advances.

## Last raw observation

`HIT #12 frame=19032 pcReported=$C09B5C targetMatch=true E=false M=1 X=0
(P=$20) PBR=$C0 DBR=$00 D=$0000 SP=$15FD A=$0002` → `ptr24Before=$CC9A79`,
`ptr24After=$CC9A7B`, `delta=2`, `deltaMatchesAlow=true`, continuation at
`$C09A6D` in the same frame.

## Active emulator state

**None.** Both runs were headless `--testrunner`, self-terminated via
`emu.stop(0)` at the frame window's end. `pgrep` clean. No GUI instance, no
resident script.

## Breakpoints/watchers

**None resident.** All watches lived in `mesen/probes/CORR-0001.lua`
(tracked, not loaded now): exec callbacks at `$C09B5C`, `$C09A6D`,
`$C09A71`, `$C09A75`, and over the range `$C0B593-$C0B6D2`.

## Evidence paths

`local_artifacts/static-analysis/CORR-0001/` (gitignored), frozen in
`hashes.sha256`:

| File | SHA-256 |
|---|---|
| `corr0001-obs0-shortwindow.log` | `8e1bc842881d5b9455891a9deed7643abf41d918f84ca3b88ee3a7c0c630bc0f` |
| `corr0001-obs1.log` | `6ded824c763ca9ecf9d4ed23ecd80bb2c2a9d5186d361c7478c047b33cdc6aca` |
| `corr0001-obs2.log` | `16529f51a49d87cf673feed82062e7491619538247b6d5c5e835fb5674d72151` |

Emulator binary used: `~/Desktop/Mesen.app`, SHA-256
`a49954ff7889e146ee1b3173ac4a0f1293ec1b741362a8f6f692937201b5f258`. Per the
bootstrap record that is the **2.2.1** build. Historical version attribution
in other records was deliberately **not** touched and remains Unknown.

## Files changed

| File | Change |
|---|---|
| `docs/correlations/CORR-0001-C09B5C.md` | created — the single correlation record |
| `mesen/probes/CORR-0001.lua` | created — instrumentation, tracked per the preserve-injected-code rule |
| `manifests/content-census.json` | CEN-EVENT-0001 only: `cpu_routines`, `next_action`, `notes` |
| `dashboards/COVERAGE.md` | regenerated by `ff6lab census sync` (one line, derived from the above) |
| `docs/checkpoints/2026-08-01-corr0001-c09b5c-pointer-advance.md` | this file |
| `docs/checkpoints/LATEST.md` | repointed |
| `dashboards/CURRENT_FOCUS.md` | next-action section |

No Go code changed. No index, experiment, session, discovery, or contradiction
record touched.

## Tests and quality gates

All green immediately before commit: `gofmt -l .` silent; `go build ./...`;
`go vet ./...`; `go test ./...` (15 packages ok); `ff6lab audit` clean;
`ff6lab census validate` clean after `census sync`; restricted-extension scan
of `git ls-files` clean.

`census validate` failed once mid-unit with `COVERAGE.md is stale` — the
staleness gate doing its job; `census sync` regenerated it.

## Git status

`main`, worktree clean after this commit, **ahead 2** of `origin/main`
(`b91f159` integration + this unit). Nothing pushed.

## Unresolved decisions

- **The immediate predecessor is unresolved.** The probe recorded the last
  execution inside the flag-handler range, which is *not* the instruction
  executed immediately before entry. The evidence supports **multiple
  distinct predecessors** but names none for 9 of the 12 hits. This is the
  pilot's one real gap.
- **Three Strong but unconfirmed interpretations**, all sharing one
  ambiguity — the value at `WRAM:+$00E5`-`+$00E7` was never observed being
  *dereferenced*:
  - that it is an **event-script pointer into ROM** (support: values sit in
    ROM banks `$CA`/`$CC`; the observed steps chain perfectly contiguously;
    splitting the ROM at those boundaries yields consistent opcode→length
    pairing, e.g. `$41` twice at length 2, `$0E`/`$0F` as matched 7-byte
    records);
  - that **A is the command length** of the just-executed command;
  - that `WRAM:+$00E3` is the **event-script frame-wait counter** (the
    countdown mechanism is Confirmed; the semantic label is not).
- Whether the invariant SP/DBR values are properties of the routine or of
  this starting state — both observations descend from the same milestone.
- **No `EXP-NNNN` record exists for this unit.** `CLAUDE.md` requires a
  falsifiable experiment record before operating Mesen. This unit was scoped
  by explicit operator instruction to exactly one `CORR` record, so the
  question, starting state, breakpoints, watches and **falsifier were
  pre-registered inside `CORR-0001-C09B5C.md`'s Runtime experiment section
  before the runs**, not fitted afterwards. No retroactive experiment record
  was created and none existed beforehand. Registering one for the runtime
  half remains available and would require a `manifests/experiments.json`
  entry plus index regeneration.

## Blockers

None for the project. One **deferred demo blocker**: the dispatcher /
predecessor investigation below gates any demonstration that claims to read
or execute event script, because until it lands the "script pointer",
"command length" and "frame wait" readings are Strong hypotheses, not
Confirmed — and a demo must not present them as settled.

Registered but **not investigated** (deferred neighbouring leads, in
CEN-EVENT-0001): the `$C09B3B-$C09B5B` dispatcher; the `ROMCPU:$C098C4`
candidate opcode table (first 64 entries all valid bank-`$C0` pointers
spanning `$C09C44-$C0A336`); the `$C09B1E` six-byte prefetch into DP
`$EA-$EF`; the `$C09B82` variant maintaining a 3-byte-stride structure at
`$05F4` indexed by `$E8`; and candidate script data in ROM banks `$CA`/`$CC`
around `$CC9A55`.

## Exact next action

**Name the immediate predecessor.** Arm an exec watch on the dispatcher's
`JMP ($002A)` at `ROMCPU:$C09B59` alongside `$C09B5C`, capturing DP
`$2A`/`$2B` and `$EA` at dispatch, and correlate each `$C09B5C` entry with
the opcode that reached it.

That single unit closes this pilot's one gap, converts "A is a command
length" from Strong hypothesis into a measured opcode→length table, and
decodes `ROMCPU:$C098C4` as a by-product.

## Recommended next command

```text
/correlate-static-runtime ROMCPU:$C09B59
```
