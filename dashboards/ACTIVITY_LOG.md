# Activity Log

- 2026-07-29 — Sessions 001–002: HP variable discovery, full decode of
  CopyCharacterFields (`ROMCPU:$C10DF3`), caller identification
  (`ROMCPU:$C101FB`), `chardata` Go package with tests. Records in
  `docs/sessions/`.
- 2026-07-29 — Session 003 (undocumented): bank-`$C2` HP/MP delta engine
  investigation; raw evidence in `mesen/out/`, implementation in
  `battle/battle.go`; interrupted by the V4 install before documentation.
- 2026-07-29 — Version 4 scaffold created.
- 2026-07-29 — `/audit-project`: 23 broken links fixed (V4 restructure),
  `mesen/bridge.lua` output path repaired, `*.mss` added to `.gitignore` and
  CI reject-list, stale doc paths in `chardata` fixed. Substantive findings
  in the audit report (chat) and migration report.
- 2026-07-29 — `/bootstrap-v4`: environment and ROM identity recorded,
  Mesen capability matrix filled from Session 002 evidence, indexes and
  dashboards initialized from canonical records, quality gates run,
  migration report and checkpoint written.
- 2026-07-29 — Post-migration decisions executed: module path restored to
  `github.com/bensabler/ff6-decompile`; `battle`/`chardata` moved under
  `internal/game/`; notation conflict resolved (CONTRA-0001); repository
  initialized in git and published publicly to GitHub.
