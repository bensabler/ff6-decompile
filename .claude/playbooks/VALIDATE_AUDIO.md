# Audio validation

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Lock trigger and state.
2. Compare command events.
3. Compare voice/DSP events.
4. Compare sample identity.
5. Compare timing.
6. Optionally compare waveform.
7. Record assumptions and result.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
