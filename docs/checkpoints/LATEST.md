# Latest Checkpoint

**[2026-08-01 — EXP-0046: the action-execution path](2026-08-01-exp0046-action-execution-path.md)**

State: the routine EXP-0045 left unnamed is named.
**`ROMCPU:$C201BE` (`INC $3219,X`)** advances `+$3218` by `$0100`,
guarded by `+$3AA0` **bit 3**, reached from an action-execution path in
`ROMCPU:$C207xx`–`$C209xx` that runs **outside** `$C21124`'s gate. Every
gated write fell on a **single frame** — the deferred completion is a
burst, not a drain.

One capture looked like a contradiction: `$C211BA`, inside the
gauge-advance routine the gate skips, wrote while gated. The stacks
resolved it rather than explaining it away — **`$C211B4` is a shared
helper** (`ORA $3AA0,X / STA / RTS`) with at least two entry points: the
scheduler at `$C211B2` with `A=$20`, and `$C208C6` with `A=$50`.

**Correction propagated:** EXP-0043 read the `+$3AA0` store at `$C211B7`
as the scheduler's threshold path. It belongs to a shared helper, so a PC
there does not imply the scheduler ran. The memory map is amended; no
earlier conclusion depended on the stronger reading.

`+$3AA0` is filling in: **bit 3** gates the increment (set by `$C20974`),
**bit 6** is the pending-action marker (EXP-0045), **bit 7** is cleared
on completion by `$C20795`. Bits 0-2, 4, 5 remain Unknown, as do
`+$3204`, the caller chain above `$C208C6`, and why the completion fired
122 frames after the gate engaged.

New CEN-BATTLE-0013. No blockers. All gates clean; no background
processes; SRAM virgin.

Exact next action: **EXP-0047 — what invokes the execution path, and
when.** Decode the captured chain (`$C208B1`, `$C21420`, `$C20EB6`) and
exec-watch the entry point across a gated interval to recover the
invocation cadence. That is the last structural piece before the action
lifecycle and queue model can be written down.
