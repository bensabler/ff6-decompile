# Mesen Capability Matrix

Fill this using the exact local Mesen build before relying on a feature.

**Local build identity:** Mesen 2.1.1, macOS x86_64 (Intel), at
`~/Desktop/Mesen.app` — version recorded from the running emulator during
[SESSION_002](../sessions/SESSION_002.md). The app bundle `Info.plist`
reports only a generic `1.0.0`, so re-verify with `eval emu.getVersion()`
through the bridge on next launch.

| Capability | Available | Export method | Automation method | Verified date | Notes |
|---|---|---|---|---|---|
| Main CPU debugger | Yes (via Lua) | log lines to `mesen/out/events.log` | `emu.addEventCallback` exec callbacks + `emu.getCpuState(emu.cpuType.snes)` | 2026-07-29 | GUI breakpoints used manually in Session 001; programmatic path verified in Session 002. Generic `getState()` does NOT expose SNES registers. |
| SPC debugger | Unknown | | | | |
| Trace logger | Unknown | | | | |
| Tile viewer | Unknown | | | | |
| Tilemap viewer | Unknown | | | | |
| Sprite viewer | Unknown | | | | |
| Event viewer / DMA filtering | Unknown | | | | |
| Memory tools | Yes (via Lua) | hex dumps to `mesen/out/resp.txt` | bridge `read wram|cpu <addr> <len>`; GUI memory search used manually in Session 001 | 2026-07-29 | WRAM and CPU-visible reads verified live against displayed values. |
| Lua scripting | Yes | file I/O from script | `mesen/bridge.lua` command loop (`cmd.txt` → `resp.txt`) | 2026-07-29 | Requires `AllowIoOsAccess: true` in `Debug.ScriptWindow` settings. Script auto-reload on file change does not work — restart Mesen after editing. |
| Save-state load (Lua) | Yes, with caveat | n/a | bridge `loadstate` | 2026-07-29 | `emu.loadSavestate` only works inside a main-CPU exec callback; the bridge queues it. |
| Screenshot (Lua) | Yes | PNG to `mesen/out/` | bridge `screenshot` | 2026-07-29 | Used for HP-window verification in Session 002. |
| Input injection (Lua) | Yes | n/a | bridge `press <buttons> <frames>` | 2026-07-29 | Used to cast Heal Force in Session 002. |
| VRAM access | **Yes (via Lua)** | hex dump to `mesen/out/resp.txt` | bridge `read vram <addr> <len>` (`emu.memType.snesVideoRam`, `bridge.lua:241`) | 2026-08-02 | Wired since the bridge's `READ_MEMTYPES` table. **Exercised**: EXP-0023 dumped the BG3 tilemap and CHR region (`mesen/out/exp23-tilemap.hex`, `exp23-chr.hex`). |
| CGRAM access | **Yes (via Lua)** | hex dump to `mesen/out/resp.txt` | bridge `read cgram <addr> <len>` (`emu.memType.snesCgRam`, `bridge.lua:242`) | 2026-08-02 | **Exercised**: EXP-0023 dumped all 512 bytes (`mesen/out/exp23-cgram.hex`). |
| OAM access | **Yes (via Lua)** | hex dump to `mesen/out/resp.txt` | bridge `read oam <addr> <len>` (`emu.memType.snesSpriteRam`, `bridge.lua:243`) | 2026-08-02 | Wired and reachable, but **never exercised** — no experiment has yet read OAM. Readiness F6 depends on it. |
| ARAM access | **Yes (via Lua)** | hex dump to `mesen/out/resp.txt` | bridge `read aram <addr> <len>` (`emu.memType.spcRam`, `bridge.lua:244`) | 2026-08-02 | **Exercised**: EXP-0024 / AUD-0001 read the SFX pack at `ARAM:$4800-$491F`. |
| DSP register inspection | **Yes (via Lua)** | hex dump to `mesen/out/resp.txt` | bridge `read dsp <addr> <len>` (`emu.memType.spcDspRegisters`, `bridge.lua:246`) | 2026-08-02 | **Exercised**: EXP-0024 captured voice/DSP state (`mesen/out/exp24-dsp.hex`, `exp24-voices.log`). |
| **Offline savestate inspection** | **Yes (no emulator)** | `ff6lab state …` | `internal/mesenstate` parses `.mss` zlib sections into named blocks | 2026-08-02 | Every preserved `.mss` carries `memoryManager.workRam` (128 KiB), `cart.saveRam` (8 KiB), `ppu.vram` (64 KiB), `ppu.cgram` (512 B), `ppu.oamRam` (544 B), `spc.ram` (64 KiB), the full PPU register set and per-channel `dmaController.*` state. `ff6lab state` currently exposes **only wram and sram**; the rest are readable but unexposed. |
| DMA register capture (live) | **No** | — | — | 2026-08-02 | `/trace-dma`, the `dma-tracer` skill, the `dma-researcher` agent and `TRACE_DMA.md` all exist; **no Lua probe or Go code reads a DMA register**. The savestate `dmaController.*` blocks are the only DMA evidence available today, and they give channel *configuration at an instant*, not a trace. |

Mesen's published code and historical documentation support substantial SNES debugging and Lua capabilities, but export behavior and UI details must be verified against the installed build.

> **Correction (2026-08-02, Unit 10):** the five rows above read `Unknown`
> until this revision, although `bridge.lua`'s `READ_MEMTYPES` has wired all of
> them since the bridge was written and four were exercised by EXP-0023 and
> EXP-0024. The matrix was understating the project's own capability, which is
> a planning hazard: work was deferred as "needs new instrumentation" when the
> instrument already existed. Rows are now marked wired-and-exercised, wired-
> but-never-exercised, or absent, because those are three different facts.

> **Hazard (observed 2026-07-29, EXP-0002):** slot file
> `SaveStates/<rom>_11.mss` is Mesen's **auto-save slot** — it was
> silently overwritten mid-session (mtime 22:08) while the emulator idled,
> destroying the Session 002/003 Narshe-battle state. Any state worth
> keeping must be copied out of `SaveStates/` immediately (the
> `mesen/out/checkpoint*.mss` practice). Never cite `_11.mss` as stable
> evidence.
