# Map system recovery

## Required inputs

- exact target map or queued question;
- ROM identity;
- preserved captures for **at least two different maps**;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Freeze captures for two or more maps. One map cannot distinguish a header
   field from a constant.
2. Record each map's PPU configuration and tile/tilemap sources
   (`ff6lab state ppu`, `ff6lab state origin`).
3. Separate per-map data from shared assets. EXP-0051 found three blocks
   common to every field scene — they are not what a header selects.
4. Locate the header record: its table base, stride, and which field selects
   the observed blocks.
5. Decode tilemap/layout storage, which is not the VRAM tilemap — that is only
   the visible window, composed at runtime.
6. Decode palette selection, collision, exits, and event-trigger placement.
7. Record layer order, priority and scroll behavior.
8. Implement decoders, extractor and a rendered scene; compare against the
   capture.

## Required outputs

- experiment and discovery records;
- header layout with per-field confidence, and unmapped bytes named as unmapped;
- the per-map values that distinguished each field;
- Go decoders, extractor, manifest entries, tests;
- confidence and alternatives;
- exact next action;
- checkpoint when the unit stops.
