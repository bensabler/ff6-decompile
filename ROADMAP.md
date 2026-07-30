# Roadmap

## Phase 0 — Lab integrity
- deterministic environment
- ROM identity
- Mesen capability inventory
- repository migration
- evidence schemas
- restart-safe checkpointing

## Phase 1 — Vertical proof
Complete one end-to-end unit in each domain:
- behavior
- graphics
- audio

## Phase 2 — Battle foundation
- party and enemy records
- battle initialization
- command queue
- target selection
- damage/status pipeline
- battle UI synchronization

## Phase 3 — Graphics foundation
- tile/palette primitives
- OAM composition
- animation representation
- battle backgrounds
- menu/font rendering
- map layers

## Phase 4 — Audio foundation
- CPU/APU protocol
- SPC driver map
- BRR decoder
- sample directory
- sequence format
- DSP state model
- sound effect and music validation

## Phase 5 — Event and world systems
- event scripts
- map data
- collision and movement
- transitions
- save data
- inventory and menus

## Phase 6 — Reconstructed game slices
Build behaviorally verified vertical slices in Go.

## Phase 7 — Release engineering
- stable CLI
- reproducible reports
- tagged releases
- public repository audit
- contributor onboarding
