# Structure recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Collect indexed accesses.
2. Infer candidate stride.
3. Find copy/init loops.
4. Build offset matrix.
5. Compare multiple instances.
6. Test boundaries.
7. Preserve unknown bytes.
8. Publish provisional schema.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
