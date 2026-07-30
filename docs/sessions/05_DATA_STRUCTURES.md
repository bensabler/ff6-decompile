# Data Structures

### CharacterFieldsSource (source region, struct-of-arrays)

- **Address:** `$2E78–$2EA7` (WRAM-relative)
- **Address space:** WRAM-relative
- **Kind:** Structure (six parallel arrays of four 16-bit entries)
- **Status:** Confirmed (layout and slot order); individual field meanings
  vary — see table
- **Previous names:** Session 001 knew only the first four arrays.
- **Observed behavior:** Read every frame by CopyCharacterFields
  (`LDA $xxxx,X`, X ∈ {0,2,4,6}). Entry order matches the displayed party
  list (verified for slots 0–2 against the battle HP window).
- **Evidence:** Full disassembly; live WRAM dumps during the Narshe intro
  battle, tracked through damage and a heal:

  | Array | Slot 0 | Slot 1 | Slot 2 | Slot 3 |
  | --- | --- | --- | --- | --- |
  | `$2E78` | 53→63 (healed)→10 | 39→29 | 33→29→3 | 0 |
  | `$2E80` | 63 | 68 | 70 | 0 |
  | `$2E88` | 24 | 0 | 0 | 0 |
  | `$2E90` | 24 | 0 | 0 | 0 |
  | `$2E98` | 8 | 8 | 8 (later `$0208`) | 0 |
  | `$2EA0` | 0 | 0 | 0 | 0 |

  Displayed HP (`?????`/WEDGE/VICKS = 53/39/33) matched `$2E78[0..2]`
  exactly; a Heal Force cast set `$2E78[0]` to exactly `$2E80[0]`.

- **Interpretation** (updated 2026-07-29 by
  [EXP-0003](../experiments/EXP-0003-2e78-producer.md): every array is a
  copy of an authoritative battle array, written by
  PartyDisplaySourceRefresh `ROMCPU:$C25D26`):

  | Array base | Copied from | Feeds record offset | Meaning | Status |
  | --- | --- | --- | --- | --- |
  | `$2E78` | `WRAM:+$3BF4` | `+$00` | Current HP | **Confirmed** (display match + delta-engine chain) |
  | `$2E80` | `WRAM:+$3C1C` | `+$02` | Max HP | **Confirmed** (operational: heal-clamp ceiling, heal-snap observation, gauge maximum) |
  | `$2E88` | `WRAM:+$3C08` | `+$04` | Current MP | Strong hypothesis (copy of the MP-path pool; no live MP observation yet) |
  | `$2E90` | `WRAM:+$3C30` | `+$06` | Max MP | Strong hypothesis (copy of the MP-path ceiling) |
  | `$2E98` | `WRAM:+$3EE4` | `+$08` (maskable `$0038`) | Status word (bit 1 = death-event suppression; bits 3–5 survive masked mode) | Partially confirmed (source identity Confirmed; bit meanings mostly Unknown) |
  | `$2EA0` | `WRAM:+$3EF8` | `+$0A` (maskable to 0) | Unknown; bit 13 drives `$61AD` mask | Source identity Confirmed; meaning Unknown |

- **Alternative explanations:** `$2E88`/`$2E90` could be any current/max
  pair; slot 3 all-zero because the intro party has 3 members — a 4-member
  party has not been observed.
- **Validation experiment:** Spend/observe MP in a battle where the party
  has MP users and watch `$2E88`; get a 4-member party and verify slot 3.
- **Go representation:** `chardata.CharacterFieldsSource`
- **Related discoveries:** CopyCharacterFields, PartySlotRecord
- **First observed in session:** [SESSION_001](SESSION_001.md)
- **Last updated:** 2026-07-29 ([SESSION_002](SESSION_002.md))

### BattleSlotHPArray (candidate authoritative layer) — `WRAM:+$3BF4`

