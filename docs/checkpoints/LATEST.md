# Latest Checkpoint

**[2026-08-01 — CORR-0001: pointer advancement at ROMCPU:$C09B5C](2026-08-01-corr0001-c09b5c-pointer-advance.md)**
(preceding session close: [the ATB research program](2026-08-01-session-atb-program-complete.md))

State: the Ghidra/Mesen static-analysis workflow is integrated and its
**pilot correlation is complete**. `ROMCPU:$C09B5C` is now understood at the
mechanism level, and the event interpreter is no longer a purely static
lead.

Confirmed at `$C09B5C`, invariant across 24 observations (12 hits × 2 runs,
byte-identical logs): entry state E = native, **M = 1**, **X = 0**,
PBR = `$C0`, DBR = `$00`, **DPR = `$0000`**, SP = `$15FD`. Because DPR was
measured rather than assumed, the direct-page operands resolve to
`WRAM:+$00E3` and `WRAM:+$00E5`-`+$00E7` — the first time this project has
measured D at a site instead of inheriting `D = 0`.

The 24-bit little-endian value at `WRAM:+$00E5`-`+$00E7` increases by
**exactly `A & $FF`** on 24/24, with deltas varying 1-7; a constant-delta
reading is refuted. Control reaches `ROMCPU:$C09A6D` in the same frame every
time, where `WRAM:+$00E3` is decremented once per frame until zero (observed
twice at exactly 30 frames). `$C09B5C` is a **genuine shared routine entry** —
preceded by an unconditional `JMP ($002A)`, and entered from multiple command
handlers, not just the flag family, which refutes the narrower reading in
ROM-0027's note.

Three interpretations remain **Strong hypothesis, not Confirmed**, sharing
one ambiguity — the value was never observed being dereferenced: that it is
an event-script pointer, that A is a command length, and that `+$00E3` is the
script frame-wait. The **immediate predecessor is unresolved**.

This unit had **no `EXP-NNNN` record**. Under explicit operator scope it was
pre-registered inside `docs/correlations/CORR-0001-C09B5C.md` — question,
starting state, breakpoints, watches and falsifier fixed before the runs. No
retroactive experiment record was created.

Evidence: `local_artifacts/static-analysis/CORR-0001/` (gitignored), frozen
in `hashes.sha256`. Nothing in flight: no emulator, no resident
instrumentation, worktree clean. All gates green (gofmt/build/vet/test,
`ff6lab audit`, `ff6lab census validate`, restricted-file scan).

Exact next action: **name the immediate predecessor** — exec-watch the
dispatcher's `JMP ($002A)` at `ROMCPU:$C09B59` alongside `$C09B5C`, capturing
DP `$2A`/`$2B` and `$EA` at dispatch. It closes the pilot's one gap, turns
"A is a command length" into a measured opcode→length table, and decodes the
candidate opcode table at `ROMCPU:$C098C4` as a by-product. It is also the
**deferred blocker** for any demonstration that claims to read or execute
event script.
