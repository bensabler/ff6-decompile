# EXP-0037: Which event flags does the opening set, and when?

- **Status:** completed (2026-08-01) — stopping condition met: two
  headless runs byte-identical on every evidence channel
- **Program:** SCN-0001 (opening-to-Whelk); serves B16 (interaction
  flags), B19 (post-battle state), persistence (CEN-SAVE-0001), and
  anchors the event engine (CEN-EVENT-0001) via CEN-EVENT-0008.

## Question

Across the established scheduled route (power-on → milestone
`05-mines-entry`), which bytes and bits of the three event-flag arrays
at `WRAM:+$1E80`, `+$1EA0`, `+$1EC0` (CONTRA-0002) are written, by what
instruction, at what frame, and in which scenario beat? This is an
**inventory of writes**, not a decoding of flag meanings.

## Starting state

Power-on under the SCN-0001 controlled lab: Mesen 2.1.1,
`RamPowerOnState=AllZeros`, virgin SRAM (Saves dir empty), ROM SHA-256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.
The EXP-0036 17-leg scheduled route drives power-on → mines interior
with no interactive input; EXP-0037 adds instrumentation only and
**does not alter the schedule** — route encoding, phase boundaries and
leg conditions are byte-identical to `mesen/probes/EXP-0036.lua`
(guarded by the existing `probe_sync_test.go` pattern: EXP-0037's ROUTE
table must match `MinesRoute()` too).

## ROM identity

`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`,
3 145 728 bytes, no copier header (verified this session via
`ff6lab archive verify`).

## Emulator identity

Mesen 2.1.1 (Session 002 identification). Exploratory pass in GUI mode
if the desktop allows it; evidence passes under `--testrunner`
(the mode that established milestones 00–05).

## Scheduled route and milestones (context frame anchors)

| Milestone | Frame | Battle entries (EXP-0034/0036) |
|---|---|---|
| `00-new-game` | 5 200 | — |
| `01-opening-cinematic` | 15 000 | — |
| `02-narshe-entry` | 30 000 | — |
| `03-first-scripted-battle` | 31 677 | b1 = 31 557 |
| `04-free-movement` | 46 375 | b2 = 34 953, b3 = 36 828, b4 = 39 500 |
| `05-mines-entry` | 51 578 | b5 = 46 802 |

## Independent variable

None — this is an observational inventory along a fixed schedule.
Run-to-run repetition (≥2 headless runs) is the control: the flag-write
timeline must be identical across runs or the divergence is reported.

## Controlled variables

AllZeros RAM, virgin SRAM before every run, same ROM, same probe file,
frame-scheduled input only, no savestates loaded mid-run.

## Instrumentation

`mesen/probes/EXP-0037.lua` — EXP-0036's route controller (unchanged)
plus:

1. **Write callback on `WRAM:+$1E80`–`+$1EDF`** (all 96 bytes of the
   three arrays). For every store: frame, writer PC, address, old
   value, incoming value, derived set/cleared bit numbers, current
   route context (phase/leg, battle state, battle count, player tile),
   emitted as one JSONL record and one human-readable log line.
   Same-value rewrites are logged (first 3 per writer-PC/address fully,
   then counted; totals reported at run end).
2. **Array snapshots** (96-byte hex) at arm time, at every milestone
   frame above, at every leg begin/end, at every battle entry/end, and
   at run end.
3. **Integrity check (analysis step):** initial snapshot + ordered
   write log must reproduce every later snapshot byte-for-byte. A
   mismatch means an unmonitored write path (e.g. DMA into WRAM) and is
   itself a reportable finding.
4. Battle detection, leg logging, alias sampling: as in EXP-0036.

Flag identifiers use the address-anchored, semantics-free scheme
`EVF-<arraybase>-$NN` (array-relative flag number `$00`–`$FF`, from
byte offset×8 + bit index, per the CONTRA-0002 decoder: `Y = flag/8`,
`X = flag&7`). Example: the CONTRA-0002 bits at `+$1EA5` are
`EVF-1EA0-$28`, `EVF-1EA0-$2A`, `EVF-1EA0-$2B`.

## Expected outcomes

- *Supports (inventory usable):* every value-changing write is
  attributable (frame, PC, bit); the reconstructed state matches every
  snapshot; ≥2 headless runs produce identical timelines
  (frame + address + value sequence); the known `+$1EA5` transitions
  (`$00→$01→$05→$0D` at ~34 298 / ~39 090 / shaft dialogue) reappear
  as bit-sets of `EVF-1EA0-$28/$2A/$2B` via `ROMCPU:$C0B5B3` — a
  built-in cross-check against EXP-0036's independent log.
- *Refutes / complicates:* runs disagree on the timeline; snapshots
  diverge from the write-reconstructed state; or writes bypass the
  monitored span.

