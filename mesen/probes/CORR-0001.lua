-- CORR-0001: bounded static/runtime correlation at ROMCPU:$C09B5C.
--
-- Question: does execution reach exactly $C09B5C, by what transfer, in what
-- 65816 state, and what does it actually do to the direct-page operands
-- $E5-$E7 and $E3?
--
-- Nothing here assumes the static reading is right. The probe records raw
-- state and lets the record decide. In particular it does NOT assume M=1;
-- it reports the P register and lets the analysis read the M bit.
--
-- Starting state comes from a pre-existing SCN-0001 savestate (CORR_STATE),
-- so the run is short and reproducible. Two observations use two
-- independently produced states from the EXP-0032 determinism pair.
--
-- Env: CORR_STATE (path to .mss), CORR_TAG, CORR_MAXHITS, CORR_FRAMES.
-- FF6_OUT is required under --testrunner: debug.getinfo returns "=[C]" there,
-- so script-directory resolution is unavailable (EXP-0037.lua:20-22).
local OUT = _G.FF6_OUT_DIR or (os.getenv and os.getenv("FF6_OUT")) or "mesen/out/"
local STATE = os.getenv and os.getenv("CORR_STATE") or nil
local TAG = (os.getenv and os.getenv("CORR_TAG")) or "obs"
local MAXHITS = tonumber((os.getenv and os.getenv("CORR_MAXHITS")) or "12")
local ENDFRAME = tonumber((os.getenv and os.getenv("CORR_FRAMES")) or "1800")

local TARGET = 0xC09B5C -- routine under test
local CONT = 0xC09A6D   -- its JMP target
local WAITPATH = 0xC09A71 -- $C09A6D fallthrough: DEX / STX $E3 / RTS
local BYPASS = 0xC09A75   -- $C09A6D BEQ target when $E3 == 0
local HANDLO, HANDHI = 0xC0B593, 0xC0B6D2 -- confirmed handler family (ROM-0027)

local logf = OUT .. "corr0001-" .. TAG .. ".log"

local function log(s)
  local f = io.open(logf, "a")
  if f then f:write(s .. "\n") f:close() end
end

local function frame() return emu.getState().frameCount end

-- Direct-page operands live in bank 0. The effective address of operand
-- <off> is (D + off) & 0xFFFF -- this is exactly why DPR is load-bearing.
local function dpAddr(d, off) return (d + off) & 0xFFFF end
local function rd(addr) return emu.read(addr & 0xFFFFFF, emu.memType.snesMemory) end

local function ptr24(d)
  local a5, a6, a7 = dpAddr(d, 0xE5), dpAddr(d, 0xE6), dpAddr(d, 0xE7)
  return rd(a5) | (rd(a6) << 8) | (rd(a7) << 16), a5, a6, a7
end

-- P register bits (native mode): M = $20, X = $10. E is a separate boolean
-- field on the CpuState (`emulationMode`), not a P bit (bridge.lua:144).
local function flags(ps, e)
  return string.format("E=%s M=%d X=%d (P=$%02X)",
    tostring(e), (ps & 0x20) ~= 0 and 1 or 0, (ps & 0x10) ~= 0 and 1 or 0, ps)
end

local hits, done = 0, false
local startFrame = nil -- set at arm(); the loaded state resumes a nonzero
                       -- frame counter, so the window must be relative
local pending = nil
local lastHandler, lastHandlerFrame = nil, -1
local refs = {}

local function cpu() return emu.getCpuState(emu.cpuType.snes) end

