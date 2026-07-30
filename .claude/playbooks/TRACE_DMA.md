# DMA provenance

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Freeze exact scene/frame.
2. Record DMA channels and registers.
3. Capture source/destination/length.
4. Trace trigger.
5. Follow buffer ancestry.
6. Find decompressor/table.
7. Hash pre/post regions.
8. Link asset IDs.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
