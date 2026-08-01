# EXP-0041: Battle configuration storage — where the Config screen's settings live

## Question

Where are the nine in-game Config settings stored in memory, and which
byte and bit encodes each — in particular **Bat.Mode (Active/Wait)** and
**Bat.Speed**?

This is the first unit of the ATB research program. It is deliberately
the cheapest one: the ACTIVE-versus-WAIT comparison cannot begin until
the controlled variable can be *set and verified*, and right now the
project can only read these settings off the Config screen with its
eyes. EXP-0040 did exactly that and got `Bat.Mode` wrong — it reported
`Active` and a switch to `Wait` when `Wait` was already selected, because
the hand cursor marks the row, not the selection. That misreading is one
of the two operator errors that produced the ATB blocker.

No ATB semantics are needed to answer this question, and no battle is
entered.

## Starting state

`local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss` —
milestone `05-mines-entry`, frame 51 578, established by EXP-0036 over
**three** byte-identical power-on runs, savestate-backed, and unrelated
to Whelk.

Deliberately **not** the pre-Whelk states: the standing directive from
the EXP-0040 checkpoint forbids reloading them, and this unit has no need
of them.

Fallback if the field menu proves unavailable at milestone 05:
`mesen/out/checkpoint2.mss`, where EXP-0026 confirmed the field menu
opens on X and the main menu exposes **Config**.

## ROM identity

