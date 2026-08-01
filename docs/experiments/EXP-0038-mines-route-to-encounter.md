# EXP-0038: Mines traversal to the first random encounter (milestone 06)

- **Status:** completed (2026-08-01) — stopping condition met:
  two scheduled runs, same encounter frame, byte-identical milestone WRAM
- **Program:** SCN-0001; serves B11 (mine traversal), B12 (encounter
  check/trigger), B13 (encounter packs) — runtime rows only.

## Question

From milestone `05-mines-entry` at (`$26`,`$1C`), where does the mines
corridor lead, and does continued scheduled walking produce a
**reproducible first random encounter** — SCN-0001 milestone
`06-random-encounter` — with its trigger context captured
(`WRAM:+$11E0` formation id, staged `+$3F44` record, position, frame)?

## Starting state

- Recon: GUI Mesen + bridge, loading the archived milestone state
  `local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss`
  (byte-provenance: promoted from EXP-0036 acceptance run 3).
- Scheduled runs: power-on under the SCN-0001 controlled lab
  (AllZeros RAM, virgin SRAM), EXP-0037's probe route (17 legs,
  unchanged) extended with the new mines legs.

## ROM identity

`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`
(unchanged).

## Emulator identity

Mesen 2.1.1; GUI for recon (operator-watchable), `--testrunner` for
evidence runs (parity verified for this schedule family by EXP-0037).

## Independent variable

None — recon then fixed-schedule reproduction, as EXP-0035→0036.

## Controlled variables

Lab controls as EXP-0036/0037. The scheduled route's leg encoding is
the only new element; phase 0 and legs 1-17 stay byte-identical.

## Instrumentation

- Recon: bridge `loadstate` + press batches + position reads
  (`+$00AF`/`+$00B0`) + screenshots at junctions; battle detection by
  position-byte behavior and screenshot.
- Scheduled: `mesen/probes/EXP-0038.lua` = EXP-0037's probe (route,
  battle detector, event-flag watch **kept running observationally**)
  + appended walk legs + a terminal leg ending on `battle_start`
  (battle 6 = the encounter). At encounter entry: `+$11E0`, staged
  `+$3F44`, position, frame; milestone capture at entry + 120 frames
  (WRAM dump, screenshot, savestate), matching the milestone-03
  pattern. Go model `MinesEncounterRoute()` extends `MinesRoute()`;
  the probe-sync test maps EXP-0038.lua to it.

## Expected outcomes

- *Supports:* recon maps a walkable corridor from (`$26`,`$1C`); the
  scheduled extension walks it and battle 6 triggers at the same frame
  and position across ≥2 power-on runs with byte-identical milestone
  WRAM; formation id recorded (EXP-0030's mines observation was
  formation 44 — same value would corroborate, a different one is
  equally valid data).
- *Refutes / complicates:* no encounter within the walk budget
  (negative result recorded — encounter rate too low for this
  corridor or state); runs disagree on encounter frame/position
  (encounter roll not schedule-determined — report the divergence
  precisely; that itself is a CEN-WORLD-0006 datum).

## Falsifiers

1. A leg times out (corridor mis-mapped) — fail, report leg.
2. Cross-run disagreement on encounter frame/formation/position —
   the "encounter is schedule-deterministic" claim is refuted;
   record what varies.
3. Milestone WRAM not byte-identical across runs.

## Evidence requirements

Recon transcript + per-leg positions + junction screenshots;
per-run events log, flags JSONL (continued), snapshots, encounter
context lines, milestone-06 WRAM/screenshot/savestate,
`hashes.sha256` — under `local_artifacts/experiments/EXP-0038/` and
`local_artifacts/scenarios/SCN-0001/06-random-encounter/`.

## Trials

1. **Recon (GUI + bridge, from `05-mines-entry.mss`).** Corridor
   mapped: straight north `(26,1C)→(26,0B)` (up blocked at `0B`),
   east `(26,0B)→(28,0B)` (a rail-trestle scene), north
   `(28,0B)→(28,09)`, east `(28,09)→(2A,09)`. A **random encounter**
   fired during the northward walk, trigger tile `(26,0B)`,
   `+$11E0 = $000E` = **formation 14** — three small enemies against
   the Magitek party (screenshot archived). Fought through with
   A-presses; control returned at `(26,0B)`. At `(2A,09)` a
   **scripted event** fired — dialogue box, the party splits, one
   walker advances (project description; game text not tracked) —
   and movement input is consumed. Recon stopped there: the event is
   beyond this unit's bound. Transcript, per-leg positions, and three
   screenshots archived under `local_artifacts/experiments/EXP-0038/`.
