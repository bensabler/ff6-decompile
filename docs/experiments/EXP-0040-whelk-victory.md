# EXP-0040: Whelk victory (branch A — head-only)

- **Status:** in progress (2026-08-01)
- **Mode:** bounded **gameplay verification**. One objective only: defeat
  Whelk correctly and capture the first stable post-battle state. This is
  *not* a Whelk data, AI, formation, graphics, audio, event or monster-id
  investigation.
- **Program:** SCN-0001 (opening-to-Whelk). Serves B18 (Whelk battle) and
  B19 (stable post-Whelk state), and establishes milestone
  `10-whelk-victory`.

## Question

Can Whelk be defeated reliably from the preserved pre-Whelk state by
healing the party first and attacking **only while the head is exposed**?

## Starting state

`local_artifacts/scenarios/SCN-0001/07-pre-whelk/pre-whelk-recon.mss`

- SHA-256 `852c82a0c65551b36931b131a685a053d52331ded904305fb0beeeeade960858`
- 129 795 bytes; preserved by EXP-0039 at mines position `(2A,09)`,
  immediately before the scripted Vicks beat that precedes Whelk.
- Party HP at preservation: **26 / 19 / 56** (EXP-0039). MP unrecorded —
  captured on load by this unit.

## ROM identity

`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1, **GUI (operator-visible)**, foreground, normal speed.
Emulator lab controls unchanged: `Snes.RamPowerOnState=AllZeros`,
virgin SRAM (verified still empty after shutdown). Exactly one Mesen
instance ran at a time. Headless was not used.

In-game configuration observed at
`local_artifacts/experiments/EXP-0040/34-config-before.png`:
**Bat.Mode `Wait`**, Bat.Speed 3, Msg.Speed 3, Cmd.Set `Window`,
Gauge `On`, Sound `Stereo`, Cursor `Reset`, Reequip `Optimum`,
Controller `Single`. This unit changed exactly one of them:
**`Cursor: Reset` → `Memory`** (see Amendments). `Bat.Mode` was already
`Wait` and was not changed.

## Independent variable

Target selection discipline: attack **only** when the head is visibly
extended; take a non-damaging action while it is retracted. EXP-0039's
defeat was the uncontrolled contrast (blind A-mashing into the shell).

## Controlled variables

Same ROM, same lab settings, same starting savestate for every attempt.
Each retry reloads the identical `.mss`, so attempts differ only in the
input sequence.

## Instrumentation

Bridge commands only (`loadstate`, `press`, `read`, `screenshot`,
`eval` for savestates). No new watches, no exec census, no disassembly.

Reads use only already-Confirmed addresses:

| Quantity | Address | Source |
|---|---|---|
| Field position | `WRAM:+$00AF`/`+$00B0` | EXP-0035 |
| Formation id (staged) | `WRAM:+$11E0` | EXP-0029/0030 |
| Staged formation record | `WRAM:+$3F44` (15 bytes) | EXP-0030 |
| Battle current HP, slots 0-9 | `WRAM:+$3BF4` (2 B × 10) | DISC-0001 |
| Battle max HP, slots 0-9 | `WRAM:+$3C1C` (2 B × 10) | DISC-0001 |
| Battle MP candidates, slots 0-9 | `WRAM:+$3C08` / `+$3C30` | DISC-0001 |
| Field HP / MP writeback | `WRAM:+$1609` / `+$160D` | EXP-0033 |

Slots 0-3 are party, 4-9 enemies (DISC-0001), so Whelk's HP is readable
live at `+$3BF4 + 8`. Enemy HP is therefore **observed**, not inferred.

## Healing strategy

Heal Force is a Magitek ability (CEN-MAGIC-0001, EXP-0039) and is
available **in battle only** — the field has no Magitek menu. The party
therefore enters at its carried HP, and healing happens inside the Whelk
fight as the first actions:

1. Load the state, confirm position and party HP/MP.
2. Enter the Whelk battle.
3. Spend the opening turns on **Heal Force** until the party is safely
   above the observed one-counter lethality band (EXP-0039: a shell
   counter killed a 19 HP member outright and took the leader from ~40s
   to 14).
4. Only then begin damaging actions.

If a field healing route turns out to exist (item, save point, or an
out-of-battle Magitek menu), take it before contact and record it. The
in-battle plan is the fallback that does not depend on that discovery.

## Attack strategy

- Advance the introduction/warning dialogue with discrete presses,
  reading the screen between them.
- Before **every** damaging action, take a screenshot and classify the
  head/shell state. Never issue an attack on an unclassified frame.
- Attack only while the head is extended.
- When the head retracts: stop attacking; use Heal Force, Defend, or
  another safe action visible in the command set; do not mash through
  command selection.
- Resume attacks only after confirming the head has returned.
- Do **not** deliberately test the shell counter — EXP-0039 already
  recorded it (branch B).
- Continue until Whelk's HP (`+$3BF4 + 8`) reaches 0 and victory
  processing runs.

## Head/shell state criteria

Classification is **visual, from the captured frame**, corroborated by
Whelk's HP response where available:

- **Head extended:** the head sprite is drawn clear of the shell body;
  the head is a selectable target in the target cursor.
- **Head retracted:** only the shell mass is drawn; the head sprite is
  absent or withdrawn against the shell.
- **Corroboration:** damage landing on the extended head should reduce
  the enemy HP word; a strike into the shell is expected to produce the
  counter observed in EXP-0039 rather than progress.

Ambiguous frames are treated as **retracted** (the safe classification)
and re-sampled.

## Expected Whelk transitions

Head extended → (retract) → shell only → (re-emerge) → head extended,
cycling. Timing, trigger and any AI-driven ordering are explicitly **not**
measured here; the unit only needs to *observe* which state is current
before each action. A deterministic timing model is deferred.

## Expected victory processing

The EXP-0033 pattern: enemy HP reaches 0, death handling runs, the
battle ends with EXP/GP award boxes, and post-battle HP/MP are written
back into the field character block (`$C2496E`/`$C24979` →
`+$1609`/`+$160D`).

## Expected stable post-battle state

Control returns to the field near `(2A,07)` with the party alive, the
battle arrays no longer authoritative, and the scenario able to continue.
EXP-0039 saw a further scripted beat at `(2A,07)` ("GUARD: We won't hand
over the Esper!!") **after a defeat-and-reload**; whether it is the
normal post-victory progression is an open question this unit's clean
run can settle incidentally, and must not be assumed.

## Expected outcomes

- *Supports:* Whelk defeated with head-only targeting; victory
  processing and a stable post-battle state captured; milestone
  `10-whelk-victory` established.
- *Partial (valid):* Whelk defeated only after retries — the failed
  attempts are preserved as evidence with their tactical cause.
- *Refutes:* head-only targeting is insufficient (see falsifiers).

## Falsifiers

1. **Head-only targeting does not win.** Repeated attempts in which
   every damaging action lands on a visibly extended head still end in
   defeat — the strategy is insufficient and the fight needs a mechanic
   not yet identified.
2. **The head/shell state is not visually classifiable.** If captured
   frames cannot be reliably sorted into extended/retracted, the stated
   criteria fail and the unit must switch to a memory-derived
   discriminator (a bounded follow-up, not an inline decode).
3. **Attacking the extended head does not reduce enemy HP.** Would
   refute the premise that the head is the damageable target.
4. **Victory does not yield a stable post-battle state** — e.g. the
   battle exits into another battle, a soft-lock, or a state that cannot
   be re-reached — leaving B19 open.

## Evidence requirements

Under `local_artifacts/experiments/EXP-0040/`, with `hashes.sha256`:

- Attempt transcript with per-action frame, party HP, Whelk HP, head
  state, command chosen, and result.
- Screenshots: pre-Whelk load, healed party, battle entry, head
  extended, head retracted, head returned, final damaging action, Whelk
  death, victory/reward processing, battle exit, first stable
  post-battle state.
- Savestates: refreshed pre-battle state, and the post-battle milestone
  state.
- Any failed attempt preserved with its cause.

Milestone artifacts additionally go to
`local_artifacts/scenarios/SCN-0001/10-whelk-victory/`.

Restricted files stay local; only metadata, addresses, sizes and hashes
are tracked.

## Stopping condition

Stop as soon as **the first correct victory and its first stable
post-battle state are captured** — do not continue playing, do not
re-fight for a better run, and do not start branch B or any deeper Whelk
investigation. Also stop if falsifier 1 or 2 is reached, recording the
exact blocking fact and the smallest experiment that would resolve it.

## Bounds (strict anti-depth rules)

Do **not** pause the fight to investigate: formation 432's high-bit
monster-id extension, the Whelk AI script, the AI interpreter, the
monster name table, graphics or audio provenance, the event-script
interpreter, map formats, encounter RNG, or the full post-Whelk story
sequence. Register each new observation briefly (census + one bounded
future question) and return to the fight. The unresolved formation-id
extension is a future research question and **must not block this
milestone**.

## Amendments made during the run

1. **Field healing exists.** The pre-registered plan assumed healing had to
   happen in battle. A field menu *is* available before contact and the
   party carried **Tonic ×4** and **Potion ×1**, so healing was done on
   the field at no turn cost. The in-battle Heal Force plan became the
   secondary channel.
2. **Battle-mode reading corrected.** This record initially stated the
   lab was in **Active** mode and that the unit switched it to Wait.
   That was an **operator-error misreading of the Config screen**: the
   hand cursor sits at the row's left edge next to the first option and
   does *not* indicate the selection. Re-inspection of the same
   screenshot (`34-config-before.png`) shows **"Active" rendered grey and
   "Wait" rendered white** — i.e. **Bat.Mode was already `Wait`**, before
   and throughout both attempts. The right-press on that row was a no-op.
   The only configuration this unit actually changed is
   **`Cursor: Reset` → `Memory`**. The correction was prompted by the
   operator and confirmed against the preserved capture.
3. **Attack element switched.** Attempt 1 used Bolt Beam; attempt 2 used
   Fire Beam after the intro dialogue described Whelk as eating lightning
   and storing it in the shell.

## Trials

Two piloted GUI attempts on 2026-08-01, both from the preserved
pre-Whelk state, both under `Bat.Mode = Wait`.

| Attempt | Origin state | Attack used | Head damage dealt | Outcome |
|---|---|---|---|---|
| 1 | `pre-whelk-recon.mss` → field-healed | Bolt Beam | 1600 → 909 (4 hits) | stopped; reloaded to inspect Config |
| 2 | `pre-whelk-healed-wait.mss` | Fire Beam | 1600 → 1246 (2 hits) | **stopped by operator directive** |

Neither attempt defeated Whelk. Attempt 2 was **frozen deliberately**,
not lost.

## Observations

### Field healing before contact (Confirmed)

The field menu is reachable at `(2A,08)`/`(2A,09)` before Whelk contact.
Inventory held **Tonic ×4** (description: recovers 50 HP) and
**Potion ×1**. Four Tonics took the party from the EXP-0039 carry-in
state **26 / 19 / 56** to **76/77, 105/105, 106/107**, verified against
`WRAM:+$1609` (`0x4C` = 76) and the on-screen party window. The Potion
was retained unused.

Menu-derived party record: Terra (`?????`) LV 4, HP 76/77, MP 24/29;
Wedge LV 4, HP 105/105; Vicks LV 4, HP 106/107; Steps 140, Gp 3720.

### Whelk is two battle entities (Confirmed)

Reading the DISC-0001 arrays for formation 432:

| Slot | Current HP | Max HP | MP candidate | Reading |
|---|---|---|---|---|
| 0-2 | party | party | — | Terra / Wedge / Vicks |
| 4 | 50000 | 50000 | 120/120 | shell |
| 5 | 1600 | 1600 | 1000/1000 | head |

The staged `+$3F44` record reproduced EXP-0039's bytes exactly
(`80 03 00 34 FF FF FF FF 48 AB 00 00 00 00 3F`), a third independent
confirmation of the formation staging.

### Head-only targeting damages the head and never the shell (Confirmed)

Every damaging action was issued only on a frame in which the head was
visibly extended, and each was verified afterwards by HP delta:

| Attempt | Ability | Head before → after | Damage | Shell |
|---|---|---|---|---|
| 1 | Bolt Beam | 1600 → 1423 | 177 | 50000 (unchanged) |
| 1 | Bolt Beam | 1423 → 1242 | 181 | unchanged |
| 1 | Bolt Beam | 1242 → 1071 | 171 | unchanged |
| 1 | (queued) | 1071 → 909 | 162 | unchanged |
| 2 | Fire Beam | 1600 → 1432 | 168 | unchanged |
| 2 | Fire Beam | 1432 → 1246 | 186 | unchanged |

Slot 4 held 50000 for the entire unit. **Falsifier 3 is not triggered**:
striking the extended head does reduce the head's HP. **No shell strike
occurred in either attempt**, so EXP-0039's counter was never
re-triggered and nothing new about it is claimed here.

### Head/shell state is visually classifiable (Confirmed)

**Falsifier 2 is not triggered.** At 4× upscale the two states separate
unambiguously: extended shows a golden serpentine head with a green eye
and mandibles clear of the shell; retracted shows only the purple spiked
shell with a dark cavity where the head was. One early frame was
initially ambiguous because the white **target cursor** sits in the same
screen region; upscaling resolved it. Classification was corroborated by
HP deltas on every attack.

### Menu, command and targeting behavior (Confirmed)

- **MagiTek ability sets are character-specific, and EXP-0039's list was
  incomplete.** Terra shows **eight**: Fire Beam, Bolt Beam, Ice Beam,
  Bio Blast, Heal Force, Confuser, X-fer, TekMissile. Wedge and Vicks
  show **four**: Fire Beam, Bolt Beam, Ice Beam, Heal Force. EXP-0039
  recorded four; that was a non-leader's list. Updates CEN-MAGIC-0001.
- Command windows match EXP-0039: Terra `MagiTek / Magic / Item`;
  Wedge and Vicks `MagiTek / Item`.
- **The ally target cursor defaults to the caster's own slot** (observed
  for all three characters), not to slot 0. Several early heals landed on
  the wrong ally before this was established.
- **Cursor = Memory** reopens a character's ability list on their
  last-used ability (observed for Terra, twice).

### Event and content observations (contextual)

- The **"GUARD: We won't hand over the Esper!!"** beat fires at
  `(2A,07)` on a clean, never-defeated run, **before** the Whelk battle,
  and reproduced on both attempts. This resolves EXP-0039's open
  alternative: it is normal pre-Whelk progression, **not** a
  post-defeat artifact. Updates CEN-EVENT-0011.
- The intro dialogue states Whelk **eats lightning** and **stores the
  energy in its shell**, and ends with an explicit **"don't attack the
  shell!"**. Registered as observed on-screen guidance; **no elemental
  absorption was tested or measured** by this unit.
- Enemy action name **"Slime"** was displayed. Whelk actions reduced
  multiple party members' HP within one action window (party-wide
  effect), magnitudes ranging roughly 10-26 per member. Recorded as
  contextual only — no AI, targeting or damage model is claimed.

### Frozen state at stop (attempt 2)

- Frame **182248**. Formation `+$11E0` = `B0 01` (432).
- Party HP: Terra **51/77**, Wedge **105/105**, Vicks **107/107**.
- Terra MP 24/29 (`+$3C08`/`+$3C30`); other party MP read 0.
- Head **1246/1600**; shell **50000/50000**.
- Head state: **retracted**.
- Menu depth: Terra's **MagiTek ability submenu open**, cursor on
  Fire Beam, no target selection open.
- Field position bytes read `00 10` — not field-meaningful during
  battle, consistent with the established note.
- Savestate `attempt2-frozen-f182248.mss`, lineage:
  `pre-whelk-recon.mss` → field healing → `pre-whelk-healed.mss` →
  Config (`Cursor` → `Memory`) → `pre-whelk-healed-wait.mss` →
  guard beat → Whelk battle attempt 2 → frozen at frame 182248.

## Interpretation

The two pre-registered mechanical premises of branch A **hold**: the
head is a distinct, damageable battle entity with 1600 HP, it is
visually distinguishable from the shell, and attacking it while extended
reduces its HP without touching the shell. In that narrow sense the
strategy is sound and the unit's falsifiers 2 and 3 were both tested and
not triggered.

What the unit could **not** do is operate the battle efficiently. Action
selection repeatedly desynchronized from the game's state: the head
changed between opening a submenu and reaching target selection, queued
actions resolved out of the order they were issued, heals landed on
unintended allies, and on one occasion the target cursor was sitting on
the shell at the moment of confirmation and had to be cancelled. Those
failures are **operator/tooling failures against an unmodelled timing
system**, not evidence about Whelk.

Critically, the project has no model of FF6's ATB: which clocks
ACTIVE/WAIT actually gate, which submenu states qualify as pausing,
how the action queue orders and resolves entries, or what drives the
head/shell cycle. Every head/shell transition observed here occurred
while the emulator was being driven through menus with operator-length
pauses between inputs, so **none of this unit's timing data can be used
to characterize Whelk's natural head/shell durations**.

## Result

**No victory. Milestone `10-whelk-victory` is NOT established, and B19
remains uncaptured.** The unit is stopped by operator directive at a
deliberately frozen, fully documented state.

Classification of this run:

- **Aborted due to missing ATB model** (primary)
- **Partial contextual capture** (substantial confirmed side-evidence)
- **Failed tactical attempt** (two attempts, neither reached victory)

The unit did establish real, reusable facts: the head/shell entity
split with exact HP, head-only damage confirmation with six measured
deltas, the field healing route, the corrected per-character MagiTek
sets, caster-default targeting, and the guard beat's true sequence
position.

### Claims explicitly NOT made by this unit

This run does **not** support, and no downstream record may assert:

- that Whelk uses any particular frame timer;
- that WAIT pauses every battle clock;
- that a visible submenu necessarily activates WAIT behavior;
- that observed enemy inactivity proves the enemy ATB is frozen;
- that holding the MagiTek menu open demonstrates anything about
  Whelk's timer;
- that the healing behavior chosen here was optimal;
- that the fight is impossible, or unusually long, under correct ATB
  operation.

## Confidence

- Head/shell entity split, HP values, and head-only damage behavior:
  **Confirmed** (six measured deltas, shell invariant throughout).
- Head/shell visual classifiability: **Confirmed**.
- Field healing route and item inventory: **Confirmed**.
- Per-character MagiTek sets and caster-default ally targeting:
  **Confirmed**.
- Guard beat precedes Whelk on a clean run: **Confirmed** (reproduced
  twice, never-defeated run).
- Cursor Memory behavior: **Confirmed** for the observed case (n=2).
- Whelk action semantics, elemental behavior, AI, and all head/shell
  **timing**: **Unknown** — deliberately not investigated, and the
  timing observations here are contaminated by menu-induced pauses.

## Blocker

**Further Whelk execution is deferred until the project establishes a
usable ATB model**, including ACTIVE/WAIT behavior, which submenu states
qualify as pausing, the relevant timer domains, and action-queue
behavior. Whelk gameplay must not resume before that research.

## Next action

Start a new session, resume from this checkpoint, audit existing battle
infrastructure and ATB evidence, and propose the first bounded ATB
experiment before operating Mesen.
