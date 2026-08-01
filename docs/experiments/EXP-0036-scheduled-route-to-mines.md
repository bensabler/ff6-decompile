# EXP-0036: Scheduled route from milestone 04 to the mines interior

- **Status:** in progress (2026-08-01)
- **Question:** can the complete route from milestone 04 through the
  Narshe exterior, the fifth scripted battle, the zigzag climb, and
  the mines transition be reproduced reliably from the project's
  scheduled power-on route — with no interactive correction?
- **Starting state:** power-on under the SCN-0001 controlled lab
  (Mesen 2.1.1 `--testrunner`, `FF6_OUT` set, AllZeros RAM, virgin
  SRAM), replaying the EXP-0034 schedule to milestone 04 (frame
  46 375), then handing control to a route controller.

## Route encoding

EXP-0035's eleven documented legs expand to **fifteen encoded legs**:
its legs 7-9 each combine a one-tile sidestep with a climb, which the
controller separates so every leg has exactly one direction and one
completion condition.

Advancement is **coordinate- or state-driven**, not duration-driven.
Comparisons are overshoot-tolerant (`>=` when moving right/down,
`<=` when moving left/up), because a leg that overshoots into a wall
still satisfies its intent. Each leg carries a timeout purely as a
divergence detector: on timeout the run **fails and reports the leg**,
it does not retry or correct.

Between legs the controller holds **8 neutral frames** (no input) so
the previous direction is released before the next is applied.

| Leg | Dir | Advance when | Expected end (X,Y) | Expected event |
|---|---|---|---|---|
| 1 | right | X ≥ `$1B` | (`$1B`,`$2A`) | — |
| 2 | up | Y ≤ `$27` | (`$1B`,`$27`) | — |
| 3 | right | X ≥ `$1E` | (`$1E`,`$27`) | guard dialogue triggers |
| 4 | A pulse | battle detected | — | dialogue → battle 5 |
| 5 | A pulse | battle ended | — | battle 5 fought |
| 6 | A pulse | 900 frames elapsed | (`$1E`,`$25`) | reward windows cleared |
| 7 | up | Y ≤ `$21` | (`$1E`,`$21`) | — |
| 8 | left | X ≤ `$1D` | (`$1D`,`$21`) | zigzag sidestep |
| 9 | up | Y ≤ `$1E` | (`$1D`,`$1E`) | — |
| 10 | right | X ≥ `$1E` | (`$1E`,`$1E`) | zigzag sidestep |
| 11 | up | Y ≤ `$18` | (`$1E`,`$18`) | — |
| 12 | right | X ≥ `$1F` | (`$1F`,`$18`) | zigzag sidestep |
| 13 | up | Y ≤ `$16` | (`$1F`,`$16`) | mine-shaft dialogue |
| 14 | A pulse | 600 frames elapsed | (`$1F`,`$16`) | shaft dialogue cleared |
| 15 | up | `+$1EA5` = `$0D` | (`$26`,`$1C`) | **map transition** |

Coordinate sources are the EXP-0035 candidates `WRAM:+$00AF` (X tile)
and `WRAM:+$00B0` (Y tile). No state variable is invented: `+$1EA5`
is used only as a *candidate* transition signal, and leg 15 also
records the coordinate jump independently so the leg is verifiable
even if `+$1EA5` turns out not to be a map id.

Leg advancement is **suspended while a battle is active** (the
coordinate bytes are not assumed meaningful in battle); the battle
legs use the battle detector instead.

## Fifth scripted battle

Leg 3-5 must handle it automatically. The unit records:

- the trigger coordinate and whether reaching it suffices
  (coordinate-triggered vs event-triggered);
- its `+$11E0` formation id and staged `+$3F44` record;
- whether the route must pause for reward/battle-exit processing;
- whether `+$00AF`/`+$00B0` stay stable across the battle;
- how control is restored afterwards.

Scope bound: identify the formation and make the route deterministic.
Enemy records and AI are **not** investigated here.

## `+$1EA5` test

A write callback on `WRAM:+$1EA5` logs every change with frame and PC;
the value is additionally sampled at every leg boundary. The unit
records the value during the scripted approach, on entering free-walk
Narshe, during each exterior leg, immediately before the transition,
and inside the mines.

Existing falsifier under test: **if `+$1EA5` changes at a boundary
inconsistent with map ownership, or does not change across the
exterior-to-mines transition, it is not a simple map-id byte.**
It is not promoted beyond Tentative unless the evidence supports it.
If it behaves like an event-state or map-mode value instead, that
interpretation is recorded separately rather than forced into the
map-id framing.

## Coordinate aliases

The four blocks that mirror `+$00AF`/`+$00B0` (`+$0541/42`,
`+$0543/44`, `+$0545/46`, `+$087A/7B`) are sampled at every leg
boundary and at battle entry/end. The unit answers only whether they
mirror the leader always, or diverge on turns, scripted movement,
battles, or the map transition. Ownership is **not** investigated
unless divergence affects route control.

## Expected outcomes

