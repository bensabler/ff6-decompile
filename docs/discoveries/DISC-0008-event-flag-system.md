# DISC-0008: Event-flag bit arrays — addressing, handler family, and opening inventory

## Status
Confirmed (addressing, masks, handler structure, opening write
timeline); implemented and tested. Flag *meanings* are deliberately
unassigned.

## Supporting experiments
CONTRA-0002 (location + decoder), EXP-0036 (first live writer),
EXP-0037 (complete opening write inventory; static decode of every
live writer PC). Records under `docs/contradictions/` and
`docs/experiments/`.

## Discovery
FF6 keeps event flags in parallel 32-byte bit arrays addressed by a
shared decoder at `ROMCPU:$C0BAED` (byte = flag/8, bit = flag&7) with
single-bit masks from `ROMCPU:$C0BAFC` (set: `01 02 04 08 10 20 40
80`) and `$C0BB04` (clear: exact complements). Sixteen uniform 20-byte
script-command handlers at `ROMCPU:$C0B593-$C0B6D2` set/clear one flag
(number in DP `$EB`) in eight bases — `$1E80/$1EA0/$1EC0` (runtime
verified) plus `$1EE0/$1F00/$1F20/$1F40/$1DC9` (static only) — and
every handler ends `LDA #$02 / JMP $C09B5C`, anchoring the event
interpreter (CEN-EVENT-0001). Boot clears `$1E80-$1EDF` (and
`$1E40-$1E6F`) via STZ loops at `$C0BB0C`/`$C0BB18`.

The opening (power-on → milestone `05-mines-entry`) touches exactly
**20 flags** in the three verified arrays: 11 latched story flags (all
set through the script handlers, seven in the boot/title era, four
during the battle/route beats), 4 transient, and 5 fast-toggling
working bits — the engine stores per-frame working state (facing
nibble, input bit, busy latch in `$1EB6`; `$1EB9`/`$1EBC-$1EBF`
neighbors) *inside the flag address space*, distinguished from story
flags purely by writer PC. The timeline is deterministic: identical
across one GUI and two headless runs at frame+address+value+PC
granularity, with the final WRAM byte-identical to milestone 05.
Tracked inventory: `data/scenarios/opening-event-flags.json`.

## Go implementation
`eventflags.ArrayBase, eventflags.Ref, eventflags.FlagAt,
eventflags.SetMasks, eventflags.ClearMasks`
(`internal/game/eventflags`)

## Tests
`TestRefDerivation, TestMaskTablesMatchROMDecode, TestFlagAtRoundTrip,
TestFlagAtRejectsOutside, TestMilestone05Value, TestID, TestIsSetBounds`
(eventflags)

## Confidence and residue
Addressing, masks, handler family and the opening timeline are
Confirmed (instruction-level static decode of every observed writer,
cross-checked by three deterministic runs and replay-consistent
snapshots). Residue, tracked in CEN-EVENT-0008's unknowns: flag
meanings; SRAM backing; runtime use of the five statically-decoded
extra bases; the `$1E40-$1E6F` cleared region; the `$0200/$0205` bit
mirrored into `EVF-1EC0-$FF`; and the flag-test callers.