## Falsifiers

1. **Callback completeness:** any snapshot ≠ initial + logged writes ⇒
   the write watch misses a path; the inventory cannot claim
   completeness until the path is found (registered, not chased).
2. **Determinism:** any cross-run timeline difference ⇒ the route's
   flag behavior is not schedule-determined; report the earliest
   divergent write.
3. **Cross-check:** if `+$1EA5` does not show the EXP-0036 sequence,
   the instrumentation itself is suspect (it disagrees with two prior
   independent observations).

## Evidence requirements

Per run: `<run>-flags.jsonl` (full write timeline),
`<run>-snapshots.log` (all 96-byte array snapshots with frame + label),
`<run>-events.log` (route/leg/battle log), final WRAM dump; plus
`hashes.sha256` over all of it. All under
`local_artifacts/experiments/EXP-0037/`. Nothing tracked carries raw
dumps; the tracked record carries the flag inventory (numbers, frames,
PCs), which is project-derived reconstruction data.

## Trials

1. **gui1** — visible GUI pass at normal speed (operator-watchable),
   probe v1. Route executed uncorrected: all five battles and all
   17 legs at the canonical EXP-0036 frames; final capture at 51 578.
   One probe defect found and fixed for later runs: leg-*begin*
   snapshots were labeled with the scheduled start frame
   (sample + 8 neutral frames), so the `$1EB6` working bits could
   toggle inside the label gap — surfaced by the replay-integrity
   check at `leg-8-begin`, not by a real missed write. gui1's JSONL
   also carries ~1 458 extra post-capture lines because GUI emulation
   continued until the process was killed; all are same-value rewrites
   and none is a value change.
2. **run1** — headless `--testrunner`, probe v2 (snapshot label fix,
   `FF6_STOP=1`). Clean exit via `emu.stop()`.
3. **run2** — headless, identical configuration, fresh virgin-SRAM
   boot (run1's clean exit wrote a `.srm` that hashes identical to the
   backed-up original `6afbcf1e…`, i.e. the game's initialized-empty
   image; deleted before run2 per the lab controls).

## Observations

### Static decode of every live-observed writer PC

The probe logs the PC *after* the storing instruction (CONTRA-0002
convention). Every after-PC observed live resolves to a decoded
instruction sequence in the ROM:

| After-PC | Store | Sequence (decoded from ROM bytes) | Role |
|---|---|---|---|
| `$C0B5B6` | `$C0B5B3 STA $1EA0,Y` | `LDA $EB / JSR $BAED / LDA $1EA0,Y / ORA $C0BAFC,X / STA` | script SET handler, array `$1EA0` |
| `$C0B5CA` | `$C0B5C7 STA $1EC0,Y` | same pattern, base `$1EC0` | script SET handler, array `$1EC0` |
| `$C0B5F2` | `$C0B5EF STA $1EA0,Y` | `AND $C0BB04,X` variant | script CLEAR handler, array `$1EA0` |
| `$C0BB11` | `$C0BB0E STZ $1E80,X` | `LDX $00 / STZ $1E80,X / INX / CPX #$0060 / BNE` | boot clear-all: 96 bytes `$1E80-$1EDF` |
| `$C0BA94` | `$C0BA91 STA $1EB6` | `LDY $0803 / LDA $087F,Y / TAX / LDA $1EB6 / AND #$F0 / ORA $C0BAFC,X / STA` | working byte `$1EB6`: low-nibble one-hot from `$087F,Y` |
| `$C0BAA0` | `$C0BA9D STA $1EB6` | `LDA $06 / BPL … / LDA $1EB6 / ORA #$10 / STA` | working bit 4 set (from `$06` sign) |
| `$C0BAAA` | `$C0BAA7 STA $1EB6` | `LDA $1EB6 / AND #$EF / STA` | working bit 4 clear |
| `$C0BFAE` | `$C0BFAB STA $1EB6` | `LDA $1EB6 / ORA #$40 / STA` | bit 6 busy-latch set |
| `$C0BFC6` | `$C0BFC3 STA $1EB6` | `LDA $1EB6 / AND #$BF / STA` | bit 6 busy-latch clear |
| `$C0BFB5`/`$C0BFB8` | `STZ $1EBE` / `STZ $1EBF` | inside the same routine, when `$58`=0 | working bytes zeroed |
| `$C0C034` | `$C0C031 STA $1EB9` | `LDA $1EB9 / ORA #$80 / STA` | direct set of flag `$CF` (cleared later by the script handler) |
| `$C0B951` | `$C0B94E STA $1EDF` | `LDA $0205 / AND #$80 / STA $1A / LDA $1EDF / AND #$7F / ORA $1A / STA` | copies bit 7 of `$0205` into flag `$FF` |
| `$C0B98A` | `$C0B987 STA $1EDF` | same idiom, source `$0200` | copies bit 7 of `$0200` into flag `$FF` |
| `$C0B9B6` | `$C0B9B3 STA $1EDF` | same idiom, source `$0200` | third copy handler (same-value only in this route) |
| `$C04A25`/`$C04A2D` | `STA $1EB6` / `STA $1EB7` | `AND #$DF` / `AND #$7F` RMW | battle-era clears of working bits `$B5`/`$BF` (same-value here) |
| `$C04AFB`-`$C04B13` | `STA $1ED8` ×4 | `AND #$EF/#$DF/#$BF/#$7F` chain gated on `$1A6D` | clears of `$1ED8` bits 4-7 (same-value here) |

