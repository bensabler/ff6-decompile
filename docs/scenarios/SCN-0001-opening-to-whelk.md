# SCN-0001: Opening scenario — New Game through the Whelk battle

- **Status:** ACTIVE (program in progress)
- **Created:** 2026-07-31
- **Machine manifest:** `data/scenarios/opening-to-whelk.json`
- **ROM revision:** SHA-256 `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2`

This is the master specification for the project's first vertical
slice: an evidence-backed reconstruction of everything the game
exercises from New Game selection through the defeat of Whelk. It is a
project-authored structural specification. It contains no game
dialogue text; dialogue is tracked by ID, location, length, hash, and
short project-authored descriptions only (see the asset policy).

## Boundary

```text
Start:
New Game is selected from a fresh boot or reproducible equivalent state.

End:
Whelk has been defeated, battle rewards and victory processing have
completed, and the game has reached the first stable post-battle state
before the frozen Esper interaction begins.
```

The frozen-Esper reaction and Terra's awakening are OUT of scope,
except for the few frames needed to establish that the Whelk scenario
has ended.

## Program shape

This is a multi-unit reconstruction program, not one experiment. Each
bounded unit gets its own `docs/experiments/EXP-NNNN` record, its own
census pass, and its own commit. This record is the integration point:
every beat below links to the experiments, census entries, ROM
regions, implementations, and tests that cover it. Update this record
and the JSON manifest whenever a unit closes.

### Golden route

A reproducible, frame-scheduled Mesen route from power-on to Whelk
victory, producing an input transcript (tracked), milestone frame
numbers (tracked), state assertions (tracked), and local savestates
plus raw captures (ignored paths only). Milestone directories under
`local_artifacts/scenarios/SCN-0001/`:

```text
00-new-game/            01-opening-cinematic/   02-narshe-entry/
03-first-scripted-battle/ 04-free-movement/     05-mines-entry/
06-random-encounter/    07-pre-whelk/           08-whelk-head-test/
09-whelk-shell-counter-test/ 10-whelk-victory/
```

### Whelk branches

Behavior coverage must not depend on one playthrough. Minimum
controlled branches from the `07-pre-whelk` state:

- **Branch A** — fight correctly: attack the exposed head to victory.
- **Branch B** — intentionally attack the shell; capture the complete
  response (counter selection, targeting, damage).
- **Branch C** — observe at least one full head → shell → head cycle
  without ending the battle.

Additional branches are added when runtime behavior shows they are
needed (record them here).

## Pre-existing evidence baseline (2026-07-31)

Four archived mid-scenario savestates exist (`mesen/out/`):
`checkpoint1.mss` (scripted dialogue beat, Narshe exterior),
`checkpoint2.mss` (free walk, Narshe exterior), `checkpoint3-mines.mss`
(mines interior), `exp10-battle.mss` (scripted opening battle, party
`?????`/WEDGE/VICKS). None starts at power-on; their input provenance
is not frame-scheduled. The golden route supersedes them as scenario
evidence; they remain useful as working states.

Confirmed systems already covering parts of this scenario:

| System | Evidence | Status |
|---|---|---|
| Battle damage pipeline (both base formulas, defense scaling, element response, display chain) | EXP-0001..0019, DISC-0001..0007 | Implemented + regression-tested |
| Action determinism under frame-exact input | EXP-0021 | Confirmed for tested window |
| Attack-record table | EXP-0017/0018/0022 | Confirmed, cross-checked |
| Battle HUD fixed font (ROMFILE 0x046FC0) | EXP-0023, GFX-0001 | Confirmed, decoder implemented |
| Battle-confirm SFX full chain (CPU→APU→DSP→BRR) | EXP-0024, AUD-0001 | Confirmed, decoder implemented |
| Menu/dialogue text encoding | EXP-0026 | Confirmed |
| Spell database $C46AC0 (54 records), name table 0x26F567, MP cost byte | EXP-0026/0027 | Confirmed, extracted |
| Spell availability array WRAM:+$1A6E | EXP-0026 | Strong hypothesis |
| Monster database $CF0000 (32-byte records) | EXP-0028 | Base Confirmed, stride Strong |
| Per-slot monster loader $C22C30 (+aux tables ×4/×8/×32/$CF8400) | EXP-0029 | Confirmed |
| Formation table $CF6200 (15-byte), staging at +$3F46, flags $CF5900 | EXP-0029/0030 | Confirmed |
| Field character-record block ~WRAM:+$1600 (HP +$09, MP +$0D offsets observed) | EXP-0027 | Strong hypothesis |
| Event-flag arrays $1E80/$1EA0/$1EC0: decoder, masks, 16-handler family, opening write inventory (20 flags), interpreter anchor $C09B5C | CONTRA-0002, EXP-0037, DISC-0008 | Confirmed (meanings unassigned) |

