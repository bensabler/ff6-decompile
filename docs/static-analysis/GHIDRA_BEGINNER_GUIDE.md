# Ghidra Beginner Guide for This Project

## What is safe

Opening, navigating, labeling, bookmarking, disassembling, and creating functions modifies the external Ghidra database, not the ROM and not the Go repository.

## What can go wrong

The usual failure is not lost source code. It is recording a confident but incorrect interpretation.

Common causes:

- importing with Raw Binary instead of the SNES loader;
- using an extension built for a different Ghidra version;
- beginning disassembly at a data byte or in the middle of an instruction;
- incorrect 8-bit/16-bit accumulator or index assumptions;
- wrong data-bank or direct-page assumptions;
- treating mirrored SNES addresses as distinct code;
- trusting the decompiler before validating the listing.

## Current first function

At `ROMCPU:$C0B593`, Ghidra now shows instructions and a C-like function after manual disassembly/function creation. This proves the toolchain can navigate and decode the mapped ROM. It does not yet prove the function boundary or purpose.

Do not rename it semantically yet. A neutral local name such as `candidate_C0B593` is appropriate.

## Simple rule

Every time Ghidra appears to answer a question, ask Mesen to verify it.
