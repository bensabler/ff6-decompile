-- EXP-0024: battle confirm SFX trigger capture (audio vertical, Unit 5).
-- Watches CPU->APU port writes ($2140-$2143) and per-frame DSP voice
-- activity (ENVX != 0) around a frame-scheduled press, so a press
-- trial and a no-press trial can be delta-compared.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")
local OUT = _G.FF6_OUT_DIR or "mesen/out/"

-- APU port write watch (always on once armed; capped logging).
_G.exp24_pw = 0
_G.exp24_pwRef = emu.addMemoryCallback(function(addr, value)
  exp24_pw = exp24_pw + 1
  if exp24_pw <= 200 then
    local c = emu.getCpuState(emu.cpuType.snes)
    local pc = (c.k << 16) | c.pc
    local f = io.open(OUT .. "exp24-ports.log", "a")
    if f then
      f:write(string.format("PORTW addr=%06X val=%02X pc=%06X frame=%d\n",
        addr, value, pc, emu.getState().frameCount))
      f:close()
    end
  end
end, emu.callbackType.write, 0x002140, 0x002143, emu.cpuType.snes, emu.memType.snesMemory)

-- Per-frame DSP voice-activity bitmap inside the observation window.
-- ENVX register is voice reg $x8; nonzero means the voice is sounding.
local function voiceBitmap()
  local bits = 0
  for v = 0, 7 do
    if emu.read(v * 0x10 + 0x08, emu.memType.spcDspRegisters) ~= 0 then
      bits = bits | (1 << v)
    end
  end
  return bits
end

_G.exp24_hold = 0
emu.addEventCallback(function()
  if exp24_hold > 0 then
    emu.setInput({ a = true }, 0)
    exp24_hold = exp24_hold - 1
  end
end, emu.eventType.inputPolled)

_G.exp24_sched = nil
emu.addEventCallback(function()
  local s = exp24_sched
  if not s or not s.anchor then return end
  local ok = pcall(function()
    local fc = emu.getState().frameCount
    local rel = fc - s.anchor
    if s.press and not s.pressed and rel >= s.pressAt then
      s.pressed = true
      exp24_hold = 10
      probelog(string.format("EXP24-PRESS %s frame=%d", s.label, fc))
    end
    if rel >= s.watchFrom and rel <= s.watchTo then
      local bm = voiceBitmap()
      if bm ~= s.lastBm then
        local f = io.open(OUT .. "exp24-voices.log", "a")
        if f then
          f:write(string.format("VOICES %s rel=%d bitmap=%02X prev=%s portw=%d\n",
            s.label, rel, bm, tostring(s.lastBm), exp24_pw))
          f:close()
        end
        s.lastBm = bm
      end
    end
    if rel >= s.watchTo + 30 then
      probelog(string.format("EXP24-TRIAL-END %s portw=%d", s.label, exp24_pw))
      exp24_sched = nil
    end
  end)
  if not ok then exp24_sched = nil end
end, emu.eventType.endFrame)

-- exp24_trial(mss, label, press): press=true injects A for 10 frames
-- at anchor+72; observation window rel 40..200 either way.
_G.exp24_trial = function(mss, label, press)
  local f = io.open(mss, "rb")
  if not f then return "cannot open " .. mss end
  local data = f:read("*a")
  f:close()
  exp24_pw = 0
  exp24_sched = {
    label = label, press = press, pressAt = 72,
    watchFrom = 40, watchTo = 200, lastBm = nil,
  }
  local ref
  ref = emu.addMemoryCallback(function()
    emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
    local okL = pcall(emu.loadSavestate, data)
    if okL and exp24_sched then
      exp24_sched.anchor = emu.getState().frameCount
      probelog(string.format("EXP24-TRIAL %s anchor=%d press=%s", label,
        exp24_sched.anchor, tostring(press)))
    else
      probelog("EXP24-LOADFAIL")
      exp24_sched = nil
    end
  end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
  return "exp24 trial " .. label .. " queued"
end

return "EXP-0024 probe armed"
