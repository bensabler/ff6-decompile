# EXP-0039: Mines-to-Whelk breadth reconnaissance

- **Status:** planned (2026-08-01)
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

(during the pass)

## Observations

(during the pass)

## Interpretation

(after the pass)

## Alternatives

(after the pass)

## Result

(after the pass)

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

(after the pass)
