-- EXP-0038: scheduled mines corridor to the first random encounter
-- (SCN-0001 milestone 06). Route = EXP-0037's 17 legs (unchanged) plus
-- nine corridor legs from the EXP-0038 recon: north (26,1C)->(26,0B),
-- east ->(28,0B), north ->(28,09), then one patrol cycle back down and
-- up to extend the step budget. The route stops short of (2A,09),
-- where recon tripped a scripted event (out of scope).
--
-- Battle 6 = the random encounter. On its entry the probe logs the
-- +$11E0 formation id, the staged +$3F44 record, position and leg,
-- then captures milestone 06 at entry+120 (WRAM + screenshot +
-- savestate, matching the milestone-03 pattern) and ends the run.
-- Completing all legs without a sixth battle is the bounded negative
-- result (logged NO-ENCOUNTER), not a route failure.
--
-- Controller fix vs EXP-0037: a random encounter can interrupt a WALK
-- leg (scripted battles only ever happened in pulse legs), so while a
-- battle owns the screen the leg's timeout window restarts
-- (legStart = fc), mirroring internal/scenario/route's Runner.
--
-- The EXP-0037 event-flag watch (span $1E80-$1EDF, JSONL + snapshots +
-- shadow integrity) rides along observationally.
--
-- Env: FF6_RUN names the run; FF6_STOP=1 stops emulation after the
-- final capture (headless runs); FF6_OUT overrides output resolution
-- (required under --testrunner, where debug.getinfo returns =[C]).
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
do
  local ok, info = pcall(debug.getinfo, 1, "S")
  if ok and type(info) == "table" and type(info.source) == "string" then
    local dir = info.source:gsub("^@", ""):match("^(.*/)")
    if dir then OUT = dir .. "../out/" end -- <repo>/mesen/probes/../out/
  end
  local env = os.getenv and os.getenv("FF6_OUT")
  if env then OUT = env end
end
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local STOP_AT_END = (os.getenv and os.getenv("FF6_STOP")) == "1"
local EVD = OUT .. "../../local_artifacts/experiments/EXP-0038/"
local MS06 = OUT .. "../../local_artifacts/scenarios/SCN-0001/06-random-encounter/"

-- phase 0 (EXP-0034 schedule) - identical to EXP-0036
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

-- event-flag arrays (CONTRA-0002)
local FLAG_LO, FLAG_HI = 0x1E80, 0x1EDF
local SNAP_MILESTONES = {
  { frame = 5200,  label = "ms-00-new-game" },
  { frame = 15000, label = "ms-01-opening-cinematic" },
  { frame = 30000, label = "ms-02-narshe-entry" },
  { frame = 31677, label = "ms-03-first-scripted-battle" },
  { frame = 46375, label = "ms-04-free-movement" },
}
local SNAP_PERIOD = 5000

-- kind: "walk" (dir + axis target), "pulse" (A presses + condition)
-- cmp: "ge" / "le" applied to the axis byte
local ROUTE = {
  { n = 1,  kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1B, timeout = 900 },
  { n = 2,  kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x27, timeout = 900 },
  { n = 3,  kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x1E, timeout = 900 },
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
  { n = 16, kind = "walk",  dir = "up",    axis = X, cmp = "ge", target = 0x26, timeout = 1800 },
  { n = 17, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x1C, timeout = 900 },
  -- EXP-0038 mines corridor (recon 2026-08-01). The random encounter is
  -- expected to interrupt one of these legs; legs 21-26 patrol the same
  -- corridor to extend the step budget without reaching the (2A,09)
  -- scripted-event trigger.
  { n = 18, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x0B, timeout = 1200 },
  { n = 19, kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x28, timeout = 900 },
  { n = 20, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x09, timeout = 900 },
  { n = 21, kind = "walk",  dir = "down",  axis = Y, cmp = "ge", target = 0x0B, timeout = 900 },
  { n = 22, kind = "walk",  dir = "left",  axis = X, cmp = "le", target = 0x26, timeout = 900 },
  { n = 23, kind = "walk",  dir = "down",  axis = Y, cmp = "ge", target = 0x14, timeout = 900 },
  { n = 24, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x0B, timeout = 1200 },
  { n = 25, kind = "walk",  dir = "right", axis = X, cmp = "ge", target = 0x28, timeout = 900 },
  { n = 26, kind = "walk",  dir = "up",    axis = Y, cmp = "le", target = 0x09, timeout = 900 },
}

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP38 " .. s .. "\n") f:close() end
end

