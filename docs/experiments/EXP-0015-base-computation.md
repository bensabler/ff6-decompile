# EXP-0015: Decode the base-computation routine

- **Status:** running (2026-07-30, overnight session 2)
- **Question (#23 completion):** Decode the routine containing the
  two-stage `+$11B0` writes (`$C22B7D`: 240, then `$C22B9A`: 450 with
  `+$11A6` = 60). Hypothesis from the values: stage 1 = power×4, stage
  2 adds a level/stat-scaled term (classic magic-damage shape).
- **Starting state:** deterministic ROM read; state-independent.
- **Observation method:** bridge `read cpu C22B40 192` (`$C22B40`–
  `$C22BFF`, covering both write sites and the enemy path `$C22BEC`);
  artifact `mesen/out/rom_C22B40_192.hex` + SHA-256; hand-disassembly
  (stores expected at write-PC −3).
- **Expected outcomes:**
  - *Supports:* multiply/shift arithmetic from `+$11A6` (power) and a
    per-attacker stat (level candidate) matching 240→450 for
    power=60; the enemy path reading a parallel stat source.
  - *Refutes the shape hypothesis:* different arithmetic — record as
    found.
- **Falsifying outcome:** `$C22B7A`/`$C22B97` are not store-anchored
  instructions (bracket error).
- **Raw evidence paths:** `mesen/out/rom_C22B40_192.hex`.
- **Result:** (artifacts `rom_C22B40_192.hex` SHA-256 `05f3eae9…14971e`,
  `rom_C20DD1_32.hex` `905eaccd…d0811d`; store sites verified
  programmatically at `$C22B7A`/`$C22B97`/`$C22BE9`, each +3 from the
  live capture PCs)
  - **Standard-path base formula, decoded and numerically closed:**

    ```asm
    C22B69  LDA $11AF     ; stat B (=4 live)
    C22B6C  STA $E8
    C22B6E  CMP #$01
    C22B70  TDC
    C22B71  LDA $11A6     ; power (=60 live)
    C22B74  REP #$20
    C22B76  BCC $C22B7A   ; statB == 0 -> skip the x4
    C22B78  ASL / ASL     ; power x4
    C22B7A  STA $11B0     ; stage 1 (=240 live)
    C22B7D  SEP #$20
    C22B7F  LDA $11AE     ; stat A (=28 live)
    C22B82  XBA
    C22B83  LDA $11A6
    C22B86  JSR $4781     ; power x statA (=1680)
    C22B89  JSR $47B7     ; x statB; full 24-bit product preserved in $EA:$E9:$E8 (=6720)
    C22B8C  LDA #$04
    C22B8E  REP #$20
    C22B90  JSR $0DD1     ; 24-bit product >> (A+1) = >>5 (=210)
    C22B93  CLC
    C22B94  ADC $11B0
    C22B97  STA $11B0     ; stage 2 (=450 live)  [wrapping add, no clamp]
    C22B9A  SEP #$20
    C22B9C  RTS
    ```

    **`base = power×4 + (power × $11AE × $11AF) >> 5`** — 240 + 210 =
    450 with the live-read operands (60, 28, 4). Exact.
  - **`$C20DD1` helper decoded:** `PHX; TAX; LDA $E8 (16-bit);
    loop X+1 times { LSR $EA; ROR A }; PLX; RTS` — shifts the
    wrapper's 24-bit full product right by A+1 (the `$C20DCB` entry
    seen in variant B prepends `JSR $47B7; LDA #$0003` → shift 4).
    The `$C247B7` wrapper's memory layout (`$E8` = residue byte,
    `$E9:$EA` = quotient) composes the full product naturally.
  - **Enemy/physical path** (`$C22B9D?`–`$C22BE9`, store verified):
    operands `$11AE`/`$11AF`/`$11A6` and double `$47B7` usage with a
    `LDA $B2 BIT #$4000`-gated stack-value scaling (≈×1.75 shape) —
    partially decoded, follow-up question.
- **Status:** completed (2026-07-30)
- **Confidence:** Standard-path formula and `$C20DD1` — **Confirmed
  (byte-exact + exact numeric closure against live memory)**.
  `$11AE`/`$11AF` semantic labels (stat pair; magic-power/level
  candidates) — Unknown/Tentative: recorded by address only. Physical
  path — partial. Damage variance — still not located (this path was
  deterministic).
- **Next action:** implement `BaseAmountStandard` + `Shift24` in Go
  with the live vector as a golden test; questions: #26 physical-path
  completion, #27 `$11AE`/`$11AF` producers/meaning, variance hunt.
