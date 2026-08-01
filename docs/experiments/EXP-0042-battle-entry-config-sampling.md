# EXP-0042: Battle-entry configuration sampling — direct reads or a copy?

## Question

Do the battle routines read `WRAM:+$1D4D`/`+$1D4E` (EXP-0041's
configuration bytes) **directly during battle**, or is the configuration
**sampled once at battle entry** into a battle-local copy?

This is the second unit of the ATB research program, and it is a
staging question rather than a mechanism question. The answer decides:

- whether changing configuration mid-battle takes effect at all, and so
  whether mid-battle toggling is a legitimate experimental technique;
- whether an ACTIVE-versus-WAIT comparison must be established **before**
  battle entry, or can be flipped within a single battle;
- where later experiments must place their controlled variable.

Getting this wrong would silently invalidate every ACTIVE/WAIT
measurement that follows, which is why it comes before any timing work.

## Starting state

`local_artifacts/scenarios/SCN-0001/05-mines-entry/05-mines-entry.mss` —
milestone `05-mines-entry`, frame 51 578, EXP-0036, three byte-identical
power-on runs. The mines corridor carries a random-encounter zone that
EXP-0038 resolved to formation 14, so a battle entry is reachable by
walking, with no Whelk state involved.

Fallback if the encounter proves slow to trigger: the scripted battle
reachable from `02-narshe-entry.mss` (battle entry at frame 31 557,
EXP-0032), which fires on a fixed schedule.

## ROM identity

