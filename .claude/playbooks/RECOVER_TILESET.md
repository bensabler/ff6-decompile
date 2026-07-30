# Tileset reconstruction

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Identify tile bank.
2. Determine bit depth.
3. Recover palette.
4. Recover tilemap format.
5. Recover flips/priority.
6. Trace loading/compression.
7. Decode in Go.
8. Render and validate.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
