# EXP-0033: Golden route segment 3 — Guard battle victory to free field movement

- **Status:** completed — **partial** (2026-07-31): victory and reward
  processing captured; milestone `04-free-movement` **not** reached
  (the opening contains further scripted battles; see Result).
- **Question:** can the schedule be extended through the first
  scripted battle to victory and on to the first stable
  player-controlled field state (milestone `04-free-movement`), and
  what does the victory/reward window write (the battle-end and
  reward processing that B07/B19 both depend on)?
- **Starting state:** the EXP-0032 route (power-on, controlled lab:
  AllZeros RAM, virgin SRAM, testrunner, frame-scheduled input),
  extended past battle entry at frame 31 557.
- **Method:**
  1. *Battle input:* continue the A cadence through the battle. The
     opening party's default command (Fight/MagiTek top entry) is
     confirmable with A; EXP-0021 established that frame-exact A
     schedules produce deterministic action content, so a fixed
     cadence is the right instrument. Cadence and window are encoded
     in `mesen/probes/EXP-0033.lua`.
  2. *Battle-end detection (state-driven):* watch the Confirmed enemy
     HP array `WRAM:+$3BF4` (EXP-0028 field map: monster record +$08
     → `$3BF4,Y`) for slot 0 reaching zero, and log the frames after.
     Also first-capture-per-PC write-watch a reward-candidate window
     during that period to see what victory processing touches
     (bounded: log only, no decoding this unit).
  3. *Milestone `04-free-movement`:* after battle end, resume the
     walk+A cadence; capture the milestone at a fixed offset past
     battle end, with WRAM dump + screenshot + savestate.
  4. *Determinism:* two fresh power-on runs; byte-compare milestone
     WRAM (the valid channel per CEN-QUIRK-0002) and screenshots.
- **Expected outcomes:**
  - *Supports:* the battle ends in victory at the same frame in both
    runs; milestone 04 WRAM byte-identical; a reward-window writer
    set is captured for later decoding.
  - *Refutes:* the battle does not end within budget (cadence
    insufficient — record and redesign), or the party is defeated
    (recorded as a route failure; the defeat flow is CEN-BATTLE-0007
    and out of scope here), or runs diverge.
- **Falsifying outcome:** the fixed schedule yields different battle
  outcomes or different milestone-04 WRAM across two power-on runs.
- **Required evidence:** probe log (input transcript, battle-end
  frame, writer captures), milestone WRAM/PNG/state + hashes, beat
  screenshots — under `local_artifacts/experiments/EXP-0033/` and
  `local_artifacts/scenarios/SCN-0001/04-free-movement/`.
- **Stopping condition:** milestone 04 asserted across two runs; or
  three failed cadence designs; or a defeat outcome (record, stop,
  redesign).
- **Bounds:** no decoding of reward or battle-end routines in this
  unit — capture only. No progression past the first stable
  free-movement state.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0033/`,
  `local_artifacts/scenarios/SCN-0001/04-free-movement/`.
- **Result:** the battle half of the question is answered; the
  free-movement half is **not** — recorded as partial, not papered
  over.
  - **Battle 1 fought and won on schedule.** With an A cadence every
    90 frames from battle+60, both enemy slots reach 0 HP at **frame
    32 706** (entry 31 557 → 1 149 frames of battle). Party finishes
    at 55 / 68 / 70 (`?????` took damage from 63).
  - **Enemy slot identification (prerequisite finding).** v1 watched
    `$3BF4` (slot 0) for the death signal and never fired: **slot 0
    is a party member.** From the milestone-03 dump, slots 0-2 are
    the party (63/68/70) and **slots 6-7 are the two enemies** —
    HP 40 / MP 15 each, matching **monster record 0**'s `+$08`/`+$0A`
    (`ROMFILE:0x0F0000`). This is what corrected EXP-0032's
    monster-identity claim (see that record's correction section).
  - **Victory sequence observed (on-screen, in order):**
    `Got 32 Exp. point(s)` → `Got 96 GP` → screen transition → **a
    second scripted battle** (two quadruped enemies). The reward
    windows are themselves input-waiting.
  - **Reward processing writes into the field character block**
    (`~WRAM:+$1600`, the block EXP-0027 located) — first-hit capture
    during the post-victory window only:

    | Writer | Address | Value |
    |---|---|---|
    | `ROMCPU:$C2626D` | `+$1613` | 0 |
    | `ROMCPU:$C26274` | `+$1611/+$1612` | $80 / 0 |
    | `ROMCPU:$C24968` | `+$1614/+$1615` | 8 / 0 |
    | `ROMCPU:$C2496E` | `+$1609/+$160A` | 55 / 0 |
    | `ROMCPU:$C24979` | `+$160D/+$160E` | 24 / 0 |
    | `ROMCPU:$C20F04` | `+$1610` | 0 |
    | `ROMCPU:$C20F15` | `+$160C` | 0 |

    `+$1609` = 55 and `+$160D` = 24 are exactly the on-screen
    post-battle HP and MP of the first party member, and are the same
    two offsets EXP-0027 identified as field current HP / current MP.
    So **battle results are written back into the field character
    records at battle end** (`$C2496E`/`$C24979` are the writeback
    stores). The other offsets are reward-window-adjacent but
    unidentified — Unknown, not guessed.
  - **Free movement not reached.** At end+600 the capture landed on
    the `Got 32 Exp.` window; at end+3000 it landed **inside the
    second scripted battle**. The opening therefore contains more
    than one scripted battle before player-controlled field movement,
    so a fixed offset cannot define milestone 04. The captured state
    was demoted out of the milestone directory (kept as
    `local_artifacts/experiments/EXP-0033/postvictory-*`);
    `04-free-movement/` is deliberately empty.
- **Confidence:** battle-end frame and outcome — Confirmed (identical
  detector, deterministic schedule; single run this unit — the
  two-run determinism check was **not** performed for segment 3 and
  is carried forward). Enemy slots 6/7 = monster record 0 —
  Confirmed (two independent record fields). Reward values 32 EXP /
  96 GP — Confirmed (on-screen). `$C2496E`/`$C24979` as the field
  HP/MP writeback — Strong hypothesis (value + offset match on one
  character, one battle). Remaining `+$16xx` writers — Unknown.
- **Next action:** EXP-0034 — segment 3b: replace the fixed
  post-victory offset with a **state-driven free-movement detector**
  (candidate: re-arm the battle-init watch and continue the walk+A
  cadence through each scripted battle until N battles have ended
  and no battle re-arms for M frames), then capture milestone
  `04-free-movement` and run the deferred two-run determinism check
  for segments 3–3b together.