`Final Fantasy III (USA).sfc`, sha256
`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1 macOS, **headless** (`--testrunner --timeout=7200`,
`FF6_OUT` set), bridge v2, one instance. Lab controls unchanged:
`Snes.RamPowerOnState=AllZeros`, virgin SRAM.

## Battle configuration

| Setting | Value | Source |
|---|---|---|
| Bat.Mode | `Wait` | memory-read |
| Bat.Speed | 3 | memory-read |
| Msg.Speed | 3 | memory-read |
| Cmd.Set | `Window` | memory-read |
| Gauge | `On` | memory-read |
| Sound | `Stereo` | memory-read |
| Cursor | `Reset` | memory-read |
| Reequip | `Optimum` | memory-read |
| Controller | `Single` | memory-read |

`WRAM:+$1D4D` = `$2A`, `WRAM:+$1D4E` = `$00`, `WRAM:+$1D54` = `$00`
(EXP-0041). Configuration is held at its route default and **not
changed** during this unit — the independent variable here is battle
entry, not a setting.

## Independent variable

Battle entry: the transition from field to battle. Nothing else changes.

## Controlled variables

Same ROM, emulator, lab settings and starting savestate. Configuration
untouched throughout. No menu is opened. Only walking input is supplied,
to trigger the encounter.

## Instrumentation

- `watchreads` added to `mesen/probes/common.lua` — a read-watch
  counterpart to the existing `watchwrites`, aggregating per reading PC.
  Read watches fire far more often than write watches, so per-PC
  aggregation with `count`, `firstFrame` and `lastFrame` is the only
  shape that stays bounded while still preserving ordering.
- `watchdump` added alongside it, so both helpers can report their
  tables; `watchwrites` previously had no dump path.
- `mesen/probes/EXP-0042.lua`: read-watch over `WRAM:+$1D4D`–`+$1D4E`,
  plus a battle-entry timestamp taken from the first `+$3B18`–`+$3BB7`
  write issued from `ROMCPU:$C22800`–`$C22FFF` (EXP-0032's detector), so
  every read can be sorted into before-entry and after-entry.

The battle-entry frame is the discriminator; without it the read PCs
cannot be attributed to field code or battle code.

## Expected outcomes

- *Sampled-at-entry:* one or more read PCs fire a **small, bounded**
  number of times in a burst at the battle-entry frame, and no PC reads
  the bytes during the battle proper.
- *Consulted-continuously:* at least one read PC keeps firing throughout
  the battle, with a count that scales with battle duration.
- *Mixed:* an entry-time sample plus a separate consumer that re-reads
  the persistent bytes — for example message speed read per dialogue box
  while battle timing uses a copy.
- *Neither:* no read of these bytes occurs across a battle entry at all,
  which would mean battle code reaches the configuration by some other
  path and would redirect the unit to finding it.

## Falsifying outcome

1. **For "sampled at entry":** any read PC whose count keeps rising
   while the battle is in progress.
2. **For "consulted continuously":** every read PC's `lastFrame` sits at
   or before the battle-entry frame.
3. If the read watch records **zero** hits across a confirmed battle
   entry, the instrumentation is suspect before the conclusion is —
   `WRAM:+$1D4D` lies below `$2000` and is therefore inside the low-RAM
   mirror, so the watch must catch mirror-bank access. EXP-0037's
   event-flag watch over `+$1E80`–`+$1EDF` establishes that
   `snesWorkRam` callbacks do fire for that region, so a zero result
   would be a real observation rather than a known artifact — but it
   must be checked against a deliberate control read before it is
   believed.

## Evidence paths

`local_artifacts/experiments/EXP-0042/` — probe output, per-phase read-PC
tables, savestates at the phase boundaries, `bridge-commands.log`,
`bridge-events.log`, `experiment.json`, and `hashes.sha256` written at
close, after which the directory is frozen.

## Trials

**Run 1** — route default (`+$1D4D` = `$2A`: Wait, Bat.Speed 3). Walked
the mines corridor from milestone 05; encounter fired at frame **54 799**,
formation `$000E` = **14**, the same formation EXP-0038 established. The
battle was then fought with A-presses so actions, animations and damage
display all ran.

**Run 2** — one deliberate configuration change before entry:
`Bat.Mode → Active`, `Bat.Speed → 6`, giving `+$1D4D` = `$25` (bits 0–2 =
5, bit 3 clear, Msg.Speed and Cmd.Set untouched). Encounter at frame
**55 922**.

## Observations

### Field baseline

Standing still in the field for ~900 frames produced **zero** reads of
either byte. They are not polled by field code.

### Run 1 — reads across the battle boundary

Battle entry at frame 54 799, detected at `ROMCPU:$C22899`.

| Reading PC | Address | Count | First | Last |
|---|---|---|---|---|
| `$C22475` | `$7E1D4D` | **1** | 54 799 | 54 799 |
| `$C22493` | `$7E1D4E` | **1** | 54 799 | 54 799 |
| `$C10FFA` | `$7E1D4E` | 1 | 54 804 | 54 804 |
| `$C159DA` | `$001D4E` | 2 | 55 039 | 60 704 |
| `$C198B0` | `$001D4D` | 2 | 60 868 | 61 217 |

The last two read through the **low-RAM mirror** (`$001D4x`), which
confirms the watch catches mirror-bank access and disposes of falsifier 3.

`$C198B0` appeared only after fighting began, and a screenshot at frame
61 508 confirms the battle was still in progress (command menu open,
three enemies, party at 4/0/20) — so those are battle-time reads, not
post-battle field code.

### Static decode of the entry-time readers

`read cpu C22460 80` gives, from `$C22472`:

```asm
C22472  AD 4D 1D    LDA $1D4D      ; config byte 1
C22475  30 03       BMI +3         ; bit 7 = Cmd.Set
C22477  9C 2E 2F    STZ $2F2E      ;   Window -> clear $2F2E
C2247A  89 08       BIT #$08       ; bit 3 = Bat.Mode
C2247C  F0 03       BEQ +3
C2247E  EE 8F 3A    INC $3A8F      ;   Wait -> increment $3A8F
C22481  29 07       AND #$07       ; bits 0-2 = Bat.Speed
C22483  0A 0A 0A    ASL ASL ASL    ; x8
C22486  85 EE       STA $EE
C22488  0A          ASL            ; x16
C22489  65 EE       ADC $EE        ; x16 + x8 = x24
C2248B  49 FF       EOR #$FF       ; one's complement
C2248D  8D 90 3A    STA $3A90
C22490  AD 4E 1D    LDA $1D4E      ; config byte 2
C22493  10 03       BPL +3         ; bit 7 = Gauge
C22495  9C 21 20    STZ $2021      ;   Off -> clear $2021
```

and at battle init, `$C10FF7 LDA $1D4E / AND #$07 / STA $2F34`.

So battle entry does not copy the bytes; it **decomposes and transforms**
them into separate battle-local cells.

### Live confirmation of the transform

Read back after each run:

| Configuration | `+$1D4D` | `$3A8F` | `$3A90` |
|---|---|---|---|
| Wait, Bat.Speed 3 (stored 2) | `$2A` | `01` | `$CF` = 255 − 24×2 |
| Active, Bat.Speed 6 (stored 5) | `$25` | `00` | `$87` = 255 − 24×5 |

Both values were **predicted from the disassembly before run 2 was
executed** and matched exactly. `$2F2E` = `00` (Cmd.Set Window),
`$2021` = `$FF` (Gauge On, the STZ not taken), `$2F34` = `00`.

