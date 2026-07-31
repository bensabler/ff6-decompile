# EXP-0021: Record-index and selection-window capture across timed trials

- **Status:** completed (2026-07-30)
- **Question (#30, second refinement):** EXP-0020 left a Tentative
  scheduling interpretation: identical-state trials with different
  menu-input timing diverge only in *how many actions resolve and
  when*, not in what each action is. Test that: at **matched action
  ordinals**, is the attack-record **index sequence** identical across
  trials? And what does the selection window read immediately before
  each setup?
- **Starting state:** fresh **headless** Mesen run (`--testrunner`,
  bridge v2, no stale probe callbacks); identical-state trials from
  `mesen/out/exp10-battle.mss`. Methodological change from
  EXP-0016..0020: the GUI could not open windows in this session (the
  machine's display is locked; Avalonia RenderTimer error −6661 on
  launch, script window never created), so trials run in testrunner
  mode with **frame-scheduled input** — the EXP-0020 wall-clock cadence
  mapped to frame offsets at ≈60 fps: trial A first press at
  anchor+72 (≈1.2 s), trial B at anchor+270 (≈4.5 s), then three more
  10-frame presses at 180-frame (≈3 s) intervals in both. This removes
  the wall-clock jitter the earlier trials tolerated; count comparisons
  are read at the equalized anchor+990 marker.
- **Method:** tracked probe `mesen/probes/EXP-0021.lua`:
  - `EXP21-SEL` — exec at `ROMCPU:$C23198` (the verified
    `LDA $3A70` in the pre-setup window; the 48-byte dump
    `rom_C23190_48.hex` disassembles to `JSR $B585 / JSR $26D3 /
    LDA $3A70 / INC A / LSR A / JSR $299F`, then a power-zero gate
    `LDA $11A6 / BNE / JMP $3275` and a `DP $B5` compare `CMP #$06`):
    logs registers plus `DP $B5`, `+$3A70`, `+$3A71`.
  - `EXP21-IDX` — exec at `ROMCPU:$C2297A` (the verified
    `MVN $C4,$7E`), gated on `A=$000D` (first byte of the 14-byte
    move): X is the source offset, so the **record index =
    (X − $6AC0) / 14** — derived from the actual source pointer, not
    an assumed register convention.
  - `EXP21-POP` — exec at `ROMCPU:$C229D1` (the verified fight-populate
    store `STA $11A6`): logs A (power) and X (attacker offset).
  - `exp21_reset()` zeroes counters between trials; trial boundaries
    marked in the log via `probelog("EXP21-TRIAL-*")`.
- **Expected outcomes:**
  - *Supports scheduling interpretation:* matched-ordinal record
    indices identical across trials (first MVN action loads the same
    record in A and B, etc.); divergence remains count/timing only.
    The RNG question then collapses into the scheduler cadence (ATB
    timer layer) — its state writers become the next frontier.
  - *Refutes:* matched-ordinal indices differ → action *content*
    varies with timing; the varying selection-window value logged at
    `$C23198` (A/X/`$B5`/`$3A70`) points at the state to trace.
- **Falsifying outcome (for the scheduling interpretation):**
  matched-ordinal indices differ across trials. If they differ while
  every logged selection-window value is identical, the varying input
  was not captured (bracket error — widen the window next).
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0021/`
  (events.log SHA-256 `496e9640…5bdb8d`, commands.log `b91f3162…9e867d`,
  experiment.json, hashes.sha256). A first trial-A run was killed by
  the testrunner's default ~2-minute timeout (partial data preserved in
  the archived events.log); trials A2/B2 ran under `--timeout=7200`.
- **Result:** **matched-ordinal action content is identical across the
  trials; the falsifier is not met.** Both trials anchored at frame
  209982 (same savestate, deterministic load).
  - *Record indices:* 8 MVN loads per trial, **every one index 238**
    (`srcX=$77C4`, `(X−$6AC0)/14 = 238` remainder 0; `mvnfires=112 =
    8×14` cross-checks the per-byte exec fire). Per-stream matched
    ordinals identical, registers included.
  - *Enemy setups:* 4 SEL reads + 4 POP stores per trial, byte-identical
    at matched ordinals in payloads and registers — powers
    **13, 0, 19, 0** (`A=$FF0D/$FF00/$FF13/$FF00` at `$C229D1`),
    attackers `X=$0008/$0009/$000A`. **The miss (power 0) lands at the
    same matched ordinal in both trials** — hit/miss is deterministic
    given the input frame schedule.
  - *What input timing does change:* (1) **when** enemy setup clusters
    resolve — phase-1 cluster at +182 (A2, press at +72) vs +180 (B2,
    no press yet); phase-2 at +334 vs +403 (69-frame shift tracking the
    press-schedule difference); (2) **where the two solo record loads
    sit** — they follow the *first press* in both trials (A2: +74/+123
    after press k=1 at +72; B2: +331/+361 after press k=1 at +270) with
    no SEL/POP attached — consistent with party-action staging
    (interpretation, not observation); (3) two microscopic
    press-history residues with no outcome effect: SEL n=4 read
    `+$3A71 = $04` (A2, two presses prior) vs `$00` (B2, one press
    prior), and POP n=3 carry flag (`PS=$35` vs `$34`).
- **Confidence:**
  - Matched-ordinal identity of record index and enemy setup content
    across the two schedules — **Confirmed (byte-exact, registers
    included)**.
  - Enemy-action timing coupling to the press schedule (+2/+69-frame
    cluster shifts) — Confirmed for this window.
  - EXP-0020's scheduling interpretation — **upgraded to Confirmed for
    this state/window**: divergence is count/timing only; content is
    schedule-deterministic.
  - Attribution of the EXP-0016/0018/0020 GUI-era variance to
    harness wall-clock jitter (uncontrolled press/load frames), not to
    hidden nondeterministic state — **Strong hypothesis** (it is the
    only remaining uncontrolled variable; not directly reproduced).
  - Whether a deterministic in-engine PRNG exists inside the schedule
    (stepped per frame/action) — **Unknown**; no *nondeterminism*
    remains, so question #30's "free RNG" framing is retired.
  - `+$3A71` as press-history-coupled state read in the selection
    window — Confirmed (correlation); semantics Unknown.
- **Limitations / alternatives:** two trials, one savestate, one
  ~16.5 s window, two schedules; headless testrunner input latching is
  assumed GUI-equivalent (untested). A schedule could exist where the
  press-coupled state (`$3A71`, carry residue) flips an outcome; not
  observed here.
- **Next action:** battle path pauses at this natural boundary (operator
  rebalance). When resumed: (a) verify GUI/testrunner input-latching
  parity with one GUI trial replayed frame-exactly; (b) trace `+$3A71`
  writers during presses to name the press-coupled state; (c) probe
  whether any frame-stepped PRNG cell exists by diffing full WRAM
  between matched ordinals.
