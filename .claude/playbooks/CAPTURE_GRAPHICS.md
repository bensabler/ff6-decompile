# Runtime graphics capture

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Recreate deterministic frame.
2. Record PPU mode/registers.
3. Dump VRAM/CGRAM/OAM.
4. Capture tile/tilemap/sprite viewers.
5. Capture isolated layers.
6. Hash all files.
7. Create runtime-capture manifest.
8. Mark source provenance unresolved.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