2. **run1 / run2** — scheduled headless power-on runs of the 26-leg
   encoding (results below).

## Observations

### Recon-derived route facts

- The scheduled route stops at `(28,09)` — one turn short of the
  `(2A,09)` event trigger — and patrols the proven corridor to extend
  the step budget rather than risk entering the event zone.
- A random encounter can interrupt a **walk** leg (all five scripted
  battles fell in pulse legs), so the controller restarts a blocked
  leg's timeout window while a battle owns the screen, mirroring
  `internal/scenario/route`'s Runner, which already modeled this.

### run1 — the encounter is reproduced on the scheduled route

| Item | Value |
|---|---|
| Encounter (battle 6) entry | **frame 51 307** |
| `+$11E0` | `$000E` = **formation 14** |
| Interrupted leg | **19** (east along the trestle, begun 51 258) |
| Position at entry / after battle | `26,00` (battle) → `(26,0B)` |
| Preceding leg 18 | begun 50 986 at `(26,1C)`, ended 51 250 at `(26,0B)` |
| Battle 6 end | frame 52 881 |

**Static cross-check passes (seventh verification of EXP-0030's
table).** The live staged `+$3F44` record
`A0 1C FF FF 13 13 13 FF 00 00 B8 4A 9E 00 23` is byte-identical to
ROM formation record 14 at `ROMFILE:0x0F62D2`. Its monster-id bytes
(record bytes 2-7) `FF FF 13 13 13 FF` decode to **three of monster
record 19** (`ROMFILE:0x0F0260`, `+$08` = 24 HP, `+$0A` = 0 MP) —
consistent with the three small enemies visible in the recon
screenshot.

**Two formations now observed in the mines interior**: 14 here and 44
in EXP-0030's earlier walk — a direct B13 datum (one area yields
multiple formations). Monster 19 is common to both.

**Both independent triggers landed in the same one-tile
neighbourhood**: recon (piloted, from the milestone-05 savestate)
triggered formation 14 arriving at `(26,0B)`; the scheduled run
(from power-on) triggered it one step later, leaving `(26,0B)`
eastward. Both followed a comparable number of steps from the
milestone-05 tile. Recorded as an observation consistent with a
step-counter check — **not** a mechanism claim; the `+$11E0` producer
remains the open CEN-WORLD-0006 question.

**No event flags are set by the mines corridor or the encounter.**
The flag timeline totals **162 value-changing writes with the last at
frame 50 880** — identical to EXP-0037's opening total, meaning legs
18-19 and the random encounter contributed **zero** new flag writes.
This is the first bounded evidence that random encounters do not
touch the three verified flag arrays.

### run1's probe defect (recorded, not hidden)

run1 logged the encounter and wrote its milestone WRAM dump, then
**stranded**: the milestone callback called `shot()`, which EXP-0037
had dropped from EXP-0036's helper set, so the call hit a nil value
and threw inside the callback. Only that callback died — the script
stayed alive and still logged `BATTLE 6 END` — but `finishRun()` and
therefore `emu.stop()` were never reached, so the process idled until
it was terminated during a background-task audit. The route itself
did not fail and the captured evidence is unaffected: the WRAM dump
completed *before* the throw, and the write log, snapshots and events
log are complete through the encounter.

Fixed for run2 by restoring `shot()` with a hardened body (screenshots
are explicitly **not** an assertion channel — CEN-QUIRK-0002 — so a
capture failure now logs `SHOT-SKIPPED` and continues) and by wrapping
the whole artifact block in `pcall` so no artifact error can ever
strand a run again.

### run2 — verification (probe fixed)

Byte-for-byte agreement with run1 on every channel:

| Channel | Result |
|---|---|
| Encounter entry | **frame 51 307** (identical) |
| Formation / leg / position | `$000E` / leg 19 / `(26,0B)` (identical) |
| Milestone-06 WRAM @ frame 51 427 | **byte-identical**, SHA-256 `c6e69ad7…` |
| Event-flag timeline | **162 value-changing writes identical** (frame+addr+old+new+PC) |
| Snapshot integrity | 66 snapshots replay-consistent, 0 shadow mismatches |

run2 exited cleanly through `emu.stop()`, and the new `pcall` guard
proved itself immediately: it caught a **second instance of the same
defect class** — `mkstate` was also dropped by EXP-0037 and was
likewise a nil global — logged
`MILESTONE-06-ARTIFACT-ERROR … attempt to call a nil value (global
'mkstate')`, and **still completed the run**. Where run1 stranded for
~20 minutes, run2 lost only one convenience artifact. `mkstate` has
now been restored in the probe alongside `shot`.

