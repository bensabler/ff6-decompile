# Open Hypotheses

| ID | Claim | Confidence | Next discriminator |
|---|---|---|---|
| H-BATTLE-0001 | ~~Destination records in the observed `$C10DF3` loop use a `$20` stride.~~ | **Resolved: Confirmed** (Session 002) | None — see [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) and [04_MEMORY_MAP.md](../docs/sessions/04_MEMORY_MAP.md) |
| H-BATTLE-0002 | ~~`WRAM:+$2E80` array holds max HP.~~ | **Resolved: Confirmed** (EXP-0003: `+$2E80` is the display copy of `+$3C1C`, the heal-clamp ceiling — operationally max HP) | None — see [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
| H-BATTLE-0003 | `ROMCPU:$C101FB` is a battle-only per-frame dispatcher. | Strong hypothesis — **Confirmed for tested contexts** (EXP-0002: title + 3 field/event states, ≈175k frames, zero fires; fires in battle within ~30 frames of entry) | Sample the world map and the in-field menu — the remaining untested contexts |
| H-BATTLE-0004 | `WRAM:+$2E88`/`+$2E90` hold current/max MP. | Strong hypothesis (upgraded by EXP-0003: they are display copies of the MP-path pool `+$3C08` and ceiling `+$3C30`) | Live MP spend/heal observation — open question #9 |
| H-BATTLE-0005 | Destination records are display/staging data (not gameplay-authoritative). | Strong hypothesis | Freeze/edit `WRAM:+$2EB5` vs `+$2E78[0]` separately; find the records' consumer — open questions #1/#2 |
| H-BATTLE-0006 | ~~`WRAM:+$3BF4` is the battle-authoritative current-HP array.~~ | **Resolved: Confirmed (battle-scoped)** — mutations originate there (delta engine) and the display path derives from it via the `$C25D26` copier (EXP-0003) | None for battle scope; out-of-battle persistence is a separate future question |
| H-BATTLE-0007 | `WRAM:+$3C08`/`+$3C30` are current/max MP; `+$3C1C` is max HP. | Strong hypothesis (code roles Confirmed, EXP-0001) | Spend/heal MP live while watching `+$3C08`; compare `+$3C1C` against known max HP values |

Claims implemented in `battle/battle.go` are deliberately **not** listed
here: they lack canonical records and cannot be tracked as hypotheses until
Session 003 is documented (see [BLOCKERS.md](BLOCKERS.md)).
