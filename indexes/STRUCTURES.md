# Structures

| ID | Title | Status | Confidence | Record |
|---|---|---|---|---|
| ST-0001 | CharacterFieldsSource — `WRAM:+$2E78`–`+$2EA7` (struct-of-arrays) | Documented; implemented (`chardata.CharacterFieldsSource`) | Confirmed (layout, slot order); field meanings vary by field — see record | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
| ST-0002 | PartySlotRecord — `WRAM:+$2EB5` + n×`$20` (array-of-structs) | Documented; implemented (`chardata.PartySlotRecord`) | Confirmed (stride, six written fields, slots 0–2); bytes `+$0C..$1F` Unknown | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md) |
| ST-0003 | BattleSlotHPArray — `WRAM:+$3BF4` (candidate authoritative layer) | Documented; implemented (`battle.PartySlots`, hypothesis encoding) | Writers Confirmed; authority Strong hypothesis; sibling arrays Unknown | [05_DATA_STRUCTURES.md](../docs/sessions/05_DATA_STRUCTURES.md), [SESSION_003.md](../docs/sessions/SESSION_003.md) |
