# EXP-0032: Golden route segment 2 — Narshe entry to the first scripted battle

- **Status:** completed (2026-07-31)
- **Question:** can the deterministic schedule be extended from the
  EXP-0031 Narshe-entry stall through the scripted approach into the
  first scripted battle — establishing milestones
  `01-opening-cinematic`, `02-narshe-entry`, and
  `03-first-scripted-battle` — and which formation id does the first
  scripted battle stage (linking B07 to the Confirmed formation
  table)?
- **Starting state:** EXP-0031's controlled lab (AllZeros RAM, virgin
  SRAM, testrunner, frame-scheduled input). Runs are fresh power-on
  replays of the segment-1 schedule plus the extension; recon may use
  `00-new-game.mss` (its frame counter matches the power-on run).
- **Method:**
  1. *Recon (bounded):* from the stall, advance dialogue with
     interactive A presses; map the beat structure (stall count,
     scene transitions, pre-battle visuals) by screenshot. No
     decoding.
  2. *Scheduled extension (probe `mesen/probes/EXP-0032.lua`):*
     segment-1 schedule + an **A metronome** (12-frame press every
     240 frames) from frame 31000, stopped by a **state-driven battle
     detector** (first write burst to the `WRAM:+$3B18` enemy-stat
     family, the Confirmed battle-init signature from EXP-0028).
     Milestone captures:
     - `01-opening-cinematic`: fixed frame 15000 (auto-run
       presentation interior; no input dependence).
     - `02-narshe-entry`: fixed frame 30000 (the EXP-0031 stall
       assertion state, re-labeled as the milestone).
     - `03-first-scripted-battle`: battle detection + 120 frames
       (state-driven, deterministic under a fixed schedule).
     At detection, log `WRAM:+$11E0` and `+$3F44..+$3F53` (staged
     formation record) — the scripted battle's formation identity.
  3. *Determinism:* two fresh power-on runs; byte-compare all
     milestone WRAM dumps and screenshots.
- **Expected outcomes:**
  - *Supports:* both runs byte-identical at every milestone; battle
    detection at the same frame; a formation id readable at staging.
  - *Refutes:* runs diverge (schedule-external entropy — locate the
    divergence before proceeding), or the metronome causes an
    unintended interaction (a choice box or menu — recorded, schedule
    revised to explicit presses).
- **Falsifying outcome:** fixed-schedule runs from power-on diverge
  at any milestone assertion, or no battle-init write burst occurs
  within the run budget (metronome fails to reach the battle).
- **Required evidence:** probe log (press transcript + detection
  frames + staged formation bytes), milestone WRAM dumps + PNGs +
  hashes, recon screenshots — all under
  `local_artifacts/experiments/EXP-0032/` and
  `local_artifacts/scenarios/SCN-0001/{01,02,03}-*/`.
- **Stopping condition:** milestones 01–03 asserted across two runs;
  or the metronome interacts with a choice/menu (record, bound out,
  redesign); or three failed schedule attempts.
- **Bounds:** stop at battle entry +120 frames. No battle commands,
  no battle-outcome capture (later segments). No event-script
  decoding. Census pass for newly visible systems only.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0032/`,
  `local_artifacts/scenarios/SCN-0001/`.
- **Result:** **milestones 01/02/03 established; the battle-entry
  falsifier is not met** — but one assertion channel is weaker than
  expected (recorded below, not hidden).
  - **Schedule (three iterations; failures preserved):**
    - *v1 (A metronome only, 240-frame period):* cleared the
      Narshe-entry box but **never reached a battle** in 60 000
      frames — the party stands still (failsafe fired). Negative
      result: the beat is not dialogue-only.
    - *v2 (single A + held Up):* the walk starts, then **stalls at a
      guard dialogue box** which persists indefinitely (frames
      32 000→48 000 identical). Negative result: walking alone is
      insufficient.
    - *v3 (held Up **and** A every 240 frames, frames 31 000-46 000):*
      reaches battle entry. Both inputs are required — the beat is a
      chain of input-waiting boxes separated by player-controlled
      walking north through the Narshe gate.
  - **Battle entry detected identically in both runs: frame 31 557**,
    first `$3B18` family write from `$C22899` (the Confirmed party-slot
    populate cluster, EXP-0028), `addr=$7E3B30 val=$82`.
  - **First scripted battle = formation id 2** (`WRAM:+$11E0` = $0002).
    Staged record at `+$3F44`: `80 0C FF FF 00 00 FF FF 00 00 4B 87
    00 00 33 10`. **Static cross-check passes:** ROM formation record
    2 (`ROMFILE:0x0F621E`, = `$CF6200 + 2x15`) reads
    `80 0C FF FF 00 00 FF FF 00 00 4B 87 00 00 33` — the first 15
    staged bytes byte-for-byte (the 16th, `$10`, is beyond the record:
    the copy loop moves 16 bytes from a 15-byte table, EXP-0030).
    Monster-id bytes are `{12}` with `$FF` empties → **a single
    enemy**, matching the on-screen single "Guard" — and monster id
    12's record is at `ROMFILE:0x0F0180`. This closes B07's formation
    identity and independently re-verifies the EXP-0030 formation
    table on a second, scripted (non-random) encounter.
  - **Determinism:** WRAM byte-identical across runs at **all three**
    milestones (01 `011588bc…`, 02 `0f4369d5…`, 03 `24302078…`);
    screenshots identical at 02 and 03.
  - **Weak channel (honest):** milestone **01's screenshot is not
    byte-stable** across runs (11 536 vs 11 339 B; same scene and
    text, sub-frame vertical difference) even though its WRAM is
    identical, and its savestate size differs (148 796 vs 148 697).
    Interpretation: frame 15 000 lands mid-presentation where
    per-scanline/HDMA effects (snowfield scroll) leave PPU-side state
    outside WRAM; the `endFrame` screenshot samples it at an unfixed
    phase. **Alternative not excluded:** genuine emulation-side
    nondeterminism in that scene. Falsifier: capture PPU/HDMA
    registers plus the frame buffer at 15 000 across runs — if the
    registers differ, the state (not the capture) is unstable.
    **Milestone 01's PNG is therefore not a valid assertion channel;
    its WRAM is.**
- **Confidence:** milestone WRAM determinism at 01/02/03 —
  **Confirmed** (byte-identical across independent power-on runs).
  Battle-entry frame and detection PC — Confirmed. First scripted
  battle = formation 2, single monster id 12 — **Confirmed** (live
  staged bytes + independent static ROM record match). Milestone 01
  visual instability cause — Tentative (two live alternatives named).
  Input-chain structure (boxes + walking) — Confirmed behaviorally
  (three-way schedule discrimination).
- **Next action:** EXP-0033 — segment 3: from milestone 03, fight the
  Guard battle to victory and continue to free field movement
  (milestone `04-free-movement`), capturing battle-end/victory
  processing on the way; separately queue the milestone-01 PPU
  falsifier and monster-record 12 extraction (B14).
