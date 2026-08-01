-- EXP-0043: what reads WRAM:+$3A90 (the battle-local Battle Speed value
-- from EXP-0042), and what does it advance?
--
-- Read-watch the two battle-local configuration cells and let the battle
-- run untouched. The reading PCs name the consumer; a separate
-- differential pass over interval savestates names the cells that
-- advance. The two must converge before anything is claimed.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

watchreads("EXP43-SPD", 0x3A8F, 0x3A90)

-- exp43_snap(base, len): hex dump of a WRAM window, for sampling a
-- candidate array at successive frames without leaving the bridge.
_G.exp43_snap = function(base, len)
  local parts = {}
  for i = 0, len - 1 do
    parts[#parts + 1] = string.format("%02X", emu.read(base + i, emu.memType.snesWorkRam))
  end
  return string.format("f=%d %04X: %s", emu.getState().frameCount, base, table.concat(parts, " "))
end

_G.exp43_report = function()
  return table.concat({
    "frame=" .. tostring(emu.getState().frameCount),
    string.format("3A8F=%02X 3A90=%02X",
      emu.read(0x3A8F, emu.memType.snesWorkRam),
      emu.read(0x3A90, emu.memType.snesWorkRam)),
    watchdump("EXP43-SPD"),
  }, "\n")
end

return "EXP-0043 probe armed"
