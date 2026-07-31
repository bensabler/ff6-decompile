-- Template probe. Copy to EXP-NNNN.lua, adjust, and load with the bridge
-- command:  probe EXP-NNNN
-- Instrumentation code stays in the repository (tracked), satisfying the
-- preserve-injected-code rule without one-line eval strings.
dofile((_G.FF6_PROBE_DIR or "mesen/probes/") .. "common.lua")

-- Example: watchwrites("EXAMPLE-WRITER", 0x11B0, 0x11B1, 30)
