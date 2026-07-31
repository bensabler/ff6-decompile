-- EXP-0020: AI/selection-layer RNG hunt (question #30).
-- Logs the $3A70 scheduler/RNG candidate at its pre-action read site and
-- census of its writers, plus action-setup markers for correlation.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

-- Value of $3A70 at the LDA $3A70 (ROMCPU:$C23198) before each setup.
_G.exp20_reads = 0
_G.exp20_readRef = emu.addMemoryCallback(function()
  exp20_reads = exp20_reads + 1
  if exp20_reads <= 24 then
    local v = emu.read(0x7E3A70, emu.memType.snesMemory)
    probelog(string.format("EXP20-READ3A70 val=%02X n=%d", v, exp20_reads))
  end
end, emu.callbackType.exec, 0xC23198, 0xC23198, emu.cpuType.snes, emu.memType.snesMemory)

-- Who advances $3A70/$3A71, with counts.
watchwrites("EXP20-3A70W", 0x3A70, 0x3A71, 24)

-- Action markers: record loader (MVN entry) and fight populate.
_G.exp20_mvn = 0
_G.exp20_mvnRef = emu.addMemoryCallback(function()
  exp20_mvn = exp20_mvn + 1
  if exp20_mvn <= 12 then probelog("EXP20-MVN-LOAD") end
end, emu.callbackType.exec, 0xC22966, 0xC22966, emu.cpuType.snes, emu.memType.snesMemory)
_G.exp20_pop = 0
_G.exp20_popRef = emu.addMemoryCallback(function()
  exp20_pop = exp20_pop + 1
  if exp20_pop <= 12 then probelog("EXP20-POPULATE") end
end, emu.callbackType.exec, 0xC229D1, 0xC229D1, emu.cpuType.snes, emu.memType.snesMemory)

return "EXP-0020 probe armed"
