# Text and dialogue recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- a capture rendering the target text on screen;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Identify which text system the target belongs to. The fixed-width menu
   encoding is Confirmed (EXP-0049, `VRAM tile = $100 + byte`); the dialogue
   stream is a different system with its own font and control codes.
2. Locate the storage region and the pointer table; record its base, stride and
   entry width.
3. Establish the character encoding against **rendered output**, not against a
   plausible-looking table.
4. Enumerate control codes with operand lengths; an unlisted byte is unmapped,
   not a guess.
5. Decode any dictionary or byte-pair expansion, and record its table.
6. Establish line-break, pagination, pause and advancement rules from observed
   rendering.
7. Implement the decoder, keeping unmapped bytes surfaced the way `textenc`
   does — never substituting a plausible character.
8. Extract the corpus locally; the bytes stay untracked, the addresses and
   hashes are tracked.

## Required outputs

- experiment and discovery records;
- pointer table and encoding with confidence per claim;
- control-code table with operand lengths;
- Go decoder, tests, fuzz target;
- the count of bytes deliberately left unmapped;
- exact next action;
- checkpoint when the unit stops.
