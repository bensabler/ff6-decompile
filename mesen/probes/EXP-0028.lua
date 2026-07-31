-- EXP-0028: battle-init enemy stat population (question: CEN-MONSTER-0001).
-- First-capture-per-PC write watch over the $14-stride battle stat
-- family block (+$3B18-+$3BB7); the battle-entry burst names the
-- populate routine(s) and their register context.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

watchwrites("EXP28-STATW", 0x3B18, 0x3BB7, 96)

-- Movement helper for the encounter walk: hold a direction for n
-- frames via the input-poll hook (frame-scheduled, headless-safe).
_G.exp28_hold = nil
emu.addEventCallback(function()
  local h = exp28_hold
  if h and h.left > 0 then
    emu.setInput(h.buttons, 0)
    h.left = h.left - 1
    if h.left == 0 then exp28_hold = nil end
  end
end, emu.eventType.inputPolled)

_G.exp28_walk = function(dir, frames)
  exp28_hold = { buttons = { [dir] = true }, left = frames }
  return "walk " .. dir .. " " .. frames
end

return "EXP-0028 probe armed"
