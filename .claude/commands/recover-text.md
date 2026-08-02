Use the variable-recovery, struct-recovery, tileset-reconstructor, and documentation-curator skills. Follow `.claude/playbooks/RECOVER_TEXT.md`.

Target: $ARGUMENTS

Record the storage region, pointer table and its stride, the character encoding, every control code with its operand length, any dictionary or byte-pair expansion, line-break and pagination rules, and the font the stream renders through. Distinguish the fixed-width menu encoding (Confirmed, EXP-0049) from the dialogue stream, which is a different system.

**Never invent a decoded string.** An unmapped byte stays unmapped, as `textenc` already does.
