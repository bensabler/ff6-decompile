# Current Focus

**Program:** SCN-0001 — reconstruct the opening scenario, New Game
through Whelk victory, as an evidence-backed vertical slice.
Master record: `docs/scenarios/SCN-0001-opening-to-whelk.md`;
machine manifest: `data/scenarios/opening-to-whelk.json`
(19 beats B01–B19, per-beat status and links live there).

## Golden route — power-on through milestone 05

One tracked probe (`mesen/probes/EXP-0036.lua`) drives power-on →
mines interior, reproducibly:

| Milestone | Frame | Established by |
|---|---|---|
| `00-new-game` | 5 200 | EXP-0031 |
| `01-opening-cinematic` | 15 000 | EXP-0032 |
| `02-narshe-entry` | 30 000 | EXP-0032 |
| `03-first-scripted-battle` | 31 677 | EXP-0032 |
| `04-free-movement` | 46 375 | EXP-0034 |
| `05-mines-entry` | 51 578 | EXP-0036 |

Milestone WRAM is byte-identical across repeated power-on runs at
every milestone (two runs each for 00–04; **three** for 05). The
route's second half is a **state-driven 17-leg controller** —
position targets on `WRAM:+$00AF`/`+$00B0`, battle edges, elapsed
settles, per-leg timeouts that name the earliest divergent leg.
Tracked model and tests: `internal/scenario/route`, including a
probe-sync test so the Lua and Go encodings cannot drift apart.

**Assertion channel is WRAM, not screenshots** — CEN-QUIRK-0002 has
now been seen at two milestones (01 and 05): identical WRAM,
non-byte-stable frame capture.

## Scenario facts established so far

- The opening runs **five scripted battles** before the mines:
  formations **2, 1, 2, 41** (before free movement) then **84** (en
  route to the shaft). Every staged `+$3F44` record matches the ROM
  formation table `$CF6200 + id×15` byte-for-byte — six independent
  verifications of EXP-0030's table.
- Pre-Whelk monsters identified so far: records **0** (40 HP/15 MP),
  **25** (27 HP/5 MP), **27** (115 HP/30 MP) from the scripted
  battles, plus **19** and **77** from the mines random encounter.
  Monster *names* remain unlocated — the name table is not found, so
  on-screen labels are observational only.
- Battle 1 rewards **32 EXP / 96 GP**, with post-battle HP/MP written
  back into the field character block (`$C2496E`/`$C24979` →
  `+$1609`/`+$160D`).
- **Event-flag system located, decoded and inventoried**
  (CONTRA-0002 → EXP-0037/DISC-0008): three bit arrays at
  `WRAM:+$1E80`/`+$1EA0`/`+$1EC0` (decoder `ROMCPU:$BAED`, masks
  `$C0BAFC`/`$C0BB04`). The whole opening touches exactly **20
  flags** — 11 latched story flags, 4 transient, 5 engine working
  bits — deterministically across GUI and headless runs. The
  16-handler script-command family covers eight bases (five beyond
  the verified three, static-only), and every handler tail anchors
  the event interpreter at candidate `ROMCPU:$C09B5C`. Inventory:
  `data/scenarios/opening-event-flags.json`.
- Position bytes are **not field-meaningful during battle**.

## Controlled lab variables (whole program)

`RamPowerOnState=AllZeros`, virgin SRAM — originals backed up under
`local_artifacts/backups/`. Unattended sessions use headless Mesen
(`--testrunner`, `FF6_OUT`, frame-scheduled input); see BLOCKERS.md.

## Whelk reached (EXP-0039 breadth pass)

The scenario's furthest point to date. The mines-to-Whelk stretch is
short and linear: one corridor with a single turn, one dead-end nub, a
**bidirectional** exterior transition, a random-encounter zone drawing
formation 14, and two scripted beats bracketing the fight. Whelk is
**formation 432**, contact-triggered from `(2A,07)` after a scripted
beat at `(2A,09)`; its introduction/warning dialogue and **shell
counterattack** were observed, and the first attempt **ended in
defeat** (which captured the defeat flow, CEN-BATTLE-0007). A
pre-Whelk savestate is preserved.

Nothing in this stretch is gated behind an unsolved system — the only
obstacle to a victory is party HP management (entering at 26/19/56).

## EXP-0040 — Whelk victory attempted, not achieved (2026-08-01)

Two piloted attempts from the preserved pre-Whelk state, both under
`Bat.Mode = Wait`. **No victory**; milestone `10-whelk-victory` and B19
remain open. Confirmed along the way:

