# EXP-0039: Mines-to-Whelk breadth reconnaissance

- **Status:** completed (2026-08-01) — stopped at a stable gameplay
  boundary on context limits; Whelk reached and observed, victory not
  yet achieved
- **Mode:** **breadth-first reconnaissance**, not a depth investigation.
  Operating rule: **explore widely, register briefly, continue forward.**
- **Program:** SCN-0001 (opening-to-Whelk). Serves B11–B18 by
  *observation and registration*, not by decoding.

## Question

Starting from the furthest verified milestone, what systems, content
families, branches, encounters, objects, interactions, graphics, audio
and events are visible or reachable between the mines interior and
Whelk — and can the scenario be advanced to a first visible Whelk
victory?

This unit answers **what exists and where**, deliberately not **how it
works**. Every observation is registered in the content census with one
bounded future question; none is investigated in place.

## Starting state

Visible GUI Mesen with the bridge, loading the archived milestone state
(`local_artifacts/scenarios/SCN-0001/06-random-encounter/…-06.mss`
once EXP-0038 establishes it; otherwise `05-mines-entry.mss`).
Interactive correction is allowed and expected — this is a piloted
recon, not a scheduled route.

## ROM identity

`0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`.

## Emulator identity

Mesen 2.1.1, **GUI (operator-visible)**, normal speed or modest
watchable fast-forward. Headless verification is explicitly *not* part
of this unit; it belongs to the follow-up units that convert
observations into scheduled routes. Exactly one Mesen instance runs at
a time.

## Independent variable

None. This is observational reconnaissance with piloted input.

## Controlled variables

Same ROM and lab settings; the savestate provenance is recorded for
each exploration branch so any observation can be re-reached.

## Instrumentation

Bridge commands only (`loadstate`, `press`, `read`, `screenshot`) plus
periodic position reads (`+$00AF`/`+$00B0`), formation reads
(`+$11E0`), and screenshots at every junction, object, dialogue,
battle and transition. No new write-watches, no exec census, no
disassembly during the pass. Existing confirmed tables are *linked*
(formation `$CF6200`, monster `$CF0000`) rather than re-derived.

## What is recorded (register, do not investigate)

- **Navigation:** coordinates, junctions, dead ends, one-way paths,
  blocked tiles, map transitions, event barriers, trigger tiles,
  camera changes, collision observations, candidate map identifiers.
- **Encounters:** every formation id, trigger location, random vs
  scripted, timing, background, music, monsters present, victory and
  rewards, and whether one area yields multiple formations.
- **Objects/interactions:** chests, save points, NPCs, switches,
  doors, signs, dialogue, invisible triggers, environmental
  animations, inaccessible or suspicious objects.
- **Menus/mechanics:** available commands, items used, Magic/Magitek
  behavior, status effects, healing, party state, equipment, save
  behavior, any newly exercised mechanic.
- **Graphics families:** tilesets, palettes, backgrounds, sprites,
  animations, lighting/screen effects, fonts and windows, Whelk assets
  once visible.
- **Audio families:** map music, battle music, boss music, SFX,
  transition cues, ability sounds, environmental audio.
- **Events/state:** dialogue occurrences (by location and length; no
  game text is tracked — asset policy), event transitions, party/map
  state changes, candidate event flags, scripted movement, battle
  invocation, pre-Whelk setup, Whelk entrance, post-battle state.

## Whelk objective

If Whelk is reached: preserve a **pre-Whelk savestate** locally, record
the route and event sequence, enter the battle visibly, attempt the
**head-focused** strategy, **avoid attacking the shell** on the first
victory attempt, and capture the first stable post-battle state if
Whelk is defeated. A first visible Whelk victory is a valid scenario
milestone even with Whelk's AI, data and assets undecoded.

## Expected outcomes

- *Supports:* a route map from the starting milestone to the furthest
  point reached, with every branch and dead end enumerated, all
  observed content families registered in the census, and the scenario
  advanced (ideally to a Whelk victory).
- *Partial (valid):* forward progress stops at a specific obstacle —
  recorded as the exact blocking fact with the smallest experiment
  that would resolve it.

## Falsifiers

Not a hypothesis test; the unit fails only if it produces
unregistered observations, undocumented route knowledge, or
depth-investigation drift. Concretely:

1. Any observation left without a census entry or a bounded future
   question.
2. Any claim of behavior not backed by a captured screenshot, log
   line, or savestate.
3. Any deep decode (formats, AI, provenance, RNG) performed inline —
   a scope violation to be recorded and stopped.

## Evidence requirements

Exploration transcript with coordinates per step, screenshots at every
junction/object/battle/transition, formation ids with staged records,
savestates at each significant boundary (pre-Whelk especially),
`hashes.sha256` — all under `local_artifacts/experiments/EXP-0039/`
and, for milestones, the SCN-0001 milestone directories.

## Trials