- *Supports:* three scheduled power-on runs each execute all fifteen
  legs with no correction and finish at (`$26`,`$1C`) inside the
  mines, with matching per-leg coordinates, battle trigger location,
  formation id, and `+$1EA5` sequence.
- *Refutes:* any run diverges — reported as the **earliest divergent
  leg** with its observed vs expected coordinates.

## Falsifiers

1. A leg times out (route not reproducible as encoded).
2. Runs disagree on any leg's end coordinates, the battle trigger
   coordinate, the formation id, or the `+$1EA5` sequence.
3. The final coordinate is not (`$26`,`$1C`), or `+$1EA5` inside the
   mines is not `$0D`.
4. (For the map-id claim, not the route) `+$1EA5` fails to change
   across the transition.

## Evidence requirements

Per-leg log lines (leg number, direction, entry/exit frame, entry/exit
coordinates, alias samples, `+$1EA5`), the `+$1EA5` write log with
frame and PC, battle entry/end frames with formation id and staged
record, milestone WRAM dump + screenshot + savestate, and
`hashes.sha256` — all under
`local_artifacts/experiments/EXP-0036/` and
`local_artifacts/scenarios/SCN-0001/05-mines-entry/`.

## Stopping condition

Three successful power-on runs (milestone 05 then created), **or** the
first reproducible leg failure (recorded with the divergent leg; the
run is not hand-corrected and counted as success), **or** three failed
route encodings.

## Bounds

No mines traversal past the settled entry state. No map-format, enemy-
record, or alias-ownership decoding. Register newly visible systems;
do not investigate them.

- **Raw evidence paths:** `local_artifacts/experiments/EXP-0036/`,
  `local_artifacts/scenarios/SCN-0001/05-mines-entry/`.
## Result

### Schedule iterations (failures preserved)

- **v1 — probe defect, not a route finding.** Battle-end detection was
  gated behind `ROUTE_START`, so after battle 1 at frame 31 557
  `inBattle` never cleared, the detector never re-armed for the
  opening's battles 2-4, and phase 0 held the battle cadence instead of
  walking. Fixed by running battle-end detection in every phase.
- **v2 — a real route finding.** Legs 1-3 executed exactly to their
  expected coordinates, then **leg 4 timed out**: the party stood on
  (`$1E`,`$27`) tapping A for 1 800 frames and no battle started.
  EXP-0035's recon had held `right` for 400 frames *into* the guard,
  whereas leg 3 released the direction the instant X reached its
  target. **The guard trigger fires on walking into it, not on
  occupying the tile.** Leg 4 was changed to hold `right` while tapping
  A, which required the route model to allow a leg that both holds a
  direction and pulses (see Implementation).

### The trigger: a transcription error in EXP-0035, then contact semantics

Two runs timed out at (`$1E`,`$27`) — standing still tapping A, then
holding `right` while tapping A. Re-reading **EXP-0035's own recon
log** (not its summary table) showed the cause: that record's condensed
leg table had **dropped an intermediate `up` step**. The recon actually
went `right → up → right → up → right`, and the guard dialogue appears
at (`$1E`,**`$25`**), one tile north of where the table implied.
EXP-0035's table has been corrected in place with the omission called
out; the corrected route inserts the missing climb as leg 4.

With the party on the correct tile, the battle triggers on **pushing
into the guard** (leg 5 holds `right` while tapping A) — 273 frames
after the leg begins. So the trigger is contact-based, and a route that
merely arrives at the tile and stops will wait forever. Both facts were
surfaced by leg timeouts naming themselves rather than by a silent
stall, which is what the timeout mechanism is for.

### Fifth scripted battle

- **Formation id 84 (`$0054`)**, staged record
  `80 0F 1B 1B 00 00 FF FF 9C 95 3C 35 00 00 30`.
- **Static cross-check passes:** ROM formation record 84 at
  `ROMFILE:0x0F66EC` is byte-identical to the live staging. This is the
  **sixth independent verification** of the EXP-0030 formation table.
- Monster-id bytes (record bytes 2-7) = `1B 1B 00 00 FF FF` →
  **two of record 27 and two of record 0**.
- **New monster reachable before Whelk: record 27**
  (`ROMFILE:0x0F0360`), `+$08` = `$0073` = **115 HP**, `+$0A` = `$001E`
  = **30 MP**. The 115 matches the `$73` enemy-HP word EXP-0035 saw
  live in this battle, so the record identification is anchored on both
  sides.
- **Trigger location:** player tile (`$1E`,`$25`), contact-triggered.
- **Control restoration:** the route pauses 900 frames after the battle
  ends (leg 7) to let the reward windows clear; position reads
  (`$1E`,`$25`) again afterwards, so control returns to the trigger
  tile and the schedule needs no other allowance.
- **Position bytes across the battle:** not stable — `1E,00` at entry
  and `00,00` at end. The route therefore suspends position-tracking
  legs while a battle owns the screen; without that, a leg targeting
  `Y <= $21` would complete spuriously on `Y = 0`.

### `+$1EA5` across the route

