# EXP-0006: How is the pending delta computed and stored?

- **Status:** running (2026-07-30, daytime session)
- **Question (#13 continuation):** Decode the routine containing the
  pending-delta store whose write-callback PC is `ROMCPU:$C20C9B`
  (expected store at `$C20C98`, the +3 convention), and its callers
  (EXP-0004 stack returns `$0436` → `JSR` at `ROMCPU:$C20434`; `$0C2A` →
  `JSR` at `ROMCPU:$C20C28`). What feeds the delta value?
- **Starting state:** deterministic ROM reads
  ([ROM_IDENTITY.md](../research/ROM_IDENTITY.md)); emulator state
  irrelevant.
- **Observation method:** bridge dumps `read cpu C20C60 128`,
  `read cpu C20420 48`, `read cpu C20C10 48`; hand-disassembly with
  arithmetic branch verification; artifacts archived as
  `mesen/out/rom_C20C60_128.hex`, `rom_C20420_48.hex`,
  `rom_C20C10_48.hex` with SHA-256s.
- **Expected outcomes:**
  - *Supports the pending-delta model:* a store to `+$33D0`-family
    around `$C20C98` fed by an arithmetic result; caller sites reveal
    the attack-resolution context.
  - *Refutes/complicates:* `$C20C98` is not a `+$33D0` store (callback
    PC convention violated) or the addresses fall mid-instruction.
- **Falsifying outcome:** no 3-byte store instruction at `$C20C98`
  targeting the pending-delta region.
- **Raw evidence paths:** the three `.hex` artifacts above.
- **Result:** (artifacts `rom_C20C60_128.hex` SHA-256 `95aa6214…73c21`,
  `rom_C20420_48.hex` `bd420802…bc4ae`, `rom_C20C10_48.hex`
  `c4c4c56233…8e2ac`)
  - **Pending-delta accumulator decoded** — routine entry
    `ROMCPU:$C20C76` (PHY/PHP; carry = mode at entry via `SEC` `$C20C6F`
    or `CLC` `$C20C75` paths): prelude `ROL / EOR $F2 / LSR` conditionally
    executes `TYA / ADC #$13 / TAY` (with carry ⇒ `Y += $14`),
    retargeting `+$33D0,Y` to `+$33E4,Y` — **one routine serves both
    pending arrays**, polarity from DP `$F2` and entry carry. Body:
    `REP #$20`; `LDA $33D0,Y`; `INC/BEQ/DEC` (`$FFFF` sentinel → 0);
    `CLC`; `ADC $F0` (**amount operand in DP `$F0`**); `BCS`→clamp;
    `CMP #$2710`; `BCC`→store; `LDA #$270F` — **accumulates and clamps
    at 9999**; `STA $33D0,Y` at `$C20C98` (callback PC `$C20C9B` −3 ✓);
    `PLP/PLY/RTS`.
  - **Caller chain:** `JSR $0C2D` at `ROMCPU:$C20C28` (pushes `$0C2A` —
    exact match to the EXP-0004 stack). `$C20C2D` reads `+$11A2`
    (`LSR` — bit 0 here, vs bit 7 in the engine dispatch), gates on
    `+$3A82`&`+$3A83` and `+$3EE4,X` `BIT #$02`, then reaches the
    accumulator. Upstream at `$C20C1A`: `ASL $F0 / ROL $F1` doubles the
    amount when `+$11A1`-related conditions hold (observed, unnamed).
  - **Correction to EXP-0004:** the supposed second caller return
    `$0436` was a misparse — those stack bytes are the accumulator's own
    pushed PHP (`$36`) and 8-bit Y (`$04`), followed by the true return
    `$0C2A`. The `$C20420` region dumped under that premise is
    attack-parameter staging into DP `$F4`–`$FC` (from `+$202E`-family
    and `+$3EE5,X`) — plausibly related, connection unproven.
- **Status:** completed (2026-07-30)
- **Confidence:** Accumulator code, 9999 cap, sentinel-init, dual-array
  retarget — **Confirmed (byte-exact)**. Caller-chain reading —
  Confirmed for the `$C20C28` JSR; gate semantics Unknown. DP `$F0`
  amount provenance — Unknown (new question).
- **Next action:** implement the accumulator in Go (confirmed,
  self-contained); new question: what computes DP `$F0` (the actual
  damage formula); explore `$C20C40`–`$C20C5F` gap and `$C2362F`.
