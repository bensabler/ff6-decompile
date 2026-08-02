---
name: static-runtime-correlator
description: Correlate Ghidra static analysis with Mesen runtime evidence for FF6 65816 code without promoting speculative disassembly.
---

# Static / Runtime Correlator

Use Ghidra to organize candidate code and Mesen to establish execution and behavior.

Read and follow:

- `../_shared/EVIDENCE_STANDARD.md`
- `../_shared/ADDRESS_SPACES.md`
- `../_shared/STOPPING_RULES.md`
- `../../../docs/static-analysis/GHIDRA_MESEN_WORKFLOW.md`
- `../../../docs/research/ROM_IDENTITY.md`

## Core rule

Ghidra suggests. Runtime evidence confirms.

## 65816 state requirements

At every proposed code entry, record the best-supported values or unknown status for:

- emulation/native mode;
- M flag and accumulator width;
- X flag and index width;
- program bank register;
- data bank register;
- direct-page register;
- stack pointer;
- relevant caller-established state.

A plausible listing is not sufficient evidence that decoding is synchronized.

## Required sequence

1. Resolve the address in canonical notation.
2. Verify ROM identity.
3. Compare bytes with Mesen or a trusted ROM dump.
4. Record Ghidra's current decoding and assumptions.
5. Obtain runtime execution evidence.
6. Bound the candidate routine conservatively.
7. Identify reads, writes, calls, and returns.
8. Design a falsification test.
9. Record confidence and alternatives.
10. Promote a semantic name only when evidence supports it.

## Naming policy

Use neutral names until confirmed, for example:

- `candidate_C0B593`
- `routine_C0B593`
- `dispatcher_candidate_C09B5C`

Do not assign names such as `EventInterpreter_MainLoop` solely from address proximity, decompiler output, or one trace.
