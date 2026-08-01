-- EXP-0046: who completes a pending action while ROMCPU:$C21124's gate
-- is shut?
--
-- EXP-0045 showed +$3AA0 bit 6 clearing and +$3218 gaining $0100 during
-- a gated interval, while the tick counter and every gauge stayed
-- frozen. The writer therefore runs outside the per-frame scheduler.
--
-- Watch both ranges and tag each write with the gate state. The
-- discriminator is countGated > 0: a PC writing these cells while the
-- scheduler is gated cannot be the scheduler.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

local M = emu.memType.snesWorkRam
local function r(a) return emu.read(a, M) end

-- gatedwatch(name, lo, hi): per-PC totals split by gate state, plus a
-- stack capture on the first few gated writes so the caller chain can be
-- recovered.
local function gatedwatch(name, lo, hi)
  _G[name .. "_seen"] = {}
  _G[name .. "_logged"] = 0
  _G[name .. "_ref"] = emu.addMemoryCallback(function(a, v)
    local c = emu.getCpuState(emu.cpuType.snes)
    local pc = (c.k << 16) | c.pc
    local k = string.format("%06X", pc)
    local seen = _G[name .. "_seen"]
    local gated = (r(0x2F41) ~= 0) and (r(0x3A8F) ~= 0)
    local f = emu.getState().frameCount
    local e = seen[k]
    if not e then
      e = { count = 0, gated = 0, firstFrame = f, lastFrame = f, lastAddr = a, lastVal = v }
      seen[k] = e
    end
    e.count = e.count + 1
    e.lastFrame = f
    e.lastAddr = a
    e.lastVal = v
    if gated then
      e.gated = e.gated + 1
      if _G[name .. "_logged"] < 12 then
        _G[name .. "_logged"] = _G[name .. "_logged"] + 1
        probelog(string.format("%s GATED-WRITE a=%06X v=%s", name, a, tostring(v)))
      end
    end
  end, emu.callbackType.write, lo, hi, emu.cpuType.snes, M)
end

gatedwatch("EXP46-FLAGS", 0x3AA0, 0x3AB3)
gatedwatch("EXP46-ACCUM", 0x3218, 0x322B)

-- Report only the PCs that wrote while gated — the whole point.
_G.exp46_gated = function(name)
  local seen = _G[name .. "_seen"] or {}
  local rows = {}
  for pc, e in pairs(seen) do
    if e.gated > 0 then
      rows[#rows + 1] = string.format(
        "pc=$%s gated=%d/%d first=%d last=%d addr=$%06X val=%s",
        pc, e.gated, e.count, e.firstFrame, e.lastFrame, e.lastAddr, tostring(e.lastVal))
    end
  end
  table.sort(rows)
  if #rows == 0 then return name .. ": no gated writes" end
  return name .. " (gated writers):\n" .. table.concat(rows, "\n")
end

_G.exp46_all = function(name) return watchdump(name) end

_G.exp46_mode = function(v)
  emu.write(0x3A8F, v, M)
  return string.format("3A8F <- %02X (%s)", v, (v ~= 0) and "Wait" or "Active")
end

_G.exp46_state = function()
  local w = function(base, slot)
    local a = base + slot * 2
    return r(a) | (r(a + 1) << 8)
  end
  return string.format(
    "f=%d gate=%02X wait=%02X tick=%04X fl=%04X,%04X,%04X ac=%04X,%04X,%04X",
    emu.getState().frameCount, r(0x2F41), r(0x3A8F),
    r(0x3A3E) | (r(0x3A3F) << 8),
    w(0x3AA0, 6), w(0x3AA0, 7), w(0x3AA0, 8),
    w(0x3218, 6), w(0x3218, 7), w(0x3218, 8))
end

-- Pending-triggered press, as in EXP-0045 trace 3: bridge round-trips
-- are hundreds of frames apart, so the arrangement has to happen here.
_G.exp46_pressLeft = 0
emu.addEventCallback(function()
  if _G.exp46_pressLeft > 0 then
    emu.setInput({ a = true }, 0)
    _G.exp46_pressLeft = _G.exp46_pressLeft - 1
  end
end, emu.eventType.inputPolled)

_G.exp46_armPending = false
emu.addEventCallback(function()
  if not _G.exp46_armPending then return end
  if (r(0x3AAC) & 0x40) ~= 0 or (r(0x3AAE) & 0x40) ~= 0 or (r(0x3AB0) & 0x40) ~= 0 then
    _G.exp46_armPending = false
    _G.exp46_pressLeft = 4
    probelog("EXP46-TRIGGER pending slot seen, opening submenu")
  end
end, emu.eventType.endFrame)

_G.exp46_arm = function()
  _G.exp46_armPending = true
  return "waiting for a pending slot"
end

return "EXP-0046 probe armed"
