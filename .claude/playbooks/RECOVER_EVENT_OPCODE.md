# Event opcode recovery

## Required inputs

- exact opcode or queued question;
- ROM identity;
- CORR-0001 and the current dispatcher evidence;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Resume from CORR-0001 rather than restarting it: `ROMCPU:$C09B5C` advances
   the 24-bit pointer at `WRAM:+$00E5..+$00E7` by `A & $FF`, Confirmed 24/24.
2. Exec-watch the dispatch predecessor `JMP ($002A)` at `ROMCPU:$C09B59`,
   capturing DP `$2A`/`$2B` and the opcode that reached it.
3. Correlate each dispatch with the entry it took from the candidate table at
   `ROMCPU:$C098C4`, which is static-only until observed.
4. Measure each opcode's operand length at dispatch. A length read statically
   is a candidate; a length measured is evidence.
5. Decode the handler: operands, side effects, flags touched, control flow.
6. Bound the work to the opcode subset the active route needs. Decoding every
   opcode in the game before integrating any is out of scope.
7. Implement the interpreter subset in Go with per-opcode tests.
8. Compare an executed transcript against the recorded event-flag timeline.

## Required outputs

- correlation and discovery records;
- opcode table with number, handler, length, operands, side effects;
- which opcodes are Confirmed, which are candidates, and which are unimplemented;
- Go interpreter subset and tests;
- confidence and alternatives;
- exact next action;
- checkpoint when the unit stops.
