# Checkpoint 2026-08-01 — EXP-0035: route to the mines, partial (Unit 35)

## Current question
SCN-0001 program. Segment 4 recon closed as **partial**. Next:
EXP-0036 — encode the leg table as a scheduled route and capture
milestone `05-mines-entry`.

## State
Golden route (scheduled, deterministic): power-on → milestones 00-04.
Beyond that, EXP-0035 walked the rest of the way to the mines
interior **interactively**, so milestone 05 is deliberately **not**
claimed — the milestone directory stays empty until a scheduled
power-on run produces it, matching how the false milestone 04 was
handled in EXP-0033.

Lab controls unchanged (AllZeros RAM, virgin SRAM; originals in
`local_artifacts/backups/`). No Mesen running. Working tree committed.

## Work completed
EXP-0035 (record: docs/experiments/EXP-0035-golden-route-seg4.md):
- **Player tile position located:** `WRAM:+$00AF` = X, `+$00B0` = Y
  (each moves by exactly 1 per step on its own axis, stable on the
  other). Aliases at `+$0541/42`, `+$0543/44`, `+$0545/46`,
  `+$087A/7B`. This turned the recon coordinate-driven.
- **Candidate map-id byte `WRAM:+$1EA5`:** `$00` on
  snowfield/opening-approach states, `$05` on all four Narshe
  exterior free-walk states, `$0D` on both mines-interior states.
  Tension preserved: milestone 02 looks like Narshe exterior but
  reads `$00`. Tentative — falsifier recorded.
- **Full route mapped** from milestone 04 ($1A,$2A) to the mines
  interior ($26,$1C), 11 legs, tabulated in the experiment record and
  the scenario manifest. The climb is a **zigzag** — a plain hold-up
  stalls at Y=$21 and again at Y=$16, so earlier segments' up-only
  cadence cannot walk it.
- **Fifth scripted battle registered** (CEN-EVENT-0007), triggered at
  tile ($1E,$27) after free movement — outside EXP-0034's four-battle
  opening chain. Fought to victory during recon; its formation id was
  not captured (no probe armed).
- Census: +CEN-WORLD-0007 (position + map-id candidates),
  +CEN-EVENT-0007 (fifth battle and shaft gate).

## What remains uncertain
- Milestone 05 uncaptured; segment 4 has **zero** determinism runs.
- Battle 5's formation id and monster records unknown.
- `+$1EA5` semantics unresolved (see falsifier).
- Producers of `+$00AF`/`+$00B0` untraced; alias blocks unidentified.
- Milestone-01 capture instability (CEN-QUIRK-0002) still open.
- Everything from mines traversal to Whelk unchanged.

## Tests and quality gates
census validate/sync clean, indexes regenerated, audit clean;
gofmt/build/vet/test at commit.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
No probe for this unit (interactive recon via the bridge).
`local_artifacts/experiments/EXP-0035/` — 57 files: recon
screenshots per leg, WRAM position diffs (`pos-*.hex`), the
`$1E00` window dumps used for the map-id comparison, a pre-battle-5
savestate, the mines-interior recon savestate, `hashes.sha256`.

## Exact next action
EXP-0036 (write record first): create `mesen/probes/EXP-0036.lua` as
EXP-0034's probe plus a leg-scheduler driven by the tabulated route.
Because walls normalize position, encode each leg as a generous
absolute frame window per direction; keep the re-arming battle
detector so battle 5 is fought and its `+$11E0` formation id and
staged `+$3F44` record are logged. Add a write-watch on `WRAM:+$1EA5`
to (a) detect the map transition for the milestone capture and
(b) test whether its writer also drives tileset/tilemap loading —
the recorded falsifier for the map-id claim. Capture milestone
`05-mines-entry` at a settled frame after the transition, then run
two fresh power-on runs and byte-compare the milestone WRAM.