Every observed change is written by the **same instruction,
`ROMCPU:$C0B5B6`**:

| Frame | Change | Context |
|---|---|---|
| ~34 298 | `$00` → `$01` | during the scripted opening chain |
| ~39 090 | `$01` → `$05` | before the fourth scripted battle |
| — | holds `$05` | every Narshe exterior route leg (1-3 verified) |

The three-value progression explains EXP-0035's milestone-02 tension:
the scripted approach runs on different map values (`$00`, then `$01`)
than the walkable town (`$05`), so a visually-exterior scripted state
legitimately reads `$00`.

**But the falsifier then fired.** Leg 16 was first encoded to complete
on `+$1EA5` reaching `$0D`. It completed in **0 frames** — and the
milestone screenshot showed the party **still standing on the Narshe
exterior at the shaft mouth**, not in the mines. `+$1EA5` had already
been set to `$0D` during the *preceding* shaft-dialogue leg, while the
exterior was still displayed and the party had not moved.

That is precisely the recorded falsifier — *"if `$1EA5` changes at a
boundary inconsistent with actual map ownership … it is not a simple
map-ID byte."* So:

- **`+$1EA5` is NOT a simple "current map" byte.** It reaches the
  destination value *before* the transition is visible and before the
  player's position changes.
- The evidence-safe reading is a **map-load target / pending-map or
  event-state value**, written by `$C0B5B6` when the event *decides* on
  a map rather than when the map becomes current. This is recorded as
  the interpretation, not asserted as fact.
- Confidence is **not promoted**; if anything the "map id" framing is
  now the weaker of the two readings. Deciding between them needs the
  `$C0B5B6` write-watch against tileset/tilemap loading — still a
  separate bounded unit.

**Route consequence:** the transition must be detected by the **player
position jump** (X reaches `$26`, which never happens on the exterior),
not by the map byte. Leg 16 was re-encoded accordingly and the invalid
milestone artifacts from that run were discarded rather than kept.

### Coordinate aliases

The four mirror blocks (`+$0541/42`, `+$0543/44`, `+$0545/46`,
`+$087A/7B`) tracked the leader exactly on legs 1 and 3. On leg 2 the
fourth block lagged by one tile at the leg boundary
(`1B/27 1B/27 1B/27 1B/28`) — consistent with a **follower chain**
rather than a mirror, and it re-converged by the next boundary. During
battles all four read the same non-field values as the leader
(`26/26`, then `26/1A`, `26/11`), i.e. they are not field coordinates
there. No divergence affected route control.

**Position bytes are not field-meaningful during battle** (they read
`26,00` at battle entry and `00,00` at battle end), which is why
position-tracking legs are suspended while a battle owns the screen.

Answering the bounded alias question directly: the blocks **mirror the
leader with a one-tile lag on the trailing member during turns**, and
they re-converge within a leg. They do not diverge in a way that
affects route control, so ownership stays a separate question
(CEN-WORLD-0007).

### Implementation

`internal/scenario/route` is the tracked model of the route the probe
executes: `Leg`, `Route`, `Runner`, the advancement predicate, and
`MinesRoute()`. Tests cover leg sequencing, direction changes,
overshoot-tolerant coordinate completion, timeout naming the earliest
divergent leg, divergence that must not advance, battle interruption
and resume, battle-edge legs, elapsed/map-change legs, and validation.

Because the Lua probe and the Go model are two encodings of one route,
`probe_sync_test.go` parses the probe's `ROUTE` table and asserts it
matches `MinesRoute()` leg for leg — so the two cannot drift silently.
That guard caught nothing yet only because every probe change in this
unit was made in the same commit as its model change.

### Route status

The 16-leg encoding reached the mines interior (verified by screenshot:
the party is inside the shaft, not at its mouth) at (`$26`,`$21`),
settling to (`$26`,`$20`). EXP-0035's documented milestone tile is
(`$26`,`$1C`), which the recon reached by walking further north after
the transition, so a 17th leg (`up` until `Y <= $1C`) was added to land
on the documented coordinate.

**The three acceptance runs of the final 17-leg encoding are still
outstanding at the time of this record's last update.** Milestone 05 is
therefore **not claimed** — see Confidence.

- **Confidence:** route legs 1-16 execute uncorrected from a scheduled
  power-on run and reach the mines interior — **Confirmed for a single
  run** of that encoding. Formation 84 identity and its ROM match —
  **Confirmed**. Monster record 27's HP/MP — Confirmed as record fields
  (HP additionally anchored to a live enemy word). Guard trigger at
  (`$1E`,`$25`), contact-based — Confirmed (three schedule variants
  discriminate it). `+$1EA5` as a simple map-id byte — **weakened, not
  promoted**: it reaches the destination value before the transition is
  visible or the position changes. Determinism of the full 17-leg route
  — **not established** (fewer than the required three runs completed).
- **Next action:** run the 17-leg encoding three times from power-on and
  byte-compare milestone-05 WRAM, per the acceptance criteria; create
  milestone 05 only if all three reach (`$26`,`$1C`) inside the mines.
