# EXP-0035: Golden route segment 4 — Narshe exterior to the mine entrance

- **Status:** completed — **partial** (2026-08-01): the route to the
  mines interior is fully mapped and reached; the scheduled encoding
  and two-run determinism check are **not** done (context boundary —
  carried to EXP-0036 with the leg table below).
- **Question:** what walking route leads from the milestone-04
  free-movement state to the mine entrance, and can the map
  transition into the mines be captured deterministically as
  milestone `05-mines-entry`?
- **Starting state:** the EXP-0034 route (controlled lab: AllZeros
  RAM, virgin SRAM, testrunner, frame-scheduled input). Recon uses
  `04-free-movement.mss` (frame counter matches the power-on runs);
  verification runs are fresh power-on replays.
- **Method:**
  1. *Recon (bounded, interactive):* from milestone 04, walk with
     directional presses guided by screenshots; find the path to the
     mine entrance; note direction legs and any dialogue stalls or
     triggers on the way. Walls normalize position (overshoot is
     safe), so legs can be generous holds.
  2. *Scheduled route (probe `mesen/probes/EXP-0035.lua`):* encode
     the legs as absolute frame windows appended to the EXP-0034
     schedule after frame 46 375. Keep the re-arming battle detector
     active throughout: if any battle (scripted or random) fires
     en route, the confirm cadence fights it and walking resumes —
     under a fixed schedule that path is deterministic too, and the
     log records it.
  3. *Milestone `05-mines-entry`:* captured at a fixed frame after
     the final leg, chosen in recon so the party has entered the
     mines interior and settled (wall-normalized). WRAM dump +
     screenshot + savestate.
  4. *Transition observation (registration only):* screenshot pairs
     across the door transition and the entry/exit frames from the
     scheduled runs — the inputs CEN-WORLD-0004 needs; locating the
     map-id variable is a separate experiment.
  5. *Determinism:* two fresh power-on runs; byte-compare milestone
     05 WRAM and screenshots.
- **Expected outcomes:**
  - *Supports:* a walkable route exists; milestone 05 byte-identical
    across runs (with or without an en-route random battle, as long
    as both runs match).
  - *Refutes:* the mines are gated by an uncleared event (recorded;
    the gate becomes the segment boundary instead), or runs diverge.
- **Falsifying outcome:** milestone-05 WRAM differs across two
  power-on runs of the identical schedule.
- **Required evidence:** recon leg screenshots, the leg schedule,
  transition screenshot pair, milestone WRAM/PNG/state + hashes,
  probe logs — under `local_artifacts/experiments/EXP-0035/` and
  `local_artifacts/scenarios/SCN-0001/05-mines-entry/`.
- **Stopping condition:** milestone 05 asserted across two runs; or
  the mines prove event-gated (record the gate, stop); or three
  failed route encodings.
