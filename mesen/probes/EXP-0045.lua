-- EXP-0045: the gate transition, one sample per emulated frame.
--
-- EXP-0044 saw a domain move once just after ROMCPU:$C21124's gate
-- engaged, then nothing for 920 frames. Bridge round-trips are hundreds
-- of frames apart, so they cannot say whether that was work resolving
-- past the gate or the last un-gated frame. Sample every frame instead.
--
-- Buffered in memory and flushed on demand: per-frame file I/O would
-- distort the timing under study.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

watchwrites("EXP45-GATEW", 0x2F41, 0x2F41)

local OUTDIR = _G.FF6_OUT_DIR or (os.getenv and os.getenv("FF6_OUT")) or "mesen/out/"

_G.exp45_buf = {}
_G.exp45_left = 0

local function r(a) return emu.read(a, emu.memType.snesWorkRam) end
local function w(base, slot)
  local a = base + slot * 2
  return r(a) | (r(a + 1) << 8)
end

emu.addEventCallback(function()
  if _G.exp45_left <= 0 then return end
  _G.exp45_left = _G.exp45_left - 1
  _G.exp45_buf[#_G.exp45_buf + 1] = string.format(
    "f=%d gate=%02X wait=%02X tick=%04X g=%04X,%04X,%04X fl=%04X,%04X,%04X ac=%04X,%04X,%04X",
    emu.getState().frameCount, r(0x2F41), r(0x3A8F),
    r(0x3A3E) | (r(0x3A3F) << 8),
    w(0x3AB4, 6), w(0x3AB4, 7), w(0x3AB4, 8),
    w(0x3AA0, 6), w(0x3AA0, 7), w(0x3AA0, 8),
    w(0x3218, 6), w(0x3218, 7), w(0x3218, 8))
end, emu.eventType.endFrame)

-- exp45_arm(n): record the next n frames.
_G.exp45_arm = function(n)
  _G.exp45_buf = {}
  _G.exp45_left = n
  return "armed for " .. n .. " frames"
end

-- exp45_flush(name): write the buffer out and report how many frames.
_G.exp45_flush = function(name)
  local f = io.open(OUTDIR .. name, "w")
  if not f then return "cannot open " .. OUTDIR .. name end
  f:write(table.concat(_G.exp45_buf, "\n"), "\n")
  f:close()
  return string.format("wrote %d frames to %s (left=%d)",
    #_G.exp45_buf, name, _G.exp45_left)
end

_G.exp45_mode = function(v)
  emu.write(0x3A8F, v, emu.memType.snesWorkRam)
  return string.format("3A8F <- %02X (%s)", v, (v ~= 0) and "Wait" or "Active")
end

_G.exp45_gate = function()
  return string.format("gate2F41=%02X wait3A8F=%02X", r(0x2F41), r(0x3A8F))
end

return "EXP-0045 probe armed"
