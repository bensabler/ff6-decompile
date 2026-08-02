# Ghidra + Mesen Workflow

## Purpose

Ghidra provides a navigable static map of the FF6 ROM. Mesen provides authoritative runtime execution and state. Neither tool alone is treated as proof of semantic meaning.

## Local layout

Keep the external workstation outside the tracked repository:

```text
FF6-Reverse-Engineering/
├── ff6-decompile/              # Git repository
├── private/roms/               # local ROM only
├── tools/                      # Ghidra and extension source
└── workspaces/ghidra/          # Ghidra .gpr and .rep files
```

The Ghidra project database is not a source artifact and should not be committed.

## Verified local baseline

The expected local setup is:

- Ghidra 12.1.2;
- SNES Loader built against the same Ghidra version;
- language `65816:LE:24:snes`;
- import format `SNES ROM Loader`, not `Raw Binary`;
- FF6 recognized as HiROM;
- canonical bank blocks such as `bank_c0_hirom`;
- known addresses such as `ROMCPU:$C09B5C` and `ROMCPU:$C0B593` resolve directly.

Record actual local versions; do not silently assume these values forever.

## The safe evidence chain

```text
existing experiment or specific question
→ canonical ROMCPU address
→ verify bytes in Ghidra and Mesen
→ record 65816 entry state
→ bounded disassembly
→ runtime breakpoint/trace
→ falsification experiment
→ correlation record
→ canonical function/discovery record
→ Go implementation only when justified
```

## Critical 65816 limitation

The meaning and length of some instructions depend on processor state. Ghidra can display plausible but incorrect instructions when M/X, bank, direct-page, or entry-boundary assumptions are wrong.

Therefore:

- do not trust the decompiler first;
- do not assume pressing `D` established a valid function boundary;
- do not recursively disassemble large unknown regions;
- recover caller-established state before confirming a routine;
- compare the listing against Mesen's executed instruction trace.

## First pilot: `ROMCPU:$C0B593`

The current Ghidra project can resolve and disassemble this address. Treat the created function as `candidate_C0B593` until the repository's existing evidence and a new runtime trace establish:

- that execution enters at exactly `$C0B593`;
- the M/X state at entry;
- the data-bank and direct-page assumptions;
- all valid return/exit paths;
- whether `$C0B593` is a true function boundary or an internal basic block;
- what behavior it actually owns.

The C-like output is a convenience view, not evidence.

## Ghidra actions and evidence status

| Action | Meaning |
|---|---|
| Navigate to an address | No evidentiary claim |
| Press `D` | Candidate disassembly under current context |
| Press `F` | Candidate function boundary |
| Auto-generated XREF | Static lead |
| Decompiler pseudocode | Derived interpretation |
| User label with no record | Unreviewed annotation |
| Breakpoint fires in controlled run | Runtime execution evidence |
| Paired experiment changes behavior | Behavioral evidence |
| Static and runtime results agree | Correlated evidence |

## Session procedure

At the start:

```text
/resume-session
/bootstrap-ghidra
```

For a precise target:

```text
/correlate-static-runtime ROMCPU:$C0B593
```

Before stopping:

```text
/census-observations
/run-quality-gates
/checkpoint
```

The checkpoint must state whether the Ghidra database contains unexported labels or comments and identify the exact address to resume.
