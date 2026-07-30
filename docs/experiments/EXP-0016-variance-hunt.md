# EXP-0016: Is damage timing-dependent? (variance hunt)

- **Status:** running (2026-07-30)
- **Question (#28):** The decoded standard pipeline contains no RNG.
  Does the same attack from the same savestate produce different damage
  when only input-frame timing differs? A divergence localizes the
  variance; identical results refute frame-timing-driven RNG for this
  path and push variance upstream (hit/crit rolls) or nowhere.
- **Starting state:** Mesen already running (EXP-0014 instance, `dlog`
  resident). Reload `checkpoint3-mines.mss`, walk to an encounter, then
  **save a local battle state** `mesen/out/exp10-battle.mss` via eval
  (`emu.createSavestate`), giving a reusable identical starting point.
- **Method — trials:** Trial A: load `exp10-battle.mss`, wait ~1.2 s,
  press A every ~3 s until the attack resolves. Trial B: same, but wait
  ~4.5 s first (different frame offset). Per trial, log **every**
  `+$11B0` write (VAR-BASE, cap 30) and every accumulator entry with DP
  `$F0`/`$F2` (VAR-QUEUE, cap 20); reset logs between trials via eval.
  Injected code preserved in `mesen/out/exp10.log` transcript (single
  `eval` lines, `dlog`-based, same pattern as EXP-0014).
- **Expected outcomes:**
  - *Divergence:* base or queued values differ between trials →
    variance exists and is state/timing-dependent; the diverging write
    site is the RNG consumer.
  - *Identity:* all values identical → this attack path is
    deterministic w.r.t. input timing; variance (if any) needs a
    different discriminator (e.g., repeated same-timing trials, or
    physical attacks).
- **Falsifying outcome:** for the "timing-driven RNG" hypothesis:
  identical values across trials.
- **Raw evidence paths:** VAR-* lines in `events.log`,
  `mesen/out/exp10.log`, `exp10-battle.mss` (local, gitignored).
- **Result:** **Divergence confirmed — damage/hit outcomes are
  timing-dependent.** (Transcript `exp10.log`; local state
  `exp10-battle.mss`, 127,332 bytes, created via the queued
  `createSavestate` pattern — it, like `loadSavestate`, requires a
  main-CPU exec callback.)
  - Trial A (wait ≈1.2 s): first enemy action at f210122 = **base 0
    (miss), nothing queued**; second at f210275 = base 12 → queued 3.
  - Trial B (wait ≈4.5 s): first action at f210162 = **base 7 → queued
    4**; then base 0 (miss), base 12 → 3, base 7 → **2**.
  - From the identical savestate, only input timing differed → the
    first-action divergence (miss vs hit) is clean evidence of
    **timing-consumed RNG**. Base values also vary per action (7/0/12)
    and the same base 7 queued 4 vs 2 on different actions (post-base
    modifier differences — RNG or state drift from our menu inputs;
    undiscriminated).
  - New writer PC: `ROMCPU:$C22C02` — zero-writes `+$11B0` on miss
    frames (miss path).
  - All observed actions were enemy attacks (`$C22BEC` path); our A
    presses never executed a party command this run.
- **Status:** completed (2026-07-30)
- **Confidence:** "Variance exists and is timing-dependent" —
  **Confirmed** (identical-state trials, divergent first action).
  RNG location — narrowed to the enemy/physical base path (the base
  value itself varies) with a possible second site downstream;
  specific RNG state/consumer — Unknown (question #28 continues).
- **Next action:** decode the enemy/physical base path precisely
  (#26 — the dump exists, `rom_C22B40_192.hex`, plus `$C22C02`
  region) and identify what it reads that varies with timing (the RNG
  state); then repeat-trial with a write-watch on that state.
