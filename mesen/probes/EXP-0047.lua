-- EXP-0047: what invokes the action-execution path, and how often?
--
-- EXP-0046 reconstructed a call chain ending at ROMCPU:$C20016. If the
-- entry fires on most gated frames while a completion happens on only
-- one of them, the delay is a countdown inside the path; if it fires
-- rarely, the delay lives upstream.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

local M = emu.memType.snesWorkRam
local function r(a) return emu.read(a, M) end
local function gated() return r(0x2F41) ~= 0 and r(0x3A8F) ~= 0 end

-- execwatch(name, addr): count executions, split by gate state, and keep
-- a bounded per-frame log so cadence can be read directly.
function execwatch(name, addr)
  _G[name .. "_n"] = 0
  _G[name .. "_ngated"] = 0
  _G[name .. "_log"] = {}
  _G[name .. "_ref"] = emu.addMemoryCallback(function()
    _G[name .. "_n"] = _G[name .. "_n"] + 1
    local g = gated()
    if g then _G[name .. "_ngated"] = _G[name .. "_ngated"] + 1 end
    local t = _G[name .. "_log"]
    if #t < 400 then
      t[#t + 1] = string.format("f=%d gate=%d", emu.getState().frameCount, g and 1 or 0)
    end
  end, emu.callbackType.exec, addr, addr, emu.cpuType.snes, emu.memType.snesMemory)
end

_G.exp47_report = function(name)
  return string.format("%s: total=%d gated=%d logged=%d",
    name, _G[name .. "_n"] or -1, _G[name .. "_ngated"] or -1, #(_G[name .. "_log"] or {}))
end

_G.exp47_dump = function(name, file)
  local f = io.open((_G.FF6_OUT_DIR or "mesen/out/") .. file, "w")
  if not f then return "cannot open " .. file end
  f:write(table.concat(_G[name .. "_log"] or {}, "\n"), "\n")
  f:close()
  return string.format("wrote %d entries", #(_G[name .. "_log"] or {}))
end

_G.exp47_state = function()
  local w = function(base, slot)
    local a = base + slot * 2
    return r(a) | (r(a + 1) << 8)
  end
  return string.format("f=%d gate=%02X wait=%02X tick=%04X fl=%04X,%04X,%04X ac=%04X,%04X,%04X",
    emu.getState().frameCount, r(0x2F41), r(0x3A8F), r(0x3A3E) | (r(0x3A3F) << 8),
    w(0x3AA0, 6), w(0x3AA0, 7), w(0x3AA0, 8), w(0x3218, 6), w(0x3218, 7), w(0x3218, 8))
end

_G.exp47_mode = function(v)
  emu.write(0x3A8F, v, M)
  return string.format("3A8F <- %02X", v)
end

-- Pending-trigger, as in EXP-0045/0046.
_G.exp47_pressLeft = 0
emu.addEventCallback(function()
  if _G.exp47_pressLeft > 0 then
    emu.setInput({ a = true }, 0)
    _G.exp47_pressLeft = _G.exp47_pressLeft - 1
  end
end, emu.eventType.inputPolled)

_G.exp47_armPending = false
emu.addEventCallback(function()
  if not _G.exp47_armPending then return end
  if (r(0x3AAC) & 0x40) ~= 0 or (r(0x3AAE) & 0x40) ~= 0 or (r(0x3AB0) & 0x40) ~= 0 then
    _G.exp47_armPending = false
    _G.exp47_pressLeft = 4
    probelog("EXP47-TRIGGER pending slot seen, opening submenu")
  end
end, emu.eventType.endFrame)

_G.exp47_arm = function()
  _G.exp47_armPending = true
  return "waiting for a pending slot"
end

return "EXP-0047 probe armed"
