-- DMA trace. Load with the bridge command:  probe dma-trace
--
-- STATUS: UNEXERCISED. Written 2026-08-02 alongside internal/platform/snesdma.
-- It has never been run in Mesen. Until an operator session runs it and the
-- output is checked against a known transfer, treat every line it produces as
-- unverified. Do not cite it as evidence before then.
--
-- WHY IT EXISTS
--
-- The project has carried DMA doctrine for a long time with nothing beneath
-- it: /trace-dma, the dma-tracer skill, the dma-researcher agent and
-- TRACE_DMA.md all existed while no probe and no Go path read a single DMA
-- register. The savestate reader (Unit 12) sees channel state at one instant,
-- which EXP-0050 showed can be the leftovers of a finished transfer: the mines
-- state's channel 0 claims $7EC180 -> $2118, but VRAM matches that buffer on
-- only 18% of bytes. Configuration at rest is not a trace.
--
-- WHAT IT CAPTURES
--
-- A transfer begins when the CPU writes $420B (MDMAEN) or $420C (HDMAEN). On
-- that write this probe reads every enabled channel's $43n0-$43nA and logs the
-- raw register block, plus the PC that kicked it and the frame. That is the
-- moment the source address is still the one that matters.
--
-- The raw bytes are logged rather than a decoded summary, deliberately: the
-- decode belongs in internal/platform/snesdma, where it is tested, and a log
-- of raw registers can be re-decoded when the decoder improves. A log of
-- someone's interpretation cannot.
--
-- OUTPUT
--
-- mesen/out/dma-trace.log, one line per enabled channel per kick:
--
--   dma frame=<n> pc=$<6hex> reg=$420B mask=$<2hex> ch=<0-7> regs=<22 hex>
--
-- Parse it with snesdma.ParseTraceLine.

dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

local OUT = _G.FF6_OUT_DIR or "mesen/out/"
do
  local ok, info = pcall(debug.getinfo, 1, "S")
  if ok and type(info) == "table" and type(info.source) == "string" then
    local dir = info.source:gsub("^@", ""):match("^(.*/)")
    if dir then OUT = dir .. "../out/" end
  end
  local env = os.getenv and os.getenv("FF6_OUT")
  if env then OUT = env end
end

local LOG = OUT .. "dma-trace.log"

-- Channel n's registers are $4300 + n*$10, eleven of them ($43n0-$43nA).
local function channelRegs(ch)
  local base = 0x4300 + ch * 0x10
  local out = {}
  for i = 0, 10 do
    -- snesMemory is the CPU bus view, which is where $43xx lives.
    out[#out + 1] = string.format("%02X", emu.read(base + i, emu.memType.snesMemory))
  end
  return table.concat(out)
end

local function logKick(regName, mask)
  -- Nothing enabled is not a transfer; $420B is written with 0 routinely.
  if mask == 0 then return end
  local c = emu.getCpuState(emu.cpuType.snes)
  local pc = (c.k << 16) | c.pc
  local frame = emu.getState().frameCount
  local f = io.open(LOG, "a")
  if not f then return end
  for ch = 0, 7 do
    if (mask & (1 << ch)) ~= 0 then
      f:write(string.format("dma frame=%d pc=$%06X reg=%s mask=$%02X ch=%d regs=%s\n",
        frame, pc, regName, mask, ch, channelRegs(ch)))
    end
  end
  f:close()
end

-- The write callback fires with the value being written, so the channel
-- registers are read while they still describe this transfer.
emu.addMemoryCallback(function(address, value)
  pcall(logKick, "$420B", value & 0xFF)
end, emu.callbackType.write, 0x420B, 0x420B, emu.cpuType.snes, emu.memType.snesMemory)

emu.addMemoryCallback(function(address, value)
  pcall(logKick, "$420C", value & 0xFF)
end, emu.callbackType.write, 0x420C, 0x420C, emu.cpuType.snes, emu.memType.snesMemory)

emu.log("dma-trace: watching $420B/$420C, appending to " .. LOG)
