# Blockers

- **Session 003 is undocumented.** `internal/game/battle/battle.go`
  implements routines (`ROMCPU:$C21323`/`$C21350`/`$C21390`) with no session
  record, no canonical doc entries, and no tests. Blocks trusting its
  confidence claims and blocks milestone M4 work. Raw evidence:
  `mesen/out/events.log`, `mesen/out/checkpoint3-mines.mss`.
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
