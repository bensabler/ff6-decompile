Use the 65816-analyst, function-recovery, static-runtime-correlator, and experiment-designer skills. Follow `.claude/playbooks/RECOVER_EVENT_OPCODE.md`.

Target: $ARGUMENTS

Resume from CORR-0001: `ROMCPU:$C09B5C` advances the 24-bit pointer at `WRAM:+$00E5..+$00E7` by `A & $FF`, Confirmed 24/24, but the value has never been observed dereferenced and the dispatch predecessor at `ROMCPU:$C09B59` is unresolved. The candidate jump table at `ROMCPU:$C098C4` is static-only and has never been exercised.

Record each opcode's number, handler address, operand length, operands, side effects, and the observation that established it. A length inferred from a static read is a candidate; a length measured at dispatch is evidence.
