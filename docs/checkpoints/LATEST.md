# Latest Checkpoint

**[2026-08-01 — session close: the ATB research program](2026-08-01-session-atb-program-complete.md)**
(per-unit detail: [EXP-0047](2026-08-01-exp0047-execution-path-invocation.md)
and its six predecessors)

State: the **ATB research program is complete for its purpose** and the
blocker EXP-0040 raised is **discharged**. Ten commits — three
infrastructure, seven experiments — established configuration storage and
encoding, battle-entry sampling, the ATB gauges and increments, the tick
counter, the exact ACTIVE/WAIT pause condition, and the action-execution
path.

The model is summarised in `dashboards/CURRENT_FOCUS.md`; the evidence
lives in the EXP-0041..0047 records and CEN-MENU-0007 /
CEN-BATTLE-0010..0013.

Two results are worth singling out because they were **predicted before
being observed**: EXP-0042 derived `+$3A8F`/`+$3A90` from a static decode
and matched both exactly at a second configuration, and EXP-0045
predicted that engaging the gate with a slot pending would produce a
deferred completion — which fired on a different slot at a different
delay.

Two corrections were propagated rather than buried: `$C211B4` is a shared
helper, not the scheduler's threshold path (amends EXP-0043); and a `JSR`
at `return − 3` is necessary but **not sufficient** to confirm a stack
frame — confirm by execution.

**Whelk is no longer blocked by an absent model.** EXP-0040's timing can
now be *scoped* — only intervals inside the ability list and target
selection were paused — though scoping is not reinterpreting, and no
record has done the latter. Whelk is a boss with its own script and every
ATB unit ran on formation 14; that is the standing caveat.

Nothing in flight: no emulator, no resident instrumentation, worktree
clean, `main` in sync. SRAM virgin throughout. All gates green
(gofmt/build/vet/test, `ff6lab audit`, census sync 66 entries, archive
verify 8/8, restricted-file scan).

Exact next action: **EXP-0048 — name the invoker of the action-execution
path**, using a different instrument (exec-watch outward from the
confirmed `ROMCPU:$C2141D`, or a trace) since stack archaeology has
failed twice. It blocks nothing — the alternative is to bank the program
and take the **Whelk decision**, now an ordinary orchestration call.
