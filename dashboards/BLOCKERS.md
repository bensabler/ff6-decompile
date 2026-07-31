# Blockers

No hard blockers (2026-07-30 maintenance audit).

Soft items:

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
