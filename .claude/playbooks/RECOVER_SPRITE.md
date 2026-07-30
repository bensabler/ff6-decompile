# Sprite reconstruction

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Select one frame.
2. Identify OAM objects.
3. Resolve tile addresses.
4. Resolve palette banks.
5. Recover origin/offsets/flips/sizes/priority.
6. Trace loading.
7. Compose in Go.
8. Compare pixels.
9. Link animation only when proven.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
