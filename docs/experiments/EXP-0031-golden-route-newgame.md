# EXP-0031: Golden route segment 1 — power-on to New Game (SCN-0001 B01)

- **Status:** completed (2026-07-31)
- **Question:** can a fully frame-scheduled input sequence from
  power-on (no savestate) deterministically reach and confirm New Game
  selection — establishing milestone `00-new-game` of the SCN-0001
  golden route — and at what frame anchors does the title flow accept
  input?
- **Starting state:** headless Mesen testrunner from power-on. No
  savestate loads. Controlled variables:
  - ROM SHA-256 `0f51b4fc…d8d5e2` (verified local file).
  - Mesen 2.1.1 macOS x86_64, `--testrunner` mode, `FF6_OUT` set.
  - **SNES `RamPowerOnState` set to `AllZeros` for this program**
    (was `Random`; recorded here as a controlled variable — random
    power-on WRAM would make from-boot determinism claims
    unfalsifiable at the full-dump level). Original-hardware
    uninitialized-RAM behavior is a QUIRK-domain question, out of
    scope.
  - **SRAM controlled to virgin (file absent).** The prior
    `Final Fantasy III (USA).srm` (8192 B, SHA-256 `6afbcf1e…52cc7a`,
    mtime 2026-07-30 21:05, produced during lab sessions) is backed
    up at `local_artifacts/backups/` before removal.
  - Input only via Lua frame-scheduled injection; no wall-clock
    scheduling.
- **Method:**
  1. *Phase A — reconnaissance:* launch bridge-only; sample
     screenshots at increasing frame counts to map the boot flow
     (splash → title/menu presentation) and find where input is
     accepted and which presses advance to New Game confirmation.
     Screen content is observed, never assumed.
  2. *Phase B — scheduled route:* encode the press schedule as
     absolute frame numbers in probe `mesen/probes/EXP-0031.lua`
     (pre-placed cmd.txt arms it at frame ~10). At the milestone
     frame: `emu.createSavestate` → local milestone dir, full 128 KB
     WRAM binary dump, screenshot, assertion line (frame count +
     schedule echo). Every injected press is logged (the input
     transcript).
  3. *Phase C — determinism:* relaunch from power-on, same probe,
     compare the two runs' assertion dumps (WRAM SHA-256, frame
     numbers) byte-for-byte.
- **Expected outcomes:**
  - *Supports:* both runs produce identical milestone assertions →
    milestone `00-new-game` established; segment 1 of the golden
    route is reproducible.
  - *Refutes:* assertion dumps differ across runs (boot-time entropy
    exists even with controlled RAM/SRAM — itself a finding: locate
    the differing bytes before proceeding).
- **Falsifying outcome:** the same absolute-frame schedule from
  power-on yields differing WRAM dumps or differing visible screens
  at the milestone across two runs.
- **Required evidence:** input transcript (frame schedule + press
  log), milestone screenshots, WRAM dumps + SHA-256 hashes,
  events.log, savestate (all local under
  `local_artifacts/scenarios/SCN-0001/00-new-game/` and
  `local_artifacts/experiments/EXP-0031/`).
- **Stopping condition:** milestone saved + determinism verified; or
  three failed schedule attempts (record the negative and the
  observed flow); or the flow is found to require anything a headless
  run cannot provide.
- **Bounds:** stop at the first stable state after New Game is
  confirmed (a settle margin of a few hundred frames is allowed to
  find a stable anchor). No decoding of initialization writes here —
  that is the separate B01 init-capture unit. No continuation into
  the opening presentation beyond identifying the boundary.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0031/`,
  `local_artifacts/scenarios/SCN-0001/00-new-game/`.
- **Result:** **milestone `00-new-game` established; determinism
  Confirmed** — the falsifier is not met. Two independent power-on
  runs of the identical schedule produced **byte-identical** 128 KB
  WRAM dumps and screenshots at both assertion frames (milestone
  WRAM `35c76d03…`, stall WRAM `0f4369d5…`; PNGs also identical).
  - *Phase A observations (title/attract flow):*
    - Power-on → title logo + © visible ~frame 2969, faded by ~3709.
    - With no input the game runs an **attract loop** (~21 000-frame
      period), alternating an opening replay (visibly **dimmed
      palette**) with a gameplay demo; attract dialogue
      auto-advances; loop returns to the title each iteration.
    - The title **ignores long Start holds and sparse 12-frame
      presses** even though injected input verifiably reaches the
      auto-read registers (`$4219=$10` observed during a hold).
      **Rapid edge toggling (6 polls on / 6 off) of start+a across
      the title window registers** — which of the two buttons
      triggers is undiscriminated (both were pressed; not needed for
      route validity).
    - With virgin SRAM, no save-select screen appears: the real
      opening (full-palette story panels) begins directly.
    - The real opening then **auto-runs without any input** (story
      panels → ridge dialogue → snowfield march credits → overlook →
      Narshe descent) until the first **input-waiting dialogue box**
      at the scripted Narshe entry (WEDGE "on point" box), reached
      well before frame 30000 and stable indefinitely; a single A
      press advances it (verified).
    - Boot performs **uninitialized WRAM reads** (`$7E7DB2`,
      `$7E7CAD/AE`, `$7E7BF7/F8`, per Mesen's boot log) — RNG/entropy
      candidates; with AllZeros RAM they read 0 (controlled).
  - *Route (tracked in `mesen/probes/EXP-0031.lua`):* toggling window
    frames 2500–4200, milestone capture at frame 5200 (WRAM +
    screenshot + savestate via one-shot exec callback — like
    loadSavestate, `createSavestate` is only legal inside a main-CPU
    exec callback), stall assertion at frame 30000.
  - *Milestone artifacts:* `00-new-game.mss` (=run2's state,
    133 754 B; load-validated — resumes in the same story panel),
    run1/run2 WRAM dumps + PNGs + hashes.sha256.
- **Confidence:** schedule determinism at both assertion points —
  **Confirmed** (byte-identical dumps across independent power-on
  processes). Title edge-press requirement — Strong hypothesis
  (holds/sparse presses failed repeatedly; toggling succeeded once;
  which button and the exact acceptance window are unmapped).
  Attract-vs-real palette difference — observed (mechanism unknown).
  Boot uninitialized-read addresses — observed from the emulator log
  (consumers undecoded).
- **Next action:** segment 2 (EXP-0032): extend the schedule from the
  stalled Narshe-entry dialogue through the scripted approach and the
  first scripted battle → milestones `01-opening-cinematic` /
  `02-narshe-entry` / `03-first-scripted-battle`. Census: register
  the title/attract flow and the boot entropy reads.
