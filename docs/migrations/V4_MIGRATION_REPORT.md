# Version 4 Migration Report

- **Date:** 2026-07-29
- **Performed by:** Claude (Fable 5) via `/audit-project` followed by `/bootstrap-v4`
- **Scope:** Merge of the Version 4 scaffold (210 generated files,
  [PACKAGE_MANIFEST.json](../../PACKAGE_MANIFEST.json)) into the existing
  research repository; migration of prior canonical records into V4
  dashboards and indexes.

## 1. Inventory

244 files at migration time: 210 scaffold files (all present and verified
against `PACKAGE_MANIFEST.json`; 0 missing) plus 32 pre-existing
research/tooling files and OS cruft (`.DS_Store`).

Pre-existing research preserved in place:

- `docs/sessions/00–08_*.md` — numbered canonical docs (moved by the V4
  restructure from `docs/` into `docs/sessions/`)
- `docs/sessions/SESSION_001.md`, `SESSION_002.md` — session records
- `chardata/` — implemented + tested reconstruction of
  `ROMCPU:$C10DF3` (CopyCharacterFields)
- `battle/` — implementation from an **undocumented** session (see §3)
- `mesen/bridge.lua` + `mesen/out/` — Lua bridge and raw runtime evidence
  (moved from `tools/mesen/` by the restructure)

## 2. Preservation of prior canonical records

No prior records were overwritten or rewritten. The restructure had broken
24 relative links inside the moved documents; 23 were repaired mechanically
during `/audit-project` (same date). The 24th points to
`FF6_Decompilation_Session_01_Summary.md`, the raw source from which
SESSION_001 was reconstructed — the file no longer exists anywhere searched
(repository, Desktop, Documents, Downloads). **Provenance gap:** SESSION_001
is now the only surviving record of Session 001.

## 3. Session 003 documentation debt (P0)

`battle/battle.go` implements the bank-`$C2` HP/MP delta engine
(`ROMCPU:$C21323`, `$C21350`, `$C21390`) and cites "Session 003"
verification, but no session record, canonical doc entries, or tests exist.
Raw evidence confirming the session happened is preserved:

| Evidence | SHA-256 |
|---|---|
| `mesen/out/events.log` (NEW-3BF4-WRITER captures at `ROMCPU:$C2134A`, `$C2133B`, `$C21399`, `$C227B4`, `$C223F6`, `$C206BC`, `$C36A5E`) | `bcfc7f4ca5034cbd13eb0eaaeaf08cd65f21467e49c1aaaed2938445f2a99d03` |
| `mesen/out/checkpoint3-mines.mss` (savestate, 20:46 local) | `61c37807f1e3dda3c7c6bdb901856c864e580167edd951da1ae90800b72168a0` |
| `mesen/out/checkpoint1.mss` | `b1b95cbeb56d122af4d7136213277c521edb3f3fe9344ad58edca4344f8030d5` |
| `mesen/out/checkpoint2.mss` | `38883818cf4b7862f68af5f5258cf2df5d56adf82719a971675f0963e807cfe4` |
| `mesen/out/routine_C10DF3.hex` | `51af2e91e6ee835a3e336cbce33dbf65fdeff5c94f6df2d8b471e693282e1f4b` |
| `mesen/out/caller_C101C0.hex` | `4e452c79a31bddacdccf0cde81d082beb697231ac182dc5e08a34eadf4b4ccf9` |
| `mesen/out/diff.txt` | `bb76a8b07841cef6a1a3c7c59dc8b27ba967ae87ec9f10319f6385c40805d0b6` |
| `mesen/out/diff2.txt` | `c734c75608d650209803b68a1a49be342197ac5fe7db2f643b81b04834ea200b` |
| `mesen/bridge.lua` (post path-fix) | `915a7aa7e863dacc6213e0054f7cc335f4c8538785eebad646787e4d309122fb` |

The session was interrupted by the V4 install (~21:04 local) without
`/checkpoint`. Tracked as SES-003 / FN-0003..FN-0005 in the indexes and as
the top item in [RESEARCH_QUEUE.md](../../dashboards/RESEARCH_QUEUE.md).

Additional provenance nuance: the NEW-3BF4-WRITER watch that produced the
`events.log` lines is not present in the current `mesen/bridge.lua` (which
watches `WRAM:+$2E78`–`+$2E7F` instead). It was evidently injected
transiently — likely via the bridge's `eval` command — so only its output
survives, not its code. SESSION_003 must state this when citing the log.

## 4. Duplicate/conflicting command and skill files

- `.claude/commands` (28) and `.claude/skills` (27 + `_shared`): no
  duplicates, no leftovers from earlier versions.
- One dangling reference: `docs/sessions/01_REVERSE_ENGINEERING_RULES.md`
  names the retired `.claude/skills/ff6-reconstruction-skill/SKILL.md`.