The field value of `$3A8F`/`$3A90` before entry was `FF FF`, so battle
entry initialises both.

### Static decode of the in-battle readers

```asm
C198AC  AF 4D 1D 00  LDA $001D4D
C198B0  4A 4A 4A 4A  LSR x4          ; >>4
C198B4  29 07        AND #$07        ; = Msg.Speed
C198B6  AA           TAX
C198B7  BF 72 98 C1  LDA $C19872,X   ; delay table
C198BB  48 / JSR $022A / PLA / DEC / BNE -8
```

```asm
C159D6  AF 4E 1D 00  LDA $001D4E
C159DA  29 40        AND #$40        ; = Cursor
C159DC  D0 0B        BNE +$0B
C159DE  ...          STZ $890F,X for $5C bytes   ; Reset -> clear cursor memory
```

`$C198AC` consumes **Msg.Speed** as a message-delay loop count via a
table at `ROMCPU:$C19872`. `$C159D6` consumes **Cursor**, clearing a
`$5C`-byte block at `WRAM:+$890F` when the setting is Reset — which is
the mechanism behind EXP-0040's observation that `Cursor = Memory`
reopens a character's list on their last-used ability.

### Run 2 — the Config menu's own read path

Run 2's table additionally contains ~20 PCs in bank `$C3`
(`$C33A8F`–`$C341CA`) firing between frames 51 771 and 52 729 — the
Config menu reading the bytes to render and toggle them. Their captured
values step `$2A` → `$22` → `$25` in exactly the order the settings were
changed. These are menu-time reads, cleanly separated from battle entry
at 55 922, and are not part of this unit's question.

In run 2 `$C198B0` does not appear at all: no battle message had been
displayed by the capture point, consistent with it being a message-pacing
consumer rather than a timing one.

## Interpretation

The answer is **mixed, and the split falls exactly along the line that
matters**.

**Battle-timing settings are sampled once, at battle entry.**
`Bat.Mode` and `Bat.Speed` are read a single time each, at the
battle-entry frame, by `$C22472`, and converted into battle-local state:

- `WRAM:+$3A8F` — incremented at entry **iff** Bat.Mode = Wait, so it is
  the battle-local WAIT flag (`01` = Wait, `00` = Active observed).
- `WRAM:+$3A90` — `255 − 24 × BatSpeed`, with BatSpeed the stored 0–5
  value. Fast (stored 0) gives `$FF`, Slow (stored 5) gives `$87`.

Nothing re-reads `+$1D4D` for timing purposes during the battle.

**Presentation settings are read live.** `Msg.Speed` and `Cursor` are
consulted from the persistent bytes while the battle runs, by `$C198AC`
and `$C159D6` respectively. So "is the configuration copied at battle
entry?" has no single answer: it depends on the setting, and the two
groups behave differently.

**Consequence for the ATB program — the reason this unit exists.**
Changing `Bat.Mode` or `Bat.Speed` during a battle cannot affect that
battle's timing, because the battle is running off `$3A8F`/`$3A90`. Every
ACTIVE/WAIT or Battle Speed condition must therefore be established
**before battle entry**. Equivalently — and far more precisely — an
experiment may patch `$3A8F`/`$3A90` directly, which is a much better
controlled handle than driving menus, and one that can be changed
mid-battle to test whether the consumers re-read them.

### Contextual observation — not a conclusion

`$3A90 = 255 − 24 × BatSpeed` is monotonically **decreasing** in the
displayed speed number, and the displayed scale is labelled "Fast" at 1
and "Slow" at 6. If `$3A90` were a threshold to be reached, that
direction would be inverted; it is consistent with an **addend or
reload value** where a larger number advances something faster. That
makes it a strong candidate for the Battle Speed input to the ATB rate.

This is registered as a **candidate correlation only**. Nothing here
observes a gauge, a tick, or any consumer of `$3A90`. Identifying its
reader is the natural next step and belongs to the next unit, not this
one.

## Alternatives

- **A consumer might re-read `+$1D4D` for timing on a path neither run
  exercised** — a battle type other than a random encounter, a status
  change, or a menu state not opened here. Both runs were single random
  encounters of formation 14; scripted, pincer, back and boss battles are
  untested.
- **`$3A8F` is an `INC`, not a store.** Both observed values (`00`, `01`)
  are consistent with "zeroed at init, incremented once when Wait", but a
  path that increments it a second time would break that reading. Only
  two configurations were observed.
