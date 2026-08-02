# SNES reconstruction (Lane 3) — not started

This directory is the home for a **matching** reconstruction: 65816 and SPC700
source that assembles to the original ROM. It is empty by design.

Lane defined by [ADR-0002](../../docs/decisions/ADR-0002-project-product-lanes.md).

## Why this exists before it has content

So that assembly-level concerns have somewhere to go **other than the Go port**.
The native Go demo (`cmd/ff6demo` and `internal/`) reproduces *behavior* by any
correct means. A matching reconstruction must reproduce the *encoding*. Those
are different correctness criteria, and mixing them degrades both: the Go code
acquires instruction-level constraints it does not need, and the matching work
inherits abstractions that make byte-identity harder.

## Correctness criterion

**Byte-matching output.** A reconstruction is correct when the assembled result
is byte-identical to the verified ROM revision, SHA-256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

That bar is strictly higher than the Go port's. Behavior-equivalent code that
assembles differently is a failure here and a success there.

## What this lane will eventually hold

- 65816 assembly for the main CPU;
- SPC700 assembly for the audio driver;
- event bytecode, as source rather than as extracted bytes;
- the linker/bank layout that places it all;
- byte-matching tests, per region, against the ROM;
- round-trip data encoders — the inverse of the Lane 1 extractors, which is
  what makes a data table *reconstructable* rather than merely readable.

## What it must never hold

- Commercial ROM bytes, in any form. The same
  [asset policy](../../docs/legal/ASSET_POLICY.md) that governs every other
  lane governs this one. A byte-matching test compares against the operator's
  local ROM; it does not vendor it.
- Work that delays DEMO-0001. This lane is explicitly lower priority than the
  playable demo until DEMO-0001 accepts.

## Prerequisites, honestly stated

ROM ownership is **0.49 %** (`indexes/ROM_REGIONS.md`). A matching
reconstruction of a 3 MiB ROM is not a near-term prospect, and nothing here
should be started in a way that implies otherwise. The realistic entry point is
a single fully-owned region — the battle damage engine in bank `$C2` is the
best-understood candidate, being byte-exact and regression-tested in Go
(DISC-0002…0006) — reconstructed as assembly and checked against the ROM bytes
it claims to reproduce.

Until then this file is the whole lane.
