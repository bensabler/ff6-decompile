# EXP-0010: Decode the formula body at `ROMCPU:$C20B83`

- **Status:** running (2026-07-30, daytime session)
- **Question (#21):** What does the routine at `ROMCPU:$C20B83` (called
  per target from the `$C23469` loop, flowing into the gate region that
  queues DP `$F0`) actually compute? Bracketed body:
  `$C20B83`–`$C20C2C`.
- **Starting state:** deterministic ROM read; state-independent.
- **Observation method:** bridge `read cpu C20B83 141`
  (`$C20B83`–`$C20C0F`, joining the already-decoded `$C20C10`+ dump);
  artifact `mesen/out/rom_C20B83_141.hex` + SHA-256; hand-disassembly
  with branch verification against the known continuation.
- **Expected outcomes:**
  - *Supports the formula hypothesis:* arithmetic producing DP `$F0`
    (multiplications via the SNES ALU registers, stat reads from the
    battle arrays, randomness) ahead of the known gate block.
  - *Alternative:* `$C20B83` is only a dispatcher and the arithmetic
    lives deeper (then its JSR/JMP targets become the next bracket).
- **Falsifying outcome:** the body neither writes the DP `$F0` family
  nor transfers control toward `$C20C0B`–`$C20C2C` (would invalidate
  the bracketing).
- **Raw evidence paths:** `mesen/out/rom_C20B83_141.hex`.
- **Result:** **Decoded — the elemental-modifier block** (artifact
  `rom_C20B83_141.hex`, SHA-256 `3a9034d3…c0a39f`; joins the EXP-0006
  `$C20C10` dump seamlessly). Verified listing:

  ```asm
  C20B83  PHP             ; the pushed-PS byte EXP-0008 predicted
  C20B84  SEP #$20
  C20B86  LDA $11A6
  C20B89  BNE $C20B8E
  C20B8B  JMP $0C2B       ; gate: $11A6==0 -> PLP/RTS
  C20B8E  LDA $11A4
  C20B91  BMI $C20B98
  C20B93  JSR $0C9E       ; base amount, variant A ($C20C9E, undumped)
  C20B96  BRA $C20B9B
  C20B98  JSR $0D87       ; base amount, variant B ($C20D87, undumped)
  C20B9B  STZ $F2
  C20B9D  LDA $3EE4,Y
  C20BA0  ASL
  C20BA1  BMI $C20BFA     ; target status bit 6 -> zero amount
  C20BA3  LDA $11A4
  C20BA6  STA $F2         ; polarity/mode byte (the $20 seen live)
  C20BA8  LDA $11A2
  C20BAB  BIT #$08
  C20BAD  BEQ $C20BD3
  C20BAF  LDA $3C95,Y
  C20BB2  BPL $C20BBF
  C20BB4  LDA $11AA
  C20BB7  BIT #$82
  C20BB9  BNE $C20C2B     ; abort -> PLP/RTS
  C20BBB  STZ $F2
  C20BBD  BRA $C20BC6
  C20BBF  LDA $3EE4,Y
  C20BC2  BIT #$02
  C20BC4  BEQ $C20BD3
  C20BC6  LDA $11A4
  C20BC9  BIT #$02
  C20BCB  BEQ $C20BD3
  C20BCD  LDA $F2
  C20BCF  EOR #$01        ; flip damage<->heal (status-driven)
  C20BD1  STA $F2
  C20BD3  LDA $11A1       ; attack element byte
  C20BD6  BEQ $C20C1E     ; no elements -> skip block
  C20BD8  LDA $3EC8
  C20BDB  EOR #$FF
  C20BDD  AND $11A1
  C20BE0  BEQ $C20BFA     ; all elements nullified -> zero
  C20BE2  LDA $3BCC,Y     ; per-target mask (absorb candidate)
  C20BE5  BIT $11A1
  C20BE8  BEQ $C20BF2
  C20BEA  LDA $F2
  C20BEC  EOR #$01        ; -> flip to heal
  C20BEE  STA $F2
  C20BF0  BRA $C20C1E
  C20BF2  LDA $3BCD,Y     ; (immune candidate)
  C20BF5  BIT $11A1
  C20BF8  BEQ $C20C00
  C20BFA  STZ $F0
  C20BFC  STZ $F1         ; -> amount = 0
  C20BFE  BRA $C20C1E
  C20C00  LDA $3BE1,Y     ; (resist candidate)
  C20C03  BIT $11A1
  C20C06  BEQ $C20C0E
  C20C08  LSR $F1
  C20C0A  ROR $F0         ; -> halve
  C20C0C  BRA $C20C1E
  C20C0E  LDA $3BE0,Y     ; (weak candidate) — joins EXP-0006 dump:
  C20C11  BIT $11A1
  C20C14  BEQ $C20C1E
  C20C16  LDA $F1
  C20C18  BMI $C20C1E     ; overflow guard (amount >= $8000)
  C20C1A  ASL $F0
  C20C1C  ROL $F1         ; -> double
  C20C1E  ...             ; falls into the decoded gate block -> queue
  ```

  **Structural bonus:** `$3BCC + $14 = $3BE0` — two more 10-entry
  family arrays; each slot's 16-bit entries pack (absorb|immune) and
  (weak|resist) element masks, read as separate low/high bytes.
- **Status:** completed (2026-07-30)
- **Confidence:** All transforms — **Confirmed (byte-exact, two joined
  dumps)**: first-match-wins order flip → zero → halve → double
  (with `$8000` guard), battle-wide nullify via `~$3EC8 & $11A1`,
  status-driven polarity flips. Semantic labels (element, absorb/
  immune/resist/weak, physical/magical bases) — Strong hypothesis from
  arithmetic shape. Base-amount routines `$C20C9E`/`$C20D87` — Unknown
  (question #22).
- **Next action:** implement the element-response transform in Go with
  behavior-derived names; open question #22 (base-amount formulas).