### The handler family is larger than the three watched arrays

The six script set/clear handlers sit inside a uniform 16-handler
family at `ROMCPU:$C0B593-$C0B6D2` (20 bytes each, every one
`LDA $EB / JSR $BAED / LDA base,Y / ORA-or-AND mask,X / STA base,Y /
LDA #$02 / JMP $C09B5C`), covering **eight bases**: `$1E80`, `$1EA0`,
`$1EC0`, `$1EE0`, `$1F00`, `$1F20`, `$1F40`, and unaligned `$1DC9`.
Six flag **test** routines (`AND $C0BAFC,X / RTS`) at
`$C0BAA9-$C0BAE9` cover `$1E80/$1EA0/$1EE0/$1F00/$1F20/$1F40`.
A second boot clear loop at `$C0BB18` zeroes `$1E40-$1E6F`. Every
handler tail (`LDA #$01|#$02 / JMP $C09B5C`) is consistent with an
event-command interpreter advancing its script pointer by the
operand length — `ROMCPU:$C09B5C` is therefore a concrete
**candidate event-interpreter advance routine** (registered against
CEN-EVENT-0001; not chased here). The extra arrays are registered as
static-only census material; this experiment's runtime watch remains
the three verified arrays.

## Interpretation

The three arrays are a **shared address space with two populations**,
separable purely by writer PC:

- **Story flags** flow exclusively through the script-command
  handlers (`$C0B5B3`/`$C0B5C7` set, `$C0B5EF` clear — after-PCs
  `$C0B5B6`/`$C0B5CA`/`$C0B5F2`), change value at most once or twice
  per run, and only in field/script contexts — **no value-changing
  write occurs while a battle owns the screen** (0 of 162).
- **Working bits** (`$1EB6` low nibble = one-hot value from
  `$087F,Y`; `$1EB6` bit 4 = `$06`-sign mirror; `$1EB6` bit 6 = busy
  latch; `$1EB9` bit 7 = engine pulse) are driven by direct RMW engine
  code at per-frame rates — 47 903 same-value rewrites vs 162 changes,
  99.7 % of all stores. `$1EB9` bit 7 (`EVF-1EA0-$CF`) is set by the
  engine but cleared through the script clear-handler, and its second
  pulse brackets the exterior-to-mines transition (frames
  50 878–50 880, leg 16) — an engine/script crossover worth noting
  and not interpreting further.

The opening's story-flag timeline divides into four eras: a
**boot/title burst** (frames 2 516–2 528: five latched flags + two
transients, between the first injected title input at 2 500 and the
title's visible appearance ~2 969), a **quiet march** (one clear at
18 981 over 25 000 frames), a **battle-era chain** (`$C1` 417 frames
before battle 1; `$28` between battles 1–2; `$2A` before battle 4;
`$30` after battle 4), and a **route chain** (`$31` in the
post-battle-5 settle; `$2B` in the shaft-dialogue settle, 214 frames
before the transition completes).

## Alternatives

- *The boot-burst flags could be input-triggered rather than timed*:
  they land 16–28 frames after title input injection begins. The
  constant schedule cannot discriminate; an input-free control run
  would. Recorded as an open question, not assumed either way.
- *The callback could miss non-CPU write paths (DMA)*: addressed by
  the shadow cross-check on every store and the snapshot replay — both
  clean in all three runs, so within this route no unmonitored path
  wrote the span. Outside this route the possibility remains open.
- *`$1EB6`-family "flags" could be a different subsystem that merely
  shares the address space*: nothing here distinguishes ownership;
  the population split is recorded as writer-PC fact, not as a
  semantic boundary.

## Result

**The complete opening (power-on → milestone 05) touches exactly 20
of the 768 watched flags — 162 value-changing writes, identical
across three runs (one GUI, two headless) at frame+address+value+PC
granularity.** The two evidence runs are **byte-identical on every
channel**: the full 48 065-line write log (including all 47 903
same-value rewrites), all 62 snapshots, and the final 128 KB WRAM —
which equals the established milestone-05 hash `c26453d3…`, extending
that milestone to **five byte-identical runs** and proving the
instrumentation perturbs nothing.