- **Address:** `WRAM:+$3BF4` family — **10-entry unified battle arrays,
  entries 0–3 party / 4–9 enemies (Confirmed, EXP-0005)**: `+$3BF4`
  (current HP), `+$3C08` (current MP cand.), `+$3C1C` (max HP), `+$3C30`
  (max MP cand.), each `$14` apart. The status family repeats the
  stride: `+$3EE4` + `$14` = `+$3EF8`. Related: `+$3C95` (flags),
  pending-delta arrays `+$33D0`/`+$33E4` (10-entry; sweeper wrote
  entry 9). Enemy entries 4–5 observed live (HP 24/35 → damaged →
  zeroed by the shared engine); entries 6–9 presumed for larger
  encounters (stride-bounded, unobserved).
- **Address space:** WRAM-relative
- **Kind:** Structure (10-entry per-slot arrays; struct-of-arrays family)
- **Status:** Writers **Confirmed** (raw captures); engine code operating
  the family **Confirmed byte-exact** (EXP-0001); role as authoritative
  current HP **Strong hypothesis** (consumer/propagation path still
  unidentified); semantic labels Strong hypothesis.
- **Observed behavior:** Damage/heal/death stores land here with
  `Y = slot×2`; `$FF` fill at boot/teardown by `ROMCPU:$C0567B`; init at
  battle start (`ROMCPU:$C223F6/$C227B4`); slot-0 write of `$0022` was
  followed by `WRAM:+$2E78[0] = $0022` in the display-source array
  (propagation mechanism unidentified).
- **Evidence:** `mesen/out/events.log`; [SESSION_003](SESSION_003.md).
- **Interpretation / Alternatives:** authoritative battle HP vs. both
  arrays mirroring a deeper store — unresolved; see SESSION_003
  interpretation 2.
- **Validation experiment:** write watch on `WRAM:+$2E78` with stack
  capture during battle damage (finds the copier); watches on the claimed
  sibling arrays during MP spend/heal.
- **Go representation:**
  [internal/game/battle/battle.go](../../internal/game/battle/battle.go)
  (`PartySlots`) — encodes claimed siblings; hypothesis encoding until
  re-verified.
- **Related discoveries:** Battle HP/MP delta engine, CharacterFieldsSource
  (downstream), [04_MEMORY_MAP.md](04_MEMORY_MAP.md)
- **First observed in session:** [SESSION_003](SESSION_003.md)
- **Last updated:** 2026-07-29

### PartySlotRecord (destination record, array-of-structs)

- **Address:** records at `$2EB5`, `$2ED5`, `$2EF5`, `$2F15` (stride `$20`)
- **Address space:** WRAM-relative
- **Kind:** Structure
- **Status:** Confirmed (stride, six written fields, slot correspondence for
  records 0–2); bytes `+$0C..$1F` Unknown
- **Previous names:** Session 001 knew four fields and hypothesized the
  stride.
- **Observed behavior:** Rewritten every frame by CopyCharacterFields at
  offsets `+0..+$B`; live dump showed records 0–2 holding the three
  displayed party members' values, and record bytes `+$0C..$1F` holding
  other live data (`08 00 … FF FF FF FF FF 0E 96 84 83 86 84 FF` for
  record 0) that the routine never writes.
- **Interpretation:** Layout relative to record base (`$2EB5 + n*$20`;
  whether that is the true record start is still open):

  | Offset | Size | Meaning | Status |
  | --- | --- | --- | --- |
  | `+$00` | 2 | Current HP | **Confirmed** (records 0–2) |
  | `+$02` | 2 | Max HP | Strong hypothesis |
  | `+$04` | 2 | Current MP? | Tentative |
  | `+$06` | 2 | Max MP? | Tentative |
  | `+$08` | 2 | Unknown (maskable with `$0038`) | Unknown |
  | `+$0A` | 2 | Unknown (maskable to 0) | Unknown |
  | `+$0C–$1F` | 20 | Written by other code; contents observed non-zero | Unknown |

- **Alternative explanations:** Records could extend before `$2EB5`; record
  3 unverified (no 4-member party observed).
- **Validation experiment:** Write breakpoints on `$2EC1–$2ED4`
  (record 0, `+$C..$1F`) to identify the other writer(s).
- **Go representation:** `chardata.PartySlotRecord`
- **Related discoveries:** CopyCharacterFields, PartySlotCurrentHP,
  UnknownSlotMask61AD
- **First observed in session:** [SESSION_001](SESSION_001.md)
- **Last updated:** 2026-07-29 ([SESSION_002](SESSION_002.md))
