# Open Hypotheses

| ID | Claim | Confidence | Next discriminator |
|---|---|---|---|
| H-BATTLE-0001 | ~~Destination records in the observed `$C10DF3` loop use a `$20` stride.~~ | **Resolved: Confirmed** (Session 002) | None — see [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) and [04_MEMORY_MAP.md](../docs/sessions/04_MEMORY_MAP.md) |
| H-BATTLE-0002 | `WRAM:+$2E80` array holds max HP. | Strong hypothesis | An independent max-HP change (level-up or max-HP-boosting effect) tracking the array — open question #11 |
| H-BATTLE-0003 | `ROMCPU:$C101FB` is a battle-only per-frame dispatcher. | Strong hypothesis | Exec-log `$C101FB` on the world map and in menus — open question #6/#7 context |
| H-BATTLE-0004 | `WRAM:+$2E88`/`+$2E90` hold current/max MP. | Tentative hypothesis | Battle with an MP-spending party member; watch `+$2E88` — open question #9 |
| H-BATTLE-0005 | Destination records are display/staging data (not gameplay-authoritative). | Strong hypothesis | Freeze/edit `WRAM:+$2EB5` vs `+$2E78[0]` separately; find the records' consumer — open questions #1/#2 |
| H-BATTLE-0006 | `WRAM:+$3BF4` is the gameplay-authoritative battle current-HP array. | Strong hypothesis | Find the `+$3BF4`→`+$2E78` copier and any reader that feeds damage formulas — open questions #1/#13 |
| H-BATTLE-0007 | `WRAM:+$3C08`/`+$3C30` are current/max MP; `+$3C1C` is max HP. | Strong hypothesis (code roles Confirmed, EXP-0001) | Spend/heal MP live while watching `+$3C08`; compare `+$3C1C` against known max HP values |

Claims implemented in `battle/battle.go` are deliberately **not** listed
here: they lack canonical records and cannot be tracked as hypotheses until
Session 003 is documented (see [BLOCKERS.md](BLOCKERS.md)).
