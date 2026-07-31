# EXP-0018: Find the action-setup layer and the RNG state

- **Status:** running (2026-07-30)
- **Question (#29):** What writes `WRAM:+$11A6` (attack power) per
  action — the setup/hit-roll layer — and what does it read that
  varies with timing (the RNG state)?
- **Starting state:** Mesen running; `exp10-battle.mss` identical-state
  trial asset.
- **Method:** `eval`-inject a `+$11A6` write watch (dlog pattern,
  first-per-PC with stacks, plus per-write value list capped 20). Two
  trials from the same state with waits 1.2 s vs 4.5 s (EXP-0016
  protocol). Compare writer PCs, values, and frames. Injected code in
  the `exp11.log` transcript.
- **Expected outcomes:** a small writer-PC set for `+$11A6`; the
  setup routine's identity from stacks; divergent values/timing
  between trials at the same writer → that routine's reads contain
  the RNG state (next bracket for dumping).
- **Falsifying outcome:** `+$11A6` written identically in both trials
  even when outcomes diverge (then the hit roll gates elsewhere, e.g.
  a separate hit flag; re-aim at `+$11A1`-block writers).
- **Raw evidence paths:** SETUP-* lines in `events.log`,
  `mesen/out/exp11.log`.
- **Result:** (SETUP-WRITER lines in `events.log`; per-trial value
  lists in `exp11.log`)
  - **Two writers:** `ROMCPU:$C2297D` — zeroes `+$11A6` as part of an
    action-block clear before every action (register shape suggests a
    block/indexed clear; deep call chain); `ROMCPU:$C229D4` — the
    **power-populate** step, X/Y = attacker slot×2 (entries 4 and 5:
    the two enemies), writing powers `$0D` (13) and `$13` (19).
  - **Populate values are identical across trials** — enemy powers are
    deterministic. The trials diverge in **which clear events are
    followed by a populate**: Trial A had clear-without-populate at
    f210063/f210112 before its first populated action (f210163);
    Trial B populated at its first action (f210162). This reproduces
    EXP-0016's miss-vs-hit divergence exactly.
  - **Conclusion: a "miss" = cleared-but-never-populated action block
    (base computes from power 0). The RNG-consuming hit/act gate sits
    between `$C2297D` and `$C229D4`.** The populate stack's first
    return is `$319F` → `JSR` at `~ROMCPU:$C2319D` — the setup
    routine's caller bracket.
- **Status:** completed (2026-07-30)
- **Confidence:** Writer identification, populate determinism, and the
  clear-vs-populate divergence — **Confirmed** (identical-state
  trials). "Gate between clear and populate consumes the RNG" —
  Strong hypothesis (the divergence lives there; the specific read
  still unlocated). Question #29 narrowed to a dump bracket.
- **Next action:** EXP-0019 — dump `ROMCPU:$C22950`–`$C229E0` (the
  clear/populate routine) and `$C23190`–`$C231C0` (the `$C2319D`
  caller); enumerate reads; the timing-varying one is the RNG state.
