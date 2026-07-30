# Activity Log

- 2026-07-29 (overnight) — Unit 4 / EXP-0003: open question #1 **answered**.
  PartyDisplaySourceRefresh (`ROMCPU:$C25D26`) decoded byte-exact: copies
  all six authoritative battle arrays into the `+$2E78` display family.
  H-BATTLE-0002 and H-BATTLE-0006 resolved (Confirmed); H-BATTLE-0004
  upgraded; `+$2E98` identified as the status-word copy (`+$3EE4`) and
  `+$2EA0` as the `+$3EF8` copy. New questions 19–20.

- 2026-07-29 (overnight) — Unit 3 / EXP-0002: PerFrameBattleUpdate fires
  only in battle across all tested contexts (≈175k non-battle frames, zero
  fires; two battles positive). Corrections: "every frame" softened to
  phase-dependent ≈0.23–0.87/frame; second one-shot caller `JSR
  ROMCPU:$C11090` discovered at battle entry. Hazard: Mesen auto-save
  destroyed the `_11.mss` battle state (recorded in capability matrix).
  Battle-init writers reproduced in an independent encounter.

- 2026-07-29 (overnight) — Unit 2 / EXP-0001: launched Mesen 2.1.1 with the
  bridge, dumped `ROMCPU:$C212F0–$C2141F`, verified every battle.go
  disassembly claim byte-exact; FN-0003..0006 upgraded to Confirmed
  (code); new unknowns recorded (`+$11A2` selector, dispatch tail, fetch
  gates); ROM mapping upgraded to Confirmed HiROM via Mesen header parse;
  Session 003 evidence archived under `mesen/out/session003/`.

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
- 2026-07-29 (overnight) — Unit 1: SESSION_003 reconstructed record written
  from raw evidence; battle delta engine promoted to 02/04/05/06/08 with
  honest confidence split (stores Confirmed, addresses Strong hypothesis,
  disassembly Unknown); `internal/game/battle` gained table-driven tests;
  battle.go provenance comment corrected.
- 2026-07-29 — Post-migration decisions executed: module path restored to
  `github.com/bensabler/ff6-decompile`; `battle`/`chardata` moved under
  `internal/game/`; notation conflict resolved (CONTRA-0001); repository
  initialized in git and published publicly to GitHub.
