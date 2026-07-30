# Project Goals

## Objective

Reconstruct the behavior and architecture of **Final Fantasy III (USA) for
SNES** (Final Fantasy VI) in clean, testable, idiomatic **Go**.

This is not a literal decompiler. The pipeline is:

```text
Original ROM
    ↓
Reverse Engineer (Mesen CE debugger)
    ↓
Understand Intent
    ↓
Recreate in Go
```

The goal is understanding, not transliteration. 65C816 instructions are not
mechanically converted to Go; the underlying data structures, state
transitions, and algorithms are recovered first.

## Fidelity target

Default target is **behavioral reconstruction**: reproduce externally
observable game behavior. Hardware-specific details (overflow, truncation,
RNG, timing) are preserved when they are known to affect results, and all
intentional deviations are marked in code comments and documentation.

## Ground rules

- No external FF6 RAM maps, symbol lists, or existing decompilation projects
  are consulted during discovery. Every conclusion must be supported by
  evidence observed in this project's own debugger sessions.
- Facts and hypotheses are always kept separate. See
  [01_REVERSE_ENGINEERING_RULES.md](01_REVERSE_ENGINEERING_RULES.md).
- Only behavior with sufficient evidence is implemented in Go, and every
  implemented behavior gets tests.

## Development environment

- macOS (Intel x86_64)
- Mesen Community Edition (SNES debugger, memory search, breakpoints)
- Target ROM: Final Fantasy III (USA) (SNES)
- Go (module `github.com/bensabler/ff6-decompile`)

No ROM data or copyrighted assets are stored in this repository.
