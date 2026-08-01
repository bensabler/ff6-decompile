-- EXP-0036: scheduled route from milestone 04 to the mines interior
-- (SCN-0001 milestone 05). Replaces the fixed-cadence walking of
-- segments 2-3b with a coordinate-driven route controller.
--
-- Phases:
--   0-46375   the EXP-0034 schedule verbatim (power-on -> milestone 04):
--             title edge toggling, then walk+A through the four scripted
--             battles, gated by the re-arming battle detector.
--   46375+    the 17-leg route controller (see ROUTE below).
--
-- Leg advancement is state-driven: a walk leg ends when its coordinate
-- target is satisfied (overshoot-tolerant), an A-pulse leg ends on a
-- battle edge or after a fixed settle. Timeouts do not retry or correct
-- - they fail the run and name the divergent leg.
--
-- Coordinates: WRAM:+$00AF (X tile) / +$00B0 (Y tile) - EXP-0035
-- candidates. +$1EA5 is logged for observation only: CONTRA-0002 decoded
-- it as byte 5 of the event-flag bit array based at +$1EA0 (set via
-- `ORA $C0BAFC,X / STA $1EA0,Y`), NOT a map indicator. Log lines before
-- that resolution used the tag `MAPC-WRITE` and the field `map=`.
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local ROOT = OUT .. "../../local_artifacts/scenarios/SCN-0001/"
local EVD = OUT .. "../../local_artifacts/experiments/EXP-0036/"

-- phase 0 (EXP-0034 schedule)
local T0, T1 = 2500, 4200
local WALK0 = 31000
local METP, METHOLD = 240, 12
local BATP, BATHOLD = 90, 10
local ENDHOLD = 30
local ROUTE_START = 46375
local NEUTRAL = 8
local SETTLE_AFTER_TRANSITION = 600

local X, Y, EVFLAG = 0x00AF, 0x00B0, 0x1EA5
local ALIASES = { 0x0541, 0x0543, 0x0545, 0x087A }

-- kind: "walk" (dir + axis target), "pulse" (A presses + condition)
-- cmp: "ge" / "le" applied to the axis byte
local ROUTE = {
  { n = 1,  kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1B, timeout = 900 },
  { n = 2,  kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x27, timeout = 900 },
  { n = 3,  kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1E, timeout = 900 },
  -- Two failed runs placed this trigger. Standing on (1E,27) tapping A
  -- did nothing, and neither did holding right there: the guard dialogue
  -- appears one tile north, at (1E,25). EXP-0035's condensed leg table
  -- had dropped this climb; its recon log has it.
  { n = 4,  kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x25, timeout = 900 },
  { n = 5,  kind = "pulse", dir = "right", until_ = "battle_start", timeout = 1800 },
  { n = 6,  kind = "pulse", until_ = "battle_end",   timeout = 7200 },
  { n = 7,  kind = "pulse", until_ = "elapsed", dur = 900,  timeout = 1200 },
  { n = 8,  kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x21, timeout = 900 },
  { n = 9,  kind = "walk",  dir = "left",  axis = X, cmp = "le", target = 0x1D, timeout = 900 },
  { n = 10, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x1E, timeout = 900 },
  { n = 11, kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1E, timeout = 900 },
  { n = 12, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x18, timeout = 900 },
  { n = 13, kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1F, timeout = 900 },
  { n = 14, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x16, timeout = 900 },
  { n = 15, kind = "pulse", until_ = "elapsed", dur = 600, timeout = 900 },
  -- Detect the transition by the position jump, NOT by +$1EA5: it reaches
  -- $0D during the shaft dialogue while the party is still visibly on the
  -- exterior. CONTRA-0002 explains why - that byte is an event-flag byte,
  -- so it never marked the transition. X only reaches $26 past it.
  { n = 16, kind = "walk",  dir = "up",    axis = X, cmp = "ge", target = 0x26, timeout = 1800 },
  -- the transition lands at Y=$21/$20; walk on to the documented
  -- milestone-05 tile (26,1C) from EXP-0035's recon.
  { n = 17, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x1C, timeout = 900 },
}

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP36 " .. s .. "\n") f:close() end
end

local function rd(a) return emu.read(a, emu.memType.snesWorkRam) end

local function coords()
  return rd(X), rd(Y)
end