- **Whelk is two battle entities** — shell slot 4 = **50000/50000 HP**,
  head slot 5 = **1600/1600 HP** (DISC-0001 arrays, formation 432).
- **Head-only targeting works**: six measured hits (162-186 damage
  each) all reduced the head; the shell's 50000 never moved. No shell
  strike occurred, so EXP-0039's counter was not re-triggered.
- **Head/shell state is visually classifiable** at 4× upscale.
- **A field healing route exists** before contact — Tonic ×4 took the
  party from 26/19/56 to **76/77, 105/105, 106/107** at no turn cost.
- **MagiTek sets are character-specific and EXP-0039's list was
  incomplete**: Terra has eight abilities, Wedge and Vicks four.
- The **guard/Esper beat precedes Whelk** on a clean run
  (CEN-EVENT-0011 resolved).

## ~~Blocked — no ATB model~~ → resolved by EXP-0041..0044

Whelk execution was deferred from 2026-08-01 because the project could
not model ACTIVE/WAIT, qualifying submenu pause states, timer domains or
action-queue ordering. **EXP-0041 through EXP-0044 supplied all of those
except the action queue**, and the pause condition is now known
precisely.

EXP-0040's head/shell transitions remain **unusable as natural timing**,
but they can now be *scoped* rather than merely dismissed: only intervals
inside the ability list and target selection were paused. See BLOCKERS.md
for the discharge note and the remaining non-blocking work.

## ATB program started — EXP-0041 (2026-08-01)

The blocker's prerequisite audit is done and the program's **first
bounded unit is complete**. Battle configuration is no longer read off
the screen with our eyes:

| Byte | Contents | Default |
|---|---|---|
| `WRAM:+$1D4D` | bits 0-2 Bat.Speed (0-5 → 1-6), bit 3 **Bat.Mode** (1 = Wait), bits 4-6 Msg.Speed, bit 7 Cmd.Set | `$2A` |
| `WRAM:+$1D4E` | bit 4 Reequip, bit 5 Sound, bit 6 Cursor, bit 7 Gauge | `$00` |
| `WRAM:+$1D54` | bit 7 Controller | `$00` |

The block is **not contiguous**. Cleared bit = the left-hand screen
option. The Config screen marks the **selected** option with tile
attribute `$20` and the unselected with `$28` — the inverse of the
intuitive reading, and exactly why EXP-0040 misread `Bat.Mode`. That
correction is now confirmed from memory: **`Wait` is where a new game
arrives**, not an operator change. Configuration is **not** SRAM-backed
before a save.

Two supporting changes landed first: the battle-configuration
fingerprint is now a required, audited record field
(`internal/audit.CheckBattleExperimentConfig`, from EXP-0041 onward), and
`ff6lab state` reads work RAM and save RAM straight out of preserved
`.mss` files, so earlier captures can be mined without an emulator.
Trial 0 used it to extract the `+$1D4E` candidate from EXP-0040's
savestate pair before Mesen was ever launched.

## EXP-0042 — configuration is sampled at battle entry (2026-08-01)

Answered **mixed, and the split falls exactly where it matters.** At
battle entry `ROMCPU:$C22472` reads the two config bytes **once each**
and decomposes them into battle-local cells:

| Setting | Battle-local cell | Value |
|---|---|---|
| Bat.Mode (bit 3) | `WRAM:+$3A8F` | `01` = Wait, `00` = Active |
| Bat.Speed (bits 0-2) | `WRAM:+$3A90` | `255 − 24 × speed` (Fast `$FF` … Slow `$87`) |
| Cmd.Set | `WRAM:+$2F2E` | cleared when Window |
| Gauge | `WRAM:+$2021` | cleared when Off |
| `+$1D4E` bits 0-2 | `WRAM:+$2F34` | at `$C10FF7` |