## Timeline beats

Ordered beats B01–B19. The division follows the required scenario
timeline; runtime evidence corrects it if the game divides the
sequence differently (corrections are recorded per beat, never
silently).

Matrix rows per beat (values `COMPLETE` / `PARTIAL` / `BLOCKED` /
`NOT_APPLICABLE`; every `PARTIAL`/`BLOCKED` is explained — a `PARTIAL`
with "nothing yet" means the row is open with zero evidence):

```text
identified · runtime · event-script · data-records · behavior ·
graphics · audio · persistence · go-impl · tests · differential
```

### B01 — New Game initialization
Power-on → title → New Game selected; WRAM/character/inventory/RNG
initialization.
- identified/runtime PARTIAL→advanced: **EXP-0031 established
  milestone `00-new-game`** — deterministic frame-scheduled route
  from power-on (title at ~2969; start+a edge toggling 2500–4200;
  milestone at 5200; two runs byte-identical at frames 5200 and
  30000). Title/attract flow registered (CEN-EVENT-0004); boot
  uninitialized-read entropy candidates registered (CEN-QUIRK-0001).
- data-records/behavior PARTIAL: the initialization writes themselves
  are uncaptured (the dedicated init unit; baseline WRAM dump
  archived). Links: CEN-SAVE-0001, EXP-0031.

### B02 — Opening cinematic / introductory presentation
Presentation between New Game and the snowfield march.
- runtime PARTIAL→advanced: the real opening **auto-runs without
  input** from the title press to the first input-waiting dialogue at
  scripted Narshe entry (EXP-0031; the attract variant runs the same
  scenes dimmed). Milestone-grade captures for the interior beats
  (`01-opening-cinematic`) still pending in segment 2.
- Remaining rows PARTIAL (event script, graphics, audio provenance
  all open).

### B03 — Snowfield march toward Narshe
Scripted Magitek walk. **All rows PARTIAL: no capture.** Event engine
unlocated (CEN-EVENT-0001).

### B04 — Character and party initialization
`?????` (pre-naming), Vicks, Wedge, Magitek armor state.
- identified/runtime PARTIAL: `?????`/WEDGE/VICKS observed on-screen
  (EXP-0025); battle-slot staging writers known (EXP-0018/0028).
- data-records PARTIAL: field block ~+$1600 located (EXP-0027); the
  new-game initializer and name store unlocated (CEN-EVENT-0003,
  CEN-CHAR-0002/0004).
- Remaining rows PARTIAL (nothing yet).

### B05 — Opening dialogue and scripted movement
- runtime PARTIAL: checkpoint1 shows a dialogue beat mid-sequence
  (EXP-0025); no full capture from B02.
- event-script PARTIAL: dispatcher unlocated (CEN-EVENT-0001);
  dialogue encoding Confirmed (EXP-0026), text/dialogue table
  unlocated (CEN-EVENT-0002).
- Remaining rows PARTIAL (nothing yet).

### B06 — Entry into Narshe
Map transition into the town proper.
- identified/runtime PARTIAL→advanced: milestone `02-narshe-entry`
  captured and deterministic (EXP-0032).
- Remaining rows PARTIAL: the transition's map records, ids, and
  loading path are still unlocated (CEN-WORLD-0004).

