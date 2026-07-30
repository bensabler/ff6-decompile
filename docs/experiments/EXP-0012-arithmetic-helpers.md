# EXP-0012: Decode the multiplier and final-transform helpers

- **Status:** running (2026-07-30, daytime session)
- **Questions (#24, #25):** What do `ROMCPU:$C247B7` (called with DP
  `$E8` operands `255−def` and `$AA` — multiplier candidate) and
  `ROMCPU:$C2370B` (final transform of `$F0` — randomness candidate)
  compute?
- **Starting state:** deterministic ROM reads; state-independent.
- **Observation method:** bridge `read cpu C247B0 80` and
  `read cpu C23700 96`; artifacts `mesen/out/rom_C247B0_80.hex`,
  `rom_C23700_96.hex` + SHA-256s; hand-disassembly. Tell-tales: SNES
  ALU registers (`$4202`/`$4203`/`$4216`) for multiply; an LFSR/table
  walk for randomness.
- **Expected outcomes:** multiplier = `A × $E8 / 256`-shape;
  final transform = bounded random variance of `$F0`. *Alternatives:*
  different arithmetic — record as found.
- **Falsifying outcome:** `$C247B7`/`$C2370B` fall mid-instruction
  (bracket error).
- **Raw evidence paths:** the two `.hex` artifacts.
- **Result:** (artifacts `rom_C247B0_80.hex` SHA-256 `94a4f755…32107c`,
  `rom_C23700_96.hex` `24ed07f7…1091a4`)
  - **`$C247B7` (multiplier wrapper):** `PHP; SEP #$20; STZ $EA;
    STA $E9; LDA $E8; JSR $4781; REP #$21; STA $EC; LDA $E8;
    JSR $4781; STA $E8; LDA $EC; ADC $E9; STA $E9; PLP; RTS` — a
    two-call composition over the core multiply at **`ROMCPU:$C24781`**
    (undumped; the code at `$C247AC` reads the SNES ALU divide result
    `$4214` long, so the ALU cluster is adjacent). Exact byte-combine
    semantics pend `$C24781` — question #25 continues one hop deeper.
  - **`$C2370B` (final transform): randomness hypothesis REFUTED.**
    Verified: `PHY; LDY $BC; BEQ done; …bypass test…; loop: A += A>>1`
    (`$FFFF` clamp) `× $BC times; STZ-equivalent $BC; PLY; RTS` — a
    **×1.5-per-count chain**, count in DP `$BC`, which the `$C23442`
    target loop increments per target under a `+$3A54` gate; count
    consumed on use. A `($B2<<1) & $11A1`-derived carry test can
    bypass the chain.
  - Bonus: `$C247D6` = 8-bit `A += A/2` clamped `$FF` (the 1.5× shape
    at byte width, stack-operand variant).
- **Status:** completed (2026-07-30)
- **Confidence:** Both decoded stretches — **Confirmed (byte-exact)**.
  `$C2370B`-as-randomness — **Refuted**; its gameplay meaning
  (multi-target/chain boost candidate) — Tentative. Core multiply
  `$C24781` — Unknown (#25 narrowed). Randomness location — folded
  into #23/#25.
- **Next action:** dump `$C24781`± (completes #25 and unlocks honest Go
  for the variant-A defense multiplier); then the `+$11B0` producer
  (#23).
