-- EXP-0031: golden route segment 1 — power-on → New Game (SCN-0001 B01).
-- Loaded from a pre-placed cmd.txt at frame ~10 of a fresh power-on run.
-- Schedule (absolute frames from power-on, Phase A anchors):
--   2500-4200  start+a edge toggling (6 polls on / 6 off) across the first
--              title window (title visible ~2969, faded by ~3709; long
--              holds and sparse presses do NOT register — edges do).
--   5200       milestone 00-new-game: WRAM dump + screenshot + savestate.
--   30000      stalled Narshe-entry dialogue: WRAM dump + screenshot
--              (downstream assertion that the REAL opening ran, not the
--              attract — the attract would be near a title loop here).
-- Run identity comes from FF6_RUN (run1/run2) for the determinism pair.
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local MIL = OUT .. "../../local_artifacts/scenarios/SCN-0001/00-new-game/"

local T0, T1 = 2500, 4200
local MILESTONE_FRAME = 5200
local STALL_FRAME = 30000

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP31 " .. s .. "\n") f:close() end
end

local function wramDump(path)
  local parts = {}
  for base = 0, 0x1FFFF, 0x1000 do
    local chunk = {}
    for i = 0, 0xFFF do
      chunk[i + 1] = string.char(emu.read(base + i, emu.memType.snesWorkRam))
    end
    parts[#parts + 1] = table.concat(chunk)
  end
  local f = io.open(path, "wb")
  if f then f:write(table.concat(parts)) f:close() end
end

local function shot(path)
  local png = emu.takeScreenshot()
  local f = io.open(path, "wb")
  if f then f:write(png) f:close() end
end

local togglingOn = false
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if fc >= T0 and fc <= T1 then
    local on = math.floor(fc / 6) % 2 == 0
    if on then emu.setInput({ start = true, a = true }, 0) end
    if on ~= togglingOn then
      togglingOn = on
      if on then log(string.format("PRESS-EDGE frame=%d", fc)) end
    end
  end
end, emu.eventType.inputPolled)

local didMilestone, didStall = false, false
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if not didMilestone and fc >= MILESTONE_FRAME then
    didMilestone = true
    wramDump(MIL .. RUN .. "-wram-milestone.bin")
    shot(MIL .. RUN .. "-milestone.png")
    -- createSavestate (like loadSavestate) is only legal inside a
    -- main-CPU exec callback; run it one-shot on the next instruction.
    local ref
    ref = emu.addMemoryCallback(function()
      emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF,
        emu.cpuType.snes, emu.memType.snesMemory)
      local ok, st = pcall(emu.createSavestate)
      if ok and type(st) == "string" and #st > 0 then
        local f = io.open(MIL .. RUN .. "-00-new-game.mss", "wb")
        if f then f:write(st) f:close() end
        log(string.format("MILESTONE00 state=%d frame=%d", #st,
          emu.getState().frameCount))
      else
        log(string.format("MILESTONE00 state=FAILED %s", tostring(st)))
      end
    end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes,
      emu.memType.snesMemory)
    log(string.format("MILESTONE00 dumps frame=%d", fc))
  end
  if not didStall and fc >= STALL_FRAME then
    didStall = true
    wramDump(MIL .. RUN .. "-wram-stall.bin")
    shot(MIL .. RUN .. "-stall.png")
    log(string.format("STALL-ASSERT frame=%d", fc))
    log("RUN-COMPLETE " .. RUN)
  end
end, emu.eventType.endFrame)

log(string.format("armed run=%s window=%d-%d milestone=%d stall=%d",
  RUN, T0, T1, MILESTONE_FRAME, STALL_FRAME))
return "EXP-0031 armed " .. RUN