- **`$3A90`'s arithmetic is Confirmed; its meaning is not.** The
  transform was verified by prediction, but no consumer of `$3A90` has
  been located, so calling it a speed *rate* remains inference.
- **`$C198AC` may not be the only Msg.Speed consumer**, and `$C159D6`
  may not be the only Cursor consumer. The watch reports readers of the
  config bytes, not consumers of the derived state.
- **`+$1D4E` bits 0–2 → `$2F34`** (`$C10FF7`) involves bits EXP-0041
  showed the nine Config settings never touch. Something else owns them;
  this unit only establishes that battle init consumes them.

## Result

**Answered: mixed, split by setting.**

- `Bat.Mode` and `Bat.Speed` — **sampled once at battle entry**
  (`ROMCPU:$C22472`) into `WRAM:+$3A8F` (Wait flag) and `WRAM:+$3A90`
  (`255 − 24 × BatSpeed`). Not re-read during battle.
- `Cmd.Set` → `$2F2E`, `Gauge` → `$2021`, `+$1D4E` bits 0–2 → `$2F34`,
  also at entry.
- `Msg.Speed` and `Cursor` — **read live during battle** from the
  persistent bytes, by `$C198AC` and `$C159D6`.

The staging rule for the ATB program follows directly: ACTIVE/WAIT and
Battle Speed conditions must be set **before battle entry**, or injected
at `$3A8F`/`$3A90`.

Falsifier 1 held for the timing settings (no rising counts on
`$C22475`/`$C22493`); falsifier 2 held for the presentation settings
(`$C198B0` and `$C159DA` fire after entry); falsifier 3 was disposed of
by the observed mirror-bank reads.

## Confidence

- `Bat.Mode`/`Bat.Speed` read exactly once at battle entry by
  `ROMCPU:$C22472`, not re-read for timing: **Confirmed** for random
  encounters (two runs, two configurations, ~6 700 battle frames with
  actions resolving).
- `WRAM:+$3A8F` = battle-local Wait flag, set at entry from bit 3:
  **Confirmed** (predicted and observed under both Active and Wait).
- `WRAM:+$3A90` = `255 − 24 × BatSpeed`: **Confirmed** (arithmetic
  decoded statically, predicted before the second run, matched exactly at
  two points).
- `$3A90` as the Battle Speed input to ATB rate: **Tentative
  hypothesis** — a contextual observation, no consumer located.
- `Msg.Speed` and `Cursor` read live during battle: **Confirmed**
  (runtime reads plus decoded consumers).
- `$2F2E`, `$2021`, `$2F34` as the Cmd.Set / Gauge / `+$1D4E`-low
  destinations: **Confirmed** (code role) — their consumers are Unknown.
- Behaviour in battle types other than a random encounter: **Unknown**.

## Stopping condition

The unit ends once the read PCs are attributed to before/after battle
entry and the sampling question is answered at Confirmed or bounded.

Out of scope, each belonging to a later unit:

- Disassembling the reading routine's full context or its callers beyond
  what is needed to identify the copy destination, if one exists.
- Locating the ATB gauges (the next unit).
- Any timing measurement, ACTIVE/WAIT comparison, or Whelk work.
- The configuration **writer** routine.

## Next action

**EXP-0043 — locate the ATB gauges and the consumer of `WRAM:+$3A90`.**

`$3A90` is now the sharpest lead the project has into battle timing: a
battle-local value derived from Battle Speed, whose reader is unknown.
Find it with a read-watch on `+$3A8F`–`+$3A90` across a battle, then
follow that routine to whatever it advances. That is very likely the ATB
tick site, and it converges with the other standing lead — the **eight
undumped callees of `ROMCPU:$C101FB`** (open question #6), one of which
should be the same routine.

Substrate needed: none new. `local_artifacts/experiments/EXP-0042/in-battle-formation14.mss`
is a live battle preserved by this unit, and the milestone-03/06 WRAM
dumps are diffable offline with `ff6lab state`.

Deferred, each with a known entry point:

- Consumers of `$2F2E`, `$2021`, `$2F34`.
- Whether `$3A8F` can be incremented more than once.
- Owner of `+$1D4E` bits 0–3 and `+$1D54` bits 0–6.
- Whether battle types other than random encounters sample the
  configuration differently — worth re-checking at the first scripted or
  boss battle rather than as its own unit.
- The Config menu's own read path (bank `$C3`, ~20 PCs captured
  incidentally in run 2) — belongs to CEN-MENU-0007 follow-up.
