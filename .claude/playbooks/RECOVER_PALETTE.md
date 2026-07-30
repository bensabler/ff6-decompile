# Palette recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Capture CGRAM.
2. Identify used indexes.
3. Trace writes.
4. Find ROM/base table.
5. Separate fades/cycles/transforms.
6. Decode BGR555.
7. Validate colors and timing.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