`Bat.Mode` and `Bat.Speed` are **never re-read for timing** during the
battle. `Msg.Speed` and `Cursor` *are* read live, by `$C198AC` (message
delay table at `ROMCPU:$C19872`) and `$C159D6` (clears the `$5C`-byte
cursor-memory block at `+$890F` when Cursor = Reset — the mechanism
behind EXP-0040's `Cursor = Memory` observation).

Both `$3A8F`/`$3A90` values were **predicted from the disassembly before
the second run** and matched exactly.

**Staging rule for the whole ATB program:** ACTIVE/WAIT and Battle Speed
must be established **before battle entry** — or injected directly at
`+$3A8F`/`+$3A90`, which is a far better controlled handle than driving
menus.

## EXP-0043 — the ATB layer is located (2026-08-01)

| Address | Role |
|---|---|
| `WRAM:+$3AB4` | **ATB gauge array** — 10 entries, **stride 2**, 16-bit; party 0-3, enemies 4-9 |
| `WRAM:+$3AC8` | per-slot **increment**; gauge += `$3AC8,X >> 1` at `ROMCPU:$C21195` |
| `WRAM:+$3AA0` | per-slot scheduler **flags** |
| `WRAM:+$3A3E` | 16-bit **battle tick counter**, one per non-gated frame |
| `ROMCPU:$C21124` | **the gate**: `LDA $2F41 / AND $3A8F / BNE skip` |

`$3A8F` is the Wait flag, so `$C21124` is **where ACTIVE and WAIT
diverge**. `$2F41` — zeroed at battle entry, `00` throughout a
free-running battle — is the untested other half.

**Battle Speed scales enemy gauges only.** The `CPX #$08 / BCC` branch at
`$C209F6` skips the `$3A90` multiply for party slots, and measurement
agrees: party increments byte-identical at Bat.Speed 3 and 6
(318/330/336) while enemy increments went 240 → 156.

Watch the stride: the ATB family is **stride 2**, not the `$14` of the
HP/stat family. DISC-0001's unified layout governs slot *assignment*, not
stride.

## EXP-0044 — the ACTIVE/WAIT pause matrix (2026-08-01)

`WRAM:+$2F41` is the **battle submenu flag**: resting `$00`, cleared
per-frame at `ROMCPU:$C17A92`, raised at `ROMCPU:$C17C01` when a
qualifying submenu opens. `ROMCPU:$C21124` ANDs it with the Wait flag and
skips the entire per-frame battle update.

| State | ACTIVE | WAIT |
|---|---|---|
| Battle running, no menu | advances | advances |
| **Main command window open** | advances | **advances** |
| Ability list open | advances* | **paused** |
| Target selection | advances* | **paused** |
| Action resolving / animation | advances* | **advances** |
| Item / Magic / Row / Defend, dialogue, damage display, victory, defeat | not sampled | not sampled |

\* structurally implied — the gate is an `AND`, so ACTIVE can never pause;
verified directly for the ability-list row.

All four located domains — tick `+$3A3E`, gauges `+$3AB4`, flags
`+$3AA0`, accumulator `+$3218` — froze and resumed **together**. No
independent clock was found among them.

**The pause is narrower than the folk model.** The command window is on
screen for most of a WAIT battle and does not pause; neither do action
animations. Verified by an in-place one-variable flip of `+$3A8F` inside
a single savestate, cross-checked against genuinely configured WAIT.

**This scopes EXP-0040.** Its Whelk intervals inside the ability list and
target selection were paused; intervals at the command window and during
action resolution were not. Whether that is enough to reinterpret those
captures or whether the fight should be re-run is an orchestration call,
not a blocked one.

## EXP-0045 — queued work resolves past the gate (2026-08-01)

Per-frame tracing settled EXP-0044's unresolved transient. The scheduler
stops dead — tick and gauges frozen across 1 245 gated frames with zero
exceptions — but **an action already pending when the gate engages still
completes**, clearing `+$3AA0` bit 6 and advancing that slot's `+$3218`
by exactly `$0100`.

| Trace | Pending at gate engage | Gated-frame changes |
|---|---|---|
| 1 | slot 6 | 1, at +78 frames |
| 2 | **none** | **0** across 438 frames |
| 3 | slot 8 — **predicted** | 1, at +119 frames |

Trace 3 was a prediction test: the discriminator from traces 1 and 2 says
*engage the gate while a slot is pending and a deferred completion will
fire*. It did, on a different slot, at a different delay. Arming had to
move into Lua — bridge round-trips are hundreds of frames apart.

**`+$3AA0` bit 6 is the pending-action marker**, resolving a semantics
question EXP-0043 and EXP-0044 both left open. And it **vindicates
EXP-0040**: "queued actions resolved out of issue order" while menus were
open was real system behaviour, not operator error.

Matrix additions: **Item** is a qualifying submenu; the **victory
presentation** is not. Magic/Row/Defend are unreachable from a Magitek
battle and stay unsampled.

## Next exact action

**EXP-0046 — name the driver of deferred completion.** Something
advances a pending action while `ROMCPU:$C21124` is shut. Read-watch
`+$3AA0` and `+$3218` through a gated interval arranged as trace 3
arranged it, and capture the writing PC. That routine is the
**action-queue execution path** — the next major piece of the ATB model,
and the one EXP-0040's out-of-order observation depends on.