-- Discover the CpuState field names rather than guessing which key holds
-- the direct-page register or the emulation flag.
local function dumpFields()
  local c = cpu()
  local keys = {}
  for k, v in pairs(c) do keys[#keys + 1] = string.format("%s=%s", k, tostring(v)) end
  table.sort(keys)
  log("CPUSTATE-FIELDS " .. table.concat(keys, " "))
end

local function arm()
  dumpFields()
  startFrame = frame()
  log(string.format("ARMED state=%s tag=%s frame=%d window=%d",
    tostring(STATE), TAG, startFrame, ENDFRAME))

  -- Predecessor evidence: the confirmed handler family. Records the last
  -- handler PC seen so a hit at TARGET can be attributed (or not).
  refs[#refs + 1] = emu.addMemoryCallback(function()
    local c = cpu()
    lastHandler = (c.k << 16) | c.pc
    lastHandlerFrame = frame()
  end, emu.callbackType.exec, HANDLO, HANDHI, emu.cpuType.snes, emu.memType.snesMemory)

  -- The routine under test.
  refs[#refs + 1] = emu.addMemoryCallback(function()
    if done then return end
    local c = cpu()
    local pc = (c.k << 16) | c.pc
    local d = c.d or -1
    local before, a5, a6, a7 = ptr24(d)
    local e3a = dpAddr(d, 0xE3)
    hits = hits + 1
    log(string.format(
      "HIT #%d frame=%d pcReported=$%06X targetMatch=%s %s PBR=$%02X DBR=$%02X D=$%04X SP=$%04X A=$%04X X=$%04X Y=$%04X",
      hits, frame(), pc, tostring(pc == TARGET), flags(c.ps, c.emulationMode), c.k, c.dbr, d, c.sp, c.a, c.x, c.y))
    log(string.format(
      "   effAddr $E5-$E7 -> $00%04X/$00%04X/$00%04X  ptr24Before=$%06X   effAddr $E3 -> $00%04X val=$%02X",
      a5, a6, a7, before, e3a, rd(e3a)))
    log(string.format("   predecessor: lastHandlerPC=%s (frame=%d, same=%s)",
      lastHandler and string.format("$%06X", lastHandler) or "none",
      lastHandlerFrame, tostring(lastHandlerFrame == frame())))
    pending = { n = hits, before = before, d = d, a = c.a, ps = c.ps, frame = frame() }
    if hits >= MAXHITS then done = true end
  end, emu.callbackType.exec, TARGET, TARGET, emu.cpuType.snes, emu.memType.snesMemory)

  -- The claimed continuation. Reached immediately after? Pointer changed?
  refs[#refs + 1] = emu.addMemoryCallback(function()
    if not pending then return end
    local c = cpu()
    local after = ptr24(pending.d)
    local delta = (after - pending.before) & 0xFFFFFF
    -- Only the low byte of A is the addend when M=1; report both readings.
    log(string.format(
      "   -> reached $C09A6D frame=%d ptr24After=$%06X delta=%d  A@entry=$%04X (low=%d)  deltaMatchesAlow=%s",
      frame(), after, delta, pending.a, pending.a & 0xFF, tostring(delta == (pending.a & 0xFF))))
    log(string.format("   $E3 at continuation = $%02X", rd(dpAddr(pending.d, 0xE3))))
    pending = nil
  end, emu.callbackType.exec, CONT, CONT, emu.cpuType.snes, emu.memType.snesMemory)

  -- Which branch $C09A6D takes.
  refs[#refs + 1] = emu.addMemoryCallback(function()
    log(string.format("   BRANCH wait-path $C09A71 (E3 nonzero) frame=%d", frame()))
  end, emu.callbackType.exec, WAITPATH, WAITPATH, emu.cpuType.snes, emu.memType.snesMemory)

  refs[#refs + 1] = emu.addMemoryCallback(function()
    log(string.format("   BRANCH bypass $C09A75 (E3 zero) frame=%d", frame()))
  end, emu.callbackType.exec, BYPASS, BYPASS, emu.cpuType.snes, emu.memType.snesMemory)
end

-- Load the starting state, then arm. loadSavestate is only legal inside a
-- main-CPU exec callback, so use the bridge's one-shot pattern.
if STATE then
  local f = io.open(STATE, "rb")
  if not f then
    log("FATAL cannot open state " .. STATE)
  else
    local data = f:read("*a")
    f:close()
    local ref
    ref = emu.addMemoryCallback(function()
      emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
      local ok, err = pcall(emu.loadSavestate, data)
      log(string.format("LOADSTATE ok=%s detail=%s bytes=%d frame=%d",
        tostring(ok), tostring(err), #data, frame()))
      arm()
    end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
  end
else
  arm()
end

emu.addEventCallback(function()
  if not startFrame then return end
  if done or frame() >= startFrame + ENDFRAME then
    log(string.format("END frame=%d elapsed=%d hits=%d", frame(), frame() - startFrame, hits))
    emu.stop(0)
  end
end, emu.eventType.endFrame)
