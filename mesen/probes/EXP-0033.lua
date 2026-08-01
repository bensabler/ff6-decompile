-- EXP-0033: golden route segment 3 — Guard battle victory → free field
-- movement (SCN-0001 milestone 04). Extends the EXP-0032 schedule.
--
-- Schedule (absolute frames from power-on):
--   2500-4200    start+a edge toggling (title)
--   15000/30000  milestones 01/02 (unchanged; captured for continuity)
--   31000-46000  up held + A every 240 frames (dialogue/walk chain)
--   battle       detected via +$3B18 write from $C22800-$C22FFF
--   battle+60..  A every 90 frames (battle confirm cadence) until both
--                enemy HP slots reach zero. The enemies of formation 2
--                occupy slots 6 and 7 (HP words $3C00/$3C02): slots 0-2
--                are the party (63/68/70 on screen), and slots 6/7 hold
--                HP 40 / MP 15, matching monster record 0's +$08/+$0A.
--                (v1 watched slot 0 = a party member and never fired.)
--   end+600      milestone 04-free-movement (post-battle settle), then
--                the walk+A cadence resumes so the state is genuinely
--                player-controlled when captured
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local ROOT = OUT .. "../../local_artifacts/scenarios/SCN-0001/"
local EVD = OUT .. "../../local_artifacts/experiments/EXP-0033/"

local T0, T1 = 2500, 4200
local WALK0, WALK1 = 31000, 46000
local METP, METHOLD = 240, 12
local BATP, BATHOLD = 90, 10
-- v2 captured milestone 04 at end+600 and landed on the reward window
-- ("Got NN Exp. point(s)"): the victory sequence is itself a chain of
-- input-waiting windows. end+3000 with the walk+A cadence running.
local SETTLE = 3000
local FAILSAFE = 80000

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP33 " .. s .. "\n") f:close() end
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

_G.exp33_battle = false
_G.exp33_end = nil
local battleFrame = nil

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
  elseif exp33_end then
    -- post-battle: resume walking so the captured state is player-driven
    if fc >= exp33_end + 200 then
      local pressA = (fc - exp33_end) % METP < METHOLD
      emu.setInput({ up = true, a = pressA or nil }, 0)
    end
  elseif battleFrame then
    -- in battle: confirm cadence
    if fc >= battleFrame + 60 and (fc - battleFrame) % BATP < BATHOLD then
      emu.setInput({ a = true }, 0)
    end
  elseif fc >= WALK0 and fc <= WALK1 then
    local pressA = (fc - WALK0) % METP < METHOLD
    emu.setInput({ up = true, a = pressA or nil }, 0)
  end
end, emu.eventType.inputPolled)

-- battle-init detector (EXP-0032 signature)
emu.addMemoryCallback(function(addr, value)
  if exp33_battle then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  if pc >= 0xC22800 and pc <= 0xC22FFF then
    exp33_battle = true
    battleFrame = emu.getState().frameCount
    log(string.format("BATTLE-DETECT frame=%d 11E0=$%04X", battleFrame,
      emu.readWord(0x11E0, emu.memType.snesWorkRam)))
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- reward-window writer capture: first hit per (PC,addr) on the field
-- character-record block ~+$1600 (EXP-0027 located current HP/MP at
-- +$1609/+$160D there) during the post-victory window only. Victory
-- processing must write experience/level state somewhere; this watch
-- says whether that block is a destination. Bounded: log, do not
-- decode. (v2 watched +$1860-+$18FF and saw nothing.)
local rewardWriters = {}
local rewardArmed = false
emu.addMemoryCallback(function(addr, value)
  if not rewardArmed then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  local k = string.format("%06X:%04X", pc, addr)
  if rewardWriters[k] then return end
  rewardWriters[k] = true
  log(string.format("REWARD-WRITER pc=$%06X addr=$%04X val=%s frame=%d",
    pc, addr, tostring(value), emu.getState().frameCount))
end, emu.callbackType.write, 0x1600, 0x16FF, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- battle end: enemy slot-0 current HP ($3BF4) reaches zero after init
local did04, failed = false, false
local nextShot = 31000
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if battleFrame and not exp33_end and fc > battleFrame + 120 then
    local hp6 = emu.readWord(0x3C00, emu.memType.snesWorkRam)
    local hp7 = emu.readWord(0x3C02, emu.memType.snesWorkRam)
    if hp6 == 0 and hp7 == 0 then
      exp33_end = fc
      rewardArmed = true
      log(string.format("BATTLE-END frame=%d (enemy slots 6/7 HP=0)", fc))
      shot(EVD .. RUN .. "-battleend.png")
    end
  end
  if fc >= nextShot then
    nextShot = nextShot + 400
    if not did04 then shot(EVD .. RUN .. "-beat-" .. fc .. ".png") end
  end
  if not did04 and exp33_end and fc >= exp33_end + SETTLE then
    did04 = true
    rewardArmed = false
    milestone(ROOT .. "04-free-movement/", "04")
    log("RUN-COMPLETE " .. RUN)
  end
  if not did04 and not failed and fc >= FAILSAFE then
    failed = true
    shot(EVD .. RUN .. "-failsafe.png")
    log(string.format("FAILSAFE frame=%d battle=%s end=%s", fc,
      tostring(battleFrame), tostring(exp33_end)))
  end
end, emu.eventType.endFrame)

log(string.format("armed run=%s batp=%d settle=%d failsafe=%d",
  RUN, BATP, SETTLE, FAILSAFE))
return "EXP-0033 armed " .. RUN
