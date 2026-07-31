# Content Taxonomy and Coverage Model

Adopted 2026-07-31 (census unit). This is the canonical breadth
inventory: every system, database, asset family, or mechanic the game
contains belongs to exactly one domain below and is tracked as a
census entry in `manifests/content-census.json`. The taxonomy makes
"everything" measurable — a system missing from the census is a
defect, not an accident.

The operating rule is:

```text
Observe broadly. Register briefly. Investigate narrowly.
```

Depth-first experiments continue as before; the census adds the
parallel breadth process (observe → register → identify candidate ROM
ownership → record unknowns → return to the focused experiment).

## Domains

Census IDs are `CEN-<DOMAIN>-<NNNN>` using the domain keys below
(the hypothesis convention, e.g. `H-BATTLE-0004`, extended).

| Key | Domain | Scope (non-exhaustive category list) |
|---|---|---|
| BATTLE | Battle engine | damage formulas, healing, targeting, elements, status effects, physical/magical accuracy, evasion, critical hits, counters, cover, row, back/pincer/side/preemptive, running, death, rewards (experience, magic points/AP, gil, drops, steals, metamorph), battle scripts, formations, transitions, backgrounds, battle UI |
| CHAR | Characters | character records, base stats, level growth, equipment permissions, commands, learned magic, innate magic, esper learning, trance/morph, blitz, bushido, tools, runic, rage, leap, lore, sketch, control, dance, slots, throw, mimic, desperation attacks, temporary/guest characters, naming, party formation, multi-party behavior |
| MAGIC | Magic and abilities | complete spell list, MP costs, power, hit rate, targeting flags, element, status application/removal, reflect behavior, runic behavior, animation, sound, effect routine, learning sources, enemy-only abilities, magitek abilities, lores, summons, item-cast spells, special command abilities |
| MONSTER | Monsters | complete monster list, stats, resistances/weaknesses/immunities, status properties, AI scripts, special attacks, rage/sketch/control data, steals, drops, experience, gil, formations, encounter packs, boss flags, graphics, palettes, animation, unused records |
| ITEM | Items and equipment | complete item list, consumables, weapons, shields, helmets, armor, relics, key items, prices, shops, equipment permissions, stat modifiers, elements, status effects, proc effects, item-use effects, throw properties, two-handed, dual-wield, special flags, inventory limits, sorting, rare/unused entries |
| WORLD | World, maps, navigation | world maps, interior maps, map headers, tilesets, tilemaps, palettes, collision, exits, entrances, transitions, encounter zones, vehicles, chocobos, airship, raft, NPC placement, object placement, treasure, hidden items, changing map state, WoB/WoR variants |
| EVENT | Events and story | event bytecode, opcodes, event flags, dialogue, branches, timers, naming events, party changes, cutscenes, scripted battles, optional scenes, failure states, mutually exclusive outcomes, world-state changes, opera, banquet, floating continent, Cid outcome, Shadow outcome, final-dungeon party switching, ending logic |
| MENU | Menus and interface | main menu, status, equipment, relics, items, magic, espers, configuration, save/load, shop menus, naming screen, party selection, colosseum, battle command menus, fonts, windows, cursor logic, sorting/filtering, help text, target selection |
| GFX | Graphics | character/enemy/NPC sprites, portraits, fonts, icons, tiles, tilemaps, palettes, battle backgrounds, spell effects, command animations, menu graphics, world-map graphics, vehicle graphics, OAM composition, DMA behavior, compression formats, animation tables |
| AUDIO | Audio | music tracks, sound effects, BRR samples, sequences, instruments, SPC driver, CPU/APU communication, DSP configuration, echo, track transitions, battle/map/event music selection, fades, tempo changes, sound priorities |
| SAVE | Persistence and global state | SRAM layout, save slots, checksums, character state, inventory state, event flags, treasure flags, map state, party state, learned spells, rage/lore state, configuration, play time, steps, RNG state, new-game initialization, game-over behavior, version fields |
| QUIRK | Compatibility and quirks | original bugs, overflow/underflow, unused data, dead code, glitches, timing-sensitive behavior, incorrect descriptions, formula discrepancies, platform-specific and hardware-dependent behavior |

Add domains only when evidence shows the taxonomy is incomplete;
record the addition here and in the schema enum in the same commit.

## Two independent status models

Every census entry carries **both** statuses. They never collapse: a
table can be fully extracted while its runtime behavior is unverified,
and a mechanic can be runtime-verified while its content records are
unextracted.

### Reconstruction status (static/data-side)

| Status | Meaning |
|---|---|
| UNMAPPED | Known to exist (taxonomy) but nothing recorded |
| OBSERVED | Seen in gameplay/evidence; no location claim |
| CANDIDATE_LOCATION | ROM/WRAM candidates named, unverified |
| LOCATED | Location confirmed by evidence |
| FORMAT_PARTIAL | Some fields/semantics decoded |
| FORMAT_DECODED | Record layout decoded (unknown fields allowed if enumerated) |
| EXTRACTED_PARTIAL | Some records extracted to local evidence |
| EXTRACTED_COMPLETE | All records extracted (count evidence required) |
| IMPLEMENTED_PARTIAL | Go implementation covers part of the format |
| IMPLEMENTED_COMPLETE | Go implementation covers the decoded format |
| REGRESSION_TESTED | Implementation has tests encoding the evidence |
| DIFFERENTIALLY_VERIFIED | Output compared against live/original behavior |

### Runtime coverage status (behavior-side)

| Status | Meaning |
|---|---|
| NOT_ENCOUNTERED | Never exercised in any captured session |
| ENCOUNTERED | Seen live at least once |
| INPUTS_CAPTURED | Triggering inputs/state captured with evidence |
| OUTPUTS_CAPTURED | Resulting outputs/state captured with evidence |
| NORMAL_PATH_VERIFIED | Main path verified against expectation |
| ALTERNATE_PATH_VERIFIED | At least one alternate path verified |
| EDGE_CASES_PARTIAL | Some boundary conditions verified |
| EDGE_CASES_COMPLETE | Known boundary conditions verified |
| EXHAUSTIVELY_VERIFIED | Input space swept or proven equivalent |

Statuses are orderable ladders; the coverage tooling
(`ff6lab coverage`) counts an entry at a threshold when its status is
at or above it. Confidence stays the project's four-level standard and
is tracked separately from both ladders.

## Record linkage

Census entries link outward: `related_experiments` (EXP ids),
`related_discoveries` (DISC ids), `related_implementations` (Go
symbols/packages), `related_tests` (test names), `rom_locations`
(ROM region ids from `manifests/rom-regions.json` or address
citations), and `runtime_scenarios` (savestates/sessions). The audit
enforces the cross-links (`ff6lab census validate`).

## ROM ownership map

`manifests/rom-regions.json` is the ledger that must eventually
account for the whole 3,145,728-byte ROM. Regions carry
classification (CODE, POINTER_TABLE, DATA_TABLE, TEXT, GRAPHICS,
PALETTE, TILEMAP, ANIMATION, MUSIC, SOUND, BRR, EVENT_SCRIPT,
AI_SCRIPT, HEADER, PADDING, UNKNOWN), status (`known` or `candidate`),
confidence, evidence, and explicit `overlaps` when two
interpretations intentionally coexist. Unknown gaps are computed by
`ff6lab rom gaps`, never papered over: do not invent ownership to
shrink the unknown percentage.
