-- EXP-0032: golden route segment 2 — Narshe entry to the first scripted
-- battle (SCN-0001 milestones 01/02/03). Extends the EXP-0031 schedule.
--
-- Schedule (absolute frames from power-on):
--   2500-4200   start+a edge toggling (segment 1, unchanged)
--   15000       milestone 01-opening-cinematic (auto-run interior)
--   30000       milestone 02-narshe-entry (the EXP-0031 stall state)
--   31000-46000 UP held continuously AND an A press (12 polls) every 240
--               frames. Both are required: recon showed the beat is a
--               chain of input-waiting dialogue boxes separated by
--               player-controlled walking north through the Narshe gate
--               (v1's A-only metronome cleared boxes but never moved;
--               walking alone stalls at the guard box).
--   detect      all input stops (state-driven, deterministic)
--   detect+120  milestone 03-first-scripted-battle
--
-- Battle-init detector: first write to WRAM $3B18-$3BB7 from a PC inside
-- the Confirmed populate range $C22800-$C22FFF (EXP-0028: party cluster
-- $C22899-$C2290F, enemy loader $C22CA3-$C22D8C). Other writers to the
-- range are logged (first hit per PC) but do not trigger.
-- At detection: logs +$11E0 (formation id word) and +$3F44..+$3F53 (the
-- staged formation record, EXP-0030).
local OUT = _G.FF6_OUT_DIR or "mesen/out/"
local RUN = (os.getenv and os.getenv("FF6_RUN")) or "run1"
local ROOT = OUT .. "../../local_artifacts/scenarios/SCN-0001/"
local EVD = OUT .. "../../local_artifacts/experiments/EXP-0032/"

local T0, T1 = 2500, 4200
local M01, M02 = 15000, 30000
local WALK0, WALK1 = 31000, 46000
local METP, METHOLD = 240, 12
local FAILSAFE = 60000

local function log(s)
  local f = io.open(OUT .. "events.log", "a")
  if f then f:write("EXP32 " .. s .. "\n") f:close() end
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

-- savestate must be created inside a main-CPU exec callback (one-shot)
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

-- input: segment-1 toggling + the A metronome
_G.exp32_battle = false
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
  elseif fc >= WALK0 and fc <= WALK1 and not exp32_battle then
    local pressA = (fc - WALK0) % METP < METHOLD
    emu.setInput({ up = true, a = pressA or nil }, 0)
    if (fc - WALK0) % METP == 0 then
      log(string.format("WALK-A frame=%d", fc))
    end
  end
end, emu.eventType.inputPolled)

-- battle-init detector
local otherWriters = {}
local detectFrame = nil
emu.addMemoryCallback(function(addr, value)
  if exp32_battle then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  if pc >= 0xC22800 and pc <= 0xC22FFF then
    exp32_battle = true
    detectFrame = emu.getState().frameCount
    local id = emu.readWord(0x11E0, emu.memType.snesWorkRam)
    local staged = {}
    for i = 0, 15 do
      staged[#staged + 1] = string.format("%02X",
        emu.read(0x3F44 + i, emu.memType.snesWorkRam))
    end
    log(string.format(
      "BATTLE-DETECT frame=%d pc=$%06X addr=$%04X val=%s 11E0=$%04X staged3F44=%s",
      detectFrame, pc, addr, tostring(value), id,
      table.concat(staged, " ")))
  else
    local k = string.format("%06X", pc)
    if not otherWriters[k] then
      otherWriters[k] = true
      log(string.format("OTHER-3B18-WRITER pc=$%06X frame=%d", pc,
        emu.getState().frameCount))
    end
  end
end, emu.callbackType.write, 0x3B18, 0x3BB7, emu.cpuType.snes,
  emu.memType.snesWorkRam)

-- frame-driven captures
local did01, did02, did03, failed = false, false, false, false
local nextShot = 30000
emu.addEventCallback(function()
  local fc = emu.getState().frameCount
  if not did01 and fc >= M01 then
    did01 = true
    milestone(ROOT .. "01-opening-cinematic/", "01")
  end
  if not did02 and fc >= M02 then
    did02 = true
    milestone(ROOT .. "02-narshe-entry/", "02")
  end
  if fc >= nextShot and not did03 then
    nextShot = nextShot + 2000
    shot(EVD .. RUN .. "-beat-" .. fc .. ".png")
  end
  if not did03 and detectFrame and fc >= detectFrame + 120 then
    did03 = true
    milestone(ROOT .. "03-first-scripted-battle/", "03")
    log("RUN-COMPLETE " .. RUN)
  end
  if not did03 and not failed and fc >= FAILSAFE then
    failed = true
    shot(EVD .. RUN .. "-failsafe.png")
    log("FAILSAFE frame=" .. fc .. " (no battle detect)")
  end
end, emu.eventType.endFrame)

log(string.format("armed run=%s met0=%d period=%d failsafe=%d",
  RUN, MET0, METP, FAILSAFE))
return "EXP-0032 armed " .. RUN
