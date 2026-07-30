# Graphics validation

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Lock emulator state.
2. Render reconstructed target.
3. Normalize dimensions and transparency.
4. Compare indexed pixels when possible.
5. Classify differences.
6. Store diff report.
7. Update confidence.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
