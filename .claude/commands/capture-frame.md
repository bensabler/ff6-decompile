Use the graphics-archaeologist, mesen-operator, and asset-cataloger skills. Follow `.claude/playbooks/CAPTURE_FRAME.md`.

Target: $ARGUMENTS

Capture one deterministic frame with the state needed to interpret it: frame identity, PPU configuration, VRAM/CGRAM/OAM, and hashes. A screenshot alone is not a capture — a PNG cannot say which VRAM span a layer read from.

Prefer a savestate: `internal/mesenstate` reads VRAM, CGRAM, OAM, ARAM and the PPU/DMA registers out of one offline, and `ff6lab state ppu` decodes it. A frame captured as a savestate stays analysable after the session ends; a frame captured as a PNG does not.