### B07 — Guard encounters / scripted battles
- runtime PARTIAL→advanced: **battle entry now captured on the golden
  route** (EXP-0032, milestone `03-first-scripted-battle`, detected
  identically at frame 31 557 in both runs). The battle interior is
  the most-tested state in the project (EXP-0001..0022 window);
  on-screen formation names Were-Rat / Repo Man registered
  (EXP-0027, defeat flow).
- data-records COMPLETE for identity: **formation id 2 = two monsters
  of record id 0** (staged `+$3F44` bytes match ROM record 2 at
  `ROMFILE:0x0F621E`; live battle slots 6/7 carry HP 40 / MP 15,
  matching monster record 0's `+$08`/`+$0A`). This corrects the first
  draft of EXP-0032, which misread record byte 1 as a monster id and
  claimed a single monster 12 — the correction is recorded in
  EXP-0032 and EXP-0033, not silently patched. Second independent
  verification of the EXP-0030 formation table, on a scripted
  encounter.
- **The opening runs exactly four scripted battles** before player
  control (EXP-0034, identical across two power-on runs):
  formations **2, 1, 2, 41** at entry frames 31 557 / 34 953 /
  36 828 / 39 500. Every staged `+$3F44` record matches the ROM
  formation table byte-for-byte — the table is now verified across
  five independent encounters.
- behavior PARTIAL: damage pipeline Confirmed+implemented; battle 1
  is **won on schedule** (frame 32 706) with rewards **32 EXP /
  96 GP** and a battle→field writeback into `~+$1600`
  (`$C2496E`/`$C24979`; EXP-0033, CEN-BATTLE-0008). The scripted-
  battle **invocation opcode** remains unknown (`$0206`/`$3A97` lead,
  EXP-0029; CEN-EVENT-0005).
- go-impl/tests PARTIAL: damage pipeline implemented and regression-
  tested; differential vs live for the standard path (EXP-0011..0016).

### B08 — Controllable field movement
- identified/runtime PARTIAL→advanced: **milestone
  `04-free-movement` established at frame 46 375 and byte-identical
  across two power-on runs** (EXP-0034) — the Narshe exterior under
  player control. Free walk and live collision were already verified
  in an archived state (EXP-0025).
- Remaining rows PARTIAL: per-step/collision logic unlocated
  (CEN-WORLD-0003).

### B09 — Narshe exterior traversal
- runtime PARTIAL→advanced: the route from milestone 04 to the mine
  entrance is **mapped leg by leg with player-tile coordinates**
  (EXP-0035). The climb is a zigzag, and it is gated by a **fifth
  scripted battle** at tile ($1E,$27) plus a shaft dialogue at
  ($1F,$16) — CEN-EVENT-0007.
- Remaining rows PARTIAL: tileset/animated tiles observed only
  (CEN-WORLD-0001); the accessible-branch and dead-end sweep is still
  not performed, and static map/object data has not been inspected.

### B10 — Mine entry and map transition
- identified/runtime **COMPLETE**: milestone **`05-mines-entry`
  established** (EXP-0036) — three power-on runs of the 17-leg
  state-driven route reach (`$26`,`$1C`) inside the mines with
  byte-identical WRAM (`c26453d3…`) and matching leg frames.
- behavior PARTIAL: the transition is observable as a player-position
  jump; `WRAM:+$1EA5` moves to `$0D` **before** the transition is
  visible — and CONTRA-0002 resolved why: the byte is part of the
  event-flag bit array at `+$1EA0` (CEN-EVENT-0008), so it never
  denoted a map at all. EXP-0037 pinned the moment exactly: flag
  `EVF-1EA0-$2B` is set at frame 50 699, during leg 15's
  shaft-dialogue settle at (`$1F`,`$16`) — 214 frames before the
  position jump completes leg 16 at 50 913, confirming EXP-0036's
  "set during the preceding dialogue leg" reading.
- Remaining rows PARTIAL: the map records, header, and tileset load
  path behind the transition are still unlocated (CEN-WORLD-0004).

### B11 — Mine traversal
**All rows PARTIAL**, but the entry corridor is now mapped and walked
on a schedule (EXP-0038): `(26,1C)` north to `(26,0B)` — up blocked
there — east across a rail trestle to `(28,0B)`, north to `(28,09)`,
and east toward `(2A,09)`, where a **scripted event** fires (dialogue,
the party splits; registered as CEN-EVENT-0009 and left uninvestigated).
Cave tileset/rails/light still observed only (EXP-0025,
CEN-WORLD-0002); the full branch sweep and collision data remain open
and are the subject of the EXP-0039 breadth pass.

### B12 — Random encounter check and trigger
- runtime PARTIAL→advanced: an encounter is now **reproduced from
  power-on on the scheduled route** (EXP-0038, frame 51 307,
  formation 14, interrupting leg 19 near tile `(26,0B)`); the roll
  chain downstream of `WRAM:+$11E0` was already Confirmed
  (EXP-0029/0030).
- behavior PARTIAL: the `+$11E0` producer (zone data + step state)
  remains the open frontier (CEN-WORLD-0006). Two independent
  triggers — piloted recon and the scheduled run — landed in the same
  one-tile neighbourhood after comparable step counts, which is
  consistent with a step-counter check but is recorded as observation
  only.
- persistence: **encounters set no event flags** — EXP-0038's flag
  timeline is identical to EXP-0037's opening total, with no writes
  after frame 50 880.
- Remaining rows PARTIAL.

### B13 — Encounter packs in accessible opening zones
- data-records PARTIAL: formation table Confirmed ($CF6200); the
  zone→pack→formation selection structures unlocated. Mines pack
  membership now known for **two** rolls: formation 44 (EXP-0030) and
  **formation 14** (EXP-0038 — three of monster record 19; staged
  record byte-identical to `ROMFILE:0x0F62D2`, the table's seventh
  independent verification). One area therefore yields multiple
  formations.
- Remaining rows PARTIAL (nothing yet).

### B14 — Monsters reachable before Whelk
- data-records PARTIAL: monster db Confirmed. Known so far — records
  **0** (40 HP / 15 MP, live-verified in battle slots 6/7) and **25**
  (27 HP / 5 MP) from the four scripted battles, plus **19** and
  **77** from the mines random encounter (EXP-0030). The full
  reachable set still depends on B13 (encounter zones). Monster name
  table unlocated; AI scripts unlocated; rewards fields unmapped;
  30 of each record's 32 bytes unmapped (CEN-MONSTER-0004).
- Remaining rows PARTIAL (nothing yet).

### B15 — Accessible menus and commands
- runtime PARTIAL: field main menu + submenus censused (EXP-0026);
  battle command window observed (EXP-0025); Magitek ability list
  never captured (CEN-MAGIC-0001; Fire Beam record ambiguous,
  EXP-0022).
- data-records PARTIAL: spell db Confirmed; command records unlocated.
- Remaining rows PARTIAL (nothing yet).

### B16 — Chests, save points, NPCs, objects, branches
**All rows PARTIAL** (CEN-WORLD-0005 treasure lead). Requires the
full accessible-branch sweep plus static object data. The flag
machinery these interactions will use is now decoded and inventoried
(EXP-0037, DISC-0008, CEN-EVENT-0008); the golden route performs no
chest/NPC interaction, and correspondingly no additional
`$1E80/$1EA0/$1EC0` flag is touched during free movement outside the
scripted beats.

### B17 — Whelk introduction
- identified/runtime PARTIAL→**REACHED** (EXP-0039, first time in the
  project): a scripted beat fires on arriving at mines tile
  `(2A,09)` (a party member steps forward; movement input is
  consumed), after which Whelk is **contact-triggered** by pushing
  north from `(2A,07)` — the same walk-into-it pattern EXP-0036 found
  for the fifth scripted battle. A multi-box introduction/warning
  dialogue precedes the fight, and the enemy name window renders
  "Whelk" on screen (observational only; the monster name table is
  still unlocated). Registered as CEN-EVENT-0010.
- Remaining rows PARTIAL: event script, graphics and audio provenance
  all untraced (breadth-pass bounds).

### B18 — Whelk battle
Head state, shell state, retreat/reappearance, warning behavior, head
attacks, deliberate shell attacks, counterattack, AI timing, victory,
rewards.
- identified/runtime PARTIAL→**FIRST CONTACT** (EXP-0039): the battle
  was entered and observed. **Formation 432 (`$01B0`)**, staged
  `+$3F44` record `80 03 00 34 FF FF FF FF 48 AB 00 00 00 00 3F`.
- behavior PARTIAL: the **shell counterattack is confirmed
  behaviorally** — with the head retracted, continued attacks drew a
  counter that killed one party member outright and then wiped the
  party. The first attempt therefore **ended in defeat**, which
  incidentally captured the defeat flow (CEN-BATTLE-0007).
- data-records PARTIAL: **Whelk's monster ids are Unknown.** Reading
  formation record bytes 2-7 as ids gives `00` — the opening guard
  record — which is implausible for a boss, so the id field must carry
  a high-bit/extension not yet decoded (FF6 exceeds 256 monsters).
  This is a new bounded question, not a decoded fact.
- A pre-Whelk savestate is preserved, so the victory attempt
  (branch A: heal first, strike only while the head is extended)
  costs nothing to set up. Downstream systems (formation load,
  monster records, damage) remain Confirmed and carry over; the boss
  AI interpreter is unlocated.

### B19 — Stable post-Whelk state
**All rows PARTIAL: not reached** — the first Whelk attempt ended in
defeat (EXP-0039), so no post-victory state exists yet. Defines the
program's stopping boundary; requires victory processing capture
(rewards, flags, music, return-to-field or next-event handoff). The
"flags" part of that capture has a concrete substrate and baseline:
the event-flag arrays and the milestone-05 inventory (EXP-0037,
DISC-0008). A further scripted beat past the Whelk point
(CEN-EVENT-0011) was glimpsed and may be the handoff, but its
sequence position is unresolved.

## Per-domain totals

See `data/scenarios/opening-to-whelk.json` (`domain_totals`), kept in
sync with the beat matrices. As of creation: 0 beats COMPLETE,
19 PARTIAL, 0 BLOCKED.

## Unresolved gaps (honest list)

1. ~~No reproducible route from power-on exists~~ — segments 1–2 done
   (EXP-0031/0032: power-on → New Game → opening presentation →
   Narshe entry → first scripted battle; milestones 00–03,
   deterministic in WRAM). Segments 3+ (battle victory → free
   movement → mines → Whelk) remain.
1a. Milestone `01`'s **frame capture** is not byte-stable across runs
   even though its WRAM is (CEN-QUIRK-0002) — a PPU/HDMA-phase
   falsifier is queued; treat only WRAM as an assertion channel at
   that milestone until resolved.
2. The event-script engine (dispatcher, PC, opcode set) is entirely
   unlocated — the largest single blocker for B02–B07, B10, B17.
3. Map system (headers, tilesets, tilemaps, collision, exits,
   encounter zones) entirely unlocated — blocks B06, B08–B13, B16.
4. Encounter roll producer (`+$11E0`) unknown — blocks B12/B13
   completeness claims.
5. Monster names, AI scripts, and reward fields unmapped — blocks B14,
   B18.
6. Magitek command/ability records unresolved (CEN-MAGIC-0001).
7. Nothing from Whelk introduction onward has ever been observed live.
8. New-game initialization and RNG seeding unobserved (CEN-SAVE-0001).
9. Music tracks (opening, field, mines, battle, boss, victory) not
   located (CEN-AUDIO-0004/0005); only the confirm SFX chain is done.
10. Map/field graphics families observed only (CEN-GFX-0002..0004,
    CEN-WORLD-0001/0002).

## Definition of complete

The scenario is COMPLETE only per the fifteen criteria recorded in the
program brief (route reproducible; every beat represented; all
accessible maps/branches/encounters/monsters accounted for; event
opcodes decoded; Whelk fully tested; graphics/audio provenance;
persistence understood; Confirmed behavior implemented or deferral
documented; regression/differential tests; census+dashboards agree;
unresolved items listed; gates pass; no restricted artifacts tracked).
Anything less is reported as PARTIAL with the highest-leverage gap
named.
