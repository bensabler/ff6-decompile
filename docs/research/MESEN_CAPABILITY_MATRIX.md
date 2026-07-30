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
| VRAM access | Unknown | | | | |
| CGRAM access | Unknown | | | | |
| OAM access | Unknown | | | | |
| ARAM access | Unknown | | | | |
| DSP register inspection | Unknown | | | | |

Mesen's published code and historical documentation support substantial SNES debugging and Lua capabilities, but export behavior and UI details must be verified against the installed build.

> **Hazard (observed 2026-07-29, EXP-0002):** slot file
> `SaveStates/<rom>_11.mss` is Mesen's **auto-save slot** — it was
> silently overwritten mid-session (mtime 22:08) while the emulator idled,
> destroying the Session 002/003 Narshe-battle state. Any state worth
> keeping must be copied out of `SaveStates/` immediately (the
> `mesen/out/checkpoint*.mss` practice). Never cite `_11.mss` as stable
> evidence.