local function jsonl(s)
  local f = io.open(EVD .. RUN .. "-flags.jsonl", "a")
  if f then f:write(s .. "\n") f:close() end
end

local function snapfile(s)
  local f = io.open(EVD .. RUN .. "-snapshots.log", "a")
  if f then f:write(s .. "\n") f:close() end
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

-- Screenshots are NOT an assertion channel (CEN-QUIRK-0002: frame capture
-- is not byte-stable while WRAM is), so a capture failure must never break
-- a run. EXP-0037 dropped EXP-0036's shot() helper because it never took
-- screenshots; EXP-0038's milestone capture called it anyway, so run1 hit
-- "attempt to call a nil value" inside the milestone callback and never
-- reached finishRun/emu.stop. Restored here, hardened, and always pcall'd.
local function shot(path)
  local ok, png = pcall(emu.takeScreenshot)
  if not ok or type(png) ~= "string" or #png == 0 then
    log("SHOT-SKIPPED " .. path .. " reason=" .. tostring(png))
    return
  end
  local f = io.open(path, "wb")
  if f then f:write(png) f:close() end
end

-- Also dropped by EXP-0037 and needed here: emu.createSavestate, like
-- loadSavestate, is only legal inside a main-CPU exec callback, so it
-- runs one-shot on the next executed instruction. run2 lost its
-- milestone savestate to this being nil (the pcall contained it and the
-- run still completed, which is why the guard is there).
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
-- event-flag instrumentation
-- ---------------------------------------------------------------------
local shadow = {}          -- offset -> last known value
local shadowReady = false
local legIdxRef = function() return 0 end -- rebound after controller setup
local battleRef = function() return false, 0 end
local changeCount, rewriteCount = 0, 0
local rewriteSeen = {}     -- "pc:addr" -> count
local mismatchCount = 0