`Final Fantasy III (USA).sfc`, sha256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2` —
verified at session start, matches the project revision.

## Emulator identity

Mesen 2.1.1 macOS, **headless** (`--testrunner --timeout=7200`, `FF6_OUT`
set), bridge v2. Exactly one instance. Lab controls unchanged:
`Snes.RamPowerOnState=AllZeros`, virgin SRAM (Saves directory verified
empty before launch).

Headless is required for unattended operation, and it is now also the
*better* choice: with a candidate config address in hand the assertions
can be WRAM reads rather than screenshots, which is the project's stated
assertion channel (CEN-QUIRK-0002).

## Battle configuration

Starting configuration, inherited from the golden route's new game and
never modified on this lineage:

| Setting | Value | Source |
|---|---|---|
| Bat.Mode | `Wait` | screen-read (EXP-0040) |
| Bat.Speed | 3 | screen-read (EXP-0040) |
| Msg.Speed | 3 | screen-read (EXP-0040) |
| Cmd.Set | `Window` | screen-read (EXP-0040) |
| Gauge | `On` | screen-read (EXP-0040) |
| Sound | `Stereo` | screen-read (EXP-0040) |
| Cursor | `Reset` | screen-read (EXP-0040) |
| Reequip | `Optimum` | screen-read (EXP-0040) |
| Controller | `Single` | screen-read (EXP-0040) |

Every value above was `screen-read` and therefore **not trustworthy** —
that is the defect this unit exists to repair.

**Resolved by this unit.** The same configuration, now `memory-read`:

`WRAM:+$1D4D` = `$2A`, `WRAM:+$1D4E` = `$00`, `WRAM:+$1D54` = `$00`
→ Bat.Speed 3, Bat.Mode **Wait**, Msg.Speed 3, Cmd.Set Window, Gauge On,
Sound Stereo, Cursor Reset, Reequip Optimum, Controller Single.

This is byte-identical to the golden route's post-new-game state and
matches EXP-0040's corrected screen reading in all nine positions,
including `Bat.Mode = Wait`.

## Independent variable

Exactly one Config setting per trial. All nine are exercised, each from
an identical reload of the starting savestate.

## Controlled variables

Same ROM, same emulator build and mode, same lab settings, same starting
savestate reloaded before every trial. Trials differ only in which
setting is changed and to what. No battle is entered. No other menu is
used.

## Instrumentation

- `mesen/probes/EXP-0041.lua`, built from `exp-template.lua`, reusing
  `watchwrites` from `probes/common.lua` over `WRAM:+$1D40`–`+$1D5F` to
  capture the writer PC for each config store, and `probelog` at each
  capture boundary.
- A full 131 072-byte WRAM dump per trial, as the wide net: if a setting
  is stored outside the watched window the diff still finds it.
- `ff6lab state diff` (added this session) for the comparisons.
- Screenshots archived as corroboration only, never as the assertion.

## Expected outcomes

- *Direct:* each of the nine settings maps to a byte and bit (or field)
  in a small contiguous block; `Bat.Mode` and `Bat.Speed` are located
  and their encodings recorded; the block's writer PC is captured.
- *Partial:* some settings map cleanly, others are packed or
  table-mapped in a way this unit can bound but not decode; those are
  registered with an explicit next action.
- *Refuted:* settings are not held in WRAM at all before a save, or the
  Config screen writes nothing until exit — either would be a real
  finding and would redirect the unit to the exit path.

## Falsifying outcome

1. **Primary.** A controlled `Cursor: Reset → Memory` toggle does **not**
   flip bit 6 of `WRAM:+$1D4E`. This refutes Trial 0's retrospective
   inference, and the block location must be re-derived from the live
   diffs alone.
2. A setting's storage lies outside `WRAM:+$1D4D`–`+$1D4F`. The
   full-WRAM diff catches this; it widens the map rather than ending the
   unit.
3. Toggling a setting changes no WRAM byte at all.

## Evidence paths

`local_artifacts/experiments/EXP-0041/` — **17 artifacts, 616 KB**, with
`hashes.sha256` over all of them, written at close, after which the
directory is frozen (`docs/research/EVIDENCE_LAYOUT.md`): 8 screenshots,
4 savestates (`t1-batmode-bit3set`, `t1-batmode-bit3clear`,
`t9-all-toggled`, `t7-cursor-memory-clean`), `trial0-cursor-diff.txt`,
`bridge-commands.log`, `bridge-events.log`, `experiment.json`.

Savestates and screenshots are ROM-derived and remain local; the tracked
record cites paths and hashes only.

Instead of per-trial raw WRAM dumps, savestates were captured at the
decision points and compared with `ff6lab state diff`, which reads the
128 KB work-RAM image straight out of the `.mss`. That is the same
evidence at a fraction of the size, and it is what made the
default-versus-all-toggled comparison practical.

## Trials

**Trial 0 — retrospective, no emulator.** EXP-0040's savestate lineage
preserved a single-config-change pair: `pre-whelk-healed.mss` →
*(Config: Cursor Reset → Memory)* → `pre-whelk-healed-wait.mss`. Diffed
offline with `ff6lab state diff`.

**Trials 1–9 — controlled toggles.** One setting each, in Config-screen
order: Bat.Mode, Bat.Speed, Msg.Speed, Cmd.Set, Gauge, Sound, Cursor,
Reequip, Controller. Each begins from a fresh reload of the starting
savestate.

## Observations

### Trial 0 — retrospective (no emulator)

`ff6lab state diff pre-whelk-healed.mss pre-whelk-healed-wait.mss wram`:
**45 differing bytes in 27 runs**, among them a single isolated bit:
`WRAM:+$1D4E` `$00` → `$40`. The block `WRAM:+$1D4D`–`+$1D4F` read
`2A 00 00` / `2A 40 00`. `cart.saveRam` was **identical** in both.

Recorded at `trial0-cursor-diff.txt`.

### Starting state

Loading milestone 05 gave player tile `(26,1C)` and
`WRAM:+$1D4D`–`+$1D4F` = `2A 00 00` — the same values as the pre-Whelk
lineage's unmodified side, from an entirely independent savestate.
Opening the field menu, moving the cursor and entering Config changed
none of them: **navigation alone writes nothing**.

### Trial 1 — Bat.Mode

| Action | `WRAM:+$1D4D` |
|---|---|
| on entry | `$2A` |
| press RIGHT | `$2A` (no change — already the right-hand option) |
| press LEFT | `$22` |
| press RIGHT | `$2A` |

Bit 3 alone. Full-WRAM diff between savestates taken in the two states —
**17 differing bytes in 7 runs** — showed `+$1D4D` plus two text runs at
`+$39A6` and `+$39B6`, which decode under EXP-0026's text encoding to
"`ctive`" and "`ait`": the Config screen's own "Active" / "Wait" cells,
each interleaved with an attribute byte that swaps `$28` ↔ `$20` when
bit 3 flips.

### Trials 2–9 — one setting each

Sequential right-presses down the screen, reading `+$1D4D`–`+$1D4F`
before and after each:

| Setting | before → after | bit moved |
|---|---|---|
| Bat.Speed | `2A 00 00` → `2B 00 00` | `+$1D4D` bit 0 |
| Msg.Speed | `2B 00 00` → `3B 00 00` | `+$1D4D` bit 4 |
| Cmd.Set | `3B 00 00` → `BB 00 00` | `+$1D4D` bit 7 |
| Gauge | `BB 00 00` → `BB 80 00` | `+$1D4E` bit 7 |
| Sound | `BB 80 00` → `BB A0 00` | `+$1D4E` bit 5 |
| Cursor | `BB A0 00` → `BB E0 00` | `+$1D4E` bit 6 |
| Reequip | `BB E0 00` → `BB F0 00` | `+$1D4E` bit 4 |
| Controller | `BB F0 00` → `BB F0 00` | **none in this window** |

Controller appeared to write nothing. Widening the read to eight bytes
found it: `WRAM:+$1D54` `$00` → `$80`, confirmed by an isolated
LEFT/RIGHT pair (`80` → `00` → `80`). The apparent null result was the
three-byte read window, not a missing write — the full-WRAM net is what
caught it.

The default-versus-all-toggled full-WRAM diff (89 bytes, 31 runs) shows
exactly three storage bytes — `+$1D4D`, `+$1D4E`, `+$1D54` — with every
other run being Config-screen text cells that decode to `indow`/`hort`,
`n`/`ff`, `tereo`/`ono`, `eset`/`emory`, `ptimum`/`mpty`,
`ingle`/`ultiple`, plus the two speed-digit attribute bytes.

### Numeric field extents

Swept rather than inferred, from a fresh reload:

- **Bat.Speed** (`+$1D4D` low bits): `2A → 29 → 28`, clamped at `28`;
  then `28 → 29 → 2A → 2B → 2C → 2D`, clamped at `2D`. Six distinct
  values, **stored 0–5**, matching the six displayed positions. Bit 3
  held set throughout.
- **Msg.Speed**: `2D → 1D → 0D`, clamped at `0D`; then up to `5D`,
  clamped. Six values, **stored 0–5**. Bit 7 and bits 0–3 undisturbed.

### Falsifier run

Fresh reload → field menu → Config → cursor to the Cursor row →
one RIGHT press, nothing else: `WRAM:+$1D4E` `$00` → `$40`. The Config
screen's "Reset" cells took attribute `$28` and its "Memory" cells `$20`.

## Interpretation

The nine Config settings occupy **three WRAM bytes**, not one contiguous
block:

**`WRAM:+$1D4D`** — default `$2A`

| Bits | Setting | Encoding |
|---|---|---|
| 0–2 | Bat.Speed | 0–5, displayed 1–6 (default 2 → "3") |
| 3 | Bat.Mode | 0 = Active, 1 = **Wait** (default 1) |
| 4–6 | Msg.Speed | 0–5, displayed 1–6 (default 2 → "3") |
| 7 | Cmd.Set | 0 = Window, 1 = Short |

**`WRAM:+$1D4E`** — default `$00`

| Bit | Setting | Encoding |
|---|---|---|
| 4 | Reequip | 0 = Optimum, 1 = Empty |
| 5 | Sound | 0 = Stereo, 1 = Mono |
| 6 | Cursor | 0 = Reset, 1 = Memory |
| 7 | Gauge | 0 = On, 1 = Off |
| 0–3 | — | untouched by these nine settings; **Unknown** |

**`WRAM:+$1D54`** — default `$00`

| Bit | Setting | Encoding |
|---|---|---|
| 7 | Controller | 0 = Single, 1 = Multiple |
| 0–6 | — | untouched by these nine settings; **Unknown** |

In every case the **cleared** bit is the left-hand screen option and the
set bit is the right-hand one, which is why a single RIGHT press from a
default state sets a bit and a LEFT press clears it.

**The Config screen's highlighting convention is now settled**, and it is
the inverse of the intuitive reading: the **selected** option's text
cells carry attribute **`$20`** and the unselected option's carry
**`$28`**. Correlating that with the bits confirms EXP-0040's correction
independently — with `+$1D4D` bit 3 set, "Wait" carries `$20` and
"Active" carries `$28`, so **Wait was indeed already selected**, and the
original "Active" reading was wrong. The hand cursor marks the row only.

`Bat.Mode = Wait` is therefore the state a new game arrives in on this
ROM revision, not an operator change.

Trial 0's retrospective inference — bit 6 of `+$1D4E` as the Cursor
setting — is **verified**: an isolated controlled toggle from a clean,
unrelated savestate produced exactly the `$00` → `$40` transition seen in
EXP-0040's preserved pair.

## Alternatives

- **The three bytes could be a battle-local copy rather than the
  persistent store.** Nothing here distinguishes them, because no battle
  was entered and no save event occurred. What is established is that the
  Config screen reads and writes *these* bytes. Whether battle code
  consults them directly or a copy is EXP-0042's question, and it is the
  reason that experiment comes next.
- **Bits 0–3 of `+$1D4E` and 0–6 of `+$1D54` may hold further settings**
  not exposed on this screen, or unrelated state. This unit only shows
  the nine settings do not touch them.
- **The `$20`/`$28` attribute pair may encode palette selection
  generally** rather than "selected/unselected" specifically. The
  correlation is exact across all seven two-option rows observed, but the
  mechanism is a menu-renderer detail this unit did not investigate.
- **Field widths for the two speed settings** are read from clamping
  behaviour, which bounds the *value* range at 0–5. The three-bit field
  boundary additionally rests on the neighbouring bits (3 and 7) being
  independently owned, which was observed directly.
- The single-byte-per-press pattern could in principle mask a second
  write that is immediately reverted within the same frame. Only
  post-hoc state was sampled; no write-watch was installed.

## Result

**All nine Config settings located and encoded.** The primary question is
answered: `Bat.Mode` is **bit 3 of `WRAM:+$1D4D`** (set = Wait) and
`Bat.Speed` is **bits 0–2 of `WRAM:+$1D4D`** (0–5, displayed 1–6).

The battle-configuration fingerprint is now machine-readable, which is
what the ACTIVE-versus-WAIT comparison needs in order to set and verify
its controlled variable. No falsifier fired.

Secondary results: the Config screen's selected/unselected attribute
convention is settled and inverts the intuitive reading; EXP-0040's
`Bat.Mode` correction is independently confirmed; the configuration is
**not** SRAM-backed before a save event (`cart.saveRam` identical across
a config change, and virgin throughout).

## Confidence

- Byte and bit for all nine settings, and their set/clear meanings:
  **Confirmed** (controlled one-variable toggles, full-WRAM verified,
  reproduced from two independent savestate lineages).
- Bat.Speed and Msg.Speed as 0–5 fields at bits 0–2 and 4–6:
  **Confirmed** (both bounds swept, clamping observed, neighbouring bits
  undisturbed).
- Default configuration `$2A/$00/$00` for a new game on this revision:
  **Confirmed** for this route lineage; new-game initialization itself
  was not captured (CEN-SAVE-0001 remains open).
- Config screen attribute `$20` = selected, `$28` = unselected:
  **Confirmed** for the seven two-option rows observed.
- Whether battle code reads these bytes directly or a copy: **Unknown** —
  deliberately not investigated.
- SRAM persistence: **Unknown** — no save event has ever been captured.
- Meaning of `+$1D4E` bits 0–3 and `+$1D54` bits 0–6: **Unknown**.

## Stopping condition

The unit ends when all nine settings are mapped or explicitly bounded.
Specifically **out of scope**, each belonging to a later unit:

- Tracing the config writer's callers or the Config menu's input loop.
- Entering a battle.
- Testing whether battle entry copies configuration into a battle-local
  cell — that is EXP-0042, and it is the next question, not this one.
- SRAM persistence, which cannot be observed until a save event is
  captured (CEN-SAVE-0001).

## Next action

**EXP-0042 — battle-entry configuration sampling.** Are `WRAM:+$1D4D`
and `+$1D4E` read directly by battle code, or copied into a battle-local
cell at battle entry? The answer decides whether toggling configuration
mid-battle is a legitimate experimental technique, and therefore how
every later ATB experiment must be staged.

Method: install a **read**-watch on `WRAM:+$1D4D`–`+$1D4E` (the bridge's
`watchwrites` helper covers writes; a read variant is needed) across a
battle entry from a scripted encounter, and capture the reading PCs. The
mines random encounter (milestone 06, EXP-0038) is reproducible on demand
and needs no Whelk state.

Deferred from this unit, each with a known entry point:

- The config writer's routine and callers — reachable from a write-watch
  on `+$1D4D` during a toggle.
- `+$1D4E` bits 0–3 and `+$1D54` bits 0–6.
- SRAM persistence, blocked on capturing a save event (CEN-SAVE-0001).
- The Config menu's input loop and the `$20`/`$28` attribute mechanism
  (CEN-MENU-0001 follow-up).
