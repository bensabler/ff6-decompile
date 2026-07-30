# Activity Log

- 2026-07-30 (overnight 2) — Units 15–16 / EXP-0014+0015: **the base
  damage formula is decoded and numerically closed** — `base = power×4 +
  (power×$11AE×$11AF)>>5` (live: 60/28/4 → 450 → ×(255−58)/256+1 →
  346, the observed Fire Beam). The `$C20DD1` 24-bit-shift helper
  explains the `$C247B7` wrapper's memory layout. Go:
  `BaseAmountStandard` + golden chain test. The standard damage
  pipeline is now byte-exact from formula to HUD. New: #26 physical
  path, #27 stat meanings, #28 variance hunt.

- 2026-07-30 (daytime) — Units 13–14 / EXP-0012+0013: arithmetic helpers
  decoded — `$C2370B` is a ×1.5-per-count chain (**randomness hypothesis
  refuted**); `$C24781` is the SNES hardware 8×8 multiply, making the
  defense scaling exactly `(amount×(255−def))/256 + 1`. Go:
  `Scale256`, `ApplyDefense`, `ChainBoost` + tables + a composition
  fuzz target. Questions #24/#25 closed; #23 (`+$11B0` producer) is the
  innermost frontier.

- 2026-07-30 (daytime) — Unit 12 / EXP-0011: **base-amount routines
  decoded** — variant A consumes a precomputed base (`+$11B0`) and
  applies defense ((255−def)/256 shape, pair at `+$3BB8,Y`), flag
  halvings, party-vs-party halving, and a final `$C2370B` transform;
  variant B = fraction-of-HP with min 1; `+$11A3` bit 7 retargets X/Y
  by `$14` (HP→MP). Family census reaches **twelve** `$14`-stride
  arrays (05 master table). No Go this unit (helpers unverified).
  Questions #23–25 opened.

- 2026-07-30 (daytime) — Units 10–11 / EXP-0009+0010: formula frame
  verified (`JSR $0B83` at `$C23469`, ten-slot target loop) and the
  **elemental-modifier block decoded byte-exact** (`$C20B83`–`$C20C2C`):
  nullify/flip-to-heal/zero/halve/double transforms over per-target
  16-bit element masks in two more `$14`-stride family arrays
  (`+$3BCC`, `+$3BE0`). Go: `battle.ApplyElementResponse` (behavior-
  derived names) + 14-case table test. Base-amount routines
  (`$C20C9E`/`$C20D87`) opened as question #22 — the innermost formula.

- 2026-07-30 (daytime) — Unit 8 / EXP-0007: **DP `$F0` = final per-hit
  damage, Confirmed on three anchors** (array arithmetic 9+4+2+6=21;
  HUD mid-values; captured popup "6"). A 346-damage Fire Beam killed the
  Were-Rat on camera. X/Y = attacker/target slot×2 (strong); `$F2`
  direction observation recorded. Question #21 narrowed: the full
  formula runs upstream of `$C20C28`; deeper-return lead `$6B16`.

- 2026-07-30 (daytime) — Unit 7 / EXP-0006: **pending-delta accumulator
  decoded byte-exact** (`ROMCPU:$C20C76`): amounts accumulate per slot
  with `$FFFF`-sentinel init and a **9999 cap** (`$270F`) — the damage
  cap located in code. One routine serves both pending arrays via a
  `Y += $14` retarget. Caller gate block at `$C20C2D` identified;
  EXP-0004's `$C20434` caller inference corrected (stack misparse).
  Go: `battle.AccumulatePending` + table tests. FN-0008 pending index.
  New question #21 (the actual damage formula computing DP `$F0`).

- 2026-07-30 (daytime) — Unit 6 / EXP-0005: **enemy slots confirmed** —
  enemy HP observed at entries 4–5 of the unified arrays (24/35), damaged
  and zeroed by the same delta-engine stores and death handler while
  VICKS won the encounter. H-BATTLE-0008 resolved Confirmed. Go model
  refactored: `battle.PartySlots` → `battle.BattleSlots` (10 entries,
  slot-uniform engine), tests extended (enemy-slot case). New writer
  `ROMCPU:$C22CCE` queued as question 19b.

- 2026-07-29 (overnight) — Unit 5 / EXP-0004: pending-delta producers
  identified at the PC level (setter `ROMCPU:$C20C9B`, sweepers
  `$C2638E`/`$C26391`); arrays shown transient (`$FFFF` between events);
  **structural discovery: the battle arrays are 10 entries wide**
  (`$14` stride family `+$3BF4/+$3C08/+$3C1C/+$3C30`; entry-9 write
  observed) — 4 party + 6 enemy candidate. Refresh trigger localized to
  the post-fetch driver (`JSR $069B` at `$C21409` → `$C2069B` → tail
  path). Questions #13/#19 narrowed; battle.go's 4-slot model flagged as
  party-view-only pending confirmation.

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