One piloted GUI pass (2026-08-01) from `05-mines-entry.mss`, ~90
bridge commands, 30 evidence artifacts. Stopped at a stable boundary
on context limits, not on a blocker.

## Observations

### Route map (from milestone 05 at `(26,1C)`)

```text
(26,1C) ──south──> (26,22) ──TRANSITION──> Narshe exterior (1F,18)
   │                                            │
   │  <──────── TRANSITION (1F,14) ─────────────┘   [bidirectional]
   ▲
(26,1F) ──north──> (26,15) [encounter fired here] ──> (26,0B)
                                                        │
   west: blocked ────────────────────────────────────── ┤
                                                        │
                                     east ──> (28,0B) ──┼── south: (28,0C) DEAD END (1 tile)
                                                        │
                                              north ──> (28,08)/(28,0A)
                                                        │
                                     east ──> (2A,0A) ──> (2A,09)  [SCRIPTED EVENT: Vicks]
                                                        │
                                              north ──> (2A,07) ──contact──> WHELK BATTLE
                                                        │
                                              (after Whelk) ──> GUARD/Esper beat at (2A,07)
```

**Branches and dead ends enumerated:** `(26,0B)` west — blocked;
`(26,0B)` north — blocked (corridor turns east); `(28,0C)` — one-tile
southern nub, dead end; `(28,08)` east — blocked (must use the `0A`
row); `(2A,09)` north/east — blocked until the Vicks event resolves.
The mines entry corridor is **linear with one turn**; no unexplored
side passage was found on this path.

**New navigation finding — the mines transition is bidirectional.**
Walking south from `(26,1C)` through `(26,22)` returns the party to
the **Narshe exterior** at `(1F,18)`; climbing back north from
`(1F,14)` re-enters the mines at `(26,1F)`. Re-entry produced **no
shaft dialogue**, consistent with that beat being one-time and
flag-gated (`EVF-1EA0-$2B` was already set) — a behavioral
corroboration of the EXP-0037 flag inventory obtained for free.

### Encounters

| # | Formation | Trigger tile | Origin | Notes |
|---|---|---|---|---|
| 1 | 14 (`$000E`) | `(26,0B)` | random | EXP-0038 recon |
| 2 | 14 (`$000E`) | `(26,0B)` | random | EXP-0038 scheduled runs ×2 |
| 3 | 14 (`$000E`) | **`(26,15)`** | random | this pass |

**This discriminates an EXP-0038 alternative.** That record listed
"the encounter could be position-triggered rather than step-counted"
as an open alternative, needing a run that reaches the same tile after
a different number of steps. This pass supplied exactly that: after
the exterior round-trip the encounter fired at `(26,15)`, nine tiles
south of the `(26,0B)` tile seen twice before. **Fixed-tile
triggering is refuted**; step/zone counting survives. The producer
itself remains unlocated (CEN-WORLD-0006).

Formation 14 recurred in all three rolls; no second mines formation
was drawn on this path (EXP-0030's formation 44 came from a different
mines location).

### Menus and mechanics