## Interpretation

The mines random encounter is **schedule-deterministic**: from a
fixed power-on schedule it fires at the same frame, in the same leg,
with the same formation, leaving byte-identical WRAM. That places the
encounter check downstream of state that the frame-exact schedule
fully determines — consistent with the step/zone-counter reading, but
this unit does not locate the producer and makes no mechanism claim
(CEN-WORLD-0006 stays open).

The corridor itself is short and linear with one turn, and the walk
plus encounter set **no event flags at all**, which is the first
direct evidence separating "story progress" writes from ordinary
traversal and combat in the three verified arrays.

## Alternatives

- *The encounter could be position-triggered rather than
  step-counted*: both observed triggers sit within one tile of
  `(26,0B)`. A run that reaches the same tile after a different number
  of steps would discriminate; not performed here.
- *Formation 14 might not be the only formation in this zone*:
  EXP-0030 saw formation 44 elsewhere in the mines. Whether the zone
  is the same and the roll simply differs is unresolved — a B13
  question, deferred to the breadth pass and the encounter-zone unit.

## Result

**Milestone `06-random-encounter` is established.** Two scheduled
power-on runs reproduce the mines random encounter at **frame
51 307** — formation **14** = three of monster record **19**,
interrupting leg 19 near tile `(26,0B)` — with **byte-identical
milestone WRAM** (`c6e69ad7…`) and identical event-flag timelines.
The route's 26-leg encoding executes uncorrected from power-on.

Supporting results: the staged formation record matches
`ROMFILE:0x0F62D2` byte-for-byte (**seventh** independent
verification of EXP-0030's table); the mines interior yields at least
two formations (14, 44); the corridor `(26,1C)→(26,0B)→(28,0B)→(28,09)`
is mapped, with a scripted event at ~`(2A,09)` registered as
CEN-EVENT-0009 and deliberately not entered; and **neither traversal
nor a random encounter writes any of the three event-flag arrays**.

**Known evidence gap (recorded, not papered over):** neither run
captured a milestone-06 **savestate**, because `mkstate` was a nil
global in both. The milestone rests on its byte-identical WRAM, which
is the project's stated assertion channel (CEN-QUIRK-0002 — frame
captures are explicitly not). A re-run purely to obtain the
convenience artifact was **not** performed; the probe is fixed, so the
next unit needing a scheduled milestone-06 state captures it for free.

## Confidence

- Encounter reproducibility (frame, formation, leg, milestone WRAM):
  **Confirmed** — two independent power-on runs, byte-identical on
  every channel, plus an independent piloted recon that hit the same
  formation in the same neighbourhood.
- Formation 14 identity and its ROM match: **Confirmed** (live
  staging vs static table, byte-identical).
- Monster record 19's HP/MP: **Confirmed** as record fields.
- Corridor geometry: **Confirmed** for the walked path only; branches
  and dead ends are **not** swept (EXP-0039).
- "Encounters set no event flags": **Confirmed for this encounter on
  this route**; not generalized.
- Encounter *mechanism* (step vs position vs zone): **Unknown**.

## Stopping condition

Stop when **two scheduled power-on runs produce the encounter at the
same frame with byte-identical milestone-06 WRAM** (milestone
established), **or** when a bounded recon walk budget (600 press-
batches ≈ ~10 min emulated) produces no encounter (negative recorded),
**or** after three failed schedule encodings, **or** on cross-run
encounter divergence (recorded as a CEN-WORLD-0006 datum — a valid
result).

## Bounds (scope control)

Primary corridor only — **no branch sweep** (B11's full sweep stays
open). No encounter-rate/zone-data decoding, no monster-record
decoding beyond ids, no map-format decoding, no treasure interaction.
Newly visible systems are registered via census observations only.
The event-flag watch rides along observationally; any flag activity
in the mines is inventory data, not a new investigation.

## Next action

EXP-0039 — **mines-to-Whelk breadth reconnaissance** in visible GUI
Mesen: explore the mines broadly from this corridor (branches, dead
ends, objects, encounters, graphics/audio families), register what is
visible, and advance the scenario toward Whelk. Depth targets this
unit deliberately left open — the `+$11E0` producer (CEN-WORLD-0006),
the `(2A,09)` scripted event (CEN-EVENT-0009), and encounter-zone
data (B13) — are queued, not chased.
