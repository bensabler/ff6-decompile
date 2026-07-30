# Blockers

- **Session 003 disassembly claims unverified.** The session is now
  documented ([SESSION_003.md](../docs/sessions/SESSION_003.md)) and
  `battle` has tests, but the routine-level detail in
  `internal/game/battle/battle.go` still rests on lost ROM dumps. M4 work
  stays gated until the `ROMCPU:$C21300–$C21410` re-dump (open question
  1b) lands. Downgraded from hard blocker: raw store evidence is preserved
  and indexed.
- Mesen exact version relies on the Session 002 live recording (2.1.1); the
  app bundle plist is generic. Re-verify via `eval emu.getVersion()` on next
  launch.

## Resolved 2026-07-29

- ~~Not a git repository~~ — repository initialized and published to
  GitHub (`github.com/bensabler/ff6-decompile`).
- ~~Module path contradiction~~ — `go.mod` restored to
  `github.com/bensabler/ff6-decompile`; all gates pass.
- ~~Notation standards conflict~~ — resolved via
  [CONTRA-0001](../docs/contradictions/CONTRA-0001-address-notation.md):
  V4 prefixes supersede for new/updated docs; historical docs unchanged.
- ~~Package placement vs ARCHITECTURE.md~~ — `battle` and `chardata` moved
  under `internal/game/`.
