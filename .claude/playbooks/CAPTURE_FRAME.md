# Deterministic frame capture

## Required inputs

- exact scene and the input schedule that reaches it;
- ROM identity;
- emulator identity;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Reach the scene on a reproducible schedule, from power-on where the route
   supports it.
2. Record frame identity: frame count, and the WRAM assertion channel. Frame
   capture is **not** byte-stable at every milestone — CEN-QUIRK-0002 has been
   seen at two — so WRAM is what proves you are at the same point.
3. Save a savestate. It carries VRAM, CGRAM, OAM, ARAM and the PPU/DMA
   registers, and stays analysable offline after the session ends; a PNG does
   not.
4. Record PPU configuration (`ff6lab state ppu`): BG mode, layer chr/tilemap
   bases, scroll, main/sub screen.
5. Take the screenshot as a human-readable companion, never as the runtime
   asset.
6. Hash every artifact into the experiment directory's `hashes.sha256`.
7. Copy artifacts out of `mesen/out/` into `local_artifacts/experiments/EXP-NNNN/`.
   `mesen/out/` is scratch, and Mesen's auto-save slot has destroyed evidence
   there before.
8. Register what the capture newly makes visible via the content census.

## Required outputs

- savestate, PPU state, and hashes under the experiment directory;
- frame identity and the schedule that reproduces it;
- census registration for newly visible systems;
- exact next action;
- checkpoint when the unit stops.
