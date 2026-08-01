# Latest Checkpoint

**[2026-08-01 — EXP-0047: the execution path is periodic and ungated](2026-08-01-exp0047-execution-path-invocation.md)**

State: the delay question is answered. `$C201BE` — the completion write —
executed at **`gate=0` and `gate=1`**, so the action-execution path was
never behind the ACTIVE/WAIT gate. It is **periodic**, firing roughly
every 100-120 frames (observed gaps 121/35/105/122) and sweeping the
battle slots. The **78/119/122-frame completion delays are the wait for
the next invocation**, scheduled upstream of `$C21124`.

This reframes EXP-0045's finding: the gate stops the scheduler, and the
execution path was never behind it to begin with.

**Two refutations, recorded rather than buried.** `$C20EB6` is not a call
frame (`TSX` at return−3). `$C20016` holds a plausible `JSR` yet **never
executed** across a gated interval in which a completion demonstrably
occurred. Method note to carry forward: **a `JSR` at return−3 is
necessary but not sufficient — confirm stack frames by execution.**

Confirmed on the live path: `$C2141D`, `$C208AE`, `$C201A0`, `$C201BE`.
Named but undecoded call targets: `$C223ED`, `$C2083F`, `$C213D3`.
The invoker itself remains **Unknown**.

No blockers. All gates clean; no background processes; SRAM virgin.

Exact next action: **EXP-0048 — name the invoker with a different
instrument.** Stack archaeology has failed twice here; exec-watch outward
from the confirmed `$C2141D`, or trace a single invocation. Narrow
question, confirmed sites to walk from.
