# Functions

| ID | Title | Status | Confidence | Record |
|---|---|---|---|---|
| FN-0001 | CopyCharacterFields — `ROMCPU:$C10DF3`–`$C10E66` | Documented; implemented in `chardata` with tests | Confirmed (behavior, full disassembly); purpose of destination buffer strong hypothesis | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0002 | PerFrameBattleUpdate (candidate) — `ROMCPU:$C101FB`–`$C10225` | Documented; callees unexplored | Strong hypothesis (dispatcher role); code Confirmed | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0003 | HP delta routine — store near `ROMCPU:$C21347`/`$C21338`; entry `$C21323` claimed | Documented; implemented + tested (`internal/game/battle`) | Stores Confirmed (raw captures); store addresses Strong hypothesis; disassembly Unknown pending re-dump | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md), [SESSION_003.md](../docs/sessions/SESSION_003.md) |
| FN-0004 | MP delta routine — `ROMCPU:$C21350` (claimed only) | Documented as claim; implemented + tested (`internal/game/battle`) | Unknown — no surviving raw evidence; pending re-dump + MP watch | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md) |
| FN-0005 | Death handler — store near `ROMCPU:$C21396`; entry `$C21390` claimed | Documented; implemented + tested (`internal/game/battle`) | Store Confirmed (raw capture); address Strong hypothesis; remainder Unknown | [02_DISCOVERED_FUNCTIONS.md](../docs/sessions/02_DISCOVERED_FUNCTIONS.md), [SESSION_003.md](../docs/sessions/SESSION_003.md) |
