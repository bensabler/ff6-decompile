-- EXP-0044: the ACTIVE/WAIT pause matrix.
--
-- ROMCPU:$C21124 skips the whole per-frame battle update when
-- ($2F41 AND $3A8F) is non-zero. $3A8F is the Wait flag (EXP-0042);
-- $2F41 has never been seen non-zero. Catch its writers, then sample
-- every located timer domain in each menu state under both modes.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

watchwrites("EXP44-GATEW", 0x2F41, 0x2F41)

local function word(base, slot)
  local a = base + slot * 2
  return emu.read(a, emu.memType.snesWorkRam)
    | (emu.read(a + 1, emu.memType.snesWorkRam) << 8)
end

local function row(label, base)
  local p = {}
  for s = 0, 9 do p[#p + 1] = string.format("%04X", word(base, s)) end
  return label .. " " .. table.concat(p, " ")
end

-- exp44_sample(): every timer domain in one response, stamped with the
-- emulator frame so the independent clock is never implicit.
_G.exp44_sample = function()
  local r = function(a) return emu.read(a, emu.memType.snesWorkRam) end
  return table.concat({
    string.format("f=%d tick=%04X gate2F41=%02X wait3A8F=%02X spd3A90=%02X",
      emu.getState().frameCount,
      r(0x3A3E) | (r(0x3A3F) << 8), r(0x2F41), r(0x3A8F), r(0x3A90)),
    row("gauge3AB4", 0x3AB4),
    row("flags3AA0", 0x3AA0),
    row("accum3218", 0x3218),
  }, "\n")
end

-- exp44_mode(v): patch the battle-local Wait flag in place. EXP-0042
-- showed battle timing runs off this cell rather than the Config bytes,
-- so writing it is the mode switch.
_G.exp44_mode = function(v)
  emu.write(0x3A8F, v, emu.memType.snesWorkRam)
  return string.format("3A8F <- %02X (%s)", v, (v ~= 0) and "Wait" or "Active")
end

_G.exp44_gate = function() return watchdump("EXP44-GATEW") end

return "EXP-0044 probe armed"
