-- EXP-0023: battle HUD font capture (graphics vertical, Unit 4).
-- exp23_capture(mss_path, offset_frames) reloads the battle savestate
-- and, at anchor+offset, atomically dumps CGRAM, the BG3 tilemap, the
-- BG3 chr region, and the relevant PPU fields to OUT files. The frame
-- is reproducible because the window is schedule-deterministic with no
-- input (EXP-0021).
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")
local OUT = _G.FF6_OUT_DIR or "mesen/out/"

-- Byte offsets derived from the phase-1 PPU read (word addresses
-- $5000/$7800 on layer 2): chr $A000, tilemap $F000.
local CHR_BASE, CHR_LEN = 0xA000, 8192
local MAP_BASE, MAP_LEN = 0xF000, 2048

local function dumphex(path, memType, base, len)
  local parts = {}
  for i = 0, len - 1 do
    parts[#parts + 1] = string.format("%02X", emu.read(base + i, memType))
  end
  local f = io.open(path, "w")
  if f then f:write(table.concat(parts, " ") .. "\n") f:close() end
end

_G.exp23_sched = nil
emu.addEventCallback(function()
  local s = exp23_sched
  if not s or not s.anchor then return end
  local fc = emu.getState().frameCount
  if fc < s.anchor + s.offset then return end
  exp23_sched = nil
  local ok, err = pcall(function()
    local st = emu.getState()
    local f = io.open(OUT .. "exp23-ppu-state.txt", "w")
    if f then
      f:write(string.format("captureFrame=%d anchor=%d offset=%d\n", fc, s.anchor, s.offset))
      local fields = {
        "ppu.bgMode", "ppu.mainScreenLayers", "ppu.subScreenLayers",
        "ppu.forcedBlank", "ppu.screenBrightness", "ppu.mode1Bg3Priority",
      }
      for i = 0, 3 do
        fields[#fields + 1] = string.format("ppu.layers[%d].chrAddress", i)
        fields[#fields + 1] = string.format("ppu.layers[%d].tilemapAddress", i)
      end
      for _, k in ipairs(fields) do
        f:write(string.format("%s=%s\n", k, tostring(st[k])))
      end
      f:close()
    end
    dumphex(OUT .. "exp23-cgram.hex", emu.memType.snesCgRam, 0, 512)
    dumphex(OUT .. "exp23-tilemap.hex", emu.memType.snesVideoRam, MAP_BASE, MAP_LEN)
    dumphex(OUT .. "exp23-chr.hex", emu.memType.snesVideoRam, CHR_BASE, CHR_LEN)
    probelog(string.format("EXP23-CAPTURE frame=%d", fc))
  end)
  if not ok then
    probelog("EXP23-CAPTURE-ERROR " .. tostring(err))
  end
end, emu.eventType.endFrame)

_G.exp23_capture = function(mss, offset)
  local f = io.open(mss, "rb")
  if not f then return "cannot open " .. mss end
  local data = f:read("*a")
  f:close()
  exp23_sched = { offset = offset }
  local ref
  ref = emu.addMemoryCallback(function()
    emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
    local okL = pcall(emu.loadSavestate, data)
    if okL and exp23_sched then
      exp23_sched.anchor = emu.getState().frameCount
      probelog(string.format("EXP23-LOAD anchor=%d offset=%d", exp23_sched.anchor, offset))
    else
      probelog("EXP23-LOADFAIL")
      exp23_sched = nil
    end
  end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
  return "exp23 capture queued offset=" .. offset
end

return "EXP-0023 probe armed"
