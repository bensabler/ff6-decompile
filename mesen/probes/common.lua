-- Shared probe helpers. Loaded by probes via: dofile("mesen/probes/common.lua")
-- Requires the bridge to be active (uses the same relative output layout).
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
do
  if _G.FF6_OUT_DIR then goto done end
  local ok, info = pcall(debug.getinfo, 1, "S")
  if ok and type(info) == "table" and type(info.source) == "string" then
    local dir = info.source:gsub("^@", ""):match("^(.*/)")
    if dir then OUT = dir .. "../out/" end -- <repo>/mesen/probes/../out/
  end
  local env = os.getenv and os.getenv("FF6_OUT")
  if env then OUT = env end
  ::done::
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
