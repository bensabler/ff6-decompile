-- bridge.lua — Mesen 2 instrumentation + file-based control bridge for the
-- ff6-decompile project.
--
-- Purpose (docs/sessions/08_OPEN_QUESTIONS.md #1): log every execution of the routine
-- at ROM CPU $C10DF3 with registers, P, DB and the caller's return address.
-- Also exposes a tiny file-based command interface so the session can be
-- driven externally while the GUI stays visible:
--   out/cmd.txt   — one command per write (file is consumed after execution)
--   out/resp.txt  — response to the last command
--   out/events.log— exec-callback log lines
--   out/api.txt   — dump of the emu.* Lua API surface (written once at load)
--
-- Commands:
--   screenshot                     -> out/screen.png
--   state                          -> CPU state + watched WRAM to resp.txt
--   read <memtype> <hexaddr> <len> -> hex dump to resp.txt (memtype: wram|cpu)
--   press <buttons> <frames>       -> hold buttons (comma list) for N frames
--   eval <lua>                     -> run lua chunk, result to resp.txt

-- Output directory: relative to Mesen's working directory (launch from the
-- repository root), overridable via mesen/bridge_config.lua returning
-- { out = "/abs/path/" }. No personal absolute paths in the script.
local OUT = "mesen/out/"
do
  local ok, cfg = pcall(dofile, "mesen/bridge_config.lua")
  if ok and type(cfg) == "table" and type(cfg.out) == "string" then OUT = cfg.out end
end

local function append(path, text)
  local f = io.open(path, "a")
  if f then f:write(text) f:close() end
end

local function writeFile(path, text, mode)
  local f = io.open(path, mode or "w")
  if f then f:write(text) f:close() end
end

