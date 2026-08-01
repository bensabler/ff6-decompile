---
name: ff6-content-census
description: Maintains the FF6 Decompile project's global content census, ROM ownership map, content inventories, and reconstruction coverage. Use after experiments, gameplay reconnaissance, static ROM discoveries, new implementations, or whenever a subsystem, asset, data table, mechanic, menu, event, map, audio behavior, graphics behavior, or persistent state becomes visible. Registers newly observed systems without interrupting the active bounded investigation, identifies coverage gaps, synchronizes census records and dashboards, and helps prioritize high-leverage breadth targets. Observe broadly, register briefly, and investigate narrowly.
---

# FF6 Content Census

## Purpose

Use this skill to ensure the FF6 Decompile project eventually reconstructs the entire game rather than repeatedly investigating only the currently active subsystem.

The project uses two complementary research modes:

```text
Depth-first:
Choose one bounded question
→ investigate
→ verify
→ implement
→ test

Breadth-first:
Observe all visible systems
→ register them
→ locate candidate ownership
→ record unknowns
→ return to the bounded question
```

This skill governs the breadth-first process.

It must not cause uncontrolled scope expansion.

The operating rule is:

```text
Observe broadly.
Register briefly.
Investigate narrowly.
```

---

# When to use this skill

Invoke this skill when any of the following occurs:

- A Mesen experiment exposes a previously unregistered system.
- Gameplay enters a new map, battle type, menu, event, or game phase.
- A new item, spell, monster, command, status, formation, asset, sound, or script becomes visible.
- Static analysis reveals a new data table, pointer table, routine family, script format, or ROM region.
- A discovery or implementation is created.
- A checkpoint is prepared.
- A research unit ends.
- The orchestrator is selecting the next high-value target.
- Several consecutive experiments have focused on one subsystem.
- The coverage dashboard, content inventory, or ROM ownership map may be stale.
- A user asks whether the project is accounting for the entire game.

Do not wait for the user to invoke this skill manually when these conditions are met.

---

# Authoritative sources

Before updating the census, inspect:

1. The latest checkpoint.
2. The active experiment.
3. Relevant evidence records.
4. Existing discoveries.
5. Existing implementation and tests.
6. Current census entries.
7. The ROM ownership map.
8. Current coverage dashboards.
9. Relevant session documentation.
10. Current Git state.

Use this precedence:

```text
Direct runtime evidence
→ reproducible static evidence
→ verified experiment records
→ confirmed discoveries
→ tested implementation
→ current checkpoint
→ indexes and dashboards
→ hypotheses
→ external references
```

External walkthroughs, databases, guides, disassemblies, and wikis may provide leads, but they are not authoritative over the verified ROM revision.

---

# Controlled status models

Every content-census entry must track two independent dimensions.

## Reconstruction status

Use only:

```text
UNMAPPED
OBSERVED
CANDIDATE_LOCATION
LOCATED
FORMAT_PARTIAL
FORMAT_DECODED
EXTRACTED_PARTIAL
EXTRACTED_COMPLETE
IMPLEMENTED_PARTIAL
IMPLEMENTED_COMPLETE
REGRESSION_TESTED
DIFFERENTIALLY_VERIFIED
```

## Runtime coverage status

Use only:

```text
NOT_ENCOUNTERED
ENCOUNTERED
INPUTS_CAPTURED
OUTPUTS_CAPTURED
NORMAL_PATH_VERIFIED
ALTERNATE_PATH_VERIFIED
EDGE_CASES_PARTIAL
EDGE_CASES_COMPLETE
EXHAUSTIVELY_VERIFIED
```

Never combine these into one status.

Examples:

- A complete spell table may be `EXTRACTED_COMPLETE` but only `ENCOUNTERED`.
- A damage formula may be `DIFFERENTIALLY_VERIFIED` even while some related spell records remain unmapped.
- A monster may have fully extracted data but untested AI branches.

---

# Canonical domains

Register systems under the most appropriate domain.

## Battle

- Damage
- Healing
- Accuracy
- Evasion
- Critical hits
- Targeting
- Elements
- Status effects
- Counters
- Cover
- Row
- Battle orientations
- Running
- Death
- Rewards
- Drops
- Steals
- Battle scripts
- Formations
- Encounters
- Battle UI
- Battle transitions

## Characters

- Character records
- Stats
- Growth
- Equipment permissions
- Commands
- Learned magic
- Innate magic
- Esper learning
- Temporary characters
- Character naming
- Party formation
- Multi-party behavior
- Unique character mechanics

## Magic and abilities

- Spells
- MP costs
- Targeting
- Power
- Hit rates
- Elements
- Status behavior
- Reflect
- Runic
- Animations
- Sound effects
- Effect dispatch
- Learning sources
- Magitek
- Lores
- Summons
- Item-cast spells
- Enemy abilities

## Monsters

- Monster records
- Stats
- Resistances
- Weaknesses
- Immunities
- AI
- Special attacks
- Rages
- Sketch
- Control
- Drops
- Steals
- Rewards
- Formations
- Encounter packs
- Graphics
- Palettes
- Unused records

