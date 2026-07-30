# Variable recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Create variable ID.
2. Collect all reads/writes.
3. Determine width/indexing.
4. Find initialization and resets.
5. Perturb one condition.
6. Track lifetime.
7. Identify owner candidates.
8. Record falsifier.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