- One standards conflict (open): 01's contextual naked-`$` address notation
  vs `docs/research/ADDRESS_NOTATION.md` mandatory domain prefixes. Needs a
  contradiction record and an explicit supersession decision; pre-V4 docs
  still use the old style. Listed in
  [BLOCKERS.md](../../dashboards/BLOCKERS.md).

## 5. Environment record

| Component | Version |
|---|---|
| OS | macOS 26.5.2 (build 25F84), Darwin 25.5.0, x86_64 |
| Go | go1.26.0 darwin/amd64 |
| Git | 2.52.0 |
| Claude Code | desktop app session, model Fable 5 (`claude-fable-5`); CLI version string not queryable from this shell |
| Mesen | 2.1.1 macOS x86_64 (Intel), `~/Desktop/Mesen.app` — recorded live in SESSION_002; bundle `Info.plist` reports generic 1.0.0; re-verify via bridge `eval emu.getVersion()` on next launch |

## 6. ROM identity

Computed in place (no copy made): SHA-256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
3,145,728 bytes, no copier header. Full record with confidence per field:
[ROM_IDENTITY.md](../research/ROM_IDENTITY.md).

## 7. `local_artifacts/` and `.gitignore` verification

- `local_artifacts/` exists and contains only its README. The working ROM
  and savestates currently live outside the repository (Desktop /
  Mesen2 app support), which is acceptable under the asset policy.
- `.gitignore` covers `local_artifacts/`, ROM/audio extensions, and `out/`
  (which matches `mesen/out/`). `*.mss` (Mesen savestate extension) was
  missing from both `.gitignore` and the CI reject-list; added during
  `/audit-project`.
- **Caveat: this directory is not a git repository.** Every git-based
  protection (ignore rules, CI, clean-repository check, branch policy) is
  inert until `git init` + initial commit. Operator decision pending.

## 8. Initialization of dashboards, indexes, manifests

- Indexes populated from existing canonical records:
  [SESSIONS](../../indexes/SESSIONS.md) (SES-001..003),
  [FUNCTIONS](../../indexes/FUNCTIONS.md) (FN-0001..0005),
  [VARIABLES](../../indexes/VARIABLES.md) (VAR-0001..0004),
  [STRUCTURES](../../indexes/STRUCTURES.md) (ST-0001..0002). Discovery,
  experiment, contradiction, and asset indexes remain empty (no canonical
  records of those types exist yet).
- Dashboards rewritten to reflect canonical state: CURRENT_FOCUS,
  ACTIVITY_LOG, BLOCKERS, RESEARCH_QUEUE (Session 002 results marked
  complete with record links), OPEN_HYPOTHESES (H-BATTLE-0001 resolved;
  H-BATTLE-0002..0005 added from existing records), MILESTONES (M1
  complete), STATISTICS.
- `manifests/*.json` left as valid empty collections. Retroactive
  experiment records for Sessions 001/002 were **not** fabricated; their
  method/falsification details live in the session records. Future
  experiments get manifest entries at execution time.

## 9. Go quality gates (2026-07-29, post-migration)

```text
gofmt -l .      clean
go build ./...  pass
go vet ./...    pass
go test ./...   pass (chardata, internal/audio/brr,
                internal/graphics/tile4bpp, internal/platform/bgr555)
make fuzz       pass (FuzzParseHeader 9.2M execs, FuzzDecode 7.0M execs,
                30s each, no failures)
```

The scaffold's `make fuzz` was broken (`-fuzz` with multiple packages is
rejected by the Go tool); rewritten to iterate per package/target. The
scaffold's own `BUILD_VERIFICATION.json` shows `go vet`/`go test` never ran
in the generator environment (network failure) — superseded by the local
results above.

## 10. Open decisions for the operator

1. `git init` + initial commit (activates all repository protections).
2. Restore module path `github.com/bensabler/ff6-decompile` in `go.mod`
   (+ one import in `cmd/ff6lab/main.go`); current placeholder is
   `example.com/ff6-reconstruction`.
3. Package placement: `battle/` and `chardata/` at repo root vs
   `ARCHITECTURE.md`'s `internal/game` — move packages or amend the doc.
4. Notation supersession decision (§4).
5. Session 003 documentation session (P0 in the research queue).

### Resolution (2026-07-29, same day, operator-approved)

Items 1–4 executed: repository initialized and published publicly to
`github.com/bensabler/ff6-decompile`; module path restored (all gates
re-run, pass); `battle` and `chardata` moved to `internal/game/`;
notation conflict resolved as
[CONTRA-0001](../contradictions/CONTRA-0001-address-notation.md) (V4
prefixes supersede for new/updated docs; historical docs keep their declared
notation). Item 5 remains the top of the research queue.

## Checkpoint

[docs/checkpoints/2026-07-29-v4-migration.md](../checkpoints/2026-07-29-v4-migration.md)
