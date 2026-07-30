# EXP-0009: Is `$C23469` the formula-path caller frame?

- **Status:** running (2026-07-30, daytime session)
- **Question (#21 continuation):** Under EXP-0008's refined stack model
  (one pushed PS byte belonging to the `$C20C28`-containing routine),
  frame 2's return is `$346B`, implying a `JSR` occupying
  `ROMCPU:$C23469`–`$C2346B`. Verify by dump; if real, decode the
  surrounding routine for DP `$F0`/`$F4`–`$FC` staging (the damage
  formula or its dispatcher).
- **Starting state:** deterministic ROM read; state-independent.
- **Observation method:** bridge `read cpu C23430 128`
  (`ROMCPU:$C23430`–`$C234AF`); artifact `mesen/out/rom_C23430_128.hex`
  + SHA-256; hand-disassembly with branch verification.
- **Expected outcomes:**
  - *Supports:* a `JSR` at `$C23469` targeting the `$C20C0B`-region (or
    an intermediary), with surrounding code reading battle stats and
    writing the DP `$F0` family.
  - *Refutes:* no JSR alignment — the one-pushed-byte model is wrong;
    switch to the live fallback (full 48-byte stack + DP `$F0`–`$FF`
    snapshot at `$C20C76` during one attack).
- **Falsifying outcome:** bytes at `$C23469` do not begin a 3-byte JSR,
  or its target is unrelated to the accumulator path.
- **Raw evidence paths:** `mesen/out/rom_C23430_128.hex`.
- **Result:** **Supported.** Dump `rom_C23430_128.hex` (SHA-256
  `a83c86d3…10d9ad`): bytes `20 83 0B` at `$C23469` — a 3-byte
  `JSR $0B83` occupying `$C23469`–`$C2346B` and pushing `$346B`,
  exactly as the refined stack model predicted. Target:
  **`ROMCPU:$C20B83`** — upstream of the already-decoded
  `$C20C0B`–`$C20C2C` gate region → the formula body is bracketed to
  `$C20B83`–`$C20C2C`. Call-site context (alignment provisional beyond
  the verified JSR): a loop from `LDY #$12` stepping −2 per iteration
  (all ten slots, enemies first), each slot gated by `LDA $3018,Y` bits
  against DP `$A4` (`TRB $A4` shape), with `CPY $33F8` and
  `JSR $35E3`/`$4406`/`$387E` calls inside — consistent with an
  action-resolution **target loop** (`$A4` = target bitmask). `A` holds
  the slot's `$3018,Y` bit value at the `$C23469` call.
- **Status:** completed (2026-07-30)
- **Confidence:** The `$C23469` JSR and its `$C20B83` target —
  **Confirmed** (dump; pushed-return arithmetic matches five live
  stacks). Target-loop interpretation — Tentative hypothesis
  (provisional alignment; gate semantics unverified). `$3018,Y` as a
  per-slot bit table — Strong hypothesis (three independent code sites:
  fetch gate `$C213B0`, dispatch tail `$C21316`, this loop).
- **Next action:** EXP-0010 — dump `ROMCPU:$C20B83`–`$C20C10` (the
  bracketed formula body) and decode what computes DP `$F0`.