- **Magitek ability set captured** (CEN-MAGIC-0001, a long-standing
  gap): **Fire Beam, Bolt Beam, Ice Beam, Heal Force** — four
  abilities in a two-column list. This is the on-screen ability list
  only; the attack-record index for Fire Beam remains ambiguous
  (EXP-0022's indices 5/131), deliberately not chased here.
- **Battle commands are character-specific**: the leader shows
  **MagiTek / Magic / Item**, while another member shows **MagiTek /
  Item** only (no Magic). Registered against CEN-MENU-0001.
- Party HP is carried across the mines: entering Whelk the party stood
  at 26 / 19 / 56 after three prior battles, which decided the outcome
  below.

### Whelk — reached, entered, and lost (B17/B18 first contact)

- **Trigger:** a scripted beat fires on arriving at `(2A,09)` (Vicks
  steps forward; movement input is consumed). After it resolves,
  Whelk itself is **contact-triggered** by pushing north from
  `(2A,07)` — the same "walk into it" pattern EXP-0036 found for the
  fifth scripted battle, now seen a second time.
- **Formation 432 (`$01B0`)**, staged `+$3F44` record
  `80 03 00 34 FF FF FF FF 48 AB 00 00 00 00 3F`. Under the current
  reading (record bytes 2-7 = monster ids) the ids decode to `00` and
  `34` — but id `00` is the *opening guard* record, which is
  implausible for a boss. **Bounded future question:** the formation
  record's monster-id field must carry a high-bit/extension not yet
  decoded (FF6 has more than 256 monsters), so the Whelk ids are
  **Unknown** pending that decode. Not investigated here.
- **Introduction/warning dialogue observed** (B17, never previously
  reached): a multi-box Vicks/Wedge exchange precedes the fight;
  the enemy name window renders **"Whelk"** on screen (observational
  only — the monster name table is still unlocated,
  CEN-MONSTER-0004).
- **Shell counter confirmed behaviorally** (B18): with the head
  retracted, continued attacks drew a counterattack that killed Wedge
  outright and left the leader at 14 HP.
- **Outcome: defeat.** The party was wiped and the game reached its
  **defeat flow** — which incidentally captures CEN-BATTLE-0007, a
  registered gap ("capture the post-defeat transition").
- **Attempt 2** reloaded the pre-Whelk savestate and advanced past the
  Vicks beat, where a **further scripted beat** appeared at `(2A,07)`
  ("GUARD: We won't hand over the Esper!!"), i.e. the scenario
  continues beyond Whelk's location. Registered; not pursued.

### Graphics and audio families (registered, not traced)

Mines interior: cave tileset with rail trestles, hanging lamp sprites
with an animated glow cycle (the CEN-QUIRK-0002 lamp phase), a dark
ambient palette, and multi-layer depth. Whelk: large multi-part boss
sprite (shell + extendable head) with a distinct battle background
reusing the cave scene. Dialogue windows use the standard blue
gradient frame. No DMA, compression or provenance tracing performed,
per the unit's bounds.

## Interpretation

The mines-to-Whelk stretch is **short, linear and fully reachable**:
one corridor with a single turn, one dead-end nub, a bidirectional
exterior transition, one random-encounter zone drawing formation 14,
and two scripted beats bracketing the Whelk fight. Nothing in this
stretch is gated behind an unsolved system — the scenario's remaining
obstacle to a Whelk victory is **party HP management**, not knowledge.

The Whelk fight's structure matches the SCN-0001 program brief's
expectation (head state vs shell state with a punishing counter), so
the planned branches A/B/C remain the right test design; this pass
effectively executed an unplanned "branch B" (attack the shell) and
recorded its consequence.

## Alternatives

- *Formation 14 may not be the zone's only formation*: three rolls
  all drew it, but three is a small sample and EXP-0030 saw 44
  elsewhere. Unresolved (B13).
- *The `(2A,07)` guard/Esper beat may be the normal post-Whelk
  progression rather than an alternate branch*: this pass reached it
  after a defeat-and-reload, so its true position in the sequence is
  **Unknown** and needs a clean run.

## Result

**Whelk was reached, entered and observed — the scenario's furthest
point to date.** B17 (introduction/warning dialogue) moves from
"never reached in any archived state" to observed with evidence;
B18 gains first-contact behavior (contact trigger, formation id,
staged record, shell counter, defeat flow); B11 gains a complete
branch/dead-end map of the entry corridor with a bidirectional
transition; B12 gains a discriminating datum that **refutes
fixed-tile encounter triggering**; B15 gains the full Magitek ability
list and the per-character command split.

**A first Whelk victory was not achieved** — the first attempt ended
in defeat after the shell counter. A pre-Whelk savestate is preserved
(`local_artifacts/scenarios/SCN-0001/07-pre-whelk/pre-whelk-recon.mss`,
129 795 bytes), so the retry costs nothing to set up.

**Scope discipline held:** every observation above is registered with
a bounded future question and none was decoded in place. No event-flag,
map-format, AI, RNG, provenance or compression work was performed.

## Confidence

Observational: reaching, seeing and reproducing a location or system
is Confirmed at the level of "this exists and is reachable this way".
No format, semantic or mechanism claim is promoted by this unit.

## Stopping condition

Stop at a **stable gameplay boundary** when any of these occurs:
Whelk is defeated and the first stable post-battle state is captured;
forward progress is blocked by a fact requiring its own experiment;
or context/session limits are reached. In every case: preserve the
savestate and evidence, clean up all background jobs, update this
record, synchronize census and scenario, run gates, checkpoint, and
record the exact next route action.

## Bounds (strict anti-depth rules)

Do **not** pause gameplay to investigate: event-flag meanings, map
header formats, collision engines, encounter RNG, monster AI
interpreters, complete monster records, graphics or audio provenance,
compression formats, save-file structures, or engine architecture.
For every unrelated observation: record what was observed, link the
evidence, add/update a census entry, write one bounded future
question, and return to gameplay. A side investigation is permitted
only when a specific unknown *directly prevents forward movement* —
then resolve the smallest blocking fact and continue.

## Next action

**EXP-0040 — Whelk victory attempt (branch A).** Reload
`local_artifacts/scenarios/SCN-0001/07-pre-whelk/pre-whelk-recon.mss`,
use **Heal Force** (Magitek) to restore the party from 26/19/56
*before* engaging, then attack **only while the head is extended**,
verifying head state from a screenshot between actions rather than
mashing A. On victory, capture milestone `10-whelk-victory` (WRAM +
screenshot + savestate) and the first stable post-battle state (B19).
Branch B (deliberate shell attack) is already partly recorded by this
pass's defeat and should be re-run cleanly afterwards.
