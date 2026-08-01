# Checkpoint 2026-07-31 — EXP-0034: golden route segment 3b (Unit 34)

## Current question
SCN-0001 program. Segment 3b closed. Next: EXP-0035 (segment 4) —
Narshe exterior navigation to the mine entrance.

## State
The golden route now runs **power-on → free field movement** as one
frame-scheduled probe, deterministically. Milestones 00, 01, 02, 03,
04 exist under `local_artifacts/scenarios/SCN-0001/`. Lab controls
unchanged (AllZeros RAM, virgin SRAM). No Mesen running.

## Work completed
EXP-0034 (record: docs/experiments/EXP-0034-golden-route-seg3b.md):
- Re-arming battle detector + slot-agnostic battle-end criterion
  (HP words in slots 3-9 all zero for 30 frames) + free-movement
  criterion (no battle re-arms for 5 400 frames).
- **The opening runs exactly four scripted battles**, identical in
  both runs: formations **2, 1, 2, 41** at entry frames 31 557 /
  34 953 / 36 828 / 39 500 (ends 32 736 / 36 004 / 38 059 / 40 975).
- Every staged `+$3F44` record matched the static ROM formation
  table (`$CF6200 + id×15`) byte-for-byte — formations 1
  (`0x0F620F`), 2 (`0x0F621E`), 41 (`0x0F6467`). The table is now
  verified across five independent encounters.
- **New monster reachable before Whelk: record 25**
  (`ROMFILE:0x0F0320`, 27 HP / 5 MP); with record 0 (40 HP / 15 MP)
  these are the only two staged by the scripted opening.
  Registered as CEN-MONSTER-0004.
- **Milestone `04-free-movement` at frame 46 375**, byte-identical
  across two power-on runs (WRAM `3e26bed9…`, screenshot
  `292013ab…`). This also **clears the determinism check deferred
  from EXP-0033**.
- Negative result preserved: a 1 200-frame quiet window captured the
  item-reward window instead of free movement (the reward chain plus
  the next battle's transition spans ~3 000 frames).

## What remains uncertain
- Monster **names** for ids 0/25 (name table unlocated); 30 of each
  32-byte record's fields unmapped.
- The scripted-battle **invocation opcode** (CEN-EVENT-0005) — the
  four entry frames are now precise watch anchors for it.
- Milestone-01 capture instability (CEN-QUIRK-0002) still open.
- Everything from mine entry onward is unchanged from program start.

## Tests and quality gates
census validate/sync clean, indexes regenerated, audit clean;
gofmt/build/vet/test at commit.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
`mesen/probes/EXP-0034.lua` (tracked schedule);
`local_artifacts/experiments/EXP-0034/` (per-battle entry/end
captures, beat screenshots, both runs' logs, hashes.sha256);
`local_artifacts/scenarios/SCN-0001/04-free-movement/` (dumps,
screens, canonical `04-free-movement.mss`, hashes.sha256).

## Exact next action
EXP-0035 (write record first): from milestone 04, replace the
"up-only" walk with a **direction-scripted** route across the Narshe
exterior toward the mine entrance. Recon first (bounded interactive
presses from `04-free-movement.mss` with screenshots) to find the
path, then encode it as absolute frame windows per direction, detect
the map transition by state (candidate: a write-watch on the map-id
word once located, else a large PPU/tileset change), capture
milestone `05-mines-entry`, and run the two-run determinism check.
This unit should also register the map-transition inputs for
CEN-WORLD-0004.