## Items and equipment

- Consumables
- Weapons
- Shields
- Helmets
- Armor
- Relics
- Key items
- Prices
- Shops
- Equipment permissions
- Stat effects
- Elements
- Status effects
- Proc effects
- Item-use effects
- Throw behavior
- Inventory behavior
- Unused records

## Maps and world

- Map headers
- Tilesets
- Tilemaps
- Palettes
- Collision
- Entrances
- Exits
- Transitions
- Encounter zones
- Vehicles
- NPCs
- Objects
- Treasure
- Hidden items
- Map-state changes
- World of Balance variants
- World of Ruin variants

## Events and story

- Event bytecode
- Event opcodes
- Event flags
- Dialogue
- Branches
- Timers
- Party changes
- Cutscenes
- Scripted battles
- Optional scenes
- Failure outcomes
- Mutually exclusive outcomes
- World-state transitions
- Ending logic

## Menus and interface

- Main menu
- Status
- Equipment
- Relics
- Inventory
- Magic
- Espers
- Configuration
- Save/load
- Shops
- Naming
- Party selection
- Colosseum
- Battle commands
- Fonts
- Windows
- Cursor logic
- Sorting
- Help text
- Target selection

## Graphics

- Character sprites
- Enemy sprites
- NPC sprites
- Portraits
- Fonts
- Icons
- Tiles
- Tilemaps
- Palettes
- Battle backgrounds
- Spell effects
- Animations
- Menus
- World maps
- Vehicles
- OAM
- DMA
- Compression

## Audio

- Music
- Sound effects
- BRR samples
- Sequences
- Instruments
- SPC driver
- CPU/APU communication
- DSP configuration
- Echo
- Track selection
- Transitions
- Fades
- Tempo
- Sound priorities

## Persistence

- SRAM
- Save slots
- Checksums
- Character state
- Inventory
- Event flags
- Treasure flags
- Map state
- Party state
- Learned spells
- Rage and Lore state
- Configuration
- Play time
- Steps
- RNG state
- New-game initialization
- Game-over behavior

## Compatibility and quirks

- Original bugs
- Overflow
- Underflow
- Dead code
- Unused data
- Glitches
- Timing behavior
- Hardware assumptions
- Formula discrepancies
- Incorrect descriptions

Create a new domain only when existing domains genuinely cannot represent the system.

---

# Required record structure

Each census record should contain:

```yaml
id:
domain:
category:
name:
description:

reconstruction_status:
runtime_status:
confidence:

rom_revision:
rom_locations:
wram_locations:
cpu_routines:
spc_locations:

record_count_expected:
record_count_found:
record_size:
pointer_tables:

producers:
consumers:

related_experiments:
related_discoveries:
related_implementations:
related_tests:
runtime_scenarios:

evidence:
unknown_fields:
contradictions:
next_action:
notes:
```

Do not fabricate values to fill the schema.

Use:

```yaml
record_count_expected: unknown
```

when the count has not been established.

Unknown fields must remain explicitly unknown.

---

# Reconnaissance procedure

After an experiment or gameplay session, perform the following review.

## Step 1 — Review what became visible

Ask:

- What commands appeared?
- What menus appeared?
- What assets were rendered?
- What audio played?
- What state changed?
- What tables must exist?
- What routines must exist?
- What scripts must exist?
- What data was displayed?
- What branches or alternate outcomes became possible?
- What system was present but unrelated to the active question?

## Step 2 — Compare against the census

For each observed system:

- Locate an existing census entry.
- Update it when new evidence exists.
- Create a new entry when none exists.
- Avoid duplicates caused by naming differences.
- Preserve uncertain identities as hypotheses.

## Step 3 — Record only bounded reconnaissance

For unrelated systems, record:

- Direct observation.
- Why the system must exist.
- Candidate ownership, if supported.
- Confidence.
- Most efficient future experiment.
- Relevant evidence path.

Do not begin a new investigation merely because the system was observed.

## Step 4 — Return to the active question

Reconnaissance should normally consume only a small portion of the research unit.

Do not abandon the current bounded experiment unless the observation reveals:

- A destructive repository problem.
- Evidence invalidating the active experiment.
- A legal-boundary problem.
- A prerequisite without which the experiment cannot continue.
- A uniquely temporary state that cannot be reproduced later.

---

# Magic-system handling

When the Magic command, spell data, or magical behavior becomes visible, ensure separate census entries exist for:

- Global spell identifiers
- Spell names
- Spell property data
- MP costs
- Targeting
- Elements
- Status behavior
- Accuracy
- Effect dispatch
- Animation
- Sound
- Character spell availability
- Learned-spell storage
- Esper learning
- Innate spell initialization
- Magic menu rendering
- MP validation
- MP deduction
- Enemy spell use
- Item-cast spells
- Lores
- Magitek abilities
- Summons

Do not assume all spell properties occupy one contiguous record.

The original structure may consist of:

- Parallel tables
- Bitfields
- Pointer tables
- Command-specific tables
- Shared effect routines
- Separate name and cost arrays
- Animation and audio lookup tables

Represent the recovered structure faithfully before designing the eventual modern Go model.