local function aliasStr()
  local p = {}
  for _, a in ipairs(ALIASES) do
    p[#p + 1] = string.format("%02X/%02X", rd(a), rd(a + 1))
  end
  return table.concat(p, " ")
end

local function wramDump(path)
  local parts = {}
  for base = 0, 0x1FFFF, 0x1000 do
    local c = {}
    for i = 0, 0xFFF do c[i + 1] = string.char(rd(base + i)) end
    parts[#parts + 1] = table.concat(c)
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

-- ---------------------------------------------------------------------
-- battle detection (EXP-0034 signature, re-arming)
-- ---------------------------------------------------------------------
local inBattle, battleCount, battleFrame = false, 0, nil
local zeroSince = nil
local battleEdgeStart, battleEdgeEnd = false, false

emu.addMemoryCallback(function()
  if inBattle then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  if pc >= 0xC22800 and pc <= 0xC22FFF then
    inBattle, battleEdgeStart = true, true
    battleCount = battleCount + 1
    battleFrame = emu.getState().frameCount
    zeroSince = nil
    local staged = {}
    for i = 0, 15 do
      staged[#staged + 1] = string.format("%02X", rd(0x3F44 + i))
    end
    local cx, cy = coords()
    log(string.format(
      "BATTLE %d ENTRY frame=%d 11E0=$%04X pos=%02X,%02X evf=%02X staged=%s",
      battleCount, battleFrame, emu.readWord(0x11E0, emu.memType.snesWorkRam),
      cx, cy, rd(EVFLAG), table.concat(staged, " ")))
    log(string.format("  alias-at-battle-entry: %s", aliasStr()))
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- ---------------------------------------------------------------------
-- +$1EA5 write watch: frame + PC for every change (event-flag byte;
-- see CONTRA-0002)
-- ---------------------------------------------------------------------
local mapPrev = nil
local mapChanged = false
emu.addMemoryCallback(function(addr, value)
  local prev = rd(EVFLAG)
  if value == prev then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  mapChanged = true
  log(string.format("EVFLAG-WRITE frame=%d pc=$%06X %02X->%s",
    emu.getState().frameCount, pc, prev, tostring(value)))
end, emu.callbackType.write, EVFLAG, EVFLAG, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- ---------------------------------------------------------------------
-- route controller
-- ---------------------------------------------------------------------
local legIdx = 0
local legStart = nil
local neutralUntil = 0
local routeDone, routeFailed = false, false
local pulseBase = nil

local function beginLeg(i, fc)
  legIdx = i
  legStart = fc
  pulseBase = fc
  local L = ROUTE[i]
  local cx, cy = coords()
  log(string.format("LEG %d BEGIN frame=%d kind=%s dir=%s pos=%02X,%02X evf=%02X alias=[%s]",
    L.n, fc, L.kind, tostring(L.dir), cx, cy, rd(EVFLAG), aliasStr()))
end

local function endLeg(fc, why)
  local L = ROUTE[legIdx]
  local cx, cy = coords()
  log(string.format("LEG %d END frame=%d dur=%d why=%s pos=%02X,%02X evf=%02X alias=[%s]",
    L.n, fc, fc - legStart, why, cx, cy, rd(EVFLAG), aliasStr()))
end

local function legSatisfied(L, fc)
  if L.kind == "walk" then
    if L.until_ == "watched_byte" then
      return rd(EVFLAG) == 0x0D
    end
    local v = rd(L.axis)
    if L.cmp == "ge" then return v >= L.target end
    return v <= L.target
  else
    if L.until_ == "battle_start" then return battleEdgeStart end
    if L.until_ == "battle_end" then return battleEdgeEnd end
    if L.until_ == "elapsed" then return fc - legStart >= L.dur end
  end
  return false
end

emu.addEventCallback(function()
  local fc = emu.getState().frameCount

  -- Battle-end detection runs in EVERY phase. (v1 gated this behind
  -- ROUTE_START, so after battle 1 at frame 31557 `inBattle` never
  -- cleared, the detector never re-armed for the opening's battles
  -- 2-4, and phase 0 held the battle cadence instead of walking.)
  if inBattle and battleFrame and fc > battleFrame + 120 then
    local allZero = true
    for slot = 3, 9 do
      if emu.readWord(0x3BF4 + slot * 2, emu.memType.snesWorkRam) ~= 0 then
        allZero = false break
      end
    end
    if allZero then
      zeroSince = zeroSince or fc
      if fc - zeroSince >= ENDHOLD then
        inBattle, battleEdgeEnd = false, true
        local cx, cy = coords()
        log(string.format("BATTLE %d END frame=%d pos=%02X,%02X alias=[%s]",
          battleCount, fc, cx, cy, aliasStr()))
      end
    else
      zeroSince = nil
    end
  end

  if fc < ROUTE_START or routeDone or routeFailed then return end

  if legIdx == 0 then beginLeg(1, fc) return end
  if fc < neutralUntil then return end

  local L = ROUTE[legIdx]
  -- Position-tracking legs cannot advance while a battle owns the screen:
  -- the position bytes are not field-meaningful there (observed reading
  -- 00,00 at battle end). Legs waiting on a battle edge still evaluate.
  local blocked = inBattle and L.kind == "walk" and L.until_ == nil

  if not blocked and legSatisfied(L, fc) then
    endLeg(fc, "satisfied")
    battleEdgeStart, battleEdgeEnd = false, false
    if legIdx == #ROUTE then
      routeDone = true
      log(string.format("ROUTE-COMPLETE frame=%d", fc))
      -- settle, then capture milestone 05
      local target = fc + SETTLE_AFTER_TRANSITION
      local ref
      ref = emu.addEventCallback(function()
        local f2 = emu.getState().frameCount
        if f2 < target then return end
        emu.removeEventCallback(ref, emu.eventType.endFrame)
        local cx, cy = coords()
        local dir = ROOT .. "05-mines-entry/"
        wramDump(dir .. RUN .. "-wram.bin")
        shot(dir .. RUN .. "-screen.png")
        mkstate(dir .. RUN .. "-05.mss", "MS-05")
        log(string.format(
          "MILESTONE 05 frame=%d pos=%02X,%02X evf=%02X alias=[%s] battles=%d",
          f2, cx, cy, rd(EVFLAG), aliasStr(), battleCount))
        log("RUN-COMPLETE " .. RUN)
      end, emu.eventType.endFrame)
      return
    end
    neutralUntil = fc + NEUTRAL
    beginLeg(legIdx + 1, fc + NEUTRAL)
    return
  end

  if fc - legStart > L.timeout and not blocked then
    routeFailed = true
    local cx, cy = coords()
    log(string.format(
      "LEG %d TIMEOUT frame=%d pos=%02X,%02X evf=%02X expected=%s alias=[%s]",
      L.n, fc, cx, cy, rd(EVFLAG),
      L.target and string.format("%s %02X", L.cmp, L.target) or tostring(L.until_),
      aliasStr()))
    shot(EVD .. RUN .. "-leg" .. L.n .. "-timeout.png")
    log("RUN-FAILED " .. RUN .. " at leg " .. L.n)
    return
  end
end, emu.eventType.endFrame)

-- ---------------------------------------------------------------------
-- input: phase 0 (EXP-0034) until ROUTE_START, then the leg controller
-- ---------------------------------------------------------------------
local togglingOn = false
emu.addEventCallback(function()
  local fc = emu.getState().frameCount

  if fc >= T0 and fc <= T1 then
    local on = math.floor(fc / 6) % 2 == 0
    if on then emu.setInput({ start = true, a = true }, 0) end
    togglingOn = on
    return
  end

  if fc < ROUTE_START then
    if inBattle then
      if battleFrame and fc >= battleFrame + 60
        and (fc - battleFrame) % BATP < BATHOLD then
        emu.setInput({ a = true }, 0)
      end
    elseif fc >= WALK0 then
      local pressA = (fc - WALK0) % METP < METHOLD
      emu.setInput({ up = true, a = pressA or nil }, 0)
    end
    return
  end

  if routeDone or routeFailed or legIdx == 0 or fc < neutralUntil then return end
  local L = ROUTE[legIdx]
  if inBattle then
    -- a battle owns the screen: fight rather than walk, whatever the leg
    if battleFrame and fc >= battleFrame + 60
      and (fc - battleFrame) % BATP < BATHOLD then
      emu.setInput({ a = true }, 0)
    end
    return
  end
  -- A leg may hold a direction, tap A, or both (leg 4 walks into the
  -- guard trigger while clearing its dialogue).
  local btn = {}
  if L.dir then btn[L.dir] = true end
  if L.kind == "pulse" and (fc - pulseBase) % 90 < 10 then btn.a = true end
  if next(btn) then emu.setInput(btn, 0) end
end, emu.eventType.inputPolled)

log(string.format("armed run=%s legs=%d routeStart=%d", RUN, #ROUTE, ROUTE_START))
return "EXP-0036 armed " .. RUN
