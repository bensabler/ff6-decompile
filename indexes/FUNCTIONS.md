# Functions

| ID | Title | Status | Confidence | Record |
|---|---|---|---|---|
| FN-0001 | CopyCharacterFields — `ROMCPU:$C10DF3`–`$C10E66` | Documented; implemented in `chardata` with tests | Confirmed (behavior, full disassembly); purpose of destination buffer strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0002 | PerFrameBattleUpdate (candidate) — `ROMCPU:$C101FB`–`$C10225` | Documented; callees unexplored | Strong hypothesis (dispatcher role); code Confirmed | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0003 | HP delta routine — `ROMCPU:$C21323` (claimed) | **Undocumented** — implemented in `internal/game/battle/battle.go`, no canonical record (Session 003 debt) | Unverifiable until documented | pending SESSION_003 |
| FN-0004 | MP delta routine — `ROMCPU:$C21350` (claimed) | **Undocumented** — implemented in `internal/game/battle/battle.go`, no canonical record (Session 003 debt) | Unverifiable until documented | pending SESSION_003 |
| FN-0005 | Death handler — `ROMCPU:$C21390` (claimed) | **Undocumented** — implemented in `internal/game/battle/battle.go`, no canonical record (Session 003 debt) | Unverifiable until documented | pending SESSION_003 |