- **Bounds:** no map-format or map-id decoding (CEN-WORLD-0004's own
  unit); no mines traversal beyond the settled entry state; register
  newly visible systems, do not investigate them.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0035/`,
  `local_artifacts/scenarios/SCN-0001/05-mines-entry/`.
- **Result:** the route exists, was walked end to end, and the mines
  interior was reached. The mines are **not** event-gated beyond one
  scripted battle. Two findings beyond the routing question.
  - **Player position bytes located.** A blind low-WRAM diff across
    single-tile moves isolated **`WRAM:+$00AF` = X tile** and
    **`WRAM:+$00B0` = Y tile**: `+$00AF` decrements by exactly 1 per
    left step (and is stable while moving vertically), `+$00B0`
    increments by 1 per down step (stable while moving
    horizontally). Aliases at `+$0541/+$0542`, `+$0543/+$0544`,
    `+$0545/+$0546` and `+$087A/+$087B` track identically (party
    followers and/or a mirror block — undecided). This made the rest
    of the recon coordinate-driven rather than screenshot-guessed.
  - **Candidate map-id byte `WRAM:+$1EA5`** — **SUPERSEDED by
    [CONTRA-0002](../contradictions/CONTRA-0002-1ea5-map-id-vs-event-flags.md):
    this byte is not a map id at all. It is byte 5 of an event-flag bit
    array based at `+$1EA0`. The reading below is preserved as recorded,
    including its unexplained tension, which the resolution explains.**
    Read across eight
    independently produced savestates it takes exactly three values,
    consistent within each context:

    | Value | States |
    |---|---|
    | `$00` | milestone 01 (snowfield), 02 (Narshe entry beat), 03 (battle) |
    | `$05` | checkpoint1, checkpoint2, milestone 04 — all Narshe exterior free-walk |
    | `$0D` | checkpoint3-mines, and this unit's mines-interior recon state |

    Note the tension worth keeping: milestone 02 is *visually* the
    Narshe exterior yet reads `$00`, not `$05` — either the scripted
    approach uses a different map record than the walkable town, or
    `+$1EA5` is not the map id. **Tentative hypothesis**, deliberately
    not upgraded. Falsifier: write-watch `+$1EA5` across the mine-entry
    transition — if its writer also drives tileset/tilemap loading it
    is the map id; if it merely correlates, discard.
  - **Route from milestone 04 (start pos X=$1A, Y=$2A):**

    > **Correction (EXP-0036).** This condensed table dropped one step:
    > the guard dialogue appears at (`$1E`,`$25`), **not** (`$1E`,`$27`).
    > This unit's own recon log has the missing `up` leg between them
    > (`right 400 → 1E 27`, `up 400 → 1E 25`, then the dialogue). Two
    > EXP-0036 runs timed out at (`$1E`,`$27`) before the omission was
    > found. The corrected sequence is below and in
    > `internal/scenario/route.MinesRoute`.

    | Leg | Input | Result |
    |---|---|---|
    | 1 | right | X $1A→$1B |
    | 2 | up | Y $2A→$27 |
    | 3 | right | X→$1E |
    | 3a | up | Y $27→$25 (**omitted from the original table**) |
    | 3b | right | pushes into the guard — **triggers dialogue** |
    | 4 | A | dialogue → **battle 5** |
    | 5 | A cadence | battle 5 won; field pos ($1E,$25) |
    | 6 | A, then up | Y $25→$21 |
    | 7 | left 1 tile, up | ($1D,$1E) |
    | 8 | right 1 tile, up | ($1E,$18) |
    | 9 | right 1 tile, up | ($1F,$16) — **mine-shaft dialogue** |
    | 10 | A ×3 | dialogue cleared |
    | 11 | up | **($26,$1C) — map transition into the mines** |

    Legs 7-9 matter: the climb is a **zigzag**, and a plain "hold up"
    stalls at Y=$21 and again at Y=$16. Blind holds in one direction
    cannot walk this route — that is why the earlier segments' up-only
    cadence had to be replaced.
  - **Battle 5 registered (a fifth scripted battle, beyond EXP-0034's
    four):** triggered by the guard-dialogue trigger at ($1E,$27) on
    the way to the mines, so it sits *after* the free-movement
    milestone rather than in the opening chain. Enemy name window
    shows two entries (a large quadruped plus guards); enemy HP words
    observed at slots 6-9 (`$73`, `$28`, `$28`). Records not
    extracted (bounds).
- **Confidence:** the route and its legs — Confirmed (walked, with
  coordinates logged at every leg). `+$00AF`/`+$00B0` as player X/Y
  tile — **Strong hypothesis** (clean per-axis single-step behaviour
  across two independent move series; the aliases are undecided and
  no producer was traced). `+$1EA5` as map id — **REFUTED by
  CONTRA-0002**; the byte is part of an event-flag array, and the
  "unexplained tension" noted above is explained by that resolution.
  Battle 5's existence —
  Confirmed (fought to victory); its formation id was **not** read
  (no probe was armed during recon).
- **Next action:** EXP-0036 — encode the leg table above as absolute
  frame windows appended to the EXP-0034 schedule (keeping the
  re-arming battle detector so battle 5 is fought and logged with its
  `+$11E0` formation id), capture milestone `05-mines-entry` at a
  settled frame after the transition, and run the two-run determinism
  check. Detect the transition by watching `+$1EA5` — which
  simultaneously tests that byte's falsifier.