local function serialize(v, depth)
  depth = depth or 0
  if depth > 3 then return "<deep>" end
  if type(v) == "table" then
    local keys = {}
    for k in pairs(v) do keys[#keys + 1] = tostring(k) end
    table.sort(keys)
    local parts = {}
    for _, k in ipairs(keys) do
      parts[#parts + 1] = k .. "=" .. serialize(v[tostring(k)] or v[tonumber(k)] or v[k], depth + 1)
    end
    return "{" .. table.concat(parts, ", ") .. "}"
  elseif type(v) == "number" then
    if v == math.floor(v) then return string.format("$%X", v) end
    return tostring(v)
  else
    return tostring(v)
  end
end

-- ---------------------------------------------------------------------------
-- API surface dump (adapt future scripts to what this build actually has)
-- ---------------------------------------------------------------------------
local function dumpApi()
  local lines = {}
  local names = {}
  for k in pairs(emu) do names[#names + 1] = k end
  table.sort(names)
  for _, k in ipairs(names) do
    local v = emu[k]
    if type(v) == "table" then
      lines[#lines + 1] = k .. " (table): " .. serialize(v, 2)
    else
      lines[#lines + 1] = k .. " (" .. type(v) .. ")"
    end
  end
  writeFile(OUT .. "api.txt", table.concat(lines, "\n") .. "\n")
end

-- ---------------------------------------------------------------------------
-- Exec logging for $C10DF3
-- ---------------------------------------------------------------------------
local TARGET = 0xC10DF3
local hitCount = 0

local function readCpu(addr)
  return emu.read(addr, emu.memType.snesMemory)
end

local function onTargetExec()
  hitCount = hitCount + 1
  -- Log the first 3 hits, then every 1000th; Session 002 established the
  -- per-frame cadence, so this is now just a liveness signal.
  if hitCount > 3 and hitCount % 1000 ~= 0 then return end
  local ok, err = pcall(function()
    local cpu = emu.getCpuState(emu.cpuType.snes)
    local sp = cpu.sp
    -- RTS-style return address: 16-bit little-endian at SP+1, +1.
    local ret = ((readCpu(sp + 2) << 8) | readCpu(sp + 1)) + 1
    append(OUT .. "events.log", string.format(
      "hit=%d frame=%d ret=$%02X:%04X A=$%04X X=$%04X Y=$%04X SP=$%04X PS=$%02X(m=%d,x=%d,e=%s) K=$%02X DB=$%02X D=$%04X | wram2E78=$%04X wram2EB5=$%04X\n",
      hitCount, emu.getState().frameCount or -1,
      cpu.k, ret & 0xFFFF,
      cpu.a, cpu.x, cpu.y, sp,
      cpu.ps, (cpu.ps >> 5) & 1, (cpu.ps >> 4) & 1, tostring(cpu.emulationMode),
      cpu.k, cpu.dbr, cpu.d,
      emu.readWord(0x2E78, emu.memType.snesWorkRam),
      emu.readWord(0x2EB5, emu.memType.snesWorkRam)))
    if hitCount <= 3 then
      -- Raw stack snippet to reconstruct the call chain by hand.
      local parts = {}
      for i = 1, 24 do
        parts[#parts + 1] = string.format("%02X", readCpu(sp + i))
      end
      append(OUT .. "events.log", string.format(
        "    stack[SP+1..SP+24] @ $%04X: %s\n", sp + 1, table.concat(parts, " ")))
    end
  end)
  if not ok then
    append(OUT .. "events.log", "hit=" .. hitCount .. " ERROR: " .. tostring(err) .. "\n")
  end
end

-- ---------------------------------------------------------------------------
-- Write-watcher on the source HP array ($2E78-$2E7F, WRAM-relative)
-- Goal (docs/08_OPEN_QUESTIONS.md #1): find the code that computes HP
-- changes. Aggregates writers by PC; first write from a new PC gets a full
-- register + stack capture in events.log.
-- ---------------------------------------------------------------------------
local hpWriters = {}
local hpWriteTotal = 0

local function onHpWrite(addr, value)
  hpWriteTotal = hpWriteTotal + 1
  local ok, err = pcall(function()
    local cpu = emu.getCpuState(emu.cpuType.snes)
    local pc = (cpu.k << 16) | cpu.pc
    local w = hpWriters[pc]
    local frame = emu.getState().frameCount
    if w then
      w.count = w.count + 1
      w.lastAddr, w.lastVal, w.lastFrame = addr, value, frame
      return
    end
    hpWriters[pc] = { count = 1, lastAddr = addr, lastVal = value, lastFrame = frame }
    local parts = {}
    for i = 1, 32 do
      parts[#parts + 1] = string.format("%02X", readCpu(cpu.sp + i))
    end
    append(OUT .. "events.log", string.format(
      "NEW-HP-WRITER pc=$%06X addr=%s val=%s A=$%04X X=$%04X Y=$%04X SP=$%04X PS=$%02X DB=$%02X frame=%d\n    stack[SP+1..+32]: %s\n",
      pc, tostring(addr), tostring(value), cpu.a, cpu.x, cpu.y, cpu.sp,
      cpu.ps, cpu.dbr, frame, table.concat(parts, " ")))
  end)
  if not ok then
    append(OUT .. "events.log", "hpwrite ERROR: " .. tostring(err) .. "\n")
  end
end

-- ---------------------------------------------------------------------------
-- Input injection
-- ---------------------------------------------------------------------------
local inputHold = nil -- {buttons = {a=true,...}, framesLeft = N}

local function onInputPoll()
  if inputHold and inputHold.framesLeft > 0 then
    emu.setInput(inputHold.buttons, 0)
    inputHold.framesLeft = inputHold.framesLeft - 1
    if inputHold.framesLeft == 0 then inputHold = nil end
  end
end

local function parseButtons(spec)
  local b = {}
  for name in string.gmatch(spec, "[^,]+") do b[name] = true end
  return b
end

-- ---------------------------------------------------------------------------
-- Command polling (every 10 frames)
-- ---------------------------------------------------------------------------
local WATCH = {
  { "wram $2E78..$2E97 (source)", 0x2E78, 0x20 },
  { "wram $2EB5..$2EBC (rec0)", 0x2EB5, 8 },
  { "wram $2ED5..$2EDC (rec1)", 0x2ED5, 8 },
}

local function hexDump(memType, base, len)
  local parts = {}
  for i = 0, len - 1 do
    parts[#parts + 1] = string.format("%02X", emu.read(base + i, memType))
  end
  return table.concat(parts, " ")
end

local function doState()
  local st = emu.getState()
  local out = { "frame=" .. tostring(st.frameCount), "cpu=" .. serialize(st.cpu or st, 1) }
  for _, w in ipairs(WATCH) do
    out[#out + 1] = w[1] .. ": " .. hexDump(emu.memType.snesWorkRam, w[2], w[3])
  end
  out[#out + 1] = "targetHits=" .. hitCount
  return table.concat(out, "\n")
end

local function handleCommand(line)
  local cmd, rest = line:match("^(%S+)%s*(.*)$")
  if cmd == "screenshot" then
    local png = emu.takeScreenshot()
    writeFile(OUT .. "screen.png", png, "wb")
    return "ok screenshot " .. tostring(#png) .. " bytes frame=" .. tostring(emu.getState().frameCount)
  elseif cmd == "state" then
    return doState()
  elseif cmd == "read" then
    local mt, addr, len = rest:match("^(%S+)%s+(%S+)%s+(%S+)$")
    local memType = (mt == "wram") and emu.memType.snesWorkRam or emu.memType.snesMemory
    return hexDump(memType, tonumber(addr, 16), tonumber(len))
  elseif cmd == "press" then
    local spec, frames = rest:match("^(%S+)%s+(%S+)$")
    inputHold = { buttons = parseButtons(spec), framesLeft = tonumber(frames) }
    return "ok press " .. spec .. " for " .. frames .. " frames"
  elseif cmd == "writers" then
    local lines = { "totalWrites=" .. hpWriteTotal }
    for pc, w in pairs(hpWriters) do
      lines[#lines + 1] = string.format(
        "pc=$%06X count=%d lastAddr=%s lastVal=%s lastFrame=%d",
        pc, w.count, tostring(w.lastAddr), tostring(w.lastVal), w.lastFrame)
    end
    return table.concat(lines, "\n")
  elseif cmd == "loadstate" then
    local f = io.open(rest, "rb")
    if not f then return "cannot open " .. rest end
    local data = f:read("*a")
    f:close()
    -- loadSavestate is only legal inside a main-CPU exec callback; run it
    -- one-shot on the next executed instruction.
    local ref
    ref = emu.addMemoryCallback(function()
      emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
      local ok, err = pcall(emu.loadSavestate, data)
      append(OUT .. "events.log", "loadstate applied ok=" .. tostring(ok) .. " detail=" .. tostring(err) .. "\n")
    end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
    return "ok loadstate queued " .. tostring(#data) .. " bytes"
  elseif cmd == "probe" then
    local name = rest:match("^(%S+)$")
    if not name or name:find("[^%w%-_]") then return "bad probe name" end
    local ok2, err2 = pcall(dofile, "mesen/probes/" .. name .. ".lua")
    return ok2 and ("ok probe " .. name) or ("probe error: " .. tostring(err2))
  elseif cmd == "eval" then
    local fn, err = load(rest)
    if not fn then return "load error: " .. tostring(err) end
    local ok, res = pcall(fn)
    return ok and ("ok eval: " .. serialize(res, 1)) or ("eval error: " .. tostring(res))
  end
  return "unknown command: " .. tostring(line)
end

-- Command protocol v2: cmd.txt holds "<id> <command...>" (id = any token;
-- clients should increase it per command). resp.txt is written atomically
-- (tmp + rename) as "<id> ok|error\n<payload>". A bare command with no id
-- is accepted with id "-" for backward compatibility. Repeated delivery of
-- the same id is ignored (duplicate-command prevention). Every command and
-- response line is appended to commands.log for the experiment transcript.
local frameTick = 0
local lastId = nil
local function onEndFrame()
  frameTick = frameTick + 1
  if frameTick % 10 ~= 0 then return end
  local f = io.open(OUT .. "cmd.txt", "r")
  if not f then return end
  local line = f:read("*l")
  f:close()
  os.remove(OUT .. "cmd.txt")
  if not line or line == "" then return end
  local id, restLine = line:match("^(%S+)%s+(.*)$")
  if id and restLine and restLine ~= "" and not restLine:match("^%s*$") and id:match("^[%w%-%.]+$") and (id:match("%d") or id == "-") then
    -- id'd form
  else
    id, restLine = "-", line
  end
  if id ~= "-" and id == lastId then return end
  lastId = id
  append(OUT .. "commands.log", os.date("%H:%M:%S") .. " > " .. line .. "\n")
  local ok, res = pcall(handleCommand, restLine)
  local status = ok and "ok" or "error"
  local payload = ok and res or ("command error: " .. tostring(res))
  writeFile(OUT .. "resp.txt.tmp", id .. " " .. status .. "\n" .. payload .. "\n")
  os.remove(OUT .. "resp.txt")
  os.rename(OUT .. "resp.txt.tmp", OUT .. "resp.txt")
  append(OUT .. "commands.log", os.date("%H:%M:%S") .. " < " .. id .. " " .. status .. "\n")
end

-- ---------------------------------------------------------------------------
-- Startup
-- ---------------------------------------------------------------------------
dumpApi()
local romName = "?"
pcall(function() romName = emu.getRomInfo().name end)
append(OUT .. "events.log", "=== bridge loaded rom=" .. romName ..
  " at " .. os.date("%Y-%m-%d %H:%M:%S") .. " ===\n")

emu.addMemoryCallback(onTargetExec, emu.callbackType.exec, TARGET, TARGET, emu.cpuType.snes, emu.memType.snesMemory)
emu.addMemoryCallback(onHpWrite, emu.callbackType.write, 0x2E78, 0x2E7F, emu.cpuType.snes, emu.memType.snesWorkRam)
emu.addEventCallback(onEndFrame, emu.eventType.endFrame)
emu.addEventCallback(onInputPoll, emu.eventType.inputPolled)
emu.displayMessage("bridge", "ff6-decompile bridge active")
