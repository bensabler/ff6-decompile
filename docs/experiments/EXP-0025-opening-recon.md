# EXP-0025: Opening-sequence reconnaissance sweep

- **Status:** completed (2026-07-31)
- **Question (census unit):** which subsystems are visibly exercised
  by the four available opening-area savestates (`checkpoint1.mss`,
  `checkpoint2.mss`, `checkpoint3-mines.mss`, `exp10-battle.mss`), and
  what census entries do they warrant? This is a breadth sweep:
  observe and register, do not decode.
- **Starting state:** headless bridge session; each state loaded in
  turn; screenshots + targeted reads only, no deep tracing.
- **Method:** per state: load, settle ~60 frames, screenshot (local
  only, never committed), record what is visibly on screen (map,
  battle, HUD elements, sprites, dialogue), and where prior evidence
  already names the backing system, link it. Register every visible
  system in `manifests/content-census.json` at OBSERVED /
  CANDIDATE_LOCATION with the observation as evidence; no table
  structures inferred from visible English names alone.
- **Expected outcomes:** a screenshot + system list per state; new
  census entries for menu/world/event/monster/save domains that
  currently have none.
- **Falsifying outcome (per state):** the state fails to load or
  render (recorded as a negative result for that state).
- **Bounds:** no decoding, no write-watches, no new WRAM semantics;
  the only reads are screenshots and already-Confirmed addresses.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0025/`
  (five screenshots, hashes.sha256: checkpoint1 `d0958e41…`,
  checkpoint2 `b21135fa…`, checkpoint2-moved `dfe5926c…`,
  checkpoint3-mines `da47b2ef…`, exp10-battle `fae1e4a7…`).
- **Result:** all four states load and render headlessly (no
  falsifier). Observations, per state:
  - **checkpoint1** — mid-event scripted beat: dialogue window
    ("GUARD: …"), variable-width dialogue font, speaker formatting,
    Narshe exterior tileset with the player armor sprite and an NPC.
  - **checkpoint2** — controllable free-walk field state (a 30-frame
    right press moves the player and structures block movement:
    collision live); animated tiles (chimney smoke, water wheels).
  - **checkpoint3-mines** — mines interior: cave tileset, rails,
    glowing light sources; party follower sprites stacked.
  - **exp10-battle** — battle UI: command window (Row / clipped
    Magitek / Item), party status window with **'?????' as the first
    member's name** (Terra pre-naming), WEDGE/VICKS named, ATB
    arrows, yellow active-row highlight; green guard enemy sprite;
    cave battle background.
  - **25 census entries registered** across MENU/EVENT/WORLD/CHAR/
    MONSTER/GFX/AUDIO/MAGIC/SAVE (10 of 12 domains now populated;
    ITEM and QUIRK remain honestly empty — nothing observed yet).
    Registrations link prior Confirmed evidence where it already
    names the backing structures (+$2E78 display family, +$3B18
    stat staging, DSP voice snapshot).
- **Confidence:** per entry in the census (screen observations
  confirmed; backing-structure claims carried at candidate/unknown).
- **Next action:** EXP-0026 (Magic menu capture) — the two UNMAPPED
  magic/menu entries are its targets.
