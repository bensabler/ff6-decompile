# EXP-0004: Who writes the pending-delta arrays, and who calls the display refresh?

- **Status:** running (2026-07-29, overnight session)
- **Primary question (#13):** What code writes the pending-delta arrays
  `WRAM:+$33E4`/`+$33D0` (per-slot 16-bit, `$FFFF` = none) consumed by the
  delta-engine fetch at `ROMCPU:$C213A7`? This is the damage/heal-formula
  layer.
- **Secondary observable (#19, declared up front):** callers of
  PartyDisplaySourceRefresh — exec watch at `ROMCPU:$C25D26` with stack
  capture. Shares the same encounter; no extra scope.
- **Starting state:** reload `mesen/out/checkpoint3-mines.mss` (mine
  tunnels; VICKS 46/70 sole survivor), walk to a random encounter, let
  enemies act ≥3 times.
- **Controlled variables:** same bridge session; watches injected before
  the encounter; first-capture-per-PC logging to bound log volume
  (repeat counts kept in memory and dumped at the end).
- **Observation method — exact injected code (three `eval` lines):**

  ```lua
  _G.dlog=function(tag) pcall(function() local c=emu.getCpuState(emu.cpuType.snes) local pc=(c.k<<16)|c.pc local s={} for i=1,24 do s[#s+1]=string.format("%02X", emu.read(0x7E0000|((c.sp+i)&0xFFFF), emu.memType.snesMemory)) end local f=io.open("/Users/benjaminsabler/Learning/GitHub/ff6-decompile/mesen/out/events.log","a") if f then f:write(string.format("%s pc=$%06X A=$%04X X=$%04X Y=$%04X SP=$%04X PS=$%02X DB=$%02X frame=%d\n    stack: %s\n", tag, pc, c.a, c.x, c.y, c.sp, c.ps, c.dbr, emu.getState().frameCount, table.concat(s," "))) f:close() end end) end return "dlog ok"
  _G.dseen={} _G.dRef=emu.addMemoryCallback(function(a,v) local c=emu.getCpuState(emu.cpuType.snes) local pc=(c.k<<16)|c.pc local k=string.format("%06X",pc) if dseen[k] then dseen[k]=dseen[k]+1 return end dseen[k]=1 dlog(string.format("DELTA-WRITER a=%06X v=%s", a, tostring(v))) end, emu.callbackType.write, 0x33D0, 0x33EB, emu.cpuType.snes, emu.memType.snesWorkRam) return "delta watch ok"
  _G.cseen=0 _G.cRef=emu.addMemoryCallback(function() cseen=cseen+1 if cseen<=5 then dlog("REFRESH-CALL") end end, emu.callbackType.exec, 0xC25D26, 0xC25D26, emu.cpuType.snes, emu.memType.snesMemory) return "refresh watch ok"
  ```

  Post-run: `eval` dumps of `dseen` counts and `cseen`; periodic
  `read wram 33D0 28`; screenshot.
- **Expected outcomes:**
  - *Primary:* a small set of writer PCs (bank `$C2` expected) whose
    stacks identify the attack-resolution call chain; delta values match
    observed HP changes (negatives as two's complement or the engine's
    sign convention).
  - *Secondary:* REFRESH-CALL stacks identify `$C25D26`'s caller(s);
    correlation of `cseen` with battle events settles event-driven vs
    sampled.
- **Falsifying outcome (primary):** no writes to `+$33D0`–`+$33EB` during
  attacks that change HP — would refute the pending-delta model of the
  fetch (`$C213A7`) despite its disassembly, forcing a re-read of the
  addressing (e.g., different DB at fetch time).
- **Raw evidence paths:** `mesen/out/events.log` (labeled lines),
  `mesen/out/exp5.log` (command transcript), screenshot
  `mesen/out/exp5-battle.png`.
- **Result:** (transcript `mesen/out/exp5.log`; labeled captures in
  `mesen/out/events.log`; screenshot `exp5-battle.png`)
  - **Writer PCs (counts):** `ROMCPU:$C20C9B` ×12 — mid-battle setter
    (first capture: `+$33D4` ← `$0004`, slot 2/VICKS, 16-bit A, deep
    JSR chain with returns `$0436`/`$0C2A`); `ROMCPU:$C2638E` ×48 and
    `$C26391` ×120 — bulk `$FFFF` sentinel sweepers (adjacent 3-byte
    stores); `$C22408` ×28 — the known battle-init writer, also
    initializing this region.
  - **Array extent discovery:** the sweeper wrote `WRAM:+$33E2` =
    `$33D0`-array entry 9 — the pending-delta arrays are **10 entries**
    (`+$33D0`–`+$33E3`, `+$33E4`–`+$33F7`), not 4. Address arithmetic
    generalizes: `+$3BF4/+$3C08/+$3C1C/+$3C30` are each `$14` apart —
    a 10-slot (4 party + 6 enemy candidate) struct-of-arrays family.
    Party slots are the first four (delta stores observed at Y∈{0,2,4};
    display copier copies Y≤6 only).
  - **Between events both arrays read all-`$FF`** (sentinels) — deltas
    are transient: set, consumed by the fetch, swept back to `$FFFF`.
  - **Secondary (#19):** `cseen=13` refresh calls in ~80 s of battle —
    event-driven confirmed. Every steady-state REFRESH-CALL stack tops
    with return `$140B`: the `JSR $069B` at `ROMCPU:$C21409` (post-fetch
    driver region of the EXP-0001 dump) calls `ROMCPU:$C2069B`, which
    must reach `$C25D26` without pushing another return (tail `JMP` or
    fallthrough — interpretation; alternative: stack-window artifact).
    First refresh fired 2 frames after the init delta writes.
  - `WRAM:+$3BF4` read `00 00 00 00 19 00 00 00` — VICKS 46→25 during
    the observation window (damage confirmed flowing through the chain).
- **Status:** completed (2026-07-29)
- **Confidence:** Writer PCs and counts — Confirmed (raw captures).
  10-entry array extent — Strong hypothesis (one entry-9 write + exact
  `$14` stride arithmetic; enemy-slot usage not yet observed directly).
  Refresh trigger chain — Strong hypothesis (consistent stacks; tail-call
  step unverified).
- **Next action:** dump the `$C20C9B` neighborhood (damage-formula layer)
  and `$C2069B` (refresh path) — next units; update structures for the
  10-slot family.
