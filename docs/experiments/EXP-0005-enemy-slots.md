# EXP-0005: Do enemy HP values live in entries 4–9 of the battle arrays?

- **Status:** running (2026-07-29, overnight session)
- **Question:** H-BATTLE-0008 — the `$14`-stride array family implies 10
  entries. Are entries 4–9 (`WRAM:+$3BFC`–`+$3C07` of the current-HP
  array) the enemy slots? Discriminator: a delta-engine store
  (`ROMCPU:$C21347`/`$C21338`/`$C21396`) writing that range with `Y ≥ 8`.
- **Starting state:** reload `mesen/out/checkpoint3-mines.mss`; walk to a
  random encounter; have VICKS attack (periodic A presses) so enemies
  take damage.
- **Observation method — exact injected code (one `eval` line; `dlog`
  from EXP-0004 already resident):**

  ```lua
  _G.eseen={} _G.eRef=emu.addMemoryCallback(function(a,v) local c=emu.getCpuState(emu.cpuType.snes) local pc=(c.k<<16)|c.pc local k=string.format("%06X",pc) if eseen[k] then eseen[k]=eseen[k]+1 return end eseen[k]=1 dlog(string.format("ENEMY-HP-WRITER a=%06X v=%s", a, tostring(v))) end, emu.callbackType.write, 0x3BFC, 0x3C07, emu.cpuType.snes, emu.memType.snesWorkRam) return "enemy watch ok"
  ```

  Plus periodic `read wram 3BF4 20` (whole 10-entry current array) and a
  final `eseen` dump + screenshot. Transcript: `mesen/out/exp6.log`.
- **Expected outcomes:**
  - *Supports:* init writes populate `+$3BFC`+ with plausible enemy HP;
    a `$C2134A`-style capture (store `$C21347`, `Y=8/10/…`) fires when
    VICKS's attack lands; the array values at those entries decrease.
  - *Refutes:* enemies visibly take damage while `+$3BFC`–`+$3C07` stays
    untouched (enemy HP elsewhere; the 10-entry family is party-padded
    or differently used).
- **Falsifying outcome:** attack lands (enemy flinches/dies on screen)
  with zero writes in the watched range.
- **Raw evidence paths:** `mesen/out/exp6.log`, labeled `events.log`
  lines, `mesen/out/exp6-battle.png`.
- **Result:** **Confirmed** (transcript `mesen/out/exp6.log`, labeled
  captures in `events.log`, screenshot `exp6-battle.png`):
  - Post-init array read (`+$3BF4`, 20 bytes): party `0,0,$3B,0`
    (VICKS 59), **enemy entries 4–5 = `$0018`/`$0023`** (24, 35),
    entries 6–9 zero (two-enemy encounter).
  - After ten A-presses: both enemy entries **zeroed** (battle won);
    VICKS 59→46 from counterattacks.
  - Writer PCs on `+$3BFC`–`+$3C07`: **`$C2134A` ×4 (the delta-engine
    damage store `$C21347`) and `$C21399` ×4 (the death-handler zero
    store `$C21396`)** — the same engine and death path operate the
    enemy entries; `$C206BC`/`$C206BF` ×4 (post-death zeroing pair, as
    seen for party slot 1 in Session 003); `$C223F6` ×12 (battle-init);
    `$C22CCE` ×4 (new writer, unexplored — victory/cleanup candidate).
- **Status:** completed (2026-07-29)
- **Confidence:** Enemy slots 4–5 in the unified arrays — **Confirmed**
  (delta-engine stores + death handler observed on them; values match
  on-screen defeat). 10-entry family width — Confirmed by base-address
  arithmetic (`$14` stride is exact); entries 6–9 usage — presumed for
  larger encounters, not yet observed. H-BATTLE-0008 resolved.
- **Next action:** extend the Go model to the unified 10-slot arrays
  (engine is slot-uniform); investigate `$C22CCE`; larger-encounter
  observation for entries 6–9 someday.