A preliminary spell inventory should track:

```yaml
id:
name:
mp_cost:
targeting:
power:
hit_rate:
element:
status_application:
status_removal:
reflect_behavior:
runic_behavior:
effect_routine:
animation:
sound:
learning_sources:
raw_property_locations:
unknown_fields:
confidence:
evidence:
```

Fields unsupported by evidence must remain unknown.

---

# ROM ownership map

Maintain the ROM-region ledger whenever a new region is discovered.

Each region should record:

```yaml
start:
end:
size:
bank:
classification:
owner:
format:
status:
confidence:
evidence:
consumers:
overlaps:
notes:
```

Valid classifications include:

```text
CODE
POINTER_TABLE
DATA_TABLE
TEXT
GRAPHICS
PALETTE
TILEMAP
ANIMATION
MUSIC
SOUND
BRR
EVENT_SCRIPT
AI_SCRIPT
PADDING
UNKNOWN
```

Rules:

- Do not assign ownership merely to reduce unknown space.
- Record overlaps explicitly.
- Distinguish candidate ownership from confirmed ownership.
- Split a region only when evidence supports the boundary.
- Preserve revision-specific addresses.

---

# Data inventories

Ensure inventory containers exist for major content families:

```text
data/census/spells.json
data/census/items.json
data/census/equipment.json
data/census/monsters.json
data/census/formations.json
data/census/encounter-packs.json
data/census/monster-ai.json
data/census/shops.json
data/census/treasures.json
data/census/maps.json
data/census/map-exits.json
data/census/events.json
data/census/dialogue.json
data/census/characters.json
data/census/commands.json
data/census/graphics.json
data/census/music.json
data/census/sound-effects.json
data/census/save-data.json
```

An inventory may be empty, but it must not contain fabricated placeholder records.

Each inventory should identify:

- Schema version
- ROM revision
- Domain
- Known records
- Evidence
- Unknowns
- Expected count, when supported
- Next extraction step

Do not commit copyrighted ROM-derived assets, dialogue, audio, or binary extracts when prohibited by the project’s legal rules.

Prefer:

- IDs
- Addresses
- Sizes
- Hashes
- Structural metadata
- Project-authored descriptions
- Generated definitions
- Legal test fixtures

---

# Coverage synchronization

After census changes, update all applicable:

```text
manifests/content-census.json
manifests/rom-regions.json
indexes/CONTENT_CENSUS.md
indexes/ROM_REGIONS.md
dashboards/COVERAGE.md
dashboards/RESEARCH_QUEUE.md
dashboards/MILESTONES.md
dashboards/OPEN_HYPOTHESES.md
```

Also update related:

- Experiments
- Discoveries
- Functions
- Variables
- Structures
- Implementations
- Tests
- Checkpoints

Generated files should be regenerated rather than manually edited when tooling exists.

---

# Coverage audit

Run the project’s census and coverage checks.

The audit should detect:

- Invalid statuses
- Duplicate IDs
- Missing required fields
- Entries without evidence
- Stale generated indexes
- Implementations without census entries
- Discoveries without domain ownership
- Experiments that exposed unregistered systems
- Inventories without ROM-region relationships
- Conflicting ROM ownership
- Unknown ROM gaps
- Extracted data without tests
- Runtime-tested behavior without implementation links
- Speculative names presented as confirmed
- Restricted material accidentally tracked

Useful commands may include:

```text
ff6lab census validate
ff6lab census sync
ff6lab coverage summary
ff6lab coverage gaps
ff6lab coverage domain <domain>
ff6lab rom gaps
ff6lab audit
```

Use the repository’s actual command names when they differ.

---

# Research prioritization

When selecting future work, prefer high-leverage questions that unlock:

- Entire data tables
- Pointer systems
- Script interpreters
- Compression formats
- Shared engines
- Reusable asset formats
- Complete content families

Do not permit endless depth in one subsystem.

After several consecutive experiments in one subsystem, perform a global coverage review.

Balance work across:

- Behavior
- Static content
- Graphics
- Audio
- Maps and events
- Menus
- Persistence
- Compatibility

A generic SNES decoder does not count as FF6-specific coverage until connected to:

- An actual FF6 ROM region
- Its consuming routine
- Runtime evidence
- A deterministic comparison

---

# Completion requirements

Before completing a census update:

1. Validate all changed census records.
2. Synchronize indexes and dashboards.
3. Validate ROM-region ownership records.
4. Run relevant Go tests and audits.
5. Confirm no restricted artifact entered Git.
6. Update the latest checkpoint.
7. Record the exact next action.
8. Commit the complete bounded unit when appropriate.

The final report should state:

- Systems newly registered
- Existing entries updated
- New ROM regions identified
- Content inventories affected
- What was directly observed
- What remains inferred
- Highest-value unanswered breadth question
- Audit and Git status

---

# Anti-drift rule

Never interpret “maintain complete coverage” as an instruction to investigate every observation immediately.

The required behavior is:

```text
Notice it.
Register it.
Preserve the lead.
Finish the current experiment.
Prioritize it later using global coverage.
```