Latched story flags (11), in first-set order:

| Flag | Byte.bit | Frame | Context (recorded, not interpreted) |
|---|---|---|---|
| `EVF-1EC0-$E0` | `$1EDC`.0 | 2 516 | boot/title era, via `$C0B5CA` |
| `EVF-1EC0-$F0` | `$1EDE`.0 | 2 516 | boot/title era |
| `EVF-1EA0-$0B` | `$1EA1`.3 | 2 520 | boot/title era, via `$C0B5B6` |
| `EVF-1EA0-$E3` | `$1EBC`.3 | 2 520 | boot/title era |
| `EVF-1EC0-$FE` | `$1EDF`.6 | 2 528 | boot/title era |
| `EVF-1EA0-$C1` | `$1EB8`.1 | 31 140 | scripted approach, 417 frames before battle 1 |
| `EVF-1EA0-$28` | `$1EA5`.0 | 34 298 | between battles 1 and 2 |
| `EVF-1EA0-$2A` | `$1EA5`.2 | 39 090 | before battle 4 |
| `EVF-1EA0-$30` | `$1EA6`.0 | 42 450 | scripted walk after battle 4, pos (26,11) |
| `EVF-1EA0-$31` | `$1EA6`.1 | 49 855 | leg 7: post-battle-5 settle at (1E,25) |
| `EVF-1EA0-$2B` | `$1EA5`.3 | 50 699 | leg 15: shaft-dialogue settle at (1F,16) |

Transient flags (4): `EVF-1EA0-$CC` (set 2 528, cleared 18 981
mid-march — the only mid-march story event), `EVF-1EA0-$CF` (engine
pulses at 3 359 and 50 878–50 880), `EVF-1EC0-$FF` (mirror of a
`$0200/$0205` bit, set 2 516 → cleared 2 527), `EVF-1EA0-$B3`
(working). Toggling working bits (5): `$B0/$B1/$B2` (facing nibble),
`$B4` (oscillator), `$B6` (busy latch) — all in `$1EB6`.

Boot behavior: `$C0BB0C` clears the whole 96-byte span at **frame
21**; New Game selection itself (milestone 00) adds **no** flag
writes; the known `+$1EA5` progression `$00→$01→$05→$0D` reproduces
exactly as bit-sets `$28`/`$2A`/`$2B` at 34 298 / 39 090 / 50 699,
all via `$C0B5B6` — the built-in cross-check against EXP-0036's
independent log **passes**.

**All three falsifiers came up clean**: snapshots replay-consistent
(0 shadow mismatches in ~48 k stores × 3 runs), timelines identical
across runs, cross-check reproduced. Tracked inventory with the full
per-flag event list: `data/scenarios/opening-event-flags.json`.

Side result: **GUI/testrunner input parity is verified for this
schedule** — the GUI run matched the headless runs to the frame in
battles, legs, flag timeline, and final WRAM (BLOCKERS updated).

## Confidence

- The 20-flag inventory, per-flag frames, writers and classifications:
  **Confirmed** (three independent runs, two byte-identical evidence
  channels each, instruction-level decode of every writer PC).
- Writer-population split (script handlers vs engine RMW):
  **Confirmed** as observed fact; the *ownership* reading of the
  `$1EB6` family is a Tentative hypothesis (see Alternatives).
- Flag **meanings**: **Unknown**, deliberately — identifiers stay
  address-anchored per the naming policy.
- Boot-burst input-dependence: **Unknown** (constant schedule cannot
  discriminate).

## Stopping condition

Stop when **two headless runs produce identical flag-write timelines**
and all snapshots pass the integrity check (successful stop), **or**
after three runs if timelines disagree (report divergence), **or** on
any route failure (leg timeout — the route, not the watch, is then the
problem and EXP-0036's encoding is re-examined before retrying).
Flag **meanings are out of scope**: no story semantics are assigned;
unresolved meanings become explicit future questions.

## Bounds (scope control)

No event-interpreter reconstruction, no event-script extraction, no
map-format decoding, no SRAM decoding, no story-state interpretation.
Newly visible systems are registered via the census workflow and left
uninvestigated.

## Next action

Golden route segment 5 (EXP-0038): extend the route controller into
the mines to milestone `06-random-encounter`, capturing the encounter
trigger context (`+$11E0` producer, CEN-WORLD-0006) on the way — the
most direct advance toward Whelk. The event-interpreter anchor
(`ROMCPU:$C09B5C`) and the extended flag bases are queued separately
(RESEARCH_QUEUE); flag meanings remain open questions tracked in
CEN-EVENT-0008.
