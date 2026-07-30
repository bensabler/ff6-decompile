# Structures

| ID | Title | Status | Confidence | Record |
|---|---|---|---|---|
| ST-0001 | CharacterFieldsSource — `WRAM:+$2E78`–`+$2EA7` (struct-of-arrays) | Documented; implemented (`chardata.CharacterFieldsSource`) | Confirmed (layout, slot order); field meanings vary by field — see record | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
| ST-0002 | PartySlotRecord — `WRAM:+$2EB5` + n×`$20` (array-of-structs) | Documented; implemented (`chardata.PartySlotRecord`) | Confirmed (stride, six written fields, slots 0–2); bytes `+$0C..$1F` Unknown | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
| ST-0003 | Unified battle-slot arrays — `WRAM:+$3BF4` family, 10 entries (party 0–3, enemies 4–9) | Documented; implemented (`battle.BattleSlots`) | **Confirmed** (authority, 10-entry extent, enemy entries 4–5 live via EXP-0003/0004/0005); MP-pair semantics Strong hypothesis; entries 6–9 unobserved | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
