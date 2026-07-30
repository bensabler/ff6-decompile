# EXP-0008: Is `~$C26B13` a caller frame on the damage-formula path?

- **Status:** running (2026-07-30, daytime session)
- **Question (#21 continuation):** Every EXP-0007 accumulator stack
  shows, beneath the `$0C2A` return, the bytes `16 6B` — candidate
  return `$6B16` (bank `$C2` presumed), implying a `JSR` at
  `ROMCPU:$C26B13`. Is that a real frame of the path that computes DP
  `$F0`, and what does the surrounding code do?
- **Starting state:** deterministic ROM read; state-independent.
- **Observation method:** bridge `read cpu C26AE0 128`
  (`ROMCPU:$C26AE0`–`$C26B5F`); artifact `mesen/out/rom_C26AE0_128.hex`
  + SHA-256; hand-disassembly with branch verification.
- **Expected outcomes:**
  - *Supports:* a 3-byte `JSR` ending at `$C26B16` (i.e. at `$C26B14`…
    wait — a `JSR` whose **last byte** is `$6B16` sits at `$C26B14`;
    alternatively the byte pair is data). Precisely: bytes at
    `$C26B14`–`$C26B16` form `JSR $xxxx`, and the surrounding routine
    stages attack parameters (DP `$F0`-family or the `$F4`–`$FC`
    block).
  - *Refutes:* mid-instruction or unrelated code — the `16 6B` bytes
    were frame data, and the fallback experiment (full-DP snapshot
    logging) takes over.
- **Falsifying outcome:** no `JSR` alignment at `$C26B14` and no
  DP-`$F0`-family access in the dumped window.
- **Raw evidence paths:** `mesen/out/rom_C26AE0_128.hex`.
- **Result:** **Lead refuted as posed.** Dump `rom_C26AE0_128.hex`
  (SHA-256 `9feb9ecb…6559e2`): bytes at `$C26B12`–`$C26B14` are
  `D0 06 2F` — no 3-byte `JSR` ends at `$C26B16`, and the window shows
  no DP-`$F0`-family staging (alignment there is itself uncertain — no
  anchor instruction found; possibly data or mid-routine).
  **Refined stack model:** the routine containing the `$C20C28` JSR
  ends `PLP/RTS` (`$C20C2B`–`$C20C2C`, EXP-0006 dump), so its prologue
  pushes a PS byte — stack byte 3 (`$16` ≙ `m=0,x=1` PS, plausible) is
  that push, not part of a return. Under that model frame 2's return is
  **`$346B`** → a `JSR` at `~ROMCPU:$C23469`. The recurring deep bytes
  (`… 25 61 C1 C0 …`) suggest the chain eventually reaches a bank-`$C1`
  JSL frame (`$C1:6125` candidate) — recorded as context only.
- **Status:** completed (2026-07-30) — negative result with a refined
  lead.
- **Confidence:** Refutation of the `$C26B14` JSR — Confirmed (dump).
  `$346B`/`$C23469` frame — Tentative hypothesis (depends on the
  one-pushed-PS model; a second pushed byte would shift the frame).
- **Next action:** dump `ROMCPU:$C23430`–`$C234B0` (EXP-0009); verify a
  `JSR` occupying `$C23469`–`$C2346B`; if absent, switch to the live
  full-stack + DP-snapshot log at `$C20C76`.
