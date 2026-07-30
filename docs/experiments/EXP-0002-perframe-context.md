# EXP-0002: Does PerFrameBattleUpdate run outside battle?

- **Status:** running (2026-07-29, overnight session)
- **Question:** Does `ROMCPU:$C101FB` (PerFrameBattleUpdate, FN-0002)
  execute at the title screen or in field navigation, or only during
  battle? Tests H-BATTLE-0003.
- **Starting state:** Mesen 2.1.1 + bridge at the title screen (post
  EXP-0001); then `mesen/out/checkpoint3-mines.mss` (believed field
  context — verified by screenshot during the run); then
  `mesen/out/checkpoint1.mss` (battle).
- **Controlled variables:** same emulator process, same injected callback
  across contexts; two samples ≥4 s apart per context so the frame delta
  is in the hundreds.
- **Observation method:** eval-injected counting exec callback —
  **exact injected code preserved:**

  ```lua
  _G.pfCount=0 _G.pfRef=emu.addMemoryCallback(function() pfCount=pfCount+1 end, emu.callbackType.exec, 0xC101FB, 0xC101FB, emu.cpuType.snes, emu.memType.snesMemory) return "armed"
  ```

  Sampled via
  `eval return string.format("<ctx> count=%d frame=%d", pfCount, emu.getState().frameCount)`
  twice per context. Screenshots per non-title context (local artifacts,
  gitignored). Transcript appended to `mesen/out/exp2.log`.
- **Expected outcomes:**
  - *Supports H-BATTLE-0003:* count delta 0 at title and in field;
    delta ≈ frame delta in battle.
  - *Refutes:* any nonzero delta at title/field (then the routine is a
    general per-frame dispatcher, and FN-0002's name needs revisiting).
- **Falsifying outcome:** nonzero count delta in a verified non-battle
  context.
- **Raw evidence paths:** `mesen/out/exp2.log`, `mesen/out/exp2-mines.png`,
  `mesen/out/exp2-battle.png`.
- **Result:** (transcript `mesen/out/exp2.log`; screenshots
  `exp2-mines.png`, `exp2-battle.png` [actually a Narshe guard event —
  checkpoint1 turned out to be field], `exp2-cp2.png` [field],
  `exp2-battle-real.png` [genuine random encounter in the mines])
  - **Non-battle: zero fires.** Counter armed at title frame ~34250 and
    stayed **0** through the title screen, three field/event states
    (checkpoint1 guard event, checkpoint2 mines entrance, checkpoint3 mine
    tunnels), all loadstates, and 10 walking moves — cumulatively
    ≈175,000 frames of non-battle emulation — until the random encounter
    at frame ~209200.
  - **Battle: fires immediately and repeatedly.** Entry phase:
    +75 fires / +320 frames (≈0.23/frame). Settled battle:
    +322 / +370 (≈0.87/frame). The bridge's own `ROMCPU:$C10DF3` counter
    tracked in lockstep (hit=1 at frame 209230; 3000+ hits by 213724).
  - **Rate correction:** "once per frame, every frame" (Sessions
    002/003 phrasing) is an overgeneralization — the archived Session 003
    log itself shows ~0.6–0.7/frame during action phases. *Interpretation:*
    skipped/lag frames during heavy phases; *alternative:* conditional
    dispatch in the bank-`$C2` frame driver. Undiscriminated.
  - **New discovery:** the **first** `$C10DF3` call at battle entry
    returned to `ROMCPU:$C11093` — a second caller, `JSR` at
    `ROMCPU:$C11090`, one-shot during battle init (every later hit
    returned to `$C10203` as before).
  - **Environment hazard:** Mesen's auto-save overwrote slot `_11.mss`
    mid-session (mtime 22:08), destroying the Session 002/003 Narshe
    battle state — recorded in
    [MESEN_CAPABILITY_MATRIX.md](../research/MESEN_CAPABILITY_MATRIX.md).
  - Battle-init writers `ROMCPU:$C22408`/`$C25D33` re-observed in this
    independent battle (frames 209184/209186), corroborating Session 003.
- **Status:** completed (2026-07-29)
- **Confidence:** "Fires only during battle" — **Confirmed for the tested
  contexts** (title, field navigation, field event scenes); world map,
  in-field menu, and other modes remain unsampled, so H-BATTLE-0003
  overall stays Strong hypothesis with a narrowed discriminator. Fire-rate
  variability — Confirmed observation; its cause — Unknown.
- **Next action:** update FN-0001/FN-0002 records (second caller, rate
  wording); sample world map/menu contexts in a future unit; investigate
  the `$C11090` battle-init call site.
