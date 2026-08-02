# DMA provenance

## Implementation status

This playbook described a procedure with **no implementation beneath it** until
2026-08-02. Nothing read a DMA register, in Lua or in Go, while the command,
the skill, the agent and this file all existed. What exists now:

| Piece | Status |
|---|---|
| `internal/platform/snesdma` — register decode, transfer patterns, source spans, trace parsing | tested, fuzzed |
| `internal/mesenstate.DMAChannels` — channel state from a preserved savestate | tested |
| `ff6lab state ppu` — prints channel setup alongside PPU config | works offline |
| `mesen/probes/dma-trace.lua` — logs registers when `$420B`/`$420C` is written | **UNEXERCISED**, never run |

**Savestate channel state is not a trace.** It is configuration at one instant:
a finished channel keeps its last setup, an unused one keeps initialisation
values. Measured falsifier — the mines savestate's channel 0 claims
`$7EC180 → $2118`, but VRAM matches that buffer on only 18% of bytes with a
49-byte leading run, so that setup is not the transfer that produced the VRAM.

Until `dma-trace.lua` has been run and checked against a known transfer, a
source address is a **lead**, never provenance.

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Freeze the exact scene/frame.
2. Record DMA channel state. Offline: `ff6lab state ppu <file.mss>`. Live:
   load `probe dma-trace` **before** the transfer you care about, since it
   captures on the enable write.
3. Capture source, destination and length. Remember a size of **zero means
   65536** — `snesdma.ByteCount` and `mesenstate.DMAChannel.ByteCount` apply
   that rule; the raw register does not.
4. Trace the trigger: the PC that wrote `$420B`/`$420C` is in the trace line.
5. Follow buffer ancestry. A WRAM source means the real question is what filled
   that buffer.
6. Find the decompressor or table — **after** confirming the data is not a
   verbatim ROM copy (`ff6lab state origin`).
7. Hash pre/post regions.
8. Link asset IDs and register newly visible systems in the census.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives, with lead distinguished from provenance;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
