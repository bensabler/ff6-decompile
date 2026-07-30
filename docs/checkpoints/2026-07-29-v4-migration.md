# Checkpoint 2026-07-29 — V4 migration complete

## Current question
None active. The last research question (Session 003: "what computes the HP
changes upstream of the display copy?") was answered in raw evidence but
never documented — documenting it is the next unit of work.

## State
Maintenance/migration state. Version 4 scaffold merged, audited
(`/audit-project`), and migrated (`/bootstrap-v4`). No emulator running, no
active experiment.

## Confirmed before this session
- CopyCharacterFields `ROMCPU:$C10DF3`–`$C10E66` fully disassembled and
  live-verified; caller `ROMCPU:$C10200` in PerFrameBattleUpdate
  (`ROMCPU:$C101FB`); records at `WRAM:+$2EB5`+n×`$20`; slot mask at
  `WRAM:+$61AD`. See indexes FN-0001/0002, VAR-0001..0004, ST-0001/0002.

## Work completed
- Audit: 23 broken links fixed; `mesen/bridge.lua` OUT path repaired
  (`tools/mesen/out/` → `mesen/out/`); `*.mss` added to `.gitignore` and CI;
  `chardata` doc paths corrected; `make fuzz` rewritten (was unconditionally
  broken with multiple packages).
- Migration: environment + ROM identity recorded; Mesen capability matrix
  filled (verified rows only); indexes SES/FN/VAR/ST populated; all eight
  dashboards rewritten from canonical records;
  `docs/migrations/V4_MIGRATION_REPORT.md` written.

## Last raw observation
(From undocumented Session 003, `mesen/out/events.log`:) NEW-3BF4-WRITER
exec captures at `ROMCPU:$C2134A`, `$C2133B`, `$C21399` (and battle-init/
teardown writers `$C227B4`, `$C223F6`, `$C206BC`, `$C36A5E`) writing the
`WRAM:+$3BF4` region, with registers and stack snapshots per line.

## Active emulator state
None running. Savestates preserved: `mesen/out/checkpoint1.mss`,
`checkpoint2.mss`, `checkpoint3-mines.mss` (SHA-256s in the migration
report). Battle savestate `Final Fantasy III (USA)_11.mss` in Mesen2 app
support (Narshe intro battle, 3-member party).

## Breakpoints/watchers
None active. When loaded, `mesen/bridge.lua` re-arms exec logging at
`ROMCPU:$C10DF3` and a write watch on `WRAM:+$2E78`–`+$2E7F` (source HP
array, logged as NEW-HP-WRITER). Note: the NEW-3BF4-WRITER instrumentation
that produced the Session 003 evidence is **not** in the current script — it
was injected transiently (likely via the bridge's `eval` command); only its
output in `events.log` survives.

## Evidence paths
`mesen/out/` (hashes in
[V4_MIGRATION_REPORT.md](../migrations/V4_MIGRATION_REPORT.md) §3);
ROM identity in [ROM_IDENTITY.md](../research/ROM_IDENTITY.md).

## Files changed
Audit: `docs/sessions/*.md` (links), `mesen/bridge.lua`, `.gitignore`,
`.github/workflows/ci.yml`, `chardata/chardata.go` (comment), `Makefile`.
Migration: `docs/research/ROM_IDENTITY.md`,
`docs/research/MESEN_CAPABILITY_MATRIX.md`, `indexes/{SESSIONS,FUNCTIONS,
VARIABLES,STRUCTURES}.md`, `dashboards/*.md`,
`docs/migrations/V4_MIGRATION_REPORT.md`, this checkpoint,
`docs/checkpoints/LATEST.md`.

## Tests and quality gates
gofmt clean; `go build`/`go vet`/`go test ./...` pass; `make fuzz` pass
(both targets, 30s each). Run 2026-07-29 after all edits.

## Git status
**Not a git repository** — nothing is version-controlled. Operator decision
pending (migration report §10.1).

## Unresolved decisions
git init; module path restore (`github.com/bensabler/ff6-decompile`);
`battle`/`chardata` package placement vs ARCHITECTURE.md; notation
supersession record. Details: migration report §10.

## Blockers
See [dashboards/BLOCKERS.md](../../dashboards/BLOCKERS.md). Top: Session 003
undocumented (blocks trusting `battle/battle.go`; blocks M4).

## Exact next action
Write `docs/sessions/SESSION_003.md` from `mesen/out/events.log` and the
interpretation in `internal/game/battle/battle.go`; promote records into
`docs/sessions/02/04/05/06/08`; add table-driven `battle` tests mirroring
the chardata pattern; then correct battle.go's evidence-path comment.

## Addendum (2026-07-29, post-checkpoint, operator-approved)
Module path restored to `github.com/bensabler/ff6-decompile`;
`battle`/`chardata` moved under `internal/game/`; notation conflict resolved
(CONTRA-0001); repository initialized in git and published publicly to
GitHub. All quality gates re-run and passing after these changes.

## Recommended next command
`/resume-session` (fresh session), then execute the exact next action above
as a documentation unit — not a new emulator experiment.
