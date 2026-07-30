# Audio cue investigation

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Create deterministic trigger.
2. Trace CPU/APU ports.
3. Capture ARAM and DSP state.
4. Identify driver dispatch.
5. Identify sequence and instruments.
6. Identify BRR samples.
7. Record timing.
8. Create cue manifest.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
