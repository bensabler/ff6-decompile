# EXP-0027: Spell-table boundary, second cost anchor, bulk extraction

- **Status:** completed (2026-07-31)
- **Question:** where does the 7-byte spell name table end (is the
  54-entry structure hypothesis right?), does a second on-screen MP
  cost match record byte 5, and can all magic-spell records be
  extracted with defensible counts?
- **Starting state:** static ROM analysis from the local ROM file
  (SHA-256 `0f51b4fc…d8d5e2`) plus, if reachable within bounds, a
  battle Magic menu capture from `exp10-battle.mss` (headless).
- **Method:**
  1. *Boundary (static):* decode the name-table region at stride 7
     from entry 36 through entry ~60. The boundary signals are: the
     class-icon byte pattern ending, and/or a structurally different
     table (different stride/charset) beginning after entry 53.
  2. *Second cost anchor (bounded live attempt):* navigate the battle
     command menu to Terra's Magic list (her battle list should show
     Fire castable with its cost). Budget: bounded navigation
     attempts guided by screenshots; if Terra's turn cannot be
     reached cleanly, record bounded-out — extraction proceeds with
     the cost field held at Strong hypothesis.
  3. *Bulk extraction (local):* write all 54 name records (decoded)
     and 14-byte data records to
     `local_artifacts/experiments/EXP-0027/` with hashes; mirror ids
     and numeric fields only into `data/census/spells.json`
     (names stay local per the asset policy).
  4. Census/status updates per outcome; ROM region confidence
     adjustments.
- **Expected outcomes:**
  - *Supports:* icon-classed 7-byte entries end exactly after entry
    53 → count 54 Confirmed; extraction complete.
  - *Refutes:* the pattern continues past entry 53 (count wrong —
    census corrected accordingly) or breaks earlier.
- **Falsifying outcome (for the 54-count hypothesis):** clean
  icon-classed entries continuing at ids 54+ with magic-style names,
  or a break before id 53.
- **Raw evidence paths:** `local_artifacts/experiments/EXP-0027/`
  (boundary-decode.txt, esper-decode.txt, spell-records.json, battle
  and field screenshots, post-cast WRAM dump; hashes.sha256,
  17 files).
- **Result:** all three questions answered; the falsifier is not met.
  - **Boundary Confirmed at exactly 54 entries:** entries 36-53
    decode cleanly (grey magic through the white-magic block, icon
    `$E8` ×9 ending at id 53); entry 54 breaks stride-7 and the
    **esper name table** begins at `0x26F6E1` — 27 names decode
    perfectly at stride 8, with a further ability-name table
    (stride ~10) starting at `0x26F7B9` (candidate, unregistered
    pending decode). Class-icon census: `$E9`×24, `$EA`×21,
    `$E8`×9 = 54. The menu bullet glyph **is** the class icon from
    the name record.
  - **Second cost anchor obtained behaviorally:** the battle route
    was unreachable (the savestate's party is at 0 HP from trial
    damage; unattended waiting produced the 'Annihilated' defeat
    flow — itself registered, along with the on-screen formation
    names Were-Rat + Repo Man). The field route worked: casting
    Cure showed a **'5 MP Needed'** targeting gate, healed
    34→77/77, and **deducted exactly 5 MP (24→19)** — record
    byte 5 = MP cost is **Confirmed for id 45** (menu readout +
    gate text + live deduction); other ids carry the extracted
    values pending per-id observation.
  - **Bulk extraction complete:** all 54 name records and 14-byte
    data records extracted to the local archive; ids + numeric
    fields mirrored into `data/census/spells.json` (names stay
    local per the asset policy).
  - **Bonus (single-byte WRAM diffs across the cast):** field
    current MP slot 0 at `WRAM:+$160D` (24→19) and field current
    HP at `WRAM:+$1609` (34→77) — the field character-record block
    (~`+$1600`) located as a by-product; registered as
    CEN-CHAR-0004.
- **Confidence:** 54-entry count and name-table extent — Confirmed.
  Esper table (27×8) — Confirmed (clean stride decode). Byte 5 =
  MP cost — Confirmed for id 45, extracted values for the rest.
  Field HP/MP store bytes — Strong hypothesis (single cast, single
  slot). Ability-name table at 0x26F7B9 — Tentative.
- **Next action:** census-observation pass complete (defeat flow,
  formation names, field record block registered). Next breadth
  targets per coverage gaps: monster stat-record source trace, or
  record bytes 0/7/8/11-13 semantics via varied casts.
