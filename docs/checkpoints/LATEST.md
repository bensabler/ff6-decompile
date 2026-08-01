# Latest Checkpoint

**[2026-08-01 — EXP-0043: the ATB layer is located](2026-08-01-exp0043-atb-layer-located.md)**

State: **the ATB layer is located.** Every timer domain the master
program asked for now has an address:

| Address | Role |
|---|---|
| `WRAM:+$3AB4` | **ATB gauge array** — 10 entries, **stride 2**, 16-bit; party 0-3, enemies 4-9 |
| `WRAM:+$3AC8` | per-slot **increment**; gauge += `$3AC8,X >> 1` at `ROMCPU:$C21195` |
| `WRAM:+$3AA0` | per-slot scheduler **flags** |
| `WRAM:+$3A3E` | 16-bit **battle tick counter**, one per non-gated frame |
| `WRAM:+$3218` | second accumulator, gated on `$3219,X` |
| `ROMCPU:$C21124` | **the gate** — `LDA $2F41 / AND $3A8F / BNE` |

Two headline results. **`$C21124` is where ACTIVE and WAIT diverge**:
that one branch skips the entire per-frame battle update, with `$3A8F`
the Wait flag from EXP-0042 and `$2F41` the untested other half.
And **Battle Speed scales enemy gauges only** — the `CPX #$08 / BCC`
branch skips the `$3A90` multiply for party slots, confirmed
numerically: party increments byte-identical at Bat.Speed 3 and 6
(318/330/336) while enemy increments went 240 → 156.

Three instruments were required to converge before anything was claimed —
read watch, static ROM decode, and live sampling showing `+$3AB4` advance
by exactly the predicted `$4E` = `$9C >> 1`, including a wrap. No
falsifier fired.

**Correction to carry forward:** the ATB family is **stride 2**, not the
`$14` of the HP/stat family. DISC-0001's unified layout governs slot
assignment, not stride.

**The blocker is now narrow rather than total.** Still unknown: the pause
condition (`$2F41` never observed non-zero), the exact increment and
threshold arithmetic, and the action queue. Whelk stays deferred until
the pause matrix exists.

8 evidence artifacts with hashes; no background processes; SRAM virgin.
All gates clean (gofmt/build/vet/test, `ff6lab audit`, census sync 64
entries, restricted-file scan).

Exact next action: **EXP-0044 — the ACTIVE/WAIT pause matrix.**
Write-watch `WRAM:+$2F41` across opening and closing each battle submenu,
then build the matrix of menu states against timer domains, reading
`+$3AB4`, `+$3A3E`, `+$3AA0` and `+$3218` directly. `$3A8F`/`$3A90` can
be patched in place, making ACTIVE versus WAIT a one-variable comparison
inside a single savestate lineage. First unit warranting
`/battle-baseline` and bounded parallel observers.
