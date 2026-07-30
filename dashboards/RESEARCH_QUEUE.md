# Research Queue

## P0 — Provenance repair

- [x] Document Session 003 —
  [SESSION_003.md](../docs/sessions/SESSION_003.md) written (reconstructed
  record), promoted to 02/04/05/06/08, `battle` tests added (2026-07-29,
  overnight session)
- [ ] Re-dump `ROMCPU:$C21300–$C21410` and verify every `battle.go`
  disassembly claim (open question 1b) — restores FN-0003..0005 provenance

## P0 — Environment integrity (completed 2026-07-29)

- [x] Complete Version 4 migration —
  [V4_MIGRATION_REPORT.md](../docs/migrations/V4_MIGRATION_REPORT.md)
- [x] Record exact Mesen capability matrix (verified rows only) —
  [MESEN_CAPABILITY_MATRIX.md](../docs/research/MESEN_CAPABILITY_MATRIX.md)
- [x] Record ROM identity and local paths —
  [ROM_IDENTITY.md](../docs/research/ROM_IDENTITY.md)
- [x] Confirm latest checkpoint from the prior active Claude session —
  no prior checkpoint existed; Session 003 was interrupted without one
  (recorded in [SESSIONS.md](../indexes/SESSIONS.md) as SES-003)

## P1 — Battle lead

- [x] Identify callers of `ROMCPU:$C10DF3` — Confirmed in
  [SESSION_002](../docs/sessions/SESSION_002.md): `JSR $0DF3` at
  `ROMCPU:$C10200` inside PerFrameBattleUpdate (`ROMCPU:$C101FB`)
- [x] Reproduce `WRAM:+$2EB5` HP match across slots and values — Confirmed
  for slots 0–2 in [SESSION_002](../docs/sessions/SESSION_002.md)
- [x] Test the candidate `$20` record stride — Confirmed in
  [SESSION_002](../docs/sessions/SESSION_002.md) (iterations Y=0/$20/$40;
  records at `WRAM:+$2EB5/+$2ED5/+$2EF5/+$2F15`)
- [ ] Identify source region ownership around `WRAM:+$2E78` (open question
  #1 in [08_OPEN_QUESTIONS.md](../docs/sessions/08_OPEN_QUESTIONS.md);
  partially advanced by undocumented Session 003 evidence — resolve SES-003
  first)

## P1 — End-to-end vertical proofs

- [ ] Reconstruct one small graphics target from runtime state to Go validation.
- [ ] Reconstruct one short sound effect from CPU trigger to sequence/sample validation.
- [x] Implement one confirmed behavior in Go with tests —
  `chardata.CopyCharacterFields`
  ([02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md),
  tests in `internal/game/chardata/chardata_test.go`)

## P2

- [ ] Recover party battle record.
- [ ] Recover a menu font/tilemap path.
- [ ] Recover sample directory and one BRR sample.
