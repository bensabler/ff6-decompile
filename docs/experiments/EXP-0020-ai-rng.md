# EXP-0020: The AI/selection layer's timing-varying state

- **Status:** running (2026-07-30)
- **Question (#30):** What state does the action-selection layer read
  that varies with input timing? Prime candidate: `WRAM:+$3A70`, read
  (`LDA/INC/LSR`) at `ROMCPU:$C23198` immediately before every action
  setup (`JSR $299F` at `$C2319D`).
- **Starting state:** fresh Mesen launch on **bridge v2** (first live
  validation of the probe/id protocol); identical-state trials from
  `mesen/out/exp10-battle.mss`, waits 1.2 s vs 4.5 s, no party input.
- **Method:** tracked probe `mesen/probes/EXP-0020.lua` (no eval
  one-liners): logs `+$3A70`'s value at each `$C23198` read, a
  first-per-PC census of `+$3A70`/`+$3A71` writers, and MVN-load vs
  fight-populate action markers. Evidence archived per
  `docs/research/EVIDENCE_LAYOUT.md` under
  `local_artifacts/experiments/EXP-0020/`.
- **Expected outcomes:**
  - *Supports `$3A70`-as-state:* its value at the read differs between
    trials at matched action ordinals, and correlates with MVN vs
    populate outcomes; writer census shows a frame-advanced updater.
  - *Refutes:* identical values across trials at matched ordinals →
    the varying state is elsewhere in the selection layer (next
    candidates from the `$C23190` region reads).
- **Falsifying outcome (for the candidate):** matched-ordinal values
  identical across both trials while action choices still diverge.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0020/`
  (events.log, commands.log, hashes.sha256).
- **Result:** **`WRAM:+$3A70` is refuted as the timing-varying state.**
  (Evidence archived under `local_artifacts/experiments/EXP-0020/`;
  probe `mesen/probes/EXP-0020.lua` is tracked, so the instrumentation
  is reproducible rather than reconstructed.)
  - **Matched-ordinal reads are identical across trials.** Trial A
    (wait 1.2 s) and trial B (wait 4.5 s) both read `val=$01` at the
    first action's `$C23198` and `val=$00` at the second. The candidate's
    falsifier is met.
  - **What did diverge: the number and timing of actions.** A:
    2 reads / 5 record loads / 2 populates; B: 6 / 11 / 6 over the same
    wall-clock window. Trial A also showed record loads at earlier
    frames (210073, 210122) with *no* accompanying read or populate.
  - **`+$3A70` writers (new):** `ROMCPU:$C21629` writes `$01`
    immediately before the read; `$C2328C` writes `$FF`; `$C22640`
    writes `$00`; `$C235B4` writes `+$3A71` = `$00`. Written-then-read
    within the same action — a per-action flag/counter shape, not an
    RNG stream.
  - **Reinterpretation:** the EXP-0016/0018 divergence is consistent
    with **action scheduling** (which battler acts, and when) rather
    than a random draw inside the damage path. Menu interaction timing
    plausibly shifts the ATB-style scheduler, changing which action
    resolves in a given window. This is an *interpretation*, not an
    observation: no scheduler state has been identified yet.
- **Status:** completed (2026-07-30) — negative result on the candidate,
  with the question refined.
- **Confidence:** `+$3A70` as RNG state — **Refuted** (identical
  matched-ordinal values under divergent timing). `+$3A70` writer set
  and write-then-read ordering — Confirmed (captures). Scheduling
  interpretation — Tentative hypothesis. Whether FF6 consumes an RNG at
  all in this path — still **Unknown**.
- **Bridge v2 validation (secondary objective, passed):** relative-path
  output resolution, id'd atomic responses, duplicate-id suppression
  (verified), command transcripts, and tracked probe loading all work
  live. Three defects were found and fixed during validation: an
  `ipairs`-over-leading-nil candidate list, CWD-relative probe paths,
  and probe-to-probe `dofile` resolution (now via `FF6_PROBE_DIR`).
- **Next action:** EXP-0021 — capture the **record index** each action
  loads (A/X at `ROMCPU:$C22966` entry) plus the scheduler-adjacent
  reads in the `$C23190` region across identical-state trials; if the
  index sequence differs while inputs match, the selection input is
  upstream state (scheduler/ATB), and the RNG question narrows to
  whether any draw feeds it.
