# EXP-0014: What writes the base amount at `WRAM:+$11B0`?

- **Status:** running (2026-07-30, overnight session 2)
- **Question (#23):** Base-amount variant A consumes a precomputed
  16-bit value at `WRAM:+$11B0` (EXP-0011), and the boost sibling
  modifies it in place. What code computes and writes it — the battle
  power × level/stat layer — and does the damage variance live there?
- **Starting state:** fresh Mesen launch (previous instance closed) with
  `mesen/bridge.lua`; reload `mesen/out/checkpoint3-mines.mss`; one
  random encounter via injected walking; let combat run with periodic A
  presses.
- **Observation method — exact injected code (two `eval` lines; `dlog`
  re-injected first, code identical to EXP-0004's record):**

  ```lua
  _G.dlog=function(tag) pcall(function() local c=emu.getCpuState(emu.cpuType.snes) local pc=(c.k<<16)|c.pc local s={} for i=1,24 do s[#s+1]=string.format("%02X", emu.read(0x7E0000|((c.sp+i)&0xFFFF), emu.memType.snesMemory)) end local f=io.open("/Users/benjaminsabler/Learning/GitHub/ff6-decompile/mesen/out/events.log","a") if f then f:write(string.format("%s pc=$%06X A=$%04X X=$%04X Y=$%04X SP=$%04X PS=$%02X DB=$%02X frame=%d\n    stack: %s\n", tag, pc, c.a, c.x, c.y, c.sp, c.ps, c.dbr, emu.getState().frameCount, table.concat(s," "))) f:close() end end) end return "dlog ok"
  _G.bseen={} _G.bRef=emu.addMemoryCallback(function(a,v) local c=emu.getCpuState(emu.cpuType.snes) local pc=(c.k<<16)|c.pc local k=string.format("%06X",pc) if bseen[k] then bseen[k]=bseen[k]+1 return end bseen[k]=1 dlog(string.format("BASE-WRITER a=%06X v=%s", a, tostring(v))) end, emu.callbackType.write, 0x11B0, 0x11B1, emu.cpuType.snes, emu.memType.snesWorkRam) return "base watch ok"
  ```

  Post-encounter: `bseen` dump, `read wram 11A1 18` snapshots, and a
  `QUEUE`-style comparison against EXP-0007's still-armed pattern is
  unavailable (fresh emulator) — correlation is via values only this
  run. Transcript `mesen/out/exp9.log`.
- **Expected outcomes:**
  - *Supports:* a small set of writer PCs; captured values that, after
    the decoded variant-A post-processing, match observed damage;
    writer stacks identifying the stat/power layer.
  - *Refutes the "precomputed base" model:* no `+$11B0` writes before
    accumulator activity in an attack (would mean `$11B0` is stale
    state reused, and the model needs revisiting).
- **Falsifying outcome:** attacks resolve with zero `+$11B0` writes.
- **Raw evidence paths:** BASE-WRITER lines in `mesen/out/events.log`,
  `mesen/out/exp9.log`, `mesen/out/exp9-battle.png`.
- **Result:** (BASE-WRITER lines in `events.log`; transcript
  `exp9.log`; screenshot `exp9-battle.png`)
  - **Writer census:** computation pair `ROMCPU:$C22B7D` (wrote 240,
    16-bit A `$00F0`, attacker X=`$04`/VICKS) then `$C22B9A` (updated
    to **450**, A=`$01C2`) at the same frame — a two-stage base
    computation; enemy-attack path `$C22BEC` (base 12, attacker
    X=`$0A`); per-target staging `$C23422` (target-loop preamble);
    init/clear writers `$C25561`/`$C20E8A`/`$C210E1` at battle start.
  - **Numeric closure:** `+$11B0` read back `$01C2` = 450 with
    `+$11A6` = 60; 240 = 60×4 (power×4 shape, second stage +210);
    450×(255−def)/256+1 ≈ 346 for def≈58 — matching EXP-0007's
    observed Fire Beam damage. `+$11A1` = `$01` during Fire Beam
    (fire = bit 0 candidate); `+$11A4` = `$20` (matches the live
    `$F2` polarity byte).
  - **Caller bracket:** both computation stacks return `$322D` → a
    `JSR` at `~ROMCPU:$C2322B`.
  - Variance: the two-stage computation looked deterministic in this
    single trial; the variance source remains unobserved (needs
    repeat-trial comparison — future experiment).
- **Status:** completed (2026-07-30)
- **Confidence:** Writer PCs, values, and the numeric chain closure —
  Confirmed (raw captures + arithmetic). "power×4 + level-scaled term"
  reading — Tentative hypothesis (single trial, shape only). Question
  #23 narrowed to the `$C22B40`–`$C22BFF` decode (EXP-0015).
- **Next action:** dump and decode the base-computation routine.
