# Latest Checkpoint

**[2026-08-01 — EXP-0041: battle configuration storage](2026-08-01-exp0041-battle-config-storage.md)**

State: the **ATB research program is open and has produced its first
result.** All nine in-game Config settings are located and bit-decoded:
`WRAM:+$1D4D` (bits 0-2 Bat.Speed 0-5, **bit 3 Bat.Mode**, bits 4-6
Msg.Speed, bit 7 Cmd.Set, default `$2A`), `WRAM:+$1D4E` (bits 4-7
Reequip/Sound/Cursor/Gauge, default `$00`), and `WRAM:+$1D54` (bit 7
Controller). The block is **not contiguous** — a narrow read window hid
Controller until the full-WRAM diff caught it. Both speed fields were
swept to their clamps rather than inferred.

The Config screen marks the **selected** option with tile attribute `$20`
and the unselected with `$28` — the inverse of the intuitive reading, and
the exact cause of EXP-0040's `Bat.Mode` misread. That correction is now
**independently confirmed from memory**: `Wait` is where a new game
arrives. Configuration is **not** SRAM-backed before a save.

Two enabling changes landed first, each its own commit: the
battle-configuration fingerprint is now a **required, audited** record
field (`internal/audit.CheckBattleExperimentConfig`, EXP-0041 onward),
and **`ff6lab state`** reads work RAM and save RAM straight out of
preserved `.mss` files — trial 0 used it to extract the Cursor candidate
from EXP-0040's savestate pair *before Mesen was launched*, and the live
falsifier run reproduced it exactly.

**The hard ATB blocker remains open**: no timer domain, pause condition,
or action-queue semantics is known. Whelk was not resumed and its
savestates were not reloaded. 17 evidence artifacts preserved with
hashes; no background processes; SRAM still virgin. All gates clean
(gofmt/build/vet/test, `ff6lab audit`, census sync 62 entries,
restricted-file scan).

Exact next action: **EXP-0042 — battle-entry configuration sampling.**
Read-watch `WRAM:+$1D4D`-`+$1D4E` across a battle entry from the mines
random encounter to determine whether battle code consults these bytes
directly or a copy taken at entry.
