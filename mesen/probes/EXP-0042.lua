-- EXP-0042: does battle code read the configuration bytes directly, or a
-- copy sampled at battle entry?
--
-- Read-watch WRAM:+$1D4D-+$1D4E (EXP-0041's config storage) and timestamp
-- the battle-entry boundary, so every reading PC can be sorted into
-- before-entry and after-entry. A PC whose count keeps rising during the
-- battle means direct consultation; a small burst at the boundary and
-- nothing after means an entry-time sample.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

watchreads("EXP42-CFGR", 0x1D4D, 0x1D4E)

-- Battle-entry timestamp: the first write into the $14-stride battler
-- stat block issued from the populate region (EXP-0032's detector, also
-- used by EXP-0033/0034/0036/0038).
_G.exp42_entry_frame = nil
_G.exp42_entry_pc = nil
_G.exp42_bref = emu.addMemoryCallback(function(a, v)
  if _G.exp42_entry_frame then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  if pc >= 0xC22800 and pc <= 0xC22FFF then
    _G.exp42_entry_frame = emu.getState().frameCount
    _G.exp42_entry_pc = pc
    probelog(string.format("EXP42-BATTLE-ENTRY pc=$%06X addr=%06X", pc, a))
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes, emu.memType.snesWorkRam)

-- Movement helper for the encounter walk (EXP-0028 pattern):
-- frame-scheduled through the input-poll hook, never wall-clock.
_G.exp42_hold = nil
emu.addEventCallback(function()
  local h = _G.exp42_hold
  if h and h.left > 0 then
    emu.setInput(h.buttons, 0)
    h.left = h.left - 1
    if h.left == 0 then _G.exp42_hold = nil end
  end
end, emu.eventType.inputPolled)

_G.exp42_walk = function(dir, frames)
  _G.exp42_hold = { buttons = { [dir] = true }, left = frames }
  return "walk " .. dir .. " " .. frames
end

-- exp42_report(): the whole picture in one response — current frame, the
-- battle-entry boundary if it has happened, and the read table.
_G.exp42_report = function()
  local f = emu.getState().frameCount
  local lines = {
    "frame=" .. tostring(f),
    "battleEntryFrame=" .. tostring(_G.exp42_entry_frame),
    "battleEntryPC=" .. (_G.exp42_entry_pc and string.format("$%06X", _G.exp42_entry_pc) or "nil"),
    string.format("cfg $1D4D=%02X $1D4E=%02X",
      emu.read(0x1D4D, emu.memType.snesWorkRam),
      emu.read(0x1D4E, emu.memType.snesWorkRam)),
    watchdump("EXP42-CFGR"),
  }
  return table.concat(lines, "\n")
end

return "EXP-0042 probe armed"
