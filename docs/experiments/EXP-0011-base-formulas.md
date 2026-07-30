# EXP-0011: Decode the base-amount formulas

- **Status:** running (2026-07-30, daytime session)
- **Question (#22):** What do the base-amount routines at
  `ROMCPU:$C20C9E` and `ROMCPU:$C20D87` compute (selected by `+$11A4`
  bit 7; physical/magical candidates)? These produce the DP `$F0` value
  that the elemental block then modifies — the innermost damage formula.
- **Starting state:** deterministic ROM reads; state-independent.
- **Observation method:** bridge `read cpu C20C9E 233` (`$C20C9E`–
  `$C20D86`, variant A) and `read cpu C20D87 128` (`$C20D87`–`$C20E06`,
  variant B start); artifacts `mesen/out/rom_C20C9E_233.hex`,
  `rom_C20D87_128.hex` + SHA-256s; hand-disassembly with branch
  verification. SNES ALU multiply/divide register usage
  (`$4202`-family) and randomness sources are the tell-tales to find.
- **Expected outcomes:**
  - *Supports:* stat reads from the battle-array family or staging DP
    (`$F4`–`$FC` from `$C20420`+), multiplications, a random component,
    result left in DP `$F0`/`$F1`.
  - *Alternative:* another dispatcher layer (re-bracket deeper).
- **Falsifying outcome:** neither routine writes the DP `$F0` family
  (would invalidate the base-amount reading of the `$C20B93`/`$C20B98`
  calls).
- **Raw evidence paths:** the two `.hex` artifacts above.
- **Result:** (artifacts `rom_C20C9E_233.hex` SHA-256 `c8c895a8…f6a83c`,
  `rom_C20D87_128.hex` `95e1f088…243cfc`; alignment anchored by the PHP
  entries, the `JMP $0D3B` landing exactly on a `PLP/RTS`, and the
  `$C20D3D` helper boundary)
  - **Variant A (`$C20C9E`, standard damage):** `REP #$20; LDA $11B0;
    STA $F0` — **the base amount arrives precomputed in `+$11B0`**.
    Then, gated by `+$3414`: defense scaling via multiplier operand DP
    `$E8` + `JSR $0D3D` (→ `F0 = $C247B7(F0)+1`), defense taken from
    the 16-bit pair at `+$3BB8,Y` (`$3BB9` vs `$3BB8` by carry —
    physical/magical defense candidates), `$FF` = no defense, value
    inverted `EOR #$FF` ((255−def)/256 shape); a `#$AA` (≈2/3)
    multiplier under a `+$3EF8,Y`-bit condition; conditional halvings
    from `+$3AA1,Y` bits 1/5 and `+$3EF9,Y` bit 3; **party-vs-party
    halving** (`$11A4` bit 0 clear ∧ `Y<8` ∧ `X<8` → `LSR $F0`);
    final `F0 = $C2370B(F0)` (unexplored — randomness candidate).
  - **Boost sibling (`$C20D4A`):** modifies `+$11B0` in place
    (+≈50% shape, `$FFFF` clamp) under `+$3C44,X` attacker flags —
    `$3C44 = $3C30 + $14`, another family array.
  - **Variant B (`$C20D87`, fraction-of-HP):** basis =
    `+$3BF4,Y − pending` (floored 0) or `+$3C1C,Y` (selected by
    `$11A2` bit 2), scaled via `$E8 = +$11A6` through helper
    `$C20DCB`/`$C247B7`, **result minimum 1**, into `$F0`.
  - **MP retarget (`$C20DDD` helper):** `$11A3` bit 7 → `X += $14` and
    `Y += $14` — every subsequent `abs,Y` family read hits the next
    array over (`$3BF4,Y` → `$3C08` at the same slot): **HP→MP
    targeting as pure index arithmetic** (Rasp/Osmose-family
    candidate).
  - **Family census:** `$14`-stride slot arrays now identified at
    `+$33D0, +$33E4, +$3BB8, +$3BCC, +$3BE0, +$3BF4, +$3C08, +$3C1C,
    +$3C30, +$3C44, +$3EE4, +$3EF8` — twelve members.
- **Status:** completed (2026-07-30)
- **Confidence:** Decoded stretches — **Confirmed (byte-exact)** with
  strong internal anchors. Helpers `$C247B7` (multiplier), `$C2370B`
  (final transform), `$C24B5A`, `$C20DCB` interior — Unknown. Semantic
  labels (defense, row/defend, boost, MP-target) — Strong hypothesis
  from arithmetic shape. **No Go implementation this unit** — the
  arithmetic routes through unverified helpers (evidence threshold not
  met).
- **Next action:** questions #23 (what computes `+$11B0` — battle
  power/level layer), #24 (`$C2370B` — randomness candidate), #25
  (`$C247B7` multiplier semantics).
