# CPU to SPC command trace

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Define trigger.
2. Break on port writes.
3. Record handshake.
4. Trace SPC port reads.
5. Find dispatcher.
6. Follow sequence/sample activation.
7. Repeat to confirm.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
