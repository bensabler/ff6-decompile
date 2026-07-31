-- EXP-0021: record-index and selection-window capture (question #30).
-- Tests EXP-0020's scheduling interpretation: do matched action ordinals
-- load the same attack record across differently-timed trials?
--
-- Trials are frame-scheduled (not wall-clock): exp21_trial() loads the
-- savestate and injects the A-presses at exact frame offsets from the
-- load anchor, so the experiment is emulation-speed independent and
-- runs identically under the GUI or --testrunner headless mode.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

-- Selection window: the verified LDA $3A70 at ROMCPU:$C23198, executed
-- immediately before each action setup (JSR $299F at $C2319D). Capture
-- the pre-setup context: registers (via probelog) plus DP $B5 and the
-- $3A70/$3A71 pair.
_G.exp21_sel = 0
_G.exp21_selRef = emu.addMemoryCallback(function()
  exp21_sel = exp21_sel + 1
  if exp21_sel <= 48 then
    local b5 = emu.read(0x7E00B5, emu.memType.snesMemory)
    local a70 = emu.read(0x7E3A70, emu.memType.snesMemory)
    local a71 = emu.read(0x7E3A71, emu.memType.snesMemory)
    probelog(string.format("EXP21-SEL n=%d b5=%02X a70=%02X a71=%02X",
      exp21_sel, b5, a70, a71))
  end
end, emu.callbackType.exec, 0xC23198, 0xC23198, emu.cpuType.snes, emu.memType.snesMemory)

-- Record index: the verified MVN $C4,$7E at ROMCPU:$C2297A. On the
-- first moved byte A=$000D (13 remaining of 14); X holds the source
-- offset, so index = (X - $6AC0) / 14 comes from the actual source
-- pointer rather than an assumed entry register.
_G.exp21_idx = 0
_G.exp21_mvnfires = 0
_G.exp21_idxRef = emu.addMemoryCallback(function()
  exp21_mvnfires = exp21_mvnfires + 1
  local c = emu.getCpuState(emu.cpuType.snes)
  if c.a ~= 0x000D then return end
  exp21_idx = exp21_idx + 1
  if exp21_idx <= 48 then
    local off = c.x - 0x6AC0
    probelog(string.format("EXP21-IDX n=%d srcX=%04X idx=%d rem=%d",
      exp21_idx, c.x, math.floor(off / 14), off % 14))
  end
end, emu.callbackType.exec, 0xC2297A, 0xC2297A, emu.cpuType.snes, emu.memType.snesMemory)

-- Fight-command populate: the verified STA $11A6 at ROMCPU:$C229D1
-- (power from +$3B68,X). Registers logged by probelog carry A (power)
-- and X (attacker offset).
_G.exp21_pop = 0
_G.exp21_popRef = emu.addMemoryCallback(function()
  exp21_pop = exp21_pop + 1
  if exp21_pop <= 48 then
    probelog(string.format("EXP21-POP n=%d", exp21_pop))
  end
end, emu.callbackType.exec, 0xC229D1, 0xC229D1, emu.cpuType.snes, emu.memType.snesMemory)

-- Zero the ordinal counters between trials so n= restarts per trial.
_G.exp21_reset = function()
  exp21_sel, exp21_idx, exp21_mvnfires, exp21_pop = 0, 0, 0, 0
  return "exp21 counters reset"
end

-- ---------------------------------------------------------------------------
-- Frame-scheduled trial driver.
-- exp21_trial(mss_path, label, first_delay): load the savestate, then
-- press A for 10 frames at anchor+first_delay and three more times at
-- 180-frame intervals (the EXP-0018/0020 cadence: 1.2s/4.5s ~= 72/270
-- frames, 3s ~= 180). Summary markers land at anchor+first_delay+720
-- and anchor+990 (the latter equalizes the count-comparison window
-- across trials).
-- ---------------------------------------------------------------------------
_G.exp21_hold = 0
emu.addEventCallback(function()
  if exp21_hold > 0 then
    emu.setInput({ a = true }, 0)
    exp21_hold = exp21_hold - 1
  end
end, emu.eventType.inputPolled)

_G.exp21_sched = nil
emu.addEventCallback(function()
  local s = exp21_sched
  if not s or not s.anchor then return end
  local ok = pcall(function()
    local fc = emu.getState().frameCount
    for i, pf in ipairs(s.presses) do
      if not s.fired[i] and fc >= s.anchor + pf then
        s.fired[i] = true
        exp21_hold = 10
        probelog(string.format("EXP21-PRESS %s k=%d frame=%d", s.label, i, fc))
      end
    end
    if not s.summed and fc >= s.anchor + s.presses[1] + 720 then
      s.summed = true
      probelog(string.format("EXP21-TRIAL-SUM %s sel=%d idx=%d pop=%d mvnfires=%d",
        s.label, exp21_sel, exp21_idx, exp21_pop, exp21_mvnfires))
    end
    if fc >= s.anchor + 990 then
      probelog(string.format("EXP21-TRIAL-END %s sel=%d idx=%d pop=%d mvnfires=%d",
        s.label, exp21_sel, exp21_idx, exp21_pop, exp21_mvnfires))
      exp21_sched = nil
    end
  end)
  if not ok then exp21_sched = nil end
end, emu.eventType.endFrame)

_G.exp21_trial = function(mss, label, first)
  local f = io.open(mss, "rb")
  if not f then return "cannot open " .. mss end
  local data = f:read("*a")
  f:close()
  exp21_reset()
  exp21_sched = {
    label = label,
    presses = { first, first + 180, first + 360, first + 540 },
    fired = {},
  }
  -- loadSavestate is only legal inside a main-CPU exec callback; anchor
  -- the schedule to the frame counter as it stands right after the load.
  local ref
  ref = emu.addMemoryCallback(function()
    emu.removeMemoryCallback(ref, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
    local okL = pcall(emu.loadSavestate, data)
    local fc = emu.getState().frameCount
    if okL and exp21_sched then
      exp21_sched.anchor = fc
      probelog(string.format("EXP21-TRIAL %s anchor=%d first=%d", label, fc, first))
    else
      probelog(string.format("EXP21-TRIAL %s LOADFAIL", label))
      exp21_sched = nil
    end
  end, emu.callbackType.exec, 0, 0xFFFFFF, emu.cpuType.snes, emu.memType.snesMemory)
  return "trial " .. label .. " queued first=" .. first
end

return "EXP-0021 probe armed"