local function arraySnapshot()
  local a = {}
  for _, base in ipairs({ 0x1E80, 0x1EA0, 0x1EC0 }) do
    local c = {}
    for i = 0, 0x1F do c[#c + 1] = string.format("%02X", rd(base + i)) end
    a[#a + 1] = string.format("%04X=%s", base, table.concat(c))
  end
  return table.concat(a, " ")
end

local function snap(label, fc)
  snapfile(string.format("SNAP %s frame=%d %s", label, fc, arraySnapshot()))
end

local function initShadow(fc)
  for off = FLAG_LO, FLAG_HI do shadow[off] = rd(off) end
  shadowReady = true
  snap("arm", fc)
end

local function flagBits(mask, byteoff)
  local out = {}
  for bit = 0, 7 do
    if (mask & (1 << bit)) ~= 0 then
      out[#out + 1] = string.format("$%02X", byteoff * 8 + bit)
    end
  end
  return out
end

emu.addMemoryCallback(function(addr, value)
  local off = addr & 0x1FFFF
  if off < FLAG_LO or off > FLAG_HI or type(value) ~= "number" then
    log(string.format("FLAG-WATCH-ODD addr=%X value=%s", addr, tostring(value)))
    return
  end
  local fc = emu.getState().frameCount
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  local old = rd(off) -- callback fires before the store lands (EXP-0036)
  if shadowReady and shadow[off] ~= old then
    mismatchCount = mismatchCount + 1
    log(string.format(
      "INTEGRITY-MISMATCH frame=%d addr=$%04X shadow=%02X actual=%02X",
      fc, off, shadow[off], old))
  end
  shadow[off] = value

  local arrayBase = 0x1E80 + 0x20 * ((off - 0x1E80) // 0x20)
  local byteoff = off - arrayBase
  local cx, cy = coords()
  local inB, bCount = battleRef()
  local leg = legIdxRef()

  if value == old then
    rewriteCount = rewriteCount + 1
    local k = string.format("%06X:%04X", pc, off)
    rewriteSeen[k] = (rewriteSeen[k] or 0) + 1
    if rewriteSeen[k] <= 3 then
      log(string.format(
        "FLAG-REWRITE frame=%d pc=$%06X addr=$%04X val=%02X leg=%d battle=%s pos=%02X,%02X",
        fc, pc, off, value, leg, tostring(inB), cx, cy))
    end
    jsonl(string.format(
      '{"frame":%d,"pc":"$%06X","addr":"$%04X","array":"%04X","byte":%d,"old":"$%02X","new":"$%02X","set":[],"cleared":[],"same":true,"leg":%d,"battle":%s,"battles":%d,"pos":"%02X,%02X"}',
      fc, pc, off, arrayBase, byteoff, old, value, leg, tostring(inB), bCount, cx, cy))
    return
  end

  changeCount = changeCount + 1
  local setB = flagBits(value & ~old & 0xFF, byteoff)
  local clrB = flagBits(old & ~value & 0xFF, byteoff)
  local function q(t)
    if #t == 0 then return "" end
    return '"' .. table.concat(t, '","') .. '"'
  end
  log(string.format(
    "FLAG-WRITE frame=%d pc=$%06X addr=$%04X %02X->%02X array=%04X set=[%s] clr=[%s] leg=%d battle=%s pos=%02X,%02X",
    fc, pc, off, old, value, arrayBase, table.concat(setB, " "),
    table.concat(clrB, " "), leg, tostring(inB), cx, cy))
  jsonl(string.format(
    '{"frame":%d,"pc":"$%06X","addr":"$%04X","array":"%04X","byte":%d,"old":"$%02X","new":"$%02X","set":[%s],"cleared":[%s],"same":false,"leg":%d,"battle":%s,"battles":%d,"pos":"%02X,%02X"}',
    fc, pc, off, arrayBase, byteoff, old, value, q(setB), q(clrB),
    leg, tostring(inB), bCount, cx, cy))
end, emu.callbackType.write, FLAG_LO, FLAG_HI, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- ---------------------------------------------------------------------
-- battle detection (EXP-0034 signature, re-arming)
-- ---------------------------------------------------------------------
local inBattle, battleCount, battleFrame = false, 0, nil
local zeroSince = nil
local battleEdgeStart, battleEdgeEnd = false, false
battleRef = function() return inBattle, battleCount end

local encounterSeen = false
local finishRun -- defined with the route controller below

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
    snap(string.format("battle-%d-entry", battleCount), battleFrame)
    -- Battle 6 is the random encounter: capture milestone 06 at
    -- entry+120 (the milestone-03 pattern) and end the run.
    if battleCount == 6 and not encounterSeen then
      encounterSeen = true
      log(string.format(
        "ENCOUNTER frame=%d 11E0=$%04X leg=%d pos=%02X,%02X",
        battleFrame, emu.readWord(0x11E0, emu.memType.snesWorkRam),
        legIdxRef(), cx, cy))
      local target = battleFrame + 120
      local ref
      ref = emu.addEventCallback(function()
        local f2 = emu.getState().frameCount
        if f2 < target then return end
        emu.removeEventCallback(ref, emu.eventType.endFrame)
        -- Artifact writing must not be able to strand the run: any error
        -- here is logged, and finishRun (which stops headless emulation)
        -- runs regardless. run1 stranded exactly this way.
        local okArt, artErr = pcall(function()
          wramDump(MS06 .. RUN .. "-wram.bin")
          shot(MS06 .. RUN .. "-screen.png")
          mkstate(MS06 .. RUN .. "-06.mss", "MS-06")
          snap("ms-06-random-encounter", f2)
        end)
        if not okArt then
          log("MILESTONE-06-ARTIFACT-ERROR " .. tostring(artErr))
        end
        log(string.format("MILESTONE 06 frame=%d battles=%d artifacts=%s",
          f2, battleCount, tostring(okArt)))
        finishRun(f2)
      end, emu.eventType.endFrame)
    end
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- ---------------------------------------------------------------------
-- route controller (EXP-0036, unchanged)
-- ---------------------------------------------------------------------
local legIdx = 0
local legStart = nil
local neutralUntil = 0
local routeDone, routeFailed = false, false
local pulseBase = nil
legIdxRef = function() return legIdx end

local function beginLeg(i, fc)
  legIdx = i
  legStart = fc
  pulseBase = fc
  local L = ROUTE[i]
  local cx, cy = coords()
  log(string.format("LEG %d BEGIN frame=%d kind=%s dir=%s pos=%02X,%02X evf=%02X alias=[%s]",
    L.n, fc, L.kind, tostring(L.dir), cx, cy, rd(EVFLAG), aliasStr()))
  -- snapshot label carries the SAMPLING frame: beginLeg's fc is the
  -- scheduled start (current+NEUTRAL), and the $1EB6 working bits can
  -- toggle inside that gap (seen in gui1's leg-8-begin replay check)
  snap(string.format("leg-%d-begin", L.n), emu.getState().frameCount)
end

local function endLeg(fc, why)
  local L = ROUTE[legIdx]
  local cx, cy = coords()
  log(string.format("LEG %d END frame=%d dur=%d why=%s pos=%02X,%02X evf=%02X alias=[%s]",
    L.n, fc, fc - legStart, why, cx, cy, rd(EVFLAG), aliasStr()))
  snap(string.format("leg-%d-end", L.n), fc)
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

function finishRun(fc)
  local cx, cy = coords()
  snap("final", fc)
  wramDump(EVD .. RUN .. "-wram-final.bin")
  local rw = {}
  for k, n in pairs(rewriteSeen) do
    rw[#rw + 1] = string.format("%s x%d", k, n)
  end
  table.sort(rw)
  log(string.format(
    "TOTALS changes=%d rewrites=%d mismatches=%d rewrite-sites=[%s]",
    changeCount, rewriteCount, mismatchCount, table.concat(rw, " ")))
  log(string.format(
    "FINAL frame=%d pos=%02X,%02X evf=%02X alias=[%s] battles=%d",
    fc, cx, cy, rd(EVFLAG), aliasStr(), battleCount))
  log("RUN-COMPLETE " .. RUN)
  if STOP_AT_END then pcall(function() emu.stop() end) end
end

emu.addEventCallback(function()
  local fc = emu.getState().frameCount

  if not shadowReady then initShadow(fc) end

  -- milestone + periodic snapshots (all phases)
  for _, m in ipairs(SNAP_MILESTONES) do
    if not m.done and fc >= m.frame then
      m.done = true
      snap(m.label, fc)
    end
  end
  if fc % SNAP_PERIOD == 0 then snap("periodic", fc) end

  -- Battle-end detection runs in EVERY phase (EXP-0036 v1 lesson).
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
        snap(string.format("battle-%d-end", battleCount), fc)
      end
    else
      zeroSince = nil
    end
  end

  if fc < ROUTE_START or routeDone or routeFailed or encounterSeen then return end

  if legIdx == 0 then beginLeg(1, fc) return end
  if fc < neutralUntil then return end

  local L = ROUTE[legIdx]
  -- Position-tracking legs cannot advance while a battle owns the screen
  -- (position bytes are not field-meaningful there - EXP-0036).
  local blocked = inBattle and L.kind == "walk" and L.until_ == nil
  -- A random encounter can interrupt a walk leg (new in EXP-0038; the
  -- scripted battles only ever fell in pulse legs). While the battle
  -- owns the screen the leg must not time out for time spent fighting:
  -- restart its window, mirroring internal/scenario/route's Runner.
  if blocked then legStart = fc end

  if not blocked and legSatisfied(L, fc) then
    endLeg(fc, "satisfied")
    battleEdgeStart, battleEdgeEnd = false, false
    if legIdx == #ROUTE then
      routeDone = true
      log(string.format("ROUTE-COMPLETE frame=%d NO-ENCOUNTER budget exhausted", fc))
      local target = fc + SETTLE_AFTER_TRANSITION
      local ref
      ref = emu.addEventCallback(function()
        local f2 = emu.getState().frameCount
        if f2 < target then return end
        emu.removeEventCallback(ref, emu.eventType.endFrame)
        finishRun(f2)
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
    log("RUN-FAILED " .. RUN .. " at leg " .. L.n)
    if STOP_AT_END then pcall(function() emu.stop() end) end
    return
  end
end, emu.eventType.endFrame)

-- ---------------------------------------------------------------------
-- input: phase 0 (EXP-0034) until ROUTE_START, then the leg controller
-- (EXP-0036, unchanged)
-- ---------------------------------------------------------------------
emu.addEventCallback(function()
  local fc = emu.getState().frameCount

  if fc >= T0 and fc <= T1 then
    local on = math.floor(fc / 6) % 2 == 0
    if on then emu.setInput({ start = true, a = true }, 0) end
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
    if battleFrame and fc >= battleFrame + 60
      and (fc - battleFrame) % BATP < BATHOLD then
      emu.setInput({ a = true }, 0)
    end
    return
  end
  local btn = {}
  if L.dir then btn[L.dir] = true end
  if L.kind == "pulse" and (fc - pulseBase) % 90 < 10 then btn.a = true end
  if next(btn) then emu.setInput(btn, 0) end
end, emu.eventType.inputPolled)

log(string.format("armed run=%s legs=%d routeStart=%d flagspan=%04X-%04X stop=%s",
  RUN, #ROUTE, ROUTE_START, FLAG_LO, FLAG_HI, tostring(STOP_AT_END)))
return "EXP-0038 armed " .. RUN
