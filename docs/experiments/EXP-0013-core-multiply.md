# EXP-0013: Decode the core multiply at `ROMCPU:$C24781`

- **Status:** running (2026-07-30, daytime session)
- **Question (#25):** What does `ROMCPU:$C24781` compute? (Called twice
  by the `$C247B7` wrapper with DP `$E8` multiplier operands; the code
  at `$C247AC` reads ALU `$4214`, so hardware multiply/divide usage is
  expected nearby.)
- **Starting state:** deterministic ROM read; state-independent.
- **Observation method:** bridge `read cpu C24770 80` (`$C24770`–
  `$C247BF`, covering the routine and joining the EXP-0012 dump);
  artifact `mesen/out/rom_C24770_80.hex` + SHA-256; hand-disassembly.
- **Expected outcomes:** writes to `$4202`/`$4203` (8×8 multiplicands)
  and a `$4216` product read, or a software multiply loop.
  *Alternative:* something else — record as found.
- **Falsifying outcome:** `$C24781` falls mid-instruction.
- **Raw evidence paths:** `mesen/out/rom_C24770_80.hex`.
- **Result:** (artifact `rom_C24770_80.hex`, SHA-256 `8b576282…516124`)
  - **`$C24781` decoded:** `PHP; REP #$20; STA $004202` (16-bit store →
    `$4202` = A.low as multiplicand, `$4203` = A.high as multiplier);
    `NOP ×4` (latency); `LDA $004216` (16-bit product); `PLP; RTS` —
    the SNES hardware 8×8 multiply of A's own two bytes.
  - **Wrapper semantics now exact:** `$C247B7` composes
    `E8×high(input) + floor(E8×low(input)/256)` with no carry-in
    (`REP #$21` clears C; inner `PHP/PLP` preserves it) and no 16-bit
    overflow (max 65280) — algebraically identical to
    `floor(input × E8 / 256)`. With the caller's `INC`:
    **`F0 = (F0 × (255−def))/256 + 1`** (defense `$FF` = skip), and the
    `$AA` path = `(F0 × 170)/256 + 1`.
  - Bonus: `$C24792` = hardware divide (`$4204/$4205` dividend from
    16-bit A, `$4206` divisor from X, 8 NOPs, remainder `$4216` → X,
    quotient `$4214` → A).
- **Status:** completed (2026-07-30)
- **Confidence:** All decoded — **Confirmed (byte-exact)**; the
  wrapper-equals-`floor(v×f/256)` identity is arithmetic, not
  interpretation. Question #25 closed. Unmodeled residue: the `#$C1`
  defense-override path under the `$3A82`&`$3A83` gate (condition
  semantics unknown).
- **Next action:** implement `Scale256`, `ApplyDefense`, `ChainBoost`
  in Go with table tests; then question #23 (the `+$11B0` producer)
  remains the innermost frontier.
