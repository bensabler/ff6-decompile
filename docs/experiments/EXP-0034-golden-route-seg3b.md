# EXP-0034: Golden route segment 3b — scripted battle chain to free movement

- **Status:** completed (2026-07-31)
- **Question:** how many scripted battles does the opening run before
  returning control to the player, which formation does each stage,
  and where is the first stable free-movement state (milestone
  `04-free-movement`)?
- **Starting state:** the EXP-0033 route (controlled lab: AllZeros
  RAM, virgin SRAM, testrunner, frame-scheduled input), with the
  battle detector made re-armable.
- **Method:** probe `mesen/probes/EXP-0034.lua`:
  1. Re-arm the Confirmed battle-init detector (`+$3B18` write from
     `$C22800-$C22FFF`) after each battle ends; log each battle's
     entry frame and `WRAM:+$11E0` formation id.
  2. Battle-end detection generalized: instead of hardcoded enemy
     slots, treat a battle as ended when the party HP words return
     to being the only non-zero HP entries — concretely, when every
     HP word in slots 3..9 reads zero for a sustained window after
     that battle's entry. (EXP-0033 established slots 0-2 = party,
     6-7 = enemies for formation 2; other formations may use other
     slots, so the detector must not assume 6/7.)
  3. Input: the walk+A cadence outside battle, the A confirm cadence
     inside battle, switched by the detector's state.
  4. Free movement = no battle re-arms for **1200 frames** after the
     last battle end AND the walk cadence is running. Capture
     milestone `04-free-movement` there.
  5. Two full power-on runs; byte-compare milestone 04 WRAM (clears
     both this unit's and EXP-0033's deferred determinism check).
- **Expected outcomes:**
  - *Supports:* a definite battle count with per-battle formation
    ids; milestone 04 captured and byte-identical across runs.
  - *Refutes:* battles never stop within budget (the walk cadence is
    driving encounters rather than the script — recorded, and the
    quiet-window criterion revisited), or runs diverge.
- **Falsifying outcome:** milestone-04 WRAM differs across two
  power-on runs of the identical schedule, or no quiet window of
  1200 frames occurs before the failsafe.
- **Required evidence:** per-battle entry frames + formation ids,
  battle-end frames, milestone WRAM/PNG/state + hashes, beat
  screenshots — under `local_artifacts/experiments/EXP-0034/` and
  `local_artifacts/scenarios/SCN-0001/04-free-movement/`.
- **Stopping condition:** milestone 04 asserted across two runs; or
  three failed detector designs; or a defeat outcome.
- **Bounds:** no decoding of the scripted-battle invocation, no
  extraction of the newly seen formations' records (register them,
  extract in the B13/B14 unit). Stop at milestone 04.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0034/`,
  `local_artifacts/scenarios/SCN-0001/04-free-movement/`.
- **Result:** **all three questions answered; milestone
  `04-free-movement` established and deterministic** — the falsifier
  is not met.
  - **The opening runs exactly four scripted battles** before player
    control, identical in both runs:

    | # | Entry frame | End frame | `+$11E0` | Formation | Staged monster-id bytes (record 2-7) |
    |---|---|---|---|---|---|
    | 1 | 31 557 | 32 736 | $0002 | 2 | `FF FF 00 00 FF FF` → two of id 0 |
    | 2 | 34 953 | 36 004 | $0001 | 1 | `19 19 FF FF FF FF` → two of id 25 |
    | 3 | 36 828 | 38 059 | $0002 | 2 | `FF FF 00 00 FF FF` → two of id 0 |
    | 4 | 39 500 | 40 975 | $0029 | 41 | `FF 19 00 00 FF FF` → id 25 + two of id 0 |

  - **Every staged record matches the ROM formation table
    statically** (`$CF6200 + id×15`, EXP-0030): formation 1 at
    `ROMFILE:0x0F620F` = `80 03 19 19 ff ff ff ff 57 8c 00 00 00 00
    3c`, formation 2 at `0x0F621E`, formation 41 at `0x0F6467` =
    `80 0e ff 19 00 00 ff ff 00 a7 bc 59 00 00 31` — byte-for-byte
    against the live staging in both runs. **Third, fourth, and fifth
    independent verifications of the formation table**, now across
    four scripted encounters plus the earlier random one.
  - **New monster reachable before Whelk: record id 25**
    (`ROMFILE:0x0F0320`), `+$08` = `$001B` = 27 HP, `+$0A` = `$0005`
    = 5 MP. Registered for the B14 extraction unit; not extracted
    here (bounds).
  - **Milestone `04-free-movement` captured at frame 46 375** — the
    Narshe exterior under player control, no window open.
    **Byte-identical across two power-on runs** (WRAM
    `3e26bed9…`, screenshot `292013ab…`, savestate size 150 574 both
    runs). This also **clears the determinism check deferred from
    EXP-0033**: every battle entry/end frame above is identical
    across the two runs.
  - *Detector design (recorded so it can be reused):* battle end is
    "every HP word in slots 3-9 zero, sustained 30 frames" — slot-
    agnostic, so it worked unchanged for formations using different
    slots. Free movement is "no battle re-arms for 5 400 frames".
    **v1 used 1 200 frames and captured the item-reward window ("Got
    Tonic ×1")** — the reward chain plus the transition to the next
    battle runs ~3 000 frames past battle end, so the short window
    was wrong; negative result preserved.
- **Confidence:** battle count (4), per-battle entry/end frames, and
  per-battle formation ids — **Confirmed** (identical across two
  independent power-on runs; every staged record independently
  matched against the static ROM table). Monster id 25's HP/MP —
  Confirmed as record fields (the live-slot cross-check that anchored
  id 0 was not repeated for id 25). Milestone-04 determinism —
  Confirmed. Monster **names** for ids 0 and 25 — Unknown (the
  monster name table is still unlocated); on-screen labels are
  observational only.
- **Next action:** EXP-0035 — segment 4: from milestone 04, traverse
  Narshe exterior to the mine entrance (milestone `05-mines-entry`),
  capturing the map transition (CEN-WORLD-0004's target) on the way;
  the walk cadence must become direction-scripted rather than
  "up-only" since free movement now requires real navigation.
