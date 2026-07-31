# EXP-0026: Magic command census — opening spell list and its sources

- **Status:** completed (2026-07-31)
- **Question (census unit, breadth proof):** what spells does the
  opening party's Magic list show (names, order, MP costs,
  availability), and where do the displayed names and costs come from
  — character-specific availability state and global spell records —
  at the candidate level?
- **Starting state:** headless bridge session; field menu attempted
  from the controllable `checkpoint2.mss` (fallback: the battle
  command menu in `exp10-battle.mss` when the field menu is locked
  during the armor sequence).
- **Method:**
  1. Screenshot-guided navigation to the Magic list; archive every
     step's screenshot.
  2. Read the visible names/order/MP costs **from the screen** (no
     external name lists).
  3. Atomic capture while the list is visible (EXP-0023 probe
     pattern): menu BG tilemap + chr region + CGRAM + a full WRAM
     dump.
  4. Derive each rendered name's tile-index sequence from the
     tilemap; attempt an affine tile-to-byte mapping against a ROM
     byte search (the technique that located the HUD font) to find
     the **name table** candidate. If menu text is composed into
     dynamic tiles (the battle compose-region pattern), fall back to
     the WRAM path: search the WRAM dump for the screen-read MP-cost
     sequence and for plausible availability structures near the
     party state; write-watch their producers in a follow-up rather
     than this unit.
  5. Register census entries (global spell database, opening
     availability, Magic menu, name text, MP display/deduction,
     targeting, effect dispatch, animation, audio) at honest
     statuses, and populate `data/census/spells.json` with the
     observed records only.
- **Expected outcomes:**
  - *Direct:* the spell list renders; names/costs recorded; at least
    one of {name table, cost source, availability structure} reaches
    CANDIDATE_LOCATION with a bounded confirmation experiment named.
  - *Bounded-out:* menu locked in all states, or text fully
    dynamic-composed with no affine match — recorded as the negative
    result; the census entries still register with the fallback path
    as next_action.
- **Falsifying outcome (for the availability question):** the Magic
  command is absent for every opening party member in both field and
  battle menus (would refute the premise that Magic is exercisable in
  the opening sequence).
- **Important constraint:** spell properties may live in parallel
  tables, pointer tables, flags, or command-specific data — the
  census records the actual observed structure; nothing is forced
  into one contiguous record layout without evidence.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0026/`
  (9 step screenshots, menu tilemap/chr/CGRAM dumps, full 128 KB WRAM
  dump in 8 chunks; hashes.sha256, 22 files).
- **Result:** the field menu opened from `checkpoint2.mss` (X press) —
  the falsifier is not met; every targeted structure reached at least
  candidate level, three reached located/decoded.
  - **Menu observations:** main menu (Item/Skills/Equip/Relic/Status/
    Config/Save-greyed) with portraits, LV/HP/MP per member, play
    time, steps, Gp. Terra's Skills submenu: Espers+Magic available;
    SwdTech/Blitz/Lore/Rage/Dance greyed. The Magic list is a
    **fixed grid with blank unlearned slots**: Cure castable
    (bullet `$E8`, MP window 5, help "Recovers HP"), Fire dim
    (bullet `$E9`, MP window 0, empty help — field-castability
    behavior, interpretation open).
  - **Text encoding derived** from the menu tilemap alone:
    A-Z=`$80+n`, a-z=`$9A+n`, `-`=`$C4`, narrow space `$FE`,
    digits `2`/`3`=`$B6`/`$B7`, blank `$FF` — and ROM strings use
    the same byte values (search hits), cross-confirmed by the
    battle-HUD tilemap ("Were-Rat").
  - **Spell name table located and format-decoded:**
    `ROMFILE:0x26F567`, 7-byte `[class icon][6-char name]` records;
    Fire=id 0, **Cure=id 45** (`0x26F6A3`); 36 entries decoded
    contiguously; class-icon byte changes exactly at the menu bullet
    values. 54-entry extent inferred from class structure
    (boundary unverified).
  - **The $C46AC0 record table is the global spell database:**
    ids match the name table; **record byte 5 = MP cost (Strong
    hypothesis)** — Cure's record holds the observed 5, Cure2/Cure3
    ascend 25/40, Fire/Ice/Bolt run 4/5/6. One of EXP-0019's seven
    unknown bytes now has a candidate meaning.
  - **Spell availability located:** a blind 128 KB WRAM search for a
    54-byte window with one repeated mark at exactly positions
    {0, 45} produced **one hit: `WRAM:+$1A6E`** ($FF = known). The
    next three 54-byte windows are all zero (the magicless
    soldiers) — per-character stride is a Strong hypothesis.
  - Preliminary inventory written to `data/census/spells.json`
    (structure + the two observed records; bulk text stays local).
- **Confidence:** encoding (observed glyphs) — Confirmed. Name-table
  location/format — Strong hypothesis (36-entry decode + id-45
  anchor; boundary open). Byte 5 = MP cost — Strong hypothesis
  (single on-screen anchor). +$1A6E availability array — Strong
  hypothesis (unique search hit + zero neighbors; single character
  observed). Field dim/MP-0 behavior — observation recorded,
  interpretation open.
- **Next action:** queued follow-ups — battle Magic menu for a second
  cost anchor and Fire's display; name-table boundary read; a second
  magic-knowing character for the availability stride; help-text
  table hunt.
