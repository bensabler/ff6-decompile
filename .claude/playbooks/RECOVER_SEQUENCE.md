# Sequence recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Trace interpreter loop.
2. Build opcode table.
3. Determine lengths/timing.
4. Recover branches/calls/loops.
5. Map note/instrument/control commands.
6. Preserve unknown opcodes.
7. Replay events.
8. Validate against trace.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
