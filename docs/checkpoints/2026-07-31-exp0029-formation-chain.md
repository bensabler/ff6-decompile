# Checkpoint 2026-07-31 — EXP-0029: formation loader chain decoded (Unit 29)

## Current question
None open. Next: EXP-0030 (+$3F46 producer write-watch at encounter
entry -> the ROM formation table).

## State
Static unit (no emulator). Evidence at
local_artifacts/experiments/EXP-0029/ (decode notes, hashes).
Census 48 entries; ROM ledger 23 regions (ROM-0023 rebounded after
the overlap audit fired — the tooling caught a real mistake).

## Work completed
EXP-0029: per-slot loader $C22C30 decoded (id x4/x8/x32 scaling into
bank-$CF tables; $CF8400 = per-monster attribute words); sole caller
$C22F22 loops six enemy slots reading ids from WRAM:+$3F46 ($FF
sentinel; $0206 alternate under the $3A97 flag). Formation record =
one producer-hop away.

## Tests and quality gates
Run at commit: gofmt/build/vet/test (10 packages)/audit clean.

## Git status
main, 3 ahead of origin after this commit. Not pushed.

## Exact next action
EXP-0030: arm write-watches on +$3F46-+$3F5F, trigger a mines
encounter (EXP-0028 walk pattern), capture the writer PC and decode
its ROM source operand — the formation table. Then extract the
opening formation record as the anchor.
