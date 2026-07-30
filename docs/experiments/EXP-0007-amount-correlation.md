# EXP-0007: Does the queued amount equal the displayed damage number?

- **Status:** running (2026-07-30, daytime session)
- **Question (#21 first step):** When PendingDeltaAccumulate
  (`ROMCPU:$C20C76`) runs during an attack, does the amount in DP `$F0`
  match the damage number shown on screen? Also capture DP `$F2`
  (polarity candidate) and Y (slot×2, possibly +`$14`-retargeted) per
  invocation.
- **Starting state:** reload `mesen/out/checkpoint3-mines.mss`; walk to
  a random encounter; let combat proceed (VICKS attacks via A presses,
  enemies counterattack) while logging every accumulator entry.
- **Observation method — exact injected code (one `eval` line; `dlog`
  resident):**

  ```lua
  _G.qseen=0 _G.qRef=emu.addMemoryCallback(function() qseen=qseen+1 if qseen<=12 then local c=emu.getCpuState(emu.cpuType.snes) local f0=emu.read(0x7E00F0,emu.memType.snesMemory)|(emu.read(0x7E00F1,emu.memType.snesMemory)<<8) local f2=emu.read(0x7E00F2,emu.memType.snesMemory) dlog(string.format("QUEUE f0=%04X f2=%02X", f0, f2)) end end, emu.callbackType.exec, 0xC20C76, 0xC20C76, emu.cpuType.snes, emu.memType.snesMemory) return "queue watch ok"
  ```

  (Direct-page reads assume `D=$0000`, as captured at every prior battle
  observation.) Screenshots taken right after attack resolutions so
  displayed numbers can be compared; final `qseen` dump.
- **Expected outcomes:**
  - *Supports:* logged `f0` values equal the damage numbers visible in
    the corresponding screenshots (within matching order); `Y` in the
    dlog line identifies target slot; `f2`/carry distinguishes
    damage vs heal invocations.
  - *Refutes:* `f0` differs systematically from displayed numbers
    (then the display value is computed elsewhere, e.g. from the
    consumed delta at fetch time).
- **Falsifying outcome:** a clean attack whose displayed number differs
  from every `f0` logged for that resolution.
- **Raw evidence paths:** labeled QUEUE lines in `mesen/out/events.log`,
  `mesen/out/exp8.log` transcript, `mesen/out/exp8-*.png` screenshots.
- **Result:** **Confirmed on three independent anchors** (QUEUE lines in
  `events.log`; screenshots `exp8-round1..6.png`; transcript `exp8.log`):
  - Five accumulator entries captured: amounts 9, 4, 2, 6 targeting
    `Y=$0004` (VICKS) and **346** (`$015A`) targeting `Y=$0008` (enemy
    entry 4), all with the same `$0C2A` caller return.
  - **Anchor 1 (authoritative array):** 9+4+2+6 = 21 = VICKS's exact HP
    loss (59→38 in `+$3BF4[2]`).
  - **Anchor 2 (HUD):** round-5 screenshot shows VICKS at 44 =
    59−(9+4+2), mid-sequence.
  - **Anchor 3 (popup):** round-6 screenshot captures the damage popup
    "6" — the final queued amount rendered on screen; the Were-Rat is
    gone from the roster (killed by the 346 Fire Beam, consistent with
    enemy HP ≤35 and entry 4 reading 0).
  - **Register semantics observed (5/5 consistent):** `X` = attacker
    slot×2, `Y` = target slot×2 at accumulator entry; DP `$F2` read
    `$20` on the party→enemy hit and `$00` on enemy→party hits
    (meaning unresolved — could encode side, command type, or
    element; one observation each).
- **Status:** completed (2026-07-30)
- **Confidence:** DP `$F0` = the final per-hit damage amount, equal to
  the displayed number and the applied delta — **Confirmed**.
  X/Y attacker/target roles — Strong hypothesis (small sample, fully
  consistent). `$F2` meaning — Unknown.
- **Next action:** the formula upstream of `$C20C28` computes `$F0`
  before the JSR — question #21 continues; stack hint: the frame below
  `$0C2A` reads `16 6B` (candidate deeper return `$6B16`, bank
  presumed `$C2` — verify by dumping around `ROMCPU:$C26B10`).
