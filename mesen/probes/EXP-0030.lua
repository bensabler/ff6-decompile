-- EXP-0030: formation staging producer (+$3F46 family writers).
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")
watchwrites("EXP30-FORMW", 0x3F40, 0x3F5F, 64)
_G.exp30_hold = nil
emu.addEventCallback(function()
  local h = exp30_hold
  if h and h.left > 0 then
    emu.setInput(h.buttons, 0)
    h.left = h.left - 1
    if h.left == 0 then exp30_hold = nil end
  end
end, emu.eventType.inputPolled)
_G.exp30_walk = function(dir, frames)
  exp30_hold = { buttons = { [dir] = true }, left = frames }
  return "walk " .. dir .. " " .. frames
end
return "EXP-0030 probe armed"
