-- Shared probe helpers. Loaded by probes via: dofile("mesen/probes/common.lua")
-- Requires the bridge to be active (uses the same relative output layout).
local OUT = "mesen/out/"
do
  local ok, cfg = pcall(dofile, "mesen/bridge_config.lua")
  if ok and type(cfg) == "table" and type(cfg.out) == "string" then OUT = cfg.out end
end

-- probelog(tag): append one labeled line with registers and a 24-byte
-- stack window to events.log. Identical layout to the historical dlog so
-- transcripts stay comparable across experiments.
function probelog(tag)
  pcall(function()
    local c = emu.getCpuState(emu.cpuType.snes)
    local pc = (c.k << 16) | c.pc
    local s = {}
    for i = 1, 24 do
      s[#s + 1] = string.format("%02X", emu.read(0x7E0000 | ((c.sp + i) & 0xFFFF), emu.memType.snesMemory))
    end
    local f = io.open(OUT .. "events.log", "a")
    if f then
      f:write(string.format(
        "%s pc=$%06X A=$%04X X=$%04X Y=$%04X SP=$%04X PS=$%02X DB=$%02X frame=%d\n    stack: %s\n",
        tag, pc, c.a, c.x, c.y, c.sp, c.ps, c.dbr, emu.getState().frameCount, table.concat(s, " ")))
      f:close()
    end
  end)
end

-- watchwrites(name, lo, hi, cap): first-capture-per-PC write watch over a
-- WRAM range with per-PC counts in _G[name.."_seen"].
function watchwrites(name, lo, hi, cap)
  _G[name .. "_seen"] = {}
  _G[name .. "_ref"] = emu.addMemoryCallback(function(a, v)
    local c = emu.getCpuState(emu.cpuType.snes)
    local pc = (c.k << 16) | c.pc
    local k = string.format("%06X", pc)
    local seen = _G[name .. "_seen"]
    if seen[k] then seen[k] = seen[k] + 1 return end
    seen[k] = 1
    probelog(string.format("%s a=%06X v=%s", name, a, tostring(v)))
  end, emu.callbackType.write, lo, hi, emu.cpuType.snes, emu.memType.snesWorkRam)
end
