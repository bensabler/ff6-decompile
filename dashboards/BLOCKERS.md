# Blockers

## Hard blocker — no ATB model (2026-08-01, EXP-0040)

**Further Whelk execution is deferred until the project establishes a
usable ATB model.** Required before resuming: ACTIVE/WAIT behavior,
which submenu states qualify as pausing, the relevant timer domains,
and action-queue ordering/resolution semantics.

EXP-0040 made two attempts on Whelk and could not operate the battle
reliably: the head changed state between opening a submenu and reaching
target selection, queued actions resolved out of issue order, heals
landed on unintended allies, and on one occasion the target cursor was
resting on the shell at the moment of confirmation and had to be
cancelled. Those are operator/tooling failures against an unmodelled
timing system, not facts about Whelk.

Consequence for evidence: every head/shell transition observed in
EXP-0040 occurred while the emulator was driven through menus with
operator-length pauses, so **none of it may be used to characterize
Whelk's natural head/shell timing**. ACTIVE-mode and WAIT-mode timing
must be treated as separate experimental conditions.

**Whelk gameplay must not resume before this research.**
Record: `docs/experiments/EXP-0040-whelk-victory.md`.

**Status 2026-08-01:** the prerequisite audit is done and the ATB program
has begun. EXP-0041 closed the configuration half — `Bat.Mode` and
`Bat.Speed` are readable and settable from memory (`WRAM:+$1D4D` bits 3
and 0-2). EXP-0042 then established the program's **staging rule**: both
are sampled **once at battle entry** into `WRAM:+$3A8F` (Wait flag) and
`WRAM:+$3A90` (`255 − 24 × speed`), so ACTIVE/WAIT and Battle Speed
conditions must be set *before* entry, or injected at those two cells.

The blocker itself **remains open**: no timer domain, pause condition, or
queue semantics is known. What has changed is that there is now a
concrete way in — `+$3A90`'s consumer is unlocated and is very likely the
ATB rate path. Next: EXP-0043.

Soft items:

- Unattended sessions must use headless Mesen
  (`--testrunner --timeout=7200` with `FF6_OUT` set and frame-scheduled
  input): with the display locked, GUI Mesen either crashes at launch
  (Avalonia RenderTimer −6661) or silently never loads the script
  window (2026-07-30, EXP-0021). ~~GUI/testrunner input-latching parity
  is assumed, not yet verified~~ — **verified for the full opening
  schedule** (2026-08-01, EXP-0037): one GUI run and two headless runs
  produce byte-identical WRAM at frame 51 578 and identical
  event-flag write timelines (frame+addr+value+PC). Parity beyond this
  schedule's input pattern remains untested.
- Live MP verification (research queue) additionally needs a battle
  savestate with an MP-consuming action — the Magitek intro battle has
  none.
- Mesen exact version rests on the Session 002 live recording (2.1.1);
  the app bundle plist is generic. Re-verify via the bridge on next
  launch (no `emu.getVersion()` in this build's Lua API — use the log
  window header).
- ~~Bridge v2 unvalidated~~ — validated live 2026-07-30 (EXP-0020);
  probe loading, id'd responses, duplicate suppression, transcripts all
  confirmed working.
- MP semantic labels (`CurrentMP`/`MaxMP`/`ApplyMPDelta`,
  `WRAM:+$3C08`/`+$3C30`) remain Strong hypothesis pending the live MP
  verification experiment (research queue).

## Resolved history
See ACTIVITY_LOG.md — Session 003 verification (EXP-0001), `$2E78`
producer (EXP-0003), pending-delta producers (EXP-0004/0006),
authoritative-HP layer (EXP-0003/0005), git/module/notation items
(2026-07-29), and the `.gitignore` collision that untracked
`cmd/ff6lab/main.go` (fixed 2026-07-30).
