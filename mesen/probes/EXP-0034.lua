-- EXP-0034: golden route segment 3b — the scripted battle chain to free
-- movement (SCN-0001 milestone 04). Extends EXP-0033 with a re-arming
-- battle detector so the route does not assume a single fight.
--
-- Schedule (absolute frames from power-on):
--   2500-4200    start+a edge toggling (title)
--   15000/30000  milestones 01/02 (continuity)
--   31000+       outside battle: up held + A every 240 frames
--                inside battle:  A 10 polls every 90 frames
--   per battle   entry detected by a +$3B18 write from
--                $C22800-$C22FFF (Confirmed signature, EXP-0032);
--                +$11E0 (formation id) logged at entry
--   battle end   every HP word in slots 3..9 reads zero (slots 0-2 are
--                the party per EXP-0033; enemy slots vary by formation,
--                so no slot is hardcoded), sustained 30 frames
--   quiet        no battle re-arms for QUIET frames after the last end
--                -> milestone 04-free-movement
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local ROOT = OUT .. "../../local_artifacts/scenarios/SCN-0001/"
local EVD = OUT .. "../../local_artifacts/experiments/EXP-0034/"

local T0, T1 = 2500, 4200
local WALK0 = 31000
local METP, METHOLD = 240, 12
local BATP, BATHOLD = 90, 10
local ENDHOLD = 30
-- v1 used QUIET=1200 and captured the item-reward window ("Got Tonic
-- x1"): the reward chain plus the transition to the NEXT scripted
-- battle runs ~3000 frames past battle end (EXP-0033 saw battle 2 by
-- frame ~35700, battle 1 having ended at 32706). A free-movement state
-- persists indefinitely, so a long quiet window is safe; a short one is
-- not.
local QUIET = 5400
local FAILSAFE = 120000

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP34 " .. s .. "\n") f:close() end
end

local function wramDump(path)
  local parts = {}
  for base = 0, 0x1FFFF, 0x1000 do
    local chunk = {}
    for i = 0, 0xFFF do
      chunk[i + 1] = string.char(emu.read(base + i, emu.memType.snesWorkRam))
    end
    parts[#parts + 1] = table.concat(chunk)
  end
  local f = io.open(path, "wb")
  if f then f:write(table.concat(parts)) f:close() end
end

local function shot(path)
  local png = emu.takeScreenshot()
  local f = io.open(path, "wb")
  if f then f:write(png) f:close() end
end

local function mkstate(path, tag)
  local ref
  ref = emu.addMemoryCallback(function()
    emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF,
      emu.cpuType.snes, emu.memType.snesMemory)
    local ok, st = pcall(emu.createSavestate)
    if ok and type(st) == "string" and #st > 0 then
      local f = io.open(path, "wb")
      if f then f:write(st) f:close() end
      log(tag .. " state=" .. #st)
    else
      log(tag .. " state=FAILED " .. tostring(st))
    end
  end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes,
    emu.memType.snesMemory)
end

local function milestone(dir, name)
  wramDump(dir .. RUN .. "-wram.bin")
  shot(dir .. RUN .. "-screen.png")
  mkstate(dir .. RUN .. "-" .. name .. ".mss", "MS-" .. name)
  log(string.format("MILESTONE %s frame=%d", name,
    emu.getState().frameCount))
end

-- battle state machine
local inBattle = false
local battleCount = 0
local battleFrame = nil
local lastEnd = nil
local zeroSince = nil

emu.addMemoryCallback(function(addr, value)
  if inBattle then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  if pc >= 0xC22800 and pc <= 0xC22FFF then
    inBattle = true
    battleCount = battleCount + 1
    battleFrame = emu.getState().frameCount
    zeroSince = nil
    local staged = {}
    for i = 0, 15 do
      staged[#staged + 1] = string.format("%02X",
        emu.read(0x3F44 + i, emu.memType.snesWorkRam))
    end
    log(string.format("BATTLE %d ENTRY frame=%d 11E0=$%04X staged=%s",
      battleCount, battleFrame,
      emu.readWord(0x11E0, emu.memType.snesWorkRam),
      table.concat(staged, " ")))
    shot(EVD .. RUN .. "-battle" .. battleCount .. "-entry.png")
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- input
local togglingOn = false
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if fc >= T0 and fc <= T1 then
    local on = math.floor(fc / 6) % 2 == 0
    if on then emu.setInput({ start = true, a = true }, 0) end
    if on ~= togglingOn then
      togglingOn = on
      if on then log(string.format("PRESS-EDGE frame=%d", fc)) end
    end
  elseif inBattle then
    if battleFrame and fc >= battleFrame + 60
      and (fc - battleFrame) % BATP < BATHOLD then
      emu.setInput({ a = true }, 0)
    end
  elseif fc >= WALK0 then
    local pressA = (fc - WALK0) % METP < METHOLD
    emu.setInput({ up = true, a = pressA or nil }, 0)
  end
end, emu.eventType.inputPolled)

-- battle end + milestone
local did04, failed = false, false
local nextShot = 31000
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if inBattle and battleFrame and fc > battleFrame + 120 then
    local allZero = true
    for slot = 3, 9 do
      if emu.readWord(0x3BF4 + slot * 2, emu.memType.snesWorkRam) ~= 0 then
        allZero = false
        break
      end
    end
    if allZero then
      zeroSince = zeroSince or fc
      if fc - zeroSince >= ENDHOLD then
        inBattle = false
        lastEnd = fc
        log(string.format("BATTLE %d END frame=%d", battleCount, fc))
        shot(EVD .. RUN .. "-battle" .. battleCount .. "-end.png")
      end
    else
      zeroSince = nil
    end
  end
  if fc >= nextShot and not did04 then
    nextShot = nextShot + 600
    shot(EVD .. RUN .. "-beat-" .. fc .. ".png")
  end
  if not did04 and lastEnd and not inBattle and fc >= lastEnd + QUIET then
    did04 = true
    milestone(ROOT .. "04-free-movement/", "04")
    log(string.format("RUN-COMPLETE %s battles=%d", RUN, battleCount))
  end
  if not did04 and not failed and fc >= FAILSAFE then
    failed = true
    shot(EVD .. RUN .. "-failsafe.png")
    log(string.format("FAILSAFE frame=%d battles=%d inBattle=%s",
      fc, battleCount, tostring(inBattle)))
  end
end, emu.eventType.endFrame)

log(string.format("armed run=%s quiet=%d failsafe=%d", RUN, QUIET,
  FAILSAFE))
return "EXP-0034 armed " .. RUN
