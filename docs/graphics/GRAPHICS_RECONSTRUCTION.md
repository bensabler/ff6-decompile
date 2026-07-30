# Graphics Reconstruction

## Runtime domains

- VRAM stores tile and tilemap data used by the PPU.
- CGRAM stores palette state.
- OAM stores sprite object state.
- PPU registers determine modes, base addresses, priorities, scroll, windowing, and color math.

A faithful reconstruction requires both bytes and interpretation.

## Asset classes

- planar tile bank;
- palette;
- tilemap;
- sprite frame;
- animation;
- font;
- menu/UI layout;
- map layer;
- battle background;
- effect;
- complete runtime scene.

## Validation levels

1. byte-level format validation;
2. indexed-pixel validation;
3. per-layer validation;
4. final-frame validation;
5. animation/timing validation.
