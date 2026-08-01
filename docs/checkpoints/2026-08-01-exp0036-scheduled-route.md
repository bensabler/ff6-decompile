# Checkpoint 2026-08-01 — EXP-0036: scheduled route to the mines, partial (Unit 36)

## Current question
SCN-0001 segment 4. The route controller is built and works; the
milestone-05 acceptance runs are outstanding.

## State
The golden route now walks **power-on → mines interior** with no manual
correction, via a 17-leg state-driven route controller. Milestone 05 is
**deliberately not claimed**: the acceptance criteria require three
scheduled power-on runs landing at (`$26`,`$1C`) in the mines, and the
final 17-leg encoding has not completed that set. The milestone
directory is empty; artifacts from the earlier, invalid 16-leg capture
were discarded rather than kept.

Lab controls unchanged (AllZeros RAM, virgin SRAM; originals in
`local_artifacts/backups/`).

## Work completed
EXP-0036 (record: docs/experiments/EXP-0036-scheduled-route-to-mines.md):

- **Route controller.** Replaced hold-one-direction cadences with legs
  advanced by state: position targets on `WRAM:+$00AF`/`+$00B0`
  (overshoot-tolerant), battle-start/battle-end edges, and elapsed
  settles. Per-leg timeouts fail the run and **name the earliest
  divergent leg** instead of retrying or self-correcting — which is how
  every finding below was surfaced.
- **Two schedule defects found and fixed, both preserved as evidence:**
  battle-end detection had been gated behind the route phase (so the
  opening's battles 2-4 never re-armed); and legs 1-3 reached
  (`$1E`,`$27`) where nothing happened.
- **EXP-0035 correction.** Its condensed leg table had dropped an
  intermediate `up` step: the guard dialogue is at (`$1E`,**`$25`**).
  That record now carries the correction inline. The trigger is
  **contact-based** — standing on the tile tapping A never fires it.
- **Battle 5 = formation 84** (`$0054`), staged record byte-identical
  to `ROMFILE:0x0F66EC` (sixth independent verification of the
  EXP-0030 formation table), monsters **{27, 27, 0, 0}**. New
  pre-Whelk monster: **record 27** (`ROMFILE:0x0F0360`), 115 HP /
  30 MP, the 115 anchored to a live enemy HP word from EXP-0035.
- **`+$1EA5` falsifier fired.** It reaches `$0D` during the shaft
  dialogue while the party is still visibly on the exterior and has
  not moved — so it is **not** a simple current-map byte. Recorded
  reading: a map-load target / event-state value written by
  `ROMCPU:$C0B5B6` (which writes every observed change:
  `$00`→`$01`→`$05`→`$0D`). Confidence **not** promoted; transition
  detection moved to the player-position jump.
- **Aliases answered within bounds:** the four mirror blocks track the
  leader with a one-tile lag on the trailing member during turns and
  re-converge within a leg (follower-chain behaviour). No divergence
  affects route control. Position bytes are **not field-meaningful
  during battle** (`1E,00` entry, `00,00` end), so position legs are
  suspended while a battle owns the screen.
- **Go implementation:** `internal/scenario/route` — `Leg`, `Route`,
  `Runner`, advancement predicate, `MinesRoute()`. Tests cover leg
  sequencing, direction changes, overshoot-tolerant completion,
  timeout naming the divergent leg, divergence that must not advance,
  battle interruption/resume, battle-edge legs, elapsed/map-change
  legs, and validation. `probe_sync_test.go` parses the Lua probe's
  `ROUTE` table and asserts it matches the model leg for leg, so the
  two encodings cannot drift silently.

## What remains uncertain
- **Milestone 05 unclaimed**: three acceptance runs outstanding.
- `+$1EA5` semantics unresolved (map id vs map-load target).
- Producers of `+$00AF`/`+$00B0` and alias-block ownership untraced.
- Battle 5's invocation opcode unknown (shared with CEN-EVENT-0005).
- Milestone-01 capture instability (CEN-QUIRK-0002) still open.

## Tests and quality gates
`go test ./internal/scenario/...` green (10 test groups incl. the
probe-sync guard); census validate/sync clean; audit clean; full gates
run at commit.

## Git status
main; unit committed and pushed.

## Active instrumentation and evidence
`mesen/probes/EXP-0036.lua` (tracked 17-leg schedule).
`local_artifacts/experiments/EXP-0036/` holds every iteration's log,
including the two timeout runs and the premature-transition run.

## Exact next action
Run `mesen/probes/EXP-0036.lua` from power-on **three times**
(`FF6_RUN=run1|run2|run3`, headless testrunner, `FF6_OUT` set), each
taking ~13 minutes. Confirm all three log `LEG 17 END … pos=26,1C` and
`MILESTONE 05 … map=0D`, then byte-compare the three
`local_artifacts/scenarios/SCN-0001/05-mines-entry/run*-wram.bin`
dumps. Only if all three agree: promote one savestate to
`05-mines-entry.mss`, write `hashes.sha256`, set segment 4 to COMPLETE
in `data/scenarios/opening-to-whelk.json` with the assertion hashes,
and record milestone 05 in the scenario record. If any run diverges,
record the earliest divergent leg — do not hand-correct and count it.
