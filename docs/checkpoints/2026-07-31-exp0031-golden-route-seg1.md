# Checkpoint 2026-07-31 — EXP-0031: golden route segment 1 (Unit 31)

## Current question
SCN-0001 program. Segment 1 closed; next is segment 2 (EXP-0032:
Narshe-entry stall → first scripted battle).

## State
Milestone `00-new-game` established and determinism-proven. Lab
controlled variables now in force for the whole SCN-0001 program:
- Mesen SNES `RamPowerOnState` = **AllZeros** (was Random; changed
  2026-07-31; original settings backed up at
  `local_artifacts/backups/mesen-settings-2026-07-31.json`).
- **Virgin SRAM**: the prior `Final Fantasy III (USA).srm` was
  removed from Mesen's Saves dir; backup at
  `local_artifacts/backups/ff3usa-2026-07-30.srm`
  (sha256 `6afbcf1e…52cc7a`). Restore it if pre-program states are
  ever needed.

## Work completed
EXP-0031 (record: docs/experiments/EXP-0031-golden-route-newgame.md):
- Phase A recon: title at ~frame 2969 (fade by ~3709); attract loop
  ~21000 frames alternating dimmed opening replay + gameplay demo;
  title ignores holds/sparse presses ($4219=$10 verified latched) but
  accepts rapid start+a edge toggling; virgin SRAM → no save-select,
  real opening starts directly and auto-runs to the input-waiting
  Narshe-entry dialogue (advances on A).
- Phase B/C: probe `mesen/probes/EXP-0031.lua` (tracked input
  schedule); two fresh power-on runs byte-identical at frame 5200
  (milestone WRAM `35c76d03…`) and frame 30000 (stall WRAM
  `0f4369d5…`); milestone savestate `00-new-game.mss` load-validated.
- Census: +CEN-EVENT-0004 (title/attract), +CEN-QUIRK-0001 (boot
  uninitialized reads at $7E7BF7/F8, $7E7CAD/AE, $7E7DB2 — first
  QUIRK entry); CEN-EVENT-0001 and CEN-SAVE-0001 annotated.
- SCN-0001 record + JSON updated (B01/B02 progress, golden_route
  segment 1 with assertion hashes).

## Gotchas recorded
`emu.createSavestate`, like loadSavestate, only works inside a
main-CPU exec callback (probe uses a one-shot callback).

## Tests and quality gates
census validate/sync clean, indexes regenerated, audit clean;
gofmt/build/vet/test at commit.

## Git status
main; this unit committed and pushed (operator push authorization for
this session).

## Active instrumentation and evidence
No Mesen running. Evidence: local_artifacts/experiments/EXP-0031/,
local_artifacts/scenarios/SCN-0001/00-new-game/ (hashes.sha256 in
both).

## Exact next action
EXP-0032 (write record first): from `00-new-game.mss` or a fresh
power-on replay, extend the frame schedule: A-press at the
Narshe-entry stall (~frame 20600+; box waits indefinitely so a late
scheduled press is determinism-safe), map subsequent stalls/beats by
recon, capture milestones `01-opening-cinematic` (a stable
mid-presentation anchor), `02-narshe-entry`, and
`03-first-scripted-battle` (battle entry state), each with WRAM
assertions; then a two-run determinism check of the full extended
schedule.
