# FF6 Reverse-Engineering Session 001

- **Date:** not recorded (before 2026-07-29; this note was reconstructed on
  2026-07-29 from
  [FF6_Decompilation_Session_01_Summary.md](../../FF6_Decompilation_Session_01_Summary.md))
- **Investigator:** Benjamin Sabler
- **ROM identity/checksum:** Final Fantasy III (USA) (SNES); checksum not recorded
- **Mesen CE version:** not recorded
- **Goal:** Discover one variable experimentally (current HP of the first
  displayed party slot) and identify the code that writes it.

## Starting state

Not recorded in detail. A running game with a visible first-party-slot HP
value; HP could be reduced on demand to drive the memory search.

## Experiment

1. Work RAM memory search: unsigned, 16-bit, exact value = displayed HP.
2. Take damage; re-search with the new displayed HP; repeat until one
   candidate remained.
3. Set a write breakpoint on the surviving candidate.
4. When the breakpoint fired, record the instruction, registers, and the
   surrounding disassembly.

## Raw observations

- Surviving search candidate: WRAM-relative `$2EB5`.
- Write breakpoint fired at ROM CPU `$C10E14`, instruction `STA $2EB5,Y`,
  with `A = $002A`, `Y = $0000`.
- Surrounding code pattern (loads paired with stores):
  `LDA $2E78,X / STA $2EB5,Y`, `LDA $2E80,X / STA $2EB7,Y`,
  `LDA $2E88,X / STA $2EB9,Y`, `LDA $2E90,X / STA $2EBB,Y`.
- Loop advance: `INX / INX` (X += 2) and `ADC #$0020 / TAY` (Y += $20).
- Loop termination: `CPX #$0008 / BNE $C10E11`.
- Routine entry `$C10DF3`; returns at `$C10E66` with `RTS`.

## Findings

### Finding 1 — `$2EB5` is first-slot current HP

- **Status:** Confirmed
- **Evidence:** Converging exact-value search; write breakpoint fires
  exactly when displayed HP changes.
- **Interpretation:** Party-slot field, not character-specific storage.
- **Alternatives:** Could be a staging/display copy rather than the
  authoritative value.
- Promoted to [03_DISCOVERED_VARIABLES.md](03_DISCOVERED_VARIABLES.md).

### Finding 2 — the routine at `$C10DF3` is a structured field copy

- **Status:** Strong hypothesis (copy structure); purpose Unknown
- **Evidence:** Four load/store pairs per iteration; X += 2, Y += $20;
  4 iterations (`CPX #$0008`).
- **Interpretation:** Struct-of-arrays source (`$2E78–$2E97`) copied into
  four `$20`-byte array-of-structs records based at `$2EB5`. Stores an
  already-computed value; does not calculate damage.
- **Alternatives:** General-purpose copy shared by several systems; records
  1–3 may not be party slots.
- Promoted to [02_DISCOVERED_FUNCTIONS.md](02_DISCOVERED_FUNCTIONS.md)
  and [05_DATA_STRUCTURES.md](05_DATA_STRUCTURES.md).

## Documentation changes

- 2026-07-29: findings promoted into `docs/00–08` (this repository's
  authoritative layout was created that day).

## Go changes

- 2026-07-29: module `github.com/bensabler/ff6-decompile` created; package
  `chardata` implements `CopyCharacterFields`, `CharacterFieldsSource`, and
  `PartySlotRecord` with table-driven tests.

## Verification

```text
gofmt:          clean (see repository root README for commands)
go test ./...:  pass (run 2026-07-29)
go vet ./...:   pass (run 2026-07-29)
```

## Open questions

- Promoted to [08_OPEN_QUESTIONS.md](08_OPEN_QUESTIONS.md) (all nine).

## Next experiment

Execution breakpoint at `$C10DF3`. Capture the call stack / return address
and the on-screen game context each time it fires (battle start, menu,
damage, per-frame). Outcome: identifies the caller and settles whether the
routine belongs to battle initialization, damage processing, or something
else — the single answer that most constrains every current hypothesis.
